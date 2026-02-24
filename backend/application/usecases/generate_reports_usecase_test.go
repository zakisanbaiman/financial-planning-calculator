package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/financial-planning-calculator/backend/domain/entities"
	"github.com/financial-planning-calculator/backend/domain/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ===========================
// GenerateFinancialSummaryReport Tests
// ===========================

func TestGenerateReportsUseCase_GenerateFinancialSummaryReport(t *testing.T) {
	ctx := context.Background()
	calcService := services.NewFinancialCalculationService()
	recService := services.NewGoalRecommendationService(calcService)

	t.Run("正常系: 財務サマリーレポートを生成できる", func(t *testing.T) {
		mockPlanRepo := new(MockFinancialPlanRepository)
		mockGoalRepo := new(MockGoalRepository)
		plan := newTestFinancialPlan("user-001")
		mockPlanRepo.On("FindByUserID", mock_anything(), entities.UserID("user-001")).Return(plan, nil)

		uc := NewGenerateReportsUseCase(mockPlanRepo, mockGoalRepo, calcService, recService)
		output, err := uc.GenerateFinancialSummaryReport(ctx, FinancialSummaryReportInput{
			UserID: "user-001",
		})

		require.NoError(t, err)
		assert.NotNil(t, output)
		assert.NotEmpty(t, output.GeneratedAt)
		mockPlanRepo.AssertExpectations(t)
	})

	t.Run("異常系: 財務計画が存在しない場合はエラー", func(t *testing.T) {
		mockPlanRepo := new(MockFinancialPlanRepository)
		mockGoalRepo := new(MockGoalRepository)
		mockPlanRepo.On("FindByUserID", mock_anything(), entities.UserID("user-999")).Return(nil, errors.New("not found"))

		uc := NewGenerateReportsUseCase(mockPlanRepo, mockGoalRepo, calcService, recService)
		_, err := uc.GenerateFinancialSummaryReport(ctx, FinancialSummaryReportInput{
			UserID: "user-999",
		})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "財務計画の取得に失敗しました")
		mockPlanRepo.AssertExpectations(t)
	})
}

// ===========================
// GenerateAssetProjectionReport Tests
// ===========================

func TestGenerateReportsUseCase_GenerateAssetProjectionReport(t *testing.T) {
	ctx := context.Background()
	calcService := services.NewFinancialCalculationService()
	recService := services.NewGoalRecommendationService(calcService)

	t.Run("正常系: 資産推移レポートを生成できる", func(t *testing.T) {
		mockPlanRepo := new(MockFinancialPlanRepository)
		mockGoalRepo := new(MockGoalRepository)
		plan := newTestFinancialPlan("user-001")
		mockPlanRepo.On("FindByUserID", mock_anything(), entities.UserID("user-001")).Return(plan, nil)

		uc := NewGenerateReportsUseCase(mockPlanRepo, mockGoalRepo, calcService, recService)
		output, err := uc.GenerateAssetProjectionReport(ctx, AssetProjectionReportInput{
			UserID: "user-001",
			Years:  10,
		})

		require.NoError(t, err)
		assert.NotNil(t, output)
		mockPlanRepo.AssertExpectations(t)
	})

	t.Run("異常系: FindByUserIDのエラーを伝播する", func(t *testing.T) {
		mockPlanRepo := new(MockFinancialPlanRepository)
		mockGoalRepo := new(MockGoalRepository)
		mockPlanRepo.On("FindByUserID", mock_anything(), entities.UserID("user-999")).Return(nil, errors.New("db error"))

		uc := NewGenerateReportsUseCase(mockPlanRepo, mockGoalRepo, calcService, recService)
		_, err := uc.GenerateAssetProjectionReport(ctx, AssetProjectionReportInput{
			UserID: "user-999",
			Years:  10,
		})

		require.Error(t, err)
		mockPlanRepo.AssertExpectations(t)
	})
}

// ===========================
// GenerateGoalsProgressReport Tests
// ===========================

func TestGenerateReportsUseCase_GenerateGoalsProgressReport(t *testing.T) {
	ctx := context.Background()
	calcService := services.NewFinancialCalculationService()
	recService := services.NewGoalRecommendationService(calcService)

	t.Run("正常系: 目標進捗レポートを生成できる", func(t *testing.T) {
		mockPlanRepo := new(MockFinancialPlanRepository)
		mockGoalRepo := new(MockGoalRepository)
		goal := newTestGoal("user-001", "goal-001")
		mockGoalRepo.On("FindByUserID", mock_anything(), entities.UserID("user-001")).Return([]*entities.Goal{goal}, nil)

		uc := NewGenerateReportsUseCase(mockPlanRepo, mockGoalRepo, calcService, recService)
		output, err := uc.GenerateGoalsProgressReport(ctx, GoalsProgressReportInput{
			UserID: "user-001",
		})

		require.NoError(t, err)
		assert.NotNil(t, output)
		mockGoalRepo.AssertExpectations(t)
	})

	t.Run("異常系: FindByUserIDのエラーを伝播する", func(t *testing.T) {
		mockPlanRepo := new(MockFinancialPlanRepository)
		mockGoalRepo := new(MockGoalRepository)
		mockGoalRepo.On("FindByUserID", mock_anything(), entities.UserID("user-999")).Return(nil, errors.New("db error"))

		uc := NewGenerateReportsUseCase(mockPlanRepo, mockGoalRepo, calcService, recService)
		_, err := uc.GenerateGoalsProgressReport(ctx, GoalsProgressReportInput{
			UserID: "user-999",
		})

		require.Error(t, err)
		mockGoalRepo.AssertExpectations(t)
	})
}
