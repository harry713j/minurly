package helper

import (
	"context"
	"log"

	"github.com/harry713j/minurly/db"
	"github.com/harry713j/minurly/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const userCollectionName string = "users"

var userCollection *mongo.Collection = db.DB.Collection(userCollectionName)

func InsertOneUser(user models.User) (any, error) {
	insertedResult, err := userCollection.InsertOne(context.Background(), user)

	if err != nil {
		log.Println("failed to create an user")
		return nil, err
	}

	return insertedResult.InsertedID, nil
}

func FindUserById(userId string) (*models.UserResponse, error) {
	var user models.UserResponse
	// fetch all the short urls docs, using aggregation pipeline
	pipeline := mongo.Pipeline{
		// Match specific user
		{{Key: "$match", Value: bson.D{{Key: "_id", Value: userId}}}},
		// Lookup with sub-pipeline
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "shorturls"},
			{Key: "let", Value: bson.D{{Key: "url_ids", Value: "$shorturls"}}}, // pass array to sub-pipeline
			{Key: "pipeline", Value: bson.A{
				bson.D{{Key: "$match", Value: bson.D{
					{Key: "$expr", Value: bson.D{
						{Key: "$in", Value: bson.A{"$_id", "$$url_ids"}},
					}},
				}}},
				bson.D{{Key: "$sort", Value: bson.D{
					{Key: "createdAt", Value: -1}, // sort newest first
				}}},
			}},
			{Key: "as", Value: "shorturls"},
		}}},
	}

	cursor, err := userCollection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		log.Println("failed to fetch the user ", err)
		return nil, err
	}

	defer cursor.Close(context.TODO())

	if cursor.Next(context.TODO()) {
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		return &user, nil
	}

	return nil, mongo.ErrNoDocuments
}

func FindUserByEmail(email string) (*models.User, bool) {
	var user models.User

	filter := bson.D{{Key: "email", Value: email}}

	err := userCollection.FindOne(context.TODO(), filter).Decode(&user)

	if err != nil {
		log.Println("Failed to get user by email")
		return nil, false
	}

	return &user, true
}

func AddShortURLToUser(userId string, shortURLId bson.ObjectID) error {
	filter := bson.D{{Key: "_id", Value: userId}}
	update := bson.D{{Key: "$push", Value: bson.D{{Key: "shorturls", Value: shortURLId}}}}

	_, err := userCollection.UpdateOne(context.TODO(), filter, update)
	return err
}

func RemoveShortURLFromUser(userId string, shortURLId bson.ObjectID) error {
	filter := bson.D{{Key: "_id", Value: userId}}
	update := bson.D{{Key: "$pull", Value: bson.D{{Key: "shorturls", Value: shortURLId}}}}

	_, err := userCollection.UpdateOne(context.TODO(), filter, update)
	return err
}
