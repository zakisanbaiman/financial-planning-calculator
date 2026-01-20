package entities

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// UserID はユーザーの一意識別子
type UserID string

// NewUserID は既存のIDからUserIDを生成する
func NewUserID(id string) (UserID, error) {
	if id == "" {
		return "", errors.New("ユーザーIDは必須です")
	}
	return UserID(id), nil
}

// String はUserIDの文字列表現を返す
func (uid UserID) String() string {
	return string(uid)
}

// Email はメールアドレスを表す値オブジェクト
type Email string

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// NewEmail は新しいEmailを作成する
func NewEmail(email string) (Email, error) {
	if email == "" {
		return "", errors.New("メールアドレスは必須です")
	}
	if !emailRegex.MatchString(email) {
		return "", errors.New("無効なメールアドレス形式です")
	}
	return Email(email), nil
}

// String はEmailの文字列表現を返す
func (e Email) String() string {
	return string(e)
}

// PasswordHash はハッシュ化されたパスワードを表す値オブジェクト
type PasswordHash string

// NewPasswordHash は平文パスワードからハッシュを生成する
func NewPasswordHash(plainPassword string) (PasswordHash, error) {
	if plainPassword == "" {
		return "", errors.New("パスワードは必須です")
	}
	if len(plainPassword) < 8 {
		return "", errors.New("パスワードは8文字以上である必要があります")
	}

	// bcryptでハッシュ化（コスト: 10）
	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("パスワードのハッシュ化に失敗しました: %w", err)
	}

	return PasswordHash(hash), nil
}

// NewPasswordHashFromHash は既存のハッシュからPasswordHashを生成する（DB読み込み用）
func NewPasswordHashFromHash(hash string) PasswordHash {
	return PasswordHash(hash)
}

// Compare は平文パスワードとハッシュを比較する
func (ph PasswordHash) Compare(plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(ph), []byte(plainPassword))
}

// String はPasswordHashの文字列表現を返す
func (ph PasswordHash) String() string {
	return string(ph)
}

// AuthProvider は認証プロバイダーを表す
type AuthProvider string

const (
	AuthProviderLocal  AuthProvider = "local"
	AuthProviderGitHub AuthProvider = "github"
	AuthProviderGoogle AuthProvider = "google"
)

// User はユーザーエンティティ
type User struct {
	id                   UserID
	email                Email
	passwordHash         PasswordHash
	provider             AuthProvider
	providerUserID       string
	name                 string
	avatarURL            string
	emailVerified        bool
	emailVerifiedAt      *time.Time
	twoFactorEnabled     bool
	twoFactorSecret      string
	twoFactorBackupCodes []string
	createdAt            time.Time
	updatedAt            time.Time
}

// NewUser は新しいユーザーを作成する（新規登録用）
func NewUser(id string, email string, plainPassword string) (*User, error) {
	// バリデーション
	userID, err := NewUserID(id)
	if err != nil {
		return nil, err
	}

	emailVO, err := NewEmail(email)
	if err != nil {
		return nil, err
	}

	passwordHash, err := NewPasswordHash(plainPassword)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	return &User{
		id:               userID,
		email:            emailVO,
		passwordHash:     passwordHash,
		provider:         AuthProviderLocal,
		emailVerified:    false, // Local users need to verify their email
		twoFactorEnabled: false,
		createdAt:        now,
		updatedAt:        now,
	}, nil
}

// ReconstructUser はDBから取得したデータからUserを再構築する（リポジトリ用）
func ReconstructUser(id string, email string, passwordHash string, emailVerified bool, emailVerifiedAt *time.Time, twoFactorEnabled bool, twoFactorSecret string, twoFactorBackupCodes []string, createdAt, updatedAt time.Time) (*User, error) {
	userID, err := NewUserID(id)
	if err != nil {
		return nil, err
	}

	emailVO, err := NewEmail(email)
	if err != nil {
		return nil, err
	}

	return &User{
		id:                   userID,
		email:                emailVO,
		passwordHash:         NewPasswordHashFromHash(passwordHash),
		provider:             AuthProviderLocal,
		emailVerified:        emailVerified,
		emailVerifiedAt:      emailVerifiedAt,
		twoFactorEnabled:     twoFactorEnabled,
		twoFactorSecret:      twoFactorSecret,
		twoFactorBackupCodes: twoFactorBackupCodes,
		createdAt:            createdAt,
		updatedAt:            updatedAt,
	}, nil
}

// ReconstructUserWithOAuth はDBから取得したOAuthユーザーデータからUserを再構築する
func ReconstructUserWithOAuth(id string, email string, passwordHash string, provider string, providerUserID string, name string, avatarURL string, emailVerified bool, emailVerifiedAt *time.Time, twoFactorEnabled bool, twoFactorSecret string, twoFactorBackupCodes []string, createdAt, updatedAt time.Time) (*User, error) {
	userID, err := NewUserID(id)
	if err != nil {
		return nil, err
	}

	emailVO, err := NewEmail(email)
	if err != nil {
		return nil, err
	}

	var pwdHash PasswordHash
	if passwordHash != "" {
		pwdHash = NewPasswordHashFromHash(passwordHash)
	}

	return &User{
		id:                   userID,
		email:                emailVO,
		passwordHash:         pwdHash,
		provider:             AuthProvider(provider),
		providerUserID:       providerUserID,
		name:                 name,
		avatarURL:            avatarURL,
		emailVerified:        emailVerified,
		emailVerifiedAt:      emailVerifiedAt,
		twoFactorEnabled:     twoFactorEnabled,
		twoFactorSecret:      twoFactorSecret,
		twoFactorBackupCodes: twoFactorBackupCodes,
		createdAt:            createdAt,
		updatedAt:            updatedAt,
	}, nil
}

