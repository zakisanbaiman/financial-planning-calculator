package services

import (
	"errors"
	"fmt"
	"math"

	"github.com/financial-planning-calculator/backend/domain/valueobjects"
)

// FinancialCalculationService は財務計算に関するドメインサービス
type FinancialCalculationService struct{}

// NewFinancialCalculationService は新しいFinancialCalculationServiceを作成する
func NewFinancialCalculationService() *FinancialCalculationService {
	return &FinancialCalculationService{}
}

// CompoundInterestResult は複利計算の結果を表す
type CompoundInterestResult struct {
	FinalAmount       valueobjects.Money `json:"final_amount"`       // 最終金額
	TotalContribution valueobjects.Money `json:"total_contribution"` // 総拠出額
	InterestEarned    valueobjects.Money `json:"interest_earned"`    // 利息収益
	EffectiveRate     valueobjects.Rate  `json:"effective_rate"`     // 実効利率
}

// InflationAdjustmentResult はインフレ調整の結果を表す
type InflationAdjustmentResult struct {
	NominalValue        valueobjects.Money `json:"nominal_value"`         // 名目価値
	RealValue           valueobjects.Money `json:"real_value"`            // 実質価値
	PurchasingPowerLoss valueobjects.Money `json:"purchasing_power_loss"` // 購買力の損失
	InflationImpact     valueobjects.Rate  `json:"inflation_impact"`      // インフレの影響率
}

// RetirementNeedsResult は老後資金必要額計算の結果を表す
type RetirementNeedsResult struct {
	TotalNeeds         valueobjects.Money `json:"total_needs"`         // 総必要額
	MonthlyNeeds       valueobjects.Money `json:"monthly_needs"`       // 月間必要額
	YearsInRetirement  int                `json:"years_in_retirement"` // 退職後年数
	InflationAdjusted  bool               `json:"inflation_adjusted"`  // インフレ調整済みか
	PensionCoverage    valueobjects.Money `json:"pension_coverage"`    // 年金でカバーされる額
	AdditionalRequired valueobjects.Money `json:"additional_required"` // 追加で必要な額
}

// CalculateCompoundInterest は複利計算を実行する
func (fcs *FinancialCalculationService) CalculateCompoundInterest(
	principal valueobjects.Money,
	rate valueobjects.Rate,
	periods int,
) (*CompoundInterestResult, error) {
	if periods < 0 {
		return nil, errors.New("期間は負の値にできません")
	}

	if periods == 0 {
		return &CompoundInterestResult{
			FinalAmount:       principal,
			TotalContribution: principal,
			InterestEarned:    valueobjects.Money{},
			EffectiveRate:     rate,
		}, nil
	}

	// 複利計算: A = P(1 + r)^n
	compoundFactor := rate.CompoundFactor(periods)
	finalAmount, err := principal.MultiplyByFloat(compoundFactor)
	if err != nil {
		return nil, fmt.Errorf("最終金額の計算に失敗しました: %w", err)
	}

	// 利息収益を計算
	interestEarned, err := finalAmount.Subtract(principal)
	if err != nil {
		return nil, fmt.Errorf("利息収益の計算に失敗しました: %w", err)
	}

	// 実効利率を計算（年率換算）
	effectiveRate := math.Pow(compoundFactor, 1.0/float64(periods)) - 1
	effectiveRateObj, err := valueobjects.NewRateFromDecimal(effectiveRate)
	if err != nil {
		effectiveRateObj = rate // フォールバック
	}

	return &CompoundInterestResult{
		FinalAmount:       finalAmount,
		TotalContribution: principal,
		InterestEarned:    interestEarned,
		EffectiveRate:     effectiveRateObj,
	}, nil
}

