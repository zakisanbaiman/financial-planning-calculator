package aggregates

import (
	"errors"
	"fmt"
	"time"

	"financial-planning-calculator/domain/entities"
	"financial-planning-calculator/domain/valueobjects"

	"github.com/google/uuid"
)

// FinancialPlanID は財務計画の一意識別子
type FinancialPlanID string

// NewFinancialPlanID は新しい財務計画IDを生成する
func NewFinancialPlanID() FinancialPlanID {
	return FinancialPlanID(uuid.New().String())
}

// ValidationError はバリデーションエラーを表す
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Error はValidationErrorのエラーメッセージを返す
func (ve ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", ve.Field, ve.Message)
}

// PlanProjection は財務計画の将来予測を表す
type PlanProjection struct {
	AssetProjections      []entities.AssetProjection      `json:"asset_projections"`
	RetirementCalculation *entities.RetirementCalculation `json:"retirement_calculation,omitempty"`
	EmergencyFundStatus   *EmergencyFundStatus            `json:"emergency_fund_status,omitempty"`
	GoalProgress          []GoalProgress                  `json:"goal_progress"`
}

// EmergencyFundStatus は緊急資金の状況を表す
type EmergencyFundStatus struct {
	RequiredAmount valueobjects.Money `json:"required_amount"`
	CurrentAmount  valueobjects.Money `json:"current_amount"`
	Shortfall      valueobjects.Money `json:"shortfall"`
	MonthsToTarget int                `json:"months_to_target"`
}

// GoalProgress は目標の進捗状況を表す
type GoalProgress struct {
	Goal     *entities.Goal        `json:"goal"`
	Progress entities.ProgressRate `json:"progress"`
	OnTrack  bool                  `json:"on_track"`
	Message  string                `json:"message"`
}

// FinancialPlan は財務計画の集約ルート
type FinancialPlan struct {
	id             FinancialPlanID
	profile        *entities.FinancialProfile
	goals          []*entities.Goal
	retirementData *entities.RetirementData
	emergencyFund  *EmergencyFundConfig
	createdAt      time.Time
	updatedAt      time.Time
}

// EmergencyFundConfig は緊急資金の設定を表す
type EmergencyFundConfig struct {
	TargetMonths int                `json:"target_months"` // 何ヶ月分の生活費を確保するか
	CurrentFund  valueobjects.Money `json:"current_fund"`  // 現在の緊急資金額
}

// NewEmergencyFundConfig は新しい緊急資金設定を作成する
func NewEmergencyFundConfig(targetMonths int, currentFund valueobjects.Money) (*EmergencyFundConfig, error) {
	if targetMonths < 0 {
		return nil, errors.New("緊急資金の目標月数は負の値にできません")
	}

	if targetMonths > 24 {
		return nil, errors.New("緊急資金の目標月数は24ヶ月以下である必要があります")
	}

	if currentFund.IsNegative() {
		return nil, errors.New("現在の緊急資金は負の値にできません")
	}

	return &EmergencyFundConfig{
		TargetMonths: targetMonths,
		CurrentFund:  currentFund,
	}, nil
}

// NewFinancialPlan は新しい財務計画を作成する
func NewFinancialPlan(profile *entities.FinancialProfile) (*FinancialPlan, error) {
	if profile == nil {
		return nil, errors.New("財務プロファイルは必須です")
	}

	// デフォルトの緊急資金設定（3ヶ月分）
	defaultEmergencyFund, err := valueobjects.NewMoneyJPY(0)
	if err != nil {
		return nil, fmt.Errorf("デフォルト緊急資金の作成に失敗しました: %w", err)
	}

	emergencyConfig, err := NewEmergencyFundConfig(3, defaultEmergencyFund)
	if err != nil {
		return nil, fmt.Errorf("緊急資金設定の作成に失敗しました: %w", err)
	}

	now := time.Now()

	return &FinancialPlan{
		id:            NewFinancialPlanID(),
		profile:       profile,
		goals:         make([]*entities.Goal, 0),
		emergencyFund: emergencyConfig,
		createdAt:     now,
		updatedAt:     now,
	}, nil
}

