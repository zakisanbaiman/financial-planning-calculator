package usecases

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/financial-planning-calculator/backend/domain/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testJWTSecret              = "test-secret-key-for-unit-tests-32chars"
	testJWTExpiration          = 15 * time.Minute
	testRefreshTokenExpiration = 7 * 24 * time.Hour
)

func newTestAuthUseCase(userRepo *MockUserRepository, tokenRepo *MockRefreshTokenRepository) AuthUseCase {
	return NewAuthUseCase(userRepo, tokenRepo, testJWTSecret, testJWTExpiration, testRefreshTokenExpiration)
}

// ===========================
// Register Tests
// ===========================

func TestAuthUseCase_Register(t *testing.T) {
	ctx := context.Background()

	t.Run("正常系: 新規ユーザーを登録できる", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockTokenRepo := new(MockRefreshTokenRepository)
		email, _ := entities.NewEmail("test@example.com")
		mockUserRepo.On("ExistsByEmail", mock_anything(), email).Return(false, nil)
		mockUserRepo.On("Save", mock_anything(), mock_anything()).Return(nil)
		mockTokenRepo.On("Save", mock_anything(), mock_anything()).Return(nil)

		uc := newTestAuthUseCase(mockUserRepo, mockTokenRepo)
		output, err := uc.Register(ctx, RegisterInput{
			Email:    "test@example.com",
			Password: "password123",
		})

		require.NoError(t, err)
		assert.NotEmpty(t, output.UserID)
		assert.Equal(t, "test@example.com", output.Email)
		assert.NotEmpty(t, output.Token)
		assert.NotEmpty(t, output.RefreshToken)
		mockUserRepo.AssertExpectations(t)
		mockTokenRepo.AssertExpectations(t)
	})

	t.Run("異常系: メールアドレスが空の場合はエラー", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockTokenRepo := new(MockRefreshTokenRepository)

		uc := newTestAuthUseCase(mockUserRepo, mockTokenRepo)
		_, err := uc.Register(ctx, RegisterInput{
			Email:    "",
			Password: "password123",
		})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "メールアドレスは必須です")
	})

	t.Run("異常系: パスワードが空の場合はエラー", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockTokenRepo := new(MockRefreshTokenRepository)

		uc := newTestAuthUseCase(mockUserRepo, mockTokenRepo)
		_, err := uc.Register(ctx, RegisterInput{
			Email:    "test@example.com",
			Password: "",
		})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "パスワードは必須です")
	})

	t.Run("異常系: 既に登録済みのメールアドレスの場合はエラー", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockTokenRepo := new(MockRefreshTokenRepository)
		email, _ := entities.NewEmail("existing@example.com")
		mockUserRepo.On("ExistsByEmail", mock_anything(), email).Return(true, nil)

		uc := newTestAuthUseCase(mockUserRepo, mockTokenRepo)
		_, err := uc.Register(ctx, RegisterInput{
			Email:    "existing@example.com",
			Password: "password123",
		})

		require.Error(t, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("異常系: ExistsByEmailでリポジトリエラーが発生した場合", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockTokenRepo := new(MockRefreshTokenRepository)
		email, _ := entities.NewEmail("test@example.com")
		mockUserRepo.On("ExistsByEmail", mock_anything(), email).Return(false, errors.New("db error"))

		uc := newTestAuthUseCase(mockUserRepo, mockTokenRepo)
		_, err := uc.Register(ctx, RegisterInput{
			Email:    "test@example.com",
			Password: "password123",
		})

		require.Error(t, err)
		mockUserRepo.AssertExpectations(t)
	})
}

// ===========================
// Login Tests
// ===========================

