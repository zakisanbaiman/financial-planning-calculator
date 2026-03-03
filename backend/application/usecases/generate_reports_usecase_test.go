package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/financial-planning-calculator/backend/domain/aggregates"
	"github.com/financial-planning-calculator/backend/domain/entities"
	"github.com/financial-planning-calculator/backend/domain/services"
	"github.com/financial-planning-calculator/backend/domain/valueobjects"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newTestFinancialPlanWithRetirementData は退職データ付きテスト用財務計画を作成するヘルパー
func newTestFinancialPlanWithRetirementData(userID entities.UserID) *aggregates.FinancialPlan {
	plan := newTestFinancialPlan(userID)
	monthlyExpenses, _ := valueobjects.NewMoneyJPY(200000)
	pension, _ := valueobjects.NewMoneyJPY(80000)
	retirement, _ := entities.NewRetirementData(userID, 40, 65, 85, monthlyExpenses, pension)
	_ = plan.SetRetirementData(retirement)
	return plan
}

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
		plan := newTestFinancialPlan("user-001")
		goal := newTestGoal("user-001", "goal-001")
		mockGoalRepo.On("FindByUserID", mock_anything(), entities.UserID("user-001")).Return([]*entities.Goal{goal}, nil)
		mockPlanRepo.On("FindByUserID", mock_anything(), entities.UserID("user-001")).Return(plan, nil)

		uc := NewGenerateReportsUseCase(mockPlanRepo, mockGoalRepo, calcService, recService)
		output, err := uc.GenerateGoalsProgressReport(ctx, GoalsProgressReportInput{
			UserID: "user-001",
		})

		require.NoError(t, err)
		assert.NotNil(t, output)
		mockGoalRepo.AssertExpectations(t)
		mockPlanRepo.AssertExpectations(t)
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
// ===========================
// GenerateRetirementPlanReport Tests
// ===========================

func TestGenerateReportsUseCase_GenerateRetirementPlanReport(t *testing.T) {
	ctx := context.Background()
	calcService := services.NewFinancialCalculationService()
	recService := services.NewGoalRecommendationService(calcService)

	t.Run("正常系: 退職計画レポートを生成できる", func(t *testing.T) {
		mockPlanRepo := new(MockFinancialPlanRepository)
		mockGoalRepo := new(MockGoalRepository)
		plan := newTestFinancialPlanWithRetirementData("user-001")
		mockPlanRepo.On("FindByUserID", mock_anything(), entities.UserID("user-001")).Return(plan, nil)

		uc := NewGenerateReportsUseCase(mockPlanRepo, mockGoalRepo, calcService, recService)
		output, err := uc.GenerateRetirementPlanReport(ctx, RetirementPlanReportInput{
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
		_, err := uc.GenerateRetirementPlanReport(ctx, RetirementPlanReportInput{
			UserID: "user-999",
		})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "財務計画の取得に失敗しました")
		mockPlanRepo.AssertExpectations(t)
	})

	t.Run("異常系: 退職データが設定されていない場合はエラー", func(t *testing.T) {
		mockPlanRepo := new(MockFinancialPlanRepository)
		mockGoalRepo := new(MockGoalRepository)
		plan := newTestFinancialPlan("user-001") // 退職データなし
		mockPlanRepo.On("FindByUserID", mock_anything(), entities.UserID("user-001")).Return(plan, nil)

		uc := NewGenerateReportsUseCase(mockPlanRepo, mockGoalRepo, calcService, recService)
		_, err := uc.GenerateRetirementPlanReport(ctx, RetirementPlanReportInput{
			UserID: "user-001",
		})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "退職データが設定されていません")
		mockPlanRepo.AssertExpectations(t)
	})
}

// ===========================
// GenerateComprehensiveReport Tests
// ===========================

func TestGenerateReportsUseCase_GenerateComprehensiveReport(t *testing.T) {
	ctx := context.Background()
	calcService := services.NewFinancialCalculationService()
	recService := services.NewGoalRecommendationService(calcService)

	t.Run("正常系: 包括的レポートを生成できる", func(t *testing.T) {
		mockPlanRepo := new(MockFinancialPlanRepository)
		mockGoalRepo := new(MockGoalRepository)
		plan := newTestFinancialPlan("user-001")
		mockPlanRepo.On("FindByUserID", mock_anything(), entities.UserID("user-001")).Return(plan, nil)
		mockGoalRepo.On("FindByUserID", mock_anything(), entities.UserID("user-001")).Return(nil, nil)

		uc := NewGenerateReportsUseCase(mockPlanRepo, mockGoalRepo, calcService, recService)
		output, err := uc.GenerateComprehensiveReport(ctx, ComprehensiveReportInput{
			UserID: "user-001",
			Years:  10,
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
		_, err := uc.GenerateComprehensiveReport(ctx, ComprehensiveReportInput{
			UserID: "user-999",
			Years:  10,
		})

		require.Error(t, err)
		mockPlanRepo.AssertExpectations(t)
	})
}