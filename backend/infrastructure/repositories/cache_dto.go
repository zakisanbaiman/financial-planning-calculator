package repositories

import (
	"fmt"
	"time"

	"github.com/financial-planning-calculator/backend/domain/aggregates"
	"github.com/financial-planning-calculator/backend/domain/entities"
	"github.com/financial-planning-calculator/backend/domain/valueobjects"
)

// --- Money / Rate のプリミティブ表現 ---

type moneyDTO struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type rateDTO struct {
	Value float64 `json:"value"` // AsPercentage() の値
}

// --- Goal DTO ---

type goalCacheDTO struct {
	ID                  string    `json:"id"`
	UserID              string    `json:"user_id"`
	GoalType            string    `json:"goal_type"`
	Title               string    `json:"title"`
	TargetAmount        moneyDTO  `json:"target_amount"`
	TargetDate          time.Time `json:"target_date"`
	CurrentAmount       moneyDTO  `json:"current_amount"`
	MonthlyContribution moneyDTO  `json:"monthly_contribution"`
	IsActive            bool      `json:"is_active"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

func goalToDTO(g *entities.Goal) goalCacheDTO {
	return goalCacheDTO{
		ID:       string(g.ID()),
		UserID:   string(g.UserID()),
		GoalType: string(g.GoalType()),
		Title:    g.Title(),
		TargetAmount: moneyDTO{
			Amount:   g.TargetAmount().Amount(),
			Currency: string(g.TargetAmount().Currency()),
		},
		TargetDate: g.TargetDate(),
		CurrentAmount: moneyDTO{
			Amount:   g.CurrentAmount().Amount(),
			Currency: string(g.CurrentAmount().Currency()),
		},
		MonthlyContribution: moneyDTO{
			Amount:   g.MonthlyContribution().Amount(),
			Currency: string(g.MonthlyContribution().Currency()),
		},
		IsActive:  g.IsActive(),
		CreatedAt: g.CreatedAt(),
		UpdatedAt: g.UpdatedAt(),
	}
}

func goalFromDTO(dto goalCacheDTO) (*entities.Goal, error) {
	targetAmount, err := valueobjects.NewMoney(dto.TargetAmount.Amount, valueobjects.Currency(dto.TargetAmount.Currency))
	if err != nil {
		return nil, fmt.Errorf("目標金額の復元に失敗しました: %w", err)
	}

	monthlyContribution, err := valueobjects.NewMoney(dto.MonthlyContribution.Amount, valueobjects.Currency(dto.MonthlyContribution.Currency))
	if err != nil {
		return nil, fmt.Errorf("月間拠出額の復元に失敗しました: %w", err)
	}

	goal, err := entities.NewGoalWithID(
		entities.GoalID(dto.ID),
		entities.UserID(dto.UserID),
		entities.GoalType(dto.GoalType),
		dto.Title,
		targetAmount,
		dto.TargetDate,
		monthlyContribution,
		dto.CreatedAt,
		dto.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("目標エンティティの復元に失敗しました: %w", err)
	}

	currentAmount, err := valueobjects.NewMoney(dto.CurrentAmount.Amount, valueobjects.Currency(dto.CurrentAmount.Currency))
	if err != nil {
		return nil, fmt.Errorf("現在の金額の復元に失敗しました: %w", err)
	}
	if err := goal.UpdateCurrentAmount(currentAmount); err != nil {
		return nil, fmt.Errorf("現在の金額の設定に失敗しました: %w", err)
	}

	if !dto.IsActive {
		goal.Deactivate()
	}

	return goal, nil
}

func goalsToDTOs(goals []*entities.Goal) []goalCacheDTO {
	dtos := make([]goalCacheDTO, len(goals))
	for i, g := range goals {
		dtos[i] = goalToDTO(g)
	}
	return dtos
}

func goalsFromDTOs(dtos []goalCacheDTO) ([]*entities.Goal, error) {
	goals := make([]*entities.Goal, 0, len(dtos))
	for _, dto := range dtos {
		g, err := goalFromDTO(dto)
		if err != nil {
			return nil, err
		}
		goals = append(goals, g)
	}
	return goals, nil
}

// --- FinancialProfile DTO ---

type expenseItemDTO struct {
	Category    string   `json:"category"`
	Amount      moneyDTO `json:"amount"`
	Description string   `json:"description,omitempty"`
}

type savingsItemDTO struct {
	Type        string   `json:"type"`
	Amount      moneyDTO `json:"amount"`
	Description string   `json:"description,omitempty"`
}

type financialProfileCacheDTO struct {
	ID               string           `json:"id"`
	UserID           string           `json:"user_id"`
	MonthlyIncome    moneyDTO         `json:"monthly_income"`
	MonthlyExpenses  []expenseItemDTO `json:"monthly_expenses"`
	CurrentSavings   []savingsItemDTO `json:"current_savings"`
	InvestmentReturn rateDTO          `json:"investment_return"`
	InflationRate    rateDTO          `json:"inflation_rate"`
	CreatedAt        time.Time        `json:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at"`
}