func TestAuthUseCase_Login(t *testing.T) {
	ctx := context.Background()

	t.Run("異常系: 存在しないメールアドレスの場合はエラー", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockTokenRepo := new(MockRefreshTokenRepository)
		email, _ := entities.NewEmail("notfound@example.com")
		mockUserRepo.On("FindByEmail", mock_anything(), email).Return(nil, errors.New("not found"))

		uc := newTestAuthUseCase(mockUserRepo, mockTokenRepo)
		_, err := uc.Login(ctx, LoginInput{
			Email:    "notfound@example.com",
			Password: "password123",
		})

		require.Error(t, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("異常系: メールアドレスが空の場合はエラー", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockTokenRepo := new(MockRefreshTokenRepository)

		uc := newTestAuthUseCase(mockUserRepo, mockTokenRepo)
		_, err := uc.Login(ctx, LoginInput{
			Email:    "",
			Password: "password123",
		})

		require.Error(t, err)
	})
}

// ===========================
// VerifyToken Tests
// ===========================

func TestAuthUseCase_VerifyToken(t *testing.T) {
	ctx := context.Background()

	t.Run("異常系: 不正なトークンの場合はエラー", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockTokenRepo := new(MockRefreshTokenRepository)

		uc := newTestAuthUseCase(mockUserRepo, mockTokenRepo)
		_, err := uc.VerifyToken(ctx, "invalid.token.value")

		require.Error(t, err)
	})

	t.Run("異常系: 空のトークンの場合はエラー", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockTokenRepo := new(MockRefreshTokenRepository)

		uc := newTestAuthUseCase(mockUserRepo, mockTokenRepo)
		_, err := uc.VerifyToken(ctx, "")

		require.Error(t, err)
	})
}

// ===========================
// RevokeRefreshToken Tests
// ===========================

func TestAuthUseCase_RevokeRefreshToken(t *testing.T) {
	ctx := context.Background()

	t.Run("正常系: リフレッシュトークンを失効できる", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockTokenRepo := new(MockRefreshTokenRepository)
		mockTokenRepo.On("RevokeByUserID", mock_anything(), entities.UserID("user-001")).Return(nil)

		uc := newTestAuthUseCase(mockUserRepo, mockTokenRepo)
		err := uc.RevokeRefreshToken(ctx, "user-001")

		require.NoError(t, err)
		mockTokenRepo.AssertExpectations(t)
	})

	t.Run("異常系: リポジトリエラーの場合はエラーを返す", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockTokenRepo := new(MockRefreshTokenRepository)
		mockTokenRepo.On("RevokeByUserID", mock_anything(), entities.UserID("user-001")).Return(errors.New("db error"))

		uc := newTestAuthUseCase(mockUserRepo, mockTokenRepo)
		err := uc.RevokeRefreshToken(ctx, "user-001")

		require.Error(t, err)
		mockTokenRepo.AssertExpectations(t)
	})
}

// ===========================
// Get2FAStatus Tests
// ===========================

func TestAuthUseCase_Get2FAStatus(t *testing.T) {
	ctx := context.Background()

	t.Run("正常系: 2FAが無効なユーザーのステータスを取得できる", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockTokenRepo := new(MockRefreshTokenRepository)
		user := newTestUser("user-001", "test@example.com")
		mockUserRepo.On("FindByID", mock_anything(), entities.UserID("user-001")).Return(user, nil)

		uc := newTestAuthUseCase(mockUserRepo, mockTokenRepo)
		output, err := uc.Get2FAStatus(ctx, "user-001")

		require.NoError(t, err)
		assert.NotNil(t, output)
		assert.False(t, output.Enabled)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("異常系: ユーザーが存在しない場合はエラー", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockTokenRepo := new(MockRefreshTokenRepository)
		mockUserRepo.On("FindByID", mock_anything(), entities.UserID("user-999")).Return(nil, errors.New("not found"))

		uc := newTestAuthUseCase(mockUserRepo, mockTokenRepo)
		_, err := uc.Get2FAStatus(ctx, "user-999")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "ユーザーが見つかりません")
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("異常系: 無効なユーザーIDの場合はエラー", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockTokenRepo := new(MockRefreshTokenRepository)

		uc := newTestAuthUseCase(mockUserRepo, mockTokenRepo)
		_, err := uc.Get2FAStatus(ctx, "")

		require.Error(t, err)
	})
}

// ===========================
// GitHubOAuthLogin Tests
// ===========================

