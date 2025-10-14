package entities

import (
	"errors"
	"fmt"
	"time"

	"financial-planning-calculator/domain/valueobjects"

	"github.com/google/uuid"
)

// GoalID は目標の一意識別子
type GoalID string

// NewGoalID は新しい目標IDを生成する
func NewGoalID() GoalID {
	return GoalID(uuid.New().String())
}

// GoalType は目標の種類を表す
type GoalType string

const (
	GoalTypeSavings    GoalType = "savings"    // 一般的な貯蓄目標
	GoalTypeRetirement GoalType = "retirement" // 退職・老後資金目標
	GoalTypeEmergency  GoalType = "emergency"  // 緊急資金目標
	GoalTypeCustom     GoalType = "custom"     // カスタム目標
)

// IsValid はGoalTypeが有効かどうかを確認する
func (gt GoalType) IsValid() bool {
	switch gt {
	case GoalTypeSavings, GoalTypeRetirement, GoalTypeEmergency, GoalTypeCustom:
		return true
	default:
		return false
	}
}

// String はGoalTypeの文字列表現を返す
func (gt GoalType) String() string {
	switch gt {
	case GoalTypeSavings:
		return "貯蓄目標"
	case GoalTypeRetirement:
		return "退職・老後資金目標"
	case GoalTypeEmergency:
		return "緊急資金目標"
	case GoalTypeCustom:
		return "カスタム目標"
	default:
		return "不明な目標タイプ"
	}
}

// ProgressRate は進捗率を表す値オブジェクト
type ProgressRate struct {
	rate valueobjects.Rate
}

// NewProgressRate は新しい進捗率を作成する
func NewProgressRate(percentage float64) (ProgressRate, error) {
	if percentage < 0 {
		percentage = 0
	}
	if percentage > 100 {
		percentage = 100
	}

	rate, err := valueobjects.NewRate(percentage)
	if err != nil {
		return ProgressRate{}, fmt.Errorf("進捗率の作成に失敗しました: %w", err)
	}

	return ProgressRate{rate: rate}, nil
}

// AsPercentage は進捗率をパーセンテージで返す
func (pr ProgressRate) AsPercentage() float64 {
	return pr.rate.AsPercentage()
}

// IsComplete は目標が完了しているかどうかを返す
func (pr ProgressRate) IsComplete() bool {
	return pr.rate.AsPercentage() >= 100.0
}

// String は進捗率の文字列表現を返す
func (pr ProgressRate) String() string {
	return fmt.Sprintf("%.1f%%", pr.rate.AsPercentage())
}

// GoalAdjustment は目標調整の提案を表す
type GoalAdjustment struct {
	Type        string      `json:"type"`        // "amount", "date", "contribution"
	Description string      `json:"description"` // 調整内容の説明
	NewValue    interface{} `json:"new_value"`   // 新しい値
	Reason      string      `json:"reason"`      // 調整理由
}

// Goal は財務目標を表すエンティティ
type Goal struct {
	id                  GoalID
	userID              UserID
	goalType            GoalType
	title               string
	targetAmount        valueobjects.Money
	targetDate          time.Time
	currentAmount       valueobjects.Money
	monthlyContribution valueobjects.Money
	isActive            bool
	createdAt           time.Time
	updatedAt           time.Time
}

// NewGoal は新しい目標を作成する
func NewGoal(
	userID UserID,
	goalType GoalType,
	title string,
	targetAmount valueobjects.Money,
	targetDate time.Time,
	monthlyContribution valueobjects.Money,
) (*Goal, error) {
	if userID == "" {
		return nil, errors.New("ユーザーIDは必須です")
	}

	if !goalType.IsValid() {
		return nil, errors.New("無効な目標タイプです")
	}

	if title == "" {
		return nil, errors.New("目標タイトルは必須です")
	}

	if !targetAmount.IsPositive() {
		return nil, errors.New("目標金額は正の値である必要があります")
	}

	if targetDate.Before(time.Now()) {
		return nil, errors.New("目標日は未来の日付である必要があります")
	}

	if monthlyContribution.IsNegative() {
		return nil, errors.New("月間拠出額は負の値にできません")
	}

	currentAmount, err := valueobjects.NewMoneyJPY(0)
	if err != nil {
		return nil, fmt.Errorf("初期金額の設定に失敗しました: %w", err)
	}

	now := time.Now()

	return &Goal{
		id:                  NewGoalID(),
		userID:              userID,
		goalType:            goalType,
		title:               title,
		targetAmount:        targetAmount,
		targetDate:          targetDate,
		currentAmount:       currentAmount,
		monthlyContribution: monthlyContribution,
		isActive:            true,
		createdAt:           now,
		updatedAt:           now,
	}, nil
}

