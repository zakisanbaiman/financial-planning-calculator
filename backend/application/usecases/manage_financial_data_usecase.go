package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/financial-planning-calculator/backend/domain/aggregates"
	"github.com/financial-planning-calculator/backend/domain/entities"
	"github.com/financial-planning-calculator/backend/domain/repositories"
	"github.com/financial-planning-calculator/backend/domain/valueobjects"
)

// ManageFinancialDataUseCase は財務データ管理のユースケース
type ManageFinancialDataUseCase interface {
	// CreateFinancialPlan は新しい財務計画を作成する
	CreateFinancialPlan(ctx context.Context, input CreateFinancialPlanInput) (*CreateFinancialPlanOutput, error)

	// GetFinancialPlan は財務計画を取得する
	GetFinancialPlan(ctx context.Context, input GetFinancialPlanInput) (*GetFinancialPlanOutput, error)

	// UpdateFinancialProfile は財務プロファイルを更新する
	UpdateFinancialProfile(ctx context.Context, input UpdateFinancialProfileInput) (*UpdateFinancialProfileOutput, error)

	// UpdateRetirementData は退職データを更新する
	UpdateRetirementData(ctx context.Context, input UpdateRetirementDataInput) (*UpdateRetirementDataOutput, error)

	// UpdateEmergencyFund は緊急資金設定を更新する
	UpdateEmergencyFund(ctx context.Context, input UpdateEmergencyFundInput) (*UpdateEmergencyFundOutput, error)

	// DeleteFinancialPlan は財務計画を削除する
	DeleteFinancialPlan(ctx context.Context, input DeleteFinancialPlanInput) error
}

// CreateFinancialPlanInput は財務計画作成の入力
type CreateFinancialPlanInput struct {
	UserID                     entities.UserID `json:"user_id"`
	MonthlyIncome              float64         `json:"monthly_income"`
	MonthlyExpenses            []ExpenseItem   `json:"monthly_expenses"`
	CurrentSavings             []SavingsItem   `json:"current_savings"`
	InvestmentReturn           float64         `json:"investment_return"`
	InflationRate              float64         `json:"inflation_rate"`
	RetirementAge              *int            `json:"retirement_age,omitempty"`
	MonthlyRetirementExpenses  *float64        `json:"monthly_retirement_expenses,omitempty"`
	PensionAmount              *float64        `json:"pension_amount,omitempty"`
	EmergencyFundTargetMonths  *int            `json:"emergency_fund_target_months,omitempty"`
	EmergencyFundCurrentAmount *float64        `json:"emergency_fund_current_amount,omitempty"`
}

// ExpenseItem は支出項目
type ExpenseItem struct {
	Category    string  `json:"category"`
	Amount      float64 `json:"amount"`
	Description *string `json:"description,omitempty"`
}

// SavingsItem は貯蓄項目
type SavingsItem struct {
	Type        string  `json:"type"`
	Amount      float64 `json:"amount"`
	Description *string `json:"description,omitempty"`
}

// CreateFinancialPlanOutput は財務計画作成の出力
type CreateFinancialPlanOutput struct {
	PlanID    aggregates.FinancialPlanID `json:"plan_id"`
	UserID    entities.UserID            `json:"user_id"`
	CreatedAt string                     `json:"created_at"`
}

// GetFinancialPlanInput は財務計画取得の入力
type GetFinancialPlanInput struct {
	UserID entities.UserID `json:"user_id"`
}

// GetFinancialPlanOutput は財務計画取得の出力
type GetFinancialPlanOutput struct {
	Plan *aggregates.FinancialPlan `json:"plan"`
}

// FinancialDataResponse はフロントエンド向けの財務データレスポンス
type FinancialDataResponse struct {
	UserID        string                 `json:"user_id"`
	Profile       map[string]interface{} `json:"profile,omitempty"`
	Retirement    map[string]interface{} `json:"retirement,omitempty"`
	EmergencyFund map[string]interface{} `json:"emergency_fund,omitempty"`
	CreatedAt     string                 `json:"created_at,omitempty"`
	UpdatedAt     string                 `json:"updated_at,omitempty"`
}

// UpdateFinancialProfileInput は財務プロファイル更新の入力
type UpdateFinancialProfileInput struct {
	UserID           entities.UserID `json:"user_id"`
	MonthlyIncome    float64         `json:"monthly_income"`
	MonthlyExpenses  []ExpenseItem   `json:"monthly_expenses"`
	CurrentSavings   []SavingsItem   `json:"current_savings"`
	InvestmentReturn float64         `json:"investment_return"`
	InflationRate    float64         `json:"inflation_rate"`
}

