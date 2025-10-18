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

func TestFinancialProfile_ValidationErrors(t *testing.T) {
	userID := UserID("test-user-123")
	monthlyIncome, _ := valueobjects.NewMoneyJPY(400000)
	expenses := ExpenseCollection{
		{Category: "住居費", Amount: mustCreateMoney(120000)},
	}
	savings := SavingsCollection{
		{Type: "deposit", Amount: mustCreateMoney(1000000)},
	}
	investmentReturn, _ := valueobjects.NewRate(5.0)
	inflationRate, _ := valueobjects.NewRate(2.0)

	// 空のユーザーID
	_, err := NewFinancialProfile("", monthlyIncome, expenses, savings, investmentReturn, inflationRate)
	if err == nil {
		t.Error("Expected error for empty user ID")
	}

	// 負の月収
	negativeIncome, _ := valueobjects.NewMoneyJPY(-1000)
	_, err = NewFinancialProfile(userID, negativeIncome, expenses, savings, investmentReturn, inflationRate)
	if err == nil {
		t.Error("Expected error for negative monthly income")
	}

	// ゼロの月収
	zeroIncome, _ := valueobjects.NewMoneyJPY(0)
	_, err = NewFinancialProfile(userID, zeroIncome, expenses, savings, investmentReturn, inflationRate)
	if err == nil {
		t.Error("Expected error for zero monthly income")
	}
}

func TestFinancialProfile_UpdateMethods(t *testing.T) {
	profile := createTestFinancialProfile(t)

	// 月収の更新
	newIncome, _ := valueobjects.NewMoneyJPY(500000)
	err := profile.UpdateMonthlyIncome(newIncome)
	if err != nil {
		t.Errorf("Failed to update monthly income: %v", err)
	}
	if profile.MonthlyIncome().Amount() != 500000 {
		t.Error("Monthly income was not updated correctly")
	}

	// 負の月収での更新（エラーになるはず）
	negativeIncome, _ := valueobjects.NewMoneyJPY(-1000)
	err = profile.UpdateMonthlyIncome(negativeIncome)
	if err == nil {
		t.Error("Expected error when updating with negative income")
	}

	// 支出の更新
	newExpenses := ExpenseCollection{
		{Category: "住居費", Amount: mustCreateMoney(150000)},
		{Category: "食費", Amount: mustCreateMoney(80000)},
	}
	err = profile.UpdateMonthlyExpenses(newExpenses)
	if err != nil {
		t.Errorf("Failed to update monthly expenses: %v", err)
	}

	// 貯蓄の更新
	newSavings := SavingsCollection{
		{Type: "deposit", Amount: mustCreateMoney(2000000)},
	}
	err = profile.UpdateCurrentSavings(newSavings)
	if err != nil {
		t.Errorf("Failed to update current savings: %v", err)
	}

	// 投資利回りの更新
	newRate, _ := valueobjects.NewRate(7.0)
	err = profile.UpdateInvestmentReturn(newRate)
	if err != nil {
		t.Errorf("Failed to update investment return: %v", err)
	}
	if profile.InvestmentReturn().AsPercentage() != 7.0 {
		t.Error("Investment return was not updated correctly")
	}

	// インフレ率の更新
	newInflationRate, _ := valueobjects.NewRate(3.0)
	err = profile.UpdateInflationRate(newInflationRate)
	if err != nil {
		t.Errorf("Failed to update inflation rate: %v", err)
	}
	if profile.InflationRate().AsPercentage() != 3.0 {
		t.Error("Inflation rate was not updated correctly")
	}
}

