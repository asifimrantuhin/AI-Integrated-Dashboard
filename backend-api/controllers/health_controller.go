package controllers

import (
	"net/http"
	"time"

	"idash-backend-api/database"

	"github.com/labstack/echo/v4"
)

// HealthCheck returns service health and DB connectivity
func HealthCheck(c echo.Context) error {
	resp := map[string]interface{}{
		"status": "ok",
		"time":   time.Now().UTC().Format(time.RFC3339),
	}

	// Check DB connectivity if initialized
	if database.DB != nil {
		sqlDB, err := database.DB.DB()
		if err == nil {
			if err = sqlDB.Ping(); err == nil {
				resp["db"] = "ok"
			} else {
				resp["db"] = "error: " + err.Error()
			}
		} else {
			resp["db"] = "error: " + err.Error()
		}
	} else {
		resp["db"] = "not_initialized"
	}

	return c.JSON(http.StatusOK, resp)
}
