package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/financial-planning-calculator/backend/application/usecases"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// setTestUserID sets a user_id in the echo context for tests that require authentication
func setTestUserID(c echo.Context, userID string) {
	c.Set("user_id", userID)
}

func TestSetup2FA(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		mockSetup      func(m *MockAuthUseCase)
		expectedStatus int
	}{
		{
			name:   "Success: setup 2FA",
			userID: "user-123",
			mockSetup: func(m *MockAuthUseCase) {
				m.On("Setup2FA", mock.Anything, "user-123").Return(&usecases.Setup2FAOutput{
					Secret:      "JBSWY3DPEHPK3PXP",
					QRCodeURL:   "otpauth://totp/...",
					BackupCodes: []string{"code1", "code2"},
				}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Error: no user in context",
			userID:         "",
			mockSetup:      func(m *MockAuthUseCase) {},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:   "Error: 2FA already enabled",
			userID: "user-123",
			mockSetup: func(m *MockAuthUseCase) {
				m.On("Setup2FA", mock.Anything, "user-123").Return(nil, errors.New("2段階認証は既に有効です"))
			},
			expectedStatus: http.StatusConflict,
		},
		{
			name:   "Error: internal server error",
			userID: "user-123",
			mockSetup: func(m *MockAuthUseCase) {
				m.On("Setup2FA", mock.Anything, "user-123").Return(nil, errors.New("database error"))
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
			controller := NewTwoFactorController(mockUseCase, newTestServerConfig())

			req := httptest.NewRequest(http.MethodPost, "/auth/2fa/setup", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			if tt.userID != "" {
				setTestUserID(c, tt.userID)
			}

			err := controller.Setup2FA(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}

func TestEnable2FA(t *testing.T) {
	tests := []struct {
		name               string
		userID             string
		requestBody        interface{}
		mockSetup          func(m *MockAuthUseCase)
		expectedStatus     int
		expectHandlerError bool
	}{
		{
			name:   "Success: enable 2FA",
			userID: "user-123",
			requestBody: Enable2FARequest{
				Code:   "123456",
				Secret: "JBSWY3DPEHPK3PXP",
			},
			mockSetup: func(m *MockAuthUseCase) {
				m.On("Enable2FA", mock.Anything, mock.MatchedBy(func(input usecases.Enable2FAInput) bool {
					return input.UserID == "user-123" && input.Code == "123456"
				})).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Error: no user in context",
			userID:         "",
			requestBody:    Enable2FARequest{Code: "123456", Secret: "secret"},
			mockSetup:      func(m *MockAuthUseCase) {},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:   "Error: validation failure (missing code)",
			userID: "user-123",
			requestBody: Enable2FARequest{
				Secret: "JBSWY3DPEHPK3PXP",
			},
			mockSetup:          func(m *MockAuthUseCase) {},
			expectHandlerError: true,
		},
		{
			name:   "Error: invalid code",
			userID: "user-123",
			requestBody: Enable2FARequest{
				Code:   "000000",
				Secret: "JBSWY3DPEHPK3PXP",
			},
			mockSetup: func(m *MockAuthUseCase) {
				m.On("Enable2FA", mock.Anything, mock.Anything).Return(errors.New("認証コードが無効です"))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Error: already enabled",
			userID: "user-123",
			requestBody: Enable2FARequest{
				Code:   "123456",
				Secret: "JBSWY3DPEHPK3PXP",
			},
			mockSetup: func(m *MockAuthUseCase) {
				m.On("Enable2FA", mock.Anything, mock.Anything).Return(errors.New("2段階認証は既に有効です"))
			},
			expectedStatus: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			e.Validator = &CustomValidator{validator: validator.New()}

			mockUseCase := new(MockAuthUseCase)
			tt.mockSetup(mockUseCase)
			controller := NewTwoFactorController(mockUseCase, newTestServerConfig())

			reqJSON, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/auth/2fa/enable", bytes.NewBuffer(reqJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			if tt.userID != "" {
				setTestUserID(c, tt.userID)
			}

			err := controller.Enable2FA(c)

			if tt.expectHandlerError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
			}
		})
	}
}

func TestVerify2FA(t *testing.T) {
	tests := []struct {
		name               string
		userID             string
		requestBody        interface{}
		mockSetup          func(m *MockAuthUseCase)
		expectedStatus     int
		expectHandlerError bool
	}{
		{
			name:   "Success: verify 2FA",
			userID: "user-123",
			requestBody: Verify2FARequest{
				Code: "123456",
			},
			mockSetup: func(m *MockAuthUseCase) {
				m.On("Verify2FA", mock.Anything, mock.MatchedBy(func(input usecases.Verify2FAInput) bool {
					return input.UserID == "user-123" && input.Code == "123456"
				})).Return(&usecases.LoginOutput{
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
			name:           "Error: no user in context",
			userID:         "",
			requestBody:    Verify2FARequest{Code: "123456"},
			mockSetup:      func(m *MockAuthUseCase) {},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:   "Error: validation failure (missing code)",
			userID: "user-123",
			requestBody: Verify2FARequest{
				// Code is empty → fails required
			},
			mockSetup:          func(m *MockAuthUseCase) {},
			expectHandlerError: true,
		},
		{
			name:   "Error: invalid code",
			userID: "user-123",
			requestBody: Verify2FARequest{
				Code: "000000",
			},
			mockSetup: func(m *MockAuthUseCase) {
				m.On("Verify2FA", mock.Anything, mock.Anything).Return(nil, errors.New("認証コードが無効です"))
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			e.Validator = &CustomValidator{validator: validator.New()}

			mockUseCase := new(MockAuthUseCase)
			tt.mockSetup(mockUseCase)
			controller := NewTwoFactorController(mockUseCase, newTestServerConfig())

			reqJSON, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/auth/2fa/verify", bytes.NewBuffer(reqJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			if tt.userID != "" {
				setTestUserID(c, tt.userID)
			}

			err := controller.Verify2FA(c)

			if tt.expectHandlerError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
			}
		})
	}
}

func TestDisable2FA(t *testing.T) {
	tests := []struct {
		name               string
		userID             string
		requestBody        interface{}
		mockSetup          func(m *MockAuthUseCase)
		expectedStatus     int
		expectHandlerError bool
	}{
		{
			name:   "Success: disable 2FA",
			userID: "user-123",
			requestBody: Disable2FARequest{
				Password: "password123",
			},
			mockSetup: func(m *MockAuthUseCase) {
				m.On("Disable2FA", mock.Anything, mock.MatchedBy(func(input usecases.Disable2FAInput) bool {
					return input.UserID == "user-123"
				})).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Error: no user in context",
			userID:         "",
			requestBody:    Disable2FARequest{Password: "password123"},
			mockSetup:      func(m *MockAuthUseCase) {},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:   "Error: validation failure (missing password)",
			userID: "user-123",
			requestBody: Disable2FARequest{
				// Password is empty → fails required
			},
			mockSetup:          func(m *MockAuthUseCase) {},
			expectHandlerError: true,
		},
		{
			name:   "Error: wrong password",
			userID: "user-123",
			requestBody: Disable2FARequest{
				Password: "wrongpassword",
			},
			mockSetup: func(m *MockAuthUseCase) {
				m.On("Disable2FA", mock.Anything, mock.Anything).Return(errors.New("パスワードが正しくありません"))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Error: 2FA not enabled",
			userID: "user-123",
			requestBody: Disable2FARequest{
				Password: "password123",
			},
			mockSetup: func(m *MockAuthUseCase) {
				m.On("Disable2FA", mock.Anything, mock.Anything).Return(errors.New("2段階認証は有効になっていません"))
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			e.Validator = &CustomValidator{validator: validator.New()}

			mockUseCase := new(MockAuthUseCase)
			tt.mockSetup(mockUseCase)
			controller := NewTwoFactorController(mockUseCase, newTestServerConfig())

			reqJSON, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodDelete, "/auth/2fa", bytes.NewBuffer(reqJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			if tt.userID != "" {
				setTestUserID(c, tt.userID)
			}

			err := controller.Disable2FA(c)

			if tt.expectHandlerError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
			}
		})
	}
}

func TestRegenerateBackupCodes(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		mockSetup      func(m *MockAuthUseCase)
		expectedStatus int
	}{
		{
			name:   "Success: regenerate backup codes",
			userID: "user-123",
			mockSetup: func(m *MockAuthUseCase) {
				m.On("RegenerateBackupCodes", mock.Anything, "user-123").Return(&usecases.RegenerateBackupCodesOutput{
					BackupCodes: []string{"code1", "code2", "code3"},
				}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Error: no user in context",
			userID:         "",
			mockSetup:      func(m *MockAuthUseCase) {},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:   "Error: 2FA not enabled",
			userID: "user-123",
			mockSetup: func(m *MockAuthUseCase) {
				m.On("RegenerateBackupCodes", mock.Anything, "user-123").Return(nil, errors.New("2段階認証が有効になっていません"))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Error: internal server error",
			userID: "user-123",
			mockSetup: func(m *MockAuthUseCase) {
				m.On("RegenerateBackupCodes", mock.Anything, "user-123").Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			mockUseCase := new(MockAuthUseCase)
			tt.mockSetup(mockUseCase)
			controller := NewTwoFactorController(mockUseCase, newTestServerConfig())

			req := httptest.NewRequest(http.MethodPost, "/auth/2fa/backup-codes", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			if tt.userID != "" {
				setTestUserID(c, tt.userID)
			}

			err := controller.RegenerateBackupCodes(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}

func TestGet2FAStatus(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		mockSetup      func(m *MockAuthUseCase)
		expectedStatus int
	}{
		{
			name:   "Success: 2FA enabled",
			userID: "user-123",
			mockSetup: func(m *MockAuthUseCase) {
				m.On("Get2FAStatus", mock.Anything, "user-123").Return(&usecases.Get2FAStatusOutput{
					Enabled: true,
				}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "Success: 2FA disabled",
			userID: "user-123",
			mockSetup: func(m *MockAuthUseCase) {
				m.On("Get2FAStatus", mock.Anything, "user-123").Return(&usecases.Get2FAStatusOutput{
					Enabled: false,
				}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Error: no user in context",
			userID:         "",
			mockSetup:      func(m *MockAuthUseCase) {},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:   "Error: internal server error",
			userID: "user-123",
			mockSetup: func(m *MockAuthUseCase) {
				m.On("Get2FAStatus", mock.Anything, "user-123").Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			mockUseCase := new(MockAuthUseCase)
			tt.mockSetup(mockUseCase)
			controller := NewTwoFactorController(mockUseCase, newTestServerConfig())

			req := httptest.NewRequest(http.MethodGet, "/auth/2fa/status", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			if tt.userID != "" {
				setTestUserID(c, tt.userID)
			}

			err := controller.Get2FAStatus(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}