// UpdateFinancialProfileOutput は財務プロファイル更新の出力
// フロントエンド向けに FinancialDataResponse を返す
type UpdateFinancialProfileOutput struct {
	*FinancialDataResponse
}

// UpdateRetirementDataInput は退職データ更新の入力
type UpdateRetirementDataInput struct {
	UserID                    entities.UserID `json:"user_id"`
	RetirementAge             int             `json:"retirement_age"`
	MonthlyRetirementExpenses float64         `json:"monthly_retirement_expenses"`
	PensionAmount             float64         `json:"pension_amount"`
}

// UpdateRetirementDataOutput は退職データ更新の出力
// フロントエンド向けに FinancialDataResponse を返す
type UpdateRetirementDataOutput struct {
	*FinancialDataResponse
}

// UpdateEmergencyFundInput は緊急資金設定更新の入力
type UpdateEmergencyFundInput struct {
	UserID        entities.UserID `json:"user_id"`
	TargetMonths  int             `json:"target_months"`
	CurrentAmount float64         `json:"current_amount"`
}

// UpdateEmergencyFundOutput は緊急資金設定更新の出力
// フロントエンド向けに FinancialDataResponse を返す
type UpdateEmergencyFundOutput struct {
	*FinancialDataResponse
}

// DeleteFinancialPlanInput は財務計画削除の入力
type DeleteFinancialPlanInput struct {
	UserID entities.UserID `json:"user_id"`
}

// manageFinancialDataUseCaseImpl はManageFinancialDataUseCaseの実装
type manageFinancialDataUseCaseImpl struct {
	financialPlanRepo repositories.FinancialPlanRepository
}

// NewManageFinancialDataUseCase は新しいManageFinancialDataUseCaseを作成する
func NewManageFinancialDataUseCase(
	financialPlanRepo repositories.FinancialPlanRepository,
) ManageFinancialDataUseCase {
	return &manageFinancialDataUseCaseImpl{
		financialPlanRepo: financialPlanRepo,
	}
}

