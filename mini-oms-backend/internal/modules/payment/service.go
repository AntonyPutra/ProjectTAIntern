package payment

import (
	"errors"
	"mini-oms-backend/internal/models"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

type CreatePaymentRequest struct {
	OrderID         uuid.UUID `json:"order_id"`
	PaymentMethod   string    `json:"payment_method"` // bank_transfer, e-wallet, credit_card
	PaymentProofURL string    `json:"payment_proof_url"`
	Notes           string    `json:"notes"`
}

func (s *Service) GetByOrderID(orderID uuid.UUID) (*models.Payment, error) {
	return s.repo.FindByOrderID(orderID)
}

func (s *Service) CreatePayment(req *CreatePaymentRequest) (*models.Payment, error) {
	// Check if order exists
	order, err := s.repo.FindOrderByID(req.OrderID)
	if err != nil {
		return nil, errors.New("order not found")
	}

	// Check if payment already exists for this order
	existingPayment, _ := s.repo.FindByOrderID(req.OrderID)
	if existingPayment != nil {
		return nil, errors.New("payment already exists for this order")
	}

	// Validate payment method
	validMethods := map[string]bool{
		"bank_transfer": true,
		"e-wallet":      true,
		"credit_card":   true,
	}
	if !validMethods[req.PaymentMethod] {
		return nil, errors.New("invalid payment method")
	}

	// Create payment
	payment := &models.Payment{
		OrderID:         req.OrderID,
		Amount:          order.TotalAmount,
		PaymentMethod:   req.PaymentMethod,
		Status:          "pending",
		PaymentProofURL: req.PaymentProofURL,
		Notes:           req.Notes,
	}

	if err := s.repo.Create(payment); err != nil {
		return nil, err
	}

	return payment, nil
}
