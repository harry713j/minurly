package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ShortUrl struct {
	ID          bson.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	OriginalUrl string        `json:"originalUrl" bson:"originalUrl"`
	ShortCode   string        `json:"shortCode" bson:"shortCode"`
	Visits      int           `json:"visits" bson:"visits"`
	UserId      bson.ObjectID `json:"userId" bson:"userId"`
	LastVisited time.Time     `json:"lastVisited" bson:"lastVisited"`
	CreatedAt   time.Time     `json:"createdAt" bson:"createdAt"`
}

// DTO
type CreateUrlPayload struct {
	OriginalUrl string `json:"originalUrl" validate:"required,url"`
}
