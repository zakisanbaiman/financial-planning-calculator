package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/financial-planning-calculator/backend/application/usecases"
	"github.com/financial-planning-calculator/backend/config"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuthUseCase is a mock implementation of AuthUseCase
type MockAuthUseCase struct {
	mock.Mock
}

func (m *MockAuthUseCase) Register(ctx context.Context, input usecases.RegisterInput) (*usecases.RegisterOutput, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecases.RegisterOutput), args.Error(1)
}

func (m *MockAuthUseCase) Login(ctx context.Context, input usecases.LoginInput) (*usecases.LoginOutput, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecases.LoginOutput), args.Error(1)
}

func (m *MockAuthUseCase) VerifyToken(ctx context.Context, tokenString string) (*usecases.TokenClaims, error) {
	args := m.Called(ctx, tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecases.TokenClaims), args.Error(1)
}

func (m *MockAuthUseCase) RefreshAccessToken(ctx context.Context, refreshToken string) (*usecases.RefreshOutput, error) {
	args := m.Called(ctx, refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecases.RefreshOutput), args.Error(1)
}

func (m *MockAuthUseCase) RevokeRefreshToken(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockAuthUseCase) GitHubOAuthLogin(ctx context.Context, input usecases.GitHubOAuthInput) (*usecases.LoginOutput, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecases.LoginOutput), args.Error(1)
}

func (m *MockAuthUseCase) Setup2FA(ctx context.Context, userID string) (*usecases.Setup2FAOutput, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecases.Setup2FAOutput), args.Error(1)
}

func (m *MockAuthUseCase) Enable2FA(ctx context.Context, input usecases.Enable2FAInput) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}

func (m *MockAuthUseCase) Verify2FA(ctx context.Context, input usecases.Verify2FAInput) (*usecases.LoginOutput, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecases.LoginOutput), args.Error(1)
}

func (m *MockAuthUseCase) Disable2FA(ctx context.Context, input usecases.Disable2FAInput) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}

func (m *MockAuthUseCase) RegenerateBackupCodes(ctx context.Context, userID string) (*usecases.RegenerateBackupCodesOutput, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecases.RegenerateBackupCodesOutput), args.Error(1)
}

func (m *MockAuthUseCase) Get2FAStatus(ctx context.Context, userID string) (*usecases.Get2FAStatusOutput, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecases.Get2FAStatusOutput), args.Error(1)
}

func (m *MockAuthUseCase) ForgotPassword(ctx context.Context, input usecases.ForgotPasswordInput) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}

func (m *MockAuthUseCase) ResetPassword(ctx context.Context, input usecases.ResetPasswordInput) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}

// newTestServerConfig creates a minimal ServerConfig for tests
func newTestServerConfig() *config.ServerConfig {
	return &config.ServerConfig{
		CookieSecure:           false,
		JWTExpiration:          time.Hour,
		RefreshTokenExpiration: 7 * 24 * time.Hour,
	}
}

