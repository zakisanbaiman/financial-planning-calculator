package entities

import (
	"errors"
	"fmt"
	"time"

	"github.com/financial-planning-calculator/backend/domain/valueobjects"

	"github.com/google/uuid"
)

// RetirementDataID は退職データの一意識別子
type RetirementDataID string

// NewRetirementDataID は新しい退職データIDを生成する
func NewRetirementDataID() RetirementDataID {
	return RetirementDataID(uuid.New().String())
}

// RetirementCalculation は老後資金計算結果を表す
type RetirementCalculation struct {
	RequiredAmount            valueobjects.Money `json:"required_amount"`             // 必要老後資金
	ProjectedAmount           valueobjects.Money `json:"projected_amount"`            // 予想達成額
	Shortfall                 valueobjects.Money `json:"shortfall"`                   // 不足額
	SufficiencyRate           valueobjects.Rate  `json:"sufficiency_rate"`            // 充足率 (%)
	RecommendedMonthlySavings valueobjects.Money `json:"recommended_monthly_savings"` // 推奨月間貯蓄額
}

// RetirementData は退職・年金情報を表すエンティティ
type RetirementData struct {
	id                        RetirementDataID
	userID                    UserID
	currentAge                int
	retirementAge             int
	lifeExpectancy            int
	monthlyRetirementExpenses valueobjects.Money
	pensionAmount             valueobjects.Money
	createdAt                 time.Time
	updatedAt                 time.Time
}

// NewRetirementData は新しい退職データを作成する
func NewRetirementData(
	userID UserID,
	currentAge int,
	retirementAge int,
	lifeExpectancy int,
	monthlyRetirementExpenses valueobjects.Money,
	pensionAmount valueobjects.Money,
) (*RetirementData, error) {
	if userID == "" {
		return nil, errors.New("ユーザーIDは必須です")
	}

	if currentAge < 0 || currentAge > 150 {
		return nil, errors.New("現在の年齢は0歳から150歳の間である必要があります")
	}

	if retirementAge < currentAge {
		return nil, errors.New("退職年齢は現在の年齢以上である必要があります")
	}

	if retirementAge > 100 {
		return nil, errors.New("退職年齢は100歳以下である必要があります")
	}

	if lifeExpectancy < retirementAge {
		return nil, errors.New("平均寿命は退職年齢以上である必要があります")
	}

	if lifeExpectancy > 150 {
		return nil, errors.New("平均寿命は150歳以下である必要があります")
	}

	if monthlyRetirementExpenses.IsNegative() {
		return nil, errors.New("月間退職後支出は負の値にできません")
	}

	if pensionAmount.IsNegative() {
		return nil, errors.New("年金額は負の値にできません")
	}

	now := time.Now()

	return &RetirementData{
		id:                        NewRetirementDataID(),
		userID:                    userID,
		currentAge:                currentAge,
		retirementAge:             retirementAge,
		lifeExpectancy:            lifeExpectancy,
		monthlyRetirementExpenses: monthlyRetirementExpenses,
		pensionAmount:             pensionAmount,
		createdAt:                 now,
		updatedAt:                 now,
	}, nil
}

// ID は退職データIDを返す
func (rd *RetirementData) ID() RetirementDataID {
	return rd.id
}

// UserID はユーザーIDを返す
func (rd *RetirementData) UserID() UserID {
	return rd.userID
}

// CurrentAge は現在の年齢を返す
func (rd *RetirementData) CurrentAge() int {
	return rd.currentAge
}

// RetirementAge は退職年齢を返す
func (rd *RetirementData) RetirementAge() int {
	return rd.retirementAge
}

// LifeExpectancy は平均寿命を返す
func (rd *RetirementData) LifeExpectancy() int {
	return rd.lifeExpectancy
}

// MonthlyRetirementExpenses は月間退職後支出を返す
func (rd *RetirementData) MonthlyRetirementExpenses() valueobjects.Money {
	return rd.monthlyRetirementExpenses
}

