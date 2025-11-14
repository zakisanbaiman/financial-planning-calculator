package pdf

import (
	"testing"

	"github.com/financial-planning-calculator/backend/application/usecases"
	"github.com/financial-planning-calculator/backend/domain/entities"
)

func TestHTMLGenerator_FormatNumber(t *testing.T) {
	generator := NewHTMLGenerator()

	tests := []struct {
		name     string
		input    float64
		expected string
	}{
		{"Small number", 123, "123"},
		{"Thousand", 1234, "1,234"},
		{"Million", 1234567, "1,234,567"},
		{"Large number", 12345678.90, "12,345,679"},
		{"Zero", 0, "0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generator.formatNumber(tt.input)
			if result != tt.expected {
				t.Errorf("formatNumber(%f) = %s; want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestHTMLGenerator_GenerateFinancialSummaryPDF(t *testing.T) {
	generator := NewHTMLGenerator()

	report := &usecases.FinancialSummaryReport{
		UserID:     entities.UserID("test-user"),
		ReportDate: "2024-11-13",
		FinancialHealth: usecases.FinancialHealth{
			OverallScore:       75,
			ScoreLevel:         "good",
			SavingsRate:        20.5,
			DebtToIncomeRatio:  0,
			EmergencyFundRatio: 4.5,
		},
		CurrentSituation: usecases.CurrentSituation{
			MonthlyIncome:    400000,
			MonthlyExpenses:  280000,
			NetSavings:       120000,
			TotalAssets:      1500000,
			InvestmentReturn: 5.0,
			InflationRate:    2.0,
		},
		KeyMetrics: []usecases.KeyMetric{
			{
				Name:        "貯蓄率",
				Value:       20.5,
				Unit:        "%",
				Description: "月収に対する純貯蓄額の割合",
				Trend:       "stable",
			},
		},
		Recommendations: []string{"月間支出を詳細に分析してください"},
		Warnings:        []string{"緊急資金が3ヶ月分の生活費を下回っています"},
	}

	html, err := generator.GenerateFinancialSummaryPDF(report)
	if err != nil {
		t.Fatalf("GenerateFinancialSummaryPDF failed: %v", err)
	}

	if len(html) == 0 {
		t.Error("Generated HTML is empty")
	}

	// HTMLの基本的な構造をチェック
	htmlStr := string(html)
	requiredElements := []string{
		"<!DOCTYPE html>",
		"財務サマリーレポート",
		"財務健全性スコア",
		"75/100",
		"良好",
		"現在の財務状況",
		"¥400,000",
		"主要指標",
		"推奨事項",
		"注意事項",
	}

	for _, element := range requiredElements {
		if !contains(htmlStr, element) {
			t.Errorf("Generated HTML does not contain expected element: %s", element)
		}
	}
}

func TestHTMLGenerator_GenerateComprehensivePDF(t *testing.T) {
	generator := NewHTMLGenerator()

	report := &usecases.ComprehensiveReport{
		UserID: entities.UserID("test-user"),
		ExecutiveSummary: usecases.ExecutiveSummary{
			OverallStatus:        "良好",
			KeyHighlights:        []string{"貯蓄率が理想的"},
			CriticalActions:      []string{"緊急資金の確保"},
			OpportunityAreas:     []string{"投資利回りの改善"},
			FinancialHealthScore: 75,
		},
		FinancialSummary: usecases.FinancialSummaryReport{
			CurrentSituation: usecases.CurrentSituation{
				MonthlyIncome:   400000,
				MonthlyExpenses: 280000,
				NetSavings:      120000,
				TotalAssets:     1500000,
			},
		},
		AssetProjection: usecases.AssetProjectionReport{
			ProjectionYears: 10,
			Projections:     []entities.AssetProjection{},
		},
		GoalsProgress: usecases.GoalsProgressReport{
			Summary: usecases.GoalsSummary{
				TotalGoals:      3,
				ActiveGoals:     2,
				OverallProgress: 65.5,
			},
		},
		ActionPlan: usecases.ActionPlan{
			ShortTerm: []usecases.ActionItem{
				{
					Priority:    "high",
					Title:       "緊急資金の確保",
					Description: "3ヶ月分の生活費を緊急資金として確保する",
					Timeline:    "3ヶ月以内",
					Impact:      "リスク軽減",
					Effort:      "medium",
				},
			},
			MediumTerm: []usecases.ActionItem{},
			LongTerm:   []usecases.ActionItem{},
		},
	}

	html, err := generator.GenerateComprehensivePDF(report)
	if err != nil {
		t.Fatalf("GenerateComprehensivePDF failed: %v", err)
	}

	if len(html) == 0 {
		t.Error("Generated HTML is empty")
	}

	htmlStr := string(html)
	requiredElements := []string{
		"<!DOCTYPE html>",
		"包括的財務レポート",
		"エグゼクティブサマリー",
		"良好",
		"財務サマリー",
		"資産推移予測",
		"目標進捗状況",
		"アクションプラン",
		"短期アクション",
	}

	for _, element := range requiredElements {
		if !contains(htmlStr, element) {
			t.Errorf("Generated HTML does not contain expected element: %s", element)
		}
	}
}

func TestJSONGenerator_GenerateFinancialSummaryPDF(t *testing.T) {
	generator := NewJSONGenerator()

	report := &usecases.FinancialSummaryReport{
		UserID:     entities.UserID("test-user"),
		ReportDate: "2024-11-13",
		FinancialHealth: usecases.FinancialHealth{
			OverallScore: 75,
			ScoreLevel:   "good",
		},
	}

	jsonData, err := generator.GenerateFinancialSummaryPDF(report)
	if err != nil {
		t.Fatalf("GenerateFinancialSummaryPDF failed: %v", err)
	}

	if len(jsonData) == 0 {
		t.Error("Generated JSON is empty")
	}

	// JSONの基本的な構造をチェック
	jsonStr := string(jsonData)
	if !contains(jsonStr, "test-user") {
		t.Error("Generated JSON does not contain user_id")
	}
	if !contains(jsonStr, "2024-11-13") {
		t.Error("Generated JSON does not contain report_date")
	}
}

// ヘルパー関数
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