// --- RetirementData DTO ---

type retirementDataCacheDTO struct {
	ID                        string    `json:"id"`
	UserID                    string    `json:"user_id"`
	CurrentAge                int       `json:"current_age"`
	RetirementAge             int       `json:"retirement_age"`
	LifeExpectancy            int       `json:"life_expectancy"`
	MonthlyRetirementExpenses moneyDTO  `json:"monthly_retirement_expenses"`
	PensionAmount             moneyDTO  `json:"pension_amount"`
	CreatedAt                 time.Time `json:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at"`
}

// --- EmergencyFundConfig DTO ---

type emergencyFundConfigDTO struct {
	TargetMonths int      `json:"target_months"`
	CurrentFund  moneyDTO `json:"current_fund"`
}

// --- FinancialPlan DTO ---

type financialPlanCacheDTO struct {
	ID             string                    `json:"id"`
	Profile        financialProfileCacheDTO  `json:"profile"`
	Goals          []goalCacheDTO            `json:"goals"`
	RetirementData *retirementDataCacheDTO   `json:"retirement_data,omitempty"`
	EmergencyFund  *emergencyFundConfigDTO   `json:"emergency_fund,omitempty"`
	CreatedAt      time.Time                 `json:"created_at"`
	UpdatedAt      time.Time                 `json:"updated_at"`
}

func financialPlanToDTO(plan *aggregates.FinancialPlan) financialPlanCacheDTO {
	profile := plan.Profile()

	expenses := make([]expenseItemDTO, len(profile.MonthlyExpenses()))
	for i, e := range profile.MonthlyExpenses() {
		expenses[i] = expenseItemDTO{
			Category:    e.Category,
			Amount:      moneyDTO{Amount: e.Amount.Amount(), Currency: string(e.Amount.Currency())},
			Description: e.Description,
		}
	}

	savings := make([]savingsItemDTO, len(profile.CurrentSavings()))
	for i, s := range profile.CurrentSavings() {
		savings[i] = savingsItemDTO{
			Type:        s.Type,
			Amount:      moneyDTO{Amount: s.Amount.Amount(), Currency: string(s.Amount.Currency())},
			Description: s.Description,
		}
	}

	profileDTO := financialProfileCacheDTO{
		ID:              string(profile.ID()),
		UserID:          string(profile.UserID()),
		MonthlyIncome:   moneyDTO{Amount: profile.MonthlyIncome().Amount(), Currency: string(profile.MonthlyIncome().Currency())},
		MonthlyExpenses: expenses,
		CurrentSavings:  savings,
		InvestmentReturn: rateDTO{Value: profile.InvestmentReturn().AsPercentage()},
		InflationRate:    rateDTO{Value: profile.InflationRate().AsPercentage()},
		CreatedAt:       profile.CreatedAt(),
		UpdatedAt:       profile.UpdatedAt(),
	}

	dto := financialPlanCacheDTO{
		ID:        string(plan.ID()),
		Profile:   profileDTO,
		Goals:     goalsToDTOs(plan.Goals()),
		CreatedAt: plan.CreatedAt(),
		UpdatedAt: plan.UpdatedAt(),
	}

	if rd := plan.RetirementData(); rd != nil {
		dto.RetirementData = &retirementDataCacheDTO{
			ID:     string(rd.ID()),
			UserID: string(rd.UserID()),
			CurrentAge:     rd.CurrentAge(),
			RetirementAge:  rd.RetirementAge(),
			LifeExpectancy: rd.LifeExpectancy(),
			MonthlyRetirementExpenses: moneyDTO{
				Amount:   rd.MonthlyRetirementExpenses().Amount(),
				Currency: string(rd.MonthlyRetirementExpenses().Currency()),
			},
			PensionAmount: moneyDTO{
				Amount:   rd.PensionAmount().Amount(),
				Currency: string(rd.PensionAmount().Currency()),
			},
			CreatedAt: rd.CreatedAt(),
			UpdatedAt: rd.UpdatedAt(),
		}
	}

	if ef := plan.EmergencyFund(); ef != nil {
		dto.EmergencyFund = &emergencyFundConfigDTO{
			TargetMonths: ef.TargetMonths,
			CurrentFund:  moneyDTO{Amount: ef.CurrentFund.Amount(), Currency: string(ef.CurrentFund.Currency())},
		}
	}

	return dto
}

