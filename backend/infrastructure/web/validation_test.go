package web

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// TestCustomValidator tests the custom validator functionality
func TestCustomValidator(t *testing.T) {
	// Create Echo instance with custom validator
	e := echo.New()
	e.Validator = NewCustomValidator()

	// Test struct for validation
	type TestRequest struct {
		MonthlyIncome    float64 `json:"monthly_income" validate:"required,gt=0"`
		InvestmentReturn float64 `json:"investment_return" validate:"required,gte=0,lte=100"`
		Email            string  `json:"email" validate:"required,email"`
		Category         string  `json:"category" validate:"required,oneof=food housing transport"`
	}

	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedError  bool
		expectedFields []string
	}{
		{
			name:           "Valid request",
			requestBody:    `{"monthly_income":400000,"investment_return":5.0,"email":"test@example.com","category":"food"}`,
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "Invalid monthly income (negative)",
			requestBody:    `{"monthly_income":-100,"investment_return":5.0,"email":"test@example.com","category":"food"}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			expectedFields: []string{"MonthlyIncome"},
		},
		{
			name:           "Invalid investment return (over 100)",
			requestBody:    `{"monthly_income":400000,"investment_return":150.0,"email":"test@example.com","category":"food"}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			expectedFields: []string{"InvestmentReturn"},
		},
		{
			name:           "Invalid email format",
			requestBody:    `{"monthly_income":400000,"investment_return":5.0,"email":"invalid-email","category":"food"}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			expectedFields: []string{"Email"},
		},
		{
			name:           "Invalid category",
			requestBody:    `{"monthly_income":400000,"investment_return":5.0,"email":"test@example.com","category":"invalid"}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			expectedFields: []string{"Category"},
		},
		{
			name:           "Multiple validation errors",
			requestBody:    `{"monthly_income":-100,"investment_return":150.0,"email":"invalid-email","category":"invalid"}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			expectedFields: []string{"MonthlyIncome", "InvestmentReturn", "Email", "Category"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test handler
			handler := func(c echo.Context) error {
				var req TestRequest
				if err := c.Bind(&req); err != nil {
					return err
				}
				if err := c.Validate(&req); err != nil {
					return err
				}
				return c.JSON(http.StatusOK, map[string]string{"status": "success"})
			}

			// Create request
			req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(tt.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Execute handler
			err := handler(c)

			if tt.expectedError {
				assert.Error(t, err)
				if httpErr, ok := err.(*echo.HTTPError); ok {
					assert.Equal(t, tt.expectedStatus, httpErr.Code)

					// Check if it's our validation error format
					if validationErr, ok := httpErr.Message.(ValidationErrorResponse); ok {
						assert.Equal(t, "入力値が無効です", validationErr.Error)
						assert.NotEmpty(t, validationErr.Details)

						// Check that expected fields are in the validation errors
						fieldMap := make(map[string]bool)
						for _, detail := range validationErr.Details {
							fieldMap[detail.Field] = true
						}

						for _, expectedField := range tt.expectedFields {
							assert.True(t, fieldMap[expectedField], "Expected field %s not found in validation errors", expectedField)
						}
					}
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
			}
		})
	}
}

// TestGetFieldDisplayName tests the field display name mapping
func TestGetFieldDisplayName(t *testing.T) {
	tests := []struct {
		field    string
		expected string
	}{
		{"MonthlyIncome", "月収"},
		{"InvestmentReturn", "投資利回り"},
		{"InflationRate", "インフレ率"},
		{"RetirementAge", "退職年齢"},
		{"UnknownField", "unknownfield"},
	}

	for _, tt := range tests {
		t.Run(tt.field, func(t *testing.T) {
			result := getFieldDisplayName(tt.field)
			assert.Equal(t, tt.expected, result)
		})
	}
}
