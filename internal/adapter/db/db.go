package db

import (
	"fmt"
	"log"

	"github.com/carfdev/microservice-go/internal/domain"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Invoice repository implementation using GORM
type InvoiceRepository struct {
	db *gorm.DB
}



func NewInvoiceRepository(db *gorm.DB) *InvoiceRepository {
	return &InvoiceRepository{db}
}

func (r *InvoiceRepository) Create(invoice *domain.Invoice) (*domain.Invoice, error) {
	if err := r.db.Create(invoice).Error; err != nil {
		return nil, fmt.Errorf("error creating invoice: %w", err)
	}
	return invoice, nil
}

func (r *InvoiceRepository) GetAll() ([]domain.Invoice, error) {
	var invoices []domain.Invoice
	if err := r.db.Find(&invoices).Error; err != nil {
		return nil, fmt.Errorf("error fetching invoices: %w", err)
	}
	return invoices, nil
}

func (r *InvoiceRepository) GetByID(id uuid.UUID) (*domain.Invoice, error) {
	var invoice domain.Invoice
	if err := r.db.First(&invoice, id).Error; err != nil {
		return nil, fmt.Errorf("error fetching invoice by ID: %w", err)
	}
	return &invoice, nil
}

func (r *InvoiceRepository) Update(id uuid.UUID, invoice *domain.Invoice) (*domain.Invoice, error) {
	if err := r.db.Model(&domain.Invoice{}).Where("id = ?", id).Updates(invoice).Error; err != nil {
		return nil, fmt.Errorf("error updating invoice: %w", err)
	}
	return invoice, nil
}

func (r *InvoiceRepository) Delete(id uuid.UUID) error {
	if err := r.db.Delete(&domain.Invoice{}, id).Error; err != nil {
		return fmt.Errorf("error deleting invoice: %w", err)
	}
	return nil
}

// Connect to the PostgreSQL database
func Connect(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

		// Auto-migrate the Invoice table
	err = db.AutoMigrate(&domain.Invoice{})
	if err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	} else {
		log.Println("✅ AutoMigrate ran without errors")
	}
	return db, nil
}