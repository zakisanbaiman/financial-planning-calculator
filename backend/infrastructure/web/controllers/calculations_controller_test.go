package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/financial-planning-calculator/backend/application/usecases"
	"github.com/financial-planning-calculator/backend/domain/entities"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCalculateProjectionUseCase is a mock implementation of CalculateProjectionUseCase
type MockCalculateProjectionUseCase struct {
	mock.Mock
}

func (m *MockCalculateProjectionUseCase) CalculateAssetProjection(ctx context.Context, input usecases.AssetProjectionInput) (*usecases.AssetProjectionOutput, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecases.AssetProjectionOutput), args.Error(1)
}

func (m *MockCalculateProjectionUseCase) CalculateRetirementProjection(ctx context.Context, input usecases.RetirementProjectionInput) (*usecases.RetirementProjectionOutput, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecases.RetirementProjectionOutput), args.Error(1)
}

func (m *MockCalculateProjectionUseCase) CalculateEmergencyFundProjection(ctx context.Context, input usecases.EmergencyFundProjectionInput) (*usecases.EmergencyFundProjectionOutput, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecases.EmergencyFundProjectionOutput), args.Error(1)
}

func (m *MockCalculateProjectionUseCase) CalculateComprehensiveProjection(ctx context.Context, input usecases.ComprehensiveProjectionInput) (*usecases.ComprehensiveProjectionOutput, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecases.ComprehensiveProjectionOutput), args.Error(1)
}

func (m *MockCalculateProjectionUseCase) CalculateGoalProjection(ctx context.Context, input usecases.GoalProjectionInput) (*usecases.GoalProjectionOutput, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecases.GoalProjectionOutput), args.Error(1)
}

// CustomValidator wraps the go-playground validator
type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return err
	}
	return nil
}

func TestAssetProjectionValidation(t *testing.T) {
	tests := []struct {
		name           string
		years          int
		expectError    bool
		expectedStatus int
	}{
		{
			name:           "Valid: 1 year",
			years:          1,
			expectError:    false,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Valid: 50 years (previously max)",
			years:          50,
			expectError:    false,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Valid: 100 years (new max)",
			years:          100,
			expectError:    false,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid: 0 years",
			years:          0,
			expectError:    true,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid: 101 years",
			years:          101,
			expectError:    true,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid: negative years",
			years:          -1,
			expectError:    true,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			e := echo.New()
			e.Validator = &CustomValidator{validator: validator.New()}
			
			mockUseCase := new(MockCalculateProjectionUseCase)
			controller := NewCalculationsController(mockUseCase)

			// Create request
			reqBody := AssetProjectionRequest{
				UserID: "test-user",
				Years:  tt.years,
			}
			reqJSON, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(http.MethodPost, "/calculations/asset-projection", bytes.NewBuffer(reqJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Mock the use case only for valid cases
			if !tt.expectError {
				mockUseCase.On("CalculateAssetProjection", mock.Anything, mock.MatchedBy(func(input usecases.AssetProjectionInput) bool {
					return input.UserID == entities.UserID("test-user") && input.Years == tt.years
				})).Return(&usecases.AssetProjectionOutput{
					Projections: []entities.AssetProjection{},
					Summary:     usecases.ProjectionSummary{},
				}, nil)
			}

			// Execute
			err := controller.CalculateAssetProjection(c)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
			}
		})
	}
}

func TestComprehensiveProjectionValidation(t *testing.T) {
	tests := []struct {
		name           string
		years          int
		expectError    bool
		expectedStatus int
	}{
		{
			name:           "Valid: 100 years",
			years:          100,
			expectError:    false,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid: 101 years",
			years:          101,
			expectError:    true,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			e := echo.New()
			e.Validator = &CustomValidator{validator: validator.New()}
			
			mockUseCase := new(MockCalculateProjectionUseCase)
			controller := NewCalculationsController(mockUseCase)

			// Create request
			reqBody := ComprehensiveProjectionRequest{
				UserID: "test-user",
				Years:  tt.years,
			}
			reqJSON, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(http.MethodPost, "/calculations/comprehensive", bytes.NewBuffer(reqJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Mock the use case only for valid cases
			if !tt.expectError {
				mockUseCase.On("CalculateComprehensiveProjection", mock.Anything, mock.MatchedBy(func(input usecases.ComprehensiveProjectionInput) bool {
					return input.UserID == entities.UserID("test-user") && input.Years == tt.years
				})).Return(&usecases.ComprehensiveProjectionOutput{}, nil)
			}

			// Execute
			err := controller.CalculateComprehensiveProjection(c)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
			}
		})
	}
}
