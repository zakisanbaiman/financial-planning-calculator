package web

import (
	"fmt"
	"net/http"
	"reflect"
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
	v := validator.New()

	// Use json tag names as field names so that validation errors report
	// camelCase field names (e.g. "email") instead of struct field names
	// (e.g. "Email").
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" || name == "" {
			return fld.Name
		}
		return name
	})

	// Register custom validation messages
	registerCustomMessages(v)

	return &CustomValidator{
		validator: v,
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

// getFieldDisplayName returns a user-friendly field name in Japanese.
// Keys are json tag names (snake_case / camelCase) as reported by the validator
// after RegisterTagNameFunc is applied.
func getFieldDisplayName(field string) string {
	fieldNames := map[string]string{
		// User and identification fields
		"user_id":   "ユーザーID",
		"goal_id":   "目標ID",
		"plan_id":   "計画ID",
		"report_id": "レポートID",

		// Financial profile fields
		"monthly_income":    "月収",
		"monthly_expenses":  "月間支出",
		"current_savings":   "現在の貯蓄",
		"investment_return": "投資利回り",
		"inflation_rate":    "インフレ率",

		// Retirement fields
		"retirement_age":              "退職年齢",
		"monthly_retirement_expenses": "老後月間生活費",
		"pension_amount":              "年金受給額",
		"current_age":                 "現在の年齢",
		"life_expectancy":             "平均寿命",

		// Emergency fund fields
		"emergency_fund_target_months":  "緊急資金目標月数",
		"emergency_fund_current_amount": "現在の緊急資金",
		"target_months":                 "目標月数",

		// Expense and savings fields
		"category":    "カテゴリ",
		"amount":      "金額",
		"type":        "種類",
		"description": "説明",

		// Goal fields
		"goal_type":            "目標タイプ",
		"title":                "タイトル",
		"target_amount":        "目標金額",
		"target_date":          "目標日",
		"current_amount":       "現在の金額",
		"monthly_contribution": "月間積立額",
		"is_active":            "アクティブ状態",
		"note":                 "メモ",

		// Calculation fields
		"years":      "年数",
		"months":     "月数",
		"percentage": "パーセンテージ",
		"rate":       "利率",
		"period":     "期間",
		"start_date": "開始日",
		"end_date":   "終了日",
		"created_at": "作成日時",
		"updated_at": "更新日時",

		// Report fields
		"report_type": "レポートタイプ",
		"format":      "フォーマット",
		"language":    "言語",
		"template":    "テンプレート",

		// Auth fields
		"password":     "パスワード",
		"new_password": "新しいパスワード",
		"token":        "トークン",

		// Common fields
		"name":    "名前",
		"value":   "値",
		"status":  "ステータス",
		"message": "メッセージ",
		"email":   "メールアドレス",
		"phone":   "電話番号",
		"address": "住所",
	}

	if displayName, exists := fieldNames[field]; exists {
		return displayName
	}

	// Return field as-is if not found in map
	return field
}