func TestFinancialProfile_ValidateFinancialHealth(t *testing.T) {
	// 健全な財務プロファイル
	healthyProfile := createTestFinancialProfile(t)
	err := healthyProfile.ValidateFinancialHealth()
	if err != nil {
		t.Errorf("Healthy profile should not have validation errors: %v", err)
	}

	// 支出が収入を上回る場合
	userID := UserID("test-user-123")
	monthlyIncome, _ := valueobjects.NewMoneyJPY(200000)
	expenses := ExpenseCollection{
		{Category: "住居費", Amount: mustCreateMoney(150000)},
		{Category: "食費", Amount: mustCreateMoney(100000)}, // 合計250000 > 200000
	}
	savings := SavingsCollection{
		{Type: "deposit", Amount: mustCreateMoney(1000000)},
	}
	investmentReturn, _ := valueobjects.NewRate(5.0)
	inflationRate, _ := valueobjects.NewRate(2.0)

	unhealthyProfile, _ := NewFinancialProfile(userID, monthlyIncome, expenses, savings, investmentReturn, inflationRate)
	err = unhealthyProfile.ValidateFinancialHealth()
	if err == nil {
		t.Error("Expected validation error for expenses exceeding income")
	}

	// 貯蓄率が低い場合（収入の10%未満）
	lowSavingsIncome, _ := valueobjects.NewMoneyJPY(300000)
	lowSavingsExpenses := ExpenseCollection{
		{Category: "住居費", Amount: mustCreateMoney(280000)}, // 貯蓄率 = 20000/300000 = 6.7%
	}
	lowSavingsProfile, _ := NewFinancialProfile(userID, lowSavingsIncome, lowSavingsExpenses, savings, investmentReturn, inflationRate)
	err = lowSavingsProfile.ValidateFinancialHealth()
	if err == nil {
		t.Error("Expected validation warning for low savings rate")
	}
}

func TestFinancialProfile_ProjectAssets(t *testing.T) {
	profile := createTestFinancialProfile(t)

	projections, err := profile.ProjectAssets(5)
	if err != nil {
		t.Fatalf("Failed to project assets: %v", err)
	}

	if len(projections) != 5 {
		t.Errorf("Expected 5 projections, got %d", len(projections))
	}

	// 各年の予測が正しく設定されているかチェック
	for i, projection := range projections {
		expectedYear := i + 1
		if projection.Year != expectedYear {
			t.Errorf("Expected year %d, got %d", expectedYear, projection.Year)
		}

		// 総資産は正の値であるはず
		if !projection.TotalAssets.IsPositive() {
			t.Errorf("Total assets should be positive for year %d", expectedYear)
		}

		// 実質価値は名目価値以下であるはず（インフレ調整のため）
		isLessOrEqual, err := projection.RealValue.LessThan(projection.TotalAssets)
		if err != nil {
			t.Errorf("Failed to compare real value and total assets: %v", err)
		}
		isEqual, err := projection.RealValue.Equal(projection.TotalAssets)
		if err != nil {
			t.Errorf("Failed to compare real value and total assets: %v", err)
		}
		if !isLessOrEqual && !isEqual {
			t.Errorf("Real value should be less than or equal to total assets for year %d", expectedYear)
		}
	}

	// 無効な年数での予測
	_, err = profile.ProjectAssets(0)
	if err == nil {
		t.Error("Expected error for zero years projection")
	}

	_, err = profile.ProjectAssets(-1)
	if err == nil {
		t.Error("Expected error for negative years projection")
	}
}

func TestExpenseCollection_Methods(t *testing.T) {
	expenses := ExpenseCollection{
		{Category: "住居費", Amount: mustCreateMoney(120000)},
		{Category: "食費", Amount: mustCreateMoney(60000)},
		{Category: "住居費", Amount: mustCreateMoney(30000)}, // 同じカテゴリ
	}

	// 合計の計算
	total, err := expenses.Total()
	if err != nil {
		t.Errorf("Failed to calculate total expenses: %v", err)
	}
	if total.Amount() != 210000 {
		t.Errorf("Expected total 210000, got %f", total.Amount())
	}

	// カテゴリ別の取得
	housingExpenses := expenses.GetByCategory("住居費")
	if len(housingExpenses) != 2 {
		t.Errorf("Expected 2 housing expenses, got %d", len(housingExpenses))
	}

	foodExpenses := expenses.GetByCategory("食費")
	if len(foodExpenses) != 1 {
		t.Errorf("Expected 1 food expense, got %d", len(foodExpenses))
	}

	// 存在しないカテゴリ
	nonExistent := expenses.GetByCategory("交通費")
	if len(nonExistent) != 0 {
		t.Errorf("Expected 0 non-existent expenses, got %d", len(nonExistent))
	}

	// 空のコレクション
	emptyExpenses := ExpenseCollection{}
	emptyTotal, err := emptyExpenses.Total()
	if err != nil {
		t.Errorf("Failed to calculate empty expenses total: %v", err)
	}
	if !emptyTotal.IsZero() {
		t.Error("Empty expenses total should be zero")
	}
}

