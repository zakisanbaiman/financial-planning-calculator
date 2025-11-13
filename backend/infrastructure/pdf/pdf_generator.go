package pdf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/financial-planning-calculator/backend/application/usecases"
)

// Generator はPDF生成インターフェース
type Generator interface {
	GenerateFinancialSummaryPDF(report *usecases.FinancialSummaryReport) ([]byte, error)
	GenerateComprehensivePDF(report *usecases.ComprehensiveReport) ([]byte, error)
	GenerateAssetProjectionPDF(report *usecases.AssetProjectionReport) ([]byte, error)
	GenerateGoalsProgressPDF(report *usecases.GoalsProgressReport) ([]byte, error)
	GenerateRetirementPlanPDF(report *usecases.RetirementPlanReport) ([]byte, error)
}

// HTMLGenerator はHTML形式でPDFを生成する（簡易実装）
type HTMLGenerator struct{}

// NewHTMLGenerator は新しいHTMLGeneratorを作成する
func NewHTMLGenerator() *HTMLGenerator {
	return &HTMLGenerator{}
}

// GenerateFinancialSummaryPDF は財務サマリーレポートのPDFを生成する
func (g *HTMLGenerator) GenerateFinancialSummaryPDF(report *usecases.FinancialSummaryReport) ([]byte, error) {
	html := g.generateFinancialSummaryHTML(report)
	return []byte(html), nil
}

// GenerateComprehensivePDF は包括的レポートのPDFを生成する
func (g *HTMLGenerator) GenerateComprehensivePDF(report *usecases.ComprehensiveReport) ([]byte, error) {
	html := g.generateComprehensiveHTML(report)
	return []byte(html), nil
}

// GenerateAssetProjectionPDF は資産推移レポートのPDFを生成する
func (g *HTMLGenerator) GenerateAssetProjectionPDF(report *usecases.AssetProjectionReport) ([]byte, error) {
	html := g.generateAssetProjectionHTML(report)
	return []byte(html), nil
}

// GenerateGoalsProgressPDF は目標進捗レポートのPDFを生成する
func (g *HTMLGenerator) GenerateGoalsProgressPDF(report *usecases.GoalsProgressReport) ([]byte, error) {
	html := g.generateGoalsProgressHTML(report)
	return []byte(html), nil
}

// GenerateRetirementPlanPDF は退職計画レポートのPDFを生成する
func (g *HTMLGenerator) GenerateRetirementPlanPDF(report *usecases.RetirementPlanReport) ([]byte, error) {
	html := g.generateRetirementPlanHTML(report)
	return []byte(html), nil
}

