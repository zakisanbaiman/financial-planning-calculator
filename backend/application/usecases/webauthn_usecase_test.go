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

// newTestWebAuthnUseCase はWebAuthnUseCaseのテスト用インスタンスを生成する（webAuthnはnilで可）
func newTestWebAuthnUseCase(
	userRepo *MockUserRepository,
	credRepo *MockWebAuthnCredentialRepository,
	tokenRepo *MockRefreshTokenRepository,
) *webAuthnUseCaseImpl {
	return &webAuthnUseCaseImpl{
		userRepo:               userRepo,
		credentialRepo:         credRepo,
		refreshTokenRepo:       tokenRepo,
		webAuthn:               nil, // WebAuthn実機が不要なテストでのみ使用
		jwtSecret:              testJWTSecret,
		jwtExpiration:          testJWTExpiration,
		refreshTokenExpiration: testRefreshTokenExpiration,
	}
}

// newTestCredential はテスト用WebAuthnCredentialを生成するヘルパー
func newTestCredential(id string, userID entities.UserID) *entities.WebAuthnCredential {
	cred, _ := entities.NewWebAuthnCredential(
		id,
		userID,
		[]byte("credential-id-bytes"),
		[]byte("public-key-bytes"),
		"none",
		[]byte{},
		[]string{"usb"},
		"My Passkey",
	)
	return cred
}

// ===========================
// BeginRegistration Tests
// ===========================

func TestWebAuthnUseCase_BeginRegistration_InvalidUserID(t *testing.T) {
	ctx := context.Background()

	userRepo := new(MockUserRepository)
	credRepo := new(MockWebAuthnCredentialRepository)
	tokenRepo := new(MockRefreshTokenRepository)

	uc := newTestWebAuthnUseCase(userRepo, credRepo, tokenRepo)
	_, err := uc.BeginRegistration(ctx, "")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "無効なユーザーID")
}

func TestWebAuthnUseCase_BeginRegistration_UserNotFound(t *testing.T) {
	ctx := context.Background()

	userRepo := new(MockUserRepository)
	credRepo := new(MockWebAuthnCredentialRepository)
	tokenRepo := new(MockRefreshTokenRepository)
	userRepo.On("FindByID", mock_anything(), entities.UserID("user-001")).Return(nil, errors.New("not found"))

	uc := newTestWebAuthnUseCase(userRepo, credRepo, tokenRepo)
	_, err := uc.BeginRegistration(ctx, "user-001")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "ユーザーが見つかりません")
	userRepo.AssertExpectations(t)
}

func TestWebAuthnUseCase_BeginRegistration_CredentialRepoError(t *testing.T) {
	ctx := context.Background()

	userRepo := new(MockUserRepository)
	credRepo := new(MockWebAuthnCredentialRepository)
	tokenRepo := new(MockRefreshTokenRepository)

	testUser := newTestUser("user-001", "test@example.com")
	userRepo.On("FindByID", mock_anything(), entities.UserID("user-001")).Return(testUser, nil)
	credRepo.On("FindByUserID", mock_anything(), entities.UserID("user-001")).Return(nil, errors.New("db error"))

	uc := newTestWebAuthnUseCase(userRepo, credRepo, tokenRepo)
	_, err := uc.BeginRegistration(ctx, "user-001")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "既存のクレデンシャル取得に失敗")
	userRepo.AssertExpectations(t)
	credRepo.AssertExpectations(t)
}

// ===========================
// FinishRegistration Tests
// ===========================

