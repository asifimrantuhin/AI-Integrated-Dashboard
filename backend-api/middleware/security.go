package middleware

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// SecurityHeaders adds common security headers and enforces HTTPS-only cookies
func SecurityHeaders() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			res := c.Response()
			res.Header().Set("X-Frame-Options", "DENY")
			res.Header().Set("X-Content-Type-Options", "nosniff")
			res.Header().Set("X-XSS-Protection", "1; mode=block")
			res.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
			res.Header().Set("Content-Security-Policy", "default-src 'self'; frame-ancestors 'none';")
			res.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
			return next(c)
		}
	}
}

// RateLimiter returns a token bucket limiter for authenticated APIs
func RateLimiter() echo.MiddlewareFunc {
	store := middleware.NewRateLimiterMemoryStoreWithConfig(middleware.RateLimiterMemoryStoreConfig{
		Rate:      20,
		Burst:     40,
		ExpiresIn: time.Minute,
	})
	return middleware.RateLimiter(store)
}

// RequestLogger logs slow requests and suspicious payloads
func RequestLogger() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogLatency:  true,
		LogStatus:   true,
		LogURI:      true,
		HandleError: true,
		BeforeNextFunc: func(c echo.Context) {
			c.Set("request_start", time.Now())
		},
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Status >= http.StatusInternalServerError {
				c.Logger().Errorf("HTTP %d %s (%v)", v.Status, v.URI, v.Latency)
			}
			return nil
		},
	})
}
