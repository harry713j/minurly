package config

import (
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	OAuthConfig  *oauth2.Config
	SessionStore = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
)

func init() {
	godotenv.Load()

	OAuthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"email", "profile"},
		RedirectURL:  os.Getenv("SERVER_URL") + "/auth/google/callback",
		Endpoint:     google.Endpoint,
	}

	SessionStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		Secure:   false,                // for dev
		SameSite: http.SameSiteLaxMode, // for dev only
	}
}