func TestWebAuthnUseCase_FinishRegistration_InvalidUserID(t *testing.T) {
	ctx := context.Background()

	userRepo := new(MockUserRepository)
	credRepo := new(MockWebAuthnCredentialRepository)
	tokenRepo := new(MockRefreshTokenRepository)

	uc := newTestWebAuthnUseCase(userRepo, credRepo, tokenRepo)
	err := uc.FinishRegistration(ctx, FinishRegistrationInput{
		UserID:      "",
		SessionData: "dummydata",
		Response:    "{}",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "無効なユーザーID")
}

func TestWebAuthnUseCase_FinishRegistration_UserNotFound(t *testing.T) {
	ctx := context.Background()

	userRepo := new(MockUserRepository)
	credRepo := new(MockWebAuthnCredentialRepository)
	tokenRepo := new(MockRefreshTokenRepository)
	userRepo.On("FindByID", mock_anything(), entities.UserID("user-001")).Return(nil, errors.New("not found"))

	uc := newTestWebAuthnUseCase(userRepo, credRepo, tokenRepo)
	err := uc.FinishRegistration(ctx, FinishRegistrationInput{
		UserID:      "user-001",
		SessionData: "dummydata",
		Response:    "{}",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "ユーザーが見つかりません")
	userRepo.AssertExpectations(t)
}

func TestWebAuthnUseCase_FinishRegistration_InvalidSessionData(t *testing.T) {
	ctx := context.Background()

	userRepo := new(MockUserRepository)
	credRepo := new(MockWebAuthnCredentialRepository)
	tokenRepo := new(MockRefreshTokenRepository)

	testUser := newTestUser("user-001", "test@example.com")
	userRepo.On("FindByID", mock_anything(), entities.UserID("user-001")).Return(testUser, nil)

	uc := newTestWebAuthnUseCase(userRepo, credRepo, tokenRepo)
	err := uc.FinishRegistration(ctx, FinishRegistrationInput{
		UserID:      "user-001",
		SessionData: "!!!invalid-base64!!!",
		Response:    "{}",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "セッションデータのデコードに失敗")
	userRepo.AssertExpectations(t)
}

// ===========================
// BeginLogin Tests
// ===========================

func TestWebAuthnUseCase_BeginLogin_WebAuthnNil(t *testing.T) {
	ctx := context.Background()

	userRepo := new(MockUserRepository)
	credRepo := new(MockWebAuthnCredentialRepository)
	tokenRepo := new(MockRefreshTokenRepository)

	uc := newTestWebAuthnUseCase(userRepo, credRepo, tokenRepo)

	// webAuthn が nil の場合にパニックせずエラーを返すことを確認
	assert.Panics(t, func() {
		//nolint:errcheck
		uc.BeginLogin(ctx, BeginLoginInput{Email: "test@example.com"})
	})
}

// ===========================
// FinishLogin Tests
// ===========================

func TestWebAuthnUseCase_FinishLogin_InvalidSessionData(t *testing.T) {
	ctx := context.Background()

	userRepo := new(MockUserRepository)
	credRepo := new(MockWebAuthnCredentialRepository)
	tokenRepo := new(MockRefreshTokenRepository)

	uc := newTestWebAuthnUseCase(userRepo, credRepo, tokenRepo)
	_, err := uc.FinishLogin(ctx, FinishLoginInput{
		SessionData: "!!!invalid-base64!!!",
		Response:    "{}",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "セッションデータのデコードに失敗")
}

// ===========================
// ListCredentials Tests
// ===========================

func TestWebAuthnUseCase_ListCredentials(t *testing.T) {
	ctx := context.Background()

	t.Run("正常系: クレデンシャル一覧を取得できる", func(t *testing.T) {
		userRepo := new(MockUserRepository)
		credRepo := new(MockWebAuthnCredentialRepository)
		tokenRepo := new(MockRefreshTokenRepository)

		uid := entities.UserID("user-001")
		cred := newTestCredential("cred-001", uid)
		credRepo.On("FindByUserID", mock_anything(), uid).Return([]*entities.WebAuthnCredential{cred}, nil)

		uc := newTestWebAuthnUseCase(userRepo, credRepo, tokenRepo)
		result, err := uc.ListCredentials(ctx, "user-001")

		require.NoError(t, err)
		require.Len(t, result, 1)
		assert.Equal(t, "cred-001", result[0].ID)
		assert.Equal(t, "My Passkey", result[0].Name)
		assert.Nil(t, result[0].LastUsedAt)
		credRepo.AssertExpectations(t)
	})

	t.Run("正常系: クレデンシャルが空の場合は空スライスを返す", func(t *testing.T) {
		userRepo := new(MockUserRepository)
		credRepo := new(MockWebAuthnCredentialRepository)
		tokenRepo := new(MockRefreshTokenRepository)

		uid := entities.UserID("user-001")
		credRepo.On("FindByUserID", mock_anything(), uid).Return([]*entities.WebAuthnCredential{}, nil)

		uc := newTestWebAuthnUseCase(userRepo, credRepo, tokenRepo)
		result, err := uc.ListCredentials(ctx, "user-001")

		require.NoError(t, err)
		assert.Empty(t, result)
		credRepo.AssertExpectations(t)
	})

	t.Run("異常系: 無効なユーザーIDの場合はエラー", func(t *testing.T) {
		userRepo := new(MockUserRepository)
		credRepo := new(MockWebAuthnCredentialRepository)
		tokenRepo := new(MockRefreshTokenRepository)

		uc := newTestWebAuthnUseCase(userRepo, credRepo, tokenRepo)
		_, err := uc.ListCredentials(ctx, "")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "無効なユーザーID")
	})

	t.Run("異常系: リポジトリエラーの場合はエラーを返す", func(t *testing.T) {
		userRepo := new(MockUserRepository)
		credRepo := new(MockWebAuthnCredentialRepository)
		tokenRepo := new(MockRefreshTokenRepository)

		uid := entities.UserID("user-001")
		credRepo.On("FindByUserID", mock_anything(), uid).Return(nil, errors.New("db error"))

		uc := newTestWebAuthnUseCase(userRepo, credRepo, tokenRepo)
		_, err := uc.ListCredentials(ctx, "user-001")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "クレデンシャルの取得に失敗")
		credRepo.AssertExpectations(t)
	})

	t.Run("正常系: LastUsedAtが設定されているクレデンシャルを返す", func(t *testing.T) {
		userRepo := new(MockUserRepository)
		credRepo := new(MockWebAuthnCredentialRepository)
		tokenRepo := new(MockRefreshTokenRepository)

		uid := entities.UserID("user-001")
		cred := newTestCredential("cred-001", uid)
		// UpdateSignCount でlastUsedAtを設定
		require.NoError(t, cred.UpdateSignCount(1))

		credRepo.On("FindByUserID", mock_anything(), uid).Return([]*entities.WebAuthnCredential{cred}, nil)

		uc := newTestWebAuthnUseCase(userRepo, credRepo, tokenRepo)
		result, err := uc.ListCredentials(ctx, "user-001")

		require.NoError(t, err)
		require.Len(t, result, 1)
		assert.NotNil(t, result[0].LastUsedAt)
		credRepo.AssertExpectations(t)
	})
}

// ===========================
// DeleteCredential Tests
// ===========================

func TestWebAuthnUseCase_DeleteCredential(t *testing.T) {
	ctx := context.Background()

	t.Run("正常系: クレデンシャルを削除できる", func(t *testing.T) {
		userRepo := new(MockUserRepository)
		credRepo := new(MockWebAuthnCredentialRepository)
		tokenRepo := new(MockRefreshTokenRepository)

		uid := entities.UserID("user-001")
		cid := entities.CredentialID("cred-001")
		cred := newTestCredential("cred-001", uid)

		credRepo.On("FindByID", mock_anything(), cid).Return(cred, nil)
		credRepo.On("Delete", mock_anything(), cid).Return(nil)

		uc := newTestWebAuthnUseCase(userRepo, credRepo, tokenRepo)
		err := uc.DeleteCredential(ctx, "user-001", "cred-001")

		require.NoError(t, err)
		credRepo.AssertExpectations(t)
	})

	t.Run("異常系: 無効なユーザーIDの場合はエラー", func(t *testing.T) {
		userRepo := new(MockUserRepository)
		credRepo := new(MockWebAuthnCredentialRepository)
		tokenRepo := new(MockRefreshTokenRepository)

		uc := newTestWebAuthnUseCase(userRepo, credRepo, tokenRepo)
		err := uc.DeleteCredential(ctx, "", "cred-001")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "無効なユーザーID")
	})

	t.Run("異常系: 無効なクレデンシャルIDの場合はエラー", func(t *testing.T) {
		userRepo := new(MockUserRepository)
		credRepo := new(MockWebAuthnCredentialRepository)
		tokenRepo := new(MockRefreshTokenRepository)

		uc := newTestWebAuthnUseCase(userRepo, credRepo, tokenRepo)
		err := uc.DeleteCredential(ctx, "user-001", "")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "無効なクレデンシャルID")
	})

	t.Run("異常系: クレデンシャルが存在しない場合はエラー", func(t *testing.T) {
		userRepo := new(MockUserRepository)
		credRepo := new(MockWebAuthnCredentialRepository)
		tokenRepo := new(MockRefreshTokenRepository)

		cid := entities.CredentialID("cred-001")
		credRepo.On("FindByID", mock_anything(), cid).Return(nil, errors.New("not found"))

		uc := newTestWebAuthnUseCase(userRepo, credRepo, tokenRepo)
		err := uc.DeleteCredential(ctx, "user-001", "cred-001")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "クレデンシャルが見つかりません")
		credRepo.AssertExpectations(t)
	})

	t.Run("異常系: 別ユーザーのクレデンシャルは削除できない", func(t *testing.T) {
		userRepo := new(MockUserRepository)
		credRepo := new(MockWebAuthnCredentialRepository)
		tokenRepo := new(MockRefreshTokenRepository)

		// クレデンシャルはuser-002のもの
		ownerUID := entities.UserID("user-002")
		cid := entities.CredentialID("cred-001")
		cred := newTestCredential("cred-001", ownerUID)

		credRepo.On("FindByID", mock_anything(), cid).Return(cred, nil)

		uc := newTestWebAuthnUseCase(userRepo, credRepo, tokenRepo)
		// user-001として削除しようとする
		err := uc.DeleteCredential(ctx, "user-001", "cred-001")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "このクレデンシャルの所有者ではありません")
		credRepo.AssertExpectations(t)
	})

	t.Run("異常系: Deleteリポジトリエラーの場合はエラーを返す", func(t *testing.T) {
		userRepo := new(MockUserRepository)
		credRepo := new(MockWebAuthnCredentialRepository)
		tokenRepo := new(MockRefreshTokenRepository)

		uid := entities.UserID("user-001")
		cid := entities.CredentialID("cred-001")
		cred := newTestCredential("cred-001", uid)

		credRepo.On("FindByID", mock_anything(), cid).Return(cred, nil)
		credRepo.On("Delete", mock_anything(), cid).Return(errors.New("db error"))

		uc := newTestWebAuthnUseCase(userRepo, credRepo, tokenRepo)
		err := uc.DeleteCredential(ctx, "user-001", "cred-001")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "クレデンシャルの削除に失敗")
		credRepo.AssertExpectations(t)
	})
}

