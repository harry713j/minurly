package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ShortUrl struct {
	ID bson.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"` 
	// The omitempty struct tag omits the corresponding field from the inserted document when left empty.
	Original string `json:"original,omitempty"`
	Short string `json:"short,omitempty"`
	VisitCount int `json:"visitCount"`
}