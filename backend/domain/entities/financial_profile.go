package entities

import (
	"errors"
	"fmt"
	"time"

	"github.com/financial-planning-calculator/backend/domain/valueobjects"

	"github.com/google/uuid"
)

// FinancialProfileID は財務プロファイルの一意識別子
type FinancialProfileID string

// NewFinancialProfileID は新しい財務プロファイルIDを生成する
func NewFinancialProfileID() FinancialProfileID {
	return FinancialProfileID(uuid.New().String())
}

// UserID はユーザーの一意識別子
type UserID string

// ExpenseItem は支出項目を表す
type ExpenseItem struct {
	Category    string             `json:"category"`
	Amount      valueobjects.Money `json:"amount"`
	Description string             `json:"description,omitempty"`
}

// ExpenseCollection は支出項目のコレクション
type ExpenseCollection []ExpenseItem

// Total は支出の合計金額を計算する
func (ec ExpenseCollection) Total() (valueobjects.Money, error) {
	if len(ec) == 0 {
		return valueobjects.NewMoneyJPY(0)
	}

	total, err := valueobjects.NewMoneyJPY(0)
	if err != nil {
		return valueobjects.Money{}, err
	}

	for _, expense := range ec {
		total, err = total.Add(expense.Amount)
		if err != nil {
			return valueobjects.Money{}, fmt.Errorf("支出合計の計算に失敗しました: %w", err)
		}
	}

	return total, nil
}

// GetByCategory は指定されたカテゴリの支出項目を取得する
func (ec ExpenseCollection) GetByCategory(category string) []ExpenseItem {
	var items []ExpenseItem
	for _, expense := range ec {
		if expense.Category == category {
			items = append(items, expense)
		}
	}
	return items
}

// SavingsItem は貯蓄項目を表す
type SavingsItem struct {
	Type        string             `json:"type"` // deposit, investment, other
	Amount      valueobjects.Money `json:"amount"`
	Description string             `json:"description,omitempty"`
}

// SavingsCollection は貯蓄項目のコレクション
type SavingsCollection []SavingsItem

// Total は貯蓄の合計金額を計算する
func (sc SavingsCollection) Total() (valueobjects.Money, error) {
	if len(sc) == 0 {
		return valueobjects.NewMoneyJPY(0)
	}

	total, err := valueobjects.NewMoneyJPY(0)
	if err != nil {
		return valueobjects.Money{}, err
	}

	for _, savings := range sc {
		total, err = total.Add(savings.Amount)
		if err != nil {
			return valueobjects.Money{}, fmt.Errorf("貯蓄合計の計算に失敗しました: %w", err)
		}
	}

	return total, nil
}

// GetByType は指定されたタイプの貯蓄項目を取得する
func (sc SavingsCollection) GetByType(savingsType string) []SavingsItem {
	var items []SavingsItem
	for _, savings := range sc {
		if savings.Type == savingsType {
			items = append(items, savings)
		}
	}
	return items
}

// AssetProjection は資産推移の予測データ
type AssetProjection struct {
	Year              int                `json:"year"`
	TotalAssets       valueobjects.Money `json:"total_assets"`
	RealValue         valueobjects.Money `json:"real_value"`
	ContributedAmount valueobjects.Money `json:"contributed_amount"`
	InvestmentGains   valueobjects.Money `json:"investment_gains"`
}

// FinancialProfile はユーザーの財務プロファイルを表すエンティティ
type FinancialProfile struct {
	id               FinancialProfileID
	userID           UserID
	monthlyIncome    valueobjects.Money
	monthlyExpenses  ExpenseCollection
	currentSavings   SavingsCollection
	investmentReturn valueobjects.Rate
	inflationRate    valueobjects.Rate
	createdAt        time.Time
	updatedAt        time.Time
}

