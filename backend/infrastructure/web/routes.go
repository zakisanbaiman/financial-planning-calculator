package web

import (
	"net/http"
	"time"

	"github.com/financial-planning-calculator/backend/infrastructure/monitoring"
	"github.com/financial-planning-calculator/backend/infrastructure/web/controllers"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// Controllers holds all controller instances
type Controllers struct {
	Auth          *controllers.AuthController
	TwoFactor     *controllers.TwoFactorController
	WebAuthn      *controllers.WebAuthnController
	FinancialData *controllers.FinancialDataController
	Calculations  *controllers.CalculationsController
	Goals         *controllers.GoalsController
	Reports       *controllers.ReportsController
}

// SetupRoutes configures all routes based on OpenAPI specification
func SetupRoutes(e *echo.Echo, controllers *Controllers, deps *ServerDependencies, rateLimitStore *CustomRateLimiterStore) {
	// Swagger UI
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// 監視・メトリクスエンドポイント
	e.GET("/metrics", monitoring.PrometheusHandler())

	// ヘルスチェック
	e.GET("/health", HealthCheckHandler)
	e.GET("/health/detailed", IntegrationHealthCheckHandler(deps))
	e.GET("/ready", APIReadinessHandler(deps))

	// CORS preflight
	e.OPTIONS("/*", CORSPreflightHandler)

	// API ルートグループ
	api := e.Group("/api")

	// Apply integration middleware
	api.Use(ErrorRecoveryMiddleware)
	api.Use(RequestValidationMiddleware)
	api.Use(ResponseEnhancementMiddleware)

	// API情報エンドポイント
	api.GET("/", APIInfoHandler)

	// レートリミットステータスエンドポイント（認証不要）
	api.GET("/rate-limit/status", RateLimitStatusHandler(rateLimitStore))

	// 認証エンドポイント（認証不要）
	setupAuthRoutes(api, controllers.Auth, deps)

	// 計算エンドポイント（ゲストモード対応のため認証不要）
	setupCalculationRoutes(api, controllers.Calculations)

	// 目標管理エンドポイント（ゲストモード対応のため認証不要）
	setupGoalRoutes(api, controllers.Goals)

	// 認証が必要なエンドポイント用グループ
	protected := api.Group("")
	if authMiddleware := deps.JWTAuthMiddlewareFunc(); authMiddleware != nil {
		protected.Use(authMiddleware)
	}

	// パスキー認証エンドポイント
	setupPasskeyRoutes(api, protected, controllers.WebAuthn)

	// 2段階認証エンドポイント（認証が必要）
	setup2FARoutes(protected, controllers.TwoFactor)

	// 財務データ管理エンドポイント
	setupFinancialDataRoutes(protected, controllers.FinancialData)

	// レポート生成エンドポイント
	setupReportRoutes(protected, controllers.Reports)
}

// setupAuthRoutes sets up authentication routes
func setupAuthRoutes(api *echo.Group, controller *controllers.AuthController, deps *ServerDependencies) {
	auth := api.Group("/auth")

	auth.POST("/register", controller.Register) // POST /api/auth/register
	auth.POST("/login", controller.Login)       // POST /api/auth/login
	auth.POST("/refresh", controller.Refresh)   // POST /api/auth/refresh
	auth.POST("/logout", controller.Logout)     // POST /api/auth/logout

	// GitHub OAuth routes with middleware (Issue: #67)
	githubOAuth := auth.Group("/github")
	githubOAuth.Use(GitHubOAuthMiddleware(deps.ServerConfig))
	githubOAuth.GET("", controller.GitHubLogin)            // GET /api/auth/github
	githubOAuth.GET("/callback", controller.GitHubCallback) // GET /api/auth/github/callback
}

// setup2FARoutes sets up two-factor authentication routes
func setup2FARoutes(api *echo.Group, controller *controllers.TwoFactorController) {
	twoFactor := api.Group("/auth/2fa")

	twoFactor.GET("/status", controller.Get2FAStatus)                   // GET /api/auth/2fa/status
	twoFactor.POST("/setup", controller.Setup2FA)                       // POST /api/auth/2fa/setup
	twoFactor.POST("/enable", controller.Enable2FA)                     // POST /api/auth/2fa/enable
	twoFactor.POST("/verify", controller.Verify2FA)                     // POST /api/auth/2fa/verify
	twoFactor.DELETE("", controller.Disable2FA)                         // DELETE /api/auth/2fa
	twoFactor.POST("/backup-codes", controller.RegenerateBackupCodes)   // POST /api/auth/2fa/backup-codes
}

// setupPasskeyRoutes sets up passkey (WebAuthn) authentication routes
func setupPasskeyRoutes(api *echo.Group, protected *echo.Group, controller *controllers.WebAuthnController) {
	// WebAuthn機能が利用できない場合はルートを設定しない
	if controller == nil {
		return
	}

	passkey := api.Group("/auth/passkey")

	// パスキーログイン（認証不要）
	passkey.POST("/login/begin", controller.BeginLogin)   // POST /api/auth/passkey/login/begin
	passkey.POST("/login/finish", controller.FinishLogin) // POST /api/auth/passkey/login/finish

	// パスキー登録と管理（認証が必要）
	passkeyProtected := protected.Group("/auth/passkey")
	passkeyProtected.POST("/register/begin", controller.BeginRegistration)      // POST /api/auth/passkey/register/begin
	passkeyProtected.POST("/register/finish", controller.FinishRegistration)    // POST /api/auth/passkey/register/finish
	passkeyProtected.GET("/credentials", controller.ListCredentials)            // GET /api/auth/passkey/credentials
	passkeyProtected.DELETE("/credentials/:credential_id", controller.DeleteCredential) // DELETE /api/auth/passkey/credentials/:credential_id
	passkeyProtected.PUT("/credentials/:credential_id", controller.RenameCredential)    // PUT /api/auth/passkey/credentials/:credential_id
}

// setupFinancialDataRoutes sets up financial data management routes
func setupFinancialDataRoutes(api *echo.Group, controller *controllers.FinancialDataController) {
	financialData := api.Group("/financial-data")

	financialData.POST("", controller.CreateFinancialData)                        // POST /api/financial-data
	financialData.GET("", controller.GetFinancialData)                            // GET /api/financial-data
	financialData.PUT("/:user_id/profile", controller.UpdateFinancialProfile)     // PUT /api/financial-data/:user_id/profile
	financialData.PUT("/:user_id/retirement", controller.UpdateRetirementData)    // PUT /api/financial-data/:user_id/retirement
	financialData.PUT("/:user_id/emergency-fund", controller.UpdateEmergencyFund) // PUT /api/financial-data/:user_id/emergency-fund
	financialData.DELETE("/:user_id", controller.DeleteFinancialData)             // DELETE /api/financial-data/:user_id
}

// setupCalculationRoutes sets up calculation routes
func setupCalculationRoutes(api *echo.Group, controller *controllers.CalculationsController) {
	calculations := api.Group("/calculations")

	calculations.POST("/asset-projection", controller.CalculateAssetProjection)       // POST /api/calculations/asset-projection
	calculations.POST("/retirement", controller.CalculateRetirementProjection)        // POST /api/calculations/retirement
	calculations.POST("/emergency-fund", controller.CalculateEmergencyFundProjection) // POST /api/calculations/emergency-fund
	calculations.POST("/comprehensive", controller.CalculateComprehensiveProjection)  // POST /api/calculations/comprehensive
	calculations.POST("/goal-projection", controller.CalculateGoalProjection)         // POST /api/calculations/goal-projection
}

// setupGoalRoutes sets up goal management routes
func setupGoalRoutes(api *echo.Group, controller *controllers.GoalsController) {
	goals := api.Group("/goals")

	goals.POST("", controller.CreateGoal)                                // POST /api/goals
	goals.GET("", controller.GetGoals)                                   // GET /api/goals
	goals.GET("/:id", controller.GetGoal)                                // GET /api/goals/:id
	goals.PUT("/:id", controller.UpdateGoal)                             // PUT /api/goals/:id
	goals.PUT("/:id/progress", controller.UpdateGoalProgress)            // PUT /api/goals/:id/progress
	goals.DELETE("/:id", controller.DeleteGoal)                          // DELETE /api/goals/:id
	goals.GET("/:id/recommendations", controller.GetGoalRecommendations) // GET /api/goals/:id/recommendations
	goals.GET("/:id/feasibility", controller.AnalyzeGoalFeasibility)     // GET /api/goals/:id/feasibility
}

// setupReportRoutes sets up report generation routes
func setupReportRoutes(api *echo.Group, controller *controllers.ReportsController) {
	reports := api.Group("/reports")

	reports.POST("/financial-summary", controller.GenerateFinancialSummaryReport) // POST /api/reports/financial-summary
	reports.POST("/asset-projection", controller.GenerateAssetProjectionReport)   // POST /api/reports/asset-projection
	reports.POST("/goals-progress", controller.GenerateGoalsProgressReport)       // POST /api/reports/goals-progress
	reports.POST("/retirement-plan", controller.GenerateRetirementPlanReport)     // POST /api/reports/retirement-plan
	reports.POST("/comprehensive", controller.GenerateComprehensiveReport)        // POST /api/reports/comprehensive
	reports.POST("/export", controller.ExportReportToPDF)                         // POST /api/reports/export
	reports.GET("/pdf", controller.GetReportPDF)                                  // GET /api/reports/pdf
	reports.GET("/download/:token", controller.DownloadReport)                    // GET /api/reports/download/:token
}

// Handler functions (placeholder implementations)

// HealthCheckHandler handles health check requests
func HealthCheckHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":    "ok",
		"message":   "財務計画計算機 API サーバーが正常に動作しています",
		"timestamp": time.Now().Format(time.RFC3339),
		"version":   "1.0.0",
	})
}

