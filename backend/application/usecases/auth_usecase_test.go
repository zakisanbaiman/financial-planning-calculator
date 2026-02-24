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
		mockUserRepo.On("ExistsByEmail", ctx, email).Return(false, nil)
		mockUserRepo.On("Save", ctx, mock_anything()).Return(nil)
		mockTokenRepo.On("Save", ctx, mock_anything()).Return(nil)

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
		mockUserRepo.On("ExistsByEmail", ctx, email).Return(true, nil)

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
		mockUserRepo.On("ExistsByEmail", ctx, email).Return(false, errors.New("db error"))

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
		mockUserRepo.On("FindByEmail", ctx, email).Return(nil, errors.New("not found"))

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
		mockTokenRepo.On("RevokeByUserID", ctx, entities.UserID("user-001")).Return(nil)

		uc := newTestAuthUseCase(mockUserRepo, mockTokenRepo)
		err := uc.RevokeRefreshToken(ctx, "user-001")

		require.NoError(t, err)
		mockTokenRepo.AssertExpectations(t)
	})

	t.Run("異常系: リポジトリエラーの場合はエラーを返す", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockTokenRepo := new(MockRefreshTokenRepository)
		mockTokenRepo.On("RevokeByUserID", ctx, entities.UserID("user-001")).Return(errors.New("db error"))

		uc := newTestAuthUseCase(mockUserRepo, mockTokenRepo)
		err := uc.RevokeRefreshToken(ctx, "user-001")

		require.Error(t, err)
		mockTokenRepo.AssertExpectations(t)
	})
}
