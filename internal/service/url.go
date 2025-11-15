package service

import (
	"context"
	"errors"
	"time"

	"github.com/harry713j/minurly/internal/apperrors"
	"github.com/harry713j/minurly/internal/models"
	"github.com/harry713j/minurly/internal/repository"
	"github.com/harry713j/minurly/internal/utils"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UrlService struct {
	repo *repository.Repository
	log  zerolog.Logger
}

func NewUrlService(repo *repository.Repository, logger zerolog.Logger) *UrlService {
	return &UrlService{
		repo: repo,
		log:  logger,
	}
}

func (u *UrlService) Create(userId string, originalUrl string) (*models.ShortUrl, error) {
	userObjId, err := bson.ObjectIDFromHex(userId)
	if err != nil {
		u.log.Warn().Str("err", err.Error()).Msg("invalid user id " + userId)
		return nil, apperrors.NewUnauthorizedErr("Unauthorized")
	}

	ctx := context.Background()

	if _, err := u.repo.User.FindById(ctx, userObjId); err != nil {
		u.log.Warn().Str("err", err.Error()).Msg("user not exists with userId " + userId)
		return nil, apperrors.NewUnauthorizedErr("user not authorized to perform this operation")
	}

	code, err := utils.GenerateRandomStrings(8)
	if err != nil {
		u.log.Err(err).Msg("failed to generate the code for creating short url ")
		return nil, apperrors.NewInternalServerErr()
	}

	url := &models.ShortUrl{
		OriginalUrl: originalUrl,
		ShortCode:   code,
		Visits:      0,
		UserId:      userObjId,
		LastVisited: time.Now(),
		CreatedAt:   time.Now(),
	}

	err = u.repo.Url.InsertOne(ctx, url)
	if err != nil {
		u.log.Err(err).Msg("failed to create the short url in database for user with userId " + userId)
		return nil, apperrors.NewInternalServerErr()
	}

	u.log.Info().Msg("short url created successfully for the user with userId " + userId)
	return url, nil
}

func (u *UrlService) Fetch(shortCode string) (*models.ShortUrl, error) {
	ctx := context.Background()

	url, err := u.repo.Url.FindByCode(ctx, shortCode)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			u.log.Warn().Str("err", err.Error()).Msg("no short url entry exists with short code " + shortCode)
			return nil, apperrors.NewBadRequestErr("invalid code provided")
		}

		u.log.Err(err).Msg("failed to fetch the short url with short code " + shortCode)
		return nil, apperrors.NewInternalServerErr()
	}

	u.log.Info().Msg("successfully fetch the short url with short code " + shortCode)
	return url, nil
}

func (u *UrlService) Remove(userId string, shortCode string) error {
	userObjId, err := bson.ObjectIDFromHex(userId)
	if err != nil {
		u.log.Warn().Str("err", err.Error()).Str("userId", userId).Msg("invalid user id ")
		return apperrors.NewUnauthorizedErr("Unauthorized")
	}

	ctx := context.Background()
	if _, err := u.repo.User.FindById(ctx, userObjId); err != nil {
		u.log.Warn().Str("err", err.Error()).Msg("user not exists with userId " + userId)
		return apperrors.NewUnauthorizedErr("user not authorized to perform this operation")
	}

	err = u.repo.Url.DeleteOne(ctx, shortCode, userObjId)
	if err != nil {
		u.log.Err(err).Msg("failed to remove the short url with code " + shortCode + " of user with userId " + userId)
		return nil
	}

	u.log.Info().Msg("short url deleted successfully with code " + shortCode + " of user with userId " + userId)
	return nil
}