func TestSavingsCollection_Methods(t *testing.T) {
	savings := SavingsCollection{
		{Type: "deposit", Amount: mustCreateMoney(1000000)},
		{Type: "investment", Amount: mustCreateMoney(500000)},
		{Type: "deposit", Amount: mustCreateMoney(200000)}, // 同じタイプ
	}

	// 合計の計算
	total, err := savings.Total()
	if err != nil {
		t.Errorf("Failed to calculate total savings: %v", err)
	}
	if total.Amount() != 1700000 {
		t.Errorf("Expected total 1700000, got %f", total.Amount())
	}

	// タイプ別の取得
	deposits := savings.GetByType("deposit")
	if len(deposits) != 2 {
		t.Errorf("Expected 2 deposits, got %d", len(deposits))
	}

	investments := savings.GetByType("investment")
	if len(investments) != 1 {
		t.Errorf("Expected 1 investment, got %d", len(investments))
	}

	// 存在しないタイプ
	nonExistent := savings.GetByType("crypto")
	if len(nonExistent) != 0 {
		t.Errorf("Expected 0 non-existent savings, got %d", len(nonExistent))
	}

	// 空のコレクション
	emptySavings := SavingsCollection{}
	emptyTotal, err := emptySavings.Total()
	if err != nil {
		t.Errorf("Failed to calculate empty savings total: %v", err)
	}
	if !emptyTotal.IsZero() {
		t.Error("Empty savings total should be zero")
	}
}

// ヘルパー関数：テスト用のFinancialProfile作成
func createTestFinancialProfile(t *testing.T) *FinancialProfile {
	userID := UserID("test-user-123")
	monthlyIncome, _ := valueobjects.NewMoneyJPY(400000)
	expenses := ExpenseCollection{
		{Category: "住居費", Amount: mustCreateMoney(120000)},
		{Category: "食費", Amount: mustCreateMoney(60000)},
	}
	savings := SavingsCollection{
		{Type: "deposit", Amount: mustCreateMoney(1000000)},
	}
	investmentReturn, _ := valueobjects.NewRate(5.0)
	inflationRate, _ := valueobjects.NewRate(2.0)

	profile, err := NewFinancialProfile(userID, monthlyIncome, expenses, savings, investmentReturn, inflationRate)
	if err != nil {
		t.Fatalf("Failed to create test financial profile: %v", err)
	}
	return profile
}

// ヘルパー関数：絶対値計算
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
func TestGoal_ValidationErrors(t *testing.T) {
	userID := UserID("test-user-123")
	targetAmount := mustCreateMoney(1000000)
	targetDate := time.Now().AddDate(2, 0, 0)
	monthlyContribution := mustCreateMoney(50000)

	// 空のユーザーID
	_, err := NewGoal("", GoalTypeSavings, "テスト目標", targetAmount, targetDate, monthlyContribution)
	if err == nil {
		t.Error("Expected error for empty user ID")
	}

	// 無効な目標タイプ
	_, err = NewGoal(userID, GoalType("invalid"), "テスト目標", targetAmount, targetDate, monthlyContribution)
	if err == nil {
		t.Error("Expected error for invalid goal type")
	}

	// 空のタイトル
	_, err = NewGoal(userID, GoalTypeSavings, "", targetAmount, targetDate, monthlyContribution)
	if err == nil {
		t.Error("Expected error for empty title")
	}

	// 負の目標金額
	negativeAmount := mustCreateMoney(-1000)
	_, err = NewGoal(userID, GoalTypeSavings, "テスト目標", negativeAmount, targetDate, monthlyContribution)
	if err == nil {
		t.Error("Expected error for negative target amount")
	}

	// 過去の目標日
	pastDate := time.Now().AddDate(-1, 0, 0)
	_, err = NewGoal(userID, GoalTypeSavings, "テスト目標", targetAmount, pastDate, monthlyContribution)
	if err == nil {
		t.Error("Expected error for past target date")
	}

	// 負の月間拠出額
	negativeContribution := mustCreateMoney(-1000)
	_, err = NewGoal(userID, GoalTypeSavings, "テスト目標", targetAmount, targetDate, negativeContribution)
	if err == nil {
		t.Error("Expected error for negative monthly contribution")
	}
}