// generateFinancialSummaryHTML は財務サマリーのHTML生成
func (g *HTMLGenerator) generateFinancialSummaryHTML(report *usecases.FinancialSummaryReport) string {
	var buf bytes.Buffer

	buf.WriteString(`<!DOCTYPE html>
<html lang="ja">
<head>
    <meta charset="UTF-8">
    <title>財務サマリーレポート</title>
    <style>
        body { font-family: 'Helvetica', 'Arial', sans-serif; margin: 40px; color: #333; }
        h1 { color: #2563eb; border-bottom: 3px solid #2563eb; padding-bottom: 10px; }
        h2 { color: #1e40af; margin-top: 30px; border-bottom: 2px solid #ddd; padding-bottom: 5px; }
        h3 { color: #1e3a8a; margin-top: 20px; }
        .header { text-align: center; margin-bottom: 40px; }
        .section { margin-bottom: 30px; }
        .metric { display: inline-block; width: 45%; margin: 10px 2%; padding: 15px; background: #f3f4f6; border-radius: 8px; }
        .metric-label { font-size: 14px; color: #6b7280; }
        .metric-value { font-size: 24px; font-weight: bold; color: #111827; }
        .score { text-align: center; padding: 20px; background: #dbeafe; border-radius: 10px; margin: 20px 0; }
        .score-value { font-size: 48px; font-weight: bold; color: #2563eb; }
        .score-level { font-size: 20px; color: #1e40af; text-transform: uppercase; }
        .list-item { padding: 10px; margin: 5px 0; background: #f9fafb; border-left: 4px solid #2563eb; }
        .warning { border-left-color: #ef4444; background: #fef2f2; }
        table { width: 100%; border-collapse: collapse; margin: 20px 0; }
        th, td { padding: 12px; text-align: left; border-bottom: 1px solid #e5e7eb; }
        th { background: #f3f4f6; font-weight: 600; }
        .footer { margin-top: 50px; text-align: center; font-size: 12px; color: #6b7280; }
    </style>
</head>
<body>
    <div class="header">
        <h1>財務サマリーレポート</h1>
        <p>作成日: ` + report.ReportDate + `</p>
    </div>

    <div class="section">
        <h2>財務健全性スコア</h2>
        <div class="score">
            <div class="score-value">` + fmt.Sprintf("%d", report.FinancialHealth.OverallScore) + `/100</div>
            <div class="score-level">` + g.getScoreLevelText(report.FinancialHealth.ScoreLevel) + `</div>
        </div>
        <div class="metric">
            <div class="metric-label">貯蓄率</div>
            <div class="metric-value">` + fmt.Sprintf("%.1f%%", report.FinancialHealth.SavingsRate) + `</div>
        </div>
        <div class="metric">
            <div class="metric-label">緊急資金比率</div>
            <div class="metric-value">` + fmt.Sprintf("%.1fヶ月", report.FinancialHealth.EmergencyFundRatio) + `</div>
        </div>
    </div>

    <div class="section">
        <h2>現在の財務状況</h2>
        <div class="metric">
            <div class="metric-label">月収</div>
            <div class="metric-value">¥` + g.formatNumber(report.CurrentSituation.MonthlyIncome) + `</div>
        </div>
        <div class="metric">
            <div class="metric-label">月間支出</div>
            <div class="metric-value">¥` + g.formatNumber(report.CurrentSituation.MonthlyExpenses) + `</div>
        </div>
        <div class="metric">
            <div class="metric-label">純貯蓄額</div>
            <div class="metric-value">¥` + g.formatNumber(report.CurrentSituation.NetSavings) + `</div>
        </div>
        <div class="metric">
            <div class="metric-label">総資産</div>
            <div class="metric-value">¥` + g.formatNumber(report.CurrentSituation.TotalAssets) + `</div>
        </div>
    </div>

    <div class="section">
        <h2>主要指標</h2>
        <table>
            <thead>
                <tr>
                    <th>指標名</th>
                    <th>値</th>
                    <th>説明</th>
                    <th>トレンド</th>
                </tr>
            </thead>
            <tbody>`)

	for _, metric := range report.KeyMetrics {
		buf.WriteString(`
                <tr>
                    <td>` + metric.Name + `</td>
                    <td>` + g.formatMetricValue(metric.Value, metric.Unit) + `</td>
                    <td>` + metric.Description + `</td>
                    <td>` + g.getTrendIcon(metric.Trend) + `</td>
                </tr>`)
	}

	buf.WriteString(`
            </tbody>
        </table>
    </div>`)

	if len(report.Recommendations) > 0 {
		buf.WriteString(`
    <div class="section">
        <h2>推奨事項</h2>`)
		for _, rec := range report.Recommendations {
			buf.WriteString(`
        <div class="list-item">` + rec + `</div>`)
		}
		buf.WriteString(`
    </div>`)
	}

	if len(report.Warnings) > 0 {
		buf.WriteString(`
    <div class="section">
        <h2>注意事項</h2>`)
		for _, warning := range report.Warnings {
			buf.WriteString(`
        <div class="list-item warning">⚠️ ` + warning + `</div>`)
		}
		buf.WriteString(`
    </div>`)
	}

	buf.WriteString(`
    <div class="footer">
        <p>このレポートは ` + time.Now().Format("2006年01月02日 15:04") + ` に生成されました</p>
        <p>Financial Planning Calculator</p>
    </div>
</body>
</html>`)

	return buf.String()
}

