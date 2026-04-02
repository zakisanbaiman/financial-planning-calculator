package repositories

import "time"

// キャッシュキーのプレフィックス
// 形式: fp:{リソース種別}:{識別子種別}:{値}
const (
	financialPlanByIDPrefix     = "fp:plan:id:"
	financialPlanByUserIDPrefix = "fp:plan:uid:"
	goalsByUserIDPrefix         = "fp:goals:uid:"
	activeGoalsByUserIDPrefix   = "fp:goals:active:uid:"

	// FinancialPlanTTL は財務計画キャッシュの有効期限
	FinancialPlanTTL = 5 * time.Minute
	// GoalTTL はゴールキャッシュの有効期限
	GoalTTL = 3 * time.Minute
)

func financialPlanByIDKey(id string) string {
	return financialPlanByIDPrefix + id
}

func financialPlanByUserIDKey(userID string) string {
	return financialPlanByUserIDPrefix + userID
}

func goalsByUserIDKey(userID string) string {
	return goalsByUserIDPrefix + userID
}

func activeGoalsByUserIDKey(userID string) string {
	return activeGoalsByUserIDPrefix + userID
}