// ID は目標IDを返す
func (g *Goal) ID() GoalID {
	return g.id
}

// UserID はユーザーIDを返す
func (g *Goal) UserID() UserID {
	return g.userID
}

// GoalType は目標タイプを返す
func (g *Goal) GoalType() GoalType {
	return g.goalType
}

// Title は目標タイトルを返す
func (g *Goal) Title() string {
	return g.title
}

// TargetAmount は目標金額を返す
func (g *Goal) TargetAmount() valueobjects.Money {
	return g.targetAmount
}

// TargetDate は目標日を返す
func (g *Goal) TargetDate() time.Time {
	return g.targetDate
}

// CurrentAmount は現在の金額を返す
func (g *Goal) CurrentAmount() valueobjects.Money {
	return g.currentAmount
}

// MonthlyContribution は月間拠出額を返す
func (g *Goal) MonthlyContribution() valueobjects.Money {
	return g.monthlyContribution
}

// IsActive は目標がアクティブかどうかを返す
func (g *Goal) IsActive() bool {
	return g.isActive
}

// CreatedAt は作成日時を返す
func (g *Goal) CreatedAt() time.Time {
	return g.createdAt
}

// UpdatedAt は更新日時を返す
func (g *Goal) UpdatedAt() time.Time {
	return g.updatedAt
}

// CalculateProgress は現在の進捗率を計算する
func (g *Goal) CalculateProgress(currentAmount valueobjects.Money) (ProgressRate, error) {
	if g.targetAmount.IsZero() {
		return NewProgressRate(100.0) // 目標金額が0の場合は100%とする
	}

	// 進捗率 = (現在の金額 / 目標金額) * 100
	progressDecimal := currentAmount.Amount() / g.targetAmount.Amount()
	progressPercentage := progressDecimal * 100

	return NewProgressRate(progressPercentage)
}

// EstimateCompletionDate は月間貯蓄額に基づいて完了予定日を推定する
func (g *Goal) EstimateCompletionDate(monthlySavings valueobjects.Money) (time.Time, error) {
	if monthlySavings.IsZero() || monthlySavings.IsNegative() {
		return time.Time{}, errors.New("月間貯蓄額は正の値である必要があります")
	}

	// 残り必要金額を計算
	remainingAmount, err := g.targetAmount.Subtract(g.currentAmount)
	if err != nil {
		return time.Time{}, fmt.Errorf("残り必要金額の計算に失敗しました: %w", err)
	}

	// 既に目標達成している場合
	if remainingAmount.IsZero() || remainingAmount.IsNegative() {
		return time.Now(), nil
	}

	// 必要な月数を計算
	monthsNeeded := remainingAmount.Amount() / monthlySavings.Amount()

	// 完了予定日を計算
	completionDate := time.Now().AddDate(0, int(monthsNeeded), 0)

	return completionDate, nil
}

// IsAchievable は財務プロファイルに基づいて目標が達成可能かどうかを判定する
func (g *Goal) IsAchievable(financialProfile *FinancialProfile) (bool, error) {
	if financialProfile == nil {
		return false, errors.New("財務プロファイルが必要です")
	}

	// 純貯蓄額を計算
	netSavings, err := financialProfile.CalculateNetSavings()
	if err != nil {
		return false, fmt.Errorf("純貯蓄額の計算に失敗しました: %w", err)
	}

	// 純貯蓄額が負の場合は達成不可能
	if netSavings.IsNegative() || netSavings.IsZero() {
		return false, nil
	}

	// 目標日までの期間を計算
	now := time.Now()
	if g.targetDate.Before(now) {
		return false, nil // 目標日が過去の場合は達成不可能
	}

	monthsUntilTarget := int(g.targetDate.Sub(now).Hours() / (24 * 30)) // 概算の月数

	if monthsUntilTarget <= 0 {
		return false, nil
	}

	// 残り必要金額を計算
	remainingAmount, err := g.targetAmount.Subtract(g.currentAmount)
	if err != nil {
		return false, fmt.Errorf("残り必要金額の計算に失敗しました: %w", err)
	}

	// 既に達成している場合
	if remainingAmount.IsZero() || remainingAmount.IsNegative() {
		return true, nil
	}

	// 必要な月間貯蓄額を計算
	requiredMonthlySavings := remainingAmount.Amount() / float64(monthsUntilTarget)

	// 現在の純貯蓄額で達成可能かチェック
	return netSavings.Amount() >= requiredMonthlySavings, nil
}

