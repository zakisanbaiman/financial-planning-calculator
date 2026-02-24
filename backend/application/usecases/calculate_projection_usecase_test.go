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
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newTestCalculateProjectionUseCase(
	planRepo *MockFinancialPlanRepository,
	goalRepo *MockGoalRepository,
) CalculateProjectionUseCase {
	calcService := services.NewFinancialCalculationService()
	recommendService := services.NewGoalRecommendationService(calcService)
	return NewCalculateProjectionUseCase(planRepo, goalRepo, calcService, recommendService)
}

func TestCalculateAssetProjection(t *testing.T) {
	tests := []struct {
		name        string
		input       AssetProjectionInput
		setupMocks  func(*MockFinancialPlanRepository, *MockGoalRepository)
		expectError bool
		errContains string
	}{
		{
			name: "正常系: 資産推移計算",
			input: AssetProjectionInput{
				UserID: "user-001",
				Years:  10,
			},
			setupMocks: func(fp *MockFinancialPlanRepository, gr *MockGoalRepository) {
				plan := newTestFinancialPlan("user-001")
				fp.On("FindByUserID", mock.Anything, entities.UserID("user-001")).Return(plan, nil)
			},
			expectError: false,
		},
		{
			name: "正常系: 1年間の資産推移",
			input: AssetProjectionInput{
				UserID: "user-002",
				Years:  1,
			},
			setupMocks: func(fp *MockFinancialPlanRepository, gr *MockGoalRepository) {
				plan := newTestFinancialPlan("user-002")
				fp.On("FindByUserID", mock.Anything, entities.UserID("user-002")).Return(plan, nil)
			},
			expectError: false,
		},
		{
			name: "異常系: 財務計画が存在しない",
			input: AssetProjectionInput{
				UserID: "user-999",
				Years:  5,
			},
			setupMocks: func(fp *MockFinancialPlanRepository, gr *MockGoalRepository) {
				fp.On("FindByUserID", mock.Anything, entities.UserID("user-999")).
					Return(nil, errors.New("財務データが見つかりません"))
			},
			expectError: true,
			errContains: "財務計画の取得に失敗しました",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			planRepo := new(MockFinancialPlanRepository)
			goalRepo := new(MockGoalRepository)
			tt.setupMocks(planRepo, goalRepo)

			uc := newTestCalculateProjectionUseCase(planRepo, goalRepo)
			output, err := uc.CalculateAssetProjection(context.Background(), tt.input)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, output)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, output)
				assert.NotEmpty(t, output.Projections)
				assert.Equal(t, tt.input.Years, len(output.Projections))
			}
			planRepo.AssertExpectations(t)
		})
	}
}

func TestCalculateRetirementProjection(t *testing.T) {
	tests := []struct {
		name        string
		input       RetirementProjectionInput
		setupMocks  func(*MockFinancialPlanRepository, *MockGoalRepository)
		expectError bool
		errContains string
	}{
		{
			name: "正常系: 退職資金予測",
			input: RetirementProjectionInput{
				UserID: "user-001",
			},
			setupMocks: func(fp *MockFinancialPlanRepository, gr *MockGoalRepository) {
				plan := newTestFinancialPlanWithRetirementData("user-001")
				fp.On("FindByUserID", mock.Anything, entities.UserID("user-001")).Return(plan, nil)
			},
			expectError: false,
		},
		{
			name: "異常系: 財務計画が存在しない",
			input: RetirementProjectionInput{
				UserID: "user-999",
			},
			setupMocks: func(fp *MockFinancialPlanRepository, gr *MockGoalRepository) {
				fp.On("FindByUserID", mock.Anything, entities.UserID("user-999")).
					Return(nil, errors.New("財務データが見つかりません"))
			},
			expectError: true,
			errContains: "財務計画の取得に失敗しました",
		},
		{
			name: "異常系: 退職データが未設定",
			input: RetirementProjectionInput{
				UserID: "user-002",
			},
			setupMocks: func(fp *MockFinancialPlanRepository, gr *MockGoalRepository) {
				// 退職データなしの財務計画
				plan := newTestFinancialPlan("user-002")
				fp.On("FindByUserID", mock.Anything, entities.UserID("user-002")).Return(plan, nil)
			},
			expectError: true,
			errContains: "退職データが設定されていません",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			planRepo := new(MockFinancialPlanRepository)
			goalRepo := new(MockGoalRepository)
			tt.setupMocks(planRepo, goalRepo)

			uc := newTestCalculateProjectionUseCase(planRepo, goalRepo)
			output, err := uc.CalculateRetirementProjection(context.Background(), tt.input)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, output)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, output)
				assert.NotNil(t, output.Calculation)
				assert.NotEmpty(t, output.SufficiencyLevel)
			}
			planRepo.AssertExpectations(t)
		})
	}
}

