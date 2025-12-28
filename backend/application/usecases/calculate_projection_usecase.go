package usecases

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/financial-planning-calculator/backend/domain/aggregates"
	"github.com/financial-planning-calculator/backend/domain/entities"
	"github.com/financial-planning-calculator/backend/domain/repositories"
	"github.com/financial-planning-calculator/backend/domain/services"
	"github.com/financial-planning-calculator/backend/infrastructure/log"
)

// CalculateProjectionUseCase は将来予測計算のユースケース
type CalculateProjectionUseCase interface {
	// CalculateAssetProjection は資産推移を計算する
	CalculateAssetProjection(ctx context.Context, input AssetProjectionInput) (*AssetProjectionOutput, error)

	// CalculateRetirementProjection は退職資金予測を計算する
	CalculateRetirementProjection(ctx context.Context, input RetirementProjectionInput) (*RetirementProjectionOutput, error)

	// CalculateEmergencyFundProjection は緊急資金予測を計算する
	CalculateEmergencyFundProjection(ctx context.Context, input EmergencyFundProjectionInput) (*EmergencyFundProjectionOutput, error)

	// CalculateComprehensiveProjection は包括的な財務予測を計算する
	CalculateComprehensiveProjection(ctx context.Context, input ComprehensiveProjectionInput) (*ComprehensiveProjectionOutput, error)

	// CalculateGoalProjection は目標達成予測を計算する
	CalculateGoalProjection(ctx context.Context, input GoalProjectionInput) (*GoalProjectionOutput, error)
}

// AssetProjectionInput は資産推移計算の入力
type AssetProjectionInput struct {
	UserID entities.UserID `json:"user_id"`
	Years  int             `json:"years"`
}

// AssetProjectionOutput は資産推移計算の出力
type AssetProjectionOutput struct {
	Projections []entities.AssetProjection `json:"projections"`
	Summary     ProjectionSummary          `json:"summary"`
}

// ProjectionSummary は予測サマリー
type ProjectionSummary struct {
	InitialAmount    float64 `json:"initial_amount"`
	FinalAmount      float64 `json:"final_amount"`
	TotalGrowth      float64 `json:"total_growth"`
	GrowthPercentage float64 `json:"growth_percentage"`
	AverageReturn    float64 `json:"average_return"`
}

// RetirementProjectionInput は退職資金予測計算の入力
type RetirementProjectionInput struct {
	UserID entities.UserID `json:"user_id"`
}

// RetirementProjectionOutput は退職資金予測計算の出力
type RetirementProjectionOutput struct {
	Calculation        *entities.RetirementCalculation `json:"calculation"`
	Recommendations    []string                        `json:"recommendations"`
	SufficiencyLevel   string                          `json:"sufficiency_level"`
	RequiredAdjustment *RequiredAdjustment             `json:"required_adjustment,omitempty"`
}

// RequiredAdjustment は必要な調整
type RequiredAdjustment struct {
	Type               string  `json:"type"` // "increase_savings", "extend_retirement", "reduce_expenses"
	Amount             float64 `json:"amount"`
	Description        string  `json:"description"`
	ImpactOnRetirement string  `json:"impact_on_retirement"`
}

// EmergencyFundProjectionInput は緊急資金予測計算の入力
type EmergencyFundProjectionInput struct {
	UserID entities.UserID `json:"user_id"`
}

// EmergencyFundProjectionOutput は緊急資金予測計算の出力
type EmergencyFundProjectionOutput struct {
	Status          *aggregates.EmergencyFundStatus `json:"status"`
	Recommendations []string                        `json:"recommendations"`
	Priority        string                          `json:"priority"`
	Timeline        *EmergencyFundTimeline          `json:"timeline"`
}

// EmergencyFundTimeline は緊急資金達成タイムライン
type EmergencyFundTimeline struct {
	MonthsToTarget     int         `json:"months_to_target"`
	MonthlySavingsGoal float64     `json:"monthly_savings_goal"`
	Milestones         []Milestone `json:"milestones"`
}

