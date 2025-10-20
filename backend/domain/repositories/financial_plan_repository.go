package repositories

import (
	"context"

	"github.com/financial-planning-calculator/backend/domain/aggregates"
	"github.com/financial-planning-calculator/backend/domain/entities"
)

// FinancialPlanRepository は財務計画の永続化を担当するリポジトリインターフェース
type FinancialPlanRepository interface {
	// Save は財務計画を保存する
	Save(ctx context.Context, plan *aggregates.FinancialPlan) error

	// FindByID は指定されたIDの財務計画を取得する
	FindByID(ctx context.Context, id aggregates.FinancialPlanID) (*aggregates.FinancialPlan, error)

	// FindByUserID は指定されたユーザーIDの財務計画を取得する
	FindByUserID(ctx context.Context, userID entities.UserID) (*aggregates.FinancialPlan, error)

	// Update は既存の財務計画を更新する
	Update(ctx context.Context, plan *aggregates.FinancialPlan) error

	// Delete は指定されたIDの財務計画を削除する
	Delete(ctx context.Context, id aggregates.FinancialPlanID) error

	// Exists は指定されたIDの財務計画が存在するかチェックする
	Exists(ctx context.Context, id aggregates.FinancialPlanID) (bool, error)

	// ExistsByUserID は指定されたユーザーIDの財務計画が存在するかチェックする
	ExistsByUserID(ctx context.Context, userID entities.UserID) (bool, error)
}