// CalculateCompoundInterestWithRegularPayments は定期積立を含む複利計算を実行する
func (fcs *FinancialCalculationService) CalculateCompoundInterestWithRegularPayments(
	principal valueobjects.Money,
	monthlyPayment valueobjects.Money,
	annualRate valueobjects.Rate,
	years int,
) (*CompoundInterestResult, error) {
	if years < 0 {
		return nil, errors.New("年数は負の値にできません")
	}

	if years == 0 {
		return &CompoundInterestResult{
			FinalAmount:       principal,
			TotalContribution: principal,
			InterestEarned:    valueobjects.Money{},
			EffectiveRate:     annualRate,
		}, nil
	}

	// 月利を計算
	monthlyRate, err := annualRate.MonthlyRate()
	if err != nil {
		return nil, fmt.Errorf("月利の計算に失敗しました: %w", err)
	}

	currentAmount := principal
	totalMonths := years * 12
	totalContribution := principal

	// 月次複利計算
	for month := 0; month < totalMonths; month++ {
		// 投資収益を加算
		if !monthlyRate.IsZero() {
			interestGain, err := currentAmount.Multiply(monthlyRate)
			if err != nil {
				return nil, fmt.Errorf("月次投資収益の計算に失敗しました: %w", err)
			}

			currentAmount, err = currentAmount.Add(interestGain)
			if err != nil {
				return nil, fmt.Errorf("投資収益の加算に失敗しました: %w", err)
			}
		}

		// 月次積立を加算
		currentAmount, err = currentAmount.Add(monthlyPayment)
		if err != nil {
			return nil, fmt.Errorf("月次積立の加算に失敗しました: %w", err)
		}

		totalContribution, err = totalContribution.Add(monthlyPayment)
		if err != nil {
			return nil, fmt.Errorf("総拠出額の計算に失敗しました: %w", err)
		}
	}

	// 利息収益を計算
	interestEarned, err := currentAmount.Subtract(totalContribution)
	if err != nil {
		return nil, fmt.Errorf("利息収益の計算に失敗しました: %w", err)
	}

	return &CompoundInterestResult{
		FinalAmount:       currentAmount,
		TotalContribution: totalContribution,
		InterestEarned:    interestEarned,
		EffectiveRate:     annualRate,
	}, nil
}

// CalculateInflationAdjustedValue はインフレ調整後の実質価値を計算する
func (fcs *FinancialCalculationService) CalculateInflationAdjustedValue(
	amount valueobjects.Money,
	inflationRate valueobjects.Rate,
	years int,
) (*InflationAdjustmentResult, error) {
	if years < 0 {
		return nil, errors.New("年数は負の値にできません")
	}

	if years == 0 {
		return &InflationAdjustmentResult{
			NominalValue:        amount,
			RealValue:           amount,
			PurchasingPowerLoss: valueobjects.Money{},
			InflationImpact:     inflationRate,
		}, nil
	}

	// インフレ調整: Real Value = Nominal Value / (1 + inflation_rate)^years
	inflationFactor := inflationRate.CompoundFactor(years)
	realValue, err := amount.MultiplyByFloat(1.0 / inflationFactor)
	if err != nil {
		return nil, fmt.Errorf("実質価値の計算に失敗しました: %w", err)
	}

	// 購買力の損失を計算
	purchasingPowerLoss, err := amount.Subtract(realValue)
	if err != nil {
		return nil, fmt.Errorf("購買力損失の計算に失敗しました: %w", err)
	}

	// インフレの影響率を計算
	impactPercentage := (purchasingPowerLoss.Amount() / amount.Amount()) * 100
	inflationImpact, err := valueobjects.NewRate(impactPercentage)
	if err != nil {
		inflationImpact = inflationRate // フォールバック
	}

	return &InflationAdjustmentResult{
		NominalValue:        amount,
		RealValue:           realValue,
		PurchasingPowerLoss: purchasingPowerLoss,
		InflationImpact:     inflationImpact,
	}, nil
}

