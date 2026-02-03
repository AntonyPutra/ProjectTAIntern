package payment

import (
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
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	payment, err := h.service.CreatePayment(&req)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	// Public response for creation
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

	payment, err := h.service.VerifyPayment(paymentID, adminID)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return utils.SuccessResponse(c, http.StatusOK, "Payment verified successfully", payment)
}
