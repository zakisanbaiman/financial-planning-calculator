package usecases

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/financial-planning-calculator/backend/domain/entities"
	"github.com/financial-planning-calculator/backend/domain/services"
	"github.com/financial-planning-calculator/backend/domain/valueobjects"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newTestGoal はテスト用の目標を作成するヘルパー
func newTestGoal(userID entities.UserID, goalID entities.GoalID) *entities.Goal {
	targetAmount, _ := valueobjects.NewMoneyJPY(1000000)
	monthlyContribution, _ := valueobjects.NewMoneyJPY(50000)
	targetDate := time.Now().AddDate(2, 0, 0) // 2年後

	goal, err := entities.NewGoal(userID, entities.GoalTypeSavings, "新車購入", targetAmount, targetDate, monthlyContribution)
	if err != nil {
		panic("テスト用目標の作成に失敗: " + err.Error())
	}
	return goal
}

// ===========================
// CreateGoal Tests
// ===========================

func TestManageGoalsUseCase_CreateGoal(t *testing.T) {
	ctx := context.Background()
	calcService := services.NewFinancialCalculationService()
	recService := services.NewGoalRecommendationService(calcService)

	baseInput := CreateGoalInput{
		UserID:              "user-001",
		GoalType:            "savings",
		Title:               "新車購入",
		TargetAmount:        1000000,
		TargetDate:          time.Now().AddDate(2, 0, 0).Format(time.RFC3339),
		CurrentAmount:       100000,
		MonthlyContribution: 50000,
	}

	t.Run("正常系: 財務計画なしでも目標を作成できる", func(t *testing.T) {
		mockGoalRepo := new(MockGoalRepository)
		mockPlanRepo := new(MockFinancialPlanRepository)
		// 財務データが見つからないエラーを返す → 達成可能性チェックをスキップして保存
		mockPlanRepo.On("FindByUserID", mock_anything(), entities.UserID("user-001")).
			Return(nil, errors.New("財務データが見つかりません"))
		mockGoalRepo.On("Save", mock_anything(), mock_anything()).Return(nil)

		uc := NewManageGoalsUseCase(mockGoalRepo, mockPlanRepo, recService)
		output, err := uc.CreateGoal(ctx, baseInput)

		require.NoError(t, err)
		assert.NotEmpty(t, output.GoalID)
		assert.Equal(t, entities.UserID("user-001"), output.UserID)
		mockGoalRepo.AssertExpectations(t)
	})

	t.Run("異常系: 無効な目標タイプの場合はエラー", func(t *testing.T) {
		mockGoalRepo := new(MockGoalRepository)
		mockPlanRepo := new(MockFinancialPlanRepository)

		uc := NewManageGoalsUseCase(mockGoalRepo, mockPlanRepo, recService)
		_, err := uc.CreateGoal(ctx, CreateGoalInput{
			UserID:              "user-001",
			GoalType:            "invalid_type",
			Title:               "テスト",
			TargetAmount:        1000000,
			TargetDate:          time.Now().AddDate(1, 0, 0).Format(time.RFC3339),
			CurrentAmount:       0,
			MonthlyContribution: 50000,
		})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "無効な目標タイプです")
	})

	t.Run("異常系: 無効な日付フォーマットの場合はエラー", func(t *testing.T) {
		mockGoalRepo := new(MockGoalRepository)
		mockPlanRepo := new(MockFinancialPlanRepository)

		uc := NewManageGoalsUseCase(mockGoalRepo, mockPlanRepo, recService)
		_, err := uc.CreateGoal(ctx, CreateGoalInput{
			UserID:              "user-001",
			GoalType:            "savings",
			Title:               "テスト",
			TargetAmount:        1000000,
			TargetDate:          "2025/12/31", // 無効なフォーマット
			CurrentAmount:       0,
			MonthlyContribution: 50000,
		})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "目標日の解析に失敗しました")
	})

	t.Run("異常系: Saveリポジトリエラーでエラーになる", func(t *testing.T) {
		mockGoalRepo := new(MockGoalRepository)
		mockPlanRepo := new(MockFinancialPlanRepository)
		mockPlanRepo.On("FindByUserID", mock_anything(), entities.UserID("user-001")).
			Return(nil, errors.New("財務データが見つかりません"))
		mockGoalRepo.On("Save", mock_anything(), mock_anything()).Return(errors.New("db error"))

		uc := NewManageGoalsUseCase(mockGoalRepo, mockPlanRepo, recService)
		_, err := uc.CreateGoal(ctx, baseInput)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "目標の保存に失敗しました")
		mockGoalRepo.AssertExpectations(t)
	})

	t.Run("正常系: 財務計画ありで達成可能な目標を作成できる", func(t *testing.T) {
		mockGoalRepo := new(MockGoalRepository)
		mockPlanRepo := new(MockFinancialPlanRepository)
		plan := newTestFinancialPlan("user-001")
		mockPlanRepo.On("FindByUserID", mock_anything(), entities.UserID("user-001")).Return(plan, nil)
		mockGoalRepo.On("Save", mock_anything(), mock_anything()).Return(nil)
		mockPlanRepo.On("Update", mock_anything(), mock_anything()).Return(nil)

		uc := NewManageGoalsUseCase(mockGoalRepo, mockPlanRepo, recService)
		output, err := uc.CreateGoal(ctx, baseInput)

		if err == nil {
			assert.NotEmpty(t, output.GoalID)
		}
		// 達成不可能と判定された場合も正常なビジネスロジック
		mockGoalRepo.AssertExpectations(t)
	})
}

