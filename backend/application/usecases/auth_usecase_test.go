package usecases

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/financial-planning-calculator/backend/domain/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	testJWTSecret  = "test-secret-key-for-unit-tests"
	testJWTExpiry  = 15 * time.Minute
	testRTExpiry   = 7 * 24 * time.Hour
)

func newTestAuthUseCase(userRepo *MockUserRepository, rtRepo *MockRefreshTokenRepository) AuthUseCase {
	return NewAuthUseCase(userRepo, rtRepo, testJWTSecret, testJWTExpiry, testRTExpiry)
}

func newTestUser(email, password string) *entities.User {
	user, err := entities.NewUser("test-user-id", email, password)
	if err != nil {
		panic(err)
	}
	return user
}

func TestAuthUseCase_Register(t *testing.T) {
	tests := []struct {
		name        string
		input       RegisterInput
		setupMock   func(*MockUserRepository, *MockRefreshTokenRepository)
		expectError bool
		errContains string
	}{
		{
			name: "正常系: ユーザー登録",
			input: RegisterInput{
				Email:    "newuser@example.com",
				Password: "Password123!",
			},
			setupMock: func(ur *MockUserRepository, rtr *MockRefreshTokenRepository) {
				email, _ := entities.NewEmail("newuser@example.com")
				ur.On("ExistsByEmail", mock.Anything, email).Return(false, nil)
				ur.On("Save", mock.Anything, mock.AnythingOfType("*entities.User")).Return(nil)
				rtr.On("Save", mock.Anything, mock.AnythingOfType("*entities.RefreshToken")).Return(nil)
			},
			expectError: false,
		},
		{
			name: "異常系: メールアドレスが空",
			input: RegisterInput{
				Email:    "",
				Password: "Password123!",
			},
			setupMock:   func(ur *MockUserRepository, rtr *MockRefreshTokenRepository) {},
			expectError: true,
			errContains: "メールアドレスは必須です",
		},
		{
			name: "異常系: パスワードが空",
			input: RegisterInput{
				Email:    "test@example.com",
				Password: "",
			},
			setupMock:   func(ur *MockUserRepository, rtr *MockRefreshTokenRepository) {},
			expectError: true,
			errContains: "パスワードは必須です",
		},
		{
			name: "異常系: 無効なメールアドレス",
			input: RegisterInput{
				Email:    "invalid-email",
				Password: "Password123!",
			},
			setupMock:   func(ur *MockUserRepository, rtr *MockRefreshTokenRepository) {},
			expectError: true,
			errContains: "無効なメールアドレスです",
		},
		{
			name: "異常系: メールアドレスが既に使用されている",
			input: RegisterInput{
				Email:    "existing@example.com",
				Password: "Password123!",
			},
			setupMock: func(ur *MockUserRepository, rtr *MockRefreshTokenRepository) {
				email, _ := entities.NewEmail("existing@example.com")
				ur.On("ExistsByEmail", mock.Anything, email).Return(true, nil)
			},
			expectError: true,
			errContains: "このメールアドレスは既に登録されています",
		},
		{
			name: "異常系: ExistsByEmailの失敗",
			input: RegisterInput{
				Email:    "check@example.com",
				Password: "Password123!",
			},
			setupMock: func(ur *MockUserRepository, rtr *MockRefreshTokenRepository) {
				email, _ := entities.NewEmail("check@example.com")
				ur.On("ExistsByEmail", mock.Anything, email).Return(false, errors.New("DBエラー"))
			},
			expectError: true,
			errContains: "メールアドレスの確認に失敗しました",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := new(MockUserRepository)
			rtRepo := new(MockRefreshTokenRepository)
			tt.setupMock(userRepo, rtRepo)

			uc := newTestAuthUseCase(userRepo, rtRepo)
			output, err := uc.Register(context.Background(), tt.input)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, output)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, output)
				assert.NotEmpty(t, output.Token)
				assert.NotEmpty(t, output.RefreshToken)
				assert.Equal(t, tt.input.Email, output.Email)
			}
			userRepo.AssertExpectations(t)
			rtRepo.AssertExpectations(t)
		})
	}
}