// Milestone はマイルストーン
type Milestone struct {
	Month       int     `json:"month"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
}

// ComprehensiveProjectionInput は包括的財務予測計算の入力
type ComprehensiveProjectionInput struct {
	UserID entities.UserID `json:"user_id"`
	Years  int             `json:"years"`
}

// ComprehensiveProjectionOutput は包括的財務予測計算の出力
type ComprehensiveProjectionOutput struct {
	PlanProjection *aggregates.PlanProjection `json:"plan_projection"`
	Insights       []FinancialInsight         `json:"insights"`
	Warnings       []FinancialWarning         `json:"warnings"`
	Opportunities  []FinancialOpportunity     `json:"opportunities"`
}

// FinancialInsight は財務洞察
type FinancialInsight struct {
	Type        string `json:"type"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Impact      string `json:"impact"`
}

// FinancialWarning は財務警告
type FinancialWarning struct {
	Type        string `json:"type"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Severity    string `json:"severity"` // "low", "medium", "high"
	Action      string `json:"action"`
}

// FinancialOpportunity は財務機会
type FinancialOpportunity struct {
	Type        string  `json:"type"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Benefit     string  `json:"benefit"`
	Effort      string  `json:"effort"` // "low", "medium", "high"
	Impact      float64 `json:"impact"`
}

// GoalProjectionInput は目標達成予測計算の入力
type GoalProjectionInput struct {
	UserID entities.UserID `json:"user_id"`
	GoalID entities.GoalID `json:"goal_id"`
}

// GoalProjectionOutput は目標達成予測計算の出力
type GoalProjectionOutput struct {
	Goal            *entities.Goal                `json:"goal"`
	Progress        entities.ProgressRate         `json:"progress"`
	Projection      []GoalProgressProjection      `json:"projection"`
	Recommendations []services.GoalRecommendation `json:"recommendations"`
	Feasibility     map[string]interface{}        `json:"feasibility"`
}

// GoalProgressProjection は目標進捗予測
type GoalProgressProjection struct {
	Month           int     `json:"month"`
	ProjectedAmount float64 `json:"projected_amount"`
	ProgressRate    float64 `json:"progress_rate"`
	OnTrack         bool    `json:"on_track"`
}

// calculateProjectionUseCaseImpl はCalculateProjectionUseCaseの実装
type calculateProjectionUseCaseImpl struct {
	financialPlanRepo     repositories.FinancialPlanRepository
	goalRepo              repositories.GoalRepository
	calculationService    *services.FinancialCalculationService
	recommendationService *services.GoalRecommendationService
	logger                *log.UseCaseLogger
}

// NewCalculateProjectionUseCase は新しいCalculateProjectionUseCaseを作成する
func NewCalculateProjectionUseCase(
	financialPlanRepo repositories.FinancialPlanRepository,
	goalRepo repositories.GoalRepository,
	calculationService *services.FinancialCalculationService,
	recommendationService *services.GoalRecommendationService,
) CalculateProjectionUseCase {
	return &calculateProjectionUseCaseImpl{
		financialPlanRepo:     financialPlanRepo,
		goalRepo:              goalRepo,
		calculationService:    calculationService,
		recommendationService: recommendationService,
		logger:                log.NewUseCaseLogger("CalculateProjectionUseCase"),
	}
}

// CalculateAssetProjection は資産推移を計算する
func (uc *calculateProjectionUseCaseImpl) CalculateAssetProjection(
	ctx context.Context,
	input AssetProjectionInput,
) (*AssetProjectionOutput, error) {
	ctx = uc.logger.StartOperation(ctx, "CalculateAssetProjection",
		slog.String("user_id", string(input.UserID)),
		slog.Int("years", input.Years),
	)

	// 財務計画を取得
	plan, err := uc.financialPlanRepo.FindByUserID(ctx, input.UserID)
	if err != nil {
		uc.logger.OperationError(ctx, "CalculateAssetProjection", err,
			slog.String("step", "find_plan"),
		)
		return nil, fmt.Errorf("財務計画の取得に失敗しました: %w", err)
	}

	// 資産推移を計算
	projections, err := plan.Profile().ProjectAssets(input.Years)
	if err != nil {
		uc.logger.OperationError(ctx, "CalculateAssetProjection", err,
			slog.String("step", "project_assets"),
		)
		return nil, fmt.Errorf("資産推移の計算に失敗しました: %w", err)
	}

	// サマリーを計算
	summary, err := uc.calculateProjectionSummary(projections)
	if err != nil {
		uc.logger.OperationError(ctx, "CalculateAssetProjection", err,
			slog.String("step", "calculate_summary"),
		)
		return nil, fmt.Errorf("予測サマリーの計算に失敗しました: %w", err)
	}

	uc.logger.EndOperation(ctx, "CalculateAssetProjection",
		slog.Int("projection_count", len(projections)),
	)

	return &AssetProjectionOutput{
		Projections: projections,
		Summary:     *summary,
	}, nil
}

