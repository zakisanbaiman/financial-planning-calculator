package services

import (
	"testing"
	"time"

	"financial-planning-calculator/domain/entities"
	"financial-planning-calculator/domain/valueobjects"
)

func TestRecommendMonthlySavings(t *testing.T) {
	calculationService := NewFinancialCalculationService()
	service := NewGoalRecommendationService(calculationService)

	// テスト用の目標を作成
	goal := createTestGoal(t)
	currentSavings, _ := valueobjects.NewMoneyJPY(100000)
	timeRemaining, _ := valueobjects.NewPeriodFromMonths(24) // 2年

	recommendation, err := service.RecommendMonthlySavings(goal, currentSavings, timeRemaining)
	if err != nil {
		t.Fatalf("月間貯蓄推奨の計算に失敗しました: %v", err)
	}

	// 検証: 推奨金額は正の値であるはず
	if !recommendation.RecommendedAmount.IsPositive() {
		t.Error("推奨月間貯蓄額が正の値ではありません")
	}

	// 優先度が設定されているはず
	if recommendation.Priority == "" {
		t.Error("優先度が設定されていません")
	}

	// 根拠が設定されているはず
	if recommendation.Rationale == "" {
		t.Error("根拠が設定されていません")
	}
}

func TestRecommendMonthlySavingsForCompletedGoal(t *testing.T) {
	calculationService := NewFinancialCalculationService()
	service := NewGoalRecommendationService(calculationService)

	// 既に達成済みの目標を作成
	targetAmount, _ := valueobjects.NewMoneyJPY(1000000)
	monthlyContribution, _ := valueobjects.NewMoneyJPY(50000)
	targetDate := time.Now().AddDate(2, 0, 0)

	goal, err := entities.NewGoal(
		"user123",
		entities.GoalTypeSavings,
		"テスト目標",
		targetAmount,
		targetDate,
		monthlyContribution,
	)
	if err != nil {
		t.Fatalf("目標の作成に失敗しました: %v", err)
	}

	// 目標金額を既に達成済みに設定
	err = goal.UpdateCurrentAmount(targetAmount)
	if err != nil {
		t.Fatalf("現在金額の更新に失敗しました: %v", err)
	}

	currentSavings, _ := valueobjects.NewMoneyJPY(1000000)
	timeRemaining, _ := valueobjects.NewPeriodFromMonths(12)

	recommendation, err := service.RecommendMonthlySavings(goal, currentSavings, timeRemaining)
	if err != nil {
		t.Fatalf("達成済み目標の推奨計算に失敗しました: %v", err)
	}

	// 検証: 達成済みの場合は推奨額は0
	if !recommendation.RecommendedAmount.IsZero() {
		t.Error("達成済み目標の推奨月間貯蓄額は0であるべきです")
	}

	// 優先度は低であるはず
	if recommendation.Priority != PriorityLow {
		t.Error("達成済み目標の優先度は低であるべきです")
	}
}

func TestSuggestGoalAdjustments(t *testing.T) {
	calculationService := NewFinancialCalculationService()
	service := NewGoalRecommendationService(calculationService)

	// 達成困難な目標を作成
	goal := createDifficultGoal(t)
	profile := createTestFinancialProfile(t)

	recommendations, err := service.SuggestGoalAdjustments(goal, profile)
	if err != nil {
		t.Fatalf("目標調整提案の計算に失敗しました: %v", err)
	}

	// 検証: 何らかの推奨事項が提案されるはず
	if len(recommendations) == 0 {
		t.Error("達成困難な目標に対して推奨事項が提案されませんでした")
	}

	// 各推奨事項の基本フィールドが設定されているかチェック
	for i, rec := range recommendations {
		if rec.Type == "" {
			t.Errorf("推奨事項[%d]のタイプが設定されていません", i)
		}
		if rec.Title == "" {
			t.Errorf("推奨事項[%d]のタイトルが設定されていません", i)
		}
		if rec.Description == "" {
			t.Errorf("推奨事項[%d]の説明が設定されていません", i)
		}
		if rec.Priority == "" {
			t.Errorf("推奨事項[%d]の優先度が設定されていません", i)
		}
	}
}

