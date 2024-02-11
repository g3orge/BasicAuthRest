package us

import "time"

type User struct {
	ID       string  `json:"id,omitempty" bson:"_id,omitempty"`
	Email    string  `json:"email" bson:"email"`
	Password string  `json:"password" bson:"password"`
	Session  Session `bson:"session,omitempty"`
}

type Session struct {
	RefreshToken string    `json:"refreshtoken" bson:"refreshToken"`
	ExpAt        time.Time `json:"expAt" bson:"expiresAt"`
}
