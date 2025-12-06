package usecases

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/financial-planning-calculator/backend/domain/entities"
	"github.com/financial-planning-calculator/backend/domain/repositories"
	"github.com/financial-planning-calculator/backend/domain/services"
	"github.com/financial-planning-calculator/backend/domain/valueobjects"
)

// ManageGoalsUseCase は目標管理のユースケース
type ManageGoalsUseCase interface {
	// CreateGoal は新しい目標を作成する
	CreateGoal(ctx context.Context, input CreateGoalInput) (*CreateGoalOutput, error)

	// GetGoal は目標を取得する
	GetGoal(ctx context.Context, input GetGoalInput) (*GetGoalOutput, error)

	// GetGoalsByUser はユーザーの目標一覧を取得する
	GetGoalsByUser(ctx context.Context, input GetGoalsByUserInput) (*GetGoalsByUserOutput, error)

	// UpdateGoal は目標を更新する
	UpdateGoal(ctx context.Context, input UpdateGoalInput) (*UpdateGoalOutput, error)

	// UpdateGoalProgress は目標の進捗を更新する
	UpdateGoalProgress(ctx context.Context, input UpdateGoalProgressInput) (*UpdateGoalProgressOutput, error)

	// DeleteGoal は目標を削除する
	DeleteGoal(ctx context.Context, input DeleteGoalInput) error

	// GetGoalRecommendations は目標の推奨事項を取得する
	GetGoalRecommendations(ctx context.Context, input GetGoalRecommendationsInput) (*GetGoalRecommendationsOutput, error)

	// AnalyzeGoalFeasibility は目標の実現可能性を分析する
	AnalyzeGoalFeasibility(ctx context.Context, input AnalyzeGoalFeasibilityInput) (*AnalyzeGoalFeasibilityOutput, error)
}

// CreateGoalInput は目標作成の入力
type CreateGoalInput struct {
	UserID              entities.UserID `json:"user_id"`
	GoalType            string          `json:"goal_type"`
	Title               string          `json:"title"`
	TargetAmount        float64         `json:"target_amount"`
	TargetDate          string          `json:"target_date"` // RFC3339 format
	CurrentAmount       float64         `json:"current_amount"`
	MonthlyContribution float64         `json:"monthly_contribution"`
	Description         *string         `json:"description,omitempty"`
}

// CreateGoalOutput は目標作成の出力
type CreateGoalOutput struct {
	GoalID    entities.GoalID `json:"goal_id"`
	UserID    entities.UserID `json:"user_id"`
	CreatedAt string          `json:"created_at"`
}

// GetGoalInput は目標取得の入力
type GetGoalInput struct {
	GoalID entities.GoalID `json:"goal_id"`
	UserID entities.UserID `json:"user_id"`
}

// GetGoalOutput は目標取得の出力
type GetGoalOutput struct {
	Goal     *entities.Goal        `json:"goal"`
	Progress entities.ProgressRate `json:"progress"`
	Status   GoalStatus            `json:"status"`
}

// GoalStatus は目標の状態
type GoalStatus struct {
	IsActive    bool   `json:"is_active"`
	IsCompleted bool   `json:"is_completed"`
	IsOverdue   bool   `json:"is_overdue"`
	DaysLeft    int    `json:"days_left"`
	Message     string `json:"message"`
}

// GetGoalsByUserInput はユーザー目標一覧取得の入力
type GetGoalsByUserInput struct {
	UserID     entities.UserID    `json:"user_id"`
	GoalType   *entities.GoalType `json:"goal_type,omitempty"`
	ActiveOnly bool               `json:"active_only"`
}

// GetGoalsByUserOutput はユーザー目標一覧取得の出力
type GetGoalsByUserOutput struct {
	Goals   []GoalWithStatus `json:"goals"`
	Summary GoalsSummary     `json:"summary"`
}

// GoalWithStatus は状態付きの目標
type GoalWithStatus struct {
	Goal     *entities.Goal        `json:"goal"`
	Progress entities.ProgressRate `json:"progress"`
	Status   GoalStatus            `json:"status"`
}

