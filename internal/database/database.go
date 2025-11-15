package database

import (
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/harry713j/minurly/internal/config"
	"github.com/harry713j/minurly/internal/logging"
)

func New(cfg *config.Config, logger *logging.LoggerService) (*mongo.Client, *mongo.Database) {
	client, err := mongo.Connect(options.Client().ApplyURI(cfg.Database.DbURL))
	if err != nil {
		logger.Logger.Fatal().Err(err).Msg("MongoDB connection failed: ")
	}

	db := client.Database(cfg.Database.DbName)
	logger.Logger.Info().Str("Database", "Connected to database successfully")
	return client, db
}