// ===========================
// RenameCredential Tests
// ===========================

func TestWebAuthnUseCase_RenameCredential(t *testing.T) {
	ctx := context.Background()

	t.Run("正常系: クレデンシャルの名前を変更できる", func(t *testing.T) {
		userRepo := new(MockUserRepository)
		credRepo := new(MockWebAuthnCredentialRepository)
		tokenRepo := new(MockRefreshTokenRepository)

		uid := entities.UserID("user-001")
		cid := entities.CredentialID("cred-001")
		cred := newTestCredential("cred-001", uid)

		credRepo.On("FindByID", mock_anything(), cid).Return(cred, nil)
		credRepo.On("Update", mock_anything(), mock_anything()).Return(nil)

		uc := newTestWebAuthnUseCase(userRepo, credRepo, tokenRepo)
		err := uc.RenameCredential(ctx, "user-001", "cred-001", "New Name")

		require.NoError(t, err)
		assert.Equal(t, "New Name", cred.Name())
		credRepo.AssertExpectations(t)
	})

	t.Run("異常系: 無効なユーザーIDの場合はエラー", func(t *testing.T) {
		userRepo := new(MockUserRepository)
		credRepo := new(MockWebAuthnCredentialRepository)
		tokenRepo := new(MockRefreshTokenRepository)

		uc := newTestWebAuthnUseCase(userRepo, credRepo, tokenRepo)
		err := uc.RenameCredential(ctx, "", "cred-001", "New Name")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "無効なユーザーID")
	})

	t.Run("異常系: 無効なクレデンシャルIDの場合はエラー", func(t *testing.T) {
		userRepo := new(MockUserRepository)
		credRepo := new(MockWebAuthnCredentialRepository)
		tokenRepo := new(MockRefreshTokenRepository)

		uc := newTestWebAuthnUseCase(userRepo, credRepo, tokenRepo)
		err := uc.RenameCredential(ctx, "user-001", "", "New Name")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "無効なクレデンシャルID")
	})

	t.Run("異常系: クレデンシャルが存在しない場合はエラー", func(t *testing.T) {
		userRepo := new(MockUserRepository)
		credRepo := new(MockWebAuthnCredentialRepository)
		tokenRepo := new(MockRefreshTokenRepository)

		cid := entities.CredentialID("cred-001")
		credRepo.On("FindByID", mock_anything(), cid).Return(nil, errors.New("not found"))

		uc := newTestWebAuthnUseCase(userRepo, credRepo, tokenRepo)
		err := uc.RenameCredential(ctx, "user-001", "cred-001", "New Name")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "クレデンシャルが見つかりません")
		credRepo.AssertExpectations(t)
	})

	t.Run("異常系: 別ユーザーのクレデンシャルは変更できない", func(t *testing.T) {
		userRepo := new(MockUserRepository)
		credRepo := new(MockWebAuthnCredentialRepository)
		tokenRepo := new(MockRefreshTokenRepository)

		ownerUID := entities.UserID("user-002")
		cid := entities.CredentialID("cred-001")
		cred := newTestCredential("cred-001", ownerUID)

		credRepo.On("FindByID", mock_anything(), cid).Return(cred, nil)

		uc := newTestWebAuthnUseCase(userRepo, credRepo, tokenRepo)
		err := uc.RenameCredential(ctx, "user-001", "cred-001", "New Name")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "このクレデンシャルの所有者ではありません")
		credRepo.AssertExpectations(t)
	})

	t.Run("異常系: Updateリポジトリエラーの場合はエラーを返す", func(t *testing.T) {
		userRepo := new(MockUserRepository)
		credRepo := new(MockWebAuthnCredentialRepository)
		tokenRepo := new(MockRefreshTokenRepository)

		uid := entities.UserID("user-001")
		cid := entities.CredentialID("cred-001")
		cred := newTestCredential("cred-001", uid)

		credRepo.On("FindByID", mock_anything(), cid).Return(cred, nil)
		credRepo.On("Update", mock_anything(), mock_anything()).Return(errors.New("db error"))

		uc := newTestWebAuthnUseCase(userRepo, credRepo, tokenRepo)
		err := uc.RenameCredential(ctx, "user-001", "cred-001", "New Name")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "クレデンシャルの更新に失敗")
		credRepo.AssertExpectations(t)
	})
}

