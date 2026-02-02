package middlewares

import (
	"mini-oms-backend/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

// AdminOnlyMiddleware checks if user has admin role
func AdminOnlyMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userRole := c.Get("user_role")
			if userRole == nil {
				return utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
			}

			role, ok := userRole.(string)
			if !ok || role != "admin" {
				return utils.ErrorResponse(c, http.StatusForbidden, "Access forbidden: admin only")
			}

			return next(c)
		}
	}
}
