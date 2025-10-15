package services

import (
	"errors"
	"fmt"
	"math"
	"time"

	"financial-planning-calculator/domain/entities"
	"financial-planning-calculator/domain/valueobjects"
)

// GoalRecommendationService は目標に関する推奨事項を提供するドメインサービス
type GoalRecommendationService struct {
	calculationService *FinancialCalculationService
}

// NewGoalRecommendationService は新しいGoalRecommendationServiceを作成する
func NewGoalRecommendationService(calculationService *FinancialCalculationService) *GoalRecommendationService {
	return &GoalRecommendationService{
		calculationService: calculationService,
	}
}

// RecommendationPriority は推奨事項の優先度を表す
type RecommendationPriority string

const (
	PriorityHigh   RecommendationPriority = "high"   // 高優先度
	PriorityMedium RecommendationPriority = "medium" // 中優先度
	PriorityLow    RecommendationPriority = "low"    // 低優先度
)

// GoalRecommendation は目標に対する推奨事項を表す
type GoalRecommendation struct {
	Type        string                 `json:"type"`        // "increase_savings", "extend_deadline", "reduce_target"
	Title       string                 `json:"title"`       // 推奨事項のタイトル
	Description string                 `json:"description"` // 詳細説明
	Priority    RecommendationPriority `json:"priority"`    // 優先度
	Impact      string                 `json:"impact"`      // 期待される効果
	NewValue    interface{}            `json:"new_value"`   // 推奨する新しい値
	Reason      string                 `json:"reason"`      // 推奨理由
}

// SavingsRecommendation は貯蓄に関する推奨事項を表す
type SavingsRecommendation struct {
	RecommendedAmount valueobjects.Money     `json:"recommended_amount"` // 推奨月間貯蓄額
	CurrentGap        valueobjects.Money     `json:"current_gap"`        // 現在の不足額
	Priority          RecommendationPriority `json:"priority"`           // 優先度
	Rationale         string                 `json:"rationale"`          // 根拠
	Achievability     string                 `json:"achievability"`      // 達成可能性の評価
}

// RecommendMonthlySavings は目標達成に必要な月間貯蓄額を推奨する
func (grs *GoalRecommendationService) RecommendMonthlySavings(
	goal *entities.Goal,
	currentSavings valueobjects.Money,
	timeRemaining valueobjects.Period,
) (*SavingsRecommendation, error) {
	if goal == nil {
		return nil, errors.New("目標は必須です")
	}

	// 残り必要金額を計算
	remainingAmount, err := goal.GetRemainingAmount()
	if err != nil {
		return nil, fmt.Errorf("残り必要金額の計算に失敗しました: %w", err)
	}

	// 既に目標達成している場合
	if remainingAmount.IsZero() || remainingAmount.IsNegative() {
		zeroAmount, _ := valueobjects.NewMoneyJPY(0)
		return &SavingsRecommendation{
			RecommendedAmount: zeroAmount,
			CurrentGap:        zeroAmount,
			Priority:          PriorityLow,
			Rationale:         "目標は既に達成されています",
			Achievability:     "達成済み",
		}, nil
	}

	// 残り期間を月数に変換
	remainingMonths := timeRemaining.ToMonths()
	if remainingMonths <= 0 {
		return &SavingsRecommendation{
			RecommendedAmount: remainingAmount,
			CurrentGap:        remainingAmount,
			Priority:          PriorityHigh,
			Rationale:         "目標期限が過ぎているため、即座の対応が必要です",
			Achievability:     "期限切れ",
		}, nil
	}

	// 推奨月間貯蓄額を計算
	recommendedMonthlySavings := remainingAmount.Amount() / float64(remainingMonths)
	recommendedAmount, err := valueobjects.NewMoneyJPY(recommendedMonthlySavings)
	if err != nil {
		return nil, fmt.Errorf("推奨月間貯蓄額の作成に失敗しました: %w", err)
	}

	// 現在の月間拠出額との差額を計算
	currentGap, err := recommendedAmount.Subtract(goal.MonthlyContribution())
	if err != nil {
		return nil, fmt.Errorf("現在の不足額の計算に失敗しました: %w", err)
	}

	// 優先度を決定
	priority := grs.determineSavingsPriority(goal, currentGap, remainingMonths)

	// 達成可能性を評価
	achievability := grs.evaluateAchievability(recommendedAmount, goal.GoalType())

	// 根拠を生成
	rationale := grs.generateSavingsRationale(goal, recommendedAmount, remainingMonths)

	return &SavingsRecommendation{
		RecommendedAmount: recommendedAmount,
		CurrentGap:        currentGap,
		Priority:          priority,
		Rationale:         rationale,
		Achievability:     achievability,
	}, nil
}