func TestAuthUseCase_Login(t *testing.T) {
	tests := []struct {
		name        string
		input       LoginInput
		setupMock   func(*MockUserRepository, *MockRefreshTokenRepository)
		expectError bool
		errContains string
	}{
		{
			name: "正常系: ログイン成功",
			input: LoginInput{
				Email:    "user@example.com",
				Password: "Password123!",
			},
			setupMock: func(ur *MockUserRepository, rtr *MockRefreshTokenRepository) {
				user := newTestUser("user@example.com", "Password123!")
				email, _ := entities.NewEmail("user@example.com")
				ur.On("FindByEmail", mock.Anything, email).Return(user, nil)
				rtr.On("Save", mock.Anything, mock.AnythingOfType("*entities.RefreshToken")).Return(nil)
			},
			expectError: false,
		},
		{
			name: "異常系: メールアドレスが空",
			input: LoginInput{
				Email:    "",
				Password: "Password123!",
			},
			setupMock:   func(ur *MockUserRepository, rtr *MockRefreshTokenRepository) {},
			expectError: true,
			errContains: "メールアドレスは必須です",
		},
		{
			name: "異常系: パスワードが空",
			input: LoginInput{
				Email:    "user@example.com",
				Password: "",
			},
			setupMock:   func(ur *MockUserRepository, rtr *MockRefreshTokenRepository) {},
			expectError: true,
			errContains: "パスワードは必須です",
		},
		{
			name: "異常系: ユーザーが存在しない",
			input: LoginInput{
				Email:    "notfound@example.com",
				Password: "Password123!",
			},
			setupMock: func(ur *MockUserRepository, rtr *MockRefreshTokenRepository) {
				email, _ := entities.NewEmail("notfound@example.com")
				ur.On("FindByEmail", mock.Anything, email).Return(nil, errors.New("ユーザーが見つかりません"))
			},
			expectError: true,
			errContains: "メールアドレスまたはパスワードが正しくありません",
		},
		{
			name: "異常系: パスワードが不正",
			input: LoginInput{
				Email:    "user@example.com",
				Password: "WrongPassword!",
			},
			setupMock: func(ur *MockUserRepository, rtr *MockRefreshTokenRepository) {
				user := newTestUser("user@example.com", "Password123!")
				email, _ := entities.NewEmail("user@example.com")
				ur.On("FindByEmail", mock.Anything, email).Return(user, nil)
			},
			expectError: true,
			errContains: "メールアドレスまたはパスワードが正しくありません",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := new(MockUserRepository)
			rtRepo := new(MockRefreshTokenRepository)
			tt.setupMock(userRepo, rtRepo)

			uc := newTestAuthUseCase(userRepo, rtRepo)
			output, err := uc.Login(context.Background(), tt.input)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, output)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, output)
				assert.NotEmpty(t, output.Token)
			}
			userRepo.AssertExpectations(t)
			rtRepo.AssertExpectations(t)
		})
	}
}

func TestAuthUseCase_VerifyToken(t *testing.T) {
	userRepo := new(MockUserRepository)
	rtRepo := new(MockRefreshTokenRepository)
	uc := newTestAuthUseCase(userRepo, rtRepo)

	// 有効なトークンを生成するために先にRegisterする
	email, _ := entities.NewEmail("token@example.com")
	userRepo.On("ExistsByEmail", mock.Anything, email).Return(false, nil)
	userRepo.On("Save", mock.Anything, mock.AnythingOfType("*entities.User")).Return(nil)
	rtRepo.On("Save", mock.Anything, mock.AnythingOfType("*entities.RefreshToken")).Return(nil)

	out, err := uc.Register(context.Background(), RegisterInput{
		Email:    "token@example.com",
		Password: "Password123!",
	})
	require.NoError(t, err)

	t.Run("正常系: トークン検証成功", func(t *testing.T) {
		claims, err := uc.VerifyToken(context.Background(), out.Token)
		require.NoError(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, "token@example.com", claims.Email)
	})

	t.Run("異常系: 無効なトークン", func(t *testing.T) {
		_, err := uc.VerifyToken(context.Background(), "invalid.token.here")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "トークンの検証に失敗しました")
	})

	t.Run("異常系: 空のトークン", func(t *testing.T) {
		_, err := uc.VerifyToken(context.Background(), "")
		assert.Error(t, err)
	})
}