// PensionAmount は年金額を返す
func (rd *RetirementData) PensionAmount() valueobjects.Money {
	return rd.pensionAmount
}

// CreatedAt は作成日時を返す
func (rd *RetirementData) CreatedAt() time.Time {
	return rd.createdAt
}

// UpdatedAt は更新日時を返す
func (rd *RetirementData) UpdatedAt() time.Time {
	return rd.updatedAt
}

// CalculateYearsUntilRetirement は退職までの年数を計算する
func (rd *RetirementData) CalculateYearsUntilRetirement() int {
	yearsUntilRetirement := rd.retirementAge - rd.currentAge
	if yearsUntilRetirement < 0 {
		return 0
	}
	return yearsUntilRetirement
}

// CalculateRetirementYears は退職後の年数を計算する
func (rd *RetirementData) CalculateRetirementYears() int {
	retirementYears := rd.lifeExpectancy - rd.retirementAge
	if retirementYears < 0 {
		return 0
	}
	return retirementYears
}

// CalculateRequiredRetirementFund は必要な老後資金を計算する
func (rd *RetirementData) CalculateRequiredRetirementFund(inflationRate valueobjects.Rate) (valueobjects.Money, error) {
	retirementYears := rd.CalculateRetirementYears()
	if retirementYears <= 0 {
		return valueobjects.NewMoneyJPY(0)
	}

	// 年金で不足する月額を計算
	monthlyShortfall, err := rd.monthlyRetirementExpenses.Subtract(rd.pensionAmount)
	if err != nil {
		return valueobjects.Money{}, fmt.Errorf("月間不足額の計算に失敗しました: %w", err)
	}

	// 年金で足りている場合は0を返す
	if monthlyShortfall.IsNegative() || monthlyShortfall.IsZero() {
		return valueobjects.NewMoneyJPY(0)
	}

	// 退職時点でのインフレ調整
	yearsUntilRetirement := rd.CalculateYearsUntilRetirement()
	inflationFactor := inflationRate.CompoundFactor(yearsUntilRetirement)

	adjustedMonthlyShortfall, err := monthlyShortfall.MultiplyByFloat(inflationFactor)
	if err != nil {
		return valueobjects.Money{}, fmt.Errorf("インフレ調整に失敗しました: %w", err)
	}

	// 退職後の総必要額を計算（月額 × 12ヶ月 × 退職後年数）
	totalMonths := retirementYears * 12
	requiredFund, err := adjustedMonthlyShortfall.MultiplyByFloat(float64(totalMonths))
	if err != nil {
		return valueobjects.Money{}, fmt.Errorf("必要老後資金の計算に失敗しました: %w", err)
	}

	return requiredFund, nil
}

