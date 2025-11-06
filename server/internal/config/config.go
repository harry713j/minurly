package config

import (
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/sessions"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Config struct {
	Primary  *PrimaryConfig
	Server   *ServerConfig
	Database *DatabaseConfig
	Auth     *AuthConfig
}

type PrimaryConfig struct {
	Env string
}

type ServerConfig struct {
	Port               string
	ReadTimeOut        int
	WriteTimeOut       int
	IdleTimeOut        int
	CORSAllowedOrigins []string
}

type AuthConfig struct {
	OAuthConfig  *oauth2.Config
	SessionStore *sessions.CookieStore
}

type DatabaseConfig struct {
	DbURL string
}

func LoadConfig() *Config {
	primary := &PrimaryConfig{Env: "development"}

	port := os.Getenv("PORT")
	whiteLists := strings.Split(os.Getenv("ALLOWED_ORIGIN"), ",")
	server := &ServerConfig{
		Port:               port,
		ReadTimeOut:        10,
		WriteTimeOut:       15,
		IdleTimeOut:        60,
		CORSAllowedOrigins: whiteLists,
	}

	oAuthConfig := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"email", "profile"},
		RedirectURL:  os.Getenv("SERVER_URL") + "/auth/google/callback",
		Endpoint:     google.Endpoint,
	}

	sessionStore := sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

	auth := &AuthConfig{
		OAuthConfig:  oAuthConfig,
		SessionStore: sessionStore,
	}

	auth.SessionStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		Secure:   false,                // for dev
		SameSite: http.SameSiteLaxMode, // for dev only
	}

	db := &DatabaseConfig{
		DbURL: os.Getenv("DATABASE_URL"),
	}

	mainConfig := &Config{
		Primary:  primary,
		Server:   server,
		Auth:     auth,
		Database: db,
	}

	return mainConfig
}