// UpdateCurrentAmount は現在の金額を更新する
func (g *Goal) UpdateCurrentAmount(newAmount valueobjects.Money) error {
	if newAmount.IsNegative() {
		return errors.New("現在の金額は負の値にできません")
	}

	g.currentAmount = newAmount
	g.updatedAt = time.Now()
	return nil
}

// UpdateMonthlyContribution は月間拠出額を更新する
func (g *Goal) UpdateMonthlyContribution(newContribution valueobjects.Money) error {
	if newContribution.IsNegative() {
		return errors.New("月間拠出額は負の値にできません")
	}

	g.monthlyContribution = newContribution
	g.updatedAt = time.Now()
	return nil
}

// UpdateTargetAmount は目標金額を更新する
func (g *Goal) UpdateTargetAmount(newAmount valueobjects.Money) error {
	if !newAmount.IsPositive() {
		return errors.New("目標金額は正の値である必要があります")
	}

	g.targetAmount = newAmount
	g.updatedAt = time.Now()
	return nil
}

// UpdateTargetDate は目標日を更新する
func (g *Goal) UpdateTargetDate(newDate time.Time) error {
	if newDate.Before(time.Now()) {
		return errors.New("目標日は未来の日付である必要があります")
	}

	g.targetDate = newDate
	g.updatedAt = time.Now()
	return nil
}

// UpdateTitle は目標タイトルを更新する
func (g *Goal) UpdateTitle(newTitle string) error {
	if newTitle == "" {
		return errors.New("目標タイトルは必須です")
	}

	g.title = newTitle
	g.updatedAt = time.Now()
	return nil
}

// Activate は目標をアクティブにする
func (g *Goal) Activate() {
	g.isActive = true
	g.updatedAt = time.Now()
}

// Deactivate は目標を非アクティブにする
func (g *Goal) Deactivate() {
	g.isActive = false
	g.updatedAt = time.Now()
}

// IsOverdue は目標が期限切れかどうかを返す
func (g *Goal) IsOverdue() bool {
	return time.Now().After(g.targetDate) && !g.IsCompleted()
}

// IsCompleted は目標が完了しているかどうかを返す
func (g *Goal) IsCompleted() bool {
	isGreaterOrEqual, err := g.currentAmount.GreaterThan(g.targetAmount)
	if err != nil {
		return false
	}

	isEqual, err := g.currentAmount.Equal(g.targetAmount)
	if err != nil {
		return false
	}

	return isGreaterOrEqual || isEqual
}

// GetRemainingAmount は残り必要金額を返す
func (g *Goal) GetRemainingAmount() (valueobjects.Money, error) {
	if g.IsCompleted() {
		return valueobjects.NewMoneyJPY(0)
	}

	return g.targetAmount.Subtract(g.currentAmount)
}

// GetRemainingDays は目標日までの残り日数を返す
func (g *Goal) GetRemainingDays() int {
	if g.targetDate.Before(time.Now()) {
		return 0
	}

	duration := time.Until(g.targetDate)
	return int(duration.Hours() / 24)
}

// CalculateRequiredMonthlySavings は目標達成に必要な月間貯蓄額を計算する
func (g *Goal) CalculateRequiredMonthlySavings() (valueobjects.Money, error) {
	remainingAmount, err := g.GetRemainingAmount()
	if err != nil {
		return valueobjects.Money{}, fmt.Errorf("残り必要金額の計算に失敗しました: %w", err)
	}

	if remainingAmount.IsZero() || remainingAmount.IsNegative() {
		return valueobjects.NewMoneyJPY(0)
	}

	remainingDays := g.GetRemainingDays()
	if remainingDays <= 0 {
		return remainingAmount, nil // 期限が過ぎている場合は全額必要
	}

	remainingMonths := float64(remainingDays) / 30.0 // 概算の月数
	if remainingMonths < 1 {
		remainingMonths = 1 // 最低1ヶ月とする
	}

	requiredMonthlySavings := remainingAmount.Amount() / remainingMonths

	return valueobjects.NewMoneyJPY(requiredMonthlySavings)
}