// SuggestGoalAdjustments は目標の調整案を提案する
func (grs *GoalRecommendationService) SuggestGoalAdjustments(
	goal *entities.Goal,
	financialProfile *entities.FinancialProfile,
) ([]GoalRecommendation, error) {
	if goal == nil {
		return nil, errors.New("目標は必須です")
	}

	if financialProfile == nil {
		return nil, errors.New("財務プロファイルは必須です")
	}

	var recommendations []GoalRecommendation

	// 目標の達成可能性をチェック
	achievable, err := goal.IsAchievable(financialProfile)
	if err != nil {
		return nil, fmt.Errorf("目標の達成可能性チェックに失敗しました: %w", err)
	}

	// 達成可能な場合は推奨事項なし
	if achievable {
		return recommendations, nil
	}

	// 純貯蓄額を取得
	netSavings, err := financialProfile.CalculateNetSavings()
	if err != nil {
		return nil, fmt.Errorf("純貯蓄額の計算に失敗しました: %w", err)
	}

	// 必要な月間貯蓄額を計算
	requiredMonthlySavings, err := goal.CalculateRequiredMonthlySavings()
	if err != nil {
		return nil, fmt.Errorf("必要月間貯蓄額の計算に失敗しました: %w", err)
	}

	// 1. 貯蓄額増加の推奨
	if netSavings.IsPositive() {
		savingsIncrease := grs.suggestSavingsIncrease(goal, netSavings, requiredMonthlySavings)
		if savingsIncrease != nil {
			recommendations = append(recommendations, *savingsIncrease)
		}
	}

	// 2. 期限延長の推奨
	deadlineExtension := grs.suggestDeadlineExtension(goal, netSavings)
	if deadlineExtension != nil {
		recommendations = append(recommendations, *deadlineExtension)
	}

	// 3. 目標金額削減の推奨
	targetReduction := grs.suggestTargetReduction(goal, netSavings)
	if targetReduction != nil {
		recommendations = append(recommendations, *targetReduction)
	}

	// 4. 支出削減の推奨
	expenseReduction := grs.suggestExpenseReduction(goal, financialProfile)
	if expenseReduction != nil {
		recommendations = append(recommendations, *expenseReduction)
	}

	// 5. 投資戦略の推奨
	investmentStrategy := grs.suggestInvestmentStrategy(goal, financialProfile)
	if investmentStrategy != nil {
		recommendations = append(recommendations, *investmentStrategy)
	}

	return recommendations, nil
}

// suggestSavingsIncrease は貯蓄額増加を推奨する
func (grs *GoalRecommendationService) suggestSavingsIncrease(
	goal *entities.Goal,
	netSavings valueobjects.Money,
	requiredMonthlySavings valueobjects.Money,
) *GoalRecommendation {
	// 現在の純貯蓄額で十分な場合はスキップ
	canAfford, err := netSavings.GreaterThan(requiredMonthlySavings)
	if err == nil && canAfford {
		return nil
	}

	// 必要な追加貯蓄額を計算
	additionalSavings, err := requiredMonthlySavings.Subtract(goal.MonthlyContribution())
	if err != nil {
		return nil
	}

	// 現在の収入に対する割合を計算
	// この情報は財務プロファイルから取得する必要があるが、ここでは簡略化

	return &GoalRecommendation{
		Type:        "increase_savings",
		Title:       "月間貯蓄額の増加",
		Description: fmt.Sprintf("目標達成のため、月間貯蓄額を%sに増加することを推奨します", requiredMonthlySavings.String()),
		Priority:    PriorityHigh,
		Impact:      "目標期日通りの達成が可能になります",
		NewValue:    requiredMonthlySavings.Amount(),
		Reason:      fmt.Sprintf("現在の貯蓄ペースでは目標達成に%s不足しています", additionalSavings.String()),
	}
}

// suggestDeadlineExtension は期限延長を推奨する
func (grs *GoalRecommendationService) suggestDeadlineExtension(
	goal *entities.Goal,
	netSavings valueobjects.Money,
) *GoalRecommendation {
	if netSavings.IsZero() || netSavings.IsNegative() {
		return nil
	}

	// 現在の純貯蓄額で目標達成に必要な期間を計算
	remainingAmount, err := goal.GetRemainingAmount()
	if err != nil {
		return nil
	}

	monthsNeeded := int(math.Ceil(remainingAmount.Amount() / netSavings.Amount()))
	newTargetDate := time.Now().AddDate(0, monthsNeeded, 0)

	// 現在の目標日と比較
	if newTargetDate.Before(goal.TargetDate()) {
		return nil // 既に十分な期間がある
	}

	extensionMonths := monthsNeeded - goal.GetRemainingDays()/30

	return &GoalRecommendation{
		Type:        "extend_deadline",
		Title:       "目標期日の延長",
		Description: fmt.Sprintf("目標期日を%dヶ月延長し、%sに設定することを推奨します", extensionMonths, newTargetDate.Format("2006年1月")),
		Priority:    PriorityMedium,
		Impact:      "現在の貯蓄ペースを維持しながら目標達成が可能になります",
		NewValue:    newTargetDate,
		Reason:      "現在の貯蓄能力に合わせた現実的な期日設定",
	}
}