func financialPlanFromDTO(dto financialPlanCacheDTO) (*aggregates.FinancialPlan, error) {
	// FinancialProfile を復元
	monthlyIncome, err := valueobjects.NewMoney(dto.Profile.MonthlyIncome.Amount, valueobjects.Currency(dto.Profile.MonthlyIncome.Currency))
	if err != nil {
		return nil, fmt.Errorf("月収の復元に失敗しました: %w", err)
	}

	expenses := make(entities.ExpenseCollection, len(dto.Profile.MonthlyExpenses))
	for i, e := range dto.Profile.MonthlyExpenses {
		amount, err := valueobjects.NewMoney(e.Amount.Amount, valueobjects.Currency(e.Amount.Currency))
		if err != nil {
			return nil, fmt.Errorf("支出項目の復元に失敗しました: %w", err)
		}
		expenses[i] = entities.ExpenseItem{
			Category:    e.Category,
			Amount:      amount,
			Description: e.Description,
		}
	}

	savings := make(entities.SavingsCollection, len(dto.Profile.CurrentSavings))
	for i, s := range dto.Profile.CurrentSavings {
		amount, err := valueobjects.NewMoney(s.Amount.Amount, valueobjects.Currency(s.Amount.Currency))
		if err != nil {
			return nil, fmt.Errorf("貯蓄項目の復元に失敗しました: %w", err)
		}
		savings[i] = entities.SavingsItem{
			Type:        s.Type,
			Amount:      amount,
			Description: s.Description,
		}
	}

	investmentReturn, err := valueobjects.NewRate(dto.Profile.InvestmentReturn.Value)
	if err != nil {
		return nil, fmt.Errorf("投資利回りの復元に失敗しました: %w", err)
	}

	inflationRate, err := valueobjects.NewRate(dto.Profile.InflationRate.Value)
	if err != nil {
		return nil, fmt.Errorf("インフレ率の復元に失敗しました: %w", err)
	}

	profile, err := entities.NewFinancialProfileWithID(
		entities.FinancialProfileID(dto.Profile.ID),
		entities.UserID(dto.Profile.UserID),
		monthlyIncome,
		expenses,
		savings,
		investmentReturn,
		inflationRate,
		dto.Profile.CreatedAt,
		dto.Profile.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("財務プロファイルの復元に失敗しました: %w", err)
	}

	plan, err := aggregates.NewFinancialPlanWithID(
		aggregates.FinancialPlanID(dto.ID),
		profile,
		dto.CreatedAt,
		dto.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("財務計画の復元に失敗しました: %w", err)
	}

	// RetirementData を復元
	if dto.RetirementData != nil {
		rd := dto.RetirementData
		monthlyExpenses, err := valueobjects.NewMoney(rd.MonthlyRetirementExpenses.Amount, valueobjects.Currency(rd.MonthlyRetirementExpenses.Currency))
		if err != nil {
			return nil, fmt.Errorf("退職後月間支出の復元に失敗しました: %w", err)
		}
		pensionAmount, err := valueobjects.NewMoney(rd.PensionAmount.Amount, valueobjects.Currency(rd.PensionAmount.Currency))
		if err != nil {
			return nil, fmt.Errorf("年金額の復元に失敗しました: %w", err)
		}
		retirementData, err := entities.NewRetirementDataWithID(
			entities.RetirementDataID(rd.ID),
			entities.UserID(rd.UserID),
			rd.CurrentAge,
			rd.RetirementAge,
			rd.LifeExpectancy,
			monthlyExpenses,
			pensionAmount,
			rd.CreatedAt,
			rd.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("退職データの復元に失敗しました: %w", err)
		}
		if err := plan.SetRetirementData(retirementData); err != nil {
			return nil, fmt.Errorf("退職データの設定に失敗しました: %w", err)
		}
	}

	// EmergencyFund を復元
	if dto.EmergencyFund != nil {
		currentFund, err := valueobjects.NewMoney(dto.EmergencyFund.CurrentFund.Amount, valueobjects.Currency(dto.EmergencyFund.CurrentFund.Currency))
		if err != nil {
			return nil, fmt.Errorf("緊急資金の復元に失敗しました: %w", err)
		}
		efConfig, err := aggregates.NewEmergencyFundConfig(dto.EmergencyFund.TargetMonths, currentFund)
		if err != nil {
			return nil, fmt.Errorf("緊急資金設定の復元に失敗しました: %w", err)
		}
		if err := plan.UpdateEmergencyFund(efConfig); err != nil {
			return nil, fmt.Errorf("緊急資金設定の適用に失敗しました: %w", err)
		}
	}

	// Goals を復元
	goals, err := goalsFromDTOs(dto.Goals)
	if err != nil {
		return nil, fmt.Errorf("目標の復元に失敗しました: %w", err)
	}
	for _, goal := range goals {
		if err := plan.AddGoal(goal); err != nil {
			// ビジネスルール違反（達成不可能など）はキャッシュからの復元時には無視
			continue
		}
	}

	return plan, nil
}