func TestGoal_UpdateMethods(t *testing.T) {
	goal := createTestGoal(t)

	// 現在金額の更新
	newAmount := mustCreateMoney(500000)
	err := goal.UpdateCurrentAmount(newAmount)
	if err != nil {
		t.Errorf("Failed to update current amount: %v", err)
	}
	if goal.CurrentAmount().Amount() != 500000 {
		t.Error("Current amount was not updated correctly")
	}

	// 負の現在金額での更新（エラーになるはず）
	negativeAmount := mustCreateMoney(-1000)
	err = goal.UpdateCurrentAmount(negativeAmount)
	if err == nil {
		t.Error("Expected error when updating with negative current amount")
	}

	// 月間拠出額の更新
	newContribution := mustCreateMoney(75000)
	err = goal.UpdateMonthlyContribution(newContribution)
	if err != nil {
		t.Errorf("Failed to update monthly contribution: %v", err)
	}
	if goal.MonthlyContribution().Amount() != 75000 {
		t.Error("Monthly contribution was not updated correctly")
	}

	// 目標金額の更新
	newTargetAmount := mustCreateMoney(1500000)
	err = goal.UpdateTargetAmount(newTargetAmount)
	if err != nil {
		t.Errorf("Failed to update target amount: %v", err)
	}
	if goal.TargetAmount().Amount() != 1500000 {
		t.Error("Target amount was not updated correctly")
	}

	// 目標日の更新
	newTargetDate := time.Now().AddDate(3, 0, 0)
	err = goal.UpdateTargetDate(newTargetDate)
	if err != nil {
		t.Errorf("Failed to update target date: %v", err)
	}

	// タイトルの更新
	newTitle := "更新されたテスト目標"
	err = goal.UpdateTitle(newTitle)
	if err != nil {
		t.Errorf("Failed to update title: %v", err)
	}
	if goal.Title() != newTitle {
		t.Error("Title was not updated correctly")
	}

	// 空のタイトルでの更新（エラーになるはず）
	err = goal.UpdateTitle("")
	if err == nil {
		t.Error("Expected error when updating with empty title")
	}
}

func TestGoal_StatusMethods(t *testing.T) {
	goal := createTestGoal(t)

	// 初期状態はアクティブ
	if !goal.IsActive() {
		t.Error("Goal should be active initially")
	}

	// 非アクティブ化
	goal.Deactivate()
	if goal.IsActive() {
		t.Error("Goal should be inactive after deactivation")
	}

	// アクティブ化
	goal.Activate()
	if !goal.IsActive() {
		t.Error("Goal should be active after activation")
	}

	// 完了状態のテスト
	if goal.IsCompleted() {
		t.Error("Goal should not be completed initially")
	}

	// 目標金額に到達
	err := goal.UpdateCurrentAmount(goal.TargetAmount())
	if err != nil {
		t.Errorf("Failed to update current amount: %v", err)
	}
	if !goal.IsCompleted() {
		t.Error("Goal should be completed when current amount equals target amount")
	}

	// 期限切れのテスト（過去の目標日を設定）
	pastDate := time.Now().AddDate(-1, 0, 0)
	goal.targetDate = pastDate              // 直接設定（テスト用）
	goal.currentAmount = mustCreateMoney(0) // 未完了状態に戻す
	if !goal.IsOverdue() {
		t.Error("Goal should be overdue when target date is in the past and not completed")
	}
}

func TestGoal_CalculationMethods(t *testing.T) {
	goal := createTestGoal(t)

	// 現在金額を設定
	currentAmount := mustCreateMoney(600000)
	err := goal.UpdateCurrentAmount(currentAmount)
	if err != nil {
		t.Errorf("Failed to update current amount: %v", err)
	}

	// 残り必要金額の計算
	remainingAmount, err := goal.GetRemainingAmount()
	if err != nil {
		t.Errorf("Failed to get remaining amount: %v", err)
	}
	expectedRemaining := goal.TargetAmount().Amount() - currentAmount.Amount()
	if remainingAmount.Amount() != expectedRemaining {
		t.Errorf("Expected remaining amount %f, got %f", expectedRemaining, remainingAmount.Amount())
	}

	// 残り日数の計算
	remainingDays := goal.GetRemainingDays()
	if remainingDays <= 0 {
		t.Error("Remaining days should be positive for future target date")
	}

	// 必要月間貯蓄額の計算
	requiredMonthlySavings, err := goal.CalculateRequiredMonthlySavings()
	if err != nil {
		t.Errorf("Failed to calculate required monthly savings: %v", err)
	}
	if !requiredMonthlySavings.IsPositive() {
		t.Error("Required monthly savings should be positive")
	}

	// 完了予定日の推定
	monthlySavings := mustCreateMoney(100000)
	completionDate, err := goal.EstimateCompletionDate(monthlySavings)
	if err != nil {
		t.Errorf("Failed to estimate completion date: %v", err)
	}
	if completionDate.Before(time.Now()) {
		t.Error("Completion date should be in the future")
	}

	// ゼロの月間貯蓄での推定（エラーになるはず）
	zeroSavings := mustCreateMoney(0)
	_, err = goal.EstimateCompletionDate(zeroSavings)
	if err == nil {
		t.Error("Expected error for zero monthly savings")
	}
}

