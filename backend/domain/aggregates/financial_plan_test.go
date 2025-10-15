package aggregates

import (
	"testing"
	"time"

	"financial-planning-calculator/domain/entities"
	"financial-planning-calculator/domain/valueobjects"
)

func TestNewFinancialPlan(t *testing.T) {
	// テスト用の財務プロファイルを作成
	monthlyIncome, _ := valueobjects.NewMoneyJPY(400000)
	expenses := entities.ExpenseCollection{
		{Category: "住居費", Amount: mustCreateMoney(120000)},
		{Category: "食費", Amount: mustCreateMoney(60000)},
	}
	savings := entities.SavingsCollection{
		{Type: "deposit", Amount: mustCreateMoney(1000000)},
	}
	investmentReturn, _ := valueobjects.NewRate(5.0)
	inflationRate, _ := valueobjects.NewRate(2.0)

	profile, err := entities.NewFinancialProfile(
		"user123",
		monthlyIncome,
		expenses,
		savings,
		investmentReturn,
		inflationRate,
	)
	if err != nil {
		t.Fatalf("財務プロファイルの作成に失敗しました: %v", err)
	}

	// 財務計画を作成
	plan, err := NewFinancialPlan(profile)
	if err != nil {
		t.Fatalf("財務計画の作成に失敗しました: %v", err)
	}

	// 検証
	if plan.ID() == "" {
		t.Error("財務計画IDが設定されていません")
	}

	if plan.Profile() != profile {
		t.Error("財務プロファイルが正しく設定されていません")
	}

	if len(plan.Goals()) != 0 {
		t.Error("初期状態では目標は空である必要があります")
	}

	if plan.EmergencyFund() == nil {
		t.Error("緊急資金設定が初期化されていません")
	}

	if plan.EmergencyFund().TargetMonths != 3 {
		t.Error("デフォルトの緊急資金目標月数が正しくありません")
	}
}

func TestAddGoal(t *testing.T) {
	plan := createTestFinancialPlan(t)

	// テスト用の目標を作成
	targetAmount, _ := valueobjects.NewMoneyJPY(1000000)
	monthlyContribution, _ := valueobjects.NewMoneyJPY(50000)
	targetDate := time.Now().AddDate(2, 0, 0) // 2年後

	goal, err := entities.NewGoal(
		"user123",
		entities.GoalTypeSavings,
		"新車購入資金",
		targetAmount,
		targetDate,
		monthlyContribution,
	)
	if err != nil {
		t.Fatalf("目標の作成に失敗しました: %v", err)
	}

	// 目標を追加
	err = plan.AddGoal(goal)
	if err != nil {
		t.Fatalf("目標の追加に失敗しました: %v", err)
	}

	// 検証
	if len(plan.Goals()) != 1 {
		t.Error("目標が正しく追加されていません")
	}

	if plan.Goals()[0].ID() != goal.ID() {
		t.Error("追加された目標が正しくありません")
	}
}

func TestAddDuplicateRetirementGoal(t *testing.T) {
	plan := createTestFinancialPlan(t)

	// 最初の退職目標を追加
	goal1 := createTestRetirementGoal(t)
	err := plan.AddGoal(goal1)
	if err != nil {
		t.Fatalf("最初の退職目標の追加に失敗しました: %v", err)
	}

	// 2つ目の退職目標を追加（エラーになるはず）
	goal2 := createTestRetirementGoal(t)
	err = plan.AddGoal(goal2)
	if err == nil {
		t.Error("重複する退職目標の追加がエラーになりませんでした")
	}
}

func TestGenerateProjection(t *testing.T) {
	plan := createTestFinancialPlan(t)

	// 予測を生成
	projection, err := plan.GenerateProjection(10)
	if err != nil {
		t.Fatalf("予測の生成に失敗しました: %v", err)
	}

	// 検証
	if len(projection.AssetProjections) != 10 {
		t.Errorf("資産推移予測の年数が正しくありません。期待値: 10, 実際: %d", len(projection.AssetProjections))
	}

	if projection.EmergencyFundStatus == nil {
		t.Error("緊急資金状況が生成されていません")
	}

	if len(projection.GoalProgress) != 0 {
		t.Error("目標がない状態では目標進捗は空である必要があります")
	}
}

func TestValidatePlan(t *testing.T) {
	plan := createTestFinancialPlan(t)

	// 緊急資金を適切に設定
	emergencyConfig, _ := NewEmergencyFundConfig(3, mustCreateMoney(540000)) // 3ヶ月分の支出
	err := plan.UpdateEmergencyFund(emergencyConfig)
	if err != nil {
		t.Fatalf("緊急資金設定の更新に失敗しました: %v", err)
	}

	// バリデーションを実行
	errors := plan.ValidatePlan()

	// 適切に設定された財務計画の場合はエラーなし
	if len(errors) != 0 {
		t.Errorf("適切に設定された財務計画でバリデーションエラーが発生しました: %v", errors)
	}
}

// ヘルパー関数
func createTestFinancialPlan(t *testing.T) *FinancialPlan {
	monthlyIncome, _ := valueobjects.NewMoneyJPY(400000)
	expenses := entities.ExpenseCollection{
		{Category: "住居費", Amount: mustCreateMoney(120000)},
		{Category: "食費", Amount: mustCreateMoney(60000)},
		{Category: "その他", Amount: mustCreateMoney(80000)},
	}
	savings := entities.SavingsCollection{
		{Type: "deposit", Amount: mustCreateMoney(1000000)},
	}
	investmentReturn, _ := valueobjects.NewRate(5.0)
	inflationRate, _ := valueobjects.NewRate(2.0)

	profile, err := entities.NewFinancialProfile(
		"user123",
		monthlyIncome,
		expenses,
		savings,
		investmentReturn,
		inflationRate,
	)
	if err != nil {
		t.Fatalf("テスト用財務プロファイルの作成に失敗しました: %v", err)
	}

	plan, err := NewFinancialPlan(profile)
	if err != nil {
		t.Fatalf("テスト用財務計画の作成に失敗しました: %v", err)
	}

	return plan
}

func createTestRetirementGoal(t *testing.T) *entities.Goal {
	targetAmount, _ := valueobjects.NewMoneyJPY(30000000)
	monthlyContribution, _ := valueobjects.NewMoneyJPY(100000)
	targetDate := time.Now().AddDate(20, 0, 0) // 20年後

	goal, err := entities.NewGoal(
		"user123",
		entities.GoalTypeRetirement,
		"老後資金",
		targetAmount,
		targetDate,
		monthlyContribution,
	)
	if err != nil {
		t.Fatalf("テスト用退職目標の作成に失敗しました: %v", err)
	}

	return goal
}

func mustCreateMoney(amount float64) valueobjects.Money {
	money, err := valueobjects.NewMoneyJPY(amount)
	if err != nil {
		panic(err)
	}
	return money
}