// CalculateRetirementSufficiency は老後資金の充足度を計算する
func (rd *RetirementData) CalculateRetirementSufficiency(
	currentSavings valueobjects.Money,
	monthlySavings valueobjects.Money,
	investmentReturn valueobjects.Rate,
	inflationRate valueobjects.Rate,
) (*RetirementCalculation, error) {
	// 必要老後資金を計算
	requiredAmount, err := rd.CalculateRequiredRetirementFund(inflationRate)
	if err != nil {
		return nil, fmt.Errorf("必要老後資金の計算に失敗しました: %w", err)
	}

	// 退職時点での予想資産額を計算
	yearsUntilRetirement := rd.CalculateYearsUntilRetirement()
	projectedAmount, err := rd.calculateProjectedAssets(currentSavings, monthlySavings, investmentReturn, yearsUntilRetirement)
	if err != nil {
		return nil, fmt.Errorf("予想資産額の計算に失敗しました: %w", err)
	}

	// 不足額を計算
	shortfall, err := requiredAmount.Subtract(projectedAmount)
	if err != nil {
		return nil, fmt.Errorf("不足額の計算に失敗しました: %w", err)
	}

	// 不足額が負の場合（余剰がある場合）は0にする
	if shortfall.IsNegative() {
		shortfall, _ = valueobjects.NewMoneyJPY(0)
	}

	// 充足率を計算
	var sufficiencyRate valueobjects.Rate
	if requiredAmount.IsZero() {
		sufficiencyRate, _ = valueobjects.NewRate(100.0) // 必要額が0の場合は100%
	} else {
		sufficiencyPercentage := (projectedAmount.Amount() / requiredAmount.Amount()) * 100
		if sufficiencyPercentage > 100 {
			sufficiencyPercentage = 100
		}
		sufficiencyRate, err = valueobjects.NewRate(sufficiencyPercentage)
		if err != nil {
			return nil, fmt.Errorf("充足率の計算に失敗しました: %w", err)
		}
	}

	// 推奨月間貯蓄額を計算
	recommendedMonthlySavings, err := rd.calculateRecommendedMonthlySavings(
		currentSavings, requiredAmount, investmentReturn, yearsUntilRetirement)
	if err != nil {
		return nil, fmt.Errorf("推奨月間貯蓄額の計算に失敗しました: %w", err)
	}

	return &RetirementCalculation{
		RequiredAmount:            requiredAmount,
		ProjectedAmount:           projectedAmount,
		Shortfall:                 shortfall,
		SufficiencyRate:           sufficiencyRate,
		RecommendedMonthlySavings: recommendedMonthlySavings,
	}, nil
}

// calculateProjectedAssets は退職時点での予想資産額を計算する
func (rd *RetirementData) calculateProjectedAssets(
	currentSavings valueobjects.Money,
	monthlySavings valueobjects.Money,
	investmentReturn valueobjects.Rate,
	years int,
) (valueobjects.Money, error) {
	if years <= 0 {
		return currentSavings, nil
	}

	// 月利を計算
	monthlyRate, err := investmentReturn.MonthlyRate()
	if err != nil {
		return valueobjects.Money{}, fmt.Errorf("月利の計算に失敗しました: %w", err)
	}

	currentAssets := currentSavings
	totalMonths := years * 12

	// 複利計算
	for month := 0; month < totalMonths; month++ {
		// 投資収益を加算
		investmentGain, err := currentAssets.Multiply(monthlyRate)
		if err != nil {
			return valueobjects.Money{}, fmt.Errorf("投資収益の計算に失敗しました: %w", err)
		}

		currentAssets, err = currentAssets.Add(investmentGain)
		if err != nil {
			return valueobjects.Money{}, fmt.Errorf("資産への投資収益加算に失敗しました: %w", err)
		}

		// 月間貯蓄を加算
		currentAssets, err = currentAssets.Add(monthlySavings)
		if err != nil {
			return valueobjects.Money{}, fmt.Errorf("資産への月間貯蓄加算に失敗しました: %w", err)
		}
	}

	return currentAssets, nil
}

// calculateRecommendedMonthlySavings は推奨月間貯蓄額を計算する
func (rd *RetirementData) calculateRecommendedMonthlySavings(
	currentSavings valueobjects.Money,
	requiredAmount valueobjects.Money,
	investmentReturn valueobjects.Rate,
	years int,
) (valueobjects.Money, error) {
	if years <= 0 {
		// 退職まで時間がない場合は、不足額をそのまま返す
		shortfall, err := requiredAmount.Subtract(currentSavings)
		if err != nil {
			return valueobjects.Money{}, err
		}
		if shortfall.IsNegative() {
			return valueobjects.NewMoneyJPY(0)
		}
		return shortfall, nil
	}

	// 現在の資産が投資収益のみで成長した場合の将来価値
	compoundFactor := investmentReturn.CompoundFactor(years)
	futureValueOfCurrentSavings, err := currentSavings.MultiplyByFloat(compoundFactor)
	if err != nil {
		// 投資収益なしの場合
		futureValueOfCurrentSavings = currentSavings
	}

	// 必要な追加資金
	additionalRequired, err := requiredAmount.Subtract(futureValueOfCurrentSavings)
	if err != nil {
		return valueobjects.Money{}, fmt.Errorf("追加必要資金の計算に失敗しました: %w", err)
	}

	// 既に十分な資産がある場合
	if additionalRequired.IsNegative() || additionalRequired.IsZero() {
		return valueobjects.NewMoneyJPY(0)
	}

	// 年金現価係数を使用して月間貯蓄額を計算
	// 簡略化のため、単純に期間で割る
	totalMonths := years * 12
	recommendedMonthlySavings := additionalRequired.Amount() / float64(totalMonths)

	return valueobjects.NewMoneyJPY(recommendedMonthlySavings)
}