// suggestTargetReduction は目標金額削減を推奨する
func (grs *GoalRecommendationService) suggestTargetReduction(
	goal *entities.Goal,
	netSavings valueobjects.Money,
) *GoalRecommendation {
	if netSavings.IsZero() || netSavings.IsNegative() {
		return nil
	}

	// 現在の期間で達成可能な金額を計算
	remainingDays := goal.GetRemainingDays()
	if remainingDays <= 0 {
		return nil
	}

	remainingMonths := remainingDays / 30
	achievableAmount := netSavings.Amount() * float64(remainingMonths)

	currentAmount := goal.CurrentAmount().Amount()
	newTargetAmount := currentAmount + achievableAmount

	// 現在の目標金額より低い場合のみ推奨
	if newTargetAmount >= goal.TargetAmount().Amount() {
		return nil
	}

	reductionAmount := goal.TargetAmount().Amount() - newTargetAmount

	newTarget, err := valueobjects.NewMoneyJPY(newTargetAmount)
	if err != nil {
		return nil
	}

	reductionMoney, err := valueobjects.NewMoneyJPY(reductionAmount)
	if err != nil {
		return nil
	}

	return &GoalRecommendation{
		Type:        "reduce_target",
		Title:       "目標金額の調整",
		Description: fmt.Sprintf("目標金額を%sに調整することを推奨します", newTarget.String()),
		Priority:    PriorityLow,
		Impact:      "現在の貯蓄能力で確実に達成可能な目標になります",
		NewValue:    newTargetAmount,
		Reason:      fmt.Sprintf("現在の貯蓄能力では%s過大な目標設定となっています", reductionMoney.String()),
	}
}

// suggestExpenseReduction は支出削減を推奨する
func (grs *GoalRecommendationService) suggestExpenseReduction(
	goal *entities.Goal,
	financialProfile *entities.FinancialProfile,
) *GoalRecommendation {
	// 必要な追加貯蓄額を計算
	requiredMonthlySavings, err := goal.CalculateRequiredMonthlySavings()
	if err != nil {
		return nil
	}

	netSavings, err := financialProfile.CalculateNetSavings()
	if err != nil {
		return nil
	}

	shortfall, err := requiredMonthlySavings.Subtract(netSavings)
	if err != nil || shortfall.IsNegative() {
		return nil
	}

	// 月収に対する支出削減の割合を計算
	monthlyIncome := financialProfile.MonthlyIncome()
	reductionPercentage := (shortfall.Amount() / monthlyIncome.Amount()) * 100

	return &GoalRecommendation{
		Type:        "reduce_expenses",
		Title:       "支出の見直し",
		Description: fmt.Sprintf("月間支出を%s（収入の%.1f%%）削減することを推奨します", shortfall.String(), reductionPercentage),
		Priority:    PriorityMedium,
		Impact:      "目標達成に必要な貯蓄額を確保できます",
		NewValue:    shortfall.Amount(),
		Reason:      "現在の収入では貯蓄額が不足しているため、支出の最適化が必要です",
	}
}

// suggestInvestmentStrategy は投資戦略を推奨する
func (grs *GoalRecommendationService) suggestInvestmentStrategy(
	goal *entities.Goal,
	financialProfile *entities.FinancialProfile,
) *GoalRecommendation {
	// 目標期間が短い場合（1年未満）は投資を推奨しない
	remainingDays := goal.GetRemainingDays()
	if remainingDays < 365 {
		return nil
	}

	// 現在の投資利回りが低い場合のみ推奨
	currentReturn := financialProfile.InvestmentReturn()
	if currentReturn.AsPercentage() >= 5.0 {
		return nil // 既に適切な利回り
	}

	// 目標タイプに応じた推奨利回りを設定
	var recommendedReturn float64
	var strategy string

	switch goal.GoalType() {
	case entities.GoalTypeRetirement:
		recommendedReturn = 6.0 // 長期投資向け
		strategy = "長期的な資産形成のため、株式中心のポートフォリオを検討してください"
	case entities.GoalTypeEmergency:
		recommendedReturn = 2.0 // 安全性重視
		strategy = "緊急資金は安全性を重視し、定期預金や国債での運用を検討してください"
	default:
		recommendedReturn = 4.0 // バランス型
		strategy = "バランス型の投資信託での運用を検討してください"
	}

	return &GoalRecommendation{
		Type:        "investment_strategy",
		Title:       "投資戦略の見直し",
		Description: fmt.Sprintf("投資利回りを%.1f%%に向上させることを推奨します", recommendedReturn),
		Priority:    PriorityMedium,
		Impact:      "複利効果により目標達成が容易になります",
		NewValue:    recommendedReturn,
		Reason:      strategy,
	}
}