// APIInfoHandler provides API information
func APIInfoHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{
		"message":     "財務計画計算機 API v1.0",
		"description": "将来の資産形成と老後の財務計画を可視化するアプリケーションのAPI",
		"docs":        "/swagger/index.html",
		"endpoints": map[string]any{
			"financial_data": map[string]any{
				"base":              "/api/financial-data",
				"create":            "POST /api/financial-data",
				"get":               "GET /api/financial-data?user_id={user_id}",
				"update_profile":    "PUT /api/financial-data/{user_id}/profile",
				"update_retirement": "PUT /api/financial-data/{user_id}/retirement",
				"update_emergency":  "PUT /api/financial-data/{user_id}/emergency-fund",
				"delete":            "DELETE /api/financial-data/{user_id}",
			},
			"calculations": map[string]any{
				"base":             "/api/calculations",
				"asset_projection": "POST /api/calculations/asset-projection",
				"retirement":       "POST /api/calculations/retirement",
				"emergency_fund":   "POST /api/calculations/emergency-fund",
				"comprehensive":    "POST /api/calculations/comprehensive",
				"goal_projection":  "POST /api/calculations/goal-projection",
			},
			"goals": map[string]any{
				"base":            "/api/goals",
				"create":          "POST /api/goals",
				"list":            "GET /api/goals?user_id={user_id}",
				"get":             "GET /api/goals/{id}?user_id={user_id}",
				"update":          "PUT /api/goals/{id}?user_id={user_id}",
				"update_progress": "PUT /api/goals/{id}/progress?user_id={user_id}",
				"delete":          "DELETE /api/goals/{id}?user_id={user_id}",
				"recommendations": "GET /api/goals/{id}/recommendations?user_id={user_id}",
				"feasibility":     "GET /api/goals/{id}/feasibility?user_id={user_id}",
			},
			"reports": map[string]any{
				"base":              "/api/reports",
				"financial_summary": "POST /api/reports/financial-summary",
				"asset_projection":  "POST /api/reports/asset-projection",
				"goals_progress":    "POST /api/reports/goals-progress",
				"retirement_plan":   "POST /api/reports/retirement-plan",
				"comprehensive":     "POST /api/reports/comprehensive",
				"export":            "POST /api/reports/export",
				"pdf":               "GET /api/reports/pdf?user_id={user_id}",
			},
			"health": "/health",
		},
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// RateLimitStatusHandler returns the current rate limit status for the caller's IP.
//
// GET /api/rate-limit/status
//
// Response:
//
//	{
//	  "limit":     50,
//	  "remaining": 47,
//	  "reset":     1739865600,
//	  "reset_at":  "2026-02-18T15:00:00Z"
//	}
func RateLimitStatusHandler(store *CustomRateLimiterStore) echo.HandlerFunc {
	return func(c echo.Context) error {
		identifier, _ := extractIdentifier(c)
		info := store.GetInfo(identifier)
		return c.JSON(http.StatusOK, info)
	}
}
