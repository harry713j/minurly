package handler

import (
	"net/http"

	"github.com/harry713j/minurly/internal/config"
	"github.com/harry713j/minurly/internal/helper"
	"github.com/harry713j/minurly/internal/models"
	"github.com/harry713j/minurly/internal/utils"
)

func HandleGetUser(w http.ResponseWriter, r *http.Request) {
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

	user, err := helper.FindUserById(userId)

	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	type responseUser struct {
		Message string              `json:"message"`
		User    models.UserResponse `json:"user"`
	}

	utils.RespondWithJSON(w, http.StatusOK, responseUser{
		Message: "Fetch user details successfully",
		User:    *user,
	})
}
