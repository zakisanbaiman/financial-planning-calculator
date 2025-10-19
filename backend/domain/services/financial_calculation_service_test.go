package services

import (
	"testing"

	"github.com/financial-planning-calculator/backend/domain/valueobjects"
)

func TestCalculateCompoundInterest(t *testing.T) {
	service := NewFinancialCalculationService()

	// テストケース: 100万円を年利5%で10年間運用
	principal, _ := valueobjects.NewMoneyJPY(1000000)
	rate, _ := valueobjects.NewRate(5.0)
	periods := 10

	result, err := service.CalculateCompoundInterest(principal, rate, periods)
	if err != nil {
		t.Fatalf("複利計算に失敗しました: %v", err)
	}

	// 検証: 10年後の金額は約162万円になるはず
	expectedAmount := 1000000 * 1.6289 // (1.05)^10 ≈ 1.6289
	if result.FinalAmount.Amount() < expectedAmount*0.99 || result.FinalAmount.Amount() > expectedAmount*1.01 {
		t.Errorf("最終金額が期待値と異なります。期待値: %.0f, 実際: %.0f", expectedAmount, result.FinalAmount.Amount())
	}

	// 利息収益の検証
	expectedInterest := expectedAmount - 1000000
	if result.InterestEarned.Amount() < expectedInterest*0.99 || result.InterestEarned.Amount() > expectedInterest*1.01 {
		t.Errorf("利息収益が期待値と異なります。期待値: %.0f, 実際: %.0f", expectedInterest, result.InterestEarned.Amount())
	}
}

func TestCalculateCompoundInterestWithRegularPayments(t *testing.T) {
	service := NewFinancialCalculationService()

	// テストケース: 初期100万円 + 月5万円積立を年利5%で10年間
	principal, _ := valueobjects.NewMoneyJPY(1000000)
	monthlyPayment, _ := valueobjects.NewMoneyJPY(50000)
	annualRate, _ := valueobjects.NewRate(5.0)
	years := 10

	result, err := service.CalculateCompoundInterestWithRegularPayments(principal, monthlyPayment, annualRate, years)
	if err != nil {
		t.Fatalf("定期積立複利計算に失敗しました: %v", err)
	}

	// 検証: 総拠出額は100万円 + 50万円×12ヶ月×10年 = 700万円
	expectedContribution := 1000000 + 50000*12*10
	if result.TotalContribution.Amount() != float64(expectedContribution) {
		t.Errorf("総拠出額が正しくありません。期待値: %d, 実際: %.0f", expectedContribution, result.TotalContribution.Amount())
	}

	// 最終金額は拠出額より大きいはず（利息があるため）
	if result.FinalAmount.Amount() <= result.TotalContribution.Amount() {
		t.Error("最終金額が総拠出額以下です。利息が正しく計算されていません")
	}

	// 利息収益は正の値であるはず
	if !result.InterestEarned.IsPositive() {
		t.Error("利息収益が正の値ではありません")
	}
}

func TestCalculateInflationAdjustedValue(t *testing.T) {
	service := NewFinancialCalculationService()

	// テストケース: 100万円を年2%のインフレで10年後の実質価値
	amount, _ := valueobjects.NewMoneyJPY(1000000)
	inflationRate, _ := valueobjects.NewRate(2.0)
	years := 10

	result, err := service.CalculateInflationAdjustedValue(amount, inflationRate, years)
	if err != nil {
		t.Fatalf("インフレ調整計算に失敗しました: %v", err)
	}

	// 検証: 10年後の実質価値は約82万円になるはず
	expectedRealValue := 1000000 / 1.2190 // 1/(1.02)^10 ≈ 0.8203
	if result.RealValue.Amount() < expectedRealValue*0.99 || result.RealValue.Amount() > expectedRealValue*1.01 {
		t.Errorf("実質価値が期待値と異なります。期待値: %.0f, 実際: %.0f", expectedRealValue, result.RealValue.Amount())
	}

	// 名目価値は元の金額と同じはず
	if result.NominalValue.Amount() != amount.Amount() {
		t.Error("名目価値が元の金額と異なります")
	}

	// 購買力の損失は正の値であるはず
	if !result.PurchasingPowerLoss.IsPositive() {
		t.Error("購買力の損失が正の値ではありません")
	}
}