func TestAuthUseCase_RevokeRefreshToken(t *testing.T) {
	tests := []struct {
		name        string
		userID      string
		setupMock   func(*MockUserRepository, *MockRefreshTokenRepository)
		expectError bool
		errContains string
	}{
		{
			name:   "正常系: リフレッシュトークン失効",
			userID: "valid-user-id",
			setupMock: func(ur *MockUserRepository, rtr *MockRefreshTokenRepository) {
				rtr.On("RevokeByUserID", mock.Anything, entities.UserID("valid-user-id")).Return(nil)
			},
			expectError: false,
		},
		{
			name:        "異常系: 空のユーザーID",
			userID:      "",
			setupMock:   func(ur *MockUserRepository, rtr *MockRefreshTokenRepository) {},
			expectError: true,
			errContains: "無効なユーザーIDです",
		},
		{
			name:   "異常系: リポジトリエラー",
			userID: "user-id-db-error",
			setupMock: func(ur *MockUserRepository, rtr *MockRefreshTokenRepository) {
				rtr.On("RevokeByUserID", mock.Anything, entities.UserID("user-id-db-error")).Return(errors.New("DBエラー"))
			},
			expectError: true,
			errContains: "リフレッシュトークンの失効に失敗しました",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := new(MockUserRepository)
			rtRepo := new(MockRefreshTokenRepository)
			tt.setupMock(userRepo, rtRepo)

			uc := newTestAuthUseCase(userRepo, rtRepo)
			err := uc.RevokeRefreshToken(context.Background(), tt.userID)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
			}
			userRepo.AssertExpectations(t)
			rtRepo.AssertExpectations(t)
		})
	}
}

func TestAuthUseCase_GitHubOAuthLogin(t *testing.T) {
	tests := []struct {
		name        string
		input       GitHubOAuthInput
		setupMock   func(*MockUserRepository, *MockRefreshTokenRepository)
		expectError bool
		errContains string
	}{
		{
			name: "正常系: 既存GitHubユーザーでログイン",
			input: GitHubOAuthInput{
				GitHubUserID: "github-123",
				Email:        "ghuser@example.com",
				Name:         "GitHub User",
				AvatarURL:    "https://example.com/avatar.png",
			},
			setupMock: func(ur *MockUserRepository, rtr *MockRefreshTokenRepository) {
				existingUser := newTestUser("ghuser@example.com", "Password123!")
				ur.On("FindByProviderUserID", mock.Anything, entities.AuthProviderGitHub, "github-123").
					Return(existingUser, nil)
				rtr.On("Save", mock.Anything, mock.AnythingOfType("*entities.RefreshToken")).Return(nil)
			},
			expectError: false,
		},
		{
			name: "異常系: GitHub user IDが空",
			input: GitHubOAuthInput{
				GitHubUserID: "",
				Email:        "ghuser@example.com",
			},
			setupMock:   func(ur *MockUserRepository, rtr *MockRefreshTokenRepository) {},
			expectError: true,
			errContains: "GitHub user IDは必須です",
		},
		{
			name: "異常系: メールアドレスが空",
			input: GitHubOAuthInput{
				GitHubUserID: "github-456",
				Email:        "",
			},
			setupMock:   func(ur *MockUserRepository, rtr *MockRefreshTokenRepository) {},
			expectError: true,
			errContains: "メールアドレスは必須です",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := new(MockUserRepository)
			rtRepo := new(MockRefreshTokenRepository)
			tt.setupMock(userRepo, rtRepo)

			uc := newTestAuthUseCase(userRepo, rtRepo)
			output, err := uc.GitHubOAuthLogin(context.Background(), tt.input)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, output)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, output)
				assert.NotEmpty(t, output.Token)
			}
			userRepo.AssertExpectations(t)
			rtRepo.AssertExpectations(t)
		})
	}
}

