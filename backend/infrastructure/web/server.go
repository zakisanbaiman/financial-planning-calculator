package web

import (
	"github.com/financial-planning-calculator/backend/application/usecases"
	"github.com/financial-planning-calculator/backend/domain/repositories"
	"github.com/financial-planning-calculator/backend/domain/services"
	"github.com/financial-planning-calculator/backend/infrastructure/web/controllers"
)

// ServerDependencies holds all dependencies needed for the web server
type ServerDependencies struct {
	// Repositories
	FinancialPlanRepo repositories.FinancialPlanRepository
	GoalRepo          repositories.GoalRepository

	// Domain Services
	CalculationService    *services.FinancialCalculationService
	RecommendationService *services.GoalRecommendationService
}

// NewControllers creates all controller instances with their dependencies
func NewControllers(deps *ServerDependencies) *Controllers {
	// Create use cases
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
		FinancialData: controllers.NewFinancialDataController(manageFinancialDataUseCase),
		Calculations:  controllers.NewCalculationsController(calculateProjectionUseCase),
		Goals:         controllers.NewGoalsController(manageGoalsUseCase),
		Reports:       controllers.NewReportsController(generateReportsUseCase),
	}
}
