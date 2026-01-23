package usecases

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/financial-planning-calculator/backend/domain/entities"
	"github.com/financial-planning-calculator/backend/domain/repositories"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
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

	// Setup2FA は2段階認証のセットアップを開始する（QRコード生成用）
	Setup2FA(ctx context.Context, userID string) (*Setup2FAOutput, error)

	// Enable2FA は2段階認証を有効化する（初回コード検証）
	Enable2FA(ctx context.Context, input Enable2FAInput) error

	// Verify2FA はログイン時の2FAコード検証を行う
	Verify2FA(ctx context.Context, input Verify2FAInput) (*LoginOutput, error)

	// Disable2FA は2段階認証を無効化する
	Disable2FA(ctx context.Context, input Disable2FAInput) error

	// RegenerateBackupCodes はバックアップコードを再生成する
	RegenerateBackupCodes(ctx context.Context, userID string) (*RegenerateBackupCodesOutput, error)

	// Get2FAStatus は2FAの有効状態を取得する
	Get2FAStatus(ctx context.Context, userID string) (*Get2FAStatusOutput, error)
}

// Get2FAStatusOutput は2FAステータス取得の出力
type Get2FAStatusOutput struct {
	Enabled bool `json:"enabled"`
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
	UserID          string `json:"user_id"`
	Email           string `json:"email"`
	Requires2FA     bool   `json:"requires_2fa,omitempty"`     // 2FA検証が必要かどうか
	TwoFactorVerify bool   `json:"two_factor_verify,omitempty"` // 2FA検証用の仮トークンかどうか
	jwt.RegisteredClaims
}

// Setup2FAOutput は2FA設定開始の出力
type Setup2FAOutput struct {
	Secret       string   `json:"secret"`
	QRCodeURL    string   `json:"qr_code_url"`
	BackupCodes  []string `json:"backup_codes"`
}

// Enable2FAInput は2FA有効化の入力
type Enable2FAInput struct {
	UserID string `json:"user_id"`
	Code   string `json:"code"`
	Secret string `json:"secret"`
}

// Verify2FAInput は2FA検証の入力
type Verify2FAInput struct {
	UserID      string `json:"user_id"`
	Code        string `json:"code"`
	UseBackup   bool   `json:"use_backup"`   // バックアップコードを使用するか
}

// Disable2FAInput は2FA無効化の入力
type Disable2FAInput struct {
	UserID   string `json:"user_id"`
	Password string `json:"password"`
}

