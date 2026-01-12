package repositories

import (
	"context"

	"github.com/financial-planning-calculator/backend/domain/entities"
)

// RefreshTokenRepository はリフレッシュトークンの永続化を担当するリポジトリインターフェース
type RefreshTokenRepository interface {
	// Save は新しいリフレッシュトークンを保存する
	Save(ctx context.Context, token *entities.RefreshToken) error

	// FindByTokenHash はトークンハッシュからリフレッシュトークンを取得する
	FindByTokenHash(ctx context.Context, tokenHash string) (*entities.RefreshToken, error)

	// FindByUserID は指定されたユーザーIDの有効なリフレッシュトークンをすべて取得する
	FindByUserID(ctx context.Context, userID entities.UserID) ([]*entities.RefreshToken, error)

	// Update は既存のリフレッシュトークン情報を更新する（最終使用日時、失効状態など）
	Update(ctx context.Context, token *entities.RefreshToken) error

	// Delete は指定されたIDのリフレッシュトークンを削除する
	Delete(ctx context.Context, id entities.RefreshTokenID) error

	// DeleteByUserID は指定されたユーザーIDのすべてのリフレッシュトークンを削除する
	DeleteByUserID(ctx context.Context, userID entities.UserID) error

	// DeleteExpired は期限切れのリフレッシュトークンをすべて削除する（定期的なクリーンアップ用）
	DeleteExpired(ctx context.Context) error

	// RevokeByUserID は指定されたユーザーIDのすべてのリフレッシュトークンを失効させる
	RevokeByUserID(ctx context.Context, userID entities.UserID) error
}
