package main

import (
	"log"

	"github.com/carfdev/microservice-go/internal/adapter/db"
	"github.com/carfdev/microservice-go/internal/adapter/nats"
	"github.com/carfdev/microservice-go/internal/application"
	"github.com/carfdev/microservice-go/internal/config"
)

func main() {
	// Load config
	cfg := config.LoadConfig()

	// Connect to PostgreSQL
	gormDB, err := db.Connect(cfg.PostgresDSN)
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
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