// RegenerateBackupCodesOutput はバックアップコード再生成の出力
type RegenerateBackupCodesOutput struct {
	BackupCodes []string `json:"backup_codes"`
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

	// 2FAが有効な場合は仮トークンを返す
	if user.TwoFactorEnabled() {
		logger.InfoContext(ctx, "2FAが有効なため仮トークンを発行します", "user_id", user.ID())
		
		// 2FA検証用の短時間有効な仮トークンを生成（5分間）
		tempToken, tempExpiresAt, err := uc.generateTempTokenFor2FA(user)
		if err != nil {
			logger.ErrorContext(ctx, "仮トークンの生成に失敗しました", "error", err)
			return nil, fmt.Errorf("認証処理に失敗しました: %w", err)
		}

		return &LoginOutput{
			UserID:       user.ID().String(),
			Email:        user.Email().String(),
			Token:        tempToken,
			RefreshToken: "", // 2FA検証前はリフレッシュトークンなし
			ExpiresAt:    tempExpiresAt.Format(time.RFC3339),
		}, nil
	}

	// 2FAが無効な場合は通常のトークンを発行
	logger.InfoContext(ctx, "通常のトークンを発行します", "user_id", user.ID())
	return uc.generateAuthTokens(ctx, user)
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

// generateTempTokenFor2FA は2FA検証用の短時間有効な仮トークンを生成する
func (uc *authUseCase) generateTempTokenFor2FA(user *entities.User) (string, time.Time, error) {
	// 5分間有効な仮トークン
	expiresAt := time.Now().Add(5 * time.Minute)

	claims := TokenClaims{
		UserID:          user.ID().String(),
		Email:           user.Email().String(),
		Requires2FA:     true,
		TwoFactorVerify: true,
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

	// GitHubユーザーが見つからない - メールアドレスで既存ユーザーを検索
	email, err := entities.NewEmail(input.Email)
	if err != nil {
		return nil, fmt.Errorf("無効なメールアドレスです: %w", err)
	}

	existingUserByEmail, err := uc.userRepo.FindByEmail(ctx, email)
	if err == nil {
		// 同一メールアドレスの既存ユーザーが見つかった
		// セキュリティ上の理由から、自動リンクは行わない
		// ユーザーは既存のアカウントでログインする必要がある
		logger.WarnContext(ctx, "同一メールアドレスの既存アカウントが見つかりました",
			"existing_user_id", existingUserByEmail.ID(),
			"existing_provider", existingUserByEmail.Provider())
		return nil, fmt.Errorf("このメールアドレスは既に登録されています。既存のアカウント（%s）でログインしてください", existingUserByEmail.Provider())
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

// Setup2FA は2段階認証のセットアップを開始する（QRコード生成用）
func (uc *authUseCase) Setup2FA(ctx context.Context, userID string) (*Setup2FAOutput, error) {
	logger := slog.With("usecase", "Setup2FA", "user_id", userID)
	logger.InfoContext(ctx, "2FA設定を開始します")

	// ユーザーを取得
	uid, err := entities.NewUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("無効なユーザーIDです: %w", err)
	}

	user, err := uc.userRepo.FindByID(ctx, uid)
	if err != nil {
		logger.ErrorContext(ctx, "ユーザーの取得に失敗しました", "error", err)
		return nil, fmt.Errorf("ユーザーが見つかりません: %w", err)
	}

	// 既に2FAが有効な場合はエラー
	if user.TwoFactorEnabled() {
		return nil, errors.New("2段階認証は既に有効です")
	}

	// TOTPシークレットを生成
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Financial Planning Calculator",
		AccountName: user.Email().String(),
		SecretSize:  32,
	})
	if err != nil {
		logger.ErrorContext(ctx, "TOTPキーの生成に失敗しました", "error", err)
		return nil, fmt.Errorf("2FAシークレットの生成に失敗しました: %w", err)
	}

	// バックアップコードを生成
	backupCodes, err := generateBackupCodes(8)
	if err != nil {
		logger.ErrorContext(ctx, "バックアップコードの生成に失敗しました", "error", err)
		return nil, fmt.Errorf("バックアップコードの生成に失敗しました: %w", err)
	}

	logger.InfoContext(ctx, "2FA設定データを生成しました")

	return &Setup2FAOutput{
		Secret:      key.Secret(),
		QRCodeURL:   key.URL(),
		BackupCodes: backupCodes,
	}, nil
}

// Enable2FA は2段階認証を有効化する（初回コード検証）
func (uc *authUseCase) Enable2FA(ctx context.Context, input Enable2FAInput) error {
	logger := slog.With("usecase", "Enable2FA", "user_id", input.UserID)
	logger.InfoContext(ctx, "2FA有効化を開始します")

	// バリデーション
	if input.UserID == "" {
		return errors.New("ユーザーIDは必須です")
	}
	if input.Code == "" {
		return errors.New("認証コードは必須です")
	}
	if input.Secret == "" {
		return errors.New("シークレットは必須です")
	}

	// ユーザーを取得
	uid, err := entities.NewUserID(input.UserID)
	if err != nil {
		return fmt.Errorf("無効なユーザーIDです: %w", err)
	}

	user, err := uc.userRepo.FindByID(ctx, uid)
	if err != nil {
		logger.ErrorContext(ctx, "ユーザーの取得に失敗しました", "error", err)
		return fmt.Errorf("ユーザーが見つかりません: %w", err)
	}

	// 既に2FAが有効な場合はエラー
	if user.TwoFactorEnabled() {
		return errors.New("2段階認証は既に有効です")
	}

	// TOTPコードを検証（時間のずれを許容するためValidateCustomを使用）
	valid, err := totp.ValidateCustom(input.Code, input.Secret, time.Now().UTC(), totp.ValidateOpts{
		Period:    30,
		Skew:      1,
		Digits:    6,
		Algorithm: otp.AlgorithmSHA1,
	})
	logger.InfoContext(ctx, "TOTP検証", "code", input.Code, "secretLength", len(input.Secret), "valid", valid, "time", time.Now().UTC())
	if err != nil || !valid {
		logger.WarnContext(ctx, "2FAコードの検証に失敗しました", "error", err)
		return errors.New("認証コードが無効です")
	}

	// バックアップコードを再生成（Enable2FAInputにバックアップコードがない場合）
	backupCodes, err := generateBackupCodes(8)
	if err != nil {
		logger.ErrorContext(ctx, "バックアップコードの生成に失敗しました", "error", err)
		return fmt.Errorf("バックアップコードの生成に失敗しました: %w", err)
	}

	// バックアップコードをハッシュ化
	hashedBackupCodes, err := hashBackupCodes(backupCodes)
	if err != nil {
		logger.ErrorContext(ctx, "バックアップコードのハッシュ化に失敗しました", "error", err)
		return fmt.Errorf("バックアップコードの保存に失敗しました: %w", err)
	}

	// 2FAを有効化
	if err := user.EnableTwoFactor(input.Secret, hashedBackupCodes); err != nil {
		logger.ErrorContext(ctx, "2FAの有効化に失敗しました", "error", err)
		return fmt.Errorf("2FAの有効化に失敗しました: %w", err)
	}

	// ユーザーを保存
	if err := uc.userRepo.Update(ctx, user); err != nil {
		logger.ErrorContext(ctx, "ユーザーの更新に失敗しました", "error", err)
		return fmt.Errorf("ユーザーの更新に失敗しました: %w", err)
	}

	logger.InfoContext(ctx, "2FAを有効化しました")
	return nil
}

// Verify2FA はログイン時の2FAコード検証を行う
func (uc *authUseCase) Verify2FA(ctx context.Context, input Verify2FAInput) (*LoginOutput, error) {
	logger := slog.With("usecase", "Verify2FA", "user_id", input.UserID)
	logger.InfoContext(ctx, "2FA検証を開始します")

	// バリデーション
	if input.UserID == "" {
		return nil, errors.New("ユーザーIDは必須です")
	}
	if input.Code == "" {
		return nil, errors.New("認証コードは必須です")
	}

	// ユーザーを取得
	uid, err := entities.NewUserID(input.UserID)
	if err != nil {
		return nil, fmt.Errorf("無効なユーザーIDです: %w", err)
	}

	user, err := uc.userRepo.FindByID(ctx, uid)
	if err != nil {
		logger.ErrorContext(ctx, "ユーザーの取得に失敗しました", "error", err)
		return nil, fmt.Errorf("ユーザーが見つかりません: %w", err)
	}

	// 2FAが有効でない場合はエラー
	if !user.TwoFactorEnabled() {
		return nil, errors.New("2段階認証が有効になっていません")
	}

	var verified bool

	if input.UseBackup {
		// バックアップコードで検証
		verified = false
		for _, hashedCode := range user.TwoFactorBackupCodes() {
			if err := bcrypt.CompareHashAndPassword([]byte(hashedCode), []byte(input.Code)); err == nil {
				verified = true
				// 使用済みバックアップコードを削除
				if err := user.RemoveBackupCode(hashedCode); err != nil {
					logger.ErrorContext(ctx, "バックアップコードの削除に失敗しました", "error", err)
				} else {
					// ユーザーを更新
					if err := uc.userRepo.Update(ctx, user); err != nil {
						logger.ErrorContext(ctx, "ユーザーの更新に失敗しました", "error", err)
					}
				}
				break
			}
		}
	} else {
		// TOTPコードで検証
		verified = totp.Validate(input.Code, user.TwoFactorSecret())
	}

	if !verified {
		logger.WarnContext(ctx, "2FAコードの検証に失敗しました")
		return nil, errors.New("認証コードが無効です")
	}

	// 認証成功 - 通常のトークンを発行
	logger.InfoContext(ctx, "2FA検証に成功しました")
	return uc.generateAuthTokens(ctx, user)
}

// Disable2FA は2段階認証を無効化する
func (uc *authUseCase) Disable2FA(ctx context.Context, input Disable2FAInput) error {
	logger := slog.With("usecase", "Disable2FA", "user_id", input.UserID)
	logger.InfoContext(ctx, "2FA無効化を開始します")

	// バリデーション
	if input.UserID == "" {
		return errors.New("ユーザーIDは必須です")
	}
	if input.Password == "" {
		return errors.New("パスワードは必須です")
	}

	// ユーザーを取得
	uid, err := entities.NewUserID(input.UserID)
	if err != nil {
		return fmt.Errorf("無効なユーザーIDです: %w", err)
	}

	user, err := uc.userRepo.FindByID(ctx, uid)
	if err != nil {
		logger.ErrorContext(ctx, "ユーザーの取得に失敗しました", "error", err)
		return fmt.Errorf("ユーザーが見つかりません: %w", err)
	}

	// ローカルユーザーの場合はパスワード検証が必要
	if user.Provider() == entities.AuthProviderLocal {
		if !user.VerifyPassword(input.Password) {
			logger.WarnContext(ctx, "パスワード検証に失敗しました")
			return errors.New("パスワードが正しくありません")
		}
	}

	// 2FAが有効でない場合はエラー
	if !user.TwoFactorEnabled() {
		return errors.New("2段階認証は有効になっていません")
	}

	// 2FAを無効化
	user.DisableTwoFactor()

	// ユーザーを保存
	if err := uc.userRepo.Update(ctx, user); err != nil {
		logger.ErrorContext(ctx, "ユーザーの更新に失敗しました", "error", err)
		return fmt.Errorf("ユーザーの更新に失敗しました: %w", err)
	}

	logger.InfoContext(ctx, "2FAを無効化しました")
	return nil
}

// RegenerateBackupCodes はバックアップコードを再生成する
func (uc *authUseCase) RegenerateBackupCodes(ctx context.Context, userID string) (*RegenerateBackupCodesOutput, error) {
	logger := slog.With("usecase", "RegenerateBackupCodes", "user_id", userID)
	logger.InfoContext(ctx, "バックアップコード再生成を開始します")

	// ユーザーを取得
	uid, err := entities.NewUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("無効なユーザーIDです: %w", err)
	}

	user, err := uc.userRepo.FindByID(ctx, uid)
	if err != nil {
		logger.ErrorContext(ctx, "ユーザーの取得に失敗しました", "error", err)
		return nil, fmt.Errorf("ユーザーが見つかりません: %w", err)
	}

	// 2FAが有効でない場合はエラー
	if !user.TwoFactorEnabled() {
		return nil, errors.New("2段階認証が有効になっていません")
	}

	// 新しいバックアップコードを生成
	backupCodes, err := generateBackupCodes(8)
	if err != nil {
		logger.ErrorContext(ctx, "バックアップコードの生成に失敗しました", "error", err)
		return nil, fmt.Errorf("バックアップコードの生成に失敗しました: %w", err)
	}

	// バックアップコードをハッシュ化
	hashedBackupCodes, err := hashBackupCodes(backupCodes)
	if err != nil {
		logger.ErrorContext(ctx, "バックアップコードのハッシュ化に失敗しました", "error", err)
		return nil, fmt.Errorf("バックアップコードの保存に失敗しました: %w", err)
	}

	// バックアップコードを再生成
	if err := user.RegenerateBackupCodes(hashedBackupCodes); err != nil {
		logger.ErrorContext(ctx, "バックアップコードの再生成に失敗しました", "error", err)
		return nil, fmt.Errorf("バックアップコードの再生成に失敗しました: %w", err)
	}

	// ユーザーを保存
	if err := uc.userRepo.Update(ctx, user); err != nil {
		logger.ErrorContext(ctx, "ユーザーの更新に失敗しました", "error", err)
		return nil, fmt.Errorf("ユーザーの更新に失敗しました: %w", err)
	}

	logger.InfoContext(ctx, "バックアップコードを再生成しました")

	return &RegenerateBackupCodesOutput{
		BackupCodes: backupCodes,
	}, nil
}

