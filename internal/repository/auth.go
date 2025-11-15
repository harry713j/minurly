package repository

import (
	"context"
	"time"

	"github.com/harry713j/minurly/internal/models"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AuthRepo struct {
	collection *mongo.Collection
	log        zerolog.Logger
}

func NewAuthRepo(db *mongo.Database, logger zerolog.Logger) *AuthRepo {
	return &AuthRepo{
		collection: db.Collection("sessions"),
		log:        logger,
	}
}

func (a *AuthRepo) InsertOne(ctx context.Context, session *models.Session) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := a.collection.InsertOne(ctx, session)
	if err != nil {
		a.log.Err(err).Msg("failed to create session on database with session string " + session.SessionId + " and user id " + session.UserId.String())
		return err
	}

	insertedIdStr := res.InsertedID.(bson.ObjectID)

	a.log.Info().Msg("session entry in database created with id " + insertedIdStr.String())
	return nil
}

func (a *AuthRepo) DeleteOne(ctx context.Context, userId bson.ObjectID) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := a.collection.DeleteOne(ctx, bson.D{{Key: "userId", Value: userId}})
	if err != nil {
		a.log.Err(err).Msg("failed to remove the session from database of user with userId " + userId.String())
		return err
	}

	a.log.Info().Msg("session entry deleted from database of user with userId " + userId.String())
	return nil
}