func TestGoal_IsAchievable(t *testing.T) {
	goal := createTestGoal(t)
	profile := createTestFinancialProfile(t)

	// 達成可能性の判定
	achievable, err := goal.IsAchievable(profile)
	if err != nil {
		t.Errorf("Failed to check achievability: %v", err)
	}
	// 具体的な値は財務プロファイルと目標の設定による

	// nilプロファイルでの判定（エラーになるはず）
	_, err = goal.IsAchievable(nil)
	if err == nil {
		t.Error("Expected error for nil financial profile")
	}

	// 支出が収入を上回るプロファイルでの判定
	userID := UserID("test-user-123")
	monthlyIncome, _ := valueobjects.NewMoneyJPY(200000)
	expenses := ExpenseCollection{
		{Category: "住居費", Amount: mustCreateMoney(250000)}, // 収入を上回る
	}
	savings := SavingsCollection{
		{Type: "deposit", Amount: mustCreateMoney(1000000)},
	}
	investmentReturn, _ := valueobjects.NewRate(5.0)
	inflationRate, _ := valueobjects.NewRate(2.0)

	unhealthyProfile, _ := NewFinancialProfile(userID, monthlyIncome, expenses, savings, investmentReturn, inflationRate)
	achievable, err = goal.IsAchievable(unhealthyProfile)
	if err != nil {
		t.Errorf("Failed to check achievability with unhealthy profile: %v", err)
	}
	if achievable {
		t.Error("Goal should not be achievable with negative net savings")
	}
}

func TestGoalType_Methods(t *testing.T) {
	// 有効なGoalTypeのテスト
	validTypes := []GoalType{GoalTypeSavings, GoalTypeRetirement, GoalTypeEmergency, GoalTypeCustom}
	for _, goalType := range validTypes {
		if !goalType.IsValid() {
			t.Errorf("GoalType %s should be valid", goalType)
		}
		if goalType.String() == "" {
			t.Errorf("GoalType %s should have a string representation", goalType)
		}
	}

	// 無効なGoalTypeのテスト
	invalidType := GoalType("invalid")
	if invalidType.IsValid() {
		t.Error("Invalid GoalType should not be valid")
	}
	if invalidType.String() == "" {
		t.Error("Invalid GoalType should still have a string representation")
	}
}

func TestProgressRate_Methods(t *testing.T) {
	// 正常な進捗率
	progress, err := NewProgressRate(75.5)
	if err != nil {
		t.Errorf("Failed to create progress rate: %v", err)
	}
	if progress.AsPercentage() != 75.5 {
		t.Errorf("Expected 75.5%%, got %f%%", progress.AsPercentage())
	}
	if progress.IsComplete() {
		t.Error("75.5% progress should not be complete")
	}

	// 100%の進捗率
	completeProgress, err := NewProgressRate(100.0)
	if err != nil {
		t.Errorf("Failed to create complete progress rate: %v", err)
	}
	if !completeProgress.IsComplete() {
		t.Error("100% progress should be complete")
	}

	// 100%を超える進捗率（100%に制限されるはず）
	overProgress, err := NewProgressRate(150.0)
	if err != nil {
		t.Errorf("Failed to create over progress rate: %v", err)
	}
	if overProgress.AsPercentage() != 100.0 {
		t.Errorf("Over 100%% progress should be capped at 100%%, got %f%%", overProgress.AsPercentage())
	}

	// 負の進捗率（0%に制限されるはず）
	negativeProgress, err := NewProgressRate(-10.0)
	if err != nil {
		t.Errorf("Failed to create negative progress rate: %v", err)
	}
	if negativeProgress.AsPercentage() != 0.0 {
		t.Errorf("Negative progress should be capped at 0%%, got %f%%", negativeProgress.AsPercentage())
	}

	// 文字列表現のテスト
	if progress.String() != "75.5%" {
		t.Errorf("Expected '75.5%%', got '%s'", progress.String())
	}
}