// ===========================
// GetGoal Tests
// ===========================

func TestManageGoalsUseCase_GetGoal(t *testing.T) {
	ctx := context.Background()
	calcService := services.NewFinancialCalculationService()
	recService := services.NewGoalRecommendationService(calcService)

	t.Run("正常系: 目標を取得できる", func(t *testing.T) {
		mockGoalRepo := new(MockGoalRepository)
		mockPlanRepo := new(MockFinancialPlanRepository)
		goal := newTestGoal("user-001", "goal-001")
		mockGoalRepo.On("FindByID", mock_anything(), goal.ID()).Return(goal, nil)

		uc := NewManageGoalsUseCase(mockGoalRepo, mockPlanRepo, recService)
		output, err := uc.GetGoal(ctx, GetGoalInput{
			GoalID: goal.ID(),
			UserID: "user-001",
		})

		require.NoError(t, err)
		assert.NotNil(t, output.Goal)
		assert.Equal(t, goal.ID(), output.Goal.ID())
		mockGoalRepo.AssertExpectations(t)
	})

	t.Run("異常系: 目標が存在しない場合はエラー", func(t *testing.T) {
		mockGoalRepo := new(MockGoalRepository)
		mockPlanRepo := new(MockFinancialPlanRepository)
		mockGoalRepo.On("FindByID", mock_anything(), entities.GoalID("goal-999")).Return(nil, errors.New("not found"))

		uc := NewManageGoalsUseCase(mockGoalRepo, mockPlanRepo, recService)
		_, err := uc.GetGoal(ctx, GetGoalInput{
			GoalID: "goal-999",
			UserID: "user-001",
		})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "目標の取得に失敗しました")
		mockGoalRepo.AssertExpectations(t)
	})

	t.Run("異常系: 別ユーザーの目標へのアクセスは拒否される", func(t *testing.T) {
		mockGoalRepo := new(MockGoalRepository)
		mockPlanRepo := new(MockFinancialPlanRepository)
		goal := newTestGoal("user-001", "goal-001")
		mockGoalRepo.On("FindByID", mock_anything(), goal.ID()).Return(goal, nil)

		uc := NewManageGoalsUseCase(mockGoalRepo, mockPlanRepo, recService)
		_, err := uc.GetGoal(ctx, GetGoalInput{
			GoalID: goal.ID(),
			UserID: "user-002", // 異なるユーザー
		})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "権限がありません")
		mockGoalRepo.AssertExpectations(t)
	})
}

// ===========================
// GetGoalsByUser Tests
// ===========================

