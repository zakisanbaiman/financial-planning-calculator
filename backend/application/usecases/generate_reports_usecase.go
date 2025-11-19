package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/financial-planning-calculator/backend/domain/aggregates"
	"github.com/financial-planning-calculator/backend/domain/entities"
	"github.com/financial-planning-calculator/backend/domain/repositories"
	"github.com/financial-planning-calculator/backend/domain/services"
)

// GenerateReportsUseCase はレポート生成のユースケース
type GenerateReportsUseCase interface {
	// GenerateFinancialSummaryReport は財務サマリーレポートを生成する
	GenerateFinancialSummaryReport(ctx context.Context, input FinancialSummaryReportInput) (*FinancialSummaryReportOutput, error)

	// GenerateAssetProjectionReport は資産推移レポートを生成する
	GenerateAssetProjectionReport(ctx context.Context, input AssetProjectionReportInput) (*AssetProjectionReportOutput, error)

	// GenerateGoalsProgressReport は目標進捗レポートを生成する
	GenerateGoalsProgressReport(ctx context.Context, input GoalsProgressReportInput) (*GoalsProgressReportOutput, error)

	// GenerateRetirementPlanReport は退職計画レポートを生成する
	GenerateRetirementPlanReport(ctx context.Context, input RetirementPlanReportInput) (*RetirementPlanReportOutput, error)

	// GenerateComprehensiveReport は包括的レポートを生成する
	GenerateComprehensiveReport(ctx context.Context, input ComprehensiveReportInput) (*ComprehensiveReportOutput, error)

	// ExportReportToPDF はレポートをPDF形式でエクスポートする
	ExportReportToPDF(ctx context.Context, input ExportReportInput) (*ExportReportOutput, error)
}

// FinancialSummaryReportInput は財務サマリーレポート生成の入力
type FinancialSummaryReportInput struct {
	UserID entities.UserID `json:"user_id"`
}

// FinancialSummaryReportOutput は財務サマリーレポート生成の出力
type FinancialSummaryReportOutput struct {
	Report      FinancialSummaryReport `json:"report"`
	GeneratedAt string                 `json:"generated_at"`
}

// FinancialSummaryReport は財務サマリーレポート
type FinancialSummaryReport struct {
	UserID           entities.UserID  `json:"user_id"`
	ReportDate       string           `json:"report_date"`
	FinancialHealth  FinancialHealth  `json:"financial_health"`
	CurrentSituation CurrentSituation `json:"current_situation"`
	KeyMetrics       []KeyMetric      `json:"key_metrics"`
	Recommendations  []string         `json:"recommendations"`
	Warnings         []string         `json:"warnings"`
}

// FinancialHealth は財務健全性
type FinancialHealth struct {
	OverallScore       int     `json:"overall_score"`        // 0-100
	ScoreLevel         string  `json:"score_level"`          // "excellent", "good", "fair", "poor"
	SavingsRate        float64 `json:"savings_rate"`         // %
	DebtToIncomeRatio  float64 `json:"debt_to_income_ratio"` // %
	EmergencyFundRatio float64 `json:"emergency_fund_ratio"` // months
}

// CurrentSituation は現在の状況
type CurrentSituation struct {
	MonthlyIncome    float64 `json:"monthly_income"`
	MonthlyExpenses  float64 `json:"monthly_expenses"`
	NetSavings       float64 `json:"net_savings"`
	TotalAssets      float64 `json:"total_assets"`
	InvestmentReturn float64 `json:"investment_return"`
	InflationRate    float64 `json:"inflation_rate"`
}

// KeyMetric は主要指標
type KeyMetric struct {
	Name        string  `json:"name"`
	Value       float64 `json:"value"`
	Unit        string  `json:"unit"`
	Description string  `json:"description"`
	Trend       string  `json:"trend"` // "up", "down", "stable"
}

// AssetProjectionReportInput は資産推移レポート生成の入力
type AssetProjectionReportInput struct {
	UserID entities.UserID `json:"user_id"`
	Years  int             `json:"years"`
}

// AssetProjectionReportOutput は資産推移レポート生成の出力
type AssetProjectionReportOutput struct {
	Report      AssetProjectionReport `json:"report"`
	GeneratedAt string                `json:"generated_at"`
}

// AssetProjectionReport は資産推移レポート
type AssetProjectionReport struct {
	UserID          entities.UserID            `json:"user_id"`
	ProjectionYears int                        `json:"projection_years"`
	Projections     []entities.AssetProjection `json:"projections"`
	Summary         ProjectionSummary          `json:"summary"`
	Scenarios       []ScenarioAnalysis         `json:"scenarios"`
	Insights        []string                   `json:"insights"`
}

// ScenarioAnalysis はシナリオ分析
type ScenarioAnalysis struct {
	Name             string  `json:"name"`
	Description      string  `json:"description"`
	InvestmentReturn float64 `json:"investment_return"`
	InflationRate    float64 `json:"inflation_rate"`
	FinalAmount      float64 `json:"final_amount"`
	RealValue        float64 `json:"real_value"`
	Impact           string  `json:"impact"`
}