// CalculateRetirementProjection は退職資金予測を計算する
func (uc *calculateProjectionUseCaseImpl) CalculateRetirementProjection(
	ctx context.Context,
	input RetirementProjectionInput,
) (*RetirementProjectionOutput, error) {
	ctx = uc.logger.StartOperation(ctx, "CalculateRetirementProjection",
		slog.String("user_id", string(input.UserID)),
	)

	// 財務計画を取得
	plan, err := uc.financialPlanRepo.FindByUserID(ctx, input.UserID)
	if err != nil {
		uc.logger.OperationError(ctx, "CalculateRetirementProjection", err,
			slog.String("step", "find_plan"),
		)
		return nil, fmt.Errorf("財務計画の取得に失敗しました: %w", err)
	}

	// 退職データが設定されているかチェック
	retirementData := plan.RetirementData()
	if retirementData == nil {
		err := fmt.Errorf("退職データが設定されていません")
		uc.logger.OperationError(ctx, "CalculateRetirementProjection", err,
			slog.String("step", "check_retirement_data"),
		)
		return nil, err
	}

	// 退職資金計算
	currentSavings, err := plan.Profile().CurrentSavings().Total()
	if err != nil {
		uc.logger.OperationError(ctx, "CalculateRetirementProjection", err,
			slog.String("step", "calculate_current_savings"),
		)
		return nil, fmt.Errorf("現在の貯蓄合計の計算に失敗しました: %w", err)
	}

	netSavings, err := plan.Profile().CalculateNetSavings()
	if err != nil {
		uc.logger.OperationError(ctx, "CalculateRetirementProjection", err,
			slog.String("step", "calculate_net_savings"),
		)
		return nil, fmt.Errorf("純貯蓄額の計算に失敗しました: %w", err)
	}

	calculation, err := retirementData.CalculateRetirementSufficiency(
		currentSavings,
		netSavings,
		plan.Profile().InvestmentReturn(),
		plan.Profile().InflationRate(),
	)
	if err != nil {
		uc.logger.OperationError(ctx, "CalculateRetirementProjection", err,
			slog.String("step", "calculate_sufficiency"),
		)
		return nil, fmt.Errorf("退職資金計算に失敗しました: %w", err)
	}

	// 推奨事項を生成
	recommendations := uc.generateRetirementRecommendations(calculation)

	// 充足レベルを評価
	sufficiencyLevel := uc.evaluateRetirementSufficiency(calculation)

	// 必要な調整を計算
	var requiredAdjustment *RequiredAdjustment
	if calculation.SufficiencyRate.AsPercentage() < 100 {
		requiredAdjustment = uc.calculateRequiredRetirementAdjustment(calculation, plan)
	}

	uc.logger.EndOperation(ctx, "CalculateRetirementProjection",
		slog.String("sufficiency_level", sufficiencyLevel),
	)

	return &RetirementProjectionOutput{
		Calculation:        calculation,
		Recommendations:    recommendations,
		SufficiencyLevel:   sufficiencyLevel,
		RequiredAdjustment: requiredAdjustment,
	}, nil
}

// CalculateEmergencyFundProjection は緊急資金予測を計算する
func (uc *calculateProjectionUseCaseImpl) CalculateEmergencyFundProjection(
	ctx context.Context,
	input EmergencyFundProjectionInput,
) (*EmergencyFundProjectionOutput, error) {
	ctx = uc.logger.StartOperation(ctx, "CalculateEmergencyFundProjection",
		slog.String("user_id", string(input.UserID)),
	)

	// 財務計画を取得
	plan, err := uc.financialPlanRepo.FindByUserID(ctx, input.UserID)
	if err != nil {
		uc.logger.OperationError(ctx, "CalculateEmergencyFundProjection", err,
			slog.String("step", "find_plan"),
		)
		return nil, fmt.Errorf("財務計画の取得に失敗しました: %w", err)
	}

	// 包括的予測を生成して緊急資金状況を取得
	projection, err := plan.GenerateProjection(1) // 1年間の予測
	if err != nil {
		uc.logger.OperationError(ctx, "CalculateEmergencyFundProjection", err,
			slog.String("step", "generate_projection"),
		)
		return nil, fmt.Errorf("財務予測の生成に失敗しました: %w", err)
	}

	if projection.EmergencyFundStatus == nil {
		err := fmt.Errorf("緊急資金状況が計算されていません")
		uc.logger.OperationError(ctx, "CalculateEmergencyFundProjection", err,
			slog.String("step", "check_emergency_status"),
		)
		return nil, err
	}

	// 推奨事項を生成
	recommendations := uc.generateEmergencyFundRecommendations(projection.EmergencyFundStatus, plan)

	// 優先度を評価
	priority := uc.evaluateEmergencyFundPriority(projection.EmergencyFundStatus)

	// タイムラインを計算
	timeline := uc.calculateEmergencyFundTimeline(projection.EmergencyFundStatus, plan)

	uc.logger.EndOperation(ctx, "CalculateEmergencyFundProjection",
		slog.String("priority", priority),
	)

	return &EmergencyFundProjectionOutput{
		Status:          projection.EmergencyFundStatus,
		Recommendations: recommendations,
		Priority:        priority,
		Timeline:        timeline,
	}, nil
}