// generateBackupCodes はランダムなバックアップコードを生成する
func generateBackupCodes(count int) ([]string, error) {
	codes := make([]string, count)
	for i := 0; i < count; i++ {
		// 10バイトのランダムデータを生成
		randomBytes := make([]byte, 10)
		if _, err := rand.Read(randomBytes); err != nil {
			return nil, fmt.Errorf("ランダムバイトの生成に失敗しました: %w", err)
		}

		// Base32エンコードして8文字のコードを生成
		code := base32.StdEncoding.EncodeToString(randomBytes)[:8]
		codes[i] = strings.ToUpper(code)
	}
	return codes, nil
}

// hashBackupCodes はバックアップコードをbcryptでハッシュ化する
func hashBackupCodes(codes []string) ([]string, error) {
	hashedCodes := make([]string, len(codes))
	for i, code := range codes {
		hash, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("バックアップコードのハッシュ化に失敗しました: %w", err)
		}
		hashedCodes[i] = string(hash)
	}
	return hashedCodes, nil
}

// Get2FAStatus は2FAの有効状態を取得する
func (uc *authUseCase) Get2FAStatus(ctx context.Context, userID string) (*Get2FAStatusOutput, error) {
	logger := slog.With("usecase", "Get2FAStatus", "user_id", userID)
	logger.InfoContext(ctx, "2FAステータス取得を開始します")

	// ユーザーを取得
	uid, err := entities.NewUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("無効なユーザーIDです: %w", err)
	}

	user, err := uc.userRepo.FindByID(ctx, uid)
	if err != nil {
		logger.ErrorContext(ctx, "ユーザーの取得に失敗しました", "error", err)
		return nil, fmt.Errorf("ユーザーが見つかりません: %w", err)
	}

	return &Get2FAStatusOutput{
		Enabled: user.TwoFactorEnabled(),
	}, nil
}
