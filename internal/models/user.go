package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID        bson.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	OAuthId   string        `json:"oauthId" bson:"oauthId"`
	Email     string        `json:"email" bson:"email"`
	Name      string        `json:"name" bson:"name"`
	Profile   string        `json:"profile" bson:"profile"`
	CreatedAt time.Time     `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt time.Time     `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
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
	SessionId string        `json:"sessionId" bson:"sessionId"`
	UserId    bson.ObjectID `json:"userId" bson:"userId"`
	CreatedAt time.Time     `json:"createdAt" bson:"createdAt"`
	ExpiresAt time.Time     `json:"expiresAt" bson:"expiresAt"`
}
