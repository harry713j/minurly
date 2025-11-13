package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/harry713j/minurly/internal/models"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UrlRepo struct {
	collection *mongo.Collection
	log        zerolog.Logger
}

func NewUrlRepo(db *mongo.Database, logger zerolog.Logger) *UrlRepo {
	return &UrlRepo{
		collection: db.Collection("shorturls"),
		log:        logger,
	}
}

func (u *UrlRepo) InsertOne(ctx context.Context, url *models.ShortUrl) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	url.CreatedAt = time.Now()
	res, err := u.collection.InsertOne(ctx, url)

	if err != nil {
		u.log.Err(err).Msg(fmt.Sprintf("failed to create url entry in database with original url %s, short code %s and userId %s",
			url.OriginalUrl, url.ShortCode, url.UserId))
		return err
	}

	insertedIdStr := res.InsertedID.(bson.ObjectID)

	u.log.Info().Msg(fmt.Sprintf("successfully created url entry in database with id %s of original url %s and userId %s",
		insertedIdStr.String(), url.OriginalUrl, url.UserId))
	return nil
}

func (u *UrlRepo) FindByCode(ctx context.Context, shortCode string) (*models.ShortUrl, error) {
	var result models.ShortUrl

	filter := bson.D{{Key: "shortCode", Value: shortCode}}
	update := bson.D{{Key: "$inc", Value: bson.D{{Key: "visits", Value: 1}}}}

	err := u.collection.FindOneAndUpdate(ctx, filter, update).Decode(&result)

	if err != nil {
		u.log.Err(err).Msg(fmt.Sprintf("failed to fetch url documents from database with code %s", shortCode))
		return nil, err
	}

	u.log.Info().Msg(fmt.Sprintf("successfully fetch url documents from database with code %s", shortCode))
	return &result, nil
}

func (u *UrlRepo) DeleteOne(ctx context.Context, shortCode string, userId bson.ObjectID) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.D{{Key: "shortCode", Value: shortCode}}

	if _, err := u.collection.DeleteOne(ctx, filter); err != nil {
		u.log.Err(err).Msg(fmt.Sprintf("failed to remove url entry from database with short code %s and userId %s", shortCode, userId.String()))
		return err
	}

	u.log.Info().Msg(fmt.Sprintf("successfully remove url entry from database with short code %s and userId %s", shortCode, userId.String()))
	return nil
}
