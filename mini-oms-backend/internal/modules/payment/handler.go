package payment

import (
	"fmt"
	"mini-oms-backend/internal/utils"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// GetByOrderID returns payment by order ID
func (h *Handler) GetByOrderID(c echo.Context) error {
	orderID, err := uuid.Parse(c.Param("orderId"))
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid order ID")
	}

	payment, err := h.service.GetByOrderID(orderID)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusNotFound, "Payment not found")
	}

	return utils.SuccessResponse(c, http.StatusOK, "Payment retrieved successfully", payment)
}

// Create creates new payment
func (h *Handler) Create(c echo.Context) error {
	var req CreatePaymentRequest
	if err := c.Bind(&req); err != nil {
		utils.LogError("PaymentService", "", "CreatePayment", err, "Failed to bind request")
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	// Log incoming request
	utils.LogInfo("PaymentService", req.OrderID.String(), "CreatePayment", fmt.Sprintf("Payment method: %s", req.PaymentMethod))

	payment, err := h.service.CreatePayment(&req)
	if err != nil {
		utils.LogError("PaymentService", req.OrderID.String(), "CreatePayment", err, "Failed to create payment")
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	// Public response for creation
	utils.LogInfo("PaymentService", payment.ID.String(), "CreatePayment", "Payment created successfully")
	return utils.SuccessResponse(c, http.StatusCreated, "Payment created successfully", payment)
}

// Verify verifies a payment (Admin only)
// @Summary Verify payment
// @Tags payments
// @Accept json
// @Produce json
// @Param id path string true "Payment ID"
// @Success 200 {object} utils.APIResponse
// @Router /api/payments/{id}/verify [post]
func (h *Handler) Verify(c echo.Context) error {
	idParam := c.Param("id")
	paymentID, err := uuid.Parse(idParam)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid payment ID")
	}

	// Get admin ID from context
	adminID := c.Get("user_id").(uuid.UUID)

	utils.LogInfo("PaymentService", paymentID.String(), "VerifyPayment", fmt.Sprintf("Verification requested by admin %s", adminID))

	payment, err := h.service.VerifyPayment(paymentID, adminID)
	if err != nil {
		utils.LogError("PaymentService", paymentID.String(), "VerifyPayment", err, "Verification failed")
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	utils.LogInfo("PaymentService", paymentID.String(), "VerifyPayment", "Payment verified successfully")
	return utils.SuccessResponse(c, http.StatusOK, "Payment verified successfully", payment)
}
