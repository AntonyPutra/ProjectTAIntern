package order

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

// GetAll returns all orders (admin) or user's orders
func (h *Handler) GetAll(c echo.Context) error {
	userRole := c.Get("user_role").(string)
	userID := c.Get("user_id").(uuid.UUID)

	utils.LogInfo("OrderService", userID.String(), "GetAllOrders", fmt.Sprintf("Fetch requested by role: %s", userRole))

	var orders []interface{}

	if userRole == "admin" {
		// Admin sees all orders
		ordersData, err := h.service.GetAll()
		if err != nil {
			utils.LogError("OrderService", userID.String(), "GetAllOrders", err, "Failed to fetch orders")
			return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch orders")
		}
		for _, order := range ordersData {
			orders = append(orders, order)
		}
	} else {
		// User sees only their orders
		ordersData, err := h.service.GetByUserID(userID)
		if err != nil {
			utils.LogError("OrderService", userID.String(), "GetAllOrders", err, "Failed to fetch orders")
			return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch orders")
		}
		for _, order := range ordersData {
			orders = append(orders, order)
		}
	}

	utils.LogInfo("OrderService", userID.String(), "GetAllOrders", fmt.Sprintf("Retrieved %d orders", len(orders)))
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
		utils.LogError("OrderService", userID.String(), "GetOrder", nil, "Unauthorized access attempt to order "+id.String())
		return utils.ErrorResponse(c, http.StatusForbidden, "Access forbidden")
	}

	utils.LogInfo("OrderService", userID.String(), "GetOrder", "Order retrieved: "+id.String())
	return utils.SuccessResponse(c, http.StatusOK, "Order retrieved successfully", order)
}

// Create creates new order
func (h *Handler) Create(c echo.Context) error {
	var req CreateOrderRequest
	if err := c.Bind(&req); err != nil {
		utils.LogError("OrderService", "", "CreateOrder", err, "Failed to bind request")
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	userID := c.Get("user_id").(uuid.UUID)

	// Log incoming request
	utils.LogInfo("OrderService", userID.String(), "CreateOrder", fmt.Sprintf("Request received from user %s", userID), fmt.Sprintf("Items count: %d", len(req.Items)))

	response, err := h.service.CreateOrder(userID, &req)
	if err != nil {
		utils.LogError("OrderService", userID.String(), "CreateOrder", err, "Failed to create order")
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	utils.LogInfo("OrderService", response.ID.String(), "CreateOrder", "Order created successfully")
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

	utils.LogInfo("OrderService", orderID.String(), "CancelOrder", fmt.Sprintf("Cancellation requested by user %s (role: %s)", userID, role))

	if err := h.service.CancelOrder(orderID, userID, role); err != nil {
		utils.LogError("OrderService", orderID.String(), "CancelOrder", err, "Cancellation failed")
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	utils.LogInfo("OrderService", orderID.String(), "CancelOrder", "Order canceled successfully")
	return utils.SuccessResponse(c, http.StatusOK, "Order canceled successfully", nil)
}

// GetStats returns admin dashboard statistics
func (h *Handler) GetStats(c echo.Context) error {
	// Log request context (Admin ID usually)
	// Note: We might need to ensure middlewares set user_id for this route
	var role string
	var userID string
	if r := c.Get("user_role"); r != nil {
		role = r.(string)
	}
	if u := c.Get("user_id"); u != nil {
		userID = u.(uuid.UUID).String()
	}

	utils.LogInfo("OrderService", userID, "GetStats", "Admin dashboard stats requested by role: "+role)

	stats, err := h.service.GetStats()
	if err != nil {
		utils.LogError("OrderService", userID, "GetStats", err, "Failed to fetch stats")
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	utils.LogInfo("OrderService", userID, "GetStats", fmt.Sprintf("Stats retrieved: Orders=%d, Revenue=%.2f", stats["total_orders"], stats["total_revenue"]))
	return utils.SuccessResponse(c, http.StatusOK, "Statistics retrieved", stats)
}
