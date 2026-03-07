package repositories

import (
	"context"

	"github.com/financial-planning-calculator/backend/domain/entities"
)

// PasswordResetTokenRepository はパスワードリセットトークンの永続化を担当するリポジトリインターフェース
type PasswordResetTokenRepository interface {
	// Save は新しいトークンを保存する
	Save(ctx context.Context, token *entities.PasswordResetToken) error

	// FindByTokenHash はトークンハッシュからトークンを取得する
	FindByTokenHash(ctx context.Context, tokenHash string) (*entities.PasswordResetToken, error)

	// FindByUserID はユーザーIDに紐づくトークン一覧を取得する
	FindByUserID(ctx context.Context, userID entities.UserID) ([]*entities.PasswordResetToken, error)

	// Update は既存のトークンを更新する（使用済みフラグの更新などに使用）
	Update(ctx context.Context, token *entities.PasswordResetToken) error

	// DeleteExpired は期限切れのトークンを全て削除する
	DeleteExpired(ctx context.Context) error

	// DeleteByUserID は指定ユーザーのトークンを全て削除する
	DeleteByUserID(ctx context.Context, userID entities.UserID) error
}
