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
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateGoal(t *testing.T) {
	futureDate := time.Now().AddDate(2, 0, 0).Format(time.RFC3339)

	tests := []struct {
		name        string
		input       CreateGoalInput
		setupMocks  func(*MockGoalRepository, *MockFinancialPlanRepository)
		expectError bool
		errContains string
	}{
		{
			name: "正常系: 目標作成（財務プロファイルなし）",
			input: CreateGoalInput{
				UserID:              "user-001",
				GoalType:            "savings",
				Title:               "新車購入",
				TargetAmount:        1000000,
				TargetDate:          futureDate,
				CurrentAmount:       100000,
				MonthlyContribution: 50000,
			},
			setupMocks: func(gr *MockGoalRepository, fp *MockFinancialPlanRepository) {
				gr.On("Save", mock.Anything, mock.AnythingOfType("*entities.Goal")).Return(nil)
				fp.On("FindByUserID", mock.Anything, entities.UserID("user-001")).
					Return(nil, errors.New("財務データが見つかりません"))
			},
			expectError: false,
		},
		{
			name: "正常系: 目標作成（財務プロファイルあり）",
			input: CreateGoalInput{
				UserID:              "user-002",
				GoalType:            "savings",
				Title:               "旅行資金",
				TargetAmount:        500000,
				TargetDate:          futureDate,
				CurrentAmount:       50000,
				MonthlyContribution: 30000,
			},
			setupMocks: func(gr *MockGoalRepository, fp *MockFinancialPlanRepository) {
				plan := newTestFinancialPlan("user-002")
				fp.On("FindByUserID", mock.Anything, entities.UserID("user-002")).Return(plan, nil)
				gr.On("Save", mock.Anything, mock.AnythingOfType("*entities.Goal")).Return(nil)
				fp.On("Update", mock.Anything, mock.AnythingOfType("*aggregates.FinancialPlan")).Return(nil)
			},
			expectError: false,
		},
		{
			name: "異常系: 無効な目標タイプ",
			input: CreateGoalInput{
				UserID:              "user-003",
				GoalType:            "invalid",
				Title:               "目標",
				TargetAmount:        1000000,
				TargetDate:          futureDate,
				CurrentAmount:       0,
				MonthlyContribution: 50000,
			},
			setupMocks:  func(gr *MockGoalRepository, fp *MockFinancialPlanRepository) {},
			expectError: true,
			errContains: "無効な目標タイプです",
		},
		{
			name: "異常系: 無効な目標日フォーマット",
			input: CreateGoalInput{
				UserID:              "user-004",
				GoalType:            "savings",
				Title:               "目標",
				TargetAmount:        1000000,
				TargetDate:          "invalid-date",
				CurrentAmount:       0,
				MonthlyContribution: 50000,
			},
			setupMocks:  func(gr *MockGoalRepository, fp *MockFinancialPlanRepository) {},
			expectError: true,
			errContains: "目標日の解析に失敗しました",
		},
		{
			name: "異常系: 退職目標が既に存在する",
			input: CreateGoalInput{
				UserID:              "user-005",
				GoalType:            "retirement",
				Title:               "退職資金",
				TargetAmount:        50000000,
				TargetDate:          futureDate,
				CurrentAmount:       1000000,
				MonthlyContribution: 100000,
			},
			setupMocks: func(gr *MockGoalRepository, fp *MockFinancialPlanRepository) {
				existingGoal := newTestGoal("user-005", entities.GoalTypeRetirement)
				gr.On("FindByUserIDAndType", mock.Anything, entities.UserID("user-005"), entities.GoalTypeRetirement).
					Return([]*entities.Goal{existingGoal}, nil)
			},
			expectError: true,
			errContains: "退職・老後資金目標の目標は既に存在します",
		},
		{
			name: "異常系: 目標保存失敗",
			input: CreateGoalInput{
				UserID:              "user-006",
				GoalType:            "custom",
				Title:               "カスタム目標",
				TargetAmount:        200000,
				TargetDate:          futureDate,
				CurrentAmount:       0,
				MonthlyContribution: 10000,
			},
			setupMocks: func(gr *MockGoalRepository, fp *MockFinancialPlanRepository) {
				fp.On("FindByUserID", mock.Anything, entities.UserID("user-006")).
					Return(nil, errors.New("財務データが見つかりません"))
				gr.On("Save", mock.Anything, mock.AnythingOfType("*entities.Goal")).Return(errors.New("保存エラー"))
			},
			expectError: true,
			errContains: "目標の保存に失敗しました",
		},
	}

	calcService := services.NewFinancialCalculationService()
	recommendService := services.NewGoalRecommendationService(calcService)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goalRepo := new(MockGoalRepository)
			planRepo := new(MockFinancialPlanRepository)
			tt.setupMocks(goalRepo, planRepo)

			uc := NewManageGoalsUseCase(goalRepo, planRepo, recommendService)
			output, err := uc.CreateGoal(context.Background(), tt.input)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, output)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, output)
				assert.NotEmpty(t, output.GoalID)
			}
			goalRepo.AssertExpectations(t)
			planRepo.AssertExpectations(t)
		})
	}
}

