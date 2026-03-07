package entities

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// PasswordResetTokenID はパスワードリセットトークンの一意識別子
type PasswordResetTokenID string

// PasswordResetToken はパスワードリセットトークンエンティティ
type PasswordResetToken struct {
	id        PasswordResetTokenID
	userID    UserID
	tokenHash string
	expiresAt time.Time
	isUsed    bool
	createdAt time.Time
}

// NewPasswordResetToken は新しいパスワードリセットトークンを生成する
// 返値: (エンティティ, 平文トークン, エラー)
func NewPasswordResetToken(userID UserID, expiresAt time.Time) (*PasswordResetToken, string, error) {
	if string(userID) == "" {
		return nil, "", errors.New("ユーザーIDは必須です")
	}

	// 32バイトのランダムトークンを生成
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return nil, "", fmt.Errorf("トークン生成に失敗しました: %w", err)
	}

	plainToken := hex.EncodeToString(randomBytes)
	tokenHash := hashToken(plainToken)

	return &PasswordResetToken{
		id:        PasswordResetTokenID(uuid.New().String()),
		userID:    userID,
		tokenHash: tokenHash,
		expiresAt: expiresAt,
		isUsed:    false,
		createdAt: time.Now(),
	}, plainToken, nil
}

// ReconstructPasswordResetToken はDBから取得したデータからエンティティを再構築する
func ReconstructPasswordResetToken(id string, userID UserID, tokenHash string, expiresAt time.Time, isUsed bool, createdAt time.Time) *PasswordResetToken {
	return &PasswordResetToken{
		id:        PasswordResetTokenID(id),
		userID:    userID,
		tokenHash: tokenHash,
		expiresAt: expiresAt,
		isUsed:    isUsed,
		createdAt: createdAt,
	}
}

// Getters

func (t *PasswordResetToken) ID() PasswordResetTokenID { return t.id }
func (t *PasswordResetToken) UserID() UserID           { return t.userID }
func (t *PasswordResetToken) TokenHash() string        { return t.tokenHash }
func (t *PasswordResetToken) ExpiresAt() time.Time     { return t.expiresAt }
func (t *PasswordResetToken) IsUsed() bool             { return t.isUsed }
func (t *PasswordResetToken) CreatedAt() time.Time     { return t.createdAt }

// IsExpired はトークンが期限切れかどうかを返す
func (t *PasswordResetToken) IsExpired() bool {
	return time.Now().After(t.expiresAt)
}

// IsValid はトークンが有効かどうかを返す（未使用かつ期限内）
func (t *PasswordResetToken) IsValid() bool {
	return !t.isUsed && !t.IsExpired()
}

// VerifyToken は平文トークンがこのエンティティのものと一致するか検証する
func (t *PasswordResetToken) VerifyToken(plainToken string) bool {
	return hashToken(plainToken) == t.tokenHash
}

// Use はトークンを使用済みにする
func (t *PasswordResetToken) Use() {
	t.isUsed = true
}


