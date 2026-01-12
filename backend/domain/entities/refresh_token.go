package entities

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// RefreshTokenID はリフレッシュトークンの一意識別子
type RefreshTokenID string

// NewRefreshTokenID は新しいリフレッシュトークンIDを生成する
func NewRefreshTokenID() RefreshTokenID {
	return RefreshTokenID(uuid.New().String())
}

// String はRefreshTokenIDの文字列表現を返す
func (id RefreshTokenID) String() string {
	return string(id)
}

// RefreshToken はJWTトークン更新用のリフレッシュトークンエンティティ
type RefreshToken struct {
	id         RefreshTokenID
	userID     UserID
	tokenHash  string
	expiresAt  time.Time
	isRevoked  bool
	createdAt  time.Time
	lastUsedAt time.Time
}

// NewRefreshToken は新しいリフレッシュトークンを生成する
// token: 平文のランダムトークン（クライアントに返却される）
// userID: トークンを所有するユーザーID
// expiresAt: トークンの有効期限
func NewRefreshToken(userID UserID, expiresAt time.Time) (*RefreshToken, string, error) {
	if userID == "" {
		return nil, "", errors.New("ユーザーIDは必須です")
	}

	if expiresAt.Before(time.Now()) {
		return nil, "", errors.New("有効期限は未来の日時である必要があります")
	}

	// ランダムトークンを生成（32バイト = 256ビット）
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, "", fmt.Errorf("トークン生成に失敗しました: %w", err)
	}
	token := hex.EncodeToString(tokenBytes)

	// トークンをハッシュ化してDBに保存
	tokenHash := hashToken(token)

	now := time.Now()
	refreshToken := &RefreshToken{
		id:         NewRefreshTokenID(),
		userID:     userID,
		tokenHash:  tokenHash,
		expiresAt:  expiresAt,
		isRevoked:  false,
		createdAt:  now,
		lastUsedAt: now,
	}

	return refreshToken, token, nil
}

// ReconstructRefreshToken は既存のデータからリフレッシュトークンを再構築する（リポジトリからの取得用）
func ReconstructRefreshToken(
	id string,
	userID UserID,
	tokenHash string,
	expiresAt time.Time,
	isRevoked bool,
	createdAt time.Time,
	lastUsedAt time.Time,
) *RefreshToken {
	return &RefreshToken{
		id:         RefreshTokenID(id),
		userID:     userID,
		tokenHash:  tokenHash,
		expiresAt:  expiresAt,
		isRevoked:  isRevoked,
		createdAt:  createdAt,
		lastUsedAt: lastUsedAt,
	}
}

// hashToken はトークンをSHA-256でハッシュ化する
func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// ID はリフレッシュトークンのIDを返す
func (rt *RefreshToken) ID() RefreshTokenID {
	return rt.id
}

// UserID はトークンを所有するユーザーIDを返す
func (rt *RefreshToken) UserID() UserID {
	return rt.userID
}

// TokenHash はトークンのハッシュ値を返す
func (rt *RefreshToken) TokenHash() string {
	return rt.tokenHash
}

// ExpiresAt はトークンの有効期限を返す
func (rt *RefreshToken) ExpiresAt() time.Time {
	return rt.expiresAt
}

// IsRevoked はトークンが失効されているかを返す
func (rt *RefreshToken) IsRevoked() bool {
	return rt.isRevoked
}

// CreatedAt はトークンの作成日時を返す
func (rt *RefreshToken) CreatedAt() time.Time {
	return rt.createdAt
}

// LastUsedAt はトークンの最終使用日時を返す
func (rt *RefreshToken) LastUsedAt() time.Time {
	return rt.lastUsedAt
}

// IsExpired はトークンが期限切れかどうかを確認する
func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.expiresAt)
}

// IsValid はトークンが有効かどうかを確認する（期限切れでなく、失効されていない）
func (rt *RefreshToken) IsValid() bool {
	return !rt.IsExpired() && !rt.isRevoked
}

// VerifyToken は提供されたトークンがこのリフレッシュトークンと一致するか検証する
func (rt *RefreshToken) VerifyToken(token string) bool {
	return rt.tokenHash == hashToken(token)
}

// Revoke はトークンを失効させる（ログアウト時などに使用）
func (rt *RefreshToken) Revoke() {
	rt.isRevoked = true
}

// UpdateLastUsedAt はトークンの最終使用日時を更新する
func (rt *RefreshToken) UpdateLastUsedAt() {
	rt.lastUsedAt = time.Now()
}