func TestGetGoal(t *testing.T) {
	calcService := services.NewFinancialCalculationService()
	recommendService := services.NewGoalRecommendationService(calcService)

	tests := []struct {
		name        string
		input       GetGoalInput
		setupMocks  func(*MockGoalRepository, *MockFinancialPlanRepository)
		expectError bool
		errContains string
	}{
		{
			name: "正常系: 目標取得",
			input: GetGoalInput{
				GoalID: "goal-001",
				UserID: "user-001",
			},
			setupMocks: func(gr *MockGoalRepository, fp *MockFinancialPlanRepository) {
				goal := newTestGoal("user-001", entities.GoalTypeSavings)
				gr.On("FindByID", mock.Anything, entities.GoalID("goal-001")).Return(goal, nil)
			},
			expectError: false,
		},
		{
			name: "異常系: 目標が存在しない",
			input: GetGoalInput{
				GoalID: "goal-999",
				UserID: "user-001",
			},
			setupMocks: func(gr *MockGoalRepository, fp *MockFinancialPlanRepository) {
				gr.On("FindByID", mock.Anything, entities.GoalID("goal-999")).
					Return(nil, errors.New("目標が見つかりません"))
			},
			expectError: true,
			errContains: "目標の取得に失敗しました",
		},
		{
			name: "異常系: 別ユーザーの目標へのアクセス",
			input: GetGoalInput{
				GoalID: "goal-001",
				UserID: "other-user",
			},
			setupMocks: func(gr *MockGoalRepository, fp *MockFinancialPlanRepository) {
				goal := newTestGoal("user-001", entities.GoalTypeSavings)
				gr.On("FindByID", mock.Anything, entities.GoalID("goal-001")).Return(goal, nil)
			},
			expectError: true,
			errContains: "指定された目標にアクセスする権限がありません",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goalRepo := new(MockGoalRepository)
			planRepo := new(MockFinancialPlanRepository)
			tt.setupMocks(goalRepo, planRepo)

			uc := NewManageGoalsUseCase(goalRepo, planRepo, recommendService)
			output, err := uc.GetGoal(context.Background(), tt.input)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, output)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, output)
				assert.NotNil(t, output.Goal)
			}
			goalRepo.AssertExpectations(t)
		})
	}
}

