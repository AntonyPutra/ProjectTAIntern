package payment

import (
	"errors"
	"mini-oms-backend/internal/models"
	"mini-oms-backend/internal/utils"
	"time"

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
	existingPayment, err := s.repo.FindByOrderID(req.OrderID)
	if err == nil && existingPayment.ID != uuid.Nil {
		// If existing payment is pending, return error
		if existingPayment.Status == "pending" || existingPayment.Status == "success" {
			return nil, errors.New("payment already exists for this order")
		}
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

	// Make PaymentProofURL optional (default to "-")
	proofURL := req.PaymentProofURL
	if proofURL == "" {
		proofURL = "-"
	}

	// Create payment
	payment := &models.Payment{
		OrderID:         req.OrderID,
		Amount:          order.TotalAmount,
		PaymentMethod:   req.PaymentMethod,
		Status:          "pending",
		PaymentProofURL: proofURL,
		Notes:           req.Notes,
	}

	if err := s.repo.Create(payment); err != nil {
		return nil, err
	}

	return payment, nil
}

func (s *Service) VerifyPayment(paymentID uuid.UUID, adminID uuid.UUID) (*models.Payment, error) {
	// Start transaction
	tx := s.repo.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Find payment
	var payment models.Payment
	if err := tx.First(&payment, "id = ?", paymentID).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("payment not found")
	}

	if payment.Status == "success" {
		tx.Rollback()
		return nil, errors.New("payment already verified")
	}

	// Update payment status
	now := time.Now()
	payment.Status = "success"
	payment.VerifiedBy = &adminID
	payment.VerifiedAt = &now

	if err := tx.Save(&payment).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Update order status to 'processing' (Paid)
	// We update directly the order table
	if err := tx.Model(&models.Order{}).Where("id = ?", payment.OrderID).Update("status", models.OrderStatusProcessing).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Log Audit
	if err := utils.LogAudit(tx, adminID, "PAYMENT_VERIFIED", "Payment", payment.ID, "Payment verified by admin"); err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &payment, nil
}