// NewOAuthUser はOAuth認証で新しいユーザーを作成する
func NewOAuthUser(id string, email string, provider AuthProvider, providerUserID string, name string, avatarURL string) (*User, error) {
	userID, err := NewUserID(id)
	if err != nil {
		return nil, err
	}

	emailVO, err := NewEmail(email)
	if err != nil {
		return nil, err
	}

	if providerUserID == "" {
		return nil, errors.New("プロバイダーユーザーIDは必須です")
	}

	now := time.Now()

	return &User{
		id:               userID,
		email:            emailVO,
		provider:         provider,
		providerUserID:   providerUserID,
		name:             name,
		avatarURL:        avatarURL,
		emailVerified:    true, // OAuth providers are trusted for email verification
		emailVerifiedAt:  &now,
		twoFactorEnabled: false,
		createdAt:        now,
		updatedAt:        now,
	}, nil
}

// Getters

// ID はユーザーIDを返す
func (u *User) ID() UserID {
	return u.id
}

// Email はメールアドレスを返す
func (u *User) Email() Email {
	return u.email
}

// PasswordHash はパスワードハッシュを返す
func (u *User) PasswordHash() PasswordHash {
	return u.passwordHash
}

// CreatedAt は作成日時を返す
func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

// UpdatedAt は更新日時を返す
func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

// Provider は認証プロバイダーを返す
func (u *User) Provider() AuthProvider {
	return u.provider
}

// ProviderUserID はプロバイダーのユーザーIDを返す
func (u *User) ProviderUserID() string {
	return u.providerUserID
}

// Name はユーザー名を返す
func (u *User) Name() string {
	return u.name
}

// AvatarURL はアバターURLを返す
func (u *User) AvatarURL() string {
	return u.avatarURL
}

// EmailVerified はメールアドレスが検証済みかどうかを返す
func (u *User) EmailVerified() bool {
	return u.emailVerified
}

// EmailVerifiedAt はメールアドレスの検証日時を返す
func (u *User) EmailVerifiedAt() *time.Time {
	return u.emailVerifiedAt
}

// IsOAuthUser はOAuthユーザーかどうかを返す
func (u *User) IsOAuthUser() bool {
	return u.provider != AuthProviderLocal
}

// VerifyPassword はパスワードが正しいか検証する
func (u *User) VerifyPassword(plainPassword string) bool {
	return u.passwordHash.Compare(plainPassword) == nil
}

// UpdatePassword はパスワードを更新する
func (u *User) UpdatePassword(newPlainPassword string) error {
	newPasswordHash, err := NewPasswordHash(newPlainPassword)
	if err != nil {
		return err
	}

	u.passwordHash = newPasswordHash
	u.updatedAt = time.Now()

	return nil
}

// TwoFactorEnabled は2FAが有効かどうかを返す
func (u *User) TwoFactorEnabled() bool {
	return u.twoFactorEnabled
}

// TwoFactorSecret は2FAシークレットを返す
func (u *User) TwoFactorSecret() string {
	return u.twoFactorSecret
}

// TwoFactorBackupCodes はバックアップコードを返す
func (u *User) TwoFactorBackupCodes() []string {
	return u.twoFactorBackupCodes
}

// EnableTwoFactor は2FAを有効化する
func (u *User) EnableTwoFactor(secret string, backupCodes []string) error {
	if secret == "" {
		return errors.New("2FAシークレットは必須です")
	}
	if len(backupCodes) == 0 {
		return errors.New("バックアップコードは必須です")
	}

	u.twoFactorEnabled = true
	u.twoFactorSecret = secret
	u.twoFactorBackupCodes = backupCodes
	u.updatedAt = time.Now()

	return nil
}

// DisableTwoFactor は2FAを無効化する
func (u *User) DisableTwoFactor() {
	u.twoFactorEnabled = false
	u.twoFactorSecret = ""
	u.twoFactorBackupCodes = nil
	u.updatedAt = time.Now()
}

// RegenerateBackupCodes はバックアップコードを再生成する
func (u *User) RegenerateBackupCodes(backupCodes []string) error {
	if !u.twoFactorEnabled {
		return errors.New("2FAが有効になっていません")
	}
	if len(backupCodes) == 0 {
		return errors.New("バックアップコードは必須です")
	}

	u.twoFactorBackupCodes = backupCodes
	u.updatedAt = time.Now()

	return nil
}

// RemoveBackupCode は使用済みのバックアップコードを削除する
func (u *User) RemoveBackupCode(usedCode string) error {
	if !u.twoFactorEnabled {
		return errors.New("2FAが有効になっていません")
	}

	for i, code := range u.twoFactorBackupCodes {
		if code == usedCode {
			u.twoFactorBackupCodes = append(u.twoFactorBackupCodes[:i], u.twoFactorBackupCodes[i+1:]...)
			u.updatedAt = time.Now()
			return nil
		}
	}

	return errors.New("指定されたバックアップコードは存在しません")
}
