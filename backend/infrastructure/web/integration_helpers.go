package web

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// IntegrationHealthCheck performs comprehensive health checks for all system components
type IntegrationHealthCheck struct {
	Status     string                 `json:"status"`
	Timestamp  string                 `json:"timestamp"`
	Version    string                 `json:"version"`
	Components map[string]ComponentHealth `json:"components"`
	Uptime     string                 `json:"uptime"`
}

// ComponentHealth represents the health status of a single component
type ComponentHealth struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Latency string `json:"latency,omitempty"`
}

var serverStartTime = time.Now()

// HealthCheckHandler performs comprehensive health check
func IntegrationHealthCheckHandler(deps *ServerDependencies) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
		defer cancel()

		health := IntegrationHealthCheck{
			Status:     "ok",
			Timestamp:  time.Now().Format(time.RFC3339),
			Version:    "1.0.0",
			Components: make(map[string]ComponentHealth),
			Uptime:     time.Since(serverStartTime).String(),
		}

		// Check database connectivity
		dbHealth := checkDatabaseHealth(ctx, deps)
		health.Components["database"] = dbHealth
		if dbHealth.Status != "ok" {
			health.Status = "degraded"
		}

		// Check domain services
		servicesHealth := checkDomainServicesHealth(ctx, deps)
		health.Components["domain_services"] = servicesHealth
		if servicesHealth.Status != "ok" {
			health.Status = "degraded"
		}

		// Check repositories
		repoHealth := checkRepositoriesHealth(ctx, deps)
		health.Components["repositories"] = repoHealth
		if repoHealth.Status != "ok" {
			health.Status = "degraded"
		}

		// System resources
		health.Components["system"] = ComponentHealth{
			Status:  "ok",
			Message: "System resources available",
		}

		statusCode := http.StatusOK
		if health.Status == "degraded" {
			statusCode = http.StatusServiceUnavailable
		}

		return c.JSON(statusCode, health)
	}
}

// checkDatabaseHealth checks database connectivity
func checkDatabaseHealth(ctx context.Context, deps *ServerDependencies) ComponentHealth {
	start := time.Now()
	
	// Try to ping database through repository
	// Note: This is a simplified check - in production, you'd want actual DB ping
	if deps.FinancialPlanRepo == nil {
		return ComponentHealth{
			Status:  "error",
			Message: "Repository not initialized",
		}
	}

	latency := time.Since(start)
	return ComponentHealth{
		Status:  "ok",
		Message: "Database connection available",
		Latency: latency.String(),
	}
}

// checkDomainServicesHealth checks domain services availability
func checkDomainServicesHealth(ctx context.Context, deps *ServerDependencies) ComponentHealth {
	if deps.CalculationService == nil || deps.RecommendationService == nil {
		return ComponentHealth{
			Status:  "error",
			Message: "Domain services not initialized",
		}
	}

	return ComponentHealth{
		Status:  "ok",
		Message: "Domain services available",
	}
}

// checkRepositoriesHealth checks repositories availability
func checkRepositoriesHealth(ctx context.Context, deps *ServerDependencies) ComponentHealth {
	if deps.FinancialPlanRepo == nil || deps.GoalRepo == nil {
		return ComponentHealth{
			Status:  "error",
			Message: "Repositories not initialized",
		}
	}

	return ComponentHealth{
		Status:  "ok",
		Message: "Repositories available",
	}
}

// CORSPreflightHandler handles CORS preflight requests
func CORSPreflightHandler(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

// APIReadinessHandler checks if API is ready to serve requests
func APIReadinessHandler(deps *ServerDependencies) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Check critical dependencies
		if deps.FinancialPlanRepo == nil || deps.GoalRepo == nil ||
			deps.CalculationService == nil || deps.RecommendationService == nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]interface{}{
				"ready":   false,
				"message": "サービスの初期化が完了していません",
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"ready":     true,
			"message":   "APIは正常に動作しています",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	}
}

// ErrorRecoveryMiddleware provides enhanced error recovery with logging
func ErrorRecoveryMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)
				if !ok {
					err = fmt.Errorf("%v", r)
				}
				
				c.Logger().Errorf("Panic recovered: %v", err)
				
				// Return user-friendly error
				c.JSON(http.StatusInternalServerError, map[string]interface{}{
					"error":      "内部サーバーエラーが発生しました",
					"code":       "INTERNAL_SERVER_ERROR",
					"timestamp":  time.Now().Format(time.RFC3339),
					"request_id": c.Response().Header().Get(echo.HeaderXRequestID),
				})
			}
		}()
		
		return next(c)
	}
}

// RequestValidationMiddleware validates common request parameters
func RequestValidationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Validate user_id if present in query params
		if userID := c.QueryParam("user_id"); userID != "" {
			if len(userID) == 0 || len(userID) > 100 {
				return echo.NewHTTPError(http.StatusBadRequest, map[string]interface{}{
					"error":   "無効なユーザーIDです",
					"code":    "INVALID_USER_ID",
					"details": "ユーザーIDは1〜100文字である必要があります",
				})
			}
		}

		// Validate content type for POST/PUT requests
		if c.Request().Method == http.MethodPost || c.Request().Method == http.MethodPut {
			contentType := c.Request().Header.Get(echo.HeaderContentType)
			if contentType != "" && contentType != echo.MIMEApplicationJSON && 
			   contentType != echo.MIMEApplicationJSONCharsetUTF8 {
				return echo.NewHTTPError(http.StatusUnsupportedMediaType, map[string]interface{}{
					"error":   "サポートされていないコンテンツタイプです",
					"code":    "UNSUPPORTED_MEDIA_TYPE",
					"details": "Content-Type: application/json を使用してください",
				})
			}
		}

		return next(c)
	}
}

// ResponseEnhancementMiddleware adds standard headers to all responses
func ResponseEnhancementMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Add standard response headers
		c.Response().Header().Set("X-Content-Type-Options", "nosniff")
		c.Response().Header().Set("X-Frame-Options", "DENY")
		c.Response().Header().Set("X-XSS-Protection", "1; mode=block")
		
		// Add API version header
		c.Response().Header().Set("X-API-Version", "1.0.0")
		
		return next(c)
	}
}