// GoalsSummary は目標のサマリー
type GoalsSummary struct {
	TotalGoals      int     `json:"total_goals"`
	ActiveGoals     int     `json:"active_goals"`
	CompletedGoals  int     `json:"completed_goals"`
	OverdueGoals    int     `json:"overdue_goals"`
	TotalTarget     float64 `json:"total_target"`
	TotalCurrent    float64 `json:"total_current"`
	OverallProgress float64 `json:"overall_progress"`
}

// UpdateGoalInput は目標更新の入力
type UpdateGoalInput struct {
	GoalID              entities.GoalID `json:"goal_id"`
	UserID              entities.UserID `json:"user_id"`
	Title               *string         `json:"title,omitempty"`
	TargetAmount        *float64        `json:"target_amount,omitempty"`
	TargetDate          *string         `json:"target_date,omitempty"` // RFC3339 format
	MonthlyContribution *float64        `json:"monthly_contribution,omitempty"`
	Description         *string         `json:"description,omitempty"`
	IsActive            *bool           `json:"is_active,omitempty"`
}

// UpdateGoalOutput は目標更新の出力
type UpdateGoalOutput struct {
	Success   bool   `json:"success"`
	UpdatedAt string `json:"updated_at"`
}

// UpdateGoalProgressInput は目標進捗更新の入力
type UpdateGoalProgressInput struct {
	GoalID        entities.GoalID `json:"goal_id"`
	UserID        entities.UserID `json:"user_id"`
	CurrentAmount float64         `json:"current_amount"`
	Note          *string         `json:"note,omitempty"`
}

// UpdateGoalProgressOutput は目標進捗更新の出力
type UpdateGoalProgressOutput struct {
	Success     bool                  `json:"success"`
	NewProgress entities.ProgressRate `json:"new_progress"`
	IsCompleted bool                  `json:"is_completed"`
	UpdatedAt   string                `json:"updated_at"`
}

// DeleteGoalInput は目標削除の入力
type DeleteGoalInput struct {
	GoalID entities.GoalID `json:"goal_id"`
	UserID entities.UserID `json:"user_id"`
}

// GetGoalRecommendationsInput は目標推奨事項取得の入力
type GetGoalRecommendationsInput struct {
	GoalID entities.GoalID `json:"goal_id"`
	UserID entities.UserID `json:"user_id"`
}

// GetGoalRecommendationsOutput は目標推奨事項取得の出力
type GetGoalRecommendationsOutput struct {
	Recommendations []services.GoalRecommendation   `json:"recommendations"`
	SavingsAdvice   *services.SavingsRecommendation `json:"savings_advice"`
}

// AnalyzeGoalFeasibilityInput は目標実現可能性分析の入力
type AnalyzeGoalFeasibilityInput struct {
	GoalID entities.GoalID `json:"goal_id"`
	UserID entities.UserID `json:"user_id"`
}

// AnalyzeGoalFeasibilityOutput は目標実現可能性分析の出力
type AnalyzeGoalFeasibilityOutput struct {
	Feasibility map[string]interface{} `json:"feasibility"`
	RiskLevel   string                 `json:"risk_level"`
	Achievable  bool                   `json:"achievable"`
	Insights    []FeasibilityInsight   `json:"insights"`
}

