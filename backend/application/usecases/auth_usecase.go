package usecases

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/financial-planning-calculator/backend/domain/entities"
	"github.com/financial-planning-calculator/backend/domain/repositories"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// AuthUseCase は認証関連のユースケース
type AuthUseCase interface {
	// Register は新しいユーザーを登録する
	Register(ctx context.Context, input RegisterInput) (*RegisterOutput, error)

	// Login はユーザー認証を行い、JWTトークンを発行する
	Login(ctx context.Context, input LoginInput) (*LoginOutput, error)

	// VerifyToken はJWTトークンを検証する
	VerifyToken(ctx context.Context, tokenString string) (*TokenClaims, error)
}

// RegisterInput はユーザー登録の入力
type RegisterInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterOutput はユーザー登録の出力
type RegisterOutput struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at"`
}

// LoginInput はログインの入力
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginOutput はログインの出力
type LoginOutput struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at"`
}

// TokenClaims はJWTトークンのクレーム
type TokenClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// authUseCase は認証ユースケースの実装
type authUseCase struct {
	userRepo      repositories.UserRepository
	jwtSecret     string
	jwtExpiration time.Duration
}

// NewAuthUseCase は新しい認証ユースケースを作成する
func NewAuthUseCase(
	userRepo repositories.UserRepository,
	jwtSecret string,
	jwtExpiration time.Duration,
) AuthUseCase {
	return &authUseCase{
		userRepo:      userRepo,
		jwtSecret:     jwtSecret,
		jwtExpiration: jwtExpiration,
	}
}

// Register は新しいユーザーを登録する
func (uc *authUseCase) Register(ctx context.Context, input RegisterInput) (*RegisterOutput, error) {
	logger := slog.With("usecase", "Register", "email", input.Email)
	logger.InfoContext(ctx, "ユーザー登録を開始します")

	// バリデーション
	if input.Email == "" {
		return nil, errors.New("メールアドレスは必須です")
	}
	if input.Password == "" {
		return nil, errors.New("パスワードは必須です")
	}

	// メールアドレスの重複チェック
	email, err := entities.NewEmail(input.Email)
	if err != nil {
		return nil, fmt.Errorf("無効なメールアドレスです: %w", err)
	}

	exists, err := uc.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		logger.ErrorContext(ctx, "メールアドレスの重複チェックに失敗しました", "error", err)
		return nil, fmt.Errorf("メールアドレスの確認に失敗しました: %w", err)
	}
	if exists {
		return nil, errors.New("このメールアドレスは既に登録されています")
	}

	// ユーザーエンティティを作成
	userID := uuid.New().String()
	user, err := entities.NewUser(userID, input.Email, input.Password)
	if err != nil {
		logger.ErrorContext(ctx, "ユーザーエンティティの作成に失敗しました", "error", err)
		return nil, fmt.Errorf("ユーザーの作成に失敗しました: %w", err)
	}

	// ユーザーを保存
	if err := uc.userRepo.Save(ctx, user); err != nil {
		logger.ErrorContext(ctx, "ユーザーの保存に失敗しました", "error", err)
		return nil, fmt.Errorf("ユーザーの保存に失敗しました: %w", err)
	}

	// JWTトークンを生成
	token, expiresAt, err := uc.generateToken(user)
	if err != nil {
		logger.ErrorContext(ctx, "トークンの生成に失敗しました", "error", err)
		return nil, fmt.Errorf("トークンの生成に失敗しました: %w", err)
	}

	logger.InfoContext(ctx, "ユーザー登録が完了しました", "user_id", userID)

	return &RegisterOutput{
		UserID:    userID,
		Email:     input.Email,
		Token:     token,
		ExpiresAt: expiresAt.Format(time.RFC3339),
	}, nil
}

// Login はユーザー認証を行い、JWTトークンを発行する
func (uc *authUseCase) Login(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	logger := slog.With("usecase", "Login", "email", input.Email)
	logger.InfoContext(ctx, "ログインを開始します")

	// バリデーション
	if input.Email == "" {
		return nil, errors.New("メールアドレスは必須です")
	}
	if input.Password == "" {
		return nil, errors.New("パスワードは必須です")
	}

	// メールアドレスでユーザーを取得
	email, err := entities.NewEmail(input.Email)
	if err != nil {
		return nil, fmt.Errorf("無効なメールアドレスです: %w", err)
	}

	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		logger.WarnContext(ctx, "ユーザーが見つかりません", "error", err)
		return nil, errors.New("メールアドレスまたはパスワードが正しくありません")
	}

	// パスワードを検証
	if !user.VerifyPassword(input.Password) {
		logger.WarnContext(ctx, "パスワードが一致しません")
		return nil, errors.New("メールアドレスまたはパスワードが正しくありません")
	}

	// JWTトークンを生成
	token, expiresAt, err := uc.generateToken(user)
	if err != nil {
		logger.ErrorContext(ctx, "トークンの生成に失敗しました", "error", err)
		return nil, fmt.Errorf("トークンの生成に失敗しました: %w", err)
	}

	logger.InfoContext(ctx, "ログインが完了しました", "user_id", user.ID())

	return &LoginOutput{
		UserID:    user.ID().String(),
		Email:     user.Email().String(),
		Token:     token,
		ExpiresAt: expiresAt.Format(time.RFC3339),
	}, nil
}

// VerifyToken はJWTトークンを検証する
func (uc *authUseCase) VerifyToken(ctx context.Context, tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 署名アルゴリズムの確認
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("予期しない署名方法です: %v", token.Header["alg"])
		}
		return []byte(uc.jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("トークンの検証に失敗しました: %w", err)
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("無効なトークンです")
}

// generateToken はユーザー情報からJWTトークンを生成する
func (uc *authUseCase) generateToken(user *entities.User) (string, time.Time, error) {
	expiresAt := time.Now().Add(uc.jwtExpiration)

	claims := TokenClaims{
		UserID: user.ID().String(),
		Email:  user.Email().String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(uc.jwtSecret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}