func TestAuthUseCase_Setup2FA(t *testing.T) {
	tests := []struct {
		name        string
		userID      string
		setupMock   func(*MockUserRepository, *MockRefreshTokenRepository)
		expectError bool
		errContains string
	}{
		{
			name:   "正常系: 2FAセットアップ開始",
			userID: "valid-user-id",
			setupMock: func(ur *MockUserRepository, rtr *MockRefreshTokenRepository) {
				user := newTestUser("setup2fa@example.com", "Password123!")
				ur.On("FindByID", mock.Anything, entities.UserID("valid-user-id")).Return(user, nil)
			},
			expectError: false,
		},
		{
			name:   "異常系: ユーザーが存在しない",
			userID: "nonexistent-user",
			setupMock: func(ur *MockUserRepository, rtr *MockRefreshTokenRepository) {
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
			rtRepo := new(MockRefreshTokenRepository)
			tt.setupMock(userRepo, rtRepo)

			uc := newTestAuthUseCase(userRepo, rtRepo)
			output, err := uc.Setup2FA(context.Background(), tt.userID)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, output)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, output)
				assert.NotEmpty(t, output.Secret)
				assert.NotEmpty(t, output.BackupCodes)
			}
			userRepo.AssertExpectations(t)
		})
	}
}

func TestAuthUseCase_Get2FAStatus(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		setupMock      func(*MockUserRepository, *MockRefreshTokenRepository)
		expectError    bool
		expectEnabled  bool
	}{
		{
			name:   "正常系: 2FAが無効",
			userID: "user-no-2fa",
			setupMock: func(ur *MockUserRepository, rtr *MockRefreshTokenRepository) {
				user := newTestUser("no2fa@example.com", "Password123!")
				ur.On("FindByID", mock.Anything, entities.UserID("user-no-2fa")).Return(user, nil)
			},
			expectError:   false,
			expectEnabled: false,
		},
		{
			name:   "異常系: ユーザーが存在しない",
			userID: "nonexistent",
			setupMock: func(ur *MockUserRepository, rtr *MockRefreshTokenRepository) {
				ur.On("FindByID", mock.Anything, entities.UserID("nonexistent")).
					Return(nil, errors.New("ユーザーが見つかりません"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := new(MockUserRepository)
			rtRepo := new(MockRefreshTokenRepository)
			tt.setupMock(userRepo, rtRepo)

			uc := newTestAuthUseCase(userRepo, rtRepo)
			output, err := uc.Get2FAStatus(context.Background(), tt.userID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, output)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, tt.expectEnabled, output.Enabled)
			}
			userRepo.AssertExpectations(t)
		})
	}
}

func TestAuthUseCase_RefreshAccessToken(t *testing.T) {
// First register a user to get a refresh token
userRepo := new(MockUserRepository)
rtRepo := new(MockRefreshTokenRepository)

email, _ := entities.NewEmail("refresh@example.com")
userRepo.On("ExistsByEmail", mock.Anything, email).Return(false, nil)
userRepo.On("Save", mock.Anything, mock.AnythingOfType("*entities.User")).Return(nil)
rtRepo.On("Save", mock.Anything, mock.AnythingOfType("*entities.RefreshToken")).Return(nil)

uc := newTestAuthUseCase(userRepo, rtRepo)
regOut, err := uc.Register(context.Background(), RegisterInput{
Email:    "refresh@example.com",
Password: "Password123!",
})
require.NoError(t, err)

t.Run("正常系: トークンリフレッシュ成功", func(t *testing.T) {
// モックのRefreshTokenを作成
user := newTestUser("refresh@example.com", "Password123!")

// FindByTokenHash は呼ばれるが、返り値は動的なトークンに依存するため、
// テストのリフレッシュトークンを手動で作成
rt, plainToken, _ := entities.NewRefreshToken(user.ID(), time.Now().Add(7*24*time.Hour))
_ = rt

rtRepo2 := new(MockRefreshTokenRepository)
userRepo2 := new(MockUserRepository)

// plainTokenのハッシュを使ってマッチさせる
rtRepo2.On("FindByTokenHash", mock.Anything, mock.Anything).Return(rt, nil)
userRepo2.On("FindByID", mock.Anything, mock.Anything).Return(user, nil)
rtRepo2.On("Update", mock.Anything, mock.AnythingOfType("*entities.RefreshToken")).Return(nil)

uc2 := newTestAuthUseCase(userRepo2, rtRepo2)
out, err := uc2.RefreshAccessToken(context.Background(), plainToken)

require.NoError(t, err)
assert.NotNil(t, out)
assert.NotEmpty(t, out.Token)
rtRepo2.AssertExpectations(t)
})

t.Run("異常系: 無効なリフレッシュトークン", func(t *testing.T) {
rtRepo2 := new(MockRefreshTokenRepository)
userRepo2 := new(MockUserRepository)
rtRepo2.On("FindByTokenHash", mock.Anything, mock.Anything).
Return(nil, errors.New("トークンが見つかりません"))

uc2 := newTestAuthUseCase(userRepo2, rtRepo2)
_, err := uc2.RefreshAccessToken(context.Background(), "invalid-refresh-token")
assert.Error(t, err)
assert.Contains(t, err.Error(), "無効なリフレッシュトークンです")
})
_ = regOut
}

