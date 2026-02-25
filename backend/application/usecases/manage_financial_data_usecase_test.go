package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/financial-planning-calculator/backend/domain/aggregates"
	"github.com/financial-planning-calculator/backend/domain/entities"
	"github.com/financial-planning-calculator/backend/domain/valueobjects"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newTestFinancialPlan はテスト用の財務計画を作成するヘルパー
func newTestFinancialPlan(userID entities.UserID) *aggregates.FinancialPlan {
	monthlyIncome, _ := valueobjects.NewMoneyJPY(400000)
	expenses := entities.ExpenseCollection{
		{Category: "住居費", Amount: mustNewMoney(120000)},
		{Category: "食費", Amount: mustNewMoney(60000)},
	}
	savings := entities.SavingsCollection{
		{Type: "deposit", Amount: mustNewMoney(1000000)},
	}
	investmentReturn, _ := valueobjects.NewRate(5.0)
	inflationRate, _ := valueobjects.NewRate(2.0)

	profile, err := entities.NewFinancialProfile(userID, monthlyIncome, expenses, savings, investmentReturn, inflationRate)
	if err != nil {
		panic("テスト用財務プロファイルの作成に失敗: " + err.Error())
	}
	plan, err := aggregates.NewFinancialPlan(profile)
	if err != nil {
		panic("テスト用財務計画の作成に失敗: " + err.Error())
	}
	return plan
}

// mustNewMoney は金額を作成するヘルパー（テスト専用）
func mustNewMoney(amount float64) valueobjects.Money {
	m, err := valueobjects.NewMoneyJPY(amount)
	if err != nil {
		panic(err)
	}
	return m
}

// ===========================
// CreateFinancialPlan Tests
// ===========================

