package web

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// CustomValidator wraps the validator instance
type CustomValidator struct {
	validator *validator.Validate
}

// ValidationError represents a validation error with detailed information
type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

// ValidationErrorResponse represents the response for validation errors
type ValidationErrorResponse struct {
	Error   string            `json:"error"`
	Details []ValidationError `json:"details"`
}

// NewCustomValidator creates a new custom validator
func NewCustomValidator() *CustomValidator {
	validator := validator.New()

	// Register custom validation messages
	registerCustomMessages(validator)

	return &CustomValidator{
		validator: validator,
	}
}

// Validate validates the struct
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		var validationErrors []ValidationError

		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			for _, validationErr := range validationErrs {
				validationErrors = append(validationErrors, ValidationError{
					Field:   validationErr.Field(),
					Tag:     validationErr.Tag(),
					Value:   fmt.Sprintf("%v", validationErr.Value()),
					Message: getCustomErrorMessage(validationErr),
				})
			}
		}

		// Return structured validation error
		return &echo.HTTPError{
			Code: http.StatusBadRequest,
			Message: ValidationErrorResponse{
				Error:   "入力値が無効です",
				Details: validationErrors,
			},
		}
	}
	return nil
}

// registerCustomMessages registers custom validation messages
func registerCustomMessages(validator *validator.Validate) {
	// Custom validation rules can be registered here if needed
}

// getCustomErrorMessage returns a custom error message for validation errors
func getCustomErrorMessage(fe validator.FieldError) string {
	field := fe.Field()
	tag := fe.Tag()
	param := fe.Param()

	switch tag {
	case "required":
		return fmt.Sprintf("%sは必須です", getFieldDisplayName(field))
	case "gt":
		return fmt.Sprintf("%sは%sより大きい値を入力してください", getFieldDisplayName(field), param)
	case "gte":
		return fmt.Sprintf("%sは%s以上の値を入力してください", getFieldDisplayName(field), param)
	case "lt":
		return fmt.Sprintf("%sは%sより小さい値を入力してください", getFieldDisplayName(field), param)
	case "lte":
		return fmt.Sprintf("%sは%s以下の値を入力してください", getFieldDisplayName(field), param)
	case "min":
		return fmt.Sprintf("%sは%s文字以上で入力してください", getFieldDisplayName(field), param)
	case "max":
		return fmt.Sprintf("%sは%s文字以下で入力してください", getFieldDisplayName(field), param)
	case "oneof":
		return fmt.Sprintf("%sは有効な値を選択してください（%s）", getFieldDisplayName(field), param)
	case "email":
		return fmt.Sprintf("%sは有効なメールアドレスを入力してください", getFieldDisplayName(field))
	case "uuid":
		return fmt.Sprintf("%sは有効なUUID形式で入力してください", getFieldDisplayName(field))
	case "dive":
		return fmt.Sprintf("%sの項目に無効な値が含まれています", getFieldDisplayName(field))
	case "numeric":
		return fmt.Sprintf("%sは数値で入力してください", getFieldDisplayName(field))
	case "alpha":
		return fmt.Sprintf("%sは英字のみで入力してください", getFieldDisplayName(field))
	case "alphanum":
		return fmt.Sprintf("%sは英数字のみで入力してください", getFieldDisplayName(field))
	case "len":
		return fmt.Sprintf("%sは%s文字で入力してください", getFieldDisplayName(field), param)
	case "url":
		return fmt.Sprintf("%sは有効なURL形式で入力してください", getFieldDisplayName(field))
	case "datetime":
		return fmt.Sprintf("%sは有効な日時形式で入力してください", getFieldDisplayName(field))
	default:
		return fmt.Sprintf("%sの値が無効です", getFieldDisplayName(field))
	}
}

// getFieldDisplayName returns a user-friendly field name in Japanese
func getFieldDisplayName(field string) string {
	fieldNames := map[string]string{
		// User and identification fields
		"UserID":   "ユーザーID",
		"GoalID":   "目標ID",
		"PlanID":   "計画ID",
		"ReportID": "レポートID",

		// Financial profile fields
		"MonthlyIncome":    "月収",
		"MonthlyExpenses":  "月間支出",
		"CurrentSavings":   "現在の貯蓄",
		"InvestmentReturn": "投資利回り",
		"InflationRate":    "インフレ率",

		// Retirement fields
		"RetirementAge":             "退職年齢",
		"MonthlyRetirementExpenses": "老後月間生活費",
		"PensionAmount":             "年金受給額",
		"CurrentAge":                "現在の年齢",
		"LifeExpectancy":            "平均寿命",

		// Emergency fund fields
		"EmergencyFundTargetMonths":  "緊急資金目標月数",
		"EmergencyFundCurrentAmount": "現在の緊急資金",
		"TargetMonths":               "目標月数",

		// Expense and savings fields
		"Category":    "カテゴリ",
		"Amount":      "金額",
		"Type":        "種類",
		"Description": "説明",

		// Goal fields
		"GoalType":            "目標タイプ",
		"Title":               "タイトル",
		"TargetAmount":        "目標金額",
		"TargetDate":          "目標日",
		"CurrentAmount":       "現在の金額",
		"MonthlyContribution": "月間積立額",
		"IsActive":            "アクティブ状態",
		"Note":                "メモ",

		// Calculation fields
		"Years":      "年数",
		"Months":     "月数",
		"Percentage": "パーセンテージ",
		"Rate":       "利率",
		"Period":     "期間",
		"StartDate":  "開始日",
		"EndDate":    "終了日",
		"CreatedAt":  "作成日時",
		"UpdatedAt":  "更新日時",

		// Report fields
		"ReportType": "レポートタイプ",
		"Format":     "フォーマット",
		"Language":   "言語",
		"Template":   "テンプレート",

		// Common fields
		"Name":    "名前",
		"Value":   "値",
		"Status":  "ステータス",
		"Message": "メッセージ",
		"Email":   "メールアドレス",
		"Phone":   "電話番号",
		"Address": "住所",
	}

	if displayName, exists := fieldNames[field]; exists {
		return displayName
	}

	// Convert camelCase to readable format if not found in map
	return strings.ToLower(field)
}
