package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ShortUrl struct {
	ID          bson.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	OriginalUrl string        `json:"originalUrl"`
	ShortCode   string        `json:"shortCode"`
	Visits      int           `json:"visits"`
	UserId      bson.ObjectID `json:"userId"`
	LastVisited time.Time     `json:"lastVisited"`
	CreatedAt   time.Time     `json:"createdAt"`
}