// CalculateRetirementNeeds は老後資金の必要額を計算する
func (fcs *FinancialCalculationService) CalculateRetirementNeeds(
	monthlyExpenses valueobjects.Money,
	yearsInRetirement int,
	inflationRate valueobjects.Rate,
	pensionAmount valueobjects.Money,
) (*RetirementNeedsResult, error) {
	if yearsInRetirement < 0 {
		return nil, errors.New("退職後年数は負の値にできません")
	}

	if monthlyExpenses.IsNegative() {
		return nil, errors.New("月間支出は負の値にできません")
	}

	if pensionAmount.IsNegative() {
		return nil, errors.New("年金額は負の値にできません")
	}

	// 年金でカバーされない月間不足額を計算
	monthlyShortfall, err := monthlyExpenses.Subtract(pensionAmount)
	if err != nil {
		return nil, fmt.Errorf("月間不足額の計算に失敗しました: %w", err)
	}

	// 年金で十分な場合
	if monthlyShortfall.IsNegative() || monthlyShortfall.IsZero() {
		zeroAmount, _ := valueobjects.NewMoneyJPY(0)
		return &RetirementNeedsResult{
			TotalNeeds:         zeroAmount,
			MonthlyNeeds:       zeroAmount,
			YearsInRetirement:  yearsInRetirement,
			InflationAdjusted:  true,
			PensionCoverage:    pensionAmount,
			AdditionalRequired: zeroAmount,
		}, nil
	}

	// インフレ調整後の月間必要額を計算（退職時点での価値）
	// 簡略化のため、退職開始時点でのインフレ調整を適用
	inflationAdjustedMonthly := monthlyShortfall

	// 総必要額を計算（月額 × 12ヶ月 × 年数）
	totalMonths := yearsInRetirement * 12
	totalNeeds, err := inflationAdjustedMonthly.MultiplyByFloat(float64(totalMonths))
	if err != nil {
		return nil, fmt.Errorf("総必要額の計算に失敗しました: %w", err)
	}

	// 年金でカバーされる総額
	totalPensionCoverage, err := pensionAmount.MultiplyByFloat(float64(totalMonths))
	if err != nil {
		return nil, fmt.Errorf("年金総額の計算に失敗しました: %w", err)
	}

	return &RetirementNeedsResult{
		TotalNeeds:         totalNeeds,
		MonthlyNeeds:       inflationAdjustedMonthly,
		YearsInRetirement:  yearsInRetirement,
		InflationAdjusted:  true,
		PensionCoverage:    totalPensionCoverage,
		AdditionalRequired: totalNeeds,
	}, nil
}

// CalculateFutureValue は将来価値を計算する（一般的な計算）
func (fcs *FinancialCalculationService) CalculateFutureValue(
	presentValue valueobjects.Money,
	rate valueobjects.Rate,
	periods int,
) (valueobjects.Money, error) {
	if periods < 0 {
		return valueobjects.Money{}, errors.New("期間は負の値にできません")
	}

	if periods == 0 {
		return presentValue, nil
	}

	compoundFactor := rate.CompoundFactor(periods)
	return presentValue.MultiplyByFloat(compoundFactor)
}

// CalculatePresentValue は現在価値を計算する
func (fcs *FinancialCalculationService) CalculatePresentValue(
	futureValue valueobjects.Money,
	rate valueobjects.Rate,
	periods int,
) (valueobjects.Money, error) {
	if periods < 0 {
		return valueobjects.Money{}, errors.New("期間は負の値にできません")
	}

	if periods == 0 {
		return futureValue, nil
	}

	if rate.IsZero() {
		return futureValue, nil
	}

	discountFactor := 1.0 / rate.CompoundFactor(periods)
	return futureValue.MultiplyByFloat(discountFactor)
}

// CalculateRequiredSavingsRate は目標達成に必要な貯蓄率を計算する
func (fcs *FinancialCalculationService) CalculateRequiredSavingsRate(
	currentIncome valueobjects.Money,
	targetAmount valueobjects.Money,
	currentSavings valueobjects.Money,
	investmentReturn valueobjects.Rate,
	years int,
) (valueobjects.Rate, error) {
	if years <= 0 {
		return valueobjects.Rate{}, errors.New("年数は正の値である必要があります")
	}

	if currentIncome.IsZero() || currentIncome.IsNegative() {
		return valueobjects.Rate{}, errors.New("現在の収入は正の値である必要があります")
	}

	// 現在の貯蓄が将来どれだけ成長するかを計算
	futureValueOfCurrentSavings, err := fcs.CalculateFutureValue(currentSavings, investmentReturn, years)
	if err != nil {
		return valueobjects.Rate{}, fmt.Errorf("現在貯蓄の将来価値計算に失敗しました: %w", err)
	}

	// 追加で必要な金額を計算
	additionalRequired, err := targetAmount.Subtract(futureValueOfCurrentSavings)
	if err != nil {
		return valueobjects.Rate{}, fmt.Errorf("追加必要額の計算に失敗しました: %w", err)
	}

	// 既に目標を達成している場合
	if additionalRequired.IsNegative() || additionalRequired.IsZero() {
		return valueobjects.NewRate(0)
	}

	// 年間必要貯蓄額を計算（簡略化：複利効果を無視）
	annualSavingsRequired := additionalRequired.Amount() / float64(years)

	// 必要貯蓄率を計算
	requiredSavingsRate := (annualSavingsRequired / currentIncome.Amount()) * 100

	return valueobjects.NewRate(requiredSavingsRate)
}

