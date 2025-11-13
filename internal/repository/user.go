package repository

import (
	"context"
	"time"

	"github.com/harry713j/minurly/internal/models"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserRepo struct {
	collection *mongo.Collection
	log        zerolog.Logger
}

func NewUserRepo(db *mongo.Database, logger zerolog.Logger) *UserRepo {
	return &UserRepo{
		collection: db.Collection("users"),
		log:        logger,
	}
}

func (u *UserRepo) InsertOne(ctx context.Context, user *models.User) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	user.CreatedAt = time.Now()

	res, err := u.collection.InsertOne(ctx, user)
	if err != nil {
		u.log.Err(err).Msg("failed to create an user entry with email " + user.Email)
		return err
	}

	insertedIdStr := res.InsertedID.(bson.ObjectID)

	u.log.Info().Msg("user entry in database created successfully with email " + user.Email + " and user id " + insertedIdStr.String())
	return nil
}

func (u *UserRepo) FindById(ctx context.Context, userId bson.ObjectID) (*models.UserResponse, error) {
	var resp models.UserResponse

	// fetch all the short urls docs, using aggregation pipeline
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "_id", Value: userId}}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "shorturls"},
			{Key: "localField", Value: "_id"},
			{Key: "foreignField", Value: "userId"},
			{Key: "as", Value: "shortUrls"},
		}}},
	}

	cursor, err := u.collection.Aggregate(ctx, pipeline)
	if err != nil {
		u.log.Err(err).Msg("failed to fetch the user response of user with id " + userId.String())
		return nil, err
	}

	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		if err := cursor.Decode(&resp); err != nil {
			u.log.Err(err).Msg("failed to decode the user response of user with id " + userId.String())
			return nil, err
		}

		u.log.Info().Msg("successfully fetch the user response of user with id " + userId.String())
		return &resp, nil
	}

	u.log.Err(err).Msg("failed to fetch the user response of user with id " + userId.String())
	return nil, mongo.ErrNoDocuments
}

func (u *UserRepo) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	filter := bson.D{{Key: "email", Value: email}}
	err := u.collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		u.log.Err(err).Msg("failed to fetch the user document from database of user with email " + email)
		return nil, err
	}

	u.log.Info().Msg("successfully fetch the user document from database of user with email " + email)
	return &user, nil
}
