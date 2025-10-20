package repositories

import (
	"context"

	"github.com/financial-planning-calculator/backend/domain/entities"
)

// GoalRepository は目標の永続化を担当するリポジトリインターフェース
type GoalRepository interface {
	// Save は目標を保存する
	Save(ctx context.Context, goal *entities.Goal) error

	// FindByID は指定されたIDの目標を取得する
	FindByID(ctx context.Context, id entities.GoalID) (*entities.Goal, error)

	// FindByUserID は指定されたユーザーIDの全ての目標を取得する
	FindByUserID(ctx context.Context, userID entities.UserID) ([]*entities.Goal, error)

	// FindActiveGoalsByUserID は指定されたユーザーIDのアクティブな目標を取得する
	FindActiveGoalsByUserID(ctx context.Context, userID entities.UserID) ([]*entities.Goal, error)

	// FindByUserIDAndType は指定されたユーザーIDと目標タイプの目標を取得する
	FindByUserIDAndType(ctx context.Context, userID entities.UserID, goalType entities.GoalType) ([]*entities.Goal, error)

	// Update は既存の目標を更新する
	Update(ctx context.Context, goal *entities.Goal) error

	// Delete は指定されたIDの目標を削除する
	Delete(ctx context.Context, id entities.GoalID) error

	// Exists は指定されたIDの目標が存在するかチェックする
	Exists(ctx context.Context, id entities.GoalID) (bool, error)

	// CountActiveGoalsByType は指定されたユーザーIDと目標タイプのアクティブな目標数を取得する
	CountActiveGoalsByType(ctx context.Context, userID entities.UserID, goalType entities.GoalType) (int, error)
}
