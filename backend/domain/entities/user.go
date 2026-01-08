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

// User はユーザーエンティティ
type User struct {
	id           UserID
	email        Email
	passwordHash PasswordHash
	createdAt    time.Time
	updatedAt    time.Time
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
		id:           userID,
		email:        emailVO,
		passwordHash: passwordHash,
		createdAt:    now,
		updatedAt:    now,
	}, nil
}

// ReconstructUser はDBから取得したデータからUserを再構築する（リポジトリ用）
func ReconstructUser(id string, email string, passwordHash string, createdAt, updatedAt time.Time) (*User, error) {
	userID, err := NewUserID(id)
	if err != nil {
		return nil, err
	}

	emailVO, err := NewEmail(email)
	if err != nil {
		return nil, err
	}

	return &User{
		id:           userID,
		email:        emailVO,
		passwordHash: NewPasswordHashFromHash(passwordHash),
		createdAt:    createdAt,
		updatedAt:    updatedAt,
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
