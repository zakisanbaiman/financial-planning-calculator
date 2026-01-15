package web

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/financial-planning-calculator/backend/config"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheckHandler(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := HealthCheckHandler(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "財務計画計算機 API サーバーが正常に動作しています")
}

func TestAPIInfoHandler(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := APIInfoHandler(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "財務計画計算機 API v1.0")
}

func TestSetupRoutes(t *testing.T) {
	e := echo.New()

	// Create mock controllers (nil for now since we're just testing route setup)
	controllers := &Controllers{
		FinancialData: nil,
		Calculations:  nil,
		Goals:         nil,
		Reports:       nil,
	}

	// Create minimal ServerDependencies for testing (all nils are acceptable here)
	deps := &ServerDependencies{
		FinancialPlanRepo:     nil,
		GoalRepo:              nil,
		CalculationService:    nil,
		RecommendationService: nil,
		SkipAuth:              true, // テスト用に認証をスキップ
		ServerConfig: &config.ServerConfig{
			GitHubClientID:       "test-client-id",
			GitHubClientSecret:   "test-client-secret",
			GitHubCallbackURL:    "http://localhost:8080/api/auth/github/callback",
			OAuthSuccessRedirect: "/auth/callback",
			OAuthFailureRedirect: "/login?error=oauth_failed",
		},
	}

	// This should not panic
	assert.NotPanics(t, func() {
		SetupRoutes(e, controllers, deps)
	})

	// Verify that routes are registered
	routes := e.Routes()
	assert.NotEmpty(t, routes)

	// Check for some key routes
	routePaths := make([]string, len(routes))
	for i, route := range routes {
		routePaths[i] = route.Path
	}

	assert.Contains(t, routePaths, "/health")
	assert.Contains(t, routePaths, "/api/")
	assert.Contains(t, routePaths, "/swagger/*")
}