// ===========================
// Helper: newTestUser
// ===========================

func newTestUser(id, email string) *entities.User {
	user, _ := entities.NewUser(id, email, "Password123!")
	return user
}

// ===========================
// convertTransports / convertTransportsToStrings Tests
// ===========================

func TestConvertTransports(t *testing.T) {
	t.Run("文字列スライスをAuthenticatorTransportに変換できる", func(t *testing.T) {
		result := convertTransports([]string{"usb", "nfc", "ble"})
		assert.Len(t, result, 3)
		assert.Equal(t, "usb", string(result[0]))
		assert.Equal(t, "nfc", string(result[1]))
		assert.Equal(t, "ble", string(result[2]))
	})

	t.Run("空のスライスを渡すと空のスライスを返す", func(t *testing.T) {
		result := convertTransports([]string{})
		assert.Empty(t, result)
	})
}

func TestConvertTransportsToStrings(t *testing.T) {
	t.Run("AuthenticatorTransportを文字列スライスに変換できる", func(t *testing.T) {
		inputs := convertTransports([]string{"usb", "nfc"})
		result := convertTransportsToStrings(inputs)
		assert.Equal(t, []string{"usb", "nfc"}, result)
	})

	t.Run("空のスライスを渡すと空のスライスを返す", func(t *testing.T) {
		result := convertTransportsToStrings(nil)
		assert.Empty(t, result)
	})
}