// CreateFinancialPlan は新しい財務計画を作成する
func (uc *manageFinancialDataUseCaseImpl) CreateFinancialPlan(
	ctx context.Context,
	input CreateFinancialPlanInput,
) (*CreateFinancialPlanOutput, error) {
	// 既存の財務計画があるかチェック
	exists, err := uc.financialPlanRepo.ExistsByUserID(ctx, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("既存財務計画の確認に失敗しました: %w", err)
	}

	if exists {
		return nil, errors.New("ユーザーの財務計画は既に存在します")
	}

	// 財務プロファイルを作成
	profile, err := uc.createFinancialProfile(input)
	if err != nil {
		return nil, fmt.Errorf("財務プロファイルの作成に失敗しました: %w", err)
	}

	// 財務計画を作成
	plan, err := aggregates.NewFinancialPlan(profile)
	if err != nil {
		return nil, fmt.Errorf("財務計画の作成に失敗しました: %w", err)
	}

	// 退職データが提供されている場合は設定
	if input.RetirementAge != nil && input.MonthlyRetirementExpenses != nil && input.PensionAmount != nil {
		retirementData, err := uc.createRetirementData(input.UserID, *input.RetirementAge, *input.MonthlyRetirementExpenses, *input.PensionAmount)
		if err != nil {
			return nil, fmt.Errorf("退職データの作成に失敗しました: %w", err)
		}

		err = plan.SetRetirementData(retirementData)
		if err != nil {
			return nil, fmt.Errorf("退職データの設定に失敗しました: %w", err)
		}
	}

	// 緊急資金設定が提供されている場合は設定
	if input.EmergencyFundTargetMonths != nil && input.EmergencyFundCurrentAmount != nil {
		currentFund, err := valueobjects.NewMoneyJPY(*input.EmergencyFundCurrentAmount)
		if err != nil {
			return nil, fmt.Errorf("緊急資金額の作成に失敗しました: %w", err)
		}

		emergencyConfig, err := aggregates.NewEmergencyFundConfig(*input.EmergencyFundTargetMonths, currentFund)
		if err != nil {
			return nil, fmt.Errorf("緊急資金設定の作成に失敗しました: %w", err)
		}

		err = plan.UpdateEmergencyFund(emergencyConfig)
		if err != nil {
			return nil, fmt.Errorf("緊急資金設定の更新に失敗しました: %w", err)
		}
	}

	// 財務計画を保存
	err = uc.financialPlanRepo.Save(ctx, plan)
	if err != nil {
		return nil, fmt.Errorf("財務計画の保存に失敗しました: %w", err)
	}

	return &CreateFinancialPlanOutput{
		PlanID:    plan.ID(),
		UserID:    input.UserID,
		CreatedAt: plan.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// GetFinancialPlan は財務計画を取得する
func (uc *manageFinancialDataUseCaseImpl) GetFinancialPlan(
	ctx context.Context,
	input GetFinancialPlanInput,
) (*GetFinancialPlanOutput, error) {
	plan, err := uc.financialPlanRepo.FindByUserID(ctx, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("財務計画の取得に失敗しました: %w", err)
	}

	return &GetFinancialPlanOutput{
		Plan: plan,
	}, nil
}

// UpdateFinancialProfile は財務プロファイルを更新する
func (uc *manageFinancialDataUseCaseImpl) UpdateFinancialProfile(
	ctx context.Context,
	input UpdateFinancialProfileInput,
) (*UpdateFinancialProfileOutput, error) {
	// 既存の財務計画を取得
	plan, err := uc.financialPlanRepo.FindByUserID(ctx, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("財務計画の取得に失敗しました: %w", err)
	}

	// 新しい財務プロファイルを作成
	profile, err := uc.createFinancialProfileFromUpdate(input)
	if err != nil {
		return nil, fmt.Errorf("財務プロファイルの作成に失敗しました: %w", err)
	}

	// 財務プロファイルを更新
	err = plan.UpdateProfile(profile)
	if err != nil {
		return nil, fmt.Errorf("財務プロファイルの更新に失敗しました: %w", err)
	}

	// 財務計画を保存
	err = uc.financialPlanRepo.Update(ctx, plan)
	if err != nil {
		return nil, fmt.Errorf("財務計画の保存に失敗しました: %w", err)
	}

	// フロントエンド向けレスポンスに変換して返す
	return convertPlanToFinancialDataResponse(plan, input.UserID), nil
}

// convertPlanToFinancialDataResponse は FinancialPlan を FinancialDataResponse に変換
func convertPlanToFinancialDataResponse(plan *aggregates.FinancialPlan, userID entities.UserID) *UpdateFinancialProfileOutput {
	if plan == nil {
		return &UpdateFinancialProfileOutput{
			FinancialDataResponse: &FinancialDataResponse{
				UserID: string(userID),
			},
		}
	}

	response := &FinancialDataResponse{
		UserID: string(userID),
	}

	// Profile を変換
	if profile := plan.Profile(); profile != nil {
		profileMap := map[string]interface{}{
			"monthly_income":    profile.MonthlyIncome(),
			"monthly_expenses":  profile.MonthlyExpenses(),
			"current_savings":   profile.CurrentSavings(),
			"investment_return": profile.InvestmentReturn(),
			"inflation_rate":    profile.InflationRate(),
		}
		response.Profile = profileMap
	}

	// RetirementData を変換
	if retirement := plan.RetirementData(); retirement != nil {
		retirementMap := map[string]interface{}{
			"retirement_age":              retirement.RetirementAge(),
			"monthly_retirement_expenses": retirement.MonthlyRetirementExpenses(),
			"pension_amount":              retirement.PensionAmount(),
		}
		response.Retirement = retirementMap
	}

	// EmergencyFund を変換
	if emergencyFund := plan.EmergencyFund(); emergencyFund != nil {
		emergencyMap := map[string]interface{}{
			"target_months": emergencyFund.TargetMonths,
			"current_fund":  emergencyFund.CurrentFund,
		}
		response.EmergencyFund = emergencyMap
	}

	return &UpdateFinancialProfileOutput{
		FinancialDataResponse: response,
	}
}

// UpdateRetirementData は退職データを更新する
func (uc *manageFinancialDataUseCaseImpl) UpdateRetirementData(
	ctx context.Context,
	input UpdateRetirementDataInput,
) (*UpdateRetirementDataOutput, error) {
	// 既存の財務計画を取得
	plan, err := uc.financialPlanRepo.FindByUserID(ctx, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("財務計画の取得に失敗しました: %w", err)
	}

	// 退職データを作成
	retirementData, err := uc.createRetirementData(input.UserID, input.RetirementAge, input.MonthlyRetirementExpenses, input.PensionAmount)
	if err != nil {
		return nil, fmt.Errorf("退職データの作成に失敗しました: %w", err)
	}

	// 退職データを設定
	err = plan.SetRetirementData(retirementData)
	if err != nil {
		return nil, fmt.Errorf("退職データの設定に失敗しました: %w", err)
	}

	// 財務計画を保存
	err = uc.financialPlanRepo.Update(ctx, plan)
	if err != nil {
		return nil, fmt.Errorf("財務計画の保存に失敗しました: %w", err)
	}

	// フロントエンド向けレスポンスに変換して返す
	return &UpdateRetirementDataOutput{
		FinancialDataResponse: convertPlanToFinancialDataResponse(plan, input.UserID).FinancialDataResponse,
	}, nil
}

// UpdateEmergencyFund は緊急資金設定を更新する
func (uc *manageFinancialDataUseCaseImpl) UpdateEmergencyFund(
	ctx context.Context,
	input UpdateEmergencyFundInput,
) (*UpdateEmergencyFundOutput, error) {
	// 既存の財務計画を取得
	plan, err := uc.financialPlanRepo.FindByUserID(ctx, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("財務計画の取得に失敗しました: %w", err)
	}

	// 緊急資金設定を作成
	currentFund, err := valueobjects.NewMoneyJPY(input.CurrentAmount)
	if err != nil {
		return nil, fmt.Errorf("緊急資金額の作成に失敗しました: %w", err)
	}

	emergencyConfig, err := aggregates.NewEmergencyFundConfig(input.TargetMonths, currentFund)
	if err != nil {
		return nil, fmt.Errorf("緊急資金設定の作成に失敗しました: %w", err)
	}

	// 緊急資金設定を更新
	err = plan.UpdateEmergencyFund(emergencyConfig)
	if err != nil {
		return nil, fmt.Errorf("緊急資金設定の更新に失敗しました: %w", err)
	}

	// 財務計画を保存
	err = uc.financialPlanRepo.Update(ctx, plan)
	if err != nil {
		return nil, fmt.Errorf("財務計画の保存に失敗しました: %w", err)
	}

	// フロントエンド向けレスポンスに変換して返す
	return &UpdateEmergencyFundOutput{
		FinancialDataResponse: convertPlanToFinancialDataResponse(plan, input.UserID).FinancialDataResponse,
	}, nil
}

// DeleteFinancialPlan は財務計画を削除する
func (uc *manageFinancialDataUseCaseImpl) DeleteFinancialPlan(
	ctx context.Context,
	input DeleteFinancialPlanInput,
) error {
	// 既存の財務計画を取得
	plan, err := uc.financialPlanRepo.FindByUserID(ctx, input.UserID)
	if err != nil {
		return fmt.Errorf("財務計画の取得に失敗しました: %w", err)
	}

	// 財務計画を削除
	err = uc.financialPlanRepo.Delete(ctx, plan.ID())
	if err != nil {
		return fmt.Errorf("財務計画の削除に失敗しました: %w", err)
	}

	return nil
}

// createFinancialProfile は財務プロファイルを作成する
func (uc *manageFinancialDataUseCaseImpl) createFinancialProfile(input CreateFinancialPlanInput) (*entities.FinancialProfile, error) {
	// 月収を作成
	monthlyIncome, err := valueobjects.NewMoneyJPY(input.MonthlyIncome)
	if err != nil {
		return nil, fmt.Errorf("月収の作成に失敗しました: %w", err)
	}

	// 月間支出を作成
	monthlyExpenses, err := uc.createExpenseCollection(input.MonthlyExpenses)
	if err != nil {
		return nil, fmt.Errorf("月間支出の作成に失敗しました: %w", err)
	}

	// 現在の貯蓄を作成
	currentSavings, err := uc.createSavingsCollection(input.CurrentSavings)
	if err != nil {
		return nil, fmt.Errorf("現在の貯蓄の作成に失敗しました: %w", err)
	}

	// 投資利回りを作成
	investmentReturn, err := valueobjects.NewRate(input.InvestmentReturn)
	if err != nil {
		return nil, fmt.Errorf("投資利回りの作成に失敗しました: %w", err)
	}

	// インフレ率を作成
	inflationRate, err := valueobjects.NewRate(input.InflationRate)
	if err != nil {
		return nil, fmt.Errorf("インフレ率の作成に失敗しました: %w", err)
	}

	// 財務プロファイルを作成
	return entities.NewFinancialProfile(
		input.UserID,
		monthlyIncome,
		*monthlyExpenses,
		*currentSavings,
		investmentReturn,
		inflationRate,
	)
}

// createFinancialProfileFromUpdate は更新用の財務プロファイルを作成する
func (uc *manageFinancialDataUseCaseImpl) createFinancialProfileFromUpdate(input UpdateFinancialProfileInput) (*entities.FinancialProfile, error) {
	// 月収を作成
	monthlyIncome, err := valueobjects.NewMoneyJPY(input.MonthlyIncome)
	if err != nil {
		return nil, fmt.Errorf("月収の作成に失敗しました: %w", err)
	}

	// 月間支出を作成
	monthlyExpenses, err := uc.createExpenseCollection(input.MonthlyExpenses)
	if err != nil {
		return nil, fmt.Errorf("月間支出の作成に失敗しました: %w", err)
	}

	// 現在の貯蓄を作成
	currentSavings, err := uc.createSavingsCollection(input.CurrentSavings)
	if err != nil {
		return nil, fmt.Errorf("現在の貯蓄の作成に失敗しました: %w", err)
	}

	// 投資利回りを作成
	investmentReturn, err := valueobjects.NewRate(input.InvestmentReturn)
	if err != nil {
		return nil, fmt.Errorf("投資利回りの作成に失敗しました: %w", err)
	}

	// インフレ率を作成
	inflationRate, err := valueobjects.NewRate(input.InflationRate)
	if err != nil {
		return nil, fmt.Errorf("インフレ率の作成に失敗しました: %w", err)
	}

	// 財務プロファイルを作成
	return entities.NewFinancialProfile(
		input.UserID,
		monthlyIncome,
		*monthlyExpenses,
		*currentSavings,
		investmentReturn,
		inflationRate,
	)
}

// createExpenseCollection は支出コレクションを作成する
func (uc *manageFinancialDataUseCaseImpl) createExpenseCollection(expenses []ExpenseItem) (*entities.ExpenseCollection, error) {
	var collection entities.ExpenseCollection

	for _, expense := range expenses {
		amount, err := valueobjects.NewMoneyJPY(expense.Amount)
		if err != nil {
			return nil, fmt.Errorf("支出額の作成に失敗しました: %w", err)
		}

		description := ""
		if expense.Description != nil {
			description = *expense.Description
		}

		expenseItem := entities.ExpenseItem{
			Category:    expense.Category,
			Amount:      amount,
			Description: description,
		}

		collection = append(collection, expenseItem)
	}

	return &collection, nil
}

// createSavingsCollection は貯蓄コレクションを作成する
func (uc *manageFinancialDataUseCaseImpl) createSavingsCollection(savings []SavingsItem) (*entities.SavingsCollection, error) {
	var collection entities.SavingsCollection

	for _, saving := range savings {
		amount, err := valueobjects.NewMoneyJPY(saving.Amount)
		if err != nil {
			return nil, fmt.Errorf("貯蓄額の作成に失敗しました: %w", err)
		}

		// 貯蓄タイプの検証（deposit, investment, other のみ許可）
		if saving.Type != "deposit" && saving.Type != "investment" && saving.Type != "other" {
			return nil, fmt.Errorf("無効な貯蓄タイプです: %s", saving.Type)
		}

		description := ""
		if saving.Description != nil {
			description = *saving.Description
		}

		savingItem := entities.SavingsItem{
			Type:        saving.Type,
			Amount:      amount,
			Description: description,
		}

		collection = append(collection, savingItem)
	}

	return &collection, nil
}

// createRetirementData は退職データを作成する
func (uc *manageFinancialDataUseCaseImpl) createRetirementData(userID entities.UserID, retirementAge int, monthlyExpenses float64, pensionAmount float64) (*entities.RetirementData, error) {
	monthlyRetirementExpenses, err := valueobjects.NewMoneyJPY(monthlyExpenses)
	if err != nil {
		return nil, fmt.Errorf("月間退職後支出の作成に失敗しました: %w", err)
	}

	pension, err := valueobjects.NewMoneyJPY(pensionAmount)
	if err != nil {
		return nil, fmt.Errorf("年金額の作成に失敗しました: %w", err)
	}

	// 現在の年齢を仮定（実際の実装では別途取得が必要）
	currentAge := 30     // デフォルト値
	lifeExpectancy := 85 // デフォルト値

	return entities.NewRetirementData(
		userID,
		currentAge,
		retirementAge,
		lifeExpectancy,
		monthlyRetirementExpenses,
		pension,
	)
}
