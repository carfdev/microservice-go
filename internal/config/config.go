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
}

// LoadConfig loads environment variables and returns a Config struct
func LoadConfig() (*Config) {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return nil
	}

	// Get variables from environment
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		log.Fatal("NATS_URL is not set in environment")
		return nil
	}

	postgresDSN := os.Getenv("DATABASE_URL")
	if postgresDSN == "" {
		log.Fatal("DATABASE_URL is not set in environment")
		return nil
	}

	return &Config{
		NatsURL:     natsURL,
		PostgresDSN: postgresDSN,
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
