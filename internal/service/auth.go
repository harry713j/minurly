package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/harry713j/minurly/internal/apperrors"
	"github.com/harry713j/minurly/internal/models"
	"github.com/harry713j/minurly/internal/repository"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AuthService struct {
	repo *repository.Repository
	log  zerolog.Logger
}

func NewAuthService(repo *repository.Repository, logger zerolog.Logger) *AuthService {
	return &AuthService{
		repo: repo,
		log:  logger,
	}
}

func (a *AuthService) CreateSession(userId bson.ObjectID) (*models.Session, error) {
	ctx := context.Background()

	_, err := a.repo.User.FindById(ctx, userId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			a.log.Warn().Str("userId", userId.Hex()).Msg("user not found")
			return nil, apperrors.NewNotFoundErr("user not found")
		}
		a.log.Err(err).Msg("database error while finding user")
		return nil, apperrors.NewInternalServerErr()
	}

	session := &models.Session{
		SessionId: uuid.NewString(),
		UserId:    userId,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	if err := a.repo.Auth.InsertOne(ctx, session); err != nil {
		a.log.Err(err).Msg("failed to insert session into database of user with userId " + userId.String())
		return nil, apperrors.NewInternalServerErr()
	}

	return session, nil
}

func (a *AuthService) DeleteSession(userId bson.ObjectID) error {
	ctx := context.Background()

	_, err := a.repo.User.FindById(ctx, userId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			a.log.Warn().Str("userId", userId.Hex()).Msg("user not found")
			return apperrors.NewNotFoundErr("user not found")
		}
		a.log.Err(err).Msg("database error while finding user")
		return apperrors.NewInternalServerErr()
	}

	if err := a.repo.Auth.DeleteOne(ctx, userId); err != nil {
		a.log.Err(err).Msg("failed to delete session from database of user with userId " + userId.String())
		return apperrors.NewInternalServerErr()
	}

	return nil
}
