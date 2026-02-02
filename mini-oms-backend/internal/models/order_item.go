package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderItem struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	OrderID      uuid.UUID `gorm:"type:uuid;not null;index" json:"order_id"`
	ProductID    uuid.UUID `gorm:"type:uuid;not null;index" json:"product_id"`
	ProductName  string    `gorm:"type:varchar(255);not null" json:"product_name"`   // Snapshot
	ProductPrice float64   `gorm:"type:decimal(12,2);not null" json:"product_price"` // Snapshot
	Quantity     int       `gorm:"type:integer;not null" json:"quantity"`
	Subtotal     float64   `gorm:"type:decimal(12,2);not null" json:"subtotal"`
	CreatedAt    time.Time `json:"created_at"`

	// Relations
	Order   *Order   `gorm:"foreignKey:OrderID" json:"-"`
	Product *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

func (oi *OrderItem) BeforeCreate(tx *gorm.DB) error {
	if oi.ID == uuid.Nil {
		oi.ID = uuid.New()
	}
	// Calculate subtotal
	if oi.Subtotal == 0 {
		oi.Subtotal = oi.ProductPrice * float64(oi.Quantity)
	}
	return nil
}
