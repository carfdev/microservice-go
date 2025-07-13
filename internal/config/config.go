package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Config structure to hold our application configurations
type Config struct {
	NatsURL        string
	PostgresDSN    string
	AppEnv         string
}

// LoadConfig loads environment variables and returns a Config struct
func LoadConfig() (*Config) {


	// Get variables from environment
	natsURL := os.Getenv("NATS_URL")
	postgresDSN := os.Getenv("DATABASE_URL")


	// Try to load from .env only if one or both are missing
	if natsURL == "" || postgresDSN == "" {
		_ = godotenv.Load() // ignore error, might not exist in Docker

		// Re-fetch after loading .env
		if natsURL == "" {
			natsURL = os.Getenv("NATS_URL")
		}
		if postgresDSN == "" {
			postgresDSN = os.Getenv("DATABASE_URL")
		}
	}

	appEnv := os.Getenv("APP_ENV")

	if natsURL == "" {
		log.Fatal("NATS_URL is not set in environment")
		return nil
	}

	
	if postgresDSN == "" {
		log.Fatal("DATABASE_URL is not set in environment")
		return nil
	}

	if appEnv == "" {
		appEnv = "development" // default fallback
	}


	return &Config{
		NatsURL:     natsURL,
		PostgresDSN: postgresDSN,
		AppEnv:      appEnv,
	}
}

// ConnectNATS connects to the NATS server and returns the connection
func ConnectNATS(natsURL string) (*nats.Conn, error) {
	// Connect to NATS
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
		return nil, err
	}
	return nc, nil
}

// ConnectDB connects to the PostgreSQL database and returns the connection
func ConnectDB(postgresDSN string) (*gorm.DB, error) {
	// Connect to PostgreSQL database using GORM
	db, err := gorm.Open(postgres.Open(postgresDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		return nil, err
	}
	return db, nil
}
