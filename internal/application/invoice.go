package application

import (
	"github.com/carfdev/microservice-go/internal/domain"
	"github.com/google/uuid"
)

// Interface for the repository
type InvoiceRepository interface {
	Create(invoice *domain.Invoice) (*domain.Invoice, error)
	GetAll() ([]domain.Invoice, error)
	GetByID(id uuid.UUID) (*domain.Invoice, error)
	Update(id uuid.UUID, invoice *domain.Invoice) (*domain.Invoice, error)
	Delete(id uuid.UUID) error
}

// Service to manage invoice use cases
type InvoiceService struct {
	repo InvoiceRepository
}

func NewInvoiceService(repo InvoiceRepository) *InvoiceService {
	return &InvoiceService{repo}
}

func (s *InvoiceService) CreateInvoice(invoice *domain.Invoice) (*domain.Invoice, error) {
	// Add your business logic here
	return s.repo.Create(invoice)
}

func (s *InvoiceService) GetAllInvoices() ([]domain.Invoice, error) {
	return s.repo.GetAll()
}

func (s *InvoiceService) GetInvoiceByID(id uuid.UUID) (*domain.Invoice, error) {
	return s.repo.GetByID(id)
}

func (s *InvoiceService) UpdateInvoice(id uuid.UUID, invoice *domain.Invoice) (*domain.Invoice, error) {
	return s.repo.Update(id, invoice)
}

func (s *InvoiceService) DeleteInvoice(id uuid.UUID) error {
	return s.repo.Delete(id)
}

