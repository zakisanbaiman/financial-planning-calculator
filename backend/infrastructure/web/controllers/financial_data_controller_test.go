package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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

// MockManageFinancialDataUseCase is a mock implementation of ManageFinancialDataUseCase
type MockManageFinancialDataUseCase struct {
	mock.Mock
}

func (m *MockManageFinancialDataUseCase) CreateFinancialPlan(ctx context.Context, input usecases.CreateFinancialPlanInput) (*usecases.CreateFinancialPlanOutput, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecases.CreateFinancialPlanOutput), args.Error(1)
}

func (m *MockManageFinancialDataUseCase) GetFinancialPlan(ctx context.Context, input usecases.GetFinancialPlanInput) (*usecases.GetFinancialPlanOutput, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecases.GetFinancialPlanOutput), args.Error(1)
}

func (m *MockManageFinancialDataUseCase) UpdateFinancialProfile(ctx context.Context, input usecases.UpdateFinancialProfileInput) (*usecases.UpdateFinancialProfileOutput, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecases.UpdateFinancialProfileOutput), args.Error(1)
}

func (m *MockManageFinancialDataUseCase) UpdateRetirementData(ctx context.Context, input usecases.UpdateRetirementDataInput) (*usecases.UpdateRetirementDataOutput, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecases.UpdateRetirementDataOutput), args.Error(1)
}

func (m *MockManageFinancialDataUseCase) UpdateEmergencyFund(ctx context.Context, input usecases.UpdateEmergencyFundInput) (*usecases.UpdateEmergencyFundOutput, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecases.UpdateEmergencyFundOutput), args.Error(1)
}

func (m *MockManageFinancialDataUseCase) DeleteFinancialPlan(ctx context.Context, input usecases.DeleteFinancialPlanInput) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}

func newFinancialDataEcho() *echo.Echo {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	return e
}

// validFinancialDataRequest returns a valid CreateFinancialDataRequest for testing
func validFinancialDataRequest() CreateFinancialDataRequest {
	return CreateFinancialDataRequest{
		UserID:           "user-123",
		MonthlyIncome:    400000,
		InvestmentReturn: 5.0,
		InflationRate:    2.0,
		MonthlyExpenses: []ExpenseItemRequest{
			{Category: "生活費", Amount: 200000},
		},
		CurrentSavings: []SavingsItemRequest{
			{Type: "deposit", Amount: 500000},
		},
	}
}

