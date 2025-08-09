package models

import (
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID        bson.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	OAuthId   string          `json:"oauthId"` // google's unique id provided by google
	Email     string          `json:"email"`
	Name      string          `json:"name"`
	Profile   string          `json:"profile"`
	CreatedAt time.Time       `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt time.Time       `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
	ShortUrls []bson.ObjectID `json:"shortUrls"`
}

type UserResponse struct {
	ID        bson.ObjectID `bson:"_id" json:"id"`
	Name      string        `bson:"name" json:"name"`
	Email     string        `bson:"email" json:"email"`
	Profile   string        `bson:"profile" json:"profile"`
	ShortURLs []ShortUrl    `bson:"shorturls" json:"shorturls"`
}

type Session struct {
	ID        bson.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	SessionId uuid.UUID     `json:"sessionId"`
	UserId    bson.ObjectID `json:"userId"`
	CreatedAt time.Time     `json:"createdAt"`
	ExpiresAt time.Time     `json:"expiresAt"`
}