func TestCalculateRetirementNeeds(t *testing.T) {
	service := NewFinancialCalculationService()

	// テストケース: 月30万円の生活費、年金15万円、20年間の老後
	monthlyExpenses, _ := valueobjects.NewMoneyJPY(300000)
	pensionAmount, _ := valueobjects.NewMoneyJPY(150000)
	yearsInRetirement := 20
	inflationRate, _ := valueobjects.NewRate(2.0)

	result, err := service.CalculateRetirementNeeds(monthlyExpenses, yearsInRetirement, inflationRate, pensionAmount)
	if err != nil {
		t.Fatalf("老後資金計算に失敗しました: %v", err)
	}

	// 検証: 月間不足額は15万円（30万円 - 15万円）
	expectedMonthlyNeeds := 150000.0
	if result.MonthlyNeeds.Amount() != expectedMonthlyNeeds {
		t.Errorf("月間必要額が正しくありません。期待値: %.0f, 実際: %.0f", expectedMonthlyNeeds, result.MonthlyNeeds.Amount())
	}

	// 総必要額は月間不足額 × 12ヶ月 × 20年 = 3600万円
	expectedTotalNeeds := expectedMonthlyNeeds * 12 * 20
	if result.TotalNeeds.Amount() != expectedTotalNeeds {
		t.Errorf("総必要額が正しくありません。期待値: %.0f, 実際: %.0f", expectedTotalNeeds, result.TotalNeeds.Amount())
	}

	// 退職後年数の確認
	if result.YearsInRetirement != yearsInRetirement {
		t.Errorf("退職後年数が正しくありません。期待値: %d, 実際: %d", yearsInRetirement, result.YearsInRetirement)
	}
}

func TestCalculateRetirementNeedsWithSufficientPension(t *testing.T) {
	service := NewFinancialCalculationService()

	// テストケース: 年金が生活費を上回る場合
	monthlyExpenses, _ := valueobjects.NewMoneyJPY(200000)
	pensionAmount, _ := valueobjects.NewMoneyJPY(250000) // 年金の方が多い
	yearsInRetirement := 20
	inflationRate, _ := valueobjects.NewRate(2.0)

	result, err := service.CalculateRetirementNeeds(monthlyExpenses, yearsInRetirement, inflationRate, pensionAmount)
	if err != nil {
		t.Fatalf("老後資金計算に失敗しました: %v", err)
	}

	// 検証: 年金で十分な場合は追加資金不要
	if !result.TotalNeeds.IsZero() {
		t.Error("年金が十分な場合は追加資金は不要のはずです")
	}

	if !result.MonthlyNeeds.IsZero() {
		t.Error("年金が十分な場合は月間追加必要額は0のはずです")
	}
}

func TestCalculateFutureValue(t *testing.T) {
	service := NewFinancialCalculationService()

	// テストケース: 100万円を年利3%で5年間
	presentValue, _ := valueobjects.NewMoneyJPY(1000000)
	rate, _ := valueobjects.NewRate(3.0)
	periods := 5

	futureValue, err := service.CalculateFutureValue(presentValue, rate, periods)
	if err != nil {
		t.Fatalf("将来価値計算に失敗しました: %v", err)
	}

	// 検証: 5年後は約115万円になるはず
	expectedValue := 1000000 * 1.1593 // (1.03)^5 ≈ 1.1593
	if futureValue.Amount() < expectedValue*0.99 || futureValue.Amount() > expectedValue*1.01 {
		t.Errorf("将来価値が期待値と異なります。期待値: %.0f, 実際: %.0f", expectedValue, futureValue.Amount())
	}
}

func TestCalculatePresentValue(t *testing.T) {
	service := NewFinancialCalculationService()

	// テストケース: 5年後の115万円の現在価値（年利3%）
	futureValue, _ := valueobjects.NewMoneyJPY(1159274)
	rate, _ := valueobjects.NewRate(3.0)
	periods := 5

	presentValue, err := service.CalculatePresentValue(futureValue, rate, periods)
	if err != nil {
		t.Fatalf("現在価値計算に失敗しました: %v", err)
	}

	// 検証: 現在価値は約100万円になるはず
	expectedValue := 1000000.0
	if presentValue.Amount() < expectedValue*0.99 || presentValue.Amount() > expectedValue*1.01 {
		t.Errorf("現在価値が期待値と異なります。期待値: %.0f, 実際: %.0f", expectedValue, presentValue.Amount())
	}
}

