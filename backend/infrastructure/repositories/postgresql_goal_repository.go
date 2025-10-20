package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/financial-planning-calculator/backend/domain/entities"
	"github.com/financial-planning-calculator/backend/domain/repositories"
	"github.com/financial-planning-calculator/backend/domain/valueobjects"
)

// PostgreSQLGoalRepository はPostgreSQLを使用した目標リポジトリの実装
type PostgreSQLGoalRepository struct {
	db *sql.DB
}

// NewPostgreSQLGoalRepository は新しいPostgreSQL目標リポジトリを作成する
func NewPostgreSQLGoalRepository(db *sql.DB) repositories.GoalRepository {
	return &PostgreSQLGoalRepository{db: db}
}

// Save は目標を保存する
func (r *PostgreSQLGoalRepository) Save(ctx context.Context, goal *entities.Goal) error {
	query := `
		INSERT INTO goals (id, user_id, type, title, target_amount, target_date, current_amount, monthly_contribution, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := r.db.ExecContext(ctx, query,
		string(goal.ID()),
		string(goal.UserID()),
		string(goal.GoalType()),
		goal.Title(),
		goal.TargetAmount().Amount(),
		goal.TargetDate(),
		goal.CurrentAmount().Amount(),
		goal.MonthlyContribution().Amount(),
		goal.IsActive(),
		goal.CreatedAt(),
		goal.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("目標の保存に失敗しました: %w", err)
	}

	return nil
}

// FindByID は指定されたIDの目標を取得する
func (r *PostgreSQLGoalRepository) FindByID(ctx context.Context, id entities.GoalID) (*entities.Goal, error) {
	var goalID, userID, goalType, title string
	var targetAmount, currentAmount, monthlyContribution float64
	var targetDate time.Time
	var isActive bool
	var createdAt, updatedAt time.Time

	query := `SELECT id, user_id, type, title, target_amount, target_date, current_amount, monthly_contribution, is_active, created_at, updated_at 
			  FROM goals WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, string(id)).Scan(
		&goalID, &userID, &goalType, &title, &targetAmount, &targetDate, &currentAmount, &monthlyContribution, &isActive, &createdAt, &updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("目標が見つかりません: %s", id)
		}
		return nil, fmt.Errorf("目標の取得に失敗しました: %w", err)
	}

	return r.buildGoalFromRow(goalID, userID, goalType, title, targetAmount, currentAmount, monthlyContribution, targetDate, isActive, createdAt, updatedAt)
}

// FindByUserID は指定されたユーザーIDの全ての目標を取得する
func (r *PostgreSQLGoalRepository) FindByUserID(ctx context.Context, userID entities.UserID) ([]*entities.Goal, error) {
	query := `SELECT id, user_id, type, title, target_amount, target_date, current_amount, monthly_contribution, is_active, created_at, updated_at 
			  FROM goals WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, string(userID))
	if err != nil {
		return nil, fmt.Errorf("目標の取得に失敗しました: %w", err)
	}
	defer rows.Close()

	return r.scanGoals(rows)
}

// FindActiveGoalsByUserID は指定されたユーザーIDのアクティブな目標を取得する
func (r *PostgreSQLGoalRepository) FindActiveGoalsByUserID(ctx context.Context, userID entities.UserID) ([]*entities.Goal, error) {
	query := `SELECT id, user_id, type, title, target_amount, target_date, current_amount, monthly_contribution, is_active, created_at, updated_at 
			  FROM goals WHERE user_id = $1 AND is_active = true ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, string(userID))
	if err != nil {
		return nil, fmt.Errorf("アクティブな目標の取得に失敗しました: %w", err)
	}
	defer rows.Close()

	return r.scanGoals(rows)
}

// FindByUserIDAndType は指定されたユーザーIDと目標タイプの目標を取得する
func (r *PostgreSQLGoalRepository) FindByUserIDAndType(ctx context.Context, userID entities.UserID, goalType entities.GoalType) ([]*entities.Goal, error) {
	query := `SELECT id, user_id, type, title, target_amount, target_date, current_amount, monthly_contribution, is_active, created_at, updated_at 
			  FROM goals WHERE user_id = $1 AND type = $2 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, string(userID), string(goalType))
	if err != nil {
		return nil, fmt.Errorf("指定タイプの目標の取得に失敗しました: %w", err)
	}
	defer rows.Close()

	return r.scanGoals(rows)
}

// Update は既存の目標を更新する
func (r *PostgreSQLGoalRepository) Update(ctx context.Context, goal *entities.Goal) error {
	query := `
		UPDATE goals SET 
			type = $2,
			title = $3,
			target_amount = $4,
			target_date = $5,
			current_amount = $6,
			monthly_contribution = $7,
			is_active = $8,
			updated_at = $9
		WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query,
		string(goal.ID()),
		string(goal.GoalType()),
		goal.Title(),
		goal.TargetAmount().Amount(),
		goal.TargetDate(),
		goal.CurrentAmount().Amount(),
		goal.MonthlyContribution().Amount(),
		goal.IsActive(),
		goal.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("目標の更新に失敗しました: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("更新結果の確認に失敗しました: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("更新対象の目標が見つかりません: %s", goal.ID())
	}

	return nil
}