// UpdateCurrentAge は現在の年齢を更新する
func (rd *RetirementData) UpdateCurrentAge(newAge int) error {
	if newAge < 0 || newAge > 150 {
		return errors.New("年齢は0歳から150歳の間である必要があります")
	}

	if newAge > rd.retirementAge {
		return errors.New("現在の年齢は退職年齢以下である必要があります")
	}

	rd.currentAge = newAge
	rd.updatedAt = time.Now()
	return nil
}

// UpdateRetirementAge は退職年齢を更新する
func (rd *RetirementData) UpdateRetirementAge(newAge int) error {
	if newAge < rd.currentAge {
		return errors.New("退職年齢は現在の年齢以上である必要があります")
	}

	if newAge > 100 {
		return errors.New("退職年齢は100歳以下である必要があります")
	}

	if newAge > rd.lifeExpectancy {
		return errors.New("退職年齢は平均寿命以下である必要があります")
	}

	rd.retirementAge = newAge
	rd.updatedAt = time.Now()
	return nil
}

// UpdateLifeExpectancy は平均寿命を更新する
func (rd *RetirementData) UpdateLifeExpectancy(newAge int) error {
	if newAge < rd.retirementAge {
		return errors.New("平均寿命は退職年齢以上である必要があります")
	}

	if newAge > 150 {
		return errors.New("平均寿命は150歳以下である必要があります")
	}

	rd.lifeExpectancy = newAge
	rd.updatedAt = time.Now()
	return nil
}

// UpdateMonthlyRetirementExpenses は月間退職後支出を更新する
func (rd *RetirementData) UpdateMonthlyRetirementExpenses(newExpenses valueobjects.Money) error {
	if newExpenses.IsNegative() {
		return errors.New("月間退職後支出は負の値にできません")
	}

	rd.monthlyRetirementExpenses = newExpenses
	rd.updatedAt = time.Now()
	return nil
}

// UpdatePensionAmount は年金額を更新する
func (rd *RetirementData) UpdatePensionAmount(newAmount valueobjects.Money) error {
	if newAmount.IsNegative() {
		return errors.New("年金額は負の値にできません")
	}

	rd.pensionAmount = newAmount
	rd.updatedAt = time.Now()
	return nil
}

// IsRetired は現在退職しているかどうかを返す
func (rd *RetirementData) IsRetired() bool {
	return rd.currentAge >= rd.retirementAge
}

// GetPensionShortfall は年金の不足額を返す
func (rd *RetirementData) GetPensionShortfall() (valueobjects.Money, error) {
	shortfall, err := rd.monthlyRetirementExpenses.Subtract(rd.pensionAmount)
	if err != nil {
		return valueobjects.Money{}, fmt.Errorf("年金不足額の計算に失敗しました: %w", err)
	}

	// 不足がない場合は0を返す
	if shortfall.IsNegative() {
		return valueobjects.NewMoneyJPY(0)
	}

	return shortfall, nil
}

// IsPensionSufficient は年金が十分かどうかを返す
func (rd *RetirementData) IsPensionSufficient() (bool, error) {
	shortfall, err := rd.GetPensionShortfall()
	if err != nil {
		return false, err
	}

	return shortfall.IsZero(), nil
}