func TestRegister(t *testing.T) {
	tests := []struct {
		name               string
		requestBody        interface{}
		mockSetup          func(m *MockAuthUseCase)
		expectedStatus     int
		expectHandlerError bool
	}{
		{
			name: "Success: valid registration",
			requestBody: RegisterRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(m *MockAuthUseCase) {
				m.On("Register", mock.Anything, mock.MatchedBy(func(input usecases.RegisterInput) bool {
					return input.Email == "test@example.com"
				})).Return(&usecases.RegisterOutput{
					UserID:       "user-123",
					Email:        "test@example.com",
					Token:        "access-token",
					RefreshToken: "refresh-token",
					ExpiresAt:    "2030-01-01T00:00:00Z",
				}, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Error: email already registered",
			requestBody: RegisterRequest{
				Email:    "existing@example.com",
				Password: "password123",
			},
			mockSetup: func(m *MockAuthUseCase) {
				m.On("Register", mock.Anything, mock.Anything).Return(nil, errors.New("このメールアドレスは既に登録されています"))
			},
			expectedStatus: http.StatusConflict,
		},
		{
			name: "Error: invalid email format",
			requestBody: RegisterRequest{
				Email:    "not-an-email",
				Password: "password123",
			},
			mockSetup:          func(m *MockAuthUseCase) {},
			expectHandlerError: true,
		},
		{
			name: "Error: password too short",
			requestBody: RegisterRequest{
				Email:    "test@example.com",
				Password: "short",
			},
			mockSetup:          func(m *MockAuthUseCase) {},
			expectHandlerError: true,
		},
		{
			name: "Error: internal server error",
			requestBody: RegisterRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(m *MockAuthUseCase) {
				m.On("Register", mock.Anything, mock.Anything).Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			e.Validator = &CustomValidator{validator: validator.New()}

			mockUseCase := new(MockAuthUseCase)
			tt.mockSetup(mockUseCase)
			controller := NewAuthController(mockUseCase, newTestServerConfig())

			reqJSON, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(reqJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := controller.Register(c)

			if tt.expectHandlerError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	tests := []struct {
		name               string
		requestBody        interface{}
		mockSetup          func(m *MockAuthUseCase)
		expectedStatus     int
		expectHandlerError bool
	}{
		{
			name: "Success: login without 2FA",
			requestBody: LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(m *MockAuthUseCase) {
				m.On("Login", mock.Anything, mock.Anything).Return(&usecases.LoginOutput{
					UserID:       "user-123",
					Email:        "test@example.com",
					Token:        "access-token",
					RefreshToken: "refresh-token",
					ExpiresAt:    "2030-01-01T00:00:00Z",
				}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Success: login with 2FA required (empty RefreshToken)",
			requestBody: LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(m *MockAuthUseCase) {
				m.On("Login", mock.Anything, mock.Anything).Return(&usecases.LoginOutput{
					UserID:       "user-123",
					Email:        "test@example.com",
					Token:        "temp-token",
					RefreshToken: "", // empty = 2FA required
					ExpiresAt:    "2030-01-01T00:00:00Z",
				}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Error: invalid credentials",
			requestBody: LoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			mockSetup: func(m *MockAuthUseCase) {
				m.On("Login", mock.Anything, mock.Anything).Return(nil, errors.New("メールアドレスまたはパスワードが正しくありません"))
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Error: validation failure (missing email)",
			requestBody: LoginRequest{
				Password: "password123",
			},
			mockSetup:          func(m *MockAuthUseCase) {},
			expectHandlerError: true,
		},
		{
			name: "Error: internal server error",
			requestBody: LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(m *MockAuthUseCase) {
				m.On("Login", mock.Anything, mock.Anything).Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			e.Validator = &CustomValidator{validator: validator.New()}

			mockUseCase := new(MockAuthUseCase)
			tt.mockSetup(mockUseCase)
			controller := NewAuthController(mockUseCase, newTestServerConfig())

			reqJSON, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(reqJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := controller.Login(c)

			if tt.expectHandlerError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
			}
		})
	}
}

func TestRefresh(t *testing.T) {
	tests := []struct {
		name           string
		cookieToken    string
		requestBody    interface{}
		mockSetup      func(m *MockAuthUseCase)
		expectedStatus int
		expectError    bool
	}{
		{
			name:        "Success: refresh from cookie",
			cookieToken: "valid-refresh-token",
			mockSetup: func(m *MockAuthUseCase) {
				m.On("RefreshAccessToken", mock.Anything, "valid-refresh-token").Return(&usecases.RefreshOutput{
					Token:     "new-access-token",
					ExpiresAt: "2030-01-01T00:00:00Z",
				}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Success: refresh from request body",
			requestBody: RefreshRequest{
				RefreshToken: "valid-refresh-token",
			},
			mockSetup: func(m *MockAuthUseCase) {
				m.On("RefreshAccessToken", mock.Anything, "valid-refresh-token").Return(&usecases.RefreshOutput{
					Token:     "new-access-token",
					ExpiresAt: "2030-01-01T00:00:00Z",
				}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Error: invalid refresh token",
			requestBody: RefreshRequest{
				RefreshToken: "invalid-token",
			},
			mockSetup: func(m *MockAuthUseCase) {
				m.On("RefreshAccessToken", mock.Anything, "invalid-token").Return(nil, errors.New("無効なリフレッシュトークンです"))
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Error: expired refresh token",
			requestBody: RefreshRequest{
				RefreshToken: "expired-token",
			},
			mockSetup: func(m *MockAuthUseCase) {
				m.On("RefreshAccessToken", mock.Anything, "expired-token").Return(nil, errors.New("リフレッシュトークンの有効期限が切れているか、失効されています"))
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Error: no refresh token",
			mockSetup:      func(m *MockAuthUseCase) {},
			expectError:    true,
			requestBody:    RefreshRequest{}, // empty refresh token → validate fail
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			e.Validator = &CustomValidator{validator: validator.New()}

			mockUseCase := new(MockAuthUseCase)
			tt.mockSetup(mockUseCase)
			controller := NewAuthController(mockUseCase, newTestServerConfig())

			var reqBody []byte
			if tt.requestBody != nil {
				reqBody, _ = json.Marshal(tt.requestBody)
			}
			req := httptest.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewBuffer(reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			if tt.cookieToken != "" {
				req.AddCookie(&http.Cookie{Name: "refresh_token", Value: tt.cookieToken})
			}
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := controller.Refresh(c)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
			}
		})
	}
}

func TestLogout(t *testing.T) {
	e := echo.New()
	mockUseCase := new(MockAuthUseCase)
	controller := NewAuthController(mockUseCase, newTestServerConfig())

	req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := controller.Logout(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}