// GoalsProgressReportInput は目標進捗レポート生成の入力
type GoalsProgressReportInput struct {
	UserID entities.UserID `json:"user_id"`
}

// GoalsProgressReportOutput は目標進捗レポート生成の出力
type GoalsProgressReportOutput struct {
	Report      GoalsProgressReport `json:"report"`
	GeneratedAt string              `json:"generated_at"`
}

// GoalsProgressReport は目標進捗レポート
type GoalsProgressReport struct {
	UserID       entities.UserID `json:"user_id"`
	Goals        []GoalProgress  `json:"goals"`
	Summary      GoalsSummary    `json:"summary"`
	Achievements []Achievement   `json:"achievements"`
	NextSteps    []string        `json:"next_steps"`
}

// GoalProgress は目標進捗
type GoalProgress struct {
	Goal            *entities.Goal        `json:"goal"`
	Progress        entities.ProgressRate `json:"progress"`
	Status          string                `json:"status"`
	DaysRemaining   int                   `json:"days_remaining"`
	OnTrack         bool                  `json:"on_track"`
	Recommendations []string              `json:"recommendations"`
}

// Achievement は達成事項
type Achievement struct {
	Type        string `json:"type"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"`
	Impact      string `json:"impact"`
}

// RetirementPlanReportInput は退職計画レポート生成の入力
type RetirementPlanReportInput struct {
	UserID entities.UserID `json:"user_id"`
}

// RetirementPlanReportOutput は退職計画レポート生成の出力
type RetirementPlanReportOutput struct {
	Report      RetirementPlanReport `json:"report"`
	GeneratedAt string               `json:"generated_at"`
}

// RetirementPlanReport は退職計画レポート
type RetirementPlanReport struct {
	UserID          entities.UserID                 `json:"user_id"`
	RetirementData  *entities.RetirementData        `json:"retirement_data"`
	Calculation     *entities.RetirementCalculation `json:"calculation"`
	Projections     []RetirementProjection          `json:"projections"`
	Strategies      []RetirementStrategy            `json:"strategies"`
	Recommendations []string                        `json:"recommendations"`
	RiskAssessment  RiskAssessment                  `json:"risk_assessment"`
}

// RetirementProjection は退職予測
type RetirementProjection struct {
	Age               int     `json:"age"`
	YearsToRetirement int     `json:"years_to_retirement"`
	ProjectedAssets   float64 `json:"projected_assets"`
	RequiredAssets    float64 `json:"required_assets"`
	SufficiencyRate   float64 `json:"sufficiency_rate"`
	MonthlyShortfall  float64 `json:"monthly_shortfall"`
}

// RetirementStrategy は退職戦略
type RetirementStrategy struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Impact      float64 `json:"impact"`
	Effort      string  `json:"effort"` // "low", "medium", "high"
	Timeline    string  `json:"timeline"`
}

// RiskAssessment はリスク評価
type RiskAssessment struct {
	OverallRisk string       `json:"overall_risk"` // "low", "medium", "high"
	RiskFactors []RiskFactor `json:"risk_factors"`
	Mitigations []string     `json:"mitigations"`
}

// RiskFactor はリスク要因
type RiskFactor struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Impact      string `json:"impact"`      // "low", "medium", "high"
	Probability string `json:"probability"` // "low", "medium", "high"
}

// ComprehensiveReportInput は包括的レポート生成の入力
type ComprehensiveReportInput struct {
	UserID entities.UserID `json:"user_id"`
	Years  int             `json:"years"`
}

// ComprehensiveReportOutput は包括的レポート生成の出力
type ComprehensiveReportOutput struct {
	Report      ComprehensiveReport `json:"report"`
	GeneratedAt string              `json:"generated_at"`
}

// ComprehensiveReport は包括的レポート
type ComprehensiveReport struct {
	UserID           entities.UserID        `json:"user_id"`
	ExecutiveSummary ExecutiveSummary       `json:"executive_summary"`
	FinancialSummary FinancialSummaryReport `json:"financial_summary"`
	AssetProjection  AssetProjectionReport  `json:"asset_projection"`
	GoalsProgress    GoalsProgressReport    `json:"goals_progress"`
	RetirementPlan   *RetirementPlanReport  `json:"retirement_plan,omitempty"`
	ActionPlan       ActionPlan             `json:"action_plan"`
}

// ExecutiveSummary はエグゼクティブサマリー
type ExecutiveSummary struct {
	OverallStatus        string   `json:"overall_status"`
	KeyHighlights        []string `json:"key_highlights"`
	CriticalActions      []string `json:"critical_actions"`
	OpportunityAreas     []string `json:"opportunity_areas"`
	FinancialHealthScore int      `json:"financial_health_score"`
}