// CalculateComprehensiveProjection は包括的な財務予測を計算する
func (uc *calculateProjectionUseCaseImpl) CalculateComprehensiveProjection(
	ctx context.Context,
	input ComprehensiveProjectionInput,
) (*ComprehensiveProjectionOutput, error) {
	ctx = uc.logger.StartOperation(ctx, "CalculateComprehensiveProjection",
		slog.String("user_id", string(input.UserID)),
		slog.Int("years", input.Years),
	)

	// 財務計画を取得
	plan, err := uc.financialPlanRepo.FindByUserID(ctx, input.UserID)
	if err != nil {
		uc.logger.OperationError(ctx, "CalculateComprehensiveProjection", err,
			slog.String("step", "find_plan"),
		)
		return nil, fmt.Errorf("財務計画の取得に失敗しました: %w", err)
	}

	// 包括的予測を生成
	projection, err := plan.GenerateProjection(input.Years)
	if err != nil {
		uc.logger.OperationError(ctx, "CalculateComprehensiveProjection", err,
			slog.String("step", "generate_projection"),
		)
		return nil, fmt.Errorf("包括的予測の生成に失敗しました: %w", err)
	}

	// 洞察を生成
	insights := uc.generateFinancialInsights(projection, plan)

	// 警告を生成
	warnings := uc.generateFinancialWarnings(projection, plan)

	// 機会を生成
	opportunities := uc.generateFinancialOpportunities(projection, plan)

	uc.logger.EndOperation(ctx, "CalculateComprehensiveProjection",
		slog.Int("insights_count", len(insights)),
		slog.Int("warnings_count", len(warnings)),
	)

	return &ComprehensiveProjectionOutput{
		PlanProjection: projection,
		Insights:       insights,
		Warnings:       warnings,
		Opportunities:  opportunities,
	}, nil
}

// CalculateGoalProjection は目標達成予測を計算する
func (uc *calculateProjectionUseCaseImpl) CalculateGoalProjection(
	ctx context.Context,
	input GoalProjectionInput,
) (*GoalProjectionOutput, error) {
	// 目標を取得
	goal, err := uc.goalRepo.FindByID(ctx, input.GoalID)
	if err != nil {
		return nil, fmt.Errorf("目標の取得に失敗しました: %w", err)
	}

	// 財務計画を取得
	plan, err := uc.financialPlanRepo.FindByUserID(ctx, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("財務計画の取得に失敗しました: %w", err)
	}

	// 進捗を計算
	progress, err := goal.CalculateProgress(goal.CurrentAmount())
	if err != nil {
		return nil, fmt.Errorf("目標進捗の計算に失敗しました: %w", err)
	}

	// 進捗予測を計算
	projection := uc.calculateGoalProgressProjection(goal, plan.Profile())

	// 推奨事項を生成
	recommendations, err := uc.recommendationService.SuggestGoalAdjustments(goal, plan.Profile())
	if err != nil {
		return nil, fmt.Errorf("推奨事項の生成に失敗しました: %w", err)
	}

	// 実現可能性を分析
	feasibility, err := uc.recommendationService.AnalyzeGoalFeasibility(goal, plan.Profile())
	if err != nil {
		return nil, fmt.Errorf("実現可能性の分析に失敗しました: %w", err)
	}

	return &GoalProjectionOutput{
		Goal:            goal,
		Progress:        progress,
		Projection:      projection,
		Recommendations: recommendations,
		Feasibility:     feasibility,
	}, nil
}

