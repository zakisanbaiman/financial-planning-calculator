package repositories

import (
	"context"

	"github.com/financial-planning-calculator/backend/domain/entities"
)

// WebAuthnCredentialRepository はWebAuthn認証情報の永続化を担当するリポジトリインターフェース
type WebAuthnCredentialRepository interface {
	// Save は新しいWebAuthn認証情報を保存する
	Save(ctx context.Context, credential *entities.WebAuthnCredential) error

	// FindByID は指定されたIDの認証情報を取得する
	FindByID(ctx context.Context, id entities.CredentialID) (*entities.WebAuthnCredential, error)

	// FindByCredentialID はWebAuthn credential IDから認証情報を取得する
	FindByCredentialID(ctx context.Context, credentialID []byte) (*entities.WebAuthnCredential, error)

	// FindByUserID は指定されたユーザーIDの全ての認証情報を取得する
	FindByUserID(ctx context.Context, userID entities.UserID) ([]*entities.WebAuthnCredential, error)

	// Update は既存の認証情報を更新する
	Update(ctx context.Context, credential *entities.WebAuthnCredential) error

	// Delete は指定されたIDの認証情報を削除する
	Delete(ctx context.Context, id entities.CredentialID) error
}
