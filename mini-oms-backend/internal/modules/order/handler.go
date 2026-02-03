package order

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

// GetAll returns all orders (admin) or user's orders
func (h *Handler) GetAll(c echo.Context) error {
	userRole := c.Get("user_role").(string)
	userID := c.Get("user_id").(uuid.UUID)

	var orders []interface{}

	if userRole == "admin" {
		// Admin sees all orders
		ordersData, err := h.service.GetAll()
		if err != nil {
			return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch orders")
		}
		for _, order := range ordersData {
			orders = append(orders, order)
		}
	} else {
		// User sees only their orders
		ordersData, err := h.service.GetByUserID(userID)
		if err != nil {
			return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch orders")
		}
		for _, order := range ordersData {
			orders = append(orders, order)
		}
	}

	return utils.SuccessResponse(c, http.StatusOK, "Orders retrieved successfully", orders)
}

// GetByID returns order by ID
func (h *Handler) GetByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid order ID")
	}

	order, err := h.service.GetByID(id)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusNotFound, "Order not found")
	}

	// Check authorization: user can only see their own orders
	userRole := c.Get("user_role").(string)
	userID := c.Get("user_id").(uuid.UUID)

	if userRole != "admin" && order.UserID != userID {
		return utils.ErrorResponse(c, http.StatusForbidden, "Access forbidden")
	}

	return utils.SuccessResponse(c, http.StatusOK, "Order retrieved successfully", order)
}

// Create creates new order
func (h *Handler) Create(c echo.Context) error {
	var req CreateOrderRequest
	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	userID := c.Get("user_id").(uuid.UUID)

	response, err := h.service.CreateOrder(userID, &req)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return utils.SuccessResponse(c, http.StatusCreated, "Order created successfully", response)
}

// Cancel cancels an order
// @Summary Cancel order
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} utils.APIResponse
// @Router /api/orders/{id}/cancel [post]
func (h *Handler) Cancel(c echo.Context) error {
	idParam := c.Param("id")
	orderID, err := uuid.Parse(idParam)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid order ID")
	}

	userID := c.Get("user_id").(uuid.UUID)
	role := c.Get("user_role").(string)

	if err := h.service.CancelOrder(orderID, userID, role); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return utils.SuccessResponse(c, http.StatusOK, "Order canceled successfully", nil)
}

// GetStats returns admin dashboard statistics
func (h *Handler) GetStats(c echo.Context) error {
	stats, err := h.service.GetStats()
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	return utils.SuccessResponse(c, http.StatusOK, "Statistics retrieved", stats)
}
