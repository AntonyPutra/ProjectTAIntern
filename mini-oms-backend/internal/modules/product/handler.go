package product

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

// GetAll returns all products
func (h *Handler) GetAll(c echo.Context) error {
	products, err := h.service.GetAll()
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch products")
	}

	return utils.SuccessResponse(c, http.StatusOK, "Products retrieved successfully", products)
}

// GetByID returns product by ID
func (h *Handler) GetByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID")
	}

	product, err := h.service.GetByID(id)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusNotFound, "Product not found")
	}

	return utils.SuccessResponse(c, http.StatusOK, "Product retrieved successfully", product)
}

// Create creates new product (admin only)
func (h *Handler) Create(c echo.Context) error {
	var req ProductRequest
	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	product, err := h.service.Create(&req)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return utils.SuccessResponse(c, http.StatusCreated, "Product created successfully", product)
}

// Update updates product (admin only)
func (h *Handler) Update(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID")
	}

	var req ProductRequest
	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	product, err := h.service.Update(id, &req)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return utils.SuccessResponse(c, http.StatusOK, "Product updated successfully", product)
}

// Delete deletes product (admin only)
func (h *Handler) Delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID")
	}

	if err := h.service.Delete(id); err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete product")
	}

	return utils.SuccessResponse(c, http.StatusOK, "Product deleted successfully", nil)
}
