package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/financial-planning-calculator/backend/domain/entities"
	"github.com/financial-planning-calculator/backend/domain/repositories"
)

// PostgreSQLRefreshTokenRepository はPostgreSQLを使用したリフレッシュトークンリポジトリの実装
type PostgreSQLRefreshTokenRepository struct {
	db *sql.DB
}

// NewPostgreSQLRefreshTokenRepository は新しいPostgreSQLリフレッシュトークンリポジトリを作成する
func NewPostgreSQLRefreshTokenRepository(db *sql.DB) repositories.RefreshTokenRepository {
	return &PostgreSQLRefreshTokenRepository{db: db}
}

// Save は新しいリフレッシュトークンを保存する
func (r *PostgreSQLRefreshTokenRepository) Save(ctx context.Context, token *entities.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, is_revoked, created_at, last_used_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.db.ExecContext(ctx, query,
		token.ID().String(),
		token.UserID().String(),
		token.TokenHash(),
		token.ExpiresAt(),
		token.IsRevoked(),
		token.CreatedAt(),
		token.LastUsedAt(),
	)
	if err != nil {
		return fmt.Errorf("リフレッシュトークンの保存に失敗しました: %w", err)
	}

	return nil
}

// FindByTokenHash はトークンハッシュからリフレッシュトークンを取得する
func (r *PostgreSQLRefreshTokenRepository) FindByTokenHash(ctx context.Context, tokenHash string) (*entities.RefreshToken, error) {
	var id, userID, storedTokenHash string
	var expiresAt, createdAt, lastUsedAt time.Time
	var isRevoked bool

	query := `
		SELECT id, user_id, token_hash, expires_at, is_revoked, created_at, last_used_at
		FROM refresh_tokens
		WHERE token_hash = $1`

	err := r.db.QueryRowContext(ctx, query, tokenHash).Scan(
		&id, &userID, &storedTokenHash, &expiresAt, &isRevoked, &createdAt, &lastUsedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("リフレッシュトークンが見つかりません")
		}
		return nil, fmt.Errorf("リフレッシュトークンの取得に失敗しました: %w", err)
	}

	userIDEntity, err := entities.NewUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("ユーザーIDの変換に失敗しました: %w", err)
	}

	return entities.ReconstructRefreshToken(id, userIDEntity, storedTokenHash, expiresAt, isRevoked, createdAt, lastUsedAt), nil
}

// FindByUserID は指定されたユーザーIDの有効なリフレッシュトークンをすべて取得する
func (r *PostgreSQLRefreshTokenRepository) FindByUserID(ctx context.Context, userID entities.UserID) ([]*entities.RefreshToken, error) {
	query := `
		SELECT id, user_id, token_hash, expires_at, is_revoked, created_at, last_used_at
		FROM refresh_tokens
		WHERE user_id = $1 AND is_revoked = false AND expires_at > NOW()
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, userID.String())
	if err != nil {
		return nil, fmt.Errorf("リフレッシュトークンの取得に失敗しました: %w", err)
	}
	defer rows.Close()

	var tokens []*entities.RefreshToken
	for rows.Next() {
		var id, userIDStr, tokenHash string
		var expiresAt, createdAt, lastUsedAt time.Time
		var isRevoked bool

		if err := rows.Scan(&id, &userIDStr, &tokenHash, &expiresAt, &isRevoked, &createdAt, &lastUsedAt); err != nil {
			return nil, fmt.Errorf("リフレッシュトークンのスキャンに失敗しました: %w", err)
		}

		userIDEntity, err := entities.NewUserID(userIDStr)
		if err != nil {
			return nil, fmt.Errorf("ユーザーIDの変換に失敗しました: %w", err)
		}

		tokens = append(tokens, entities.ReconstructRefreshToken(id, userIDEntity, tokenHash, expiresAt, isRevoked, createdAt, lastUsedAt))
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("リフレッシュトークンの取得中にエラーが発生しました: %w", err)
	}

	return tokens, nil
}

// Update は既存のリフレッシュトークン情報を更新する
func (r *PostgreSQLRefreshTokenRepository) Update(ctx context.Context, token *entities.RefreshToken) error {
	query := `
		UPDATE refresh_tokens
		SET is_revoked = $1, last_used_at = $2
		WHERE id = $3`

	result, err := r.db.ExecContext(ctx, query, token.IsRevoked(), token.LastUsedAt(), token.ID().String())
	if err != nil {
		return fmt.Errorf("リフレッシュトークンの更新に失敗しました: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("更新結果の確認に失敗しました: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("リフレッシュトークンが見つかりません: %s", token.ID())
	}

	return nil
}

// Delete は指定されたIDのリフレッシュトークンを削除する
func (r *PostgreSQLRefreshTokenRepository) Delete(ctx context.Context, id entities.RefreshTokenID) error {
	query := `DELETE FROM refresh_tokens WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id.String())
	if err != nil {
		return fmt.Errorf("リフレッシュトークンの削除に失敗しました: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("削除結果の確認に失敗しました: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("リフレッシュトークンが見つかりません: %s", id)
	}

	return nil
}

// DeleteByUserID は指定されたユーザーIDのすべてのリフレッシュトークンを削除する
func (r *PostgreSQLRefreshTokenRepository) DeleteByUserID(ctx context.Context, userID entities.UserID) error {
	query := `DELETE FROM refresh_tokens WHERE user_id = $1`

	_, err := r.db.ExecContext(ctx, query, userID.String())
	if err != nil {
		return fmt.Errorf("リフレッシュトークンの削除に失敗しました: %w", err)
	}

	return nil
}

// DeleteExpired は期限切れのリフレッシュトークンをすべて削除する
func (r *PostgreSQLRefreshTokenRepository) DeleteExpired(ctx context.Context) error {
	query := `DELETE FROM refresh_tokens WHERE expires_at < NOW()`

	result, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("期限切れリフレッシュトークンの削除に失敗しました: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("削除結果の確認に失敗しました: %w", err)
	}

	// ログ出力用に削除された行数を返すことも可能
	_ = rowsAffected

	return nil
}

// RevokeByUserID は指定されたユーザーIDのすべてのリフレッシュトークンを失効させる
func (r *PostgreSQLRefreshTokenRepository) RevokeByUserID(ctx context.Context, userID entities.UserID) error {
	query := `UPDATE refresh_tokens SET is_revoked = true WHERE user_id = $1 AND is_revoked = false`

	_, err := r.db.ExecContext(ctx, query, userID.String())
	if err != nil {
		return fmt.Errorf("リフレッシュトークンの失効に失敗しました: %w", err)
	}

	return nil
}