// calculateProjectionSummary は予測サマリーを計算する
func (uc *calculateProjectionUseCaseImpl) calculateProjectionSummary(projections []entities.AssetProjection) (*ProjectionSummary, error) {
	if len(projections) == 0 {
		return &ProjectionSummary{}, nil
	}

	initialAmount := projections[0].TotalAssets.Amount()
	finalAmount := projections[len(projections)-1].TotalAssets.Amount()
	totalGrowth := finalAmount - initialAmount
	growthPercentage := (totalGrowth / initialAmount) * 100
	averageReturn := growthPercentage / float64(len(projections))

	return &ProjectionSummary{
		InitialAmount:    initialAmount,
		FinalAmount:      finalAmount,
		TotalGrowth:      totalGrowth,
		GrowthPercentage: growthPercentage,
		AverageReturn:    averageReturn,
	}, nil
}

// generateRetirementRecommendations は退職資金の推奨事項を生成する
func (uc *calculateProjectionUseCaseImpl) generateRetirementRecommendations(calculation *entities.RetirementCalculation) []string {
	var recommendations []string

	sufficiencyRate := calculation.SufficiencyRate.AsPercentage()

	switch {
	case sufficiencyRate >= 100:
		recommendations = append(recommendations, "退職資金は十分に確保されています")
		recommendations = append(recommendations, "余剰資金を他の目標に振り分けることを検討してください")
	case sufficiencyRate >= 80:
		recommendations = append(recommendations, "退職資金はほぼ十分ですが、さらなる貯蓄を推奨します")
		recommendations = append(recommendations, "月間貯蓄額を少し増やすことを検討してください")
	case sufficiencyRate >= 60:
		recommendations = append(recommendations, "退職資金が不足しています。貯蓄額の増加が必要です")
		recommendations = append(recommendations, "支出の見直しや副収入の検討をお勧めします")
	default:
		recommendations = append(recommendations, "退職資金が大幅に不足しています。緊急の対策が必要です")
		recommendations = append(recommendations, "退職年齢の延長や生活費の大幅な見直しを検討してください")
	}

	return recommendations
}

// evaluateRetirementSufficiency は退職資金の充足レベルを評価する
func (uc *calculateProjectionUseCaseImpl) evaluateRetirementSufficiency(calculation *entities.RetirementCalculation) string {
	sufficiencyRate := calculation.SufficiencyRate.AsPercentage()

	switch {
	case sufficiencyRate >= 120:
		return "十分以上"
	case sufficiencyRate >= 100:
		return "十分"
	case sufficiencyRate >= 80:
		return "ほぼ十分"
	case sufficiencyRate >= 60:
		return "不足"
	default:
		return "大幅不足"
	}
}

// calculateRequiredRetirementAdjustment は必要な退職資金調整を計算する
func (uc *calculateProjectionUseCaseImpl) calculateRequiredRetirementAdjustment(calculation *entities.RetirementCalculation, plan *aggregates.FinancialPlan) *RequiredAdjustment {
	shortfall := calculation.Shortfall.Amount()

	// 月間貯蓄増加による調整
	monthsToRetirement := 12 * 30 // 仮定：30年
	requiredMonthlySavingsIncrease := shortfall / float64(monthsToRetirement)

	return &RequiredAdjustment{
		Type:               "increase_savings",
		Amount:             requiredMonthlySavingsIncrease,
		Description:        fmt.Sprintf("月間貯蓄額を%.0f円増加させる必要があります", requiredMonthlySavingsIncrease),
		ImpactOnRetirement: "目標通りの退職が可能になります",
	}
}

// generateEmergencyFundRecommendations は緊急資金の推奨事項を生成する
func (uc *calculateProjectionUseCaseImpl) generateEmergencyFundRecommendations(status *aggregates.EmergencyFundStatus, plan *aggregates.FinancialPlan) []string {
	var recommendations []string

	if status.Shortfall.IsZero() || status.Shortfall.IsNegative() {
		recommendations = append(recommendations, "緊急資金は十分に確保されています")
		return recommendations
	}

	shortfallRatio := status.Shortfall.Amount() / status.RequiredAmount.Amount()

	switch {
	case shortfallRatio > 0.8:
		recommendations = append(recommendations, "緊急資金が大幅に不足しています。最優先で確保してください")
		recommendations = append(recommendations, "他の投資を一時停止して緊急資金の確保を優先することを検討してください")
	case shortfallRatio > 0.5:
		recommendations = append(recommendations, "緊急資金が不足しています。計画的な積立が必要です")
		recommendations = append(recommendations, "月間支出を見直して緊急資金への拠出を増やしてください")
	default:
		recommendations = append(recommendations, "緊急資金をもう少し増やすことを推奨します")
		recommendations = append(recommendations, "安全性の高い預金商品での積立を検討してください")
	}

	return recommendations
}