// FeasibilityInsight は実現可能性の洞察
type FeasibilityInsight struct {
	Type        string `json:"type"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Impact      string `json:"impact"`
	Severity    string `json:"severity"` // "info", "warning", "error"
}

// manageGoalsUseCaseImpl はManageGoalsUseCaseの実装
type manageGoalsUseCaseImpl struct {
	goalRepo              repositories.GoalRepository
	financialPlanRepo     repositories.FinancialPlanRepository
	recommendationService *services.GoalRecommendationService
}

// NewManageGoalsUseCase は新しいManageGoalsUseCaseを作成する
func NewManageGoalsUseCase(
	goalRepo repositories.GoalRepository,
	financialPlanRepo repositories.FinancialPlanRepository,
	recommendationService *services.GoalRecommendationService,
) ManageGoalsUseCase {
	return &manageGoalsUseCaseImpl{
		goalRepo:              goalRepo,
		financialPlanRepo:     financialPlanRepo,
		recommendationService: recommendationService,
	}
}

// CreateGoal は新しい目標を作成する
func (uc *manageGoalsUseCaseImpl) CreateGoal(
	ctx context.Context,
	input CreateGoalInput,
) (*CreateGoalOutput, error) {
	// 目標タイプを解析
	var goalType entities.GoalType
	switch input.GoalType {
	case "savings":
		goalType = entities.GoalTypeSavings
	case "retirement":
		goalType = entities.GoalTypeRetirement
	case "emergency":
		goalType = entities.GoalTypeEmergency
	case "custom":
		goalType = entities.GoalTypeCustom
	default:
		return nil, fmt.Errorf("無効な目標タイプです: %s", input.GoalType)
	}

	// 目標日を解析
	targetDate, err := time.Parse(time.RFC3339, input.TargetDate)
	if err != nil {
		return nil, fmt.Errorf("目標日の解析に失敗しました: %w", err)
	}

	// 金額を作成
	targetAmount, err := valueobjects.NewMoneyJPY(input.TargetAmount)
	if err != nil {
		return nil, fmt.Errorf("目標金額の作成に失敗しました: %w", err)
	}

	currentAmount, err := valueobjects.NewMoneyJPY(input.CurrentAmount)
	if err != nil {
		return nil, fmt.Errorf("現在金額の作成に失敗しました: %w", err)
	}

	monthlyContribution, err := valueobjects.NewMoneyJPY(input.MonthlyContribution)
	if err != nil {
		return nil, fmt.Errorf("月間拠出額の作成に失敗しました: %w", err)
	}

	// 同じタイプの目標が既に存在するかチェック（退職・緊急資金目標は1つまで）
	if goalType == entities.GoalTypeRetirement || goalType == entities.GoalTypeEmergency {
		existingGoals, err := uc.goalRepo.FindByUserIDAndType(ctx, input.UserID, goalType)
		if err != nil {
			return nil, fmt.Errorf("既存目標の確認に失敗しました: %w", err)
		}

		for _, existingGoal := range existingGoals {
			if existingGoal.IsActive() {
				return nil, fmt.Errorf("%sの目標は既に存在します", goalType.String())
			}
		}
	}

	// 目標を作成
	goal, err := entities.NewGoal(
		input.UserID,
		goalType,
		input.Title,
		targetAmount,
		targetDate,
		monthlyContribution,
	)
	if err != nil {
		return nil, fmt.Errorf("目標の作成に失敗しました: %w", err)
	}

	// 現在金額を設定
	err = goal.UpdateCurrentAmount(currentAmount)
	if err != nil {
		return nil, fmt.Errorf("現在金額の設定に失敗しました: %w", err)
	}

	// 財務計画を取得して達成可能性をチェック
	plan, err := uc.financialPlanRepo.FindByUserID(ctx, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("財務計画の取得に失敗しました: %w", err)
	}

	achievable, err := goal.IsAchievable(plan.Profile())
	if err != nil {
		return nil, fmt.Errorf("目標の達成可能性チェックに失敗しました: %w", err)
	}

	if !achievable {
		return nil, errors.New("現在の財務状況では目標の達成が困難です。目標金額または期日の調整を検討してください")
	}

	// 目標を保存
	err = uc.goalRepo.Save(ctx, goal)
	if err != nil {
		return nil, fmt.Errorf("目標の保存に失敗しました: %w", err)
	}

	// 財務計画に目標を追加
	err = plan.AddGoal(goal)
	if err != nil {
		return nil, fmt.Errorf("財務計画への目標追加に失敗しました: %w", err)
	}

	err = uc.financialPlanRepo.Update(ctx, plan)
	if err != nil {
		return nil, fmt.Errorf("財務計画の更新に失敗しました: %w", err)
	}

	return &CreateGoalOutput{
		GoalID:    goal.ID(),
		UserID:    input.UserID,
		CreatedAt: goal.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// GetGoal は目標を取得する
func (uc *manageGoalsUseCaseImpl) GetGoal(
	ctx context.Context,
	input GetGoalInput,
) (*GetGoalOutput, error) {
	// 目標を取得
	goal, err := uc.goalRepo.FindByID(ctx, input.GoalID)
	if err != nil {
		return nil, fmt.Errorf("目標の取得に失敗しました: %w", err)
	}

	// ユーザーIDが一致するかチェック
	if goal.UserID() != input.UserID {
		return nil, errors.New("指定された目標にアクセスする権限がありません")
	}

	// 進捗を計算
	progress, err := goal.CalculateProgress(goal.CurrentAmount())
	if err != nil {
		return nil, fmt.Errorf("進捗の計算に失敗しました: %w", err)
	}

	// 状態を生成
	status := uc.generateGoalStatus(goal)

	return &GetGoalOutput{
		Goal:     goal,
		Progress: progress,
		Status:   status,
	}, nil
}

// GetGoalsByUser はユーザーの目標一覧を取得する
func (uc *manageGoalsUseCaseImpl) GetGoalsByUser(
	ctx context.Context,
	input GetGoalsByUserInput,
) (*GetGoalsByUserOutput, error) {
	var goals []*entities.Goal
	var err error

	// 目標を取得
	if input.GoalType != nil {
		goals, err = uc.goalRepo.FindByUserIDAndType(ctx, input.UserID, *input.GoalType)
	} else if input.ActiveOnly {
		goals, err = uc.goalRepo.FindActiveGoalsByUserID(ctx, input.UserID)
	} else {
		goals, err = uc.goalRepo.FindByUserID(ctx, input.UserID)
	}

	if err != nil {
		return nil, fmt.Errorf("目標の取得に失敗しました: %w", err)
	}

	// 状態付きの目標リストを作成
	var goalsWithStatus []GoalWithStatus
	var summary GoalsSummary

	for _, goal := range goals {
		progress, err := goal.CalculateProgress(goal.CurrentAmount())
		if err != nil {
			// エラーが発生しても処理を止めずにログを出力し、進捗は0として扱う
			slog.Error("failed to calculate goal progress", "goal_id", goal.ID(), "error", err)
			progress, _ = entities.NewProgressRate(0) // 0% で進捗を初期化 (エラーは無視し、0%とする)
		}

		status := uc.generateGoalStatus(goal)

		goalsWithStatus = append(goalsWithStatus, GoalWithStatus{
			Goal:     goal,
			Progress: progress,
			Status:   status,
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

	return &GetGoalsByUserOutput{
		Goals:   goalsWithStatus,
		Summary: summary,
	}, nil
}

// UpdateGoal は目標を更新する
func (uc *manageGoalsUseCaseImpl) UpdateGoal(
	ctx context.Context,
	input UpdateGoalInput,
) (*UpdateGoalOutput, error) {
	// 目標を取得
	goal, err := uc.goalRepo.FindByID(ctx, input.GoalID)
	if err != nil {
		return nil, fmt.Errorf("目標の取得に失敗しました: %w", err)
	}

	// ユーザーIDが一致するかチェック
	if goal.UserID() != input.UserID {
		return nil, errors.New("指定された目標にアクセスする権限がありません")
	}

	// 更新処理
	if input.Title != nil {
		err = goal.UpdateTitle(*input.Title)
		if err != nil {
			return nil, fmt.Errorf("タイトルの更新に失敗しました: %w", err)
		}
	}

	if input.TargetAmount != nil {
		targetAmount, err := valueobjects.NewMoneyJPY(*input.TargetAmount)
		if err != nil {
			return nil, fmt.Errorf("目標金額の作成に失敗しました: %w", err)
		}

		err = goal.UpdateTargetAmount(targetAmount)
		if err != nil {
			return nil, fmt.Errorf("目標金額の更新に失敗しました: %w", err)
		}
	}

	if input.TargetDate != nil {
		targetDate, err := time.Parse(time.RFC3339, *input.TargetDate)
		if err != nil {
			return nil, fmt.Errorf("目標日の解析に失敗しました: %w", err)
		}

		err = goal.UpdateTargetDate(targetDate)
		if err != nil {
			return nil, fmt.Errorf("目標日の更新に失敗しました: %w", err)
		}
	}

	if input.MonthlyContribution != nil {
		monthlyContribution, err := valueobjects.NewMoneyJPY(*input.MonthlyContribution)
		if err != nil {
			return nil, fmt.Errorf("月間拠出額の作成に失敗しました: %w", err)
		}

		err = goal.UpdateMonthlyContribution(monthlyContribution)
		if err != nil {
			return nil, fmt.Errorf("月間拠出額の更新に失敗しました: %w", err)
		}
	}

	// Note: Description update is not available in the current Goal entity
	// This would need to be added to the Goal entity if required

	if input.IsActive != nil {
		if *input.IsActive {
			goal.Activate()
		} else {
			goal.Deactivate()
		}
	}

	// 目標を保存
	err = uc.goalRepo.Update(ctx, goal)
	if err != nil {
		return nil, fmt.Errorf("目標の保存に失敗しました: %w", err)
	}

	return &UpdateGoalOutput{
		Success:   true,
		UpdatedAt: goal.UpdatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// UpdateGoalProgress は目標の進捗を更新する
func (uc *manageGoalsUseCaseImpl) UpdateGoalProgress(
	ctx context.Context,
	input UpdateGoalProgressInput,
) (*UpdateGoalProgressOutput, error) {
	// 目標を取得
	goal, err := uc.goalRepo.FindByID(ctx, input.GoalID)
	if err != nil {
		return nil, fmt.Errorf("目標の取得に失敗しました: %w", err)
	}

	// ユーザーIDが一致するかチェック
	if goal.UserID() != input.UserID {
		return nil, errors.New("指定された目標にアクセスする権限がありません")
	}

	// 現在金額を更新
	currentAmount, err := valueobjects.NewMoneyJPY(input.CurrentAmount)
	if err != nil {
		return nil, fmt.Errorf("現在金額の作成に失敗しました: %w", err)
	}

	err = goal.UpdateCurrentAmount(currentAmount)
	if err != nil {
		return nil, fmt.Errorf("現在金額の更新に失敗しました: %w", err)
	}

	// 進捗を計算
	progress, err := goal.CalculateProgress(currentAmount)
	if err != nil {
		return nil, fmt.Errorf("進捗の計算に失敗しました: %w", err)
	}

	// 完了チェック
	isCompleted := goal.IsCompleted()

	// 目標を保存
	err = uc.goalRepo.Update(ctx, goal)
	if err != nil {
		return nil, fmt.Errorf("目標の保存に失敗しました: %w", err)
	}

	return &UpdateGoalProgressOutput{
		Success:     true,
		NewProgress: progress,
		IsCompleted: isCompleted,
		UpdatedAt:   goal.UpdatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// DeleteGoal は目標を削除する
func (uc *manageGoalsUseCaseImpl) DeleteGoal(
	ctx context.Context,
	input DeleteGoalInput,
) error {
	// 目標を取得
	goal, err := uc.goalRepo.FindByID(ctx, input.GoalID)
	if err != nil {
		return fmt.Errorf("目標の取得に失敗しました: %w", err)
	}

	// ユーザーIDが一致するかチェック
	if goal.UserID() != input.UserID {
		return errors.New("指定された目標にアクセスする権限がありません")
	}

	// 財務計画から目標を削除
	plan, err := uc.financialPlanRepo.FindByUserID(ctx, input.UserID)
	if err != nil {
		return fmt.Errorf("財務計画の取得に失敗しました: %w", err)
	}

	err = plan.RemoveGoal(input.GoalID)
	if err != nil {
		return fmt.Errorf("財務計画からの目標削除に失敗しました: %w", err)
	}

	err = uc.financialPlanRepo.Update(ctx, plan)
	if err != nil {
		return fmt.Errorf("財務計画の更新に失敗しました: %w", err)
	}

	// 目標を削除
	err = uc.goalRepo.Delete(ctx, input.GoalID)
	if err != nil {
		return fmt.Errorf("目標の削除に失敗しました: %w", err)
	}

	return nil
}

// GetGoalRecommendations は目標の推奨事項を取得する
func (uc *manageGoalsUseCaseImpl) GetGoalRecommendations(
	ctx context.Context,
	input GetGoalRecommendationsInput,
) (*GetGoalRecommendationsOutput, error) {
	// 目標を取得
	goal, err := uc.goalRepo.FindByID(ctx, input.GoalID)
	if err != nil {
		return nil, fmt.Errorf("目標の取得に失敗しました: %w", err)
	}

	// ユーザーIDが一致するかチェック
	if goal.UserID() != input.UserID {
		return nil, errors.New("指定された目標にアクセスする権限がありません")
	}

	// 財務計画を取得
	plan, err := uc.financialPlanRepo.FindByUserID(ctx, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("財務計画の取得に失敗しました: %w", err)
	}

	// 推奨事項を生成
	recommendations, err := uc.recommendationService.SuggestGoalAdjustments(goal, plan.Profile())
	if err != nil {
		return nil, fmt.Errorf("推奨事項の生成に失敗しました: %w", err)
	}

	// 貯蓄推奨を生成
	remainingDays := goal.GetRemainingDays()
	remainingMonths := remainingDays / 30 // 概算
	timeRemaining, err := valueobjects.NewPeriodFromMonths(remainingMonths)
	if err != nil {
		return nil, fmt.Errorf("残り期間の計算に失敗しました: %w", err)
	}

	currentSavings, err := plan.Profile().CurrentSavings().Total()
	if err != nil {
		return nil, fmt.Errorf("現在の貯蓄合計の計算に失敗しました: %w", err)
	}

	savingsAdvice, err := uc.recommendationService.RecommendMonthlySavings(goal, currentSavings, timeRemaining)
	if err != nil {
		return nil, fmt.Errorf("貯蓄推奨の生成に失敗しました: %w", err)
	}

	return &GetGoalRecommendationsOutput{
		Recommendations: recommendations,
		SavingsAdvice:   savingsAdvice,
	}, nil
}

// AnalyzeGoalFeasibility は目標の実現可能性を分析する
func (uc *manageGoalsUseCaseImpl) AnalyzeGoalFeasibility(
	ctx context.Context,
	input AnalyzeGoalFeasibilityInput,
) (*AnalyzeGoalFeasibilityOutput, error) {
	// 目標を取得
	goal, err := uc.goalRepo.FindByID(ctx, input.GoalID)
	if err != nil {
		return nil, fmt.Errorf("目標の取得に失敗しました: %w", err)
	}

	// ユーザーIDが一致するかチェック
	if goal.UserID() != input.UserID {
		return nil, errors.New("指定された目標にアクセスする権限がありません")
	}

	// 財務計画を取得
	plan, err := uc.financialPlanRepo.FindByUserID(ctx, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("財務計画の取得に失敗しました: %w", err)
	}

	// 実現可能性を分析
	feasibility, err := uc.recommendationService.AnalyzeGoalFeasibility(goal, plan.Profile())
	if err != nil {
		return nil, fmt.Errorf("実現可能性の分析に失敗しました: %w", err)
	}

	// 達成可能性を判定
	achievable, err := goal.IsAchievable(plan.Profile())
	if err != nil {
		return nil, fmt.Errorf("達成可能性の判定に失敗しました: %w", err)
	}

	// リスクレベルを取得
	riskLevel, ok := feasibility["risk_level"].(string)
	if !ok {
		riskLevel = "不明"
	}

	// 洞察を生成
	insights := uc.generateFeasibilityInsights(goal, plan.Profile(), feasibility)

	return &AnalyzeGoalFeasibilityOutput{
		Feasibility: feasibility,
		RiskLevel:   riskLevel,
		Achievable:  achievable,
		Insights:    insights,
	}, nil
}

// generateGoalStatus は目標の状態を生成する
func (uc *manageGoalsUseCaseImpl) generateGoalStatus(goal *entities.Goal) GoalStatus {
	isActive := goal.IsActive()
	isCompleted := goal.IsCompleted()
	isOverdue := goal.IsOverdue()
	daysLeft := goal.GetRemainingDays()

	var message string
	switch {
	case isCompleted:
		message = "目標を達成しました！"
	case isOverdue:
		message = "目標期限を過ぎています"
	case daysLeft <= 30:
		message = "目標期限が近づいています"
	case !isActive:
		message = "目標は非アクティブです"
	default:
		message = "順調に進行中です"
	}

	return GoalStatus{
		IsActive:    isActive,
		IsCompleted: isCompleted,
		IsOverdue:   isOverdue,
		DaysLeft:    daysLeft,
		Message:     message,
	}
}

// generateFeasibilityInsights は実現可能性の洞察を生成する
func (uc *manageGoalsUseCaseImpl) generateFeasibilityInsights(
	goal *entities.Goal,
	profile *entities.FinancialProfile,
	feasibility map[string]interface{},
) []FeasibilityInsight {
	var insights []FeasibilityInsight

	// 進捗率の洞察
	if progressPercentage, ok := feasibility["progress_percentage"].(float64); ok {
		if progressPercentage < 25 {
			insights = append(insights, FeasibilityInsight{
				Type:        "progress",
				Title:       "進捗が遅れています",
				Description: fmt.Sprintf("現在の進捗率は%.1f%%です", progressPercentage),
				Impact:      "目標達成のためにはペースアップが必要です",
				Severity:    "warning",
			})
		} else if progressPercentage > 75 {
			insights = append(insights, FeasibilityInsight{
				Type:        "progress",
				Title:       "順調に進捗しています",
				Description: fmt.Sprintf("現在の進捗率は%.1f%%です", progressPercentage),
				Impact:      "このペースを維持すれば目標達成が期待できます",
				Severity:    "info",
			})
		}
	}

	// 必要貯蓄額の洞察
	if requiredMonthlySavings, ok := feasibility["required_monthly_savings"].(float64); ok {
		netSavings, err := profile.CalculateNetSavings()
		if err == nil {
			if requiredMonthlySavings > netSavings.Amount() {
				shortfall := requiredMonthlySavings - netSavings.Amount()
				insights = append(insights, FeasibilityInsight{
					Type:        "savings",
					Title:       "貯蓄額が不足しています",
					Description: fmt.Sprintf("月間%.0f円の追加貯蓄が必要です", shortfall),
					Impact:      "支出の見直しまたは収入の増加を検討してください",
					Severity:    "error",
				})
			}
		}
	}

	// 残り日数の洞察
	remainingDays := goal.GetRemainingDays()
	if remainingDays <= 90 && remainingDays > 0 {
		insights = append(insights, FeasibilityInsight{
			Type:        "timeline",
			Title:       "目標期限が近づいています",
			Description: fmt.Sprintf("残り%d日です", remainingDays),
			Impact:      "最終的な調整や集中的な取り組みが必要です",
			Severity:    "warning",
		})
	} else if remainingDays <= 0 {
		insights = append(insights, FeasibilityInsight{
			Type:        "timeline",
			Title:       "目標期限を過ぎています",
			Description: "期限の延長または目標の見直しが必要です",
			Impact:      "新しい計画を立て直すことをお勧めします",
			Severity:    "error",
		})
	}

	return insights
}