// determineSavingsPriority は貯蓄推奨の優先度を決定する
func (grs *GoalRecommendationService) determineSavingsPriority(
	goal *entities.Goal,
	currentGap valueobjects.Money,
	remainingMonths int,
) RecommendationPriority {
	// 緊急資金目標は高優先度
	if goal.GoalType() == entities.GoalTypeEmergency {
		return PriorityHigh
	}

	// 期限が近い場合は高優先度
	if remainingMonths <= 6 {
		return PriorityHigh
	}

	// 不足額が大きい場合は高優先度
	if currentGap.IsPositive() && currentGap.Amount() > 50000 {
		return PriorityHigh
	}

	// 退職目標は中優先度
	if goal.GoalType() == entities.GoalTypeRetirement {
		return PriorityMedium
	}

	return PriorityLow
}

// evaluateAchievability は達成可能性を評価する
func (grs *GoalRecommendationService) evaluateAchievability(
	recommendedAmount valueobjects.Money,
	goalType entities.GoalType,
) string {
	// 簡略化された評価ロジック
	amount := recommendedAmount.Amount()

	switch {
	case amount <= 10000:
		return "容易に達成可能"
	case amount <= 50000:
		return "努力により達成可能"
	case amount <= 100000:
		return "計画的な取り組みが必要"
	default:
		return "大幅な生活スタイルの変更が必要"
	}
}

// generateSavingsRationale は貯蓄推奨の根拠を生成する
func (grs *GoalRecommendationService) generateSavingsRationale(
	goal *entities.Goal,
	recommendedAmount valueobjects.Money,
	remainingMonths int,
) string {
	goalTypeStr := goal.GoalType().String()

	return fmt.Sprintf(
		"%sの達成のため、残り%dヶ月で月額%sの貯蓄が必要です。",
		goalTypeStr,
		remainingMonths,
		recommendedAmount.String(),
	)
}

// AnalyzeGoalFeasibility は目標の実現可能性を分析する
func (grs *GoalRecommendationService) AnalyzeGoalFeasibility(
	goal *entities.Goal,
	financialProfile *entities.FinancialProfile,
) (map[string]interface{}, error) {
	if goal == nil || financialProfile == nil {
		return nil, errors.New("目標と財務プロファイルは必須です")
	}

	analysis := make(map[string]interface{})

	// 基本情報
	analysis["goal_type"] = goal.GoalType().String()
	analysis["target_amount"] = goal.TargetAmount().Amount()
	analysis["current_amount"] = goal.CurrentAmount().Amount()
	analysis["remaining_days"] = goal.GetRemainingDays()

	// 財務状況
	netSavings, err := financialProfile.CalculateNetSavings()
	if err != nil {
		return nil, fmt.Errorf("純貯蓄額の計算に失敗しました: %w", err)
	}
	analysis["net_savings"] = netSavings.Amount()

	// 必要貯蓄額
	requiredMonthlySavings, err := goal.CalculateRequiredMonthlySavings()
	if err != nil {
		return nil, fmt.Errorf("必要月間貯蓄額の計算に失敗しました: %w", err)
	}
	analysis["required_monthly_savings"] = requiredMonthlySavings.Amount()

	// 達成可能性
	achievable, err := goal.IsAchievable(financialProfile)
	if err != nil {
		return nil, fmt.Errorf("達成可能性の判定に失敗しました: %w", err)
	}
	analysis["achievable"] = achievable

	// 進捗率
	progress, err := goal.CalculateProgress(goal.CurrentAmount())
	if err != nil {
		return nil, fmt.Errorf("進捗率の計算に失敗しました: %w", err)
	}
	analysis["progress_percentage"] = progress.AsPercentage()

	// リスク評価
	analysis["risk_level"] = grs.assessRiskLevel(goal, financialProfile)

	return analysis, nil
}

// assessRiskLevel はリスクレベルを評価する
func (grs *GoalRecommendationService) assessRiskLevel(
	goal *entities.Goal,
	financialProfile *entities.FinancialProfile,
) string {
	// 簡略化されたリスク評価
	netSavings, err := financialProfile.CalculateNetSavings()
	if err != nil || netSavings.IsNegative() {
		return "高リスク"
	}

	requiredMonthlySavings, err := goal.CalculateRequiredMonthlySavings()
	if err != nil {
		return "評価不可"
	}

	ratio := requiredMonthlySavings.Amount() / netSavings.Amount()

	switch {
	case ratio <= 0.5:
		return "低リスク"
	case ratio <= 0.8:
		return "中リスク"
	default:
		return "高リスク"
	}
}