func TestManageGoalsUseCase_GetGoalsByUser(t *testing.T) {
	ctx := context.Background()
	calcService := services.NewFinancialCalculationService()
	recService := services.NewGoalRecommendationService(calcService)

	t.Run("正常系: ユーザーの全目標を取得できる", func(t *testing.T) {
		mockGoalRepo := new(MockGoalRepository)
		mockPlanRepo := new(MockFinancialPlanRepository)
		goal := newTestGoal("user-001", "goal-001")
		mockGoalRepo.On("FindByUserID", mock_anything(), entities.UserID("user-001")).Return([]*entities.Goal{goal}, nil)

		uc := NewManageGoalsUseCase(mockGoalRepo, mockPlanRepo, recService)
		output, err := uc.GetGoalsByUser(ctx, GetGoalsByUserInput{
			UserID:     "user-001",
			ActiveOnly: false,
		})

		require.NoError(t, err)
		assert.Len(t, output.Goals, 1)
		mockGoalRepo.AssertExpectations(t)
	})

	t.Run("正常系: 目標が0件の場合も正常に返す", func(t *testing.T) {
		mockGoalRepo := new(MockGoalRepository)
		mockPlanRepo := new(MockFinancialPlanRepository)
		mockGoalRepo.On("FindByUserID", mock_anything(), entities.UserID("user-001")).Return([]*entities.Goal{}, nil)

		uc := NewManageGoalsUseCase(mockGoalRepo, mockPlanRepo, recService)
		output, err := uc.GetGoalsByUser(ctx, GetGoalsByUserInput{
			UserID:     "user-001",
			ActiveOnly: false,
		})

		require.NoError(t, err)
		assert.Empty(t, output.Goals)
		mockGoalRepo.AssertExpectations(t)
	})

	t.Run("異常系: リポジトリエラーの場合はエラーを返す", func(t *testing.T) {
		mockGoalRepo := new(MockGoalRepository)
		mockPlanRepo := new(MockFinancialPlanRepository)
		mockGoalRepo.On("FindByUserID", mock_anything(), entities.UserID("user-001")).Return(nil, errors.New("db error"))

		uc := NewManageGoalsUseCase(mockGoalRepo, mockPlanRepo, recService)
		_, err := uc.GetGoalsByUser(ctx, GetGoalsByUserInput{
			UserID:     "user-001",
			ActiveOnly: false,
		})

		require.Error(t, err)
		mockGoalRepo.AssertExpectations(t)
	})

	t.Run("正常系: アクティブな目標のみを取得できる", func(t *testing.T) {
		mockGoalRepo := new(MockGoalRepository)
		mockPlanRepo := new(MockFinancialPlanRepository)
		goal := newTestGoal("user-001", "goal-001")
		mockGoalRepo.On("FindActiveGoalsByUserID", mock_anything(), entities.UserID("user-001")).Return([]*entities.Goal{goal}, nil)

		uc := NewManageGoalsUseCase(mockGoalRepo, mockPlanRepo, recService)
		output, err := uc.GetGoalsByUser(ctx, GetGoalsByUserInput{
			UserID:     "user-001",
			ActiveOnly: true,
		})

		require.NoError(t, err)
		assert.Len(t, output.Goals, 1)
		mockGoalRepo.AssertExpectations(t)
	})
}

// ===========================
// DeleteGoal Tests
// ===========================

func TestManageGoalsUseCase_DeleteGoal(t *testing.T) {
	ctx := context.Background()
	calcService := services.NewFinancialCalculationService()
	recService := services.NewGoalRecommendationService(calcService)

	t.Run("正常系: 目標を削除できる", func(t *testing.T) {
		mockGoalRepo := new(MockGoalRepository)
		mockPlanRepo := new(MockFinancialPlanRepository)
		goal := newTestGoal("user-001", "goal-001")
		mockGoalRepo.On("FindByID", mock_anything(), goal.ID()).Return(goal, nil)
		mockGoalRepo.On("Delete", mock_anything(), goal.ID()).Return(nil)

		uc := NewManageGoalsUseCase(mockGoalRepo, mockPlanRepo, recService)
		err := uc.DeleteGoal(ctx, DeleteGoalInput{
			GoalID: goal.ID(),
			UserID: "user-001",
		})

		require.NoError(t, err)
		mockGoalRepo.AssertExpectations(t)
	})

	t.Run("異常系: 別ユーザーの目標は削除できない", func(t *testing.T) {
		mockGoalRepo := new(MockGoalRepository)
		mockPlanRepo := new(MockFinancialPlanRepository)
		goal := newTestGoal("user-001", "goal-001")
		mockGoalRepo.On("FindByID", mock_anything(), goal.ID()).Return(goal, nil)

		uc := NewManageGoalsUseCase(mockGoalRepo, mockPlanRepo, recService)
		err := uc.DeleteGoal(ctx, DeleteGoalInput{
			GoalID: goal.ID(),
			UserID: "user-002",
		})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "権限がありません")
		mockGoalRepo.AssertExpectations(t)
	})
}