func TestCalculateEmergencyFundTarget(t *testing.T) {
	service := NewFinancialCalculationService()

	// テストケース: 月20万円の支出、6ヶ月分の緊急資金
	monthlyExpenses, _ := valueobjects.NewMoneyJPY(200000)
	targetMonths := 6
	inflationRate, _ := valueobjects.NewRate(2.0)
	yearsToTarget := 0 // 即座に必要

	target, err := service.CalculateEmergencyFundTarget(monthlyExpenses, targetMonths, inflationRate, yearsToTarget)
	if err != nil {
		t.Fatalf("緊急資金目標計算に失敗しました: %v", err)
	}

	// 検証: 6ヶ月分なので120万円
	expectedTarget := 200000.0 * 6
	if target.Amount() != expectedTarget {
		t.Errorf("緊急資金目標額が正しくありません。期待値: %.0f, 実際: %.0f", expectedTarget, target.Amount())
	}
}

func TestCalculateDebtPayoffTime(t *testing.T) {
	service := NewFinancialCalculationService()

	// テストケース: 100万円の債務、月5万円返済、年利3%
	debtAmount, _ := valueobjects.NewMoneyJPY(1000000)
	monthlyPayment, _ := valueobjects.NewMoneyJPY(50000)
	interestRate, _ := valueobjects.NewRate(3.0)

	months, err := service.CalculateDebtPayoffTime(debtAmount, monthlyPayment, interestRate)
	if err != nil {
		t.Fatalf("債務返済期間計算に失敗しました: %v", err)
	}

	// 検証: 返済期間は正の値であるはず
	if months <= 0 {
		t.Error("返済期間が正の値ではありません")
	}

	// 利息なしの場合は20ヶ月、利息ありの場合はそれより長いはず
	if months < 20 {
		t.Error("利息を考慮した返済期間が短すぎます")
	}
}

func TestCalculateDebtPayoffTimeWithInsufficientPayment(t *testing.T) {
	service := NewFinancialCalculationService()

	// テストケース: 返済額が利息を下回る場合
	debtAmount, _ := valueobjects.NewMoneyJPY(1000000)
	monthlyPayment, _ := valueobjects.NewMoneyJPY(1000) // 非常に少ない返済額
	interestRate, _ := valueobjects.NewRate(12.0)       // 高い利率

	months, err := service.CalculateDebtPayoffTime(debtAmount, monthlyPayment, interestRate)
	if err == nil {
		t.Error("返済額が不十分な場合はエラーになるはずです")
	}

	if months != -1 {
		t.Error("返済不可能な場合は-1を返すはずです")
	}
}
func TestFinancialCalculationServiceEdgeCases(t *testing.T) {
	service := NewFinancialCalculationService()

	// 非常に高い利率でのテスト
	highRate, _ := valueobjects.NewRate(50.0) // 50%
	principal, _ := valueobjects.NewMoneyJPY(1000000)

	result, err := service.CalculateCompoundInterest(principal, highRate, 5)
	if err != nil {
		t.Fatalf("高利率での複利計算に失敗しました: %v", err)
	}

	// 最終金額は元本より大幅に大きくなるはず
	if result.FinalAmount.Amount() <= principal.Amount()*2 {
		t.Error("高利率では最終金額が大幅に増加するはずです")
	}

	// 非常に小さな金額でのテスト
	smallAmount, _ := valueobjects.NewMoneyJPY(1.0) // 1円
	rate, _ := valueobjects.NewRate(5.0)
	smallResult, err := service.CalculateCompoundInterest(smallAmount, rate, 10)
	if err != nil {
		t.Fatalf("小額での複利計算に失敗しました: %v", err)
	}

	// 結果は正の値であるはず
	if !smallResult.FinalAmount.IsPositive() {
		t.Error("小額でも複利計算結果は正の値になるはずです")
	}

	// ゼロ期間でのテスト
	zeroResult, err := service.CalculateCompoundInterest(principal, rate, 0)
	if err != nil {
		t.Fatalf("ゼロ期間での複利計算に失敗しました: %v", err)
	}

	// ゼロ期間では元本と同じになるはず
	if zeroResult.FinalAmount.Amount() != principal.Amount() {
		t.Error("ゼロ期間では最終金額は元本と同じになるはずです")
	}
}
