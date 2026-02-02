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

	return utils.SuccessResponse(c, http.StatusCreated, "Payment created successfully", payment)
}
