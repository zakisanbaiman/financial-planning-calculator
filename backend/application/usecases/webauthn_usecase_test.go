package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/financial-planning-calculator/backend/domain/entities"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// newTestWebAuthnConfig はテスト用のWebAuthn設定を作成する
func newTestWebAuthnConfig() *webauthn.WebAuthn {
	wauthn, err := webauthn.New(&webauthn.Config{
		RPDisplayName: "Test App",
		RPID:          "localhost",
		RPOrigins:     []string{"http://localhost:3000"},
	})
	if err != nil {
		panic(err)
	}
	return wauthn
}

// newTestWebAuthnUseCase はテスト用のWebAuthnUseCaseを作成する
func newTestWebAuthnUseCase(
	userRepo *MockUserRepository,
	credRepo *MockWebAuthnCredentialRepository,
	rtRepo *MockRefreshTokenRepository,
) WebAuthnUseCase {
	mockAuthUC := new(MockAuthUseCase)
	return NewWebAuthnUseCase(
		userRepo,
		credRepo,
		rtRepo,
		newTestWebAuthnConfig(),
		mockAuthUC,
		testJWTSecret,
		testJWTExpiry,
		testRTExpiry,
	)
}

// MockAuthUseCase はAuthUseCaseのモック実装（WebAuthnテスト用）
type MockAuthUseCase struct {
	mock.Mock
}

func (m *MockAuthUseCase) Register(ctx context.Context, input RegisterInput) (*RegisterOutput, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*RegisterOutput), args.Error(1)
}

func (m *MockAuthUseCase) Login(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*LoginOutput), args.Error(1)
}

func (m *MockAuthUseCase) VerifyToken(ctx context.Context, tokenString string) (*TokenClaims, error) {
	args := m.Called(ctx, tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*TokenClaims), args.Error(1)
}

func (m *MockAuthUseCase) RefreshAccessToken(ctx context.Context, refreshToken string) (*RefreshOutput, error) {
	args := m.Called(ctx, refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*RefreshOutput), args.Error(1)
}

func (m *MockAuthUseCase) RevokeRefreshToken(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockAuthUseCase) GitHubOAuthLogin(ctx context.Context, input GitHubOAuthInput) (*LoginOutput, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*LoginOutput), args.Error(1)
}

func (m *MockAuthUseCase) Setup2FA(ctx context.Context, userID string) (*Setup2FAOutput, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Setup2FAOutput), args.Error(1)
}

func (m *MockAuthUseCase) Enable2FA(ctx context.Context, input Enable2FAInput) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}

func (m *MockAuthUseCase) Verify2FA(ctx context.Context, input Verify2FAInput) (*LoginOutput, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*LoginOutput), args.Error(1)
}

func (m *MockAuthUseCase) Disable2FA(ctx context.Context, input Disable2FAInput) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}

func (m *MockAuthUseCase) RegenerateBackupCodes(ctx context.Context, userID string) (*RegenerateBackupCodesOutput, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*RegenerateBackupCodesOutput), args.Error(1)
}

func (m *MockAuthUseCase) Get2FAStatus(ctx context.Context, userID string) (*Get2FAStatusOutput, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Get2FAStatusOutput), args.Error(1)
}

func TestWebAuthnUseCase_ListCredentials(t *testing.T) {
	tests := []struct {
		name        string
		userID      string
		setupMock   func(*MockUserRepository, *MockWebAuthnCredentialRepository, *MockRefreshTokenRepository)
		expectError bool
		errContains string
		expectCount int
	}{
		{
			name:   "正常系: クレデンシャル一覧取得（空）",
			userID: "valid-user-id",
			setupMock: func(ur *MockUserRepository, cr *MockWebAuthnCredentialRepository, rtr *MockRefreshTokenRepository) {
				cr.On("FindByUserID", mock.Anything, entities.UserID("valid-user-id")).
					Return([]*entities.WebAuthnCredential{}, nil)
			},
			expectError: false,
			expectCount: 0,
		},
		{
			name:   "正常系: クレデンシャル一覧取得（1件）",
			userID: "user-with-credentials",
			setupMock: func(ur *MockUserRepository, cr *MockWebAuthnCredentialRepository, rtr *MockRefreshTokenRepository) {
				cred, _ := entities.NewWebAuthnCredential(
					"cred-id-001",
					entities.UserID("user-with-credentials"),
					[]byte("credential-id-bytes"),
					[]byte("public-key-bytes"),
					"none",
					[]byte("aaguid"),
					[]string{"internal"},
					"My Passkey",
				)
				cr.On("FindByUserID", mock.Anything, entities.UserID("user-with-credentials")).
					Return([]*entities.WebAuthnCredential{cred}, nil)
			},
			expectError: false,
			expectCount: 1,
		},
		{
			name:        "異常系: 空のユーザーID",
			userID:      "",
			setupMock:   func(ur *MockUserRepository, cr *MockWebAuthnCredentialRepository, rtr *MockRefreshTokenRepository) {},
			expectError: true,
			errContains: "無効なユーザーID",
		},
		{
			name:   "異常系: リポジトリエラー",
			userID: "user-db-error",
			setupMock: func(ur *MockUserRepository, cr *MockWebAuthnCredentialRepository, rtr *MockRefreshTokenRepository) {
				cr.On("FindByUserID", mock.Anything, entities.UserID("user-db-error")).
					Return(nil, errors.New("DBエラー"))
			},
			expectError: true,
			errContains: "クレデンシャルの取得に失敗",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := new(MockUserRepository)
			credRepo := new(MockWebAuthnCredentialRepository)
			rtRepo := new(MockRefreshTokenRepository)
			tt.setupMock(userRepo, credRepo, rtRepo)

			uc := newTestWebAuthnUseCase(userRepo, credRepo, rtRepo)
			result, err := uc.ListCredentials(context.Background(), tt.userID)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectCount, len(result))
			}
			credRepo.AssertExpectations(t)
		})
	}
}