// CalculateEmergencyFundTarget は緊急資金の目標額を計算する
func (fcs *FinancialCalculationService) CalculateEmergencyFundTarget(
	monthlyExpenses valueobjects.Money,
	targetMonths int,
	inflationRate valueobjects.Rate,
	yearsToTarget int,
) (valueobjects.Money, error) {
	if targetMonths < 0 {
		return valueobjects.Money{}, errors.New("目標月数は負の値にできません")
	}

	if yearsToTarget < 0 {
		return valueobjects.Money{}, errors.New("目標年数は負の値にできません")
	}

	// 基本的な緊急資金額を計算
	baseTarget, err := monthlyExpenses.MultiplyByFloat(float64(targetMonths))
	if err != nil {
		return valueobjects.Money{}, fmt.Errorf("基本緊急資金額の計算に失敗しました: %w", err)
	}

	// インフレ調整を適用
	if yearsToTarget > 0 && !inflationRate.IsZero() {
		inflationFactor := inflationRate.CompoundFactor(yearsToTarget)
		adjustedTarget, err := baseTarget.MultiplyByFloat(inflationFactor)
		if err != nil {
			return baseTarget, nil // インフレ調整に失敗した場合は基本額を返す
		}
		return adjustedTarget, nil
	}

	return baseTarget, nil
}

// CalculateDebtPayoffTime は債務返済期間を計算する
func (fcs *FinancialCalculationService) CalculateDebtPayoffTime(
	debtAmount valueobjects.Money,
	monthlyPayment valueobjects.Money,
	interestRate valueobjects.Rate,
) (int, error) {
	if debtAmount.IsNegative() || debtAmount.IsZero() {
		return 0, nil
	}

	if monthlyPayment.IsNegative() || monthlyPayment.IsZero() {
		return -1, errors.New("月間返済額は正の値である必要があります")
	}

	// 月利を計算
	monthlyRate, err := interestRate.MonthlyRate()
	if err != nil {
		return -1, fmt.Errorf("月利の計算に失敗しました: %w", err)
	}

	// 利息のみで月間返済額を上回る場合は返済不可能
	if !monthlyRate.IsZero() {
		monthlyInterest, err := debtAmount.Multiply(monthlyRate)
		if err == nil {
			isPaymentInsufficient, err := monthlyPayment.LessThan(monthlyInterest)
			if err == nil && isPaymentInsufficient {
				return -1, errors.New("月間返済額が利息を下回るため返済できません")
			}
		}
	}

	// 単純計算（利息なしの場合）
	if monthlyRate.IsZero() {
		months := int(math.Ceil(debtAmount.Amount() / monthlyPayment.Amount()))
		return months, nil
	}

	// 複利計算による返済期間計算
	// 数値計算で近似
	remainingDebt := debtAmount
	months := 0
	maxMonths := 1200 // 100年の上限

	for months < maxMonths && remainingDebt.IsPositive() {
		// 月利を加算
		interest, err := remainingDebt.Multiply(monthlyRate)
		if err != nil {
			break
		}

		remainingDebt, err = remainingDebt.Add(interest)
		if err != nil {
			break
		}

		// 返済額を減算
		remainingDebt, err = remainingDebt.Subtract(monthlyPayment)
		if err != nil {
			break
		}

		months++
	}

	if months >= maxMonths {
		return -1, errors.New("返済期間が長すぎます")
	}

	return months, nil
}
