package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Payment struct {
	ID              uuid.UUID  `gorm:"type:uuid;primary_key;" json:"id"`
	OrderID         uuid.UUID  `gorm:"type:uuid;uniqueIndex;not null" json:"order_id"` // One payment per order
	PaymentNumber   string     `gorm:"type:varchar(50);uniqueIndex;not null" json:"payment_number"`
	Amount          float64    `gorm:"type:decimal(12,2);not null" json:"amount"`
	PaymentMethod   string     `gorm:"type:varchar(50);not null" json:"payment_method"`           // bank_transfer, e-wallet, credit_card
	Status          string     `gorm:"type:varchar(20);not null;default:'pending'" json:"status"` // pending, paid, failed
	PaymentProofURL string     `gorm:"type:varchar(500)" json:"payment_proof_url"`
	VerifiedBy      *uuid.UUID `gorm:"type:uuid" json:"verified_by"`
	VerifiedAt      *time.Time `json:"verified_at"`
	Notes           string     `gorm:"type:text" json:"notes"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`

	// Relations
	Order    *Order `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	Verifier *User  `gorm:"foreignKey:VerifiedBy" json:"verifier,omitempty"`
}

func (p *Payment) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	if p.PaymentNumber == "" {
		p.PaymentNumber = generatePaymentNumber()
	}
	if p.Status == "" {
		p.Status = "pending"
	}
	return nil
}

// generatePaymentNumber creates unique payment number
func generatePaymentNumber() string {
	return fmt.Sprintf("PAY-%s-%d", time.Now().Format("20060102"), time.Now().Unix()%10000)
}

// IsPending checks if payment is still pending
func (p *Payment) IsPending() bool {
	return p.Status == "pending"
}

// MarkAsPaid marks payment as paid
func (p *Payment) MarkAsPaid(adminID uuid.UUID) {
	p.Status = "paid"
	p.VerifiedBy = &adminID
	now := time.Now()
	p.VerifiedAt = &now
}

// MarkAsFailed marks payment as failed
func (p *Payment) MarkAsFailed(adminID uuid.UUID, reason string) {
	p.Status = "failed"
	p.VerifiedBy = &adminID
	now := time.Now()
	p.VerifiedAt = &now
	p.Notes = reason
}