func TestAuthUseCase_GitHubOAuthLogin(t *testing.T) {
	ctx := context.Background()

	t.Run("異常系: GitHubUserIDが空の場合はエラー", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockTokenRepo := new(MockRefreshTokenRepository)

		uc := newTestAuthUseCase(mockUserRepo, mockTokenRepo)
		_, err := uc.GitHubOAuthLogin(ctx, GitHubOAuthInput{
			GitHubUserID: "",
			Email:        "test@example.com",
			Name:         "Test User",
		})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "GitHub user IDは必須です")
	})

	t.Run("異常系: メールアドレスが空の場合はエラー", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockTokenRepo := new(MockRefreshTokenRepository)

		uc := newTestAuthUseCase(mockUserRepo, mockTokenRepo)
		_, err := uc.GitHubOAuthLogin(ctx, GitHubOAuthInput{
			GitHubUserID: "github-123",
			Email:        "",
			Name:         "Test User",
		})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "メールアドレスは必須です")
	})

	t.Run("正常系: 既存のGitHubユーザーでログインできる", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockTokenRepo := new(MockRefreshTokenRepository)
		user := newTestUser("user-001", "github@example.com")
		mockUserRepo.On("FindByProviderUserID", mock_anything(), entities.AuthProviderGitHub, "github-123").Return(user, nil)
		mockTokenRepo.On("Save", mock_anything(), mock_anything()).Return(nil)

		uc := newTestAuthUseCase(mockUserRepo, mockTokenRepo)
		output, err := uc.GitHubOAuthLogin(ctx, GitHubOAuthInput{
			GitHubUserID: "github-123",
			Email:        "github@example.com",
			Name:         "GitHub User",
		})

		require.NoError(t, err)
		assert.NotEmpty(t, output.Token)
		mockUserRepo.AssertExpectations(t)
		mockTokenRepo.AssertExpectations(t)
	})

	t.Run("異常系: 同一メールアドレスの既存アカウントが存在する場合はエラー", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockTokenRepo := new(MockRefreshTokenRepository)
		email, _ := entities.NewEmail("existing@example.com")
		existingUser := newTestUser("user-existing", "existing@example.com")
		mockUserRepo.On("FindByProviderUserID", mock_anything(), entities.AuthProviderGitHub, "github-new").Return(nil, errors.New("not found"))
		mockUserRepo.On("FindByEmail", mock_anything(), email).Return(existingUser, nil)

		uc := newTestAuthUseCase(mockUserRepo, mockTokenRepo)
		_, err := uc.GitHubOAuthLogin(ctx, GitHubOAuthInput{
			GitHubUserID: "github-new",
			Email:        "existing@example.com",
			Name:         "New User",
		})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "既に登録されています")
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("正常系: 新規GitHubユーザーを作成してログインできる", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockTokenRepo := new(MockRefreshTokenRepository)
		email, _ := entities.NewEmail("new@example.com")
		mockUserRepo.On("FindByProviderUserID", mock_anything(), entities.AuthProviderGitHub, "github-brand-new").Return(nil, errors.New("not found"))
		mockUserRepo.On("FindByEmail", mock_anything(), email).Return(nil, errors.New("not found"))
		mockUserRepo.On("Save", mock_anything(), mock_anything()).Return(nil)
		mockTokenRepo.On("Save", mock_anything(), mock_anything()).Return(nil)

		uc := newTestAuthUseCase(mockUserRepo, mockTokenRepo)
		output, err := uc.GitHubOAuthLogin(ctx, GitHubOAuthInput{
			GitHubUserID: "github-brand-new",
			Email:        "new@example.com",
			Name:         "Brand New User",
		})

		require.NoError(t, err)
		assert.NotEmpty(t, output.Token)
		mockUserRepo.AssertExpectations(t)
		mockTokenRepo.AssertExpectations(t)
	})
}

// ===========================
// RefreshAccessToken Tests
// ===========================

func TestAuthUseCase_RefreshAccessToken(t *testing.T) {
	ctx := context.Background()

	t.Run("異常系: 無効なリフレッシュトークンの場合はエラー", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockTokenRepo := new(MockRefreshTokenRepository)
		mockTokenRepo.On("FindByTokenHash", mock_anything(), mock_anything()).Return(nil, errors.New("not found"))

		uc := newTestAuthUseCase(mockUserRepo, mockTokenRepo)
		_, err := uc.RefreshAccessToken(ctx, "invalid-refresh-token")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "無効なリフレッシュトークンです")
		mockTokenRepo.AssertExpectations(t)
	})
}
// ===========================
// Setup2FA Tests
// ===========================

