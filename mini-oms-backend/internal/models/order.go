package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Order struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;index" json:"user_id"`
	OrderNumber string         `gorm:"type:varchar(50);uniqueIndex" json:"order_number"` // New field
	TotalAmount float64        `gorm:"type:decimal(12,2);not null" json:"total_amount"`
	Status      string         `gorm:"type:varchar(20);default:'created'" json:"status"` // created, processing, completed, canceled
	Notes       string         `gorm:"type:text" json:"notes"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User       *User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	OrderItems []OrderItem `gorm:"foreignKey:OrderID" json:"items,omitempty"`
	Payment    *Payment    `gorm:"foreignKey:OrderID" json:"payment,omitempty"`
}

// Order Status Constants
const (
	OrderStatusCreated    = "created"
	OrderStatusProcessing = "processing" // Paid but not yet shipped
	OrderStatusCompleted  = "completed"
	OrderStatusCanceled   = "canceled"
)

func (o *Order) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	if o.OrderNumber == "" {
		o.OrderNumber = generateOrderNumber()
	}
	if o.Status == "" {
		o.Status = "created"
	}
	return nil
}

// generateOrderNumber creates unique order number
func generateOrderNumber() string {
	return fmt.Sprintf("ORD-%s-%d", time.Now().Format("20060102"), time.Now().Unix()%10000)
}

// CanBeCanceled checks if order can be canceled
func (o *Order) CanBeCanceled() bool {
	// Only 'created' (Unpaid) orders can be canceled by user/admin
	return o.Status == "created"
}

// CanBeCompleted checks if order can be marked as completed
func (o *Order) CanBeCompleted() bool {
	return o.Status == "processing"
}
