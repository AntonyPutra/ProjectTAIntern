package utils

import (
	"github.com/labstack/echo/v4"
)

// APIResponse is standard API response structure
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

// SuccessResponse returns success response
func SuccessResponse(c echo.Context, statusCode int, message string, data interface{}) error {
	return c.JSON(statusCode, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SuccessResponseWithMeta returns success response with pagination meta
func SuccessResponseWithMeta(c echo.Context, statusCode int, message string, data interface{}, meta interface{}) error {
	return c.JSON(statusCode, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

// ErrorResponse returns error response
func ErrorResponse(c echo.Context, statusCode int, message string) error {
	return c.JSON(statusCode, APIResponse{
		Success: false,
		Message: message,
	})
}

// ValidationErrorResponse returns validation error response
func ValidationErrorResponse(c echo.Context, statusCode int, message string, errors interface{}) error {
	return c.JSON(statusCode, APIResponse{
		Success: false,
		Message: message,
		Errors:  errors,
	})
}