func TestAuthUseCase_Setup2FA(t *testing.T) {
	ctx := context.Background()

	t.Run("正常系: 2FAセットアップデータを取得できる", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockTokenRepo := new(MockRefreshTokenRepository)
		user := newTestUser("user-001", "test@example.com")
		mockUserRepo.On("FindByID", mock_anything(), entities.UserID("user-001")).Return(user, nil)

		uc := newTestAuthUseCase(mockUserRepo, mockTokenRepo)
		output, err := uc.Setup2FA(ctx, "user-001")

		require.NoError(t, err)
		assert.NotEmpty(t, output.Secret)
		assert.NotEmpty(t, output.QRCodeURL)
		assert.Len(t, output.BackupCodes, 8)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("異常系: ユーザーが存在しない場合はエラー", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockTokenRepo := new(MockRefreshTokenRepository)
		mockUserRepo.On("FindByID", mock_anything(), entities.UserID("user-999")).Return(nil, errors.New("not found"))

		uc := newTestAuthUseCase(mockUserRepo, mockTokenRepo)
		_, err := uc.Setup2FA(ctx, "user-999")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "ユーザーが見つかりません")
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("異常系: 無効なユーザーIDの場合はエラー", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockTokenRepo := new(MockRefreshTokenRepository)

		uc := newTestAuthUseCase(mockUserRepo, mockTokenRepo)
		_, err := uc.Setup2FA(ctx, "")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "無効なユーザーID")
	})
}

// ===========================
// Enable2FA Tests
// ===========================

func TestAuthUseCase_Enable2FA(t *testing.T) {
	ctx := context.Background()

	t.Run("異常系: ユーザーIDが空の場合はエラー", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockTokenRepo := new(MockRefreshTokenRepository)

		uc := newTestAuthUseCase(mockUserRepo, mockTokenRepo)
		err := uc.Enable2FA(ctx, Enable2FAInput{
			UserID: "",
			Code:   "123456",
			Secret: "TESTSECRET",
		})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "ユーザーIDは必須です")
	})

	t.Run("異常系: 認証コードが空の場合はエラー", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockTokenRepo := new(MockRefreshTokenRepository)

		uc := newTestAuthUseCase(mockUserRepo, mockTokenRepo)
		err := uc.Enable2FA(ctx, Enable2FAInput{
			UserID: "user-001",
			Code:   "",
			Secret: "TESTSECRET",
		})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "認証コードは必須です")
	})

	t.Run("異常系: シークレットが空の場合はエラー", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockTokenRepo := new(MockRefreshTokenRepository)

		uc := newTestAuthUseCase(mockUserRepo, mockTokenRepo)
		err := uc.Enable2FA(ctx, Enable2FAInput{
			UserID: "user-001",
			Code:   "123456",
			Secret: "",
		})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "シークレットは必須です")
	})

	t.Run("異常系: ユーザーが存在しない場合はエラー", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockTokenRepo := new(MockRefreshTokenRepository)
		mockUserRepo.On("FindByID", mock_anything(), entities.UserID("user-999")).Return(nil, errors.New("not found"))

		uc := newTestAuthUseCase(mockUserRepo, mockTokenRepo)
		err := uc.Enable2FA(ctx, Enable2FAInput{
			UserID: "user-999",
			Code:   "123456",
			Secret: "TESTSECRET",
		})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "ユーザーが見つかりません")
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("異常系: 無効なTOTPコードの場合はエラー", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockTokenRepo := new(MockRefreshTokenRepository)
		user := newTestUser("user-001", "test@example.com")
		mockUserRepo.On("FindByID", mock_anything(), entities.UserID("user-001")).Return(user, nil)

		uc := newTestAuthUseCase(mockUserRepo, mockTokenRepo)
		err := uc.Enable2FA(ctx, Enable2FAInput{
			UserID: "user-001",
			Code:   "000000",     // 無効なコード
			Secret: "TESTSECRET", // 無効なシークレット
		})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "認証コードが無効です")
		mockUserRepo.AssertExpectations(t)
	})
}

// ===========================
// Verify2FA Tests
// ===========================

