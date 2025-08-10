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
const urlCollectionName string = "shorturls"

var urlCollection *mongo.Collection = db.DB.Collection(urlCollectionName)

func InsertOneUrl(url models.ShortUrl) (any, error) {
	insertedResult, err := urlCollection.InsertOne(context.Background(), url)

	if err != nil {
		log.Println("Failed to insert the url")
		return nil, err
	}

	err = AddShortURLToUser(url.UserId.String(), insertedResult.InsertedID.(bson.ObjectID))

	if err != nil {
		log.Println("Error while create url ", err)
		return nil, err
	}

	return insertedResult.InsertedID, nil
}

func FindOneUrlByShort(shortCode string) (*models.ShortUrl, error) {
	var result models.ShortUrl

	filter := bson.D{{Key: "shortCode", Value: shortCode}}
	update := bson.D{{Key: "$inc", Value: bson.D{{Key: "visits", Value: 1}}}}

	err := urlCollection.FindOneAndUpdate(context.Background(), filter, update).Decode(&result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func DeleteUrlByShort(shortCode string, userId bson.ObjectID) (int, error) {
	filter := bson.D{{Key: "shortCode", Value: shortCode}}

	shortUrl, err := FindOneUrlByShort(shortCode)

	if err != nil {
		log.Println("Unable to fetch short url")
		return 0, err
	}

	result, err := urlCollection.DeleteOne(context.TODO(), filter)

	if err != nil {
		return 0, err
	}

	err = RemoveShortURLFromUser(userId.String(), shortUrl.ID)

	if err != nil {
		log.Println("Error while remove url ", err)
		return 0, err
	}

	return int(result.DeletedCount), nil
}
