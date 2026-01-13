package usecases

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
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

	// RefreshAccessToken はリフレッシュトークンを使用して新しいアクセストークンを発行する
	RefreshAccessToken(ctx context.Context, refreshToken string) (*RefreshOutput, error)

	// RevokeRefreshToken はリフレッシュトークンを失効させる（ログアウト時に使用）
	RevokeRefreshToken(ctx context.Context, userID string) error

	// GitHubOAuthLogin はGitHubからのユーザー情報でログイン/登録を行う（Issue: #67）
	GitHubOAuthLogin(ctx context.Context, input GitHubOAuthInput) (*LoginOutput, error)
}

// GitHubOAuthInput はGitHub OAuthログインの入力
type GitHubOAuthInput struct {
	GitHubUserID string `json:"github_user_id"`
	Email        string `json:"email"`
	Name         string `json:"name"`
	AvatarURL    string `json:"avatar_url"`
}

// RegisterInput はユーザー登録の入力
type RegisterInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterOutput はユーザー登録の出力
type RegisterOutput struct {
	UserID       string `json:"user_id"`
	Email        string `json:"email"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    string `json:"expires_at"`
}

// LoginInput はログインの入力
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginOutput はログインの出力
type LoginOutput struct {
	UserID       string `json:"user_id"`
	Email        string `json:"email"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    string `json:"expires_at"`
}

// RefreshOutput はトークンリフレッシュの出力
type RefreshOutput struct {
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
	userRepo                repositories.UserRepository
	refreshTokenRepo        repositories.RefreshTokenRepository
	jwtSecret               string
	jwtExpiration           time.Duration
	refreshTokenExpiration  time.Duration
}

