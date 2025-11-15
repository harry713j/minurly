package service

import (
	"context"
	"fmt"

	"github.com/harry713j/minurly/internal/apperrors"
	"github.com/harry713j/minurly/internal/config"
	"github.com/harry713j/minurly/internal/models"
	"github.com/harry713j/minurly/internal/repository"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
	"google.golang.org/api/idtoken"
)

type UserService struct {
	cfg  *config.Config
	repo *repository.Repository
	log  zerolog.Logger
}

func NewUserService(cfg *config.Config, repo *repository.Repository, logger zerolog.Logger) *UserService {
	return &UserService{
		cfg:  cfg,
		repo: repo,
		log:  logger,
	}
}

func (u *UserService) Create(authCode string) (*models.User, error) {
	ctx := context.Background()

	token, err := u.cfg.Auth.OAuthConfig.Exchange(ctx, authCode)
	if err != nil {
		u.log.Warn().Str("err", err.Error()).Msg("failed to exchange the oauth token")
		return nil, apperrors.NewInternalServerErr()
	}

	rawIdToken, ok := token.Extra("id_token").(string) // jwt token
	if !ok {
		u.log.Warn().Msg("no id_token found in oauth token")
		return nil, apperrors.NewInternalServerErr()
	}

	payload, err := idtoken.Validate(ctx, rawIdToken, u.cfg.Auth.GoogleClientId)
	if err != nil {
		u.log.Warn().Str("err", err.Error()).Msg("failed to validate from provided id token")
		return nil, apperrors.NewUnauthorizedErr("Invalid Token")
	}

	email := payload.Claims["email"].(string)
	// check user with this email already exists or not
	existedUser, err := u.repo.User.FindByEmail(ctx, email)
	if err == nil {
		u.log.Info().Msg(fmt.Sprintf("user with email %v already exists in database", email))
		return existedUser, nil
	}

	user := &models.User{
		Email:   email,
		Name:    payload.Claims["name"].(string),
		Profile: payload.Claims["picture"].(string),
		OAuthId: payload.Subject,
	}

	id, err := u.repo.User.InsertOne(ctx, user)
	if err != nil {
		u.log.Err(err).Msg("failed to create user in database with email " + email)
		return nil, apperrors.NewInternalServerErr()
	}

	user.ID = *id
	u.log.Info().Msg("user creation in database successful with user email " + email)

	return user, nil
}

func (u *UserService) Fetch(userId string) (*models.UserResponse, error) {
	userObjId, err := bson.ObjectIDFromHex(userId)
	if err != nil {
		u.log.Warn().Str("err", err.Error()).Msg("invalid user id " + userId)
		return nil, apperrors.NewUnauthorizedErr("Unauthorized")
	}
	ctx := context.Background()

	resp, err := u.repo.User.FindById(ctx, userObjId)
	if err != nil {
		u.log.Err(err).Msg("failed to get user url data of user with userId " + userId)
		return nil, apperrors.NewInternalServerErr()
	}

	u.log.Info().Msg("successfully get the user url data of user with userId " + userId)
	return resp, nil
}