func TestCalculateEmergencyFundProjection(t *testing.T) {
	tests := []struct {
		name        string
		input       EmergencyFundProjectionInput
		setupMocks  func(*MockFinancialPlanRepository, *MockGoalRepository)
		expectError bool
		errContains string
	}{
		{
			name: "正常系: 緊急資金予測",
			input: EmergencyFundProjectionInput{
				UserID: "user-001",
			},
			setupMocks: func(fp *MockFinancialPlanRepository, gr *MockGoalRepository) {
				plan := newTestFinancialPlan("user-001")
				fp.On("FindByUserID", mock.Anything, entities.UserID("user-001")).Return(plan, nil)
			},
			expectError: false,
		},
		{
			name: "異常系: 財務計画が存在しない",
			input: EmergencyFundProjectionInput{
				UserID: "user-999",
			},
			setupMocks: func(fp *MockFinancialPlanRepository, gr *MockGoalRepository) {
				fp.On("FindByUserID", mock.Anything, entities.UserID("user-999")).
					Return(nil, errors.New("財務データが見つかりません"))
			},
			expectError: true,
			errContains: "財務計画の取得に失敗しました",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			planRepo := new(MockFinancialPlanRepository)
			goalRepo := new(MockGoalRepository)
			tt.setupMocks(planRepo, goalRepo)

			uc := newTestCalculateProjectionUseCase(planRepo, goalRepo)
			output, err := uc.CalculateEmergencyFundProjection(context.Background(), tt.input)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, output)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, output)
				assert.NotNil(t, output.Status)
				assert.NotEmpty(t, output.Priority)
			}
			planRepo.AssertExpectations(t)
		})
	}
}

func TestCalculateComprehensiveProjection(t *testing.T) {
	tests := []struct {
		name        string
		input       ComprehensiveProjectionInput
		setupMocks  func(*MockFinancialPlanRepository, *MockGoalRepository)
		expectError bool
		errContains string
	}{
		{
			name: "正常系: 包括的財務予測",
			input: ComprehensiveProjectionInput{
				UserID: "user-001",
				Years:  5,
			},
			setupMocks: func(fp *MockFinancialPlanRepository, gr *MockGoalRepository) {
				plan := newTestFinancialPlan("user-001")
				fp.On("FindByUserID", mock.Anything, entities.UserID("user-001")).Return(plan, nil)
			},
			expectError: false,
		},
		{
			name: "異常系: 財務計画が存在しない",
			input: ComprehensiveProjectionInput{
				UserID: "user-999",
				Years:  5,
			},
			setupMocks: func(fp *MockFinancialPlanRepository, gr *MockGoalRepository) {
				fp.On("FindByUserID", mock.Anything, entities.UserID("user-999")).
					Return(nil, errors.New("財務データが見つかりません"))
			},
			expectError: true,
			errContains: "財務計画の取得に失敗しました",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			planRepo := new(MockFinancialPlanRepository)
			goalRepo := new(MockGoalRepository)
			tt.setupMocks(planRepo, goalRepo)

			uc := newTestCalculateProjectionUseCase(planRepo, goalRepo)
			output, err := uc.CalculateComprehensiveProjection(context.Background(), tt.input)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, output)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, output)
				assert.NotNil(t, output.PlanProjection)
			}
			planRepo.AssertExpectations(t)
		})
	}
}