func TestCreateFinancialData(t *testing.T) {
	tests := []struct {
		name               string
		requestBody        interface{}
		mockSetup          func(m *MockManageFinancialDataUseCase)
		expectedStatus     int
		expectHandlerError bool
	}{
		{
			name:        "Success: create financial data",
			requestBody: validFinancialDataRequest(),
			mockSetup: func(m *MockManageFinancialDataUseCase) {
				m.On("CreateFinancialPlan", mock.Anything, mock.MatchedBy(func(input usecases.CreateFinancialPlanInput) bool {
					return input.UserID == entities.UserID("user-123")
				})).Return(&usecases.CreateFinancialPlanOutput{
					UserID:    entities.UserID("user-123"),
					CreatedAt: "2030-01-01T00:00:00Z",
				}, nil)
				m.On("GetFinancialPlan", mock.Anything, mock.Anything).Return(&usecases.GetFinancialPlanOutput{
					Plan: nil,
				}, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Error: missing user_id",
			requestBody: CreateFinancialDataRequest{
				MonthlyIncome:    400000,
				InvestmentReturn: 5.0,
				InflationRate:    2.0,
			},
			mockSetup:          func(m *MockManageFinancialDataUseCase) {},
			expectHandlerError: true,
		},
		{
			name: "Error: monthly expenses exceed income (business logic)",
			requestBody: CreateFinancialDataRequest{
				UserID:           "user-123",
				MonthlyIncome:    100000,
				InvestmentReturn: 5.0,
				InflationRate:    2.0,
				MonthlyExpenses: []ExpenseItemRequest{
					{Category: "生活費", Amount: 300000}, // exceeds income
				},
				CurrentSavings: []SavingsItemRequest{
					{Type: "deposit", Amount: 500000},
				},
			},
			// ValidateBusinessLogic writes 400 and returns nil, so the controller
			// continues and calls CreateFinancialPlan. We mock it to avoid panics.
			// The recorder already has status 400 from the first write.
			mockSetup: func(m *MockManageFinancialDataUseCase) {
				m.On("CreateFinancialPlan", mock.Anything, mock.Anything).Return(&usecases.CreateFinancialPlanOutput{
					UserID: entities.UserID("user-123"),
				}, nil)
				m.On("GetFinancialPlan", mock.Anything, mock.Anything).Return(&usecases.GetFinancialPlanOutput{Plan: nil}, nil)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "Error: internal server error",
			requestBody: validFinancialDataRequest(),
			mockSetup: func(m *MockManageFinancialDataUseCase) {
				m.On("CreateFinancialPlan", mock.Anything, mock.Anything).Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newFinancialDataEcho()
			mockUseCase := new(MockManageFinancialDataUseCase)
			tt.mockSetup(mockUseCase)
			controller := NewFinancialDataController(mockUseCase)

			reqJSON, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/financial-data", bytes.NewBuffer(reqJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := controller.CreateFinancialData(c)

			if tt.expectHandlerError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
			}
		})
	}
}

func TestGetFinancialData(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		mockSetup      func(m *MockManageFinancialDataUseCase)
		expectedStatus int
	}{
		{
			name:   "Success: get financial data",
			userID: "user-123",
			mockSetup: func(m *MockManageFinancialDataUseCase) {
				m.On("GetFinancialPlan", mock.Anything, usecases.GetFinancialPlanInput{
					UserID: entities.UserID("user-123"),
				}).Return(&usecases.GetFinancialPlanOutput{
					Plan: nil, // nil plan returns empty response gracefully
				}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Error: missing user_id",
			userID:         "",
			mockSetup:      func(m *MockManageFinancialDataUseCase) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Error: financial data not found",
			userID: "user-123",
			mockSetup: func(m *MockManageFinancialDataUseCase) {
				m.On("GetFinancialPlan", mock.Anything, mock.Anything).Return(nil, errors.New("財務データが見つかりません"))
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:   "Error: internal server error",
			userID: "user-123",
			mockSetup: func(m *MockManageFinancialDataUseCase) {
				m.On("GetFinancialPlan", mock.Anything, mock.Anything).Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newFinancialDataEcho()
			mockUseCase := new(MockManageFinancialDataUseCase)
			tt.mockSetup(mockUseCase)
			controller := NewFinancialDataController(mockUseCase)

			target := "/financial-data"
			if tt.userID != "" {
				target += "?user_id=" + tt.userID
			}
			req := httptest.NewRequest(http.MethodGet, target, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := controller.GetFinancialData(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}

func TestUpdateFinancialProfile(t *testing.T) {
	validUpdateRequest := UpdateFinancialProfileRequest{
		MonthlyIncome:    400000,
		InvestmentReturn: 5.0,
		InflationRate:    2.0,
		MonthlyExpenses: []ExpenseItemRequest{
			{Category: "生活費", Amount: 200000},
		},
		CurrentSavings: []SavingsItemRequest{
			{Type: "deposit", Amount: 500000},
		},
	}

	tests := []struct {
		name               string
		userID             string
		requestBody        interface{}
		mockSetup          func(m *MockManageFinancialDataUseCase)
		expectedStatus     int
		expectHandlerError bool
	}{
		{
			name:        "Success: update financial profile",
			userID:      "user-123",
			requestBody: validUpdateRequest,
			mockSetup: func(m *MockManageFinancialDataUseCase) {
				m.On("UpdateFinancialProfile", mock.Anything, mock.MatchedBy(func(input usecases.UpdateFinancialProfileInput) bool {
					return input.UserID == entities.UserID("user-123")
				})).Return(&usecases.UpdateFinancialProfileOutput{
					FinancialDataResponse: &usecases.FinancialDataResponse{UserID: "user-123"},
				}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Error: missing user_id in path",
			userID:         "",
			requestBody:    validUpdateRequest,
			mockSetup:      func(m *MockManageFinancialDataUseCase) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "Error: not found - fallback to create",
			userID:      "user-123",
			requestBody: validUpdateRequest,
			mockSetup: func(m *MockManageFinancialDataUseCase) {
				m.On("UpdateFinancialProfile", mock.Anything, mock.Anything).Return(nil, errors.New("財務データが見つかりません"))
				m.On("CreateFinancialPlan", mock.Anything, mock.Anything).Return(&usecases.CreateFinancialPlanOutput{
					UserID:    entities.UserID("user-123"),
					CreatedAt: "2030-01-01T00:00:00Z",
				}, nil)
				m.On("GetFinancialPlan", mock.Anything, mock.Anything).Return(&usecases.GetFinancialPlanOutput{
					Plan: nil,
				}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "Error: internal server error",
			userID:      "user-123",
			requestBody: validUpdateRequest,
			mockSetup: func(m *MockManageFinancialDataUseCase) {
				m.On("UpdateFinancialProfile", mock.Anything, mock.Anything).Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newFinancialDataEcho()
			mockUseCase := new(MockManageFinancialDataUseCase)
			tt.mockSetup(mockUseCase)
			controller := NewFinancialDataController(mockUseCase)

			reqJSON, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPut, "/financial-data/"+tt.userID+"/profile", bytes.NewBuffer(reqJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			if tt.userID != "" {
				c.SetParamNames("user_id")
				c.SetParamValues(tt.userID)
			}

			err := controller.UpdateFinancialProfile(c)

			if tt.expectHandlerError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
			}
		})
	}
}

func TestUpdateRetirementData(t *testing.T) {
	validRetirementRequest := UpdateRetirementDataRequest{
		RetirementAge:             65,
		MonthlyRetirementExpenses: 200000,
		PensionAmount:             100000,
	}

	tests := []struct {
		name               string
		userID             string
		requestBody        interface{}
		mockSetup          func(m *MockManageFinancialDataUseCase)
		expectedStatus     int
		expectHandlerError bool
	}{
		{
			name:        "Success: update retirement data",
			userID:      "user-123",
			requestBody: validRetirementRequest,
			mockSetup: func(m *MockManageFinancialDataUseCase) {
				m.On("UpdateRetirementData", mock.Anything, mock.MatchedBy(func(input usecases.UpdateRetirementDataInput) bool {
					return input.UserID == entities.UserID("user-123") && input.RetirementAge == 65
				})).Return(&usecases.UpdateRetirementDataOutput{
					FinancialDataResponse: &usecases.FinancialDataResponse{UserID: "user-123"},
				}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Error: missing user_id in path",
			userID:         "",
			requestBody:    validRetirementRequest,
			mockSetup:      func(m *MockManageFinancialDataUseCase) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Error: invalid retirement age (below minimum)",
			userID: "user-123",
			requestBody: UpdateRetirementDataRequest{
				RetirementAge:             30, // below 50
				MonthlyRetirementExpenses: 200000,
				PensionAmount:             100000,
			},
			mockSetup:          func(m *MockManageFinancialDataUseCase) {},
			expectHandlerError: true,
		},
		{
			name:        "Error: financial data not found",
			userID:      "user-123",
			requestBody: validRetirementRequest,
			mockSetup: func(m *MockManageFinancialDataUseCase) {
				m.On("UpdateRetirementData", mock.Anything, mock.Anything).Return(nil, errors.New("財務データが見つかりません"))
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:        "Error: internal server error",
			userID:      "user-123",
			requestBody: validRetirementRequest,
			mockSetup: func(m *MockManageFinancialDataUseCase) {
				m.On("UpdateRetirementData", mock.Anything, mock.Anything).Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newFinancialDataEcho()
			mockUseCase := new(MockManageFinancialDataUseCase)
			tt.mockSetup(mockUseCase)
			controller := NewFinancialDataController(mockUseCase)

			reqJSON, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPut, "/financial-data/"+tt.userID+"/retirement", bytes.NewBuffer(reqJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			if tt.userID != "" {
				c.SetParamNames("user_id")
				c.SetParamValues(tt.userID)
			}

			err := controller.UpdateRetirementData(c)

			if tt.expectHandlerError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
			}
		})
	}
}

func TestUpdateEmergencyFund(t *testing.T) {
	validEmergencyFundRequest := UpdateEmergencyFundRequest{
		TargetMonths:  6,
		CurrentAmount: 300000,
	}

	tests := []struct {
		name               string
		userID             string
		requestBody        interface{}
		mockSetup          func(m *MockManageFinancialDataUseCase)
		expectedStatus     int
		expectHandlerError bool
	}{
		{
			name:        "Success: update emergency fund",
			userID:      "user-123",
			requestBody: validEmergencyFundRequest,
			mockSetup: func(m *MockManageFinancialDataUseCase) {
				m.On("UpdateEmergencyFund", mock.Anything, mock.MatchedBy(func(input usecases.UpdateEmergencyFundInput) bool {
					return input.UserID == entities.UserID("user-123") && input.TargetMonths == 6
				})).Return(&usecases.UpdateEmergencyFundOutput{
					FinancialDataResponse: &usecases.FinancialDataResponse{UserID: "user-123"},
				}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Error: missing user_id in path",
			userID:         "",
			requestBody:    validEmergencyFundRequest,
			mockSetup:      func(m *MockManageFinancialDataUseCase) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Error: target months exceeds maximum",
			userID: "user-123",
			requestBody: UpdateEmergencyFundRequest{
				TargetMonths:  25, // exceeds 24
				CurrentAmount: 300000,
			},
			mockSetup:          func(m *MockManageFinancialDataUseCase) {},
			expectHandlerError: true,
		},
		{
			name:        "Error: financial data not found",
			userID:      "user-123",
			requestBody: validEmergencyFundRequest,
			mockSetup: func(m *MockManageFinancialDataUseCase) {
				m.On("UpdateEmergencyFund", mock.Anything, mock.Anything).Return(nil, errors.New("財務データが見つかりません"))
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:        "Error: internal server error",
			userID:      "user-123",
			requestBody: validEmergencyFundRequest,
			mockSetup: func(m *MockManageFinancialDataUseCase) {
				m.On("UpdateEmergencyFund", mock.Anything, mock.Anything).Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newFinancialDataEcho()
			mockUseCase := new(MockManageFinancialDataUseCase)
			tt.mockSetup(mockUseCase)
			controller := NewFinancialDataController(mockUseCase)

			reqJSON, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPut, "/financial-data/"+tt.userID+"/emergency-fund", bytes.NewBuffer(reqJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			if tt.userID != "" {
				c.SetParamNames("user_id")
				c.SetParamValues(tt.userID)
			}

			err := controller.UpdateEmergencyFund(c)

			if tt.expectHandlerError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
			}
		})
	}
}

func TestDeleteFinancialData(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		mockSetup      func(m *MockManageFinancialDataUseCase)
		expectedStatus int
	}{
		{
			name:   "Success: delete financial data",
			userID: "user-123",
			mockSetup: func(m *MockManageFinancialDataUseCase) {
				m.On("DeleteFinancialPlan", mock.Anything, usecases.DeleteFinancialPlanInput{
					UserID: entities.UserID("user-123"),
				}).Return(nil)
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "Error: missing user_id",
			userID:         "",
			mockSetup:      func(m *MockManageFinancialDataUseCase) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Error: financial data not found",
			userID: "user-123",
			mockSetup: func(m *MockManageFinancialDataUseCase) {
				m.On("DeleteFinancialPlan", mock.Anything, mock.Anything).Return(errors.New("財務データが見つかりません"))
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:   "Error: internal server error",
			userID: "user-123",
			mockSetup: func(m *MockManageFinancialDataUseCase) {
				m.On("DeleteFinancialPlan", mock.Anything, mock.Anything).Return(errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newFinancialDataEcho()
			mockUseCase := new(MockManageFinancialDataUseCase)
			tt.mockSetup(mockUseCase)
			controller := NewFinancialDataController(mockUseCase)

			req := httptest.NewRequest(http.MethodDelete, "/financial-data/"+tt.userID, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			if tt.userID != "" {
				c.SetParamNames("user_id")
				c.SetParamValues(tt.userID)
			}

			err := controller.DeleteFinancialData(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}
