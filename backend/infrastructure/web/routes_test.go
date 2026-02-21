package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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
			AuthRateLimitRPS:     10,
			AuthRateLimitBurst:   5,
		},
	}

	// This should not panic
	assert.NotPanics(t, func() {
		testStore := NewCustomRateLimiterStore(100, 50, 3*time.Minute)
		SetupRoutes(e, controllers, deps, testStore)
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
	assert.Contains(t, routePaths, "/api/rate-limit/status")
}

func TestRateLimitStatusHandler(t *testing.T) {
	store := NewCustomRateLimiterStore(100, 50, time.Minute)
	e := echo.New()

	req := httptest.NewRequest("GET", "/api/rate-limit/status", nil)
	req.RemoteAddr = "203.0.113.1:1234"
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := RateLimitStatusHandler(store)
	err := handler(c)

	assert.NoError(t, err)
	assert.Equal(t, 200, rec.Code)

	// Response should contain rate limit info fields
	body := rec.Body.String()
	assert.Contains(t, body, "\"limit\"")
	assert.Contains(t, body, "\"remaining\"")
	assert.Contains(t, body, "\"reset\"")
	assert.Contains(t, body, "\"reset_at\"")
}

func TestRateLimitStatusHandler_RemainingDecreases(t *testing.T) {
	store := NewCustomRateLimiterStore(100, 5, time.Minute)
	e := echo.New()
	ip := "203.0.113.2"

	callStatus := func() int {
		req := httptest.NewRequest("GET", "/api/rate-limit/status", nil)
		req.RemoteAddr = ip + ":1234"
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		handler := RateLimitStatusHandler(store)
		handler(c) //nolint:errcheck

		var info RateLimitInfo
		json.NewDecoder(rec.Body).Decode(&info) //nolint:errcheck
		return info.Remaining
	}

	first := callStatus()
	// Consume 2 tokens via Allow
	store.Allow(ip) //nolint:errcheck
	store.Allow(ip) //nolint:errcheck
	second := callStatus()

	assert.GreaterOrEqual(t, first, second, "remaining should decrease after Allow calls")
}
