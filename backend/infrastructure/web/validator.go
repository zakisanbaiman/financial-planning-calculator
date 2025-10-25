package web

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// CustomValidator wraps the validator instance
type CustomValidator struct {
	validator *validator.Validate
}

// NewCustomValidator creates a new custom validator
func NewCustomValidator() *CustomValidator {
	return &CustomValidator{
		validator: validator.New(),
	}
}

// Validate validates the struct
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Return a custom error response
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