// NewAuthUseCase は新しい認証ユースケースを作成する
func NewAuthUseCase(
	userRepo repositories.UserRepository,
	refreshTokenRepo repositories.RefreshTokenRepository,
	jwtSecret string,
	jwtExpiration time.Duration,
	refreshTokenExpiration time.Duration,
) AuthUseCase {
	return &authUseCase{
		userRepo:               userRepo,
		refreshTokenRepo:       refreshTokenRepo,
		jwtSecret:              jwtSecret,
		jwtExpiration:          jwtExpiration,
		refreshTokenExpiration: refreshTokenExpiration,
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

	// リフレッシュトークンを生成
	refreshToken, err := uc.generateRefreshToken(ctx, user.ID())
	if err != nil {
		logger.ErrorContext(ctx, "リフレッシュトークンの生成に失敗しました", "error", err)
		return nil, fmt.Errorf("リフレッシュトークンの生成に失敗しました: %w", err)
	}

	logger.InfoContext(ctx, "ユーザー登録が完了しました", "user_id", userID)

	return &RegisterOutput{
		UserID:       userID,
		Email:        input.Email,
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt.Format(time.RFC3339),
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

	// リフレッシュトークンを生成
	refreshToken, err := uc.generateRefreshToken(ctx, user.ID())
	if err != nil {
		logger.ErrorContext(ctx, "リフレッシュトークンの生成に失敗しました", "error", err)
		return nil, fmt.Errorf("リフレッシュトークンの生成に失敗しました: %w", err)
	}

	logger.InfoContext(ctx, "ログインが完了しました", "user_id", user.ID())

	return &LoginOutput{
		UserID:       user.ID().String(),
		Email:        user.Email().String(),
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt.Format(time.RFC3339),
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

// generateRefreshToken はリフレッシュトークンを生成してDBに保存する
func (uc *authUseCase) generateRefreshToken(ctx context.Context, userID entities.UserID) (string, error) {
	expiresAt := time.Now().Add(uc.refreshTokenExpiration)

	refreshToken, token, err := entities.NewRefreshToken(userID, expiresAt)
	if err != nil {
		return "", fmt.Errorf("リフレッシュトークンの生成に失敗しました: %w", err)
	}

	if err := uc.refreshTokenRepo.Save(ctx, refreshToken); err != nil {
		return "", fmt.Errorf("リフレッシュトークンの保存に失敗しました: %w", err)
	}

	return token, nil
}

// RefreshAccessToken はリフレッシュトークンを使用して新しいアクセストークンを発行する
func (uc *authUseCase) RefreshAccessToken(ctx context.Context, refreshTokenString string) (*RefreshOutput, error) {
	logger := slog.With("usecase", "RefreshAccessToken")
	logger.InfoContext(ctx, "トークンリフレッシュを開始します")

	// リフレッシュトークンをハッシュ化して検索
	tokenHash := hashRefreshToken(refreshTokenString)
	refreshToken, err := uc.refreshTokenRepo.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		logger.WarnContext(ctx, "リフレッシュトークンが見つかりません", "error", err)
		return nil, errors.New("無効なリフレッシュトークンです")
	}

	// トークンを検証
	if !refreshToken.VerifyToken(refreshTokenString) {
		logger.WarnContext(ctx, "リフレッシュトークンの検証に失敗しました")
		return nil, errors.New("無効なリフレッシュトークンです")
	}

	if !refreshToken.IsValid() {
		logger.WarnContext(ctx, "リフレッシュトークンが無効です", "expired", refreshToken.IsExpired(), "revoked", refreshToken.IsRevoked())
		return nil, errors.New("リフレッシュトークンの有効期限が切れているか、失効されています")
	}

	// ユーザーを取得
	user, err := uc.userRepo.FindByID(ctx, refreshToken.UserID())
	if err != nil {
		logger.ErrorContext(ctx, "ユーザーが見つかりません", "error", err)
		return nil, errors.New("ユーザーが見つかりません")
	}

	// 新しいアクセストークンを生成
	token, expiresAt, err := uc.generateToken(user)
	if err != nil {
		logger.ErrorContext(ctx, "トークンの生成に失敗しました", "error", err)
		return nil, fmt.Errorf("トークンの生成に失敗しました: %w", err)
	}

	// リフレッシュトークンの最終使用日時を更新
	refreshToken.UpdateLastUsedAt()
	if err := uc.refreshTokenRepo.Update(ctx, refreshToken); err != nil {
		logger.ErrorContext(ctx, "リフレッシュトークンの更新に失敗しました", "error", err)
		// エラーをログに記録するが、処理は続行
	}

	logger.InfoContext(ctx, "トークンリフレッシュが完了しました", "user_id", user.ID())

	return &RefreshOutput{
		Token:     token,
		ExpiresAt: expiresAt.Format(time.RFC3339),
	}, nil
}

// RevokeRefreshToken はリフレッシュトークンを失効させる（ログアウト時に使用）
func (uc *authUseCase) RevokeRefreshToken(ctx context.Context, userID string) error {
	logger := slog.With("usecase", "RevokeRefreshToken", "user_id", userID)
	logger.InfoContext(ctx, "リフレッシュトークンの失効を開始します")

	uid, err := entities.NewUserID(userID)
	if err != nil {
		return fmt.Errorf("無効なユーザーIDです: %w", err)
	}

	if err := uc.refreshTokenRepo.RevokeByUserID(ctx, uid); err != nil {
		logger.ErrorContext(ctx, "リフレッシュトークンの失効に失敗しました", "error", err)
		return fmt.Errorf("リフレッシュトークンの失効に失敗しました: %w", err)
	}

	logger.InfoContext(ctx, "リフレッシュトークンの失効が完了しました")
	return nil
}

// hashRefreshToken はリフレッシュトークンをハッシュ化する（entities.RefreshTokenと同じロジック）
func hashRefreshToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// GitHubOAuthLogin はGitHubからのユーザー情報でログイン/登録を行う（Issue: #67）
func (uc *authUseCase) GitHubOAuthLogin(ctx context.Context, input GitHubOAuthInput) (*LoginOutput, error) {
	logger := slog.With("usecase", "GitHubOAuthLogin", "github_user_id", input.GitHubUserID, "email", input.Email)
	logger.InfoContext(ctx, "GitHub OAuthログインを開始します")

	// バリデーション
	if input.GitHubUserID == "" {
		return nil, errors.New("GitHub user IDは必須です")
	}
	if input.Email == "" {
		return nil, errors.New("メールアドレスは必須です")
	}

	// GitHubプロバイダーIDで既存ユーザーを検索
	existingUser, err := uc.userRepo.FindByProviderUserID(ctx, entities.AuthProviderGitHub, input.GitHubUserID)
	if err == nil {
		// 既存のGitHubユーザーが見つかった - ログイン処理
		logger.InfoContext(ctx, "既存のGitHubユーザーでログインします", "user_id", existingUser.ID())
		return uc.generateAuthTokens(ctx, existingUser)
	}

	// GitHubユーザーが見つからない - メールアドレスで既存ユーザーを検索（自動リンク）
	email, err := entities.NewEmail(input.Email)
	if err != nil {
		return nil, fmt.Errorf("無効なメールアドレスです: %w", err)
	}

	existingUserByEmail, err := uc.userRepo.FindByEmail(ctx, email)
	if err == nil {
		// 同一メールアドレスの既存ユーザーが見つかった
		// GitHubのメールは検証済みなので、既存アカウントでログインを許可
		logger.InfoContext(ctx, "同一メールアドレスの既存ユーザーでログインします", "existing_user_id", existingUserByEmail.ID())
		return uc.generateAuthTokens(ctx, existingUserByEmail)
	}

	// 新規ユーザーを作成
	userID := uuid.New().String()
	newUser, err := entities.NewOAuthUser(
		userID,
		input.Email,
		entities.AuthProviderGitHub,
		input.GitHubUserID,
		input.Name,
		input.AvatarURL,
	)
	if err != nil {
		logger.ErrorContext(ctx, "ユーザーエンティティの作成に失敗しました", "error", err)
		return nil, fmt.Errorf("ユーザーの作成に失敗しました: %w", err)
	}

	// ユーザーを保存
	if err := uc.userRepo.Save(ctx, newUser); err != nil {
		logger.ErrorContext(ctx, "ユーザーの保存に失敗しました", "error", err)
		return nil, fmt.Errorf("ユーザーの保存に失敗しました: %w", err)
	}

	logger.InfoContext(ctx, "新規GitHubユーザーを作成しました", "user_id", newUser.ID())

	// トークンを生成して返す
	return uc.generateAuthTokens(ctx, newUser)
}

// generateAuthTokens はユーザーの認証トークンを生成する（共通処理）
func (uc *authUseCase) generateAuthTokens(ctx context.Context, user *entities.User) (*LoginOutput, error) {
	// JWTトークンを生成
	token, expiresAt, err := uc.generateToken(user)
	if err != nil {
		return nil, fmt.Errorf("トークンの生成に失敗しました: %w", err)
	}

	// リフレッシュトークンを生成してDBに保存
	refreshTokenValue, err := uc.generateRefreshToken(ctx, user.ID())
	if err != nil {
		return nil, fmt.Errorf("リフレッシュトークンの生成に失敗しました: %w", err)
	}

	return &LoginOutput{
		UserID:       user.ID().String(),
		Email:        user.Email().String(),
		Token:        token,
		RefreshToken: refreshTokenValue,
		ExpiresAt:    expiresAt.Format(time.RFC3339),
	}, nil
}
