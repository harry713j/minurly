package db

import (
	"log"
	"os"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/harry713j/minurly/internal/utils"
	"github.com/joho/godotenv"
)

// exporting for global use
var (
	Client *mongo.Client
	DB     *mongo.Database
)

func connectDB() {
	DBURL := os.Getenv("DBURL")

	if DBURL == "" {
		log.Fatal("Database URI not found!!")
		return
	}

	Client, err := mongo.Connect(options.Client().ApplyURI(DBURL))

	if err != nil {
		log.Fatal("Failed to connect to the Database")
	}

	DB = Client.Database(utils.DBNAME)
	log.Println("Connect to database successfully")
}

func init() {
	godotenv.Load()
	connectDB()
}
