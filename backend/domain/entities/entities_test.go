package entities

import (
	"testing"
	"time"

	"financial-planning-calculator/domain/valueobjects"
)

func TestFinancialProfile_Creation(t *testing.T) {
	// テスト用のデータを準備
	userID := UserID("test-user-123")
	monthlyIncome, _ := valueobjects.NewMoneyJPY(400000)

	expenses := ExpenseCollection{
		{Category: "住居費", Amount: mustCreateMoney(120000), Description: "家賃"},
		{Category: "食費", Amount: mustCreateMoney(60000), Description: "食事代"},
	}

	savings := SavingsCollection{
		{Type: "deposit", Amount: mustCreateMoney(1000000), Description: "普通預金"},
		{Type: "investment", Amount: mustCreateMoney(500000), Description: "投資信託"},
	}

	investmentReturn, _ := valueobjects.NewRate(5.0)
	inflationRate, _ := valueobjects.NewRate(2.0)

	// FinancialProfileを作成
	profile, err := NewFinancialProfile(userID, monthlyIncome, expenses, savings, investmentReturn, inflationRate)
	if err != nil {
		t.Fatalf("FinancialProfile作成に失敗しました: %v", err)
	}

	// 基本的な検証
	if profile.UserID() != userID {
		t.Errorf("UserIDが期待値と異なります。期待値: %s, 実際: %s", userID, profile.UserID())
	}

	if profile.MonthlyIncome().Amount() != 400000 {
		t.Errorf("月収が期待値と異なります。期待値: 400000, 実際: %f", profile.MonthlyIncome().Amount())
	}

	// 純貯蓄額の計算をテスト
	netSavings, err := profile.CalculateNetSavings()
	if err != nil {
		t.Fatalf("純貯蓄額の計算に失敗しました: %v", err)
	}

	expectedNetSavings := 400000.0 - 120000.0 - 60000.0 // 220000
	if netSavings.Amount() != expectedNetSavings {
		t.Errorf("純貯蓄額が期待値と異なります。期待値: %f, 実際: %f", expectedNetSavings, netSavings.Amount())
	}
}

func TestGoal_Creation(t *testing.T) {
	userID := UserID("test-user-123")
	goalType := GoalTypeSavings
	title := "新車購入資金"
	targetAmount := mustCreateMoney(3000000)
	targetDate := time.Now().AddDate(2, 0, 0) // 2年後
	monthlyContribution := mustCreateMoney(100000)

	goal, err := NewGoal(userID, goalType, title, targetAmount, targetDate, monthlyContribution)
	if err != nil {
		t.Fatalf("Goal作成に失敗しました: %v", err)
	}

	// 基本的な検証
	if goal.UserID() != userID {
		t.Errorf("UserIDが期待値と異なります。期待値: %s, 実際: %s", userID, goal.UserID())
	}

	if goal.Title() != title {
		t.Errorf("タイトルが期待値と異なります。期待値: %s, 実際: %s", title, goal.Title())
	}

	if goal.TargetAmount().Amount() != 3000000 {
		t.Errorf("目標金額が期待値と異なります。期待値: 3000000, 実際: %f", goal.TargetAmount().Amount())
	}

	// 進捗率の計算をテスト
	currentAmount := mustCreateMoney(1000000)
	progress, err := goal.CalculateProgress(currentAmount)
	if err != nil {
		t.Fatalf("進捗率の計算に失敗しました: %v", err)
	}

	expectedProgress := (1000000.0 / 3000000.0) * 100 // 33.33%
	tolerance := 0.01                                 // 0.01%の許容誤差
	if abs(progress.AsPercentage()-expectedProgress) > tolerance {
		t.Errorf("進捗率が期待値と異なります。期待値: %f, 実際: %f", expectedProgress, progress.AsPercentage())
	}
}

func TestRetirementData_Creation(t *testing.T) {
	userID := UserID("test-user-123")
	currentAge := 35
	retirementAge := 65
	lifeExpectancy := 85
	monthlyRetirementExpenses := mustCreateMoney(250000)
	pensionAmount := mustCreateMoney(150000)

	retirementData, err := NewRetirementData(
		userID, currentAge, retirementAge, lifeExpectancy,
		monthlyRetirementExpenses, pensionAmount)
	if err != nil {
		t.Fatalf("RetirementData作成に失敗しました: %v", err)
	}

	// 基本的な検証
	if retirementData.UserID() != userID {
		t.Errorf("UserIDが期待値と異なります。期待値: %s, 実際: %s", userID, retirementData.UserID())
	}

	if retirementData.CurrentAge() != currentAge {
		t.Errorf("現在年齢が期待値と異なります。期待値: %d, 実際: %d", currentAge, retirementData.CurrentAge())
	}

	// 退職までの年数計算をテスト
	yearsUntilRetirement := retirementData.CalculateYearsUntilRetirement()
	expectedYears := retirementAge - currentAge // 30年
	if yearsUntilRetirement != expectedYears {
		t.Errorf("退職までの年数が期待値と異なります。期待値: %d, 実際: %d", expectedYears, yearsUntilRetirement)
	}

	// 退職後年数計算をテスト
	retirementYears := retirementData.CalculateRetirementYears()
	expectedRetirementYears := lifeExpectancy - retirementAge // 20年
	if retirementYears != expectedRetirementYears {
		t.Errorf("退職後年数が期待値と異なります。期待値: %d, 実際: %d", expectedRetirementYears, retirementYears)
	}

	// 年金不足額計算をテスト
	shortfall, err := retirementData.GetPensionShortfall()
	if err != nil {
		t.Fatalf("年金不足額の計算に失敗しました: %v", err)
	}

	expectedShortfall := 250000.0 - 150000.0 // 100000
	if shortfall.Amount() != expectedShortfall {
		t.Errorf("年金不足額が期待値と異なります。期待値: %f, 実際: %f", expectedShortfall, shortfall.Amount())
	}
}

// ヘルパー関数：テスト用のMoney作成
func mustCreateMoney(amount float64) valueobjects.Money {
	money, err := valueobjects.NewMoneyJPY(amount)
	if err != nil {
		panic(err)
	}
	return money
}

// ヘルパー関数：絶対値計算
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