// evaluateEmergencyFundPriority は緊急資金の優先度を評価する
func (uc *calculateProjectionUseCaseImpl) evaluateEmergencyFundPriority(status *aggregates.EmergencyFundStatus) string {
	if status.Shortfall.IsZero() || status.Shortfall.IsNegative() {
		return "低"
	}

	shortfallRatio := status.Shortfall.Amount() / status.RequiredAmount.Amount()

	switch {
	case shortfallRatio > 0.8:
		return "最高"
	case shortfallRatio > 0.5:
		return "高"
	case shortfallRatio > 0.2:
		return "中"
	default:
		return "低"
	}
}

// calculateEmergencyFundTimeline は緊急資金のタイムラインを計算する
func (uc *calculateProjectionUseCaseImpl) calculateEmergencyFundTimeline(status *aggregates.EmergencyFundStatus, plan *aggregates.FinancialPlan) *EmergencyFundTimeline {
	if status.MonthsToTarget <= 0 {
		return &EmergencyFundTimeline{
			MonthsToTarget:     0,
			MonthlySavingsGoal: 0,
			Milestones:         []Milestone{},
		}
	}

	monthlySavingsGoal := status.Shortfall.Amount() / float64(status.MonthsToTarget)

	var milestones []Milestone
	quarterlyAmount := status.RequiredAmount.Amount() / 4

	for i := 1; i <= 4; i++ {
		month := status.MonthsToTarget / 4 * i
		amount := quarterlyAmount * float64(i)
		description := fmt.Sprintf("緊急資金の%d%%達成", 25*i)

		milestones = append(milestones, Milestone{
			Month:       month,
			Amount:      amount,
			Description: description,
		})
	}

	return &EmergencyFundTimeline{
		MonthsToTarget:     status.MonthsToTarget,
		MonthlySavingsGoal: monthlySavingsGoal,
		Milestones:         milestones,
	}
}

// generateFinancialInsights は財務洞察を生成する
func (uc *calculateProjectionUseCaseImpl) generateFinancialInsights(projection *aggregates.PlanProjection, plan *aggregates.FinancialPlan) []FinancialInsight {
	var insights []FinancialInsight

	// 複利効果の洞察
	if len(projection.AssetProjections) > 0 {
		finalProjection := projection.AssetProjections[len(projection.AssetProjections)-1]
		compoundEffect := finalProjection.InvestmentGains.Amount() / finalProjection.ContributedAmount.Amount() * 100

		if compoundEffect > 50 {
			insights = append(insights, FinancialInsight{
				Type:        "compound_interest",
				Title:       "複利効果が大きく働いています",
				Description: fmt.Sprintf("投資収益が元本の%.1f%%に達し、複利効果が顕著に現れています", compoundEffect),
				Impact:      "長期的な資産形成に大きく貢献します",
			})
		}
	}

	// 貯蓄率の洞察
	netSavings, err := plan.Profile().CalculateNetSavings()
	if err == nil {
		monthlyIncome := plan.Profile().MonthlyIncome()
		savingsRate := netSavings.Amount() / monthlyIncome.Amount() * 100

		if savingsRate > 20 {
			insights = append(insights, FinancialInsight{
				Type:        "savings_rate",
				Title:       "優秀な貯蓄率を維持しています",
				Description: fmt.Sprintf("貯蓄率%.1f%%は理想的な水準です", savingsRate),
				Impact:      "財務目標の早期達成が期待できます",
			})
		}
	}

	return insights
}

