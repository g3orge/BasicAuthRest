package storage

import (
	"context"
	"fmt"
	"os/user"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storage struct {
	db *mongo.Collection
}

func New() (*Storage, error) {
	url := "mongodb://localhost:27017"
	databaseName := "bAuth"
	collectionName := "simpleAuth"

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(url))
	if err != nil {
		return nil, fmt.Errorf("cannot connect to database: %v", err)
	}

	if err = client.Ping(context.Background(), nil); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return &Storage{db: client.Database(databaseName).Collection(collectionName)}, nil
}

func (s *Storage) Create(ctx context.Context, user user.User) (string, error) {
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

func (s *Storage) FindOne(ctx context.Context, id string) (user user.User, err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return user, fmt.Errorf("failed to convers id to objectid: %v", err)
	}

	filter := bson.M{"_id": oid}
	result := s.db.FindOne(ctx, filter)
	if result.Err() != nil {
		return user, fmt.Errorf("failed to find one user by id: %s", id)
	}

	if err = result.Decode(&user); err != nil {
		return user, fmt.Errorf("failed to decode user: %v", err)
	}

	return user, nil
}