func TestAuthUseCase_Verify2FA(t *testing.T) {
	ctx := context.Background()

	t.Run("異常系: ユーザーIDが空の場合はエラー", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockTokenRepo := new(MockRefreshTokenRepository)

		uc := newTestAuthUseCase(mockUserRepo, mockTokenRepo)
		_, err := uc.Verify2FA(ctx, Verify2FAInput{
			UserID: "",
			Code:   "123456",
		})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "ユーザーIDは必須です")
	})

	t.Run("異常系: 認証コードが空の場合はエラー", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockTokenRepo := new(MockRefreshTokenRepository)

		uc := newTestAuthUseCase(mockUserRepo, mockTokenRepo)
		_, err := uc.Verify2FA(ctx, Verify2FAInput{
			UserID: "user-001",
			Code:   "",
		})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "認証コードは必須です")
	})

	t.Run("異常系: ユーザーが存在しない場合はエラー", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockTokenRepo := new(MockRefreshTokenRepository)
		mockUserRepo.On("FindByID", mock_anything(), entities.UserID("user-999")).Return(nil, errors.New("not found"))

		uc := newTestAuthUseCase(mockUserRepo, mockTokenRepo)
		_, err := uc.Verify2FA(ctx, Verify2FAInput{
			UserID: "user-999",
			Code:   "123456",
		})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "ユーザーが見つかりません")
		mockUserRepo.AssertExpectations(t)
	})
}

// ===========================
// Disable2FA Tests
// ===========================

func TestAuthUseCase_Disable2FA(t *testing.T) {
	ctx := context.Background()

	t.Run("異常系: ユーザーIDが空の場合はエラー", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockTokenRepo := new(MockRefreshTokenRepository)

		uc := newTestAuthUseCase(mockUserRepo, mockTokenRepo)
		err := uc.Disable2FA(ctx, Disable2FAInput{
			UserID:   "",
			Password: "password123",
		})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "ユーザーIDは必須です")
	})

	t.Run("異常系: パスワードが空の場合はエラー", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockTokenRepo := new(MockRefreshTokenRepository)

		uc := newTestAuthUseCase(mockUserRepo, mockTokenRepo)
		err := uc.Disable2FA(ctx, Disable2FAInput{
			UserID:   "user-001",
			Password: "",
		})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "パスワードは必須です")
	})

	t.Run("異常系: ユーザーが存在しない場合はエラー", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockTokenRepo := new(MockRefreshTokenRepository)
		mockUserRepo.On("FindByID", mock_anything(), entities.UserID("user-999")).Return(nil, errors.New("not found"))

		uc := newTestAuthUseCase(mockUserRepo, mockTokenRepo)
		err := uc.Disable2FA(ctx, Disable2FAInput{
			UserID:   "user-999",
			Password: "password123",
		})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "ユーザーが見つかりません")
		mockUserRepo.AssertExpectations(t)
	})
}

// ===========================
// RegenerateBackupCodes Tests
// ===========================

func TestAuthUseCase_RegenerateBackupCodes(t *testing.T) {
	ctx := context.Background()

	t.Run("異常系: 無効なユーザーIDの場合はエラー", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockTokenRepo := new(MockRefreshTokenRepository)

		uc := newTestAuthUseCase(mockUserRepo, mockTokenRepo)
		_, err := uc.RegenerateBackupCodes(ctx, "")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "無効なユーザーID")
	})

	t.Run("異常系: ユーザーが存在しない場合はエラー", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockTokenRepo := new(MockRefreshTokenRepository)
		mockUserRepo.On("FindByID", mock_anything(), entities.UserID("user-999")).Return(nil, errors.New("not found"))

		uc := newTestAuthUseCase(mockUserRepo, mockTokenRepo)
		_, err := uc.RegenerateBackupCodes(ctx, "user-999")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "ユーザーが見つかりません")
		mockUserRepo.AssertExpectations(t)
	})
}

// ===========================
// VerifyToken Additional Tests
// ===========================

func TestAuthUseCase_VerifyToken_Invalid(t *testing.T) {
	ctx := context.Background()

	t.Run("異常系: 無効なJWTトークンの場合はエラー", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockTokenRepo := new(MockRefreshTokenRepository)

		uc := newTestAuthUseCase(mockUserRepo, mockTokenRepo)
		_, err := uc.VerifyToken(ctx, "invalid.jwt.token")

		require.Error(t, err)
	})

	t.Run("異常系: 空のトークンの場合はエラー", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockTokenRepo := new(MockRefreshTokenRepository)

		uc := newTestAuthUseCase(mockUserRepo, mockTokenRepo)
		_, err := uc.VerifyToken(ctx, "")

		require.Error(t, err)
	})
}