// ID は財務計画IDを返す
func (fp *FinancialPlan) ID() FinancialPlanID {
	return fp.id
}

// Profile は財務プロファイルを返す
func (fp *FinancialPlan) Profile() *entities.FinancialProfile {
	return fp.profile
}

// Goals は目標一覧を返す
func (fp *FinancialPlan) Goals() []*entities.Goal {
	return fp.goals
}

// RetirementData は退職データを返す
func (fp *FinancialPlan) RetirementData() *entities.RetirementData {
	return fp.retirementData
}

// EmergencyFund は緊急資金設定を返す
func (fp *FinancialPlan) EmergencyFund() *EmergencyFundConfig {
	return fp.emergencyFund
}

// CreatedAt は作成日時を返す
func (fp *FinancialPlan) CreatedAt() time.Time {
	return fp.createdAt
}

// UpdatedAt は更新日時を返す
func (fp *FinancialPlan) UpdatedAt() time.Time {
	return fp.updatedAt
}

// AddGoal は新しい目標を追加する
func (fp *FinancialPlan) AddGoal(goal *entities.Goal) error {
	if goal == nil {
		return errors.New("目標は必須です")
	}

	// 同じタイプの目標が既に存在するかチェック（退職・緊急資金目標は1つまで）
	if goal.GoalType() == entities.GoalTypeRetirement || goal.GoalType() == entities.GoalTypeEmergency {
		for _, existingGoal := range fp.goals {
			if existingGoal.GoalType() == goal.GoalType() && existingGoal.IsActive() {
				return fmt.Errorf("%sの目標は既に存在します", goal.GoalType().String())
			}
		}
	}

	// 目標の達成可能性をチェック
	achievable, err := goal.IsAchievable(fp.profile)
	if err != nil {
		return fmt.Errorf("目標の達成可能性チェックに失敗しました: %w", err)
	}

	if !achievable {
		return errors.New("現在の財務状況では目標の達成が困難です。目標金額または期日の調整を検討してください")
	}

	fp.goals = append(fp.goals, goal)
	fp.updatedAt = time.Now()
	return nil
}

// RemoveGoal は目標を削除する
func (fp *FinancialPlan) RemoveGoal(goalID entities.GoalID) error {
	for i, goal := range fp.goals {
		if goal.ID() == goalID {
			// スライスから要素を削除
			fp.goals = append(fp.goals[:i], fp.goals[i+1:]...)
			fp.updatedAt = time.Now()
			return nil
		}
	}

	return errors.New("指定された目標が見つかりません")
}

// UpdateProfile は財務プロファイルを更新する
func (fp *FinancialPlan) UpdateProfile(profile *entities.FinancialProfile) error {
	if profile == nil {
		return errors.New("財務プロファイルは必須です")
	}

	fp.profile = profile
	fp.updatedAt = time.Now()
	return nil
}

// SetRetirementData は退職データを設定する
func (fp *FinancialPlan) SetRetirementData(retirementData *entities.RetirementData) error {
	if retirementData == nil {
		return errors.New("退職データは必須です")
	}

	fp.retirementData = retirementData
	fp.updatedAt = time.Now()
	return nil
}

// UpdateEmergencyFund は緊急資金設定を更新する
func (fp *FinancialPlan) UpdateEmergencyFund(config *EmergencyFundConfig) error {
	if config == nil {
		return errors.New("緊急資金設定は必須です")
	}

	fp.emergencyFund = config
	fp.updatedAt = time.Now()
	return nil
}

// GenerateProjection は財務計画の将来予測を生成する
func (fp *FinancialPlan) GenerateProjection(years int) (*PlanProjection, error) {
	if years <= 0 {
		return nil, errors.New("予測年数は正の値である必要があります")
	}

	projection := &PlanProjection{
		GoalProgress: make([]GoalProgress, 0),
	}

	// 資産推移予測
	assetProjections, err := fp.profile.ProjectAssets(years)
	if err != nil {
		return nil, fmt.Errorf("資産推移予測の生成に失敗しました: %w", err)
	}
	projection.AssetProjections = assetProjections

	// 退職資金計算
	if fp.retirementData != nil {
		currentSavings, err := fp.profile.CurrentSavings().Total()
		if err != nil {
			return nil, fmt.Errorf("現在の貯蓄合計の計算に失敗しました: %w", err)
		}

		netSavings, err := fp.profile.CalculateNetSavings()
		if err != nil {
			return nil, fmt.Errorf("純貯蓄額の計算に失敗しました: %w", err)
		}

		retirementCalc, err := fp.retirementData.CalculateRetirementSufficiency(
			currentSavings,
			netSavings,
			fp.profile.InvestmentReturn(),
			fp.profile.InflationRate(),
		)
		if err != nil {
			return nil, fmt.Errorf("退職資金計算に失敗しました: %w", err)
		}
		projection.RetirementCalculation = retirementCalc
	}

	// 緊急資金状況
	if fp.emergencyFund != nil {
		emergencyStatus, err := fp.calculateEmergencyFundStatus()
		if err != nil {
			return nil, fmt.Errorf("緊急資金状況の計算に失敗しました: %w", err)
		}
		projection.EmergencyFundStatus = emergencyStatus
	}

	// 目標進捗
	for _, goal := range fp.goals {
		if !goal.IsActive() {
			continue
		}

		progress, err := goal.CalculateProgress(goal.CurrentAmount())
		if err != nil {
			return nil, fmt.Errorf("目標進捗の計算に失敗しました: %w", err)
		}

		onTrack, message := fp.evaluateGoalProgress(goal)

		projection.GoalProgress = append(projection.GoalProgress, GoalProgress{
			Goal:     goal,
			Progress: progress,
			OnTrack:  onTrack,
			Message:  message,
		})
	}

	return projection, nil
}

// calculateEmergencyFundStatus は緊急資金の状況を計算する
func (fp *FinancialPlan) calculateEmergencyFundStatus() (*EmergencyFundStatus, error) {
	// 月間支出を計算
	monthlyExpenses, err := fp.profile.MonthlyExpenses().Total()
	if err != nil {
		return nil, fmt.Errorf("月間支出の計算に失敗しました: %w", err)
	}

	// 必要緊急資金を計算
	requiredAmount, err := monthlyExpenses.MultiplyByFloat(float64(fp.emergencyFund.TargetMonths))
	if err != nil {
		return nil, fmt.Errorf("必要緊急資金の計算に失敗しました: %w", err)
	}

	// 不足額を計算
	shortfall, err := requiredAmount.Subtract(fp.emergencyFund.CurrentFund)
	if err != nil {
		return nil, fmt.Errorf("緊急資金不足額の計算に失敗しました: %w", err)
	}

	// 不足がない場合は0にする
	if shortfall.IsNegative() {
		shortfall, _ = valueobjects.NewMoneyJPY(0)
	}

	// 目標達成までの月数を計算
	monthsToTarget := 0
	if shortfall.IsPositive() {
		netSavings, err := fp.profile.CalculateNetSavings()
		if err == nil && netSavings.IsPositive() {
			monthsToTarget = int(shortfall.Amount() / netSavings.Amount())
		}
	}

	return &EmergencyFundStatus{
		RequiredAmount: requiredAmount,
		CurrentAmount:  fp.emergencyFund.CurrentFund,
		Shortfall:      shortfall,
		MonthsToTarget: monthsToTarget,
	}, nil
}

// evaluateGoalProgress は目標の進捗を評価する
func (fp *FinancialPlan) evaluateGoalProgress(goal *entities.Goal) (bool, string) {
	// 目標達成可能性をチェック
	achievable, err := goal.IsAchievable(fp.profile)
	if err != nil {
		return false, "進捗評価中にエラーが発生しました"
	}

	if !achievable {
		return false, "現在のペースでは目標達成が困難です"
	}

	// 期限チェック
	if goal.IsOverdue() {
		return false, "目標期限を過ぎています"
	}

	// 完了チェック
	if goal.IsCompleted() {
		return true, "目標を達成しました！"
	}

	// 進捗率チェック
	progress, err := goal.CalculateProgress(goal.CurrentAmount())
	if err != nil {
		return false, "進捗計算中にエラーが発生しました"
	}

	remainingDays := goal.GetRemainingDays()
	if remainingDays <= 0 {
		return false, "目標期限を過ぎています"
	}

	// 期待進捗率を計算（時間ベース）
	totalDays := int(goal.TargetDate().Sub(goal.CreatedAt()).Hours() / 24)
	elapsedDays := totalDays - remainingDays
	expectedProgress := float64(elapsedDays) / float64(totalDays) * 100

	actualProgress := progress.AsPercentage()

	if actualProgress >= expectedProgress {
		return true, "順調に進捗しています"
	} else if actualProgress >= expectedProgress*0.8 {
		return true, "概ね順調です"
	} else {
		return false, "進捗が遅れています。貯蓄額の見直しを検討してください"
	}
}

// ValidatePlan は財務計画全体の妥当性をチェックする
func (fp *FinancialPlan) ValidatePlan() []ValidationError {
	var errors []ValidationError

	// 財務プロファイルの健全性チェック
	if err := fp.profile.ValidateFinancialHealth(); err != nil {
		errors = append(errors, ValidationError{
			Field:   "financial_profile",
			Message: err.Error(),
		})
	}

	// 目標の妥当性チェック
	for i, goal := range fp.goals {
		if !goal.IsActive() {
			continue
		}

		achievable, err := goal.IsAchievable(fp.profile)
		if err != nil {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("goals[%d]", i),
				Message: fmt.Sprintf("目標の達成可能性チェックに失敗しました: %s", err.Error()),
			})
		} else if !achievable {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("goals[%d]", i),
				Message: "現在の財務状況では目標の達成が困難です",
			})
		}
	}

	// 緊急資金の妥当性チェック
	if fp.emergencyFund != nil {
		monthlyExpenses, err := fp.profile.MonthlyExpenses().Total()
		if err == nil {
			requiredAmount, err := monthlyExpenses.MultiplyByFloat(float64(fp.emergencyFund.TargetMonths))
			if err == nil {
				shortfall, err := requiredAmount.Subtract(fp.emergencyFund.CurrentFund)
				if err == nil && shortfall.IsPositive() {
					// 緊急資金が不足している場合の警告
					shortfallRatio := shortfall.Amount() / requiredAmount.Amount()
					if shortfallRatio > 0.5 {
						errors = append(errors, ValidationError{
							Field:   "emergency_fund",
							Message: "緊急資金が大幅に不足しています。目標額の確保を優先してください",
						})
					}
				}
			}
		}
	}

	return errors
}

// GetActiveGoals はアクティブな目標一覧を返す
func (fp *FinancialPlan) GetActiveGoals() []*entities.Goal {
	var activeGoals []*entities.Goal
	for _, goal := range fp.goals {
		if goal.IsActive() {
			activeGoals = append(activeGoals, goal)
		}
	}
	return activeGoals
}

// GetGoalsByType は指定されたタイプの目標一覧を返す
func (fp *FinancialPlan) GetGoalsByType(goalType entities.GoalType) []*entities.Goal {
	var goals []*entities.Goal
	for _, goal := range fp.goals {
		if goal.GoalType() == goalType {
			goals = append(goals, goal)
		}
	}
	return goals
}

// HasRetirementGoal は退職目標が設定されているかどうかを返す
func (fp *FinancialPlan) HasRetirementGoal() bool {
	retirementGoals := fp.GetGoalsByType(entities.GoalTypeRetirement)
	for _, goal := range retirementGoals {
		if goal.IsActive() {
			return true
		}
	}
	return false
}

// HasEmergencyGoal は緊急資金目標が設定されているかどうかを返す
func (fp *FinancialPlan) HasEmergencyGoal() bool {
	emergencyGoals := fp.GetGoalsByType(entities.GoalTypeEmergency)
	for _, goal := range emergencyGoals {
		if goal.IsActive() {
			return true
		}
	}
	return false
}