func TestGetGoalsByUser(t *testing.T) {
	calcService := services.NewFinancialCalculationService()
	recommendService := services.NewGoalRecommendationService(calcService)

	tests := []struct {
		name        string
		input       GetGoalsByUserInput
		setupMocks  func(*MockGoalRepository, *MockFinancialPlanRepository)
		expectError bool
		expectCount int
	}{
		{
			name: "正常系: 全目標取得",
			input: GetGoalsByUserInput{
				UserID:     "user-001",
				ActiveOnly: false,
			},
			setupMocks: func(gr *MockGoalRepository, fp *MockFinancialPlanRepository) {
				goals := []*entities.Goal{
					newTestGoal("user-001", entities.GoalTypeSavings),
					newTestGoal("user-001", entities.GoalTypeCustom),
				}
				gr.On("FindByUserID", mock.Anything, entities.UserID("user-001")).Return(goals, nil)
			},
			expectError: false,
			expectCount: 2,
		},
		{
			name: "正常系: アクティブ目標のみ取得",
			input: GetGoalsByUserInput{
				UserID:     "user-001",
				ActiveOnly: true,
			},
			setupMocks: func(gr *MockGoalRepository, fp *MockFinancialPlanRepository) {
				goals := []*entities.Goal{newTestGoal("user-001", entities.GoalTypeSavings)}
				gr.On("FindActiveGoalsByUserID", mock.Anything, entities.UserID("user-001")).Return(goals, nil)
			},
			expectError: false,
			expectCount: 1,
		},
		{
			name: "正常系: 目標タイプ指定取得",
			input: GetGoalsByUserInput{
				UserID:   "user-001",
				GoalType: func() *entities.GoalType { t := entities.GoalTypeSavings; return &t }(),
			},
			setupMocks: func(gr *MockGoalRepository, fp *MockFinancialPlanRepository) {
				goals := []*entities.Goal{newTestGoal("user-001", entities.GoalTypeSavings)}
				gr.On("FindByUserIDAndType", mock.Anything, entities.UserID("user-001"), entities.GoalTypeSavings).
					Return(goals, nil)
			},
			expectError: false,
			expectCount: 1,
		},
		{
			name: "正常系: 目標が存在しない",
			input: GetGoalsByUserInput{
				UserID:     "user-002",
				ActiveOnly: false,
			},
			setupMocks: func(gr *MockGoalRepository, fp *MockFinancialPlanRepository) {
				gr.On("FindByUserID", mock.Anything, entities.UserID("user-002")).Return([]*entities.Goal{}, nil)
			},
			expectError: false,
			expectCount: 0,
		},
		{
			name: "異常系: リポジトリエラー",
			input: GetGoalsByUserInput{
				UserID:     "user-999",
				ActiveOnly: false,
			},
			setupMocks: func(gr *MockGoalRepository, fp *MockFinancialPlanRepository) {
				gr.On("FindByUserID", mock.Anything, entities.UserID("user-999")).
					Return(nil, errors.New("DBエラー"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goalRepo := new(MockGoalRepository)
			planRepo := new(MockFinancialPlanRepository)
			tt.setupMocks(goalRepo, planRepo)

			uc := NewManageGoalsUseCase(goalRepo, planRepo, recommendService)
			output, err := uc.GetGoalsByUser(context.Background(), tt.input)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, output)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, tt.expectCount, len(output.Goals))
			}
			goalRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateGoal(t *testing.T) {
	calcService := services.NewFinancialCalculationService()
	recommendService := services.NewGoalRecommendationService(calcService)
	futureDate := time.Now().AddDate(3, 0, 0).Format(time.RFC3339)

	tests := []struct {
		name        string
		input       UpdateGoalInput
		setupMocks  func(*MockGoalRepository, *MockFinancialPlanRepository)
		expectError bool
		errContains string
	}{
		{
			name: "正常系: タイトル更新",
			input: UpdateGoalInput{
				GoalID: "goal-001",
				UserID: "user-001",
				Title:  func() *string { s := "新しいタイトル"; return &s }(),
			},
			setupMocks: func(gr *MockGoalRepository, fp *MockFinancialPlanRepository) {
				goal := newTestGoal("user-001", entities.GoalTypeSavings)
				gr.On("FindByID", mock.Anything, entities.GoalID("goal-001")).Return(goal, nil)
				gr.On("Update", mock.Anything, mock.AnythingOfType("*entities.Goal")).Return(nil)
			},
			expectError: false,
		},
		{
			name: "正常系: 目標金額と目標日更新",
			input: UpdateGoalInput{
				GoalID:       "goal-001",
				UserID:       "user-001",
				TargetAmount: func() *float64 { v := 2000000.0; return &v }(),
				TargetDate:   &futureDate,
			},
			setupMocks: func(gr *MockGoalRepository, fp *MockFinancialPlanRepository) {
				goal := newTestGoal("user-001", entities.GoalTypeSavings)
				gr.On("FindByID", mock.Anything, entities.GoalID("goal-001")).Return(goal, nil)
				gr.On("Update", mock.Anything, mock.AnythingOfType("*entities.Goal")).Return(nil)
			},
			expectError: false,
		},
		{
			name: "異常系: 目標が存在しない",
			input: UpdateGoalInput{
				GoalID: "goal-999",
				UserID: "user-001",
			},
			setupMocks: func(gr *MockGoalRepository, fp *MockFinancialPlanRepository) {
				gr.On("FindByID", mock.Anything, entities.GoalID("goal-999")).
					Return(nil, errors.New("目標が見つかりません"))
			},
			expectError: true,
			errContains: "目標の取得に失敗しました",
		},
		{
			name: "異常系: 別ユーザーの目標更新",
			input: UpdateGoalInput{
				GoalID: "goal-001",
				UserID: "other-user",
				Title:  func() *string { s := "タイトル"; return &s }(),
			},
			setupMocks: func(gr *MockGoalRepository, fp *MockFinancialPlanRepository) {
				goal := newTestGoal("user-001", entities.GoalTypeSavings)
				gr.On("FindByID", mock.Anything, entities.GoalID("goal-001")).Return(goal, nil)
			},
			expectError: true,
			errContains: "指定された目標にアクセスする権限がありません",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goalRepo := new(MockGoalRepository)
			planRepo := new(MockFinancialPlanRepository)
			tt.setupMocks(goalRepo, planRepo)

			uc := NewManageGoalsUseCase(goalRepo, planRepo, recommendService)
			output, err := uc.UpdateGoal(context.Background(), tt.input)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, output)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, output)
				assert.True(t, output.Success)
			}
			goalRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateGoalProgress(t *testing.T) {
	calcService := services.NewFinancialCalculationService()
	recommendService := services.NewGoalRecommendationService(calcService)

	tests := []struct {
		name        string
		input       UpdateGoalProgressInput
		setupMocks  func(*MockGoalRepository, *MockFinancialPlanRepository)
		expectError bool
		errContains string
	}{
		{
			name: "正常系: 進捗更新",
			input: UpdateGoalProgressInput{
				GoalID:        "goal-001",
				UserID:        "user-001",
				CurrentAmount: 500000,
			},
			setupMocks: func(gr *MockGoalRepository, fp *MockFinancialPlanRepository) {
				goal := newTestGoal("user-001", entities.GoalTypeSavings)
				gr.On("FindByID", mock.Anything, entities.GoalID("goal-001")).Return(goal, nil)
				gr.On("Update", mock.Anything, mock.AnythingOfType("*entities.Goal")).Return(nil)
			},
			expectError: false,
		},
		{
			name: "異常系: 別ユーザーの目標進捗更新",
			input: UpdateGoalProgressInput{
				GoalID:        "goal-001",
				UserID:        "other-user",
				CurrentAmount: 500000,
			},
			setupMocks: func(gr *MockGoalRepository, fp *MockFinancialPlanRepository) {
				goal := newTestGoal("user-001", entities.GoalTypeSavings)
				gr.On("FindByID", mock.Anything, entities.GoalID("goal-001")).Return(goal, nil)
			},
			expectError: true,
			errContains: "指定された目標にアクセスする権限がありません",
		},
		{
			name: "異常系: 負の金額",
			input: UpdateGoalProgressInput{
				GoalID:        "goal-001",
				UserID:        "user-001",
				CurrentAmount: -100000,
			},
			setupMocks: func(gr *MockGoalRepository, fp *MockFinancialPlanRepository) {
				goal := newTestGoal("user-001", entities.GoalTypeSavings)
				gr.On("FindByID", mock.Anything, entities.GoalID("goal-001")).Return(goal, nil)
			},
			expectError: true,
			errContains: "現在金額の更新に失敗しました",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goalRepo := new(MockGoalRepository)
			planRepo := new(MockFinancialPlanRepository)
			tt.setupMocks(goalRepo, planRepo)

			uc := NewManageGoalsUseCase(goalRepo, planRepo, recommendService)
			output, err := uc.UpdateGoalProgress(context.Background(), tt.input)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, output)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, output)
				assert.True(t, output.Success)
			}
			goalRepo.AssertExpectations(t)
		})
	}
}

