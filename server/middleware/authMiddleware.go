package middleware

import (
	"net/http"

	"github.com/harry713j/minurly/config"
	"github.com/harry713j/minurly/helper"
	"github.com/harry713j/minurly/utils"
)

func VerifyLogin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get session from the store
		session, err := config.SessionStore.Get(r, "sessionId")
		if err != nil {
			http.Error(w, "Unable to get session", http.StatusInternalServerError)
			return
		}

		userId, ok := session.Values["userId"].(string)

		if !ok || userId == "" {
			utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		user, err := helper.FindUserById(userId)

		if err != nil || user == nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid or expired session")
			return
		}

		next(w, r)
	}
}
