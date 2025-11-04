package controllers

import (
	"time"

	"github.com/labstack/echo/v4"
)

// ErrorResponse は統一されたエラーレスポンス形式
type ErrorResponse struct {
	Error     string      `json:"error"`
	Details   interface{} `json:"details,omitempty"`
	Timestamp string      `json:"timestamp"`
	RequestID string      `json:"request_id,omitempty"`
	Code      string      `json:"code,omitempty"`
}

// ErrorCode represents different types of errors
type ErrorCode string

const (
	ErrorCodeValidation         ErrorCode = "VALIDATION_ERROR"
	ErrorCodeNotFound           ErrorCode = "NOT_FOUND"
	ErrorCodeInternalServer     ErrorCode = "INTERNAL_SERVER_ERROR"
	ErrorCodeBadRequest         ErrorCode = "BAD_REQUEST"
	ErrorCodeBusinessLogic      ErrorCode = "BUSINESS_LOGIC_ERROR"
	ErrorCodeUnauthorized       ErrorCode = "UNAUTHORIZED"
	ErrorCodeForbidden          ErrorCode = "FORBIDDEN"
	ErrorCodeConflict           ErrorCode = "CONFLICT"
	ErrorCodeTooManyRequests    ErrorCode = "TOO_MANY_REQUESTS"
	ErrorCodeServiceUnavailable ErrorCode = "SERVICE_UNAVAILABLE"
	ErrorCodeTimeout            ErrorCode = "TIMEOUT"
	ErrorCodeDataIntegrity      ErrorCode = "DATA_INTEGRITY_ERROR"
	ErrorCodeCalculation        ErrorCode = "CALCULATION_ERROR"
	ErrorCodeInsufficientData   ErrorCode = "INSUFFICIENT_DATA"
)

// BusinessLogicError represents business logic validation errors
type BusinessLogicError struct {
	Type          string      `json:"type"`
	Field         string      `json:"field,omitempty"`
	Message       string      `json:"message"`
	Suggestion    string      `json:"suggestion,omitempty"`
	CurrentValue  interface{} `json:"current_value,omitempty"`
	ExpectedValue interface{} `json:"expected_value,omitempty"`
	Severity      string      `json:"severity,omitempty"` // "error", "warning", "info"
	HelpURL       string      `json:"help_url,omitempty"`
}

// NewErrorResponse creates a new error response with timestamp and request ID
func NewErrorResponse(ctx echo.Context, code ErrorCode, message string, details interface{}) ErrorResponse {
	requestID := ctx.Response().Header().Get(echo.HeaderXRequestID)
	if requestID == "" {
		requestID = ctx.Request().Header.Get("X-Request-ID")
	}

	return ErrorResponse{
		Error:     message,
		Details:   details,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		RequestID: requestID,
		Code:      string(code),
	}
}

// NewValidationErrorResponse creates a validation error response
func NewValidationErrorResponse(ctx echo.Context, details interface{}) ErrorResponse {
	return NewErrorResponse(ctx, ErrorCodeValidation, "入力値が無効です", details)
}

// NewBusinessLogicErrorResponse creates a business logic error response
func NewBusinessLogicErrorResponse(ctx echo.Context, errors []BusinessLogicError) ErrorResponse {
	return NewErrorResponse(ctx, ErrorCodeBusinessLogic, "ビジネスロジックエラーが発生しました", errors)
}

// NewNotFoundErrorResponse creates a not found error response
func NewNotFoundErrorResponse(ctx echo.Context, resource string) ErrorResponse {
	return NewErrorResponse(ctx, ErrorCodeNotFound, resource+"が見つかりません", nil)
}

// NewInternalServerErrorResponse creates an internal server error response
func NewInternalServerErrorResponse(ctx echo.Context, details string) ErrorResponse {
	return NewErrorResponse(ctx, ErrorCodeInternalServer, "内部サーバーエラーが発生しました", details)
}

// NewConflictErrorResponse creates a conflict error response
func NewConflictErrorResponse(ctx echo.Context, resource string) ErrorResponse {
	return NewErrorResponse(ctx, ErrorCodeConflict, resource+"が既に存在します", nil)
}

// NewCalculationErrorResponse creates a calculation error response
func NewCalculationErrorResponse(ctx echo.Context, details string) ErrorResponse {
	return NewErrorResponse(ctx, ErrorCodeCalculation, "計算処理でエラーが発生しました", details)
}

// NewInsufficientDataErrorResponse creates an insufficient data error response
func NewInsufficientDataErrorResponse(ctx echo.Context, missingData string) ErrorResponse {
	return NewErrorResponse(ctx, ErrorCodeInsufficientData, "計算に必要なデータが不足しています", map[string]string{
		"missing_data": missingData,
		"suggestion":   "必要なデータを入力してから再度お試しください",
	})
}

// NewDataIntegrityErrorResponse creates a data integrity error response
func NewDataIntegrityErrorResponse(ctx echo.Context, details string) ErrorResponse {
	return NewErrorResponse(ctx, ErrorCodeDataIntegrity, "データの整合性エラーが発生しました", details)
}

// ValidateBusinessLogic validates business logic and returns errors if any
func ValidateBusinessLogic(ctx echo.Context, validations ...func() *BusinessLogicError) error {
	var errors []BusinessLogicError

	for _, validation := range validations {
		if err := validation(); err != nil {
			errors = append(errors, *err)
		}
	}

	if len(errors) > 0 {
		response := NewBusinessLogicErrorResponse(ctx, errors)
		return ctx.JSON(400, response)
	}

	return nil
}

// CreateBusinessLogicError creates a business logic error
func CreateBusinessLogicError(errorType, message, suggestion string, currentValue, expectedValue interface{}) *BusinessLogicError {
	return &BusinessLogicError{
		Type:          errorType,
		Message:       message,
		Suggestion:    suggestion,
		CurrentValue:  currentValue,
		ExpectedValue: expectedValue,
		Severity:      "error",
	}
}

// CreateBusinessLogicErrorWithField creates a business logic error with field information
func CreateBusinessLogicErrorWithField(errorType, field, message, suggestion string, currentValue, expectedValue interface{}) *BusinessLogicError {
	return &BusinessLogicError{
		Type:          errorType,
		Field:         field,
		Message:       message,
		Suggestion:    suggestion,
		CurrentValue:  currentValue,
		ExpectedValue: expectedValue,
		Severity:      "error",
	}
}

// CreateBusinessLogicWarning creates a business logic warning
func CreateBusinessLogicWarning(errorType, message, suggestion string, currentValue, expectedValue interface{}) *BusinessLogicError {
	return &BusinessLogicError{
		Type:          errorType,
		Message:       message,
		Suggestion:    suggestion,
		CurrentValue:  currentValue,
		ExpectedValue: expectedValue,
		Severity:      "warning",
	}
}

// CreateBusinessLogicInfo creates a business logic info message
func CreateBusinessLogicInfo(errorType, message, suggestion string, currentValue, expectedValue interface{}) *BusinessLogicError {
	return &BusinessLogicError{
		Type:          errorType,
		Message:       message,
		Suggestion:    suggestion,
		CurrentValue:  currentValue,
		ExpectedValue: expectedValue,
		Severity:      "info",
	}
}