// NewFinancialProfile は新しい財務プロファイルを作成する
func NewFinancialProfile(
	userID UserID,
	monthlyIncome valueobjects.Money,
	monthlyExpenses ExpenseCollection,
	currentSavings SavingsCollection,
	investmentReturn valueobjects.Rate,
	inflationRate valueobjects.Rate,
) (*FinancialProfile, error) {
	if userID == "" {
		return nil, errors.New("ユーザーIDは必須です")
	}

	if !monthlyIncome.IsPositive() {
		return nil, errors.New("月収は正の値である必要があります")
	}

	// 支出の合計を計算してバリデーション
	totalExpenses, err := monthlyExpenses.Total()
	if err != nil {
		return nil, fmt.Errorf("支出の合計計算に失敗しました: %w", err)
	}

	if totalExpenses.IsNegative() {
		return nil, errors.New("支出の合計は負の値にできません")
	}

	// 貯蓄の合計を計算してバリデーション
	totalSavings, err := currentSavings.Total()
	if err != nil {
		return nil, fmt.Errorf("貯蓄の合計計算に失敗しました: %w", err)
	}

	if totalSavings.IsNegative() {
		return nil, errors.New("貯蓄の合計は負の値にできません")
	}

	now := time.Now()

	return &FinancialProfile{
		id:               NewFinancialProfileID(),
		userID:           userID,
		monthlyIncome:    monthlyIncome,
		monthlyExpenses:  monthlyExpenses,
		currentSavings:   currentSavings,
		investmentReturn: investmentReturn,
		inflationRate:    inflationRate,
		createdAt:        now,
		updatedAt:        now,
	}, nil
}

// ID は財務プロファイルIDを返す
func (fp *FinancialProfile) ID() FinancialProfileID {
	return fp.id
}

// UserID はユーザーIDを返す
func (fp *FinancialProfile) UserID() UserID {
	return fp.userID
}

// MonthlyIncome は月収を返す
func (fp *FinancialProfile) MonthlyIncome() valueobjects.Money {
	return fp.monthlyIncome
}

// MonthlyExpenses は月間支出を返す
func (fp *FinancialProfile) MonthlyExpenses() ExpenseCollection {
	return fp.monthlyExpenses
}

// CurrentSavings は現在の貯蓄を返す
func (fp *FinancialProfile) CurrentSavings() SavingsCollection {
	return fp.currentSavings
}

// InvestmentReturn は投資利回りを返す
func (fp *FinancialProfile) InvestmentReturn() valueobjects.Rate {
	return fp.investmentReturn
}

// InflationRate はインフレ率を返す
func (fp *FinancialProfile) InflationRate() valueobjects.Rate {
	return fp.inflationRate
}

// CreatedAt は作成日時を返す
func (fp *FinancialProfile) CreatedAt() time.Time {
	return fp.createdAt
}

// UpdatedAt は更新日時を返す
func (fp *FinancialProfile) UpdatedAt() time.Time {
	return fp.updatedAt
}

// CalculateNetSavings は月間純貯蓄額を計算する（収入 - 支出）
func (fp *FinancialProfile) CalculateNetSavings() (valueobjects.Money, error) {
	totalExpenses, err := fp.monthlyExpenses.Total()
	if err != nil {
		return valueobjects.Money{}, fmt.Errorf("支出合計の計算に失敗しました: %w", err)
	}

	netSavings, err := fp.monthlyIncome.Subtract(totalExpenses)
	if err != nil {
		return valueobjects.Money{}, fmt.Errorf("純貯蓄額の計算に失敗しました: %w", err)
	}

	return netSavings, nil
}

// ValidateFinancialHealth は財務健全性をチェックする
func (fp *FinancialProfile) ValidateFinancialHealth() error {
	netSavings, err := fp.CalculateNetSavings()
	if err != nil {
		return fmt.Errorf("財務健全性の検証に失敗しました: %w", err)
	}

	// 純貯蓄額が負の場合は警告
	if netSavings.IsNegative() {
		return errors.New("月間支出が収入を上回っています。支出の見直しが必要です")
	}

	// 貯蓄率が低い場合の警告（収入の10%未満）
	savingsRate, err := netSavings.MultiplyByFloat(1.0)
	if err != nil {
		return fmt.Errorf("貯蓄率の計算に失敗しました: %w", err)
	}

	minimumSavingsTarget, err := fp.monthlyIncome.MultiplyByFloat(0.1) // 収入の10%
	if err != nil {
		return fmt.Errorf("最低貯蓄目標の計算に失敗しました: %w", err)
	}

	isLowSavings, err := savingsRate.LessThan(minimumSavingsTarget)
	if err != nil {
		return fmt.Errorf("貯蓄率の比較に失敗しました: %w", err)
	}

	if isLowSavings {
		return errors.New("貯蓄率が低すぎます。収入の10%以上の貯蓄を推奨します")
	}

	return nil
}