// generateFinancialWarnings は財務警告を生成する
func (uc *calculateProjectionUseCaseImpl) generateFinancialWarnings(projection *aggregates.PlanProjection, plan *aggregates.FinancialPlan) []FinancialWarning {
	var warnings []FinancialWarning

	// 緊急資金の警告
	if projection.EmergencyFundStatus != nil && projection.EmergencyFundStatus.Shortfall.IsPositive() {
		shortfallRatio := projection.EmergencyFundStatus.Shortfall.Amount() / projection.EmergencyFundStatus.RequiredAmount.Amount()

		if shortfallRatio > 0.5 {
			warnings = append(warnings, FinancialWarning{
				Type:        "emergency_fund",
				Title:       "緊急資金が不足しています",
				Description: "緊急時に対応できる資金が不足している可能性があります",
				Severity:    "high",
				Action:      "緊急資金の確保を最優先で進めてください",
			})
		}
	}

	// 退職資金の警告
	if projection.RetirementCalculation != nil {
		sufficiencyRate := projection.RetirementCalculation.SufficiencyRate.AsPercentage()

		if sufficiencyRate < 80 {
			severity := "medium"
			if sufficiencyRate < 60 {
				severity = "high"
			}

			warnings = append(warnings, FinancialWarning{
				Type:        "retirement_fund",
				Title:       "退職資金が不足する可能性があります",
				Description: fmt.Sprintf("現在のペースでは退職資金が%.1f%%しか確保できません", sufficiencyRate),
				Severity:    severity,
				Action:      "貯蓄額の増加または退職計画の見直しを検討してください",
			})
		}
	}

	return warnings
}

// generateFinancialOpportunities は財務機会を生成する
func (uc *calculateProjectionUseCaseImpl) generateFinancialOpportunities(projection *aggregates.PlanProjection, plan *aggregates.FinancialPlan) []FinancialOpportunity {
	var opportunities []FinancialOpportunity

	// 投資利回り改善の機会
	currentReturn := plan.Profile().InvestmentReturn().AsPercentage()
	if currentReturn < 5 {
		currentSavingsTotal, err := plan.Profile().CurrentSavings().Total()
		if err == nil {
			potentialGain := (5 - currentReturn) / 100 * currentSavingsTotal.Amount()

			opportunities = append(opportunities, FinancialOpportunity{
				Type:        "investment_optimization",
				Title:       "投資利回りの改善機会",
				Description: "投資ポートフォリオの見直しにより利回り向上が期待できます",
				Benefit:     fmt.Sprintf("年間約%.0f円の追加収益が見込めます", potentialGain),
				Effort:      "medium",
				Impact:      potentialGain,
			})
		}
	}

	// 支出最適化の機会
	monthlyExpenses, err := plan.Profile().MonthlyExpenses().Total()
	if err == nil {
		monthlyIncome := plan.Profile().MonthlyIncome()
		expenseRatio := monthlyExpenses.Amount() / monthlyIncome.Amount()

		if expenseRatio > 0.7 {
			potentialSavings := monthlyExpenses.Amount() * 0.1 * 12 // 10%削減を1年間

			opportunities = append(opportunities, FinancialOpportunity{
				Type:        "expense_optimization",
				Title:       "支出最適化の機会",
				Description: "支出の見直しにより貯蓄額を増やすことができます",
				Benefit:     fmt.Sprintf("年間約%.0f円の追加貯蓄が可能です", potentialSavings),
				Effort:      "low",
				Impact:      potentialSavings,
			})
		}
	}

	return opportunities
}

// calculateGoalProgressProjection は目標進捗予測を計算する
func (uc *calculateProjectionUseCaseImpl) calculateGoalProgressProjection(goal *entities.Goal, profile *entities.FinancialProfile) []GoalProgressProjection {
	var projection []GoalProgressProjection

	remainingDays := goal.GetRemainingDays()
	if remainingDays <= 0 {
		return projection
	}

	remainingMonths := remainingDays / 30
	if remainingMonths <= 0 {
		remainingMonths = 1
	}

	currentAmount := goal.CurrentAmount().Amount()
	monthlyContribution := goal.MonthlyContribution().Amount()
	targetAmount := goal.TargetAmount().Amount()

	for month := 1; month <= remainingMonths; month++ {
		projectedAmount := currentAmount + (monthlyContribution * float64(month))
		progressRate := (projectedAmount / targetAmount) * 100
		onTrack := progressRate >= (float64(month)/float64(remainingMonths))*100

		projection = append(projection, GoalProgressProjection{
			Month:           month,
			ProjectedAmount: projectedAmount,
			ProgressRate:    progressRate,
			OnTrack:         onTrack,
		})
	}

	return projection
}