func TestWebAuthnUseCase_DeleteCredential(t *testing.T) {
	tests := []struct {
		name         string
		userID       string
		credentialID string
		setupMock    func(*MockUserRepository, *MockWebAuthnCredentialRepository, *MockRefreshTokenRepository)
		expectError  bool
		errContains  string
	}{
		{
			name:         "正常系: クレデンシャル削除",
			userID:       "user-001",
			credentialID: "cred-id-001",
			setupMock: func(ur *MockUserRepository, cr *MockWebAuthnCredentialRepository, rtr *MockRefreshTokenRepository) {
				cred, _ := entities.NewWebAuthnCredential(
					"cred-id-001",
					entities.UserID("user-001"),
					[]byte("credential-id"),
					[]byte("public-key"),
					"none",
					[]byte("aaguid"),
					[]string{"internal"},
					"My Passkey",
				)
				cr.On("FindByID", mock.Anything, entities.CredentialID("cred-id-001")).Return(cred, nil)
				cr.On("Delete", mock.Anything, entities.CredentialID("cred-id-001")).Return(nil)
			},
			expectError: false,
		},
		{
			name:         "異常系: 空のユーザーID",
			userID:       "",
			credentialID: "cred-id-001",
			setupMock:    func(ur *MockUserRepository, cr *MockWebAuthnCredentialRepository, rtr *MockRefreshTokenRepository) {},
			expectError:  true,
			errContains:  "無効なユーザーID",
		},
		{
			name:         "異常系: クレデンシャルが存在しない",
			userID:       "user-001",
			credentialID: "cred-id-999",
			setupMock: func(ur *MockUserRepository, cr *MockWebAuthnCredentialRepository, rtr *MockRefreshTokenRepository) {
				cr.On("FindByID", mock.Anything, entities.CredentialID("cred-id-999")).
					Return(nil, errors.New("クレデンシャルが見つかりません"))
			},
			expectError: true,
			errContains: "クレデンシャルが見つかりません",
		},
		{
			name:         "異常系: 別ユーザーのクレデンシャル削除",
			userID:       "other-user",
			credentialID: "cred-id-001",
			setupMock: func(ur *MockUserRepository, cr *MockWebAuthnCredentialRepository, rtr *MockRefreshTokenRepository) {
				cred, _ := entities.NewWebAuthnCredential(
					"cred-id-001",
					entities.UserID("user-001"),
					[]byte("credential-id"),
					[]byte("public-key"),
					"none",
					[]byte("aaguid"),
					[]string{"internal"},
					"My Passkey",
				)
				cr.On("FindByID", mock.Anything, entities.CredentialID("cred-id-001")).Return(cred, nil)
			},
			expectError: true,
			errContains: "このクレデンシャルの所有者ではありません",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := new(MockUserRepository)
			credRepo := new(MockWebAuthnCredentialRepository)
			rtRepo := new(MockRefreshTokenRepository)
			tt.setupMock(userRepo, credRepo, rtRepo)

			uc := newTestWebAuthnUseCase(userRepo, credRepo, rtRepo)
			err := uc.DeleteCredential(context.Background(), tt.userID, tt.credentialID)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
			}
			credRepo.AssertExpectations(t)
		})
	}
}

