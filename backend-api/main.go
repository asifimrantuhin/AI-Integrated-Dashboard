package main

import (
	"idash-backend-api/config"
	"idash-backend-api/database"
	"idash-backend-api/routes"
	"log"
	"net/http"

	"idash-backend-api/middleware"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	// Middleware
	e.Use(echoMiddleware.Recover())
	e.Use(middleware.SecurityHeaders())
	e.Use(middleware.RequestLogger())
	e.Use(middleware.RateLimiter())
	// Allow common headers including a custom x-request-id used by the frontend
	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins:     []string{"https://idash.example.com", "http://localhost:3000"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderAccept, echo.HeaderContentType, echo.HeaderAuthorization, "X-Requested-With", "x-request-id"},
		AllowCredentials: true,
	}))

	// Load configuration
	config.LoadConfig()

	// Initialize Database
	database.InitDB()

	// Register routes
	routes.SetupRoutes(e)

	// Start Server
	log.Printf("Starting server on port %s", config.AppConfig.APIPort)
	e.Logger.Fatal(e.Start(":" + config.AppConfig.APIPort))
}
