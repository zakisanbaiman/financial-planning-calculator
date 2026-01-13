package web

import (
	"time"

	"github.com/financial-planning-calculator/backend/application/usecases"
	"github.com/financial-planning-calculator/backend/config"
	"github.com/financial-planning-calculator/backend/domain/repositories"
	"github.com/financial-planning-calculator/backend/domain/services"
	"github.com/financial-planning-calculator/backend/infrastructure/web/controllers"
	"github.com/labstack/echo/v4"
)

// ServerDependencies holds all dependencies needed for the web server
type ServerDependencies struct {
	// Repositories
	UserRepo          repositories.UserRepository
	RefreshTokenRepo  repositories.RefreshTokenRepository
	FinancialPlanRepo repositories.FinancialPlanRepository
	GoalRepo          repositories.GoalRepository

	// Domain Services
	CalculationService    *services.FinancialCalculationService
	RecommendationService *services.GoalRecommendationService

	// Auth Config
	JWTSecret              string
	JWTExpiration          time.Duration
	RefreshTokenExpiration time.Duration

	// Server Config (for OAuth)
	ServerConfig *config.ServerConfig

	// AuthUseCase (ミドルウェア用、NewControllersで初期化される)
	AuthUseCase usecases.AuthUseCase

	// SkipAuth テスト用：認証をスキップする
	SkipAuth bool
}

// NewControllers creates all controller instances with their dependencies
func NewControllers(deps *ServerDependencies) *Controllers {
	// Create use cases
	authUseCase := usecases.NewAuthUseCase(
		deps.UserRepo,
		deps.RefreshTokenRepo,
		deps.JWTSecret,
		deps.JWTExpiration,
		deps.RefreshTokenExpiration,
	)

	// Store auth use case for middleware
	deps.AuthUseCase = authUseCase

	manageFinancialDataUseCase := usecases.NewManageFinancialDataUseCase(
		deps.FinancialPlanRepo,
	)

	manageGoalsUseCase := usecases.NewManageGoalsUseCase(
		deps.GoalRepo,
		deps.FinancialPlanRepo,
		deps.RecommendationService,
	)

	calculateProjectionUseCase := usecases.NewCalculateProjectionUseCase(
		deps.FinancialPlanRepo,
		deps.GoalRepo,
		deps.CalculationService,
		deps.RecommendationService,
	)

	generateReportsUseCase := usecases.NewGenerateReportsUseCase(
		deps.FinancialPlanRepo,
		deps.GoalRepo,
		deps.CalculationService,
		deps.RecommendationService,
	)

	// Create controllers
	return &Controllers{
		Auth:          controllers.NewAuthController(authUseCase),
		FinancialData: controllers.NewFinancialDataController(manageFinancialDataUseCase),
		Calculations:  controllers.NewCalculationsController(calculateProjectionUseCase),
		Goals:         controllers.NewGoalsController(manageGoalsUseCase),
		Reports:       controllers.NewReportsController(generateReportsUseCase),
	}
}

// JWTAuthMiddlewareFunc returns the JWT authentication middleware
// Returns nil if SkipAuth is true (for testing)
func (deps *ServerDependencies) JWTAuthMiddlewareFunc() echo.MiddlewareFunc {
	if deps.SkipAuth {
		return nil
	}
	return JWTAuthMiddleware(deps.AuthUseCase)
}