// Delete は指定されたIDの目標を削除する
func (r *PostgreSQLGoalRepository) Delete(ctx context.Context, id entities.GoalID) error {
	query := `DELETE FROM goals WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, string(id))
	if err != nil {
		return fmt.Errorf("目標の削除に失敗しました: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("削除結果の確認に失敗しました: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("削除対象の目標が見つかりません: %s", id)
	}

	return nil
}

// Exists は指定されたIDの目標が存在するかチェックする
func (r *PostgreSQLGoalRepository) Exists(ctx context.Context, id entities.GoalID) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM goals WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, string(id)).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("目標の存在確認に失敗しました: %w", err)
	}
	return count > 0, nil
}

// CountActiveGoalsByType は指定されたユーザーIDと目標タイプのアクティブな目標数を取得する
func (r *PostgreSQLGoalRepository) CountActiveGoalsByType(ctx context.Context, userID entities.UserID, goalType entities.GoalType) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM goals WHERE user_id = $1 AND type = $2 AND is_active = true`
	err := r.db.QueryRowContext(ctx, query, string(userID), string(goalType)).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("アクティブな目標数の取得に失敗しました: %w", err)
	}
	return count, nil
}

// scanGoals は複数の目標をスキャンする
func (r *PostgreSQLGoalRepository) scanGoals(rows *sql.Rows) ([]*entities.Goal, error) {
	var goals []*entities.Goal

	for rows.Next() {
		var goalID, userID, goalType, title string
		var targetAmount, currentAmount, monthlyContribution float64
		var targetDate time.Time
		var isActive bool
		var createdAt, updatedAt time.Time

		if err := rows.Scan(&goalID, &userID, &goalType, &title, &targetAmount, &targetDate, &currentAmount, &monthlyContribution, &isActive, &createdAt, &updatedAt); err != nil {
			return nil, fmt.Errorf("目標の読み取りに失敗しました: %w", err)
		}

		goal, err := r.buildGoalFromRow(goalID, userID, goalType, title, targetAmount, currentAmount, monthlyContribution, targetDate, isActive, createdAt, updatedAt)
		if err != nil {
			return nil, err
		}

		goals = append(goals, goal)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("目標の読み取り中にエラーが発生しました: %w", err)
	}

	return goals, nil
}

// buildGoalFromRow は行データから目標エンティティを構築する
func (r *PostgreSQLGoalRepository) buildGoalFromRow(
	goalID, userID, goalType, title string,
	targetAmount, currentAmount, monthlyContribution float64,
	targetDate time.Time,
	isActive bool,
	createdAt, updatedAt time.Time,
) (*entities.Goal, error) {
	// 値オブジェクトを作成
	targetAmountVO, err := valueobjects.NewMoneyJPY(targetAmount)
	if err != nil {
		return nil, fmt.Errorf("目標金額の作成に失敗しました: %w", err)
	}

	monthlyContributionVO, err := valueobjects.NewMoneyJPY(monthlyContribution)
	if err != nil {
		return nil, fmt.Errorf("月間拠出額の作成に失敗しました: %w", err)
	}

	// 目標を作成
	goal, err := entities.NewGoalWithID(
		entities.GoalID(goalID),
		entities.UserID(userID),
		entities.GoalType(goalType),
		title,
		targetAmountVO,
		targetDate,
		monthlyContributionVO,
		createdAt,
		updatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("目標の作成に失敗しました: %w", err)
	}

	// 現在の金額を設定
	currentAmountVO, err := valueobjects.NewMoneyJPY(currentAmount)
	if err != nil {
		return nil, fmt.Errorf("現在の金額の作成に失敗しました: %w", err)
	}
	if err := goal.UpdateCurrentAmount(currentAmountVO); err != nil {
		return nil, fmt.Errorf("現在の金額の設定に失敗しました: %w", err)
	}

	// アクティブ状態を設定
	if !isActive {
		goal.Deactivate()
	}

	return goal, nil
}