func TestDeleteGoal(t *testing.T) {
	calcService := services.NewFinancialCalculationService()
	recommendService := services.NewGoalRecommendationService(calcService)

	// 固定のゴールIDを使ったシナリオ（削除成功）
	t.Run("正常系: 目標削除", func(t *testing.T) {
		goal := newTestGoal("user-001", entities.GoalTypeSavings)
		plan := newTestFinancialPlanWithGoal("user-001", goal)
		goalID := goal.ID()

		goalRepo := new(MockGoalRepository)
		planRepo := new(MockFinancialPlanRepository)
		goalRepo.On("FindByID", mock.Anything, goalID).Return(goal, nil)
		planRepo.On("FindByUserID", mock.Anything, entities.UserID("user-001")).Return(plan, nil)
		planRepo.On("Update", mock.Anything, mock.AnythingOfType("*aggregates.FinancialPlan")).Return(nil)
		goalRepo.On("Delete", mock.Anything, goalID).Return(nil)

		uc := NewManageGoalsUseCase(goalRepo, planRepo, recommendService)
		err := uc.DeleteGoal(context.Background(), DeleteGoalInput{GoalID: goalID, UserID: "user-001"})

		require.NoError(t, err)
		goalRepo.AssertExpectations(t)
		planRepo.AssertExpectations(t)
	})

	tests := []struct {
		name        string
		input       DeleteGoalInput
		setupMocks  func(*MockGoalRepository, *MockFinancialPlanRepository)
		expectError bool
		errContains string
	}{
		{
			name: "異常系: 目標が存在しない",
			input: DeleteGoalInput{
				GoalID: "goal-999",
				UserID: "user-001",
			},
			setupMocks: func(gr *MockGoalRepository, fp *MockFinancialPlanRepository) {
				gr.On("FindByID", mock.Anything, entities.GoalID("goal-999")).
					Return(nil, errors.New("目標が見つかりません"))
			},
			expectError: true,
			errContains: "目標の取得に失敗しました",
		},
		{
			name: "異常系: 別ユーザーの目標削除",
			input: DeleteGoalInput{
				GoalID: "goal-001",
				UserID: "other-user",
			},
			setupMocks: func(gr *MockGoalRepository, fp *MockFinancialPlanRepository) {
				goal := newTestGoal("user-001", entities.GoalTypeSavings)
				gr.On("FindByID", mock.Anything, entities.GoalID("goal-001")).Return(goal, nil)
			},
			expectError: true,
			errContains: "指定された目標にアクセスする権限がありません",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goalRepo := new(MockGoalRepository)
			planRepo := new(MockFinancialPlanRepository)
			tt.setupMocks(goalRepo, planRepo)

			uc := NewManageGoalsUseCase(goalRepo, planRepo, recommendService)
			err := uc.DeleteGoal(context.Background(), tt.input)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
			}
			goalRepo.AssertExpectations(t)
			planRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateGoalProgress_Completed(t *testing.T) {
	calcService := services.NewFinancialCalculationService()
	recommendService := services.NewGoalRecommendationService(calcService)

	goalRepo := new(MockGoalRepository)
	planRepo := new(MockFinancialPlanRepository)

	// 目標金額と同額を入力 → 完了状態になるはず
	goal := newTestGoal("user-001", entities.GoalTypeSavings)
	targetAmount, _ := valueobjects.NewMoneyJPY(1000000)
	_ = goal.UpdateCurrentAmount(targetAmount)

	goalRepo.On("FindByID", mock.Anything, entities.GoalID("goal-001")).Return(goal, nil)
	goalRepo.On("Update", mock.Anything, mock.AnythingOfType("*entities.Goal")).Return(nil)

	uc := NewManageGoalsUseCase(goalRepo, planRepo, recommendService)
	output, err := uc.UpdateGoalProgress(context.Background(), UpdateGoalProgressInput{
		GoalID:        "goal-001",
		UserID:        "user-001",
		CurrentAmount: 1000000,
	})

	require.NoError(t, err)
	assert.True(t, output.Success)
	assert.True(t, output.IsCompleted)
	goalRepo.AssertExpectations(t)
}

func TestGetGoalRecommendations(t *testing.T) {
	calcService := services.NewFinancialCalculationService()
	recommendService := services.NewGoalRecommendationService(calcService)

	tests := []struct {
		name        string
		input       GetGoalRecommendationsInput
		setupMocks  func(*MockGoalRepository, *MockFinancialPlanRepository)
		expectError bool
		errContains string
	}{
		{
			name: "正常系: 目標推奨事項取得",
			input: GetGoalRecommendationsInput{
				GoalID: "goal-001",
				UserID: "user-001",
			},
			setupMocks: func(gr *MockGoalRepository, fp *MockFinancialPlanRepository) {
				goal := newTestGoal("user-001", entities.GoalTypeSavings)
				plan := newTestFinancialPlan("user-001")
				gr.On("FindByID", mock.Anything, entities.GoalID("goal-001")).Return(goal, nil)
				fp.On("FindByUserID", mock.Anything, entities.UserID("user-001")).Return(plan, nil)
			},
			expectError: false,
		},
		{
			name: "異常系: 目標が存在しない",
			input: GetGoalRecommendationsInput{
				GoalID: "goal-999",
				UserID: "user-001",
			},
			setupMocks: func(gr *MockGoalRepository, fp *MockFinancialPlanRepository) {
				gr.On("FindByID", mock.Anything, entities.GoalID("goal-999")).
					Return(nil, errors.New("目標が見つかりません"))
			},
			expectError: true,
			errContains: "目標の取得に失敗しました",
		},
		{
			name: "異常系: 別ユーザーのアクセス",
			input: GetGoalRecommendationsInput{
				GoalID: "goal-001",
				UserID: "other-user",
			},
			setupMocks: func(gr *MockGoalRepository, fp *MockFinancialPlanRepository) {
				goal := newTestGoal("user-001", entities.GoalTypeSavings)
				gr.On("FindByID", mock.Anything, entities.GoalID("goal-001")).Return(goal, nil)
			},
			expectError: true,
			errContains: "指定された目標にアクセスする権限がありません",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goalRepo := new(MockGoalRepository)
			planRepo := new(MockFinancialPlanRepository)
			tt.setupMocks(goalRepo, planRepo)

			uc := NewManageGoalsUseCase(goalRepo, planRepo, recommendService)
			output, err := uc.GetGoalRecommendations(context.Background(), tt.input)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, output)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, output)
			}
			goalRepo.AssertExpectations(t)
			planRepo.AssertExpectations(t)
		})
	}
}