// ProjectAssets は指定年数の資産推移を予測する
func (fp *FinancialProfile) ProjectAssets(years int) ([]AssetProjection, error) {
	if years <= 0 {
		return nil, errors.New("予測年数は正の値である必要があります")
	}

	netSavings, err := fp.CalculateNetSavings()
	if err != nil {
		return nil, fmt.Errorf("純貯蓄額の計算に失敗しました: %w", err)
	}

	currentSavingsTotal, err := fp.currentSavings.Total()
	if err != nil {
		return nil, fmt.Errorf("現在の貯蓄合計の計算に失敗しました: %w", err)
	}

	projections := make([]AssetProjection, years)

	// 月利を計算
	monthlyInvestmentRate, err := fp.investmentReturn.MonthlyRate()
	if err != nil {
		return nil, fmt.Errorf("月利の計算に失敗しました: %w", err)
	}

	// 月間インフレ率を計算（後で使用）
	_, err = fp.inflationRate.MonthlyRate()
	if err != nil {
		return nil, fmt.Errorf("月間インフレ率の計算に失敗しました: %w", err)
	}

	currentAssets := currentSavingsTotal
	totalContributed := currentSavingsTotal

	for year := 1; year <= years; year++ {
		// 年間の複利計算
		for month := 1; month <= 12; month++ {
			// 投資収益を加算
			investmentGain, err := currentAssets.Multiply(monthlyInvestmentRate)
			if err != nil {
				return nil, fmt.Errorf("投資収益の計算に失敗しました: %w", err)
			}

			currentAssets, err = currentAssets.Add(investmentGain)
			if err != nil {
				return nil, fmt.Errorf("資産への投資収益加算に失敗しました: %w", err)
			}

			// 月間貯蓄を加算
			currentAssets, err = currentAssets.Add(netSavings)
			if err != nil {
				return nil, fmt.Errorf("資産への月間貯蓄加算に失敗しました: %w", err)
			}

			totalContributed, err = totalContributed.Add(netSavings)
			if err != nil {
				return nil, fmt.Errorf("総拠出額の計算に失敗しました: %w", err)
			}
		}

		// 投資収益を計算
		investmentGains, err := currentAssets.Subtract(totalContributed)
		if err != nil {
			return nil, fmt.Errorf("投資収益の計算に失敗しました: %w", err)
		}

		// インフレ調整後の実質価値を計算
		inflationFactor := fp.inflationRate.CompoundFactor(year)
		realValue, err := currentAssets.MultiplyByFloat(1.0 / inflationFactor)
		if err != nil {
			return nil, fmt.Errorf("実質価値の計算に失敗しました: %w", err)
		}

		projections[year-1] = AssetProjection{
			Year:              year,
			TotalAssets:       currentAssets,
			RealValue:         realValue,
			ContributedAmount: totalContributed,
			InvestmentGains:   investmentGains,
		}
	}

	return projections, nil
}

// UpdateMonthlyIncome は月収を更新する
func (fp *FinancialProfile) UpdateMonthlyIncome(newIncome valueobjects.Money) error {
	if !newIncome.IsPositive() {
		return errors.New("月収は正の値である必要があります")
	}

	fp.monthlyIncome = newIncome
	fp.updatedAt = time.Now()
	return nil
}

// UpdateMonthlyExpenses は月間支出を更新する
func (fp *FinancialProfile) UpdateMonthlyExpenses(newExpenses ExpenseCollection) error {
	totalExpenses, err := newExpenses.Total()
	if err != nil {
		return fmt.Errorf("支出の合計計算に失敗しました: %w", err)
	}

	if totalExpenses.IsNegative() {
		return errors.New("支出の合計は負の値にできません")
	}

	fp.monthlyExpenses = newExpenses
	fp.updatedAt = time.Now()
	return nil
}

// UpdateCurrentSavings は現在の貯蓄を更新する
func (fp *FinancialProfile) UpdateCurrentSavings(newSavings SavingsCollection) error {
	totalSavings, err := newSavings.Total()
	if err != nil {
		return fmt.Errorf("貯蓄の合計計算に失敗しました: %w", err)
	}

	if totalSavings.IsNegative() {
		return errors.New("貯蓄の合計は負の値にできません")
	}

	fp.currentSavings = newSavings
	fp.updatedAt = time.Now()
	return nil
}

// UpdateInvestmentReturn は投資利回りを更新する
func (fp *FinancialProfile) UpdateInvestmentReturn(newRate valueobjects.Rate) error {
	fp.investmentReturn = newRate
	fp.updatedAt = time.Now()
	return nil
}

// UpdateInflationRate はインフレ率を更新する
func (fp *FinancialProfile) UpdateInflationRate(newRate valueobjects.Rate) error {
	fp.inflationRate = newRate
	fp.updatedAt = time.Now()
	return nil
}