// ActionPlan はアクションプラン
type ActionPlan struct {
	ShortTerm  []ActionItem `json:"short_term"`  // 3ヶ月以内
	MediumTerm []ActionItem `json:"medium_term"` // 1年以内
	LongTerm   []ActionItem `json:"long_term"`   // 1年以上
}

// ActionItem はアクション項目
type ActionItem struct {
	Priority    string `json:"priority"` // "high", "medium", "low"
	Title       string `json:"title"`
	Description string `json:"description"`
	Timeline    string `json:"timeline"`
	Impact      string `json:"impact"`
	Effort      string `json:"effort"`
}

// ExportReportInput はレポートエクスポートの入力
type ExportReportInput struct {
	UserID     entities.UserID `json:"user_id"`
	ReportType string          `json:"report_type"` // "financial_summary", "comprehensive", etc.
	Format     string          `json:"format"`      // "pdf", "excel", "csv"
	ReportData interface{}     `json:"report_data"`
}

// ExportReportOutput はレポートエクスポートの出力
type ExportReportOutput struct {
	FileName    string `json:"file_name"`
	FileSize    int64  `json:"file_size"`
	DownloadURL string `json:"download_url"`
	ExpiresAt   string `json:"expires_at"`
}

// generateReportsUseCaseImpl はGenerateReportsUseCaseの実装
type generateReportsUseCaseImpl struct {
	financialPlanRepo     repositories.FinancialPlanRepository
	goalRepo              repositories.GoalRepository
	calculationService    *services.FinancialCalculationService
	recommendationService *services.GoalRecommendationService
}

// NewGenerateReportsUseCase は新しいGenerateReportsUseCaseを作成する
func NewGenerateReportsUseCase(
	financialPlanRepo repositories.FinancialPlanRepository,
	goalRepo repositories.GoalRepository,
	calculationService *services.FinancialCalculationService,
	recommendationService *services.GoalRecommendationService,
) GenerateReportsUseCase {
	return &generateReportsUseCaseImpl{
		financialPlanRepo:     financialPlanRepo,
		goalRepo:              goalRepo,
		calculationService:    calculationService,
		recommendationService: recommendationService,
	}
}

