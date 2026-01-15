package repositories

import (
	"context"

	"github.com/financial-planning-calculator/backend/domain/entities"
)

// UserRepository はユーザーの永続化を担当するリポジトリインターフェース
type UserRepository interface {
	// Save は新しいユーザーを保存する
	Save(ctx context.Context, user *entities.User) error

	// FindByID は指定されたIDのユーザーを取得する
	FindByID(ctx context.Context, id entities.UserID) (*entities.User, error)

	// FindByEmail はメールアドレスからユーザーを取得する（ログイン時に使用）
	FindByEmail(ctx context.Context, email entities.Email) (*entities.User, error)

	// Update は既存のユーザー情報を更新する
	Update(ctx context.Context, user *entities.User) error

	// Delete は指定されたIDのユーザーを削除する
	Delete(ctx context.Context, id entities.UserID) error

	// Exists は指定されたIDのユーザーが存在するか確認する
	Exists(ctx context.Context, id entities.UserID) (bool, error)

	// ExistsByEmail はメールアドレスが既に使用されているか確認する
	ExistsByEmail(ctx context.Context, email entities.Email) (bool, error)

	// FindByProviderUserID はOAuthプロバイダーのユーザーIDからユーザーを取得する
	FindByProviderUserID(ctx context.Context, provider entities.AuthProvider, providerUserID string) (*entities.User, error)
}