// ヘルパー関数：テスト用のGoal作成
func createTestGoal(t *testing.T) *Goal {
	userID := UserID("test-user-123")
	targetAmount := mustCreateMoney(2000000)
	monthlyContribution := mustCreateMoney(50000)
	targetDate := time.Now().AddDate(3, 0, 0) // 3年後

	goal, err := NewGoal(userID, GoalTypeSavings, "テスト目標", targetAmount, targetDate, monthlyContribution)
	if err != nil {
		t.Fatalf("Failed to create test goal: %v", err)
	}
	return goal
}
func TestRetirementData_ValidationErrors(t *testing.T) {
	userID := UserID("test-user-123")
	monthlyRetirementExpenses := mustCreateMoney(250000)
	pensionAmount := mustCreateMoney(150000)

	// 空のユーザーID
	_, err := NewRetirementData("", 35, 65, 85, monthlyRetirementExpenses, pensionAmount)
	if err == nil {
		t.Error("Expected error for empty user ID")
	}

	// 無効な年齢（現在年齢が負）
	_, err = NewRetirementData(userID, -1, 65, 85, monthlyRetirementExpenses, pensionAmount)
	if err == nil {
		t.Error("Expected error for negative current age")
	}

	// 無効な年齢（退職年齢が現在年齢以下）
	_, err = NewRetirementData(userID, 65, 60, 85, monthlyRetirementExpenses, pensionAmount)
	if err == nil {
		t.Error("Expected error for retirement age less than current age")
	}

	// 無効な年齢（平均寿命が退職年齢以下）
	_, err = NewRetirementData(userID, 35, 65, 60, monthlyRetirementExpenses, pensionAmount)
	if err == nil {
		t.Error("Expected error for life expectancy less than retirement age")
	}

	// 負の月間退職後支出
	negativeExpenses := mustCreateMoney(-1000)
	_, err = NewRetirementData(userID, 35, 65, 85, negativeExpenses, pensionAmount)
	if err == nil {
		t.Error("Expected error for negative monthly retirement expenses")
	}

	// 負の年金額
	negativePension := mustCreateMoney(-1000)
	_, err = NewRetirementData(userID, 35, 65, 85, monthlyRetirementExpenses, negativePension)
	if err == nil {
		t.Error("Expected error for negative pension amount")
	}
}

func TestRetirementData_CalculationMethods(t *testing.T) {
	userID := UserID("test-user-123")
	currentAge := 35
	retirementAge := 65
	lifeExpectancy := 85
	monthlyRetirementExpenses := mustCreateMoney(250000)
	pensionAmount := mustCreateMoney(150000)

	retirementData, err := NewRetirementData(userID, currentAge, retirementAge, lifeExpectancy, monthlyRetirementExpenses, pensionAmount)
	if err != nil {
		t.Fatalf("Failed to create retirement data: %v", err)
	}

	// 退職までの年数計算
	yearsUntilRetirement := retirementData.CalculateYearsUntilRetirement()
	expectedYears := retirementAge - currentAge
	if yearsUntilRetirement != expectedYears {
		t.Errorf("Expected %d years until retirement, got %d", expectedYears, yearsUntilRetirement)
	}

	// 退職後年数計算
	retirementYears := retirementData.CalculateRetirementYears()
	expectedRetirementYears := lifeExpectancy - retirementAge
	if retirementYears != expectedRetirementYears {
		t.Errorf("Expected %d retirement years, got %d", expectedRetirementYears, retirementYears)
	}

	// 年金不足額計算
	shortfall, err := retirementData.GetPensionShortfall()
	if err != nil {
		t.Errorf("Failed to calculate pension shortfall: %v", err)
	}
	expectedShortfall := monthlyRetirementExpenses.Amount() - pensionAmount.Amount()
	if shortfall.Amount() != expectedShortfall {
		t.Errorf("Expected shortfall %f, got %f", expectedShortfall, shortfall.Amount())
	}

	// 年金が十分な場合のテスト
	sufficientPension := mustCreateMoney(300000) // 支出を上回る年金
	retirementDataSufficient, _ := NewRetirementData(userID, currentAge, retirementAge, lifeExpectancy, monthlyRetirementExpenses, sufficientPension)
	shortfallSufficient, err := retirementDataSufficient.GetPensionShortfall()
	if err != nil {
		t.Errorf("Failed to calculate pension shortfall for sufficient pension: %v", err)
	}
	if !shortfallSufficient.IsZero() {
		t.Error("Shortfall should be zero when pension is sufficient")
	}
}

