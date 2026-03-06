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

// MockGenerateReportsUseCase is a mock implementation of GenerateReportsUseCase
type MockGenerateReportsUseCase struct {
	mock.Mock
}

func (m *MockGenerateReportsUseCase) GenerateFinancialSummaryReport(ctx context.Context, input usecases.FinancialSummaryReportInput) (*usecases.FinancialSummaryReportOutput, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecases.FinancialSummaryReportOutput), args.Error(1)
}

func (m *MockGenerateReportsUseCase) GenerateAssetProjectionReport(ctx context.Context, input usecases.AssetProjectionReportInput) (*usecases.AssetProjectionReportOutput, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecases.AssetProjectionReportOutput), args.Error(1)
}

func (m *MockGenerateReportsUseCase) GenerateGoalsProgressReport(ctx context.Context, input usecases.GoalsProgressReportInput) (*usecases.GoalsProgressReportOutput, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecases.GoalsProgressReportOutput), args.Error(1)
}

func (m *MockGenerateReportsUseCase) GenerateRetirementPlanReport(ctx context.Context, input usecases.RetirementPlanReportInput) (*usecases.RetirementPlanReportOutput, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecases.RetirementPlanReportOutput), args.Error(1)
}

func (m *MockGenerateReportsUseCase) GenerateComprehensiveReport(ctx context.Context, input usecases.ComprehensiveReportInput) (*usecases.ComprehensiveReportOutput, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecases.ComprehensiveReportOutput), args.Error(1)
}

func (m *MockGenerateReportsUseCase) ExportReportToPDF(ctx context.Context, input usecases.ExportReportInput) (*usecases.ExportReportOutput, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecases.ExportReportOutput), args.Error(1)
}