// generateComprehensiveHTML は包括的レポートのHTML生成
func (g *HTMLGenerator) generateComprehensiveHTML(report *usecases.ComprehensiveReport) string {
	var buf bytes.Buffer

	buf.WriteString(`<!DOCTYPE html>
<html lang="ja">
<head>
    <meta charset="UTF-8">
    <title>包括的財務レポート</title>
    <style>
        body { font-family: 'Helvetica', 'Arial', sans-serif; margin: 40px; color: #333; line-height: 1.6; }
        h1 { color: #2563eb; border-bottom: 3px solid #2563eb; padding-bottom: 10px; }
        h2 { color: #1e40af; margin-top: 30px; border-bottom: 2px solid #ddd; padding-bottom: 5px; page-break-before: always; }
        h3 { color: #1e3a8a; margin-top: 20px; }
        .header { text-align: center; margin-bottom: 40px; }
        .executive-summary { background: #dbeafe; padding: 20px; border-radius: 10px; margin: 20px 0; }
        .highlight { background: #fef3c7; padding: 10px; margin: 10px 0; border-left: 4px solid #f59e0b; }
        .action-item { padding: 15px; margin: 10px 0; background: #f3f4f6; border-radius: 8px; border-left: 4px solid #10b981; }
        .priority-high { border-left-color: #ef4444; }
        .priority-medium { border-left-color: #f59e0b; }
        .priority-low { border-left-color: #10b981; }
        .metric { display: inline-block; width: 45%; margin: 10px 2%; padding: 15px; background: #f3f4f6; border-radius: 8px; }
        .metric-label { font-size: 14px; color: #6b7280; }
        .metric-value { font-size: 24px; font-weight: bold; color: #111827; }
        .footer { margin-top: 50px; text-align: center; font-size: 12px; color: #6b7280; page-break-before: always; }
        table { width: 100%; border-collapse: collapse; margin: 20px 0; }
        th, td { padding: 12px; text-align: left; border-bottom: 1px solid #e5e7eb; }
        th { background: #f3f4f6; font-weight: 600; }
    </style>
</head>
<body>
    <div class="header">
        <h1>包括的財務レポート</h1>
        <p>作成日: ` + time.Now().Format("2006年01月02日") + `</p>
    </div>

    <div class="executive-summary">
        <h2 style="border: none; margin-top: 0;">エグゼクティブサマリー</h2>
        <p><strong>総合ステータス:</strong> ` + report.ExecutiveSummary.OverallStatus + `</p>
        <p><strong>財務健全性スコア:</strong> ` + fmt.Sprintf("%d/100", report.ExecutiveSummary.FinancialHealthScore) + `</p>
        
        <h3>主要ハイライト</h3>`)

	for _, highlight := range report.ExecutiveSummary.KeyHighlights {
		buf.WriteString(`
        <div class="highlight">✓ ` + highlight + `</div>`)
	}

	buf.WriteString(`
        
        <h3>重要アクション</h3>`)

	for _, action := range report.ExecutiveSummary.CriticalActions {
		buf.WriteString(`
        <div class="highlight">⚡ ` + action + `</div>`)
	}

	buf.WriteString(`
    </div>`)

	// 財務サマリーセクション
	buf.WriteString(`
    <h2>財務サマリー</h2>
    <div class="metric">
        <div class="metric-label">月収</div>
        <div class="metric-value">¥` + g.formatNumber(report.FinancialSummary.CurrentSituation.MonthlyIncome) + `</div>
    </div>
    <div class="metric">
        <div class="metric-label">月間支出</div>
        <div class="metric-value">¥` + g.formatNumber(report.FinancialSummary.CurrentSituation.MonthlyExpenses) + `</div>
    </div>
    <div class="metric">
        <div class="metric-label">純貯蓄額</div>
        <div class="metric-value">¥` + g.formatNumber(report.FinancialSummary.CurrentSituation.NetSavings) + `</div>
    </div>
    <div class="metric">
        <div class="metric-label">総資産</div>
        <div class="metric-value">¥` + g.formatNumber(report.FinancialSummary.CurrentSituation.TotalAssets) + `</div>
    </div>`)

	// 資産推移セクション
	buf.WriteString(`
    <h2>資産推移予測</h2>
    <p>予測期間: ` + fmt.Sprintf("%d年", report.AssetProjection.ProjectionYears) + `</p>
    <table>
        <thead>
            <tr>
                <th>年</th>
                <th>総資産</th>
                <th>実質価値</th>
                <th>積立元本</th>
                <th>投資収益</th>
            </tr>
        </thead>
        <tbody>`)

	// 最初、中間、最後の年のみ表示
	projections := report.AssetProjection.Projections
	displayYears := []int{0}
	if len(projections) > 2 {
		displayYears = append(displayYears, len(projections)/2)
	}
	if len(projections) > 1 {
		displayYears = append(displayYears, len(projections)-1)
	}

	for _, idx := range displayYears {
		if idx < len(projections) {
			p := projections[idx]
			buf.WriteString(`
            <tr>
                <td>` + fmt.Sprintf("%d年後", p.Year) + `</td>
                <td>¥` + g.formatNumber(p.TotalAssets.Amount()) + `</td>
                <td>¥` + g.formatNumber(p.RealValue.Amount()) + `</td>
                <td>¥` + g.formatNumber(p.ContributedAmount.Amount()) + `</td>
                <td>¥` + g.formatNumber(p.InvestmentGains.Amount()) + `</td>
            </tr>`)
		}
	}

	buf.WriteString(`
        </tbody>
    </table>`)

	// 目標進捗セクション
	buf.WriteString(`
    <h2>目標進捗状況</h2>
    <p>総目標数: ` + fmt.Sprintf("%d", report.GoalsProgress.Summary.TotalGoals) + ` (アクティブ: ` + fmt.Sprintf("%d", report.GoalsProgress.Summary.ActiveGoals) + `)</p>
    <p>全体進捗率: ` + fmt.Sprintf("%.1f%%", report.GoalsProgress.Summary.OverallProgress) + `</p>
    
    <table>
        <thead>
            <tr>
                <th>目標</th>
                <th>目標金額</th>
                <th>現在額</th>
                <th>進捗率</th>
                <th>ステータス</th>
            </tr>
        </thead>
        <tbody>`)

	for _, goalProgress := range report.GoalsProgress.Goals {
		buf.WriteString(`
            <tr>
                <td>` + goalProgress.Goal.Title() + `</td>
                <td>¥` + g.formatNumber(goalProgress.Goal.TargetAmount().Amount()) + `</td>
                <td>¥` + g.formatNumber(goalProgress.Goal.CurrentAmount().Amount()) + `</td>
                <td>` + fmt.Sprintf("%.1f%%", goalProgress.Progress.AsPercentage()) + `</td>
                <td>` + goalProgress.Status + `</td>
            </tr>`)
	}

	buf.WriteString(`
        </tbody>
    </table>`)

	// アクションプラン
	buf.WriteString(`
    <h2>アクションプラン</h2>
    
    <h3>短期アクション（3ヶ月以内）</h3>`)

	for _, action := range report.ActionPlan.ShortTerm {
		priorityClass := "priority-" + action.Priority
		buf.WriteString(`
    <div class="action-item ` + priorityClass + `">
        <strong>` + action.Title + `</strong> [優先度: ` + action.Priority + `]<br>
        ` + action.Description + `<br>
        <small>期限: ` + action.Timeline + ` | 影響: ` + action.Impact + ` | 労力: ` + action.Effort + `</small>
    </div>`)
	}

	buf.WriteString(`
    
    <h3>中期アクション（1年以内）</h3>`)

	for _, action := range report.ActionPlan.MediumTerm {
		priorityClass := "priority-" + action.Priority
		buf.WriteString(`
    <div class="action-item ` + priorityClass + `">
        <strong>` + action.Title + `</strong> [優先度: ` + action.Priority + `]<br>
        ` + action.Description + `<br>
        <small>期限: ` + action.Timeline + ` | 影響: ` + action.Impact + ` | 労力: ` + action.Effort + `</small>
    </div>`)
	}

	buf.WriteString(`
    
    <h3>長期アクション（1年以上）</h3>`)

	for _, action := range report.ActionPlan.LongTerm {
		priorityClass := "priority-" + action.Priority
		buf.WriteString(`
    <div class="action-item ` + priorityClass + `">
        <strong>` + action.Title + `</strong> [優先度: ` + action.Priority + `]<br>
        ` + action.Description + `<br>
        <small>期限: ` + action.Timeline + ` | 影響: ` + action.Impact + ` | 労力: ` + action.Effort + `</small>
    </div>`)
	}

	buf.WriteString(`
    <div class="footer">
        <p>このレポートは ` + time.Now().Format("2006年01月02日 15:04") + ` に生成されました</p>
        <p>Financial Planning Calculator</p>
    </div>
</body>
</html>`)

	return buf.String()
}