func TestWebAuthnUseCase_RenameCredential(t *testing.T) {
	tests := []struct {
		name         string
		userID       string
		credentialID string
		newName      string
		setupMock    func(*MockUserRepository, *MockWebAuthnCredentialRepository, *MockRefreshTokenRepository)
		expectError  bool
		errContains  string
	}{
		{
			name:         "正常系: クレデンシャル名変更",
			userID:       "user-001",
			credentialID: "cred-id-001",
			newName:      "New Passkey Name",
			setupMock: func(ur *MockUserRepository, cr *MockWebAuthnCredentialRepository, rtr *MockRefreshTokenRepository) {
				cred, _ := entities.NewWebAuthnCredential(
					"cred-id-001",
					entities.UserID("user-001"),
					[]byte("credential-id"),
					[]byte("public-key"),
					"none",
					[]byte("aaguid"),
					[]string{"internal"},
					"Old Name",
				)
				cr.On("FindByID", mock.Anything, entities.CredentialID("cred-id-001")).Return(cred, nil)
				cr.On("Update", mock.Anything, mock.AnythingOfType("*entities.WebAuthnCredential")).Return(nil)
			},
			expectError: false,
		},
		{
			name:         "異常系: クレデンシャルが存在しない",
			userID:       "user-001",
			credentialID: "cred-id-999",
			newName:      "New Name",
			setupMock: func(ur *MockUserRepository, cr *MockWebAuthnCredentialRepository, rtr *MockRefreshTokenRepository) {
				cr.On("FindByID", mock.Anything, entities.CredentialID("cred-id-999")).
					Return(nil, errors.New("クレデンシャルが見つかりません"))
			},
			expectError: true,
			errContains: "クレデンシャルが見つかりません",
		},
		{
			name:         "異常系: 別ユーザーのクレデンシャル名変更",
			userID:       "other-user",
			credentialID: "cred-id-001",
			newName:      "New Name",
			setupMock: func(ur *MockUserRepository, cr *MockWebAuthnCredentialRepository, rtr *MockRefreshTokenRepository) {
				cred, _ := entities.NewWebAuthnCredential(
					"cred-id-001",
					entities.UserID("user-001"),
					[]byte("credential-id"),
					[]byte("public-key"),
					"none",
					[]byte("aaguid"),
					[]string{"internal"},
					"Old Name",
				)
				cr.On("FindByID", mock.Anything, entities.CredentialID("cred-id-001")).Return(cred, nil)
			},
			expectError: true,
			errContains: "このクレデンシャルの所有者ではありません",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := new(MockUserRepository)
			credRepo := new(MockWebAuthnCredentialRepository)
			rtRepo := new(MockRefreshTokenRepository)
			tt.setupMock(userRepo, credRepo, rtRepo)

			uc := newTestWebAuthnUseCase(userRepo, credRepo, rtRepo)
			err := uc.RenameCredential(context.Background(), tt.userID, tt.credentialID, tt.newName)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
			}
			credRepo.AssertExpectations(t)
		})
	}
}

func TestWebAuthnUseCase_BeginRegistration(t *testing.T) {
	tests := []struct {
		name        string
		userID      string
		setupMock   func(*MockUserRepository, *MockWebAuthnCredentialRepository, *MockRefreshTokenRepository)
		expectError bool
		errContains string
	}{
		{
			name:   "正常系: パスキー登録開始",
			userID: "valid-user-id",
			setupMock: func(ur *MockUserRepository, cr *MockWebAuthnCredentialRepository, rtr *MockRefreshTokenRepository) {
				user := newTestUser("passkey@example.com", "Password123!")
				ur.On("FindByID", mock.Anything, entities.UserID("valid-user-id")).Return(user, nil)
				cr.On("FindByUserID", mock.Anything, entities.UserID("valid-user-id")).
					Return([]*entities.WebAuthnCredential{}, nil)
			},
			expectError: false,
		},
		{
			name:        "異常系: 空のユーザーID",
			userID:      "",
			setupMock:   func(ur *MockUserRepository, cr *MockWebAuthnCredentialRepository, rtr *MockRefreshTokenRepository) {},
			expectError: true,
			errContains: "無効なユーザーID",
		},
		{
			name:   "異常系: ユーザーが存在しない",
			userID: "nonexistent-user",
			setupMock: func(ur *MockUserRepository, cr *MockWebAuthnCredentialRepository, rtr *MockRefreshTokenRepository) {
				ur.On("FindByID", mock.Anything, entities.UserID("nonexistent-user")).
					Return(nil, errors.New("ユーザーが見つかりません"))
			},
			expectError: true,
			errContains: "ユーザーが見つかりません",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := new(MockUserRepository)
			credRepo := new(MockWebAuthnCredentialRepository)
			rtRepo := new(MockRefreshTokenRepository)
			tt.setupMock(userRepo, credRepo, rtRepo)

			uc := newTestWebAuthnUseCase(userRepo, credRepo, rtRepo)
			output, err := uc.BeginRegistration(context.Background(), tt.userID)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, output)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, output)
				assert.NotEmpty(t, output.PublicKeyOptions)
				assert.NotEmpty(t, output.SessionData)
			}
			userRepo.AssertExpectations(t)
			credRepo.AssertExpectations(t)
		})
	}
}

func TestWebAuthnUseCase_BeginLogin(t *testing.T) {
t.Run("正常系: パスキーログイン開始", func(t *testing.T) {
userRepo := new(MockUserRepository)
credRepo := new(MockWebAuthnCredentialRepository)
rtRepo := new(MockRefreshTokenRepository)

uc := newTestWebAuthnUseCase(userRepo, credRepo, rtRepo)
output, err := uc.BeginLogin(context.Background(), BeginLoginInput{
Email: "",
})

require.NoError(t, err)
assert.NotNil(t, output)
assert.NotEmpty(t, output.PublicKeyOptions)
assert.NotEmpty(t, output.SessionData)
})
}
