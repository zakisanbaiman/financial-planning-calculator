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

func TestCalculateProjectionUseCase_CalculateAssetProjection(t *testing.T) {
	ctx := context.Background()
	calcService := services.NewFinancialCalculationService()
	recService := services.NewGoalRecommendationService()

	t.Run("正常系: 資産推移を計算できる", func(t *testing.T) {
		mockPlanRepo := new(MockFinancialPlanRepository)
		mockGoalRepo := new(MockGoalRepository)
		plan := newTestFinancialPlan("user-001")
		mockPlanRepo.On("FindByUserID", ctx, entities.UserID("user-001")).Return(plan, nil)

		uc := NewCalculateProjectionUseCase(mockPlanRepo, mockGoalRepo, calcService, recService)
		output, err := uc.CalculateAssetProjection(ctx, AssetProjectionInput{
			UserID: "user-001",
			Years:  10,
		})

		require.NoError(t, err)
		assert.NotNil(t, output)
		assert.Len(t, output.Projections, 10)
		assert.Greater(t, output.Summary.FinalAmount, output.Summary.InitialAmount)
		mockPlanRepo.AssertExpectations(t)
	})

	t.Run("異常系: 財務計画が存在しない場合はエラー", func(t *testing.T) {
		mockPlanRepo := new(MockFinancialPlanRepository)
		mockGoalRepo := new(MockGoalRepository)
		mockPlanRepo.On("FindByUserID", ctx, entities.UserID("user-999")).Return(nil, errors.New("not found"))

		uc := NewCalculateProjectionUseCase(mockPlanRepo, mockGoalRepo, calcService, recService)
		_, err := uc.CalculateAssetProjection(ctx, AssetProjectionInput{
			UserID: "user-999",
			Years:  10,
		})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "財務計画の取得に失敗しました")
		mockPlanRepo.AssertExpectations(t)
	})

	t.Run("異常系: Yearsが0の場合は空スライスを返す", func(t *testing.T) {
		mockPlanRepo := new(MockFinancialPlanRepository)
		mockGoalRepo := new(MockGoalRepository)
		plan := newTestFinancialPlan("user-001")
		mockPlanRepo.On("FindByUserID", ctx, entities.UserID("user-001")).Return(plan, nil)

		uc := NewCalculateProjectionUseCase(mockPlanRepo, mockGoalRepo, calcService, recService)
		output, err := uc.CalculateAssetProjection(ctx, AssetProjectionInput{
			UserID: "user-001",
			Years:  0,
		})

		// 0年の場合はエラーまたは空データになる
		if err == nil {
			assert.Empty(t, output.Projections)
		}
		mockPlanRepo.AssertExpectations(t)
	})
}

func TestCalculateProjectionUseCase_CalculateRetirementProjection(t *testing.T) {
	ctx := context.Background()
	calcService := services.NewFinancialCalculationService()
	recService := services.NewGoalRecommendationService()

	t.Run("異常系: FindByUserIDのエラーを伝播する", func(t *testing.T) {
		mockPlanRepo := new(MockFinancialPlanRepository)
		mockGoalRepo := new(MockGoalRepository)
		mockPlanRepo.On("FindByUserID", ctx, entities.UserID("user-999")).Return(nil, errors.New("db connection error"))

		uc := NewCalculateProjectionUseCase(mockPlanRepo, mockGoalRepo, calcService, recService)
		_, err := uc.CalculateRetirementProjection(ctx, RetirementProjectionInput{
			UserID: "user-999",
		})

		require.Error(t, err)
		mockPlanRepo.AssertExpectations(t)
	})
}

func TestCalculateProjectionUseCase_CalculateGoalProjection(t *testing.T) {
	ctx := context.Background()
	calcService := services.NewFinancialCalculationService()
	recService := services.NewGoalRecommendationService()

	t.Run("異常系: 目標が存在しない場合はエラー", func(t *testing.T) {
		mockPlanRepo := new(MockFinancialPlanRepository)
		mockGoalRepo := new(MockGoalRepository)
		mockGoalRepo.On("FindByID", ctx, entities.GoalID("goal-999")).Return(nil, errors.New("not found"))

		uc := NewCalculateProjectionUseCase(mockPlanRepo, mockGoalRepo, calcService, recService)
		_, err := uc.CalculateGoalProjection(ctx, GoalProjectionInput{
			UserID: "user-001",
			GoalID: "goal-999",
		})

		require.Error(t, err)
		mockGoalRepo.AssertExpectations(t)
	})
}

func TestCalculateProjectionUseCase_CalculateComprehensiveProjection(t *testing.T) {
	ctx := context.Background()
	calcService := services.NewFinancialCalculationService()
	recService := services.NewGoalRecommendationService()

	t.Run("異常系: FindByUserIDのエラーを伝播する", func(t *testing.T) {
		mockPlanRepo := new(MockFinancialPlanRepository)
		mockGoalRepo := new(MockGoalRepository)
		mockPlanRepo.On("FindByUserID", ctx, entities.UserID("user-999")).Return(nil, errors.New("not found"))

		uc := NewCalculateProjectionUseCase(mockPlanRepo, mockGoalRepo, calcService, recService)
		_, err := uc.CalculateComprehensiveProjection(ctx, ComprehensiveProjectionInput{
			UserID: "user-999",
			Years:  10,
		})

		require.Error(t, err)
		mockPlanRepo.AssertExpectations(t)
	})
}
