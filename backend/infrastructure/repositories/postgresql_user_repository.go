package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/financial-planning-calculator/backend/domain/entities"
	"github.com/financial-planning-calculator/backend/domain/repositories"
)

// PostgreSQLUserRepository はPostgreSQLを使用したユーザーリポジトリの実装
type PostgreSQLUserRepository struct {
	db *sql.DB
}

// NewPostgreSQLUserRepository は新しいPostgreSQLユーザーリポジトリを作成する
func NewPostgreSQLUserRepository(db *sql.DB) repositories.UserRepository {
	return &PostgreSQLUserRepository{db: db}
}

// Save は新しいユーザーを保存する
func (r *PostgreSQLUserRepository) Save(ctx context.Context, user *entities.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.ExecContext(ctx, query,
		user.ID().String(),
		user.Email().String(),
		user.PasswordHash().String(),
		user.CreatedAt(),
		user.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("ユーザーの保存に失敗しました: %w", err)
	}

	return nil
}

// FindByID は指定されたIDのユーザーを取得する
func (r *PostgreSQLUserRepository) FindByID(ctx context.Context, id entities.UserID) (*entities.User, error) {
	var userID, email, passwordHash string
	var createdAt, updatedAt time.Time

	query := `SELECT id, email, password_hash, created_at, updated_at FROM users WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(
		&userID, &email, &passwordHash, &createdAt, &updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("ユーザーが見つかりません: %s", id)
		}
		return nil, fmt.Errorf("ユーザーの取得に失敗しました: %w", err)
	}

	return entities.ReconstructUser(userID, email, passwordHash, createdAt, updatedAt)
}

// FindByEmail はメールアドレスからユーザーを取得する
func (r *PostgreSQLUserRepository) FindByEmail(ctx context.Context, email entities.Email) (*entities.User, error) {
	var userID, emailStr, passwordHash string
	var createdAt, updatedAt time.Time

	query := `SELECT id, email, password_hash, created_at, updated_at FROM users WHERE email = $1`
	err := r.db.QueryRowContext(ctx, query, email.String()).Scan(
		&userID, &emailStr, &passwordHash, &createdAt, &updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("ユーザーが見つかりません: %s", email)
		}
		return nil, fmt.Errorf("ユーザーの取得に失敗しました: %w", err)
	}

	return entities.ReconstructUser(userID, emailStr, passwordHash, createdAt, updatedAt)
}

// Update は既存のユーザー情報を更新する
func (r *PostgreSQLUserRepository) Update(ctx context.Context, user *entities.User) error {
	query := `
		UPDATE users 
		SET email = $1, password_hash = $2, updated_at = $3
		WHERE id = $4`

	result, err := r.db.ExecContext(ctx, query,
		user.Email().String(),
		user.PasswordHash().String(),
		user.UpdatedAt(),
		user.ID().String(),
	)
	if err != nil {
		return fmt.Errorf("ユーザーの更新に失敗しました: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("更新結果の確認に失敗しました: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("ユーザーが見つかりません: %s", user.ID())
	}

	return nil
}

// Delete は指定されたIDのユーザーを削除する
func (r *PostgreSQLUserRepository) Delete(ctx context.Context, id entities.UserID) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id.String())
	if err != nil {
		return fmt.Errorf("ユーザーの削除に失敗しました: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("削除結果の確認に失敗しました: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("ユーザーが見つかりません: %s", id)
	}

	return nil
}

// Exists は指定されたIDのユーザーが存在するか確認する
func (r *PostgreSQLUserRepository) Exists(ctx context.Context, id entities.UserID) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`

	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("ユーザーの存在確認に失敗しました: %w", err)
	}

	return exists, nil
}

// ExistsByEmail はメールアドレスが既に使用されているか確認する
func (r *PostgreSQLUserRepository) ExistsByEmail(ctx context.Context, email entities.Email) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	err := r.db.QueryRowContext(ctx, query, email.String()).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("メールアドレスの存在確認に失敗しました: %w", err)
	}

	return exists, nil
}