func TestAuthUseCase_GitHubOAuthLogin_NewUser(t *testing.T) {
t.Run("正常系: 新規GitHubユーザー作成", func(t *testing.T) {
userRepo := new(MockUserRepository)
rtRepo := new(MockRefreshTokenRepository)

// GitHubプロバイダーIDによる検索失敗（新規ユーザー）
userRepo.On("FindByProviderUserID", mock.Anything, entities.AuthProviderGitHub, "github-new-user").
Return(nil, errors.New("ユーザーが見つかりません"))

email, _ := entities.NewEmail("newgithub@example.com")
userRepo.On("FindByEmail", mock.Anything, email).
Return(nil, errors.New("ユーザーが見つかりません"))

userRepo.On("Save", mock.Anything, mock.AnythingOfType("*entities.User")).Return(nil)
rtRepo.On("Save", mock.Anything, mock.AnythingOfType("*entities.RefreshToken")).Return(nil)

uc := newTestAuthUseCase(userRepo, rtRepo)
output, err := uc.GitHubOAuthLogin(context.Background(), GitHubOAuthInput{
GitHubUserID: "github-new-user",
Email:        "newgithub@example.com",
Name:         "New GitHub User",
AvatarURL:    "https://example.com/avatar.png",
})

require.NoError(t, err)
assert.NotNil(t, output)
assert.NotEmpty(t, output.Token)
userRepo.AssertExpectations(t)
rtRepo.AssertExpectations(t)
})

t.Run("異常系: 同一メールアドレスの既存アカウント", func(t *testing.T) {
userRepo := new(MockUserRepository)
rtRepo := new(MockRefreshTokenRepository)

existingUser := newTestUser("existing@example.com", "Password123!")
email, _ := entities.NewEmail("existing@example.com")

userRepo.On("FindByProviderUserID", mock.Anything, entities.AuthProviderGitHub, "github-conflict").
Return(nil, errors.New("ユーザーが見つかりません"))
userRepo.On("FindByEmail", mock.Anything, email).Return(existingUser, nil)

uc := newTestAuthUseCase(userRepo, rtRepo)
_, err := uc.GitHubOAuthLogin(context.Background(), GitHubOAuthInput{
GitHubUserID: "github-conflict",
Email:        "existing@example.com",
})

assert.Error(t, err)
assert.Contains(t, err.Error(), "このメールアドレスは既に登録されています")
userRepo.AssertExpectations(t)
})
}