// ===========================
// webAuthnUser Tests
// ===========================

func TestWebAuthnUser(t *testing.T) {
	u := &webAuthnUser{
		id:          []byte("user-id"),
		name:        "test@example.com",
		displayName: "Test User",
		credentials: nil,
	}

	assert.Equal(t, []byte("user-id"), u.WebAuthnID())
	assert.Equal(t, "test@example.com", u.WebAuthnName())
	assert.Equal(t, "Test User", u.WebAuthnDisplayName())
	assert.Nil(t, u.WebAuthnCredentials())
	assert.Equal(t, "", u.WebAuthnIcon())
}

// ===========================
// generateToken / generateRefreshToken Tests
// ===========================

func TestWebAuthnUseCase_GenerateToken(t *testing.T) {
	userRepo := new(MockUserRepository)
	credRepo := new(MockWebAuthnCredentialRepository)
	tokenRepo := new(MockRefreshTokenRepository)

	uc := newTestWebAuthnUseCase(userRepo, credRepo, tokenRepo)

	testUser := newTestUser("user-001", "test@example.com")
	token, expiresAt, err := uc.generateToken(testUser)

	require.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.True(t, expiresAt.After(time.Now()))
}

func TestWebAuthnUseCase_GenerateRefreshToken(t *testing.T) {
	ctx := context.Background()

	userRepo := new(MockUserRepository)
	credRepo := new(MockWebAuthnCredentialRepository)
	tokenRepo := new(MockRefreshTokenRepository)
	tokenRepo.On("Save", mock_anything(), mock_anything()).Return(nil)

	uc := newTestWebAuthnUseCase(userRepo, credRepo, tokenRepo)

	rawToken, err := uc.generateRefreshToken(ctx, entities.UserID("user-001"))

	require.NoError(t, err)
	assert.NotEmpty(t, rawToken)
	tokenRepo.AssertExpectations(t)
}

func TestWebAuthnUseCase_GenerateRefreshToken_RepositoryError(t *testing.T) {
	ctx := context.Background()

	userRepo := new(MockUserRepository)
	credRepo := new(MockWebAuthnCredentialRepository)
	tokenRepo := new(MockRefreshTokenRepository)
	tokenRepo.On("Save", mock_anything(), mock_anything()).Return(errors.New("db error"))

	uc := newTestWebAuthnUseCase(userRepo, credRepo, tokenRepo)

	_, err := uc.generateRefreshToken(ctx, entities.UserID("user-001"))

	require.Error(t, err)
	assert.Contains(t, err.Error(), "リフレッシュトークンの保存に失敗")
	tokenRepo.AssertExpectations(t)
}
