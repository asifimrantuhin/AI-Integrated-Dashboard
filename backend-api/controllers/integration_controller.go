package controllers

import (
	"net/http"

	"idash-backend-api/services"

	"github.com/labstack/echo/v4"
)

// TriggerIntegrationSync kicks off the ETL + forecast regeneration pipeline
func TriggerIntegrationSync(c echo.Context) error {
	if services.Integration().IsRunning() {
		return c.JSON(http.StatusConflict, map[string]string{"error": "integration job already running"})
	}

	userID := c.Get("user_id")
	go services.Integration().RunFullSync(userID)

	return c.JSON(http.StatusAccepted, map[string]string{"status": "integration sync started"})
}

// GetIntegrationStatus returns the last known status for the integration pipeline
func GetIntegrationStatus(c echo.Context) error {
	status := services.Integration().Status()
	return c.JSON(http.StatusOK, status)
}
