package storage

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	us "github.com/g3orge/BasicAuthRest/internal/user"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type Storage struct {
	db *mongo.Collection
}

func New(url, databaseName, collectionName string) (*Storage, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(url))
	if err != nil {
		return nil, fmt.Errorf("cannot connect to database: %v", err)
	}

	if err = client.Ping(context.Background(), nil); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return &Storage{db: client.Database(databaseName).Collection(collectionName)}, nil
}

func (s *Storage) Create(ctx context.Context, user us.User) (string, error) {
	filter := bson.M{"email": user.Email}

	r := s.db.FindOne(ctx, filter)
	if r.Err() == nil {
		return "", fmt.Errorf("user already exsist")
	}

	res, err := s.db.InsertOne(ctx, user)
	if err != nil {
		return "", fmt.Errorf("failed when creating a user: %v", err)
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if ok {
		return oid.Hex(), nil
	}

	return "", fmt.Errorf("failed to convert oid to hex")
}

func (s *Storage) GenerateTokens(ctx context.Context, guid string) (string, string, error) {
	oid, err := primitive.ObjectIDFromHex(guid)
	if err != nil {
		return "", "", fmt.Errorf("failed to convers id to objectid: %v", err)
	}

	filter := bson.M{"_id": oid}

	r := s.db.FindOne(ctx, filter)
	if r.Err() != nil {
		return "", "", fmt.Errorf("failed to find user")
	}

	var user us.User

	if err := r.Decode(&user); err != nil {
		return "", "", fmt.Errorf("failed when decoding user")
	}

	secretKey := []byte(user.Email + user.Password)

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 1).Unix(),
		"sub": user.ID,
	})

	at, err := token.SignedString(secretKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to create jwt token: %v", err)
	}

	rand := uuid.New()
	rtBase64 := base64.StdEncoding.EncodeToString([]byte(rand.String()))

	hashedRT, err1 := createTokenHash(rtBase64)
	if err1 != nil {
		return "", "", fmt.Errorf("error in creating token hash: %v", err1)
	}

	fmt.Println(user.ID)
	session := us.Session{
		RefreshToken: hashedRT,
		ExpAt:        time.Now().Add(time.Hour * 2),
	}

	_, err2 := s.db.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": bson.M{"session": session}})
	if err2 != nil {
		return "", "", fmt.Errorf("cannot insert refresh token: %v", err2)
	}

	return at, rtBase64, nil
}

func (s *Storage) RefreshToken(ctx context.Context, rt string, guid string) (string, string, error) {
	oid, err := primitive.ObjectIDFromHex(guid)
	if err != nil {
		return "", "", fmt.Errorf("failed to convers id to objectid: %v", err)
	}

	filter := bson.M{"_id": oid}

	r := s.db.FindOne(ctx, filter)
	if r.Err() != nil {
		return "", "", fmt.Errorf("failed to find user")
	}

	var user us.User

	if err := r.Decode(&user); err != nil {
		return "", "", fmt.Errorf("failed when decoding user")
	}

	if err = checkTokenHash(rt, user.Session.RefreshToken); err != nil {
		return "", "", fmt.Errorf("refreshedToken does not match: %v", err)
	}

	rand := uuid.New()
	rtBase64 := base64.StdEncoding.EncodeToString([]byte(rand.String()))

	hashedRT, err1 := createTokenHash(rtBase64)
	if err1 != nil {
		return "", "", fmt.Errorf("error in creating token hash: %v", err1)
	}

	session := us.Session{
		RefreshToken: hashedRT,
		ExpAt:        time.Now().Add(time.Hour * 2),
	}

	_, err2 := s.db.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": bson.M{"session": session}})

	if err2 != nil {
		return "", "", fmt.Errorf("cannot insert refresh token: %v", err2)
	}

	secretKey := []byte(user.Email + user.Password)

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 1).Unix(),
		"sub": user.ID,
	})

	at, err3 := token.SignedString(secretKey)
	if err3 != nil {
		return "", "", fmt.Errorf("failed to create jwt token: %v", err3)
	}

	return at, rtBase64, nil
}

func checkTokenHash(refreshToken, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(refreshToken))

	return err
}

func createTokenHash(rt string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(rt), 12)

	return string(bytes), err
}
