package helper

import (
	"context"
	"log"

	"github.com/harry713j/minurly/db"
	"github.com/harry713j/minurly/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// define mongodb operations
const collectionName string = "shorturls"
var collection *mongo.Collection = db.DB.Collection(collectionName)


func InsertOneUrl(url models.ShortUrl) (any, error){
   insertedResult, err := collection.InsertOne(context.Background(), url)

   if err != nil {
	log.Println("Failed to insert the url")
	return nil, err
   }

   return insertedResult.InsertedID, nil
}

func FindOneUrlByShort(short string) (*models.ShortUrl, error){
	var result models.ShortUrl

	filter := bson.M{"short": short}
	update := bson.D{{Key: "$inc", Value: bson.D{{Key: "visitcount", Value: 1}}}}

	err := collection.FindOneAndUpdate(context.Background(), filter, update).Decode(&result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}