package logging

import (
	"io"
	"os"

	"github.com/harry713j/minurly/internal/config"
	"github.com/rs/zerolog"
)

type LoggerService struct {
	Logger zerolog.Logger
}

func NewLogger(cfg *config.Config) *LoggerService {
	logLvlStr := cfg.Log.Level
	logLvl, err := zerolog.ParseLevel(logLvlStr)

	if err != nil {
		if cfg.Primary.Env == "production" {
			logLvl = zerolog.InfoLevel
		} else {
			logLvl = zerolog.DebugLevel
		}
	}

	var writer io.Writer
	if cfg.Primary.Env == "production" {
		writer = os.Stdout
	} else {
		consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006-01-02 15:04:05"}
		writer = consoleWriter
	}

	logger := zerolog.New(writer).Level(logLvl).With().Timestamp().Str("environment", cfg.Primary.Env).Logger()

	if cfg.Primary.Env != "production" {
		logger = logger.With().Stack().Logger()
	}

	return &LoggerService{Logger: logger}
}