func TestSuggestGoalAdjustmentsForAchievableGoal(t *testing.T) {
	calculationService := NewFinancialCalculationService()
	service := NewGoalRecommendationService(calculationService)

	// 達成可能な目標を作成
	goal := createAchievableGoal(t)
	profile := createTestFinancialProfile(t)

	recommendations, err := service.SuggestGoalAdjustments(goal, profile)
	if err != nil {
		t.Fatalf("達成可能目標の調整提案計算に失敗しました: %v", err)
	}

	// 検証: 達成可能な目標には推奨事項なし
	if len(recommendations) != 0 {
		t.Error("達成可能な目標に対して不要な推奨事項が提案されました")
	}
}

func TestAnalyzeGoalFeasibility(t *testing.T) {
	calculationService := NewFinancialCalculationService()
	service := NewGoalRecommendationService(calculationService)

	goal := createTestGoal(t)
	profile := createTestFinancialProfile(t)

	analysis, err := service.AnalyzeGoalFeasibility(goal, profile)
	if err != nil {
		t.Fatalf("目標実現可能性分析に失敗しました: %v", err)
	}

	// 検証: 必要なフィールドが含まれているかチェック
	requiredFields := []string{
		"goal_type",
		"target_amount",
		"current_amount",
		"remaining_days",
		"net_savings",
		"required_monthly_savings",
		"achievable",
		"progress_percentage",
		"risk_level",
	}

	for _, field := range requiredFields {
		if _, exists := analysis[field]; !exists {
			t.Errorf("分析結果に必要なフィールド '%s' が含まれていません", field)
		}
	}

	// 型チェック
	if _, ok := analysis["achievable"].(bool); !ok {
		t.Error("achievableフィールドがbool型ではありません")
	}

	if _, ok := analysis["target_amount"].(float64); !ok {
		t.Error("target_amountフィールドがfloat64型ではありません")
	}

	if _, ok := analysis["risk_level"].(string); !ok {
		t.Error("risk_levelフィールドがstring型ではありません")
	}
}

// ヘルパー関数
func createTestGoal(t *testing.T) *entities.Goal {
	targetAmount, _ := valueobjects.NewMoneyJPY(2000000)
	monthlyContribution, _ := valueobjects.NewMoneyJPY(50000)
	targetDate := time.Now().AddDate(3, 0, 0) // 3年後

	goal, err := entities.NewGoal(
		"user123",
		entities.GoalTypeSavings,
		"テスト目標",
		targetAmount,
		targetDate,
		monthlyContribution,
	)
	if err != nil {
		t.Fatalf("テスト目標の作成に失敗しました: %v", err)
	}

	return goal
}

func createDifficultGoal(t *testing.T) *entities.Goal {
	// 非常に高額で短期間の目標（達成困難）
	targetAmount, _ := valueobjects.NewMoneyJPY(10000000) // 1000万円
	monthlyContribution, _ := valueobjects.NewMoneyJPY(50000)
	targetDate := time.Now().AddDate(1, 0, 0) // 1年後

	goal, err := entities.NewGoal(
		"user123",
		entities.GoalTypeSavings,
		"困難な目標",
		targetAmount,
		targetDate,
		monthlyContribution,
	)
	if err != nil {
		t.Fatalf("困難な目標の作成に失敗しました: %v", err)
	}

	return goal
}

func createAchievableGoal(t *testing.T) *entities.Goal {
	// 現実的な金額と期間の目標
	targetAmount, _ := valueobjects.NewMoneyJPY(1200000) // 120万円
	monthlyContribution, _ := valueobjects.NewMoneyJPY(50000)
	targetDate := time.Now().AddDate(2, 0, 0) // 2年後

	goal, err := entities.NewGoal(
		"user123",
		entities.GoalTypeSavings,
		"達成可能な目標",
		targetAmount,
		targetDate,
		monthlyContribution,
	)
	if err != nil {
		t.Fatalf("達成可能な目標の作成に失敗しました: %v", err)
	}

	return goal
}

func createTestFinancialProfile(t *testing.T) *entities.FinancialProfile {
	monthlyIncome, _ := valueobjects.NewMoneyJPY(400000)
	expenses := entities.ExpenseCollection{
		{Category: "住居費", Amount: mustCreateMoneyForTest(120000)},
		{Category: "食費", Amount: mustCreateMoneyForTest(60000)},
		{Category: "その他", Amount: mustCreateMoneyForTest(80000)},
	}
	savings := entities.SavingsCollection{
		{Type: "deposit", Amount: mustCreateMoneyForTest(1000000)},
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

	return profile
}

func mustCreateMoneyForTest(amount float64) valueobjects.Money {
	money, err := valueobjects.NewMoneyJPY(amount)
	if err != nil {
		panic(err)
	}
	return money
}
