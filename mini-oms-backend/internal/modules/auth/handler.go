package auth

import (
	"mini-oms-backend/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Register handles user registration
// @Summary Register new user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Register Request"
// @Success 201 {object} utils.APIResponse
// @Router /api/auth/register [post]
func (h *Handler) Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	// Basic validation
	if req.Name == "" || req.Email == "" || req.Password == "" {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Name, email, and password are required")
	}

	if len(req.Password) < 6 {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Password must be at least 6 characters")
	}

	// Register user
	response, err := h.service.Register(&req)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return utils.SuccessResponse(c, http.StatusCreated, "User registered successfully", response)
}

// Login handles user login
// @Summary Login user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login Request"
// @Success 200 {object} utils.APIResponse
// @Router /api/auth/login [post]
func (h *Handler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	// Basic validation
	if req.Email == "" || req.Password == "" {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Email and password are required")
	}

	// Login user
	response, err := h.service.Login(&req)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusUnauthorized, err.Error())
	}

	return utils.SuccessResponse(c, http.StatusOK, "Login successful", response)
}