func TestAuthUseCase_Enable2FA(t *testing.T) {
tests := []struct {
name        string
input       Enable2FAInput
setupMock   func(*MockUserRepository, *MockRefreshTokenRepository)
expectError bool
errContains string
}{
{
name: "異常系: ユーザーが存在しない",
input: Enable2FAInput{
UserID: "nonexistent",
Code:   "123456",
Secret: "JBSWY3DPEHPK3PXP",
},
setupMock: func(ur *MockUserRepository, rtr *MockRefreshTokenRepository) {
ur.On("FindByID", mock.Anything, entities.UserID("nonexistent")).
Return(nil, errors.New("ユーザーが見つかりません"))
},
expectError: true,
errContains: "ユーザーが見つかりません",
},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
userRepo := new(MockUserRepository)
rtRepo := new(MockRefreshTokenRepository)
tt.setupMock(userRepo, rtRepo)

uc := newTestAuthUseCase(userRepo, rtRepo)
err := uc.Enable2FA(context.Background(), tt.input)

if tt.expectError {
assert.Error(t, err)
if tt.errContains != "" {
assert.Contains(t, err.Error(), tt.errContains)
}
} else {
require.NoError(t, err)
}
userRepo.AssertExpectations(t)
})
}
}

func TestAuthUseCase_Disable2FA(t *testing.T) {
tests := []struct {
name        string
input       Disable2FAInput
setupMock   func(*MockUserRepository, *MockRefreshTokenRepository)
expectError bool
errContains string
}{
{
name: "異常系: ユーザーが存在しない",
input: Disable2FAInput{
UserID:   "nonexistent",
Password: "Password123!",
},
setupMock: func(ur *MockUserRepository, rtr *MockRefreshTokenRepository) {
ur.On("FindByID", mock.Anything, entities.UserID("nonexistent")).
Return(nil, errors.New("ユーザーが見つかりません"))
},
expectError: true,
errContains: "ユーザーが見つかりません",
},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
userRepo := new(MockUserRepository)
rtRepo := new(MockRefreshTokenRepository)
tt.setupMock(userRepo, rtRepo)

uc := newTestAuthUseCase(userRepo, rtRepo)
err := uc.Disable2FA(context.Background(), tt.input)

if tt.expectError {
assert.Error(t, err)
if tt.errContains != "" {
assert.Contains(t, err.Error(), tt.errContains)
}
} else {
require.NoError(t, err)
}
userRepo.AssertExpectations(t)
})
}
}

func TestAuthUseCase_Verify2FA(t *testing.T) {
tests := []struct {
name        string
input       Verify2FAInput
setupMock   func(*MockUserRepository, *MockRefreshTokenRepository)
expectError bool
errContains string
}{
{
name: "異常系: ユーザーが存在しない",
input: Verify2FAInput{
UserID: "nonexistent",
Code:  "123456",
},
setupMock: func(ur *MockUserRepository, rtr *MockRefreshTokenRepository) {
ur.On("FindByID", mock.Anything, entities.UserID("nonexistent")).
Return(nil, errors.New("ユーザーが見つかりません"))
},
expectError: true,
errContains: "ユーザーが見つかりません",
},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
userRepo := new(MockUserRepository)
rtRepo := new(MockRefreshTokenRepository)
tt.setupMock(userRepo, rtRepo)

uc := newTestAuthUseCase(userRepo, rtRepo)
output, err := uc.Verify2FA(context.Background(), tt.input)

if tt.expectError {
assert.Error(t, err)
if tt.errContains != "" {
assert.Contains(t, err.Error(), tt.errContains)
}
assert.Nil(t, output)
} else {
require.NoError(t, err)
assert.NotNil(t, output)
}
userRepo.AssertExpectations(t)
})
}
}

func TestAuthUseCase_RegenerateBackupCodes(t *testing.T) {
tests := []struct {
name        string
userID      string
setupMock   func(*MockUserRepository, *MockRefreshTokenRepository)
expectError bool
errContains string
}{
{
name:   "異常系: ユーザーが存在しない",
userID: "nonexistent",
setupMock: func(ur *MockUserRepository, rtr *MockRefreshTokenRepository) {
ur.On("FindByID", mock.Anything, entities.UserID("nonexistent")).
Return(nil, errors.New("ユーザーが見つかりません"))
},
expectError: true,
errContains: "ユーザーが見つかりません",
},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
userRepo := new(MockUserRepository)
rtRepo := new(MockRefreshTokenRepository)
tt.setupMock(userRepo, rtRepo)

uc := newTestAuthUseCase(userRepo, rtRepo)
output, err := uc.RegenerateBackupCodes(context.Background(), tt.userID)

if tt.expectError {
assert.Error(t, err)
if tt.errContains != "" {
assert.Contains(t, err.Error(), tt.errContains)
}
assert.Nil(t, output)
} else {
require.NoError(t, err)
assert.NotNil(t, output)
assert.NotEmpty(t, output.BackupCodes)
}
userRepo.AssertExpectations(t)
})
}
}