func TestCalculateGoalProjection(t *testing.T) {
	tests := []struct {
		name        string
		input       GoalProjectionInput
		setupMocks  func(*MockFinancialPlanRepository, *MockGoalRepository)
		expectError bool
		errContains string
	}{
		{
			name: "正常系: 目標達成予測",
			input: GoalProjectionInput{
				UserID: "user-001",
				GoalID: "goal-001",
			},
			setupMocks: func(fp *MockFinancialPlanRepository, gr *MockGoalRepository) {
				goal := newTestGoal("user-001", entities.GoalTypeSavings)
				plan := newTestFinancialPlan("user-001")
				gr.On("FindByID", mock.Anything, entities.GoalID("goal-001")).Return(goal, nil)
				fp.On("FindByUserID", mock.Anything, entities.UserID("user-001")).Return(plan, nil)
			},
			expectError: false,
		},
		{
			name: "異常系: 目標が存在しない",
			input: GoalProjectionInput{
				UserID: "user-001",
				GoalID: "goal-999",
			},
			setupMocks: func(fp *MockFinancialPlanRepository, gr *MockGoalRepository) {
				gr.On("FindByID", mock.Anything, entities.GoalID("goal-999")).
					Return(nil, errors.New("目標が見つかりません"))
			},
			expectError: true,
			errContains: "目標の取得に失敗しました",
		},
		{
			name: "異常系: 財務計画が存在しない",
			input: GoalProjectionInput{
				UserID: "user-999",
				GoalID: "goal-001",
			},
			setupMocks: func(fp *MockFinancialPlanRepository, gr *MockGoalRepository) {
				goal := newTestGoal("user-999", entities.GoalTypeSavings)
				gr.On("FindByID", mock.Anything, entities.GoalID("goal-001")).Return(goal, nil)
				fp.On("FindByUserID", mock.Anything, entities.UserID("user-999")).
					Return(nil, errors.New("財務データが見つかりません"))
			},
			expectError: true,
			errContains: "財務計画の取得に失敗しました",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			planRepo := new(MockFinancialPlanRepository)
			goalRepo := new(MockGoalRepository)
			tt.setupMocks(planRepo, goalRepo)

			uc := newTestCalculateProjectionUseCase(planRepo, goalRepo)
			output, err := uc.CalculateGoalProjection(context.Background(), tt.input)

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
			planRepo.AssertExpectations(t)
			goalRepo.AssertExpectations(t)
		})
	}
}

