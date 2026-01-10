package controllers

import (
	"net/http"

	"github.com/financial-planning-calculator/backend/application/usecases"
	"github.com/labstack/echo/v4"
)

// AuthController は認証関連のコントローラー
type AuthController struct {
	authUseCase usecases.AuthUseCase
}

// NewAuthController は新しいAuthControllerを作成する
func NewAuthController(authUseCase usecases.AuthUseCase) *AuthController {
	return &AuthController{
		authUseCase: authUseCase,
	}
}

// RegisterRequest はユーザー登録リクエスト
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// LoginRequest はログインリクエスト
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// AuthResponse は認証レスポンス
type AuthResponse struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at"`
}

// Register は新しいユーザーを登録する
// @Summary ユーザー登録
// @Description 新しいユーザーを登録し、JWTトークンを発行します
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "ユーザー登録リクエスト"
// @Success 201 {object} AuthResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse "メールアドレスが既に登録されています"
// @Failure 500 {object} ErrorResponse
// @Router /auth/register [post]
func (c *AuthController) Register(ctx echo.Context) error {
	var req RegisterRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, NewErrorResponse(ctx, ErrorCodeBadRequest, "リクエストの解析に失敗しました", err.Error()))
	}

	if err := ctx.Validate(&req); err != nil {
		return err // Validator already returns proper error response
	}

	// ユーザー登録
	input := usecases.RegisterInput{
		Email:    req.Email,
		Password: req.Password,
	}

	output, err := c.authUseCase.Register(ctx.Request().Context(), input)
	if err != nil {
		// メールアドレス重複エラー
		if err.Error() == "このメールアドレスは既に登録されています" {
			return ctx.JSON(http.StatusConflict, NewErrorResponse(ctx, ErrorCodeConflict, err.Error(), nil))
		}
		// パスワードやメールアドレスのバリデーションエラー
		if err.Error() == "無効なメールアドレス形式です" || err.Error() == "パスワードは8文字以上である必要があります" {
			return ctx.JSON(http.StatusBadRequest, NewErrorResponse(ctx, ErrorCodeValidation, err.Error(), nil))
		}
		// その他のエラー
		return ctx.JSON(http.StatusInternalServerError, NewErrorResponse(ctx, ErrorCodeInternalServer, "ユーザー登録に失敗しました", err.Error()))
	}

	response := AuthResponse{
		UserID:    output.UserID,
		Email:     output.Email,
		Token:     output.Token,
		ExpiresAt: output.ExpiresAt,
	}

	return ctx.JSON(http.StatusCreated, response)
}

// Login はユーザー認証を行う
// @Summary ログイン
// @Description ユーザー認証を行い、JWTトークンを発行します
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "ログインリクエスト"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse "メールアドレスまたはパスワードが正しくありません"
// @Failure 500 {object} ErrorResponse
// @Router /auth/login [post]
func (c *AuthController) Login(ctx echo.Context) error {
	var req LoginRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, NewErrorResponse(ctx, ErrorCodeBadRequest, "リクエストの解析に失敗しました", err.Error()))
	}

	if err := ctx.Validate(&req); err != nil {
		return err // Validator already returns proper error response
	}

	// ログイン
	input := usecases.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	}

	output, err := c.authUseCase.Login(ctx.Request().Context(), input)
	if err != nil {
		// 認証エラー（メールアドレスまたはパスワードが間違っている）
		if err.Error() == "メールアドレスまたはパスワードが正しくありません" {
			return ctx.JSON(http.StatusUnauthorized, NewErrorResponse(ctx, ErrorCodeUnauthorized, err.Error(), nil))
		}
		// その他のエラー
		return ctx.JSON(http.StatusInternalServerError, NewErrorResponse(ctx, ErrorCodeInternalServer, "ログインに失敗しました", err.Error()))
	}

	response := AuthResponse{
		UserID:    output.UserID,
		Email:     output.Email,
		Token:     output.Token,
		ExpiresAt: output.ExpiresAt,
	}

	return ctx.JSON(http.StatusOK, response)
}