func TestAuthUseCase_Enable2FA_MoreCases(t *testing.T) {
tests := []struct {
name        string
input       Enable2FAInput
setupMock   func(*MockUserRepository, *MockRefreshTokenRepository)
expectError bool
errContains string
}{
{
name: "異常系: ユーザーIDが空",
input: Enable2FAInput{
UserID: "",
Code:   "123456",
Secret: "JBSWY3DPEHPK3PXP",
},
setupMock:   func(ur *MockUserRepository, rtr *MockRefreshTokenRepository) {},
expectError: true,
errContains: "ユーザーIDは必須です",
},
{
name: "異常系: 認証コードが空",
input: Enable2FAInput{
UserID: "valid-user-id",
Code:   "",
Secret: "JBSWY3DPEHPK3PXP",
},
setupMock:   func(ur *MockUserRepository, rtr *MockRefreshTokenRepository) {},
expectError: true,
errContains: "認証コードは必須です",
},
{
name: "異常系: 2FAが既に有効",
input: Enable2FAInput{
UserID: "valid-user-with-2fa",
Code:   "123456",
Secret: "JBSWY3DPEHPK3PXP",
},
setupMock: func(ur *MockUserRepository, rtr *MockRefreshTokenRepository) {
user := newTestUser("with2fa@example.com", "Password123!")
// 2FAを有効化状態にする
_ = user.EnableTwoFactor("JBSWY3DPEHPK3PXP", []string{"hash1", "hash2"})
ur.On("FindByID", mock.Anything, entities.UserID("valid-user-with-2fa")).Return(user, nil)
},
expectError: true,
errContains: "2段階認証は既に有効です",
},
{
name: "異常系: 無効なTOTPコード",
input: Enable2FAInput{
UserID: "valid-user-id",
Code:   "000000", // 無効なコード
Secret: "JBSWY3DPEHPK3PXP",
},
setupMock: func(ur *MockUserRepository, rtr *MockRefreshTokenRepository) {
user := newTestUser("no2fa@example.com", "Password123!")
ur.On("FindByID", mock.Anything, entities.UserID("valid-user-id")).Return(user, nil)
},
expectError: true,
errContains: "認証コードが無効です",
},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
userRepo := new(MockUserRepository)
rtRepo := new(MockRefreshTokenRepository)
tt.setupMock(userRepo, rtRepo)

uc := newTestAuthUseCase(userRepo, rtRepo)
err := uc.Enable2FA(context.Background(), tt.input)

assert.Error(t, err)
if tt.errContains != "" {
assert.Contains(t, err.Error(), tt.errContains)
}
userRepo.AssertExpectations(t)
})
}
}

func TestAuthUseCase_Verify2FA_MoreCases(t *testing.T) {
tests := []struct {
name        string
input       Verify2FAInput
setupMock   func(*MockUserRepository, *MockRefreshTokenRepository)
expectError bool
errContains string
}{
{
name: "異常系: 2FAが有効でない",
input: Verify2FAInput{
UserID: "user-no-2fa",
Code:   "123456",
},
setupMock: func(ur *MockUserRepository, rtr *MockRefreshTokenRepository) {
user := newTestUser("no2fa@example.com", "Password123!")
ur.On("FindByID", mock.Anything, entities.UserID("user-no-2fa")).Return(user, nil)
},
expectError: true,
errContains: "2段階認証が有効になっていません",
},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
userRepo := new(MockUserRepository)
rtRepo := new(MockRefreshTokenRepository)
tt.setupMock(userRepo, rtRepo)

uc := newTestAuthUseCase(userRepo, rtRepo)
output, err := uc.Verify2FA(context.Background(), tt.input)

assert.Error(t, err)
if tt.errContains != "" {
assert.Contains(t, err.Error(), tt.errContains)
}
assert.Nil(t, output)
userRepo.AssertExpectations(t)
})
}
}