// GenerateFinancialSummaryReport は財務サマリーレポートを生成する
func (uc *generateReportsUseCaseImpl) GenerateFinancialSummaryReport(
	ctx context.Context,
	input FinancialSummaryReportInput,
) (*FinancialSummaryReportOutput, error) {
	// 財務計画を取得
	plan, err := uc.financialPlanRepo.FindByUserID(ctx, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("財務計画の取得に失敗しました: %w", err)
	}

	// 財務健全性を計算
	financialHealth, err := uc.calculateFinancialHealth(plan)
	if err != nil {
		return nil, fmt.Errorf("財務健全性の計算に失敗しました: %w", err)
	}

	// 現在の状況を取得
	currentSituation, err := uc.getCurrentSituation(plan)
	if err != nil {
		return nil, fmt.Errorf("現在の状況の取得に失敗しました: %w", err)
	}

	// 主要指標を計算
	keyMetrics, err := uc.calculateKeyMetrics(plan)
	if err != nil {
		return nil, fmt.Errorf("主要指標の計算に失敗しました: %w", err)
	}

	// 推奨事項と警告を生成
	recommendations, warnings := uc.generateRecommendationsAndWarnings(plan)

	report := FinancialSummaryReport{
		UserID:           input.UserID,
		ReportDate:       time.Now().Format("2006-01-02"),
		FinancialHealth:  *financialHealth,
		CurrentSituation: *currentSituation,
		KeyMetrics:       keyMetrics,
		Recommendations:  recommendations,
		Warnings:         warnings,
	}

	return &FinancialSummaryReportOutput{
		Report:      report,
		GeneratedAt: time.Now().Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// GenerateAssetProjectionReport は資産推移レポートを生成する
func (uc *generateReportsUseCaseImpl) GenerateAssetProjectionReport(
	ctx context.Context,
	input AssetProjectionReportInput,
) (*AssetProjectionReportOutput, error) {
	// 財務計画を取得
	plan, err := uc.financialPlanRepo.FindByUserID(ctx, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("財務計画の取得に失敗しました: %w", err)
	}

	// 資産推移を計算
	projections, err := plan.Profile().ProjectAssets(input.Years)
	if err != nil {
		return nil, fmt.Errorf("資産推移の計算に失敗しました: %w", err)
	}

	// サマリーを計算
	summary, err := uc.calculateProjectionSummary(projections)
	if err != nil {
		return nil, fmt.Errorf("予測サマリーの計算に失敗しました: %w", err)
	}

	// シナリオ分析を実行
	scenarios := uc.generateScenarioAnalysis(plan, input.Years)

	// 洞察を生成
	insights := uc.generateProjectionInsights(projections, scenarios)

	report := AssetProjectionReport{
		UserID:          input.UserID,
		ProjectionYears: input.Years,
		Projections:     projections,
		Summary:         *summary,
		Scenarios:       scenarios,
		Insights:        insights,
	}

	return &AssetProjectionReportOutput{
		Report:      report,
		GeneratedAt: time.Now().Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// GenerateGoalsProgressReport は目標進捗レポートを生成する
func (uc *generateReportsUseCaseImpl) GenerateGoalsProgressReport(
	ctx context.Context,
	input GoalsProgressReportInput,
) (*GoalsProgressReportOutput, error) {
	// 目標一覧を取得
	goals, err := uc.goalRepo.FindByUserID(ctx, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("目標の取得に失敗しました: %w", err)
	}

	// 財務計画を取得
	plan, err := uc.financialPlanRepo.FindByUserID(ctx, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("財務計画の取得に失敗しました: %w", err)
	}

	// 目標進捗を計算
	var goalProgresses []GoalProgress
	var summary GoalsSummary

	for _, goal := range goals {
		progress, err := goal.CalculateProgress(goal.CurrentAmount())
		if err != nil {
			return nil, fmt.Errorf("目標進捗の計算に失敗しました: %w", err)
		}

		// 推奨事項を生成
		recommendations, err := uc.recommendationService.SuggestGoalAdjustments(goal, plan.Profile())
		if err != nil {
			return nil, fmt.Errorf("推奨事項の生成に失敗しました: %w", err)
		}

		var recommendationTexts []string
		for _, rec := range recommendations {
			recommendationTexts = append(recommendationTexts, rec.Description)
		}

		status := uc.getGoalStatusText(goal)
		onTrack, _ := goal.IsAchievable(plan.Profile())

		goalProgresses = append(goalProgresses, GoalProgress{
			Goal:            goal,
			Progress:        progress,
			Status:          status,
			DaysRemaining:   goal.GetRemainingDays(),
			OnTrack:         onTrack,
			Recommendations: recommendationTexts,
		})

		// サマリーを更新
		summary.TotalGoals++
		summary.TotalTarget += goal.TargetAmount().Amount()
		summary.TotalCurrent += goal.CurrentAmount().Amount()

		if goal.IsActive() {
			summary.ActiveGoals++
		}
		if goal.IsCompleted() {
			summary.CompletedGoals++
		}
		if goal.IsOverdue() {
			summary.OverdueGoals++
		}
	}

	// 全体進捗を計算
	if summary.TotalTarget > 0 {
		summary.OverallProgress = (summary.TotalCurrent / summary.TotalTarget) * 100
	}

	// 達成事項を生成
	achievements := uc.generateAchievements(goals)

	// 次のステップを生成
	nextSteps := uc.generateNextSteps(goalProgresses)

	report := GoalsProgressReport{
		UserID:       input.UserID,
		Goals:        goalProgresses,
		Summary:      summary,
		Achievements: achievements,
		NextSteps:    nextSteps,
	}

	return &GoalsProgressReportOutput{
		Report:      report,
		GeneratedAt: time.Now().Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// GenerateRetirementPlanReport は退職計画レポートを生成する
func (uc *generateReportsUseCaseImpl) GenerateRetirementPlanReport(
	ctx context.Context,
	input RetirementPlanReportInput,
) (*RetirementPlanReportOutput, error) {
	// 財務計画を取得
	plan, err := uc.financialPlanRepo.FindByUserID(ctx, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("財務計画の取得に失敗しました: %w", err)
	}

	// 退職データが設定されているかチェック
	retirementData := plan.RetirementData()
	if retirementData == nil {
		return nil, fmt.Errorf("退職データが設定されていません")
	}

	// 退職資金計算
	currentSavings, err := plan.Profile().CurrentSavings().Total()
	if err != nil {
		return nil, fmt.Errorf("現在の貯蓄合計の計算に失敗しました: %w", err)
	}

	netSavings, err := plan.Profile().CalculateNetSavings()
	if err != nil {
		return nil, fmt.Errorf("純貯蓄額の計算に失敗しました: %w", err)
	}

	calculation, err := retirementData.CalculateRetirementSufficiency(
		currentSavings,
		netSavings,
		plan.Profile().InvestmentReturn(),
		plan.Profile().InflationRate(),
	)
	if err != nil {
		return nil, fmt.Errorf("退職資金計算に失敗しました: %w", err)
	}

	// 退職予測を生成
	projections := uc.generateRetirementProjections(plan, retirementData)

	// 退職戦略を生成
	strategies := uc.generateRetirementStrategies(calculation, plan)

	// 推奨事項を生成
	recommendations := uc.generateRetirementRecommendations(calculation)

	// リスク評価を実行
	riskAssessment := uc.assessRetirementRisks(plan, calculation)

	report := RetirementPlanReport{
		UserID:          input.UserID,
		RetirementData:  retirementData,
		Calculation:     calculation,
		Projections:     projections,
		Strategies:      strategies,
		Recommendations: recommendations,
		RiskAssessment:  riskAssessment,
	}

	return &RetirementPlanReportOutput{
		Report:      report,
		GeneratedAt: time.Now().Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// GenerateComprehensiveReport は包括的レポートを生成する
func (uc *generateReportsUseCaseImpl) GenerateComprehensiveReport(
	ctx context.Context,
	input ComprehensiveReportInput,
) (*ComprehensiveReportOutput, error) {
	// 各種レポートを生成
	financialSummary, err := uc.GenerateFinancialSummaryReport(ctx, FinancialSummaryReportInput{
		UserID: input.UserID,
	})
	if err != nil {
		return nil, fmt.Errorf("財務サマリーレポートの生成に失敗しました: %w", err)
	}

	assetProjection, err := uc.GenerateAssetProjectionReport(ctx, AssetProjectionReportInput(input))
	if err != nil {
		return nil, fmt.Errorf("資産推移レポートの生成に失敗しました: %w", err)
	}

	goalsProgress, err := uc.GenerateGoalsProgressReport(ctx, GoalsProgressReportInput{
		UserID: input.UserID,
	})
	if err != nil {
		return nil, fmt.Errorf("目標進捗レポートの生成に失敗しました: %w", err)
	}

	// 退職計画レポート（オプション）
	var retirementPlan *RetirementPlanReport
	retirementReport, err := uc.GenerateRetirementPlanReport(ctx, RetirementPlanReportInput{
		UserID: input.UserID,
	})
	if err == nil {
		retirementPlan = &retirementReport.Report
	}

	// エグゼクティブサマリーを生成
	executiveSummary := uc.generateExecutiveSummary(
		&financialSummary.Report,
		&assetProjection.Report,
		&goalsProgress.Report,
		retirementPlan,
	)

	// アクションプランを生成
	actionPlan := uc.generateActionPlan(
		&financialSummary.Report,
		&goalsProgress.Report,
		retirementPlan,
	)

	report := ComprehensiveReport{
		UserID:           input.UserID,
		ExecutiveSummary: executiveSummary,
		FinancialSummary: financialSummary.Report,
		AssetProjection:  assetProjection.Report,
		GoalsProgress:    goalsProgress.Report,
		RetirementPlan:   retirementPlan,
		ActionPlan:       actionPlan,
	}

	return &ComprehensiveReportOutput{
		Report:      report,
		GeneratedAt: time.Now().Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// calculateFinancialHealth は財務健全性を計算する
func (uc *generateReportsUseCaseImpl) calculateFinancialHealth(plan *aggregates.FinancialPlan) (*FinancialHealth, error) {
	// 貯蓄率を計算
	netSavings, err := plan.Profile().CalculateNetSavings()
	if err != nil {
		return nil, err
	}

	monthlyIncome := plan.Profile().MonthlyIncome()
	savingsRate := (netSavings.Amount() / monthlyIncome.Amount()) * 100

	// 緊急資金比率を計算
	monthlyExpenses, err := plan.Profile().MonthlyExpenses().Total()
	if err != nil {
		return nil, err
	}

	emergencyFundRatio := 0.0
	if plan.EmergencyFund() != nil {
		emergencyFundRatio = plan.EmergencyFund().CurrentFund.Amount() / monthlyExpenses.Amount()
	}

	// 総合スコアを計算（簡略化）
	score := 0
	if savingsRate >= 20 {
		score += 30
	} else if savingsRate >= 10 {
		score += 20
	} else if savingsRate >= 5 {
		score += 10
	}

	if emergencyFundRatio >= 6 {
		score += 30
	} else if emergencyFundRatio >= 3 {
		score += 20
	} else if emergencyFundRatio >= 1 {
		score += 10
	}

	// 投資利回りによる加点
	investmentReturn := plan.Profile().InvestmentReturn().AsPercentage()
	if investmentReturn >= 5 {
		score += 20
	} else if investmentReturn >= 3 {
		score += 15
	} else if investmentReturn >= 1 {
		score += 10
	}

	// 債務対収入比率（簡略化：0と仮定）
	debtToIncomeRatio := 0.0

	// スコアレベルを決定
	var scoreLevel string
	switch {
	case score >= 80:
		scoreLevel = "excellent"
	case score >= 60:
		scoreLevel = "good"
	case score >= 40:
		scoreLevel = "fair"
	default:
		scoreLevel = "poor"
	}

	return &FinancialHealth{
		OverallScore:       score,
		ScoreLevel:         scoreLevel,
		SavingsRate:        savingsRate,
		DebtToIncomeRatio:  debtToIncomeRatio,
		EmergencyFundRatio: emergencyFundRatio,
	}, nil
}

// getCurrentSituation は現在の状況を取得する
func (uc *generateReportsUseCaseImpl) getCurrentSituation(plan *aggregates.FinancialPlan) (*CurrentSituation, error) {
	monthlyExpenses, err := plan.Profile().MonthlyExpenses().Total()
	if err != nil {
		return nil, err
	}

	netSavings, err := plan.Profile().CalculateNetSavings()
	if err != nil {
		return nil, err
	}

	totalAssets, err := plan.Profile().CurrentSavings().Total()
	if err != nil {
		return nil, err
	}

	return &CurrentSituation{
		MonthlyIncome:    plan.Profile().MonthlyIncome().Amount(),
		MonthlyExpenses:  monthlyExpenses.Amount(),
		NetSavings:       netSavings.Amount(),
		TotalAssets:      totalAssets.Amount(),
		InvestmentReturn: plan.Profile().InvestmentReturn().AsPercentage(),
		InflationRate:    plan.Profile().InflationRate().AsPercentage(),
	}, nil
}

// calculateKeyMetrics は主要指標を計算する
func (uc *generateReportsUseCaseImpl) calculateKeyMetrics(plan *aggregates.FinancialPlan) ([]KeyMetric, error) {
	var metrics []KeyMetric

	// 貯蓄率
	netSavings, err := plan.Profile().CalculateNetSavings()
	if err != nil {
		return nil, err
	}

	monthlyIncome := plan.Profile().MonthlyIncome()
	savingsRate := (netSavings.Amount() / monthlyIncome.Amount()) * 100

	metrics = append(metrics, KeyMetric{
		Name:        "貯蓄率",
		Value:       savingsRate,
		Unit:        "%",
		Description: "月収に対する純貯蓄額の割合",
		Trend:       "stable", // 実際の実装では履歴データから計算
	})

	// 投資利回り
	metrics = append(metrics, KeyMetric{
		Name:        "投資利回り",
		Value:       plan.Profile().InvestmentReturn().AsPercentage(),
		Unit:        "%",
		Description: "年間の期待投資収益率",
		Trend:       "stable",
	})

	// 総資産
	totalAssets, err := plan.Profile().CurrentSavings().Total()
	if err != nil {
		return nil, err
	}

	metrics = append(metrics, KeyMetric{
		Name:        "総資産",
		Value:       totalAssets.Amount(),
		Unit:        "円",
		Description: "現在の総貯蓄・投資額",
		Trend:       "up",
	})

	return metrics, nil
}

// generateRecommendationsAndWarnings は推奨事項と警告を生成する
func (uc *generateReportsUseCaseImpl) generateRecommendationsAndWarnings(plan *aggregates.FinancialPlan) ([]string, []string) {
	var recommendations []string
	var warnings []string

	// 貯蓄率チェック
	netSavings, err := plan.Profile().CalculateNetSavings()
	if err == nil {
		monthlyIncome := plan.Profile().MonthlyIncome()
		savingsRate := (netSavings.Amount() / monthlyIncome.Amount()) * 100

		if savingsRate < 10 {
			warnings = append(warnings, "貯蓄率が10%を下回っています。支出の見直しを検討してください")
			recommendations = append(recommendations, "月間支出を詳細に分析し、削減可能な項目を特定してください")
		} else if savingsRate > 30 {
			recommendations = append(recommendations, "優秀な貯蓄率です。投資商品の多様化を検討してください")
		}
	}

	// 緊急資金チェック
	if plan.EmergencyFund() != nil {
		monthlyExpenses, err := plan.Profile().MonthlyExpenses().Total()
		if err == nil {
			emergencyFundRatio := plan.EmergencyFund().CurrentFund.Amount() / monthlyExpenses.Amount()

			if emergencyFundRatio < 3 {
				warnings = append(warnings, "緊急資金が3ヶ月分の生活費を下回っています")
				recommendations = append(recommendations, "緊急資金として3-6ヶ月分の生活費を確保してください")
			}
		}
	}

	// 投資利回りチェック
	investmentReturn := plan.Profile().InvestmentReturn().AsPercentage()
	if investmentReturn < 3 {
		recommendations = append(recommendations, "投資利回りが低めです。ポートフォリオの見直しを検討してください")
	}

	return recommendations, warnings
}

// その他のヘルパーメソッドは簡略化のため省略
// 実際の実装では以下のメソッドも必要：
// - calculateProjectionSummary
// - generateScenarioAnalysis
// - generateProjectionInsights
// - getGoalStatusText
// - generateAchievements
// - generateNextSteps
// - generateRetirementProjections
// - generateRetirementStrategies
// - generateRetirementRecommendations
// - assessRetirementRisks
// - generateExecutiveSummary
// - generateActionPlan

// calculateProjectionSummary は予測サマリーを計算する（簡略版）
func (uc *generateReportsUseCaseImpl) calculateProjectionSummary(projections []entities.AssetProjection) (*ProjectionSummary, error) {
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

// generateScenarioAnalysis はシナリオ分析を生成する（簡略版）
func (uc *generateReportsUseCaseImpl) generateScenarioAnalysis(plan *aggregates.FinancialPlan, years int) []ScenarioAnalysis {
	// 楽観的、標準、悲観的シナリオを生成
	scenarios := []ScenarioAnalysis{
		{
			Name:             "楽観的シナリオ",
			Description:      "市場が好調で高い投資収益が期待できる場合",
			InvestmentReturn: plan.Profile().InvestmentReturn().AsPercentage() + 2,
			InflationRate:    plan.Profile().InflationRate().AsPercentage(),
			Impact:           "資産形成が加速します",
		},
		{
			Name:             "標準シナリオ",
			Description:      "現在の前提条件が継続する場合",
			InvestmentReturn: plan.Profile().InvestmentReturn().AsPercentage(),
			InflationRate:    plan.Profile().InflationRate().AsPercentage(),
			Impact:           "計画通りの資産形成が期待できます",
		},
		{
			Name:             "悲観的シナリオ",
			Description:      "市場が低迷し投資収益が低下する場合",
			InvestmentReturn: plan.Profile().InvestmentReturn().AsPercentage() - 2,
			InflationRate:    plan.Profile().InflationRate().AsPercentage() + 1,
			Impact:           "目標達成が困難になる可能性があります",
		},
	}

	// 各シナリオの最終金額を計算（簡略化）
	for i := range scenarios {
		// 実際の計算ロジックを実装
		scenarios[i].FinalAmount = 1000000 // プレースホルダー
		scenarios[i].RealValue = 900000    // プレースホルダー
	}

	return scenarios
}

// generateProjectionInsights は予測洞察を生成する（簡略版）
func (uc *generateReportsUseCaseImpl) generateProjectionInsights(projections []entities.AssetProjection, scenarios []ScenarioAnalysis) []string {
	var insights []string

	if len(projections) > 0 {
		finalProjection := projections[len(projections)-1]
		compoundEffect := finalProjection.InvestmentGains.Amount() / finalProjection.ContributedAmount.Amount() * 100

		if compoundEffect > 100 {
			insights = append(insights, "複利効果により投資収益が元本を上回る見込みです")
		}

		insights = append(insights, "長期投資により安定した資産形成が期待できます")
	}

	return insights
}

// getGoalStatusText は目標の状態テキストを取得する（簡略版）
func (uc *generateReportsUseCaseImpl) getGoalStatusText(goal *entities.Goal) string {
	if goal.IsCompleted() {
		return "達成済み"
	}
	if goal.IsOverdue() {
		return "期限切れ"
	}
	if !goal.IsActive() {
		return "非アクティブ"
	}
	return "進行中"
}

// generateAchievements は達成事項を生成する（簡略版）
func (uc *generateReportsUseCaseImpl) generateAchievements(goals []*entities.Goal) []Achievement {
	var achievements []Achievement

	for _, goal := range goals {
		if goal.IsCompleted() {
			achievements = append(achievements, Achievement{
				Type:        "goal_completion",
				Title:       fmt.Sprintf("%s達成", goal.Title()),
				Description: fmt.Sprintf("目標金額%sを達成しました", goal.TargetAmount().String()),
				Date:        goal.UpdatedAt().Format("2006-01-02"),
				Impact:      "財務目標の達成により安心感が向上しました",
			})
		}
	}

	return achievements
}

// generateNextSteps は次のステップを生成する（簡略版）
func (uc *generateReportsUseCaseImpl) generateNextSteps(goalProgresses []GoalProgress) []string {
	var nextSteps []string

	for _, progress := range goalProgresses {
		if !progress.OnTrack && progress.Goal.IsActive() {
			nextSteps = append(nextSteps, fmt.Sprintf("%sの進捗改善が必要です", progress.Goal.Title()))
		}
	}

	if len(nextSteps) == 0 {
		nextSteps = append(nextSteps, "現在の計画を継続してください")
	}

	return nextSteps
}

// generateRetirementProjections は退職予測を生成する（簡略版）
func (uc *generateReportsUseCaseImpl) generateRetirementProjections(plan *aggregates.FinancialPlan, retirementData *entities.RetirementData) []RetirementProjection {
	// 簡略化された実装
	return []RetirementProjection{
		{
			Age:               65,
			YearsToRetirement: 30,
			ProjectedAssets:   50000000,
			RequiredAssets:    60000000,
			SufficiencyRate:   83.3,
			MonthlyShortfall:  50000,
		},
	}
}

// generateRetirementStrategies は退職戦略を生成する（簡略版）
func (uc *generateReportsUseCaseImpl) generateRetirementStrategies(calculation *entities.RetirementCalculation, plan *aggregates.FinancialPlan) []RetirementStrategy {
	return []RetirementStrategy{
		{
			Name:        "貯蓄額増加",
			Description: "月間貯蓄額を増やして退職資金を確保する",
			Impact:      100000,
			Effort:      "medium",
			Timeline:    "即座に開始可能",
		},
	}
}

// generateRetirementRecommendations は退職推奨事項を生成する（簡略版）
func (uc *generateReportsUseCaseImpl) generateRetirementRecommendations(calculation *entities.RetirementCalculation) []string {
	return []string{
		"月間貯蓄額の増加を検討してください",
		"投資ポートフォリオの見直しを行ってください",
	}
}

// assessRetirementRisks は退職リスクを評価する（簡略版）
func (uc *generateReportsUseCaseImpl) assessRetirementRisks(plan *aggregates.FinancialPlan, calculation *entities.RetirementCalculation) RiskAssessment {
	return RiskAssessment{
		OverallRisk: "medium",
		RiskFactors: []RiskFactor{
			{
				Type:        "longevity_risk",
				Description: "予想より長生きした場合の資金不足リスク",
				Impact:      "high",
				Probability: "medium",
			},
		},
		Mitigations: []string{
			"健康管理による医療費削減",
			"副収入源の確保",
		},
	}
}

// generateExecutiveSummary はエグゼクティブサマリーを生成する（簡略版）
func (uc *generateReportsUseCaseImpl) generateExecutiveSummary(
	financialSummary *FinancialSummaryReport,
	assetProjection *AssetProjectionReport,
	goalsProgress *GoalsProgressReport,
	retirementPlan *RetirementPlanReport,
) ExecutiveSummary {
	return ExecutiveSummary{
		OverallStatus:        "良好",
		KeyHighlights:        []string{"貯蓄率が理想的", "目標進捗が順調"},
		CriticalActions:      []string{"緊急資金の確保"},
		OpportunityAreas:     []string{"投資利回りの改善"},
		FinancialHealthScore: financialSummary.FinancialHealth.OverallScore,
	}
}

// generateActionPlan はアクションプランを生成する（簡略版）
func (uc *generateReportsUseCaseImpl) generateActionPlan(
	financialSummary *FinancialSummaryReport,
	goalsProgress *GoalsProgressReport,
	retirementPlan *RetirementPlanReport,
) ActionPlan {
	return ActionPlan{
		ShortTerm: []ActionItem{
			{
				Priority:    "high",
				Title:       "緊急資金の確保",
				Description: "3ヶ月分の生活費を緊急資金として確保する",
				Timeline:    "3ヶ月以内",
				Impact:      "リスク軽減",
				Effort:      "medium",
			},
		},
		MediumTerm: []ActionItem{
			{
				Priority:    "medium",
				Title:       "投資ポートフォリオの見直し",
				Description: "リスク分散と利回り向上のためのポートフォリオ最適化",
				Timeline:    "6ヶ月以内",
				Impact:      "収益向上",
				Effort:      "low",
			},
		},
		LongTerm: []ActionItem{
			{
				Priority:    "medium",
				Title:       "退職計画の詳細化",
				Description: "具体的な退職後の生活設計と資金計画の策定",
				Timeline:    "1年以内",
				Impact:      "安心感向上",
				Effort:      "high",
			},
		},
	}
}

// ExportReportToPDF はレポートをPDF形式でエクスポートする
func (uc *generateReportsUseCaseImpl) ExportReportToPDF(
	ctx context.Context,
	input ExportReportInput,
) (*ExportReportOutput, error) {
	// TODO: 実際のPDF生成ロジックを実装
	// ここでは簡易的なダミーPDFを生成
	pdfContent := uc.generateDummyPDF(input)

	// 一時ファイルストレージに保存（実際の実装では依存性注入で渡す）
	// ここでは簡易的な実装として、ファイル名とサイズのみ返す
	fileName := fmt.Sprintf("%s_report_%s.pdf", input.ReportType, time.Now().Format("20060102_150405"))
	fileSize := int64(len(pdfContent))

	// 有効期限を24時間後に設定
	expiresAt := time.Now().Add(24 * time.Hour)

	// ダウンロードURLを生成（実際の実装では署名付きURLを生成）
	// トークンを生成（簡易版）
	token := uc.generateDownloadToken(fileName, expiresAt)
	downloadURL := fmt.Sprintf("/api/reports/download/%s", token)

	return &ExportReportOutput{
		FileName:    fileName,
		FileSize:    fileSize,
		DownloadURL: downloadURL,
		ExpiresAt:   expiresAt.Format(time.RFC3339),
	}, nil
}

// generateDummyPDF はダミーのPDFコンテンツを生成する
func (uc *generateReportsUseCaseImpl) generateDummyPDF(input ExportReportInput) []byte {
	// 実際の実装では、PDFライブラリ（例: gopdf, gofpdf）を使用
	content := fmt.Sprintf("PDF Report\nType: %s\nFormat: %s\nGenerated: %s\n",
		input.ReportType,
		input.Format,
		time.Now().Format(time.RFC3339),
	)
	return []byte(content)
}

// generateDownloadToken はダウンロード用のトークンを生成する
func (uc *generateReportsUseCaseImpl) generateDownloadToken(fileName string, expiresAt time.Time) string {
	// 実際の実装では、HMAC-SHA256などで署名付きトークンを生成
	// ここでは簡易的な実装
	return fmt.Sprintf("%s_%d", fileName, expiresAt.Unix())
}
