package controller

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/harry713j/minurly/config"
	"github.com/harry713j/minurly/helper"
	"github.com/harry713j/minurly/models"
	"github.com/harry713j/minurly/utils"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/oauth2"
	"google.golang.org/api/idtoken"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	state, err := utils.GenerateRandomStrings(16)

	if err != nil {
		http.Error(w, "Failed to login", http.StatusInternalServerError)
		return
	}

	// url will be created by oauth client
	// look for this state
	url := config.OAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusFound)
}

func HandleLoginCallback(w http.ResponseWriter, r *http.Request) {
	godotenv.Load()
	code := r.URL.Query().Get("code") // oauth client send a ?code=random in query
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	token, err := config.OAuthConfig.Exchange(context.Background(), code)

	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	rawIdToken, ok := token.Extra("id_token").(string)

	if !ok {
		http.Error(w, "No id_token found", http.StatusInternalServerError)
		return
	}

	payload, err := idtoken.Validate(context.Background(), rawIdToken, os.Getenv("GOOGLE_CLIENT_ID"))

	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	email := payload.Claims["email"].(string)

	user, _ := helper.FindUserByEmail(email)

	if user == nil {
		// user not exists
		user = &models.User{
			Email:     email,
			Name:      payload.Claims["name"].(string),
			Profile:   payload.Claims["picture"].(string),
			OAuthId:   payload.Subject,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			ShortUrls: []bson.ObjectID{},
		}

		id, err := helper.InsertOneUser(*user)

		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "failed to create the user")
			return
		}

		user.ID = id.(bson.ObjectID) // type assertion
	}

	// create the session
	sessionId := uuid.NewString()

	sessionDoc := &models.Session{
		SessionId: sessionId,
		UserId:    user.ID,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	res, err := helper.InsertOneSession(*sessionDoc)

	if err != nil || res == nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to create session")
		return
	}

	session, _ := config.SessionStore.Get(r, "sessionId")

	session.Values["sessionId"] = sessionId
	session.Values["userId"] = user.ID.Hex()
	session.Values["email"] = email

	if err := session.Save(r, w); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to save session")
		return
	}

	http.Redirect(w, r, os.Getenv("ALLOWED_ORIGIN")+"/dashboard", http.StatusSeeOther)
}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	// delete the session and clear the cookies
	session, err := config.SessionStore.Get(r, "sessionId")

	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userId, ok := session.Values["userId"].(string)

	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid or expired session")
		return
	}

	// check user exists or not
	user, err := helper.FindUserById(userId)

	if err != nil || user == nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid session")
		return
	}

	// remove the session entry from database and remove all the cookies
	row, err := helper.DeleteSession(userId)

	if err != nil || row == 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "Failed to remove the session")
		return
	}

	// Clear session values in memory
	session.Options.MaxAge = -1 // Marks it as expired
	err = session.Save(r, w)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to clear session")
		return
	}

	// Clear cookie explicitly (optional, but safe)
	http.SetCookie(w, &http.Cookie{
		Name:     "sessionId",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Log out successful",
	})
}