func newReportsTestContext(method, target string, body interface{}) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	var reqBody []byte
	if body != nil {
		reqBody, _ = json.Marshal(body)
	}
	req := httptest.NewRequest(method, target, bytes.NewBuffer(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return c, rec
}

func TestGenerateFinancialSummaryReport(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(m *MockGenerateReportsUseCase)
		expectedStatus int
	}{
		{
			name:        "Success: generate financial summary report",
			requestBody: FinancialSummaryReportRequest{UserID: "user-123"},
			mockSetup: func(m *MockGenerateReportsUseCase) {
				m.On("GenerateFinancialSummaryReport", mock.Anything, usecases.FinancialSummaryReportInput{
					UserID: entities.UserID("user-123"),
				}).Return(&usecases.FinancialSummaryReportOutput{
					Report:      usecases.FinancialSummaryReport{},
					GeneratedAt: "2030-01-01T00:00:00Z",
				}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Error: missing user_id",
			requestBody:    FinancialSummaryReportRequest{},
			mockSetup:      func(m *MockGenerateReportsUseCase) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "Error: internal server error",
			requestBody: FinancialSummaryReportRequest{UserID: "user-123"},
			mockSetup: func(m *MockGenerateReportsUseCase) {
				m.On("GenerateFinancialSummaryReport", mock.Anything, mock.Anything).Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := new(MockGenerateReportsUseCase)
			tt.mockSetup(mockUseCase)
			controller := NewReportsController(mockUseCase)

			c, rec := newReportsTestContext(http.MethodPost, "/reports/financial-summary", tt.requestBody)

			err := controller.GenerateFinancialSummaryReport(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}

func TestGenerateAssetProjectionReport(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(m *MockGenerateReportsUseCase)
		expectedStatus int
	}{
		{
			name:        "Success: valid years",
			requestBody: AssetProjectionReportRequest{UserID: "user-123", Years: 10},
			mockSetup: func(m *MockGenerateReportsUseCase) {
				m.On("GenerateAssetProjectionReport", mock.Anything, usecases.AssetProjectionReportInput{
					UserID: entities.UserID("user-123"),
					Years:  10,
				}).Return(&usecases.AssetProjectionReportOutput{
					Report:      usecases.AssetProjectionReport{},
					GeneratedAt: "2030-01-01T00:00:00Z",
				}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Error: years exceeds maximum (51)",
			requestBody:    AssetProjectionReportRequest{UserID: "user-123", Years: 51},
			mockSetup:      func(m *MockGenerateReportsUseCase) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Error: years is zero (fails required)",
			requestBody:    AssetProjectionReportRequest{UserID: "user-123", Years: 0},
			mockSetup:      func(m *MockGenerateReportsUseCase) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "Error: internal server error",
			requestBody: AssetProjectionReportRequest{UserID: "user-123", Years: 10},
			mockSetup: func(m *MockGenerateReportsUseCase) {
				m.On("GenerateAssetProjectionReport", mock.Anything, mock.Anything).Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := new(MockGenerateReportsUseCase)
			tt.mockSetup(mockUseCase)
			controller := NewReportsController(mockUseCase)

			c, rec := newReportsTestContext(http.MethodPost, "/reports/asset-projection", tt.requestBody)

			err := controller.GenerateAssetProjectionReport(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}

func TestGenerateGoalsProgressReport(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(m *MockGenerateReportsUseCase)
		expectedStatus int
	}{
		{
			name:        "Success: generate goals progress report",
			requestBody: GoalsProgressReportRequest{UserID: "user-123"},
			mockSetup: func(m *MockGenerateReportsUseCase) {
				m.On("GenerateGoalsProgressReport", mock.Anything, mock.Anything).Return(&usecases.GoalsProgressReportOutput{
					Report:      usecases.GoalsProgressReport{},
					GeneratedAt: "2030-01-01T00:00:00Z",
				}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Error: missing user_id",
			requestBody:    GoalsProgressReportRequest{},
			mockSetup:      func(m *MockGenerateReportsUseCase) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "Error: internal server error",
			requestBody: GoalsProgressReportRequest{UserID: "user-123"},
			mockSetup: func(m *MockGenerateReportsUseCase) {
				m.On("GenerateGoalsProgressReport", mock.Anything, mock.Anything).Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := new(MockGenerateReportsUseCase)
			tt.mockSetup(mockUseCase)
			controller := NewReportsController(mockUseCase)

			c, rec := newReportsTestContext(http.MethodPost, "/reports/goals-progress", tt.requestBody)

			err := controller.GenerateGoalsProgressReport(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}

func TestGenerateRetirementPlanReport(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(m *MockGenerateReportsUseCase)
		expectedStatus int
	}{
		{
			name:        "Success: generate retirement plan report",
			requestBody: RetirementPlanReportRequest{UserID: "user-123"},
			mockSetup: func(m *MockGenerateReportsUseCase) {
				m.On("GenerateRetirementPlanReport", mock.Anything, mock.Anything).Return(&usecases.RetirementPlanReportOutput{
					Report:      usecases.RetirementPlanReport{},
					GeneratedAt: "2030-01-01T00:00:00Z",
				}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Error: missing user_id",
			requestBody:    RetirementPlanReportRequest{},
			mockSetup:      func(m *MockGenerateReportsUseCase) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "Error: internal server error",
			requestBody: RetirementPlanReportRequest{UserID: "user-123"},
			mockSetup: func(m *MockGenerateReportsUseCase) {
				m.On("GenerateRetirementPlanReport", mock.Anything, mock.Anything).Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := new(MockGenerateReportsUseCase)
			tt.mockSetup(mockUseCase)
			controller := NewReportsController(mockUseCase)

			c, rec := newReportsTestContext(http.MethodPost, "/reports/retirement-plan", tt.requestBody)

			err := controller.GenerateRetirementPlanReport(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}

func TestGenerateComprehensiveReport(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(m *MockGenerateReportsUseCase)
		expectedStatus int
	}{
		{
			name:        "Success: valid request",
			requestBody: ComprehensiveReportRequest{UserID: "user-123", Years: 10},
			mockSetup: func(m *MockGenerateReportsUseCase) {
				m.On("GenerateComprehensiveReport", mock.Anything, usecases.ComprehensiveReportInput{
					UserID: entities.UserID("user-123"),
					Years:  10,
				}).Return(&usecases.ComprehensiveReportOutput{
					Report:      usecases.ComprehensiveReport{},
					GeneratedAt: "2030-01-01T00:00:00Z",
				}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Error: years exceeds maximum",
			requestBody:    ComprehensiveReportRequest{UserID: "user-123", Years: 51},
			mockSetup:      func(m *MockGenerateReportsUseCase) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "Error: internal server error",
			requestBody: ComprehensiveReportRequest{UserID: "user-123", Years: 10},
			mockSetup: func(m *MockGenerateReportsUseCase) {
				m.On("GenerateComprehensiveReport", mock.Anything, mock.Anything).Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := new(MockGenerateReportsUseCase)
			tt.mockSetup(mockUseCase)
			controller := NewReportsController(mockUseCase)

			c, rec := newReportsTestContext(http.MethodPost, "/reports/comprehensive", tt.requestBody)

			err := controller.GenerateComprehensiveReport(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}

func TestExportReportToPDF(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(m *MockGenerateReportsUseCase)
		expectedStatus int
	}{
		{
			name: "Success: export report",
			requestBody: ExportReportRequest{
				UserID:     "user-123",
				ReportType: "financial_summary",
				Format:     "pdf",
				ReportData: map[string]interface{}{"key": "value"},
			},
			mockSetup: func(m *MockGenerateReportsUseCase) {
				m.On("ExportReportToPDF", mock.Anything, mock.Anything).Return(&usecases.ExportReportOutput{
					FileName:    "report.pdf",
					FileSize:    1024,
					DownloadURL: "https://example.com/report.pdf",
					ExpiresAt:   "2030-01-01T00:00:00Z",
				}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Error: invalid report_type",
			requestBody: ExportReportRequest{
				UserID:     "user-123",
				ReportType: "invalid_type",
				Format:     "pdf",
				ReportData: map[string]interface{}{"key": "value"},
			},
			mockSetup:      func(m *MockGenerateReportsUseCase) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Error: missing user_id",
			requestBody: ExportReportRequest{
				ReportType: "financial_summary",
				Format:     "pdf",
				ReportData: map[string]interface{}{"key": "value"},
			},
			mockSetup:      func(m *MockGenerateReportsUseCase) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Error: internal server error",
			requestBody: ExportReportRequest{
				UserID:     "user-123",
				ReportType: "financial_summary",
				Format:     "pdf",
				ReportData: map[string]interface{}{"key": "value"},
			},
			mockSetup: func(m *MockGenerateReportsUseCase) {
				m.On("ExportReportToPDF", mock.Anything, mock.Anything).Return(nil, errors.New("export error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := new(MockGenerateReportsUseCase)
			tt.mockSetup(mockUseCase)
			controller := NewReportsController(mockUseCase)

			c, rec := newReportsTestContext(http.MethodPost, "/reports/export", tt.requestBody)

			err := controller.ExportReportToPDF(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}

func TestGetReportPDF(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    map[string]string
		mockSetup      func(m *MockGenerateReportsUseCase)
		expectedStatus int
	}{
		{
			name: "Success: comprehensive report (default)",
			queryParams: map[string]string{
				"user_id": "user-123",
			},
			mockSetup: func(m *MockGenerateReportsUseCase) {
				m.On("GenerateComprehensiveReport", mock.Anything, mock.Anything).Return(&usecases.ComprehensiveReportOutput{
					Report:      usecases.ComprehensiveReport{},
					GeneratedAt: "2030-01-01T00:00:00Z",
				}, nil)
				m.On("ExportReportToPDF", mock.Anything, mock.Anything).Return(&usecases.ExportReportOutput{
					FileName:    "report.pdf",
					DownloadURL: "https://example.com/report.pdf",
					ExpiresAt:   "2030-01-01T00:00:00Z",
				}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Success: financial_summary report",
			queryParams: map[string]string{
				"user_id":     "user-123",
				"report_type": "financial_summary",
			},
			mockSetup: func(m *MockGenerateReportsUseCase) {
				m.On("GenerateFinancialSummaryReport", mock.Anything, mock.Anything).Return(&usecases.FinancialSummaryReportOutput{
					Report:      usecases.FinancialSummaryReport{},
					GeneratedAt: "2030-01-01T00:00:00Z",
				}, nil)
				m.On("ExportReportToPDF", mock.Anything, mock.Anything).Return(&usecases.ExportReportOutput{
					FileName:    "report.pdf",
					DownloadURL: "https://example.com/report.pdf",
					ExpiresAt:   "2030-01-01T00:00:00Z",
				}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Error: missing user_id",
			queryParams:    map[string]string{},
			mockSetup:      func(m *MockGenerateReportsUseCase) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Error: unsupported report type",
			queryParams: map[string]string{
				"user_id":     "user-123",
				"report_type": "unsupported_type",
			},
			mockSetup:      func(m *MockGenerateReportsUseCase) {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			e.Validator = &CustomValidator{validator: validator.New()}
			mockUseCase := new(MockGenerateReportsUseCase)
			tt.mockSetup(mockUseCase)
			controller := NewReportsController(mockUseCase)

			target := "/reports/pdf"
			if len(tt.queryParams) > 0 {
				target += "?"
				for k, v := range tt.queryParams {
					target += k + "=" + v + "&"
				}
			}
			req := httptest.NewRequest(http.MethodGet, target, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := controller.GetReportPDF(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}

func TestDownloadReport(t *testing.T) {
	tests := []struct {
		name           string
		token          string
		expectedStatus int
	}{
		{
			name:           "Success: valid token",
			token:          "valid-download-token",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Error: empty token",
			token:          "",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			mockUseCase := new(MockGenerateReportsUseCase)
			controller := NewReportsController(mockUseCase)

			req := httptest.NewRequest(http.MethodGet, "/reports/download/"+tt.token, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			if tt.token != "" {
				c.SetParamNames("token")
				c.SetParamValues(tt.token)
			}

			err := controller.DownloadReport(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}