// TestCalculateRetirementProjection_Scenarios tests different retirement sufficiency scenarios
func TestCalculateRetirementProjection_Scenarios(t *testing.T) {
// These tests exercise different sufficiency calculation paths
tests := []struct {
name               string
setupRetirement    func(userID string) *aggregates.FinancialPlan
expectedSufficiency string
}{
{
name: "シナリオ: 現在年齢が退職年齢に近い",
setupRetirement: func(userID string) *aggregates.FinancialPlan {
plan := newTestFinancialPlan(userID)
monthlyExpenses := mustCreateMoneyUsecase(250000)
pension := mustCreateMoneyUsecase(200000)
rd, _ := entities.NewRetirementData(
entities.UserID(userID),
63, // 現在年齢（退職年齢に近い）
65,
85,
monthlyExpenses,
pension,
)
plan.SetRetirementData(rd)
return plan
},
},
{
name: "シナリオ: 年金が多い（充足率高）",
setupRetirement: func(userID string) *aggregates.FinancialPlan {
plan := newTestFinancialPlan(userID)
monthlyExpenses := mustCreateMoneyUsecase(150000)
// 年金が支出をカバー
pension := mustCreateMoneyUsecase(180000)
rd, _ := entities.NewRetirementData(
entities.UserID(userID),
40,
65,
85,
monthlyExpenses,
pension,
)
plan.SetRetirementData(rd)
return plan
},
},
{
name: "シナリオ: 長寿命・少年金（不足）",
setupRetirement: func(userID string) *aggregates.FinancialPlan {
plan := newTestFinancialPlan(userID)
monthlyExpenses := mustCreateMoneyUsecase(350000)
pension := mustCreateMoneyUsecase(50000)
rd, _ := entities.NewRetirementData(
entities.UserID(userID),
30,
60, // 早期退職
95, // 長寿命
monthlyExpenses,
pension,
)
plan.SetRetirementData(rd)
return plan
},
},
}

calcService := services.NewFinancialCalculationService()
recommendService := services.NewGoalRecommendationService(calcService)

for i, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
userID := "scenario-user-" + string(rune('A'+i))
planRepo := new(MockFinancialPlanRepository)
goalRepo := new(MockGoalRepository)

plan := tt.setupRetirement(userID)
planRepo.On("FindByUserID", mock.Anything, entities.UserID(userID)).Return(plan, nil)

uc := NewCalculateProjectionUseCase(planRepo, goalRepo, calcService, recommendService)
output, err := uc.CalculateRetirementProjection(context.Background(), RetirementProjectionInput{
UserID: entities.UserID(userID),
})

require.NoError(t, err)
assert.NotNil(t, output)
assert.NotNil(t, output.Calculation)
assert.NotEmpty(t, output.SufficiencyLevel)
assert.NotEmpty(t, output.Recommendations)
planRepo.AssertExpectations(t)
})
}
}

// TestCalculateEmergencyFundProjection_Scenarios tests different emergency fund scenarios
func TestCalculateEmergencyFundProjection_Scenarios(t *testing.T) {
calcService := services.NewFinancialCalculationService()
recommendService := services.NewGoalRecommendationService(calcService)

tests := []struct {
name    string
userID  string
monthly float64
savings float64
}{
{
name:    "十分な緊急資金",
userID:  "ef-user-1",
monthly: 200000,
savings: 2400000, // 12か月分
},
{
name:    "部分的な緊急資金",
userID:  "ef-user-2",
monthly: 300000,
savings: 600000, // 2か月分
},
{
name:    "緊急資金なし",
userID:  "ef-user-3",
monthly: 400000,
savings: 0,
},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
planRepo := new(MockFinancialPlanRepository)
goalRepo := new(MockGoalRepository)

income, _ := valueobjects.NewMoneyJPY(tt.monthly + 100000)
expenses := entities.ExpenseCollection{
{Category: "生活費", Amount: mustCreateMoneyUsecase(tt.monthly)},
}
var savings entities.SavingsCollection
if tt.savings > 0 {
savings = entities.SavingsCollection{
{Type: "deposit", Amount: mustCreateMoneyUsecase(tt.savings)},
}
} else {
savings = entities.SavingsCollection{}
}
investReturn, _ := valueobjects.NewRate(3.0)
inflation, _ := valueobjects.NewRate(2.0)

profile, _ := entities.NewFinancialProfile(
entities.UserID(tt.userID),
income,
expenses,
savings,
investReturn,
inflation,
)
plan, _ := aggregates.NewFinancialPlan(profile)

planRepo.On("FindByUserID", mock.Anything, entities.UserID(tt.userID)).Return(plan, nil)

uc := NewCalculateProjectionUseCase(planRepo, goalRepo, calcService, recommendService)
output, err := uc.CalculateEmergencyFundProjection(context.Background(), EmergencyFundProjectionInput{
UserID: entities.UserID(tt.userID),
})

require.NoError(t, err)
assert.NotNil(t, output)
assert.NotNil(t, output.Status)
assert.NotEmpty(t, output.Priority)
planRepo.AssertExpectations(t)
})
}
}
