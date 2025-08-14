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

	userObjId, err := bson.ObjectIDFromHex(userId)

	if err != nil {
		log.Println("Failed to convert to object id")
		return nil, err
	}

	// fetch all the short urls docs, using aggregation pipeline
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "_id", Value: userObjId}}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "shorturls"},
			{Key: "localField", Value: "_id"},
			{Key: "foreignField", Value: "userId"},
			{Key: "as", Value: "shortUrls"},
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
		return nil, false
	}

	return &user, true
}