func TestAnalyzeGoalFeasibility(t *testing.T) {
	calcService := services.NewFinancialCalculationService()
	recommendService := services.NewGoalRecommendationService(calcService)

	tests := []struct {
		name        string
		input       AnalyzeGoalFeasibilityInput
		setupMocks  func(*MockGoalRepository, *MockFinancialPlanRepository)
		expectError bool
		errContains string
	}{
		{
			name: "正常系: 目標実現可能性分析",
			input: AnalyzeGoalFeasibilityInput{
				GoalID: "goal-001",
				UserID: "user-001",
			},
			setupMocks: func(gr *MockGoalRepository, fp *MockFinancialPlanRepository) {
				goal := newTestGoal("user-001", entities.GoalTypeSavings)
				plan := newTestFinancialPlan("user-001")
				gr.On("FindByID", mock.Anything, entities.GoalID("goal-001")).Return(goal, nil)
				fp.On("FindByUserID", mock.Anything, entities.UserID("user-001")).Return(plan, nil)
			},
			expectError: false,
		},
		{
			name: "異常系: 目標が存在しない",
			input: AnalyzeGoalFeasibilityInput{
				GoalID: "goal-999",
				UserID: "user-001",
			},
			setupMocks: func(gr *MockGoalRepository, fp *MockFinancialPlanRepository) {
				gr.On("FindByID", mock.Anything, entities.GoalID("goal-999")).
					Return(nil, errors.New("目標が見つかりません"))
			},
			expectError: true,
			errContains: "目標の取得に失敗しました",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goalRepo := new(MockGoalRepository)
			planRepo := new(MockFinancialPlanRepository)
			tt.setupMocks(goalRepo, planRepo)

			uc := NewManageGoalsUseCase(goalRepo, planRepo, recommendService)
			output, err := uc.AnalyzeGoalFeasibility(context.Background(), tt.input)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, output)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, output)
			}
			goalRepo.AssertExpectations(t)
			planRepo.AssertExpectations(t)
		})
	}
}