func TestRetirementData_UpdateMethods(t *testing.T) {
	retirementData := createTestRetirementData(t)

	// 退職年齢の更新
	newRetirementAge := 67
	err := retirementData.UpdateRetirementAge(newRetirementAge)
	if err != nil {
		t.Errorf("Failed to update retirement age: %v", err)
	}
	if retirementData.RetirementAge() != newRetirementAge {
		t.Error("Retirement age was not updated correctly")
	}

	// 無効な退職年齢での更新（現在年齢以下）
	err = retirementData.UpdateRetirementAge(30)
	if err == nil {
		t.Error("Expected error when updating with retirement age less than current age")
	}

	// 平均寿命の更新
	newLifeExpectancy := 90
	err = retirementData.UpdateLifeExpectancy(newLifeExpectancy)
	if err != nil {
		t.Errorf("Failed to update life expectancy: %v", err)
	}
	if retirementData.LifeExpectancy() != newLifeExpectancy {
		t.Error("Life expectancy was not updated correctly")
	}

	// 無効な平均寿命での更新（退職年齢以下）
	err = retirementData.UpdateLifeExpectancy(60)
	if err == nil {
		t.Error("Expected error when updating with life expectancy less than retirement age")
	}

	// 月間退職後支出の更新
	newExpenses := mustCreateMoney(300000)
	err = retirementData.UpdateMonthlyRetirementExpenses(newExpenses)
	if err != nil {
		t.Errorf("Failed to update monthly retirement expenses: %v", err)
	}
	if retirementData.MonthlyRetirementExpenses().Amount() != 300000 {
		t.Error("Monthly retirement expenses were not updated correctly")
	}

	// 年金額の更新
	newPensionAmount := mustCreateMoney(200000)
	err = retirementData.UpdatePensionAmount(newPensionAmount)
	if err != nil {
		t.Errorf("Failed to update pension amount: %v", err)
	}
	if retirementData.PensionAmount().Amount() != 200000 {
		t.Error("Pension amount was not updated correctly")
	}
}

func TestRetirementData_EdgeCases(t *testing.T) {
	userID := UserID("test-user-123")

	// 現在年齢と退職年齢が同じ場合
	currentAge := 65
	retirementAge := 65
	lifeExpectancy := 85
	monthlyRetirementExpenses := mustCreateMoney(250000)
	pensionAmount := mustCreateMoney(150000)

	retirementData, err := NewRetirementData(userID, currentAge, retirementAge, lifeExpectancy, monthlyRetirementExpenses, pensionAmount)
	if err != nil {
		t.Errorf("Should allow current age equal to retirement age: %v", err)
	}

	yearsUntilRetirement := retirementData.CalculateYearsUntilRetirement()
	if yearsUntilRetirement != 0 {
		t.Errorf("Expected 0 years until retirement, got %d", yearsUntilRetirement)
	}

	// 退職年齢と平均寿命が同じ場合
	retirementAge = 85
	lifeExpectancy = 85
	retirementDataSame, err := NewRetirementData(userID, currentAge, retirementAge, lifeExpectancy, monthlyRetirementExpenses, pensionAmount)
	if err != nil {
		t.Errorf("Should allow retirement age equal to life expectancy: %v", err)
	}

	retirementYears := retirementDataSame.CalculateRetirementYears()
	if retirementYears != 0 {
		t.Errorf("Expected 0 retirement years, got %d", retirementYears)
	}
}

// ヘルパー関数：テスト用のRetirementData作成
func createTestRetirementData(t *testing.T) *RetirementData {
	userID := UserID("test-user-123")
	currentAge := 35
	retirementAge := 65
	lifeExpectancy := 85
	monthlyRetirementExpenses := mustCreateMoney(250000)
	pensionAmount := mustCreateMoney(150000)

	retirementData, err := NewRetirementData(userID, currentAge, retirementAge, lifeExpectancy, monthlyRetirementExpenses, pensionAmount)
	if err != nil {
		t.Fatalf("Failed to create test retirement data: %v", err)
	}
	return retirementData
}