// generateAssetProjectionHTML は資産推移レポートのHTML生成（簡略版）
func (g *HTMLGenerator) generateAssetProjectionHTML(report *usecases.AssetProjectionReport) string {
	// 簡略化のため、基本的な構造のみ実装
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="ja">
<head><meta charset="UTF-8"><title>資産推移レポート</title></head>
<body>
<h1>資産推移レポート</h1>
<p>予測期間: %d年</p>
<p>レポート生成日: %s</p>
</body>
</html>`, report.ProjectionYears, time.Now().Format("2006-01-02"))
}

// generateGoalsProgressHTML は目標進捗レポートのHTML生成（簡略版）
func (g *HTMLGenerator) generateGoalsProgressHTML(report *usecases.GoalsProgressReport) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="ja">
<head><meta charset="UTF-8"><title>目標進捗レポート</title></head>
<body>
<h1>目標進捗レポート</h1>
<p>総目標数: %d</p>
<p>レポート生成日: %s</p>
</body>
</html>`, report.Summary.TotalGoals, time.Now().Format("2006-01-02"))
}

// generateRetirementPlanHTML は退職計画レポートのHTML生成（簡略版）
func (g *HTMLGenerator) generateRetirementPlanHTML(report *usecases.RetirementPlanReport) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="ja">
<head><meta charset="UTF-8"><title>退職計画レポート</title></head>
<body>
<h1>退職計画レポート</h1>
<p>レポート生成日: %s</p>
</body>
</html>`, time.Now().Format("2006-01-02"))
}

// ヘルパー関数

func (g *HTMLGenerator) formatNumber(num float64) string {
	// カンマ区切りで数値をフォーマット
	str := fmt.Sprintf("%.0f", num)

	// 3桁ごとにカンマを挿入
	if len(str) <= 3 {
		return str
	}

	var result string
	for i, c := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result += ","
		}
		result += string(c)
	}
	return result
}

func (g *HTMLGenerator) formatMetricValue(value float64, unit string) string {
	if unit == "%" {
		return fmt.Sprintf("%.1f%%", value)
	} else if unit == "円" {
		return "¥" + g.formatNumber(value)
	}
	return fmt.Sprintf("%.2f %s", value, unit)
}

func (g *HTMLGenerator) getScoreLevelText(level string) string {
	switch level {
	case "excellent":
		return "優秀"
	case "good":
		return "良好"
	case "fair":
		return "普通"
	case "poor":
		return "要改善"
	default:
		return level
	}
}

func (g *HTMLGenerator) getTrendIcon(trend string) string {
	switch trend {
	case "up":
		return "↑ 上昇"
	case "down":
		return "↓ 下降"
	case "stable":
		return "→ 安定"
	default:
		return trend
	}
}

// JSONGenerator はJSON形式でレポートを生成する
type JSONGenerator struct{}

// NewJSONGenerator は新しいJSONGeneratorを作成する
func NewJSONGenerator() *JSONGenerator {
	return &JSONGenerator{}
}

// GenerateFinancialSummaryPDF は財務サマリーレポートのJSONを生成する
func (g *JSONGenerator) GenerateFinancialSummaryPDF(report *usecases.FinancialSummaryReport) ([]byte, error) {
	return json.MarshalIndent(report, "", "  ")
}

// GenerateComprehensivePDF は包括的レポートのJSONを生成する
func (g *JSONGenerator) GenerateComprehensivePDF(report *usecases.ComprehensiveReport) ([]byte, error) {
	return json.MarshalIndent(report, "", "  ")
}

// GenerateAssetProjectionPDF は資産推移レポートのJSONを生成する
func (g *JSONGenerator) GenerateAssetProjectionPDF(report *usecases.AssetProjectionReport) ([]byte, error) {
	return json.MarshalIndent(report, "", "  ")
}

// GenerateGoalsProgressPDF は目標進捗レポートのJSONを生成する
func (g *JSONGenerator) GenerateGoalsProgressPDF(report *usecases.GoalsProgressReport) ([]byte, error) {
	return json.MarshalIndent(report, "", "  ")
}

// GenerateRetirementPlanPDF は退職計画レポートのJSONを生成する
func (g *JSONGenerator) GenerateRetirementPlanPDF(report *usecases.RetirementPlanReport) ([]byte, error) {
	return json.MarshalIndent(report, "", "  ")
}
