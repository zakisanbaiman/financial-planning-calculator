package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/financial-planning-calculator/backend/domain/entities"
	"github.com/financial-planning-calculator/backend/domain/repositories"
)

// PostgreSQLPasswordResetTokenRepository はPostgreSQLを使ったパスワードリセットトークンリポジトリ
type PostgreSQLPasswordResetTokenRepository struct {
	db *sql.DB
}

// NewPostgreSQLPasswordResetTokenRepository は新しいリポジトリを作成する
func NewPostgreSQLPasswordResetTokenRepository(db *sql.DB) repositories.PasswordResetTokenRepository {
	return &PostgreSQLPasswordResetTokenRepository{db: db}
}

// Save は新しいトークンを保存する
func (r *PostgreSQLPasswordResetTokenRepository) Save(ctx context.Context, token *entities.PasswordResetToken) error {
	query := `
		INSERT INTO password_reset_tokens (id, user_id, token_hash, expires_at, is_used, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query,
		string(token.ID()),
		string(token.UserID()),
		token.TokenHash(),
		token.ExpiresAt(),
		token.IsUsed(),
		token.CreatedAt(),
	)
	if err != nil {
		return fmt.Errorf("パスワードリセットトークンの保存に失敗しました: %w", err)
	}
	return nil
}

// FindByTokenHash はトークンハッシュからトークンを取得する
func (r *PostgreSQLPasswordResetTokenRepository) FindByTokenHash(ctx context.Context, tokenHash string) (*entities.PasswordResetToken, error) {
	query := `
		SELECT id, user_id, token_hash, expires_at, is_used, created_at
		FROM password_reset_tokens
		WHERE token_hash = $1
	`
	row := r.db.QueryRowContext(ctx, query, tokenHash)
	return scanPasswordResetToken(row)
}

// FindByUserID はユーザーIDに紐づくトークン一覧を取得する
func (r *PostgreSQLPasswordResetTokenRepository) FindByUserID(ctx context.Context, userID entities.UserID) ([]*entities.PasswordResetToken, error) {
	query := `
		SELECT id, user_id, token_hash, expires_at, is_used, created_at
		FROM password_reset_tokens
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, string(userID))
	if err != nil {
		return nil, fmt.Errorf("パスワードリセットトークンの取得に失敗しました: %w", err)
	}
	defer rows.Close()

	var tokens []*entities.PasswordResetToken
	for rows.Next() {
		token, err := scanPasswordResetTokenRows(rows)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}
	return tokens, rows.Err()
}

// Update は既存のトークンを更新する
func (r *PostgreSQLPasswordResetTokenRepository) Update(ctx context.Context, token *entities.PasswordResetToken) error {
	query := `
		UPDATE password_reset_tokens
		SET is_used = $1
		WHERE id = $2
	`
	result, err := r.db.ExecContext(ctx, query, token.IsUsed(), string(token.ID()))
	if err != nil {
		return fmt.Errorf("パスワードリセットトークンの更新に失敗しました: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("更新行数の取得に失敗しました: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("パスワードリセットトークンが見つかりません: %s", string(token.ID()))
	}
	return nil
}

// DeleteExpired は期限切れのトークンを全て削除する
func (r *PostgreSQLPasswordResetTokenRepository) DeleteExpired(ctx context.Context) error {
	query := `DELETE FROM password_reset_tokens WHERE expires_at < $1`
	_, err := r.db.ExecContext(ctx, query, time.Now())
	if err != nil {
		return fmt.Errorf("期限切れトークンの削除に失敗しました: %w", err)
	}
	return nil
}

// DeleteByUserID は指定ユーザーのトークンを全て削除する
func (r *PostgreSQLPasswordResetTokenRepository) DeleteByUserID(ctx context.Context, userID entities.UserID) error {
	query := `DELETE FROM password_reset_tokens WHERE user_id = $1`
	_, err := r.db.ExecContext(ctx, query, string(userID))
	if err != nil {
		return fmt.Errorf("ユーザーのトークン削除に失敗しました: %w", err)
	}
	return nil
}

// scanPasswordResetToken は単一行をスキャンしてエンティティを返す
func scanPasswordResetToken(row *sql.Row) (*entities.PasswordResetToken, error) {
	var (
		id        string
		userID    string
		tokenHash string
		expiresAt time.Time
		isUsed    bool
		createdAt time.Time
	)
	err := row.Scan(&id, &userID, &tokenHash, &expiresAt, &isUsed, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("パスワードリセットトークンのスキャンに失敗しました: %w", err)
	}
	return entities.ReconstructPasswordResetToken(id, entities.UserID(userID), tokenHash, expiresAt, isUsed, createdAt), nil
}

// scanPasswordResetTokenRows は複数行から1行をスキャンしてエンティティを返す
func scanPasswordResetTokenRows(rows *sql.Rows) (*entities.PasswordResetToken, error) {
	var (
		id        string
		userID    string
		tokenHash string
		expiresAt time.Time
		isUsed    bool
		createdAt time.Time
	)
	err := rows.Scan(&id, &userID, &tokenHash, &expiresAt, &isUsed, &createdAt)
	if err != nil {
		return nil, fmt.Errorf("パスワードリセットトークンのスキャンに失敗しました: %w", err)
	}
	return entities.ReconstructPasswordResetToken(id, entities.UserID(userID), tokenHash, expiresAt, isUsed, createdAt), nil
}
