package middlewares

import (
	"mini-oms-backend/internal/config"
	"mini-oms-backend/internal/utils"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// JWTMiddleware validates JWT token
func JWTMiddleware(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Extract token from Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return utils.ErrorResponse(c, http.StatusUnauthorized, "Missing authorization header")
			}

			// Check Bearer prefix
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid authorization format")
			}

			tokenString := parts[1]

			// Validate token
			claims, err := utils.ValidateJWT(cfg, tokenString)
			if err != nil {
				return utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid or expired token")
			}

			// Set user context
			c.Set("user_id", claims.UserID)
			c.Set("user_email", claims.Email)
			c.Set("user_role", claims.Role)

			return next(c)
		}
	}
}
