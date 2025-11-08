package middleware

import (
	"idash-backend-api/utils"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Authorization header required"})
			}

			// Extract token from "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid authorization header"})
			}

			token := parts[1]
			claims, err := utils.ValidateJWT(token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
			}

			// Set user info in context
			c.Set("user_id", claims.UserID)
			c.Set("user_roles", claims.Roles)
			c.Set("user_permissions", claims.Permissions)
			c.Set("default_company_id", claims.DefaultCompanyID)
			c.Set("company_ids", claims.CompanyIDs)

			return next(c)
		}
	}
}

// RequireRoles ensures the authenticated user has at least one of the required roles
func RequireRoles(requiredRoles ...string) echo.MiddlewareFunc {
	required := make(map[string]struct{}, len(requiredRoles))
	for _, role := range requiredRoles {
		required[strings.ToLower(role)] = struct{}{}
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			rolesAny := c.Get("user_roles")
			roles, ok := rolesAny.([]string)
			if !ok {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "Insufficient permissions"})
			}

			for _, role := range roles {
				if _, exists := required[strings.ToLower(role)]; exists {
					return next(c)
				}
			}

			return c.JSON(http.StatusForbidden, map[string]string{"error": "Insufficient permissions"})
		}
	}
}

// RequirePermissions ensures the authenticated user has at least one of the required permissions
func RequirePermissions(requiredPermissions ...string) echo.MiddlewareFunc {
	required := make(map[string]struct{}, len(requiredPermissions))
	for _, permission := range requiredPermissions {
		required[strings.ToLower(permission)] = struct{}{}
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			permissionsAny := c.Get("user_permissions")
			permissions, ok := permissionsAny.([]string)
			if !ok {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "Insufficient permissions"})
			}

			for _, permission := range permissions {
				if _, exists := required[strings.ToLower(permission)]; exists {
					return next(c)
				}
			}

			return c.JSON(http.StatusForbidden, map[string]string{"error": "Insufficient permissions"})
		}
	}
}