func TestAuthUseCase_Disable2FA_MoreCases(t *testing.T) {
tests := []struct {
name        string
input       Disable2FAInput
setupMock   func(*MockUserRepository, *MockRefreshTokenRepository)
expectError bool
errContains string
}{
{
name: "異常系: 2FAが有効でない",
input: Disable2FAInput{
UserID:   "user-no-2fa",
Password: "Password123!",
},
setupMock: func(ur *MockUserRepository, rtr *MockRefreshTokenRepository) {
user := newTestUser("no2fa@example.com", "Password123!")
ur.On("FindByID", mock.Anything, entities.UserID("user-no-2fa")).Return(user, nil)
},
expectError: true,
errContains: "2段階認証は有効になっていません",
},
{
name: "異常系: パスワードが正しくない",
input: Disable2FAInput{
UserID:   "user-with-2fa",
Password: "WrongPassword!",
},
setupMock: func(ur *MockUserRepository, rtr *MockRefreshTokenRepository) {
user := newTestUser("with2fa@example.com", "Password123!")
_ = user.EnableTwoFactor("JBSWY3DPEHPK3PXP", []string{"hash1", "hash2"})
ur.On("FindByID", mock.Anything, entities.UserID("user-with-2fa")).Return(user, nil)
},
expectError: true,
errContains: "パスワードが正しくありません",
},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
userRepo := new(MockUserRepository)
rtRepo := new(MockRefreshTokenRepository)
tt.setupMock(userRepo, rtRepo)

uc := newTestAuthUseCase(userRepo, rtRepo)
err := uc.Disable2FA(context.Background(), tt.input)

assert.Error(t, err)
if tt.errContains != "" {
assert.Contains(t, err.Error(), tt.errContains)
}
userRepo.AssertExpectations(t)
})
}
}

func TestAuthUseCase_RegenerateBackupCodes_MoreCases(t *testing.T) {
tests := []struct {
name        string
userID      string
setupMock   func(*MockUserRepository, *MockRefreshTokenRepository)
expectError bool
errContains string
}{
{
name:   "異常系: 2FAが有効でない",
userID: "user-no-2fa",
setupMock: func(ur *MockUserRepository, rtr *MockRefreshTokenRepository) {
user := newTestUser("no2fa@example.com", "Password123!")
ur.On("FindByID", mock.Anything, entities.UserID("user-no-2fa")).Return(user, nil)
},
expectError: true,
errContains: "2段階認証が有効になっていません",
},
{
name:   "正常系: バックアップコード再生成",
userID: "user-with-2fa",
setupMock: func(ur *MockUserRepository, rtr *MockRefreshTokenRepository) {
user := newTestUser("with2fa@example.com", "Password123!")
_ = user.EnableTwoFactor("JBSWY3DPEHPK3PXP", []string{"hash1", "hash2"})
ur.On("FindByID", mock.Anything, entities.UserID("user-with-2fa")).Return(user, nil)
ur.On("Update", mock.Anything, mock.AnythingOfType("*entities.User")).Return(nil)
},
expectError: false,
},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
userRepo := new(MockUserRepository)
rtRepo := new(MockRefreshTokenRepository)
tt.setupMock(userRepo, rtRepo)

uc := newTestAuthUseCase(userRepo, rtRepo)
output, err := uc.RegenerateBackupCodes(context.Background(), tt.userID)

if tt.expectError {
assert.Error(t, err)
if tt.errContains != "" {
assert.Contains(t, err.Error(), tt.errContains)
}
assert.Nil(t, output)
} else {
require.NoError(t, err)
assert.NotNil(t, output)
assert.NotEmpty(t, output.BackupCodes)
}
userRepo.AssertExpectations(t)
})
}
}