func TestManageFinancialDataUseCase_CreateFinancialPlan(t *testing.T) {
	ctx := context.Background()
	baseInput := CreateFinancialPlanInput{
		UserID:           "user-001",
		MonthlyIncome:    400000,
		MonthlyExpenses:  []ExpenseItem{{Category: "住居費", Amount: 120000}},
		CurrentSavings:   []SavingsItem{{Type: "deposit", Amount: 1000000}},
		InvestmentReturn: 5.0,
		InflationRate:    2.0,
	}

	t.Run("正常系: 財務計画を新規作成できる", func(t *testing.T) {
		mockRepo := new(MockFinancialPlanRepository)
		mockRepo.On("ExistsByUserID", mock_anything(), entities.UserID("user-001")).Return(false, nil)
		mockRepo.On("Save", mock_anything(), mock_anything()).Return(nil)

		uc := NewManageFinancialDataUseCase(mockRepo)
		output, err := uc.CreateFinancialPlan(ctx, baseInput)

		require.NoError(t, err)
		assert.NotEmpty(t, output.PlanID)
		assert.Equal(t, entities.UserID("user-001"), output.UserID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("異常系: 財務計画が既に存在する場合はエラー", func(t *testing.T) {
		mockRepo := new(MockFinancialPlanRepository)
		mockRepo.On("ExistsByUserID", mock_anything(), entities.UserID("user-001")).Return(true, nil)

		uc := NewManageFinancialDataUseCase(mockRepo)
		_, err := uc.CreateFinancialPlan(ctx, baseInput)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "既に存在します")
		mockRepo.AssertExpectations(t)
	})

	t.Run("異常系: ExistsByUserIDでリポジトリエラーが発生した場合", func(t *testing.T) {
		mockRepo := new(MockFinancialPlanRepository)
		mockRepo.On("ExistsByUserID", mock_anything(), entities.UserID("user-001")).Return(false, errors.New("db error"))

		uc := NewManageFinancialDataUseCase(mockRepo)
		_, err := uc.CreateFinancialPlan(ctx, baseInput)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "既存財務計画の確認に失敗しました")
		mockRepo.AssertExpectations(t)
	})

	t.Run("異常系: Saveでリポジトリエラーが発生した場合", func(t *testing.T) {
		mockRepo := new(MockFinancialPlanRepository)
		mockRepo.On("ExistsByUserID", mock_anything(), entities.UserID("user-001")).Return(false, nil)
		mockRepo.On("Save", mock_anything(), mock_anything()).Return(errors.New("db error"))

		uc := NewManageFinancialDataUseCase(mockRepo)
		_, err := uc.CreateFinancialPlan(ctx, baseInput)

		require.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

// ===========================
// GetFinancialPlan Tests
// ===========================

func TestManageFinancialDataUseCase_GetFinancialPlan(t *testing.T) {
	ctx := context.Background()

	t.Run("正常系: 財務計画を取得できる", func(t *testing.T) {
		mockRepo := new(MockFinancialPlanRepository)
		plan := newTestFinancialPlan("user-001")
		mockRepo.On("FindByUserID", mock_anything(), entities.UserID("user-001")).Return(plan, nil)

		uc := NewManageFinancialDataUseCase(mockRepo)
		output, err := uc.GetFinancialPlan(ctx, GetFinancialPlanInput{UserID: "user-001"})

		require.NoError(t, err)
		assert.NotNil(t, output.Plan)
		mockRepo.AssertExpectations(t)
	})

	t.Run("異常系: 財務計画が存在しない場合はエラー", func(t *testing.T) {
		mockRepo := new(MockFinancialPlanRepository)
		mockRepo.On("FindByUserID", mock_anything(), entities.UserID("user-999")).Return(nil, errors.New("not found"))

		uc := NewManageFinancialDataUseCase(mockRepo)
		_, err := uc.GetFinancialPlan(ctx, GetFinancialPlanInput{UserID: "user-999"})

		require.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

// ===========================
// UpdateFinancialProfile Tests
// ===========================

func TestManageFinancialDataUseCase_UpdateFinancialProfile(t *testing.T) {
	ctx := context.Background()
	input := UpdateFinancialProfileInput{
		UserID:           "user-001",
		MonthlyIncome:    500000,
		MonthlyExpenses:  []ExpenseItem{{Category: "住居費", Amount: 150000}},
		CurrentSavings:   []SavingsItem{{Type: "deposit", Amount: 2000000}},
		InvestmentReturn: 6.0,
		InflationRate:    2.5,
	}

	t.Run("正常系: 財務プロファイルを更新できる", func(t *testing.T) {
		mockRepo := new(MockFinancialPlanRepository)
		plan := newTestFinancialPlan("user-001")
		mockRepo.On("FindByUserID", mock_anything(), entities.UserID("user-001")).Return(plan, nil)
		mockRepo.On("Update", mock_anything(), mock_anything()).Return(nil)

		uc := NewManageFinancialDataUseCase(mockRepo)
		output, err := uc.UpdateFinancialProfile(ctx, input)

		require.NoError(t, err)
		assert.NotNil(t, output)
		mockRepo.AssertExpectations(t)
	})

	t.Run("異常系: FindByUserIDでエラーが発生した場合", func(t *testing.T) {
		mockRepo := new(MockFinancialPlanRepository)
		mockRepo.On("FindByUserID", mock_anything(), entities.UserID("user-001")).Return(nil, errors.New("not found"))

		uc := NewManageFinancialDataUseCase(mockRepo)
		_, err := uc.UpdateFinancialProfile(ctx, input)

		require.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("異常系: Updateでエラーが発生した場合", func(t *testing.T) {
		mockRepo := new(MockFinancialPlanRepository)
		plan := newTestFinancialPlan("user-001")
		mockRepo.On("FindByUserID", mock_anything(), entities.UserID("user-001")).Return(plan, nil)
		mockRepo.On("Update", mock_anything(), mock_anything()).Return(errors.New("db error"))

		uc := NewManageFinancialDataUseCase(mockRepo)
		_, err := uc.UpdateFinancialProfile(ctx, input)

		require.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

// ===========================
// DeleteFinancialPlan Tests
// ===========================

func TestManageFinancialDataUseCase_DeleteFinancialPlan(t *testing.T) {
	ctx := context.Background()

	t.Run("正常系: 財務計画を削除できる", func(t *testing.T) {
		mockRepo := new(MockFinancialPlanRepository)
		plan := newTestFinancialPlan("user-001")
		mockRepo.On("FindByUserID", mock_anything(), entities.UserID("user-001")).Return(plan, nil)
		mockRepo.On("Delete", mock_anything(), plan.ID()).Return(nil)

		uc := NewManageFinancialDataUseCase(mockRepo)
		err := uc.DeleteFinancialPlan(ctx, DeleteFinancialPlanInput{UserID: "user-001"})

		require.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("異常系: FindByUserIDでエラーが発生した場合", func(t *testing.T) {
		mockRepo := new(MockFinancialPlanRepository)
		mockRepo.On("FindByUserID", mock_anything(), entities.UserID("user-001")).Return(nil, errors.New("not found"))

		uc := NewManageFinancialDataUseCase(mockRepo)
		err := uc.DeleteFinancialPlan(ctx, DeleteFinancialPlanInput{UserID: "user-001"})

		require.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("異常系: Deleteでエラーが発生した場合", func(t *testing.T) {
		mockRepo := new(MockFinancialPlanRepository)
		plan := newTestFinancialPlan("user-001")
		mockRepo.On("FindByUserID", mock_anything(), entities.UserID("user-001")).Return(plan, nil)
		mockRepo.On("Delete", mock_anything(), plan.ID()).Return(errors.New("db error"))

		uc := NewManageFinancialDataUseCase(mockRepo)
		err := uc.DeleteFinancialPlan(ctx, DeleteFinancialPlanInput{UserID: "user-001"})

		require.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}


