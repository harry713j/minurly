package helper

import (
	"context"
	"log"

	"github.com/harry713j/minurly/db"
	"github.com/harry713j/minurly/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const sessionCollectionName = "sessions"

var sessionCollection *mongo.Collection = db.DB.Collection(sessionCollectionName)

func InsertOneSession(session models.Session) (any, error) {
	inserted, err := sessionCollection.InsertOne(context.Background(), session)

	if err != nil {
		log.Println("Failed to insert session to db: ", err)
		return nil, err
	}

	return inserted.InsertedID, nil
}

func DeleteSession(userId string) (int, error) {
	filter := bson.D{{Key: "userId", Value: userId}}

	result, err := sessionCollection.DeleteOne(context.TODO(), filter)

	if err != nil {
		log.Println("failed to delete the session")
		return 0, err
	}

	return int(result.DeletedCount), err
}
