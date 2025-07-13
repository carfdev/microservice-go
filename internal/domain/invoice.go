package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Invoice entity
type Invoice struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Amount    float64   `gorm:"not null" json:"amount"`
	Customer  string    `gorm:"not null" json:"customer"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
// Methods for Invoice (like validations, etc.)
func (i *Invoice) Validate() error {
	if i.Amount <= 0 {
		return fmt.Errorf("invalid amount")
	}
	// Add more business logic here
	return nil
}

// BeforeCreate hook to set UUID before insert if not set
func (i *Invoice) BeforeCreate(tx *gorm.DB) (err error) {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return
}

