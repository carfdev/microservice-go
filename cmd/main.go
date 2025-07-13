package main

import (
	"log"

	"github.com/carfdev/microservice-go/internal/adapter/db"
	"github.com/carfdev/microservice-go/internal/adapter/nats"
	"github.com/carfdev/microservice-go/internal/application"
	"github.com/carfdev/microservice-go/internal/config"
	"gorm.io/gorm"
)

func main() {
	// Load config
	cfg := config.LoadConfig()

	// Connect to PostgreSQL
	var gormDB *gorm.DB
	var err error

	if cfg.AppEnv == "development" {
		log.Println("Using development DB connection with migration")
		gormDB, err = db.ConnectAndMigrate(cfg.PostgresDSN)
	} else {
		log.Println("Using production DB connection")
		gormDB, err = config.ConnectDB(cfg.PostgresDSN)
	}

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Connect to NATS
	nc, err := config.ConnectNATS(cfg.NatsURL)
	if err != nil {
		log.Fatalf("NATS connection failed: %v", err)
	}
	defer nc.Drain() // Will run only if process is terminated

	// Set up the application
	invoiceRepo := db.NewInvoiceRepository(gormDB)
	invoiceService := application.NewInvoiceService(invoiceRepo)
	natsHandler := nats.NewInvoiceNATSAdapter(nc, invoiceService)

	// Start listening to NATS subjects
	natsHandler.ListenForMessages()

	// Keep the service running forever
	select {}
}
