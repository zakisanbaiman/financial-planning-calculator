package controllers

import (
	"net/http"

	"github.com/financial-planning-calculator/backend/application/usecases"
	"github.com/financial-planning-calculator/backend/config"
	"github.com/labstack/echo/v4"
)

// AuthController は認証関連のコントローラー
type AuthController struct {
	authUseCase  usecases.AuthUseCase
	serverConfig *config.ServerConfig
}

// NewAuthController は新しいAuthControllerを作成する
func NewAuthController(authUseCase usecases.AuthUseCase, serverConfig *config.ServerConfig) *AuthController {
	return &AuthController{
		authUseCase:  authUseCase,
		serverConfig: serverConfig,
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
	UserID       string `json:"user_id"`
	Email        string `json:"email"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    string `json:"expires_at"`
}

// RefreshRequest はトークンリフレッシュリクエスト
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// RefreshResponse はトークンリフレッシュレスポンス
type RefreshResponse struct {
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

	// トークンをhttpOnly Cookieに設定
	setAuthCookies(ctx, output.Token, output.RefreshToken, c.serverConfig)

	response := AuthResponse{
		UserID:       output.UserID,
		Email:        output.Email,
		Token:        output.Token,
		RefreshToken: output.RefreshToken,
		ExpiresAt:    output.ExpiresAt,
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

	// 2FA検証が必要な場合（RefreshTokenが空）は仮トークンのみをCookieに設定
	if output.RefreshToken == "" {
		// 2FA仮トークンをアクセストークンCookieに設定（5分間有効）
		setAccessTokenCookie(ctx, output.Token, c.serverConfig)
	} else {
		// 通常のトークンをhttpOnly Cookieに設定
		setAuthCookies(ctx, output.Token, output.RefreshToken, c.serverConfig)
	}

	response := AuthResponse{
		UserID:       output.UserID,
		Email:        output.Email,
		Token:        output.Token,
		RefreshToken: output.RefreshToken,
		ExpiresAt:    output.ExpiresAt,
	}

	return ctx.JSON(http.StatusOK, response)
}

// Refresh はリフレッシュトークンを使用して新しいアクセストークンを発行する
// @Summary トークンリフレッシュ
// @Description リフレッシュトークンを使用して新しいアクセストークンを発行します（Cookieまたはリクエストボディから取得）
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshRequest false "リフレッシュリクエスト（Cookieがない場合のみ必要）"
// @Success 200 {object} RefreshResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse "無効なリフレッシュトークンです"
// @Failure 500 {object} ErrorResponse
// @Router /auth/refresh [post]
func (c *AuthController) Refresh(ctx echo.Context) error {
	var refreshToken string

	// まずCookieからリフレッシュトークンを取得
	cookie, err := ctx.Cookie("refresh_token")
	if err == nil && cookie.Value != "" {
		refreshToken = cookie.Value
	} else {
		// Cookieにない場合はリクエストボディから取得（後方互換性のため）
		var req RefreshRequest
		if err := ctx.Bind(&req); err != nil {
			return ctx.JSON(http.StatusBadRequest, NewErrorResponse(ctx, ErrorCodeBadRequest, "リクエストの解析に失敗しました", err.Error()))
		}

		if err := ctx.Validate(&req); err != nil {
			return err // Validator already returns proper error response
		}

		refreshToken = req.RefreshToken
	}

	if refreshToken == "" {
		return ctx.JSON(http.StatusBadRequest, NewErrorResponse(ctx, ErrorCodeBadRequest, "リフレッシュトークンが必要です", nil))
	}

	// トークンリフレッシュ
	output, err := c.authUseCase.RefreshAccessToken(ctx.Request().Context(), refreshToken)
	if err != nil {
		// リフレッシュトークンが無効または期限切れ
		if err.Error() == "無効なリフレッシュトークンです" || err.Error() == "リフレッシュトークンの有効期限が切れているか、失効されています" {
			return ctx.JSON(http.StatusUnauthorized, NewErrorResponse(ctx, ErrorCodeUnauthorized, err.Error(), nil))
		}
		// その他のエラー
		return ctx.JSON(http.StatusInternalServerError, NewErrorResponse(ctx, ErrorCodeInternalServer, "トークンリフレッシュに失敗しました", err.Error()))
	}

	// トークンをhttpOnly Cookieに設定（アクセストークンのみ更新）
	setAccessTokenCookie(ctx, output.Token, c.serverConfig)

	response := RefreshResponse{
		Token:     output.Token,
		ExpiresAt: output.ExpiresAt,
	}

	return ctx.JSON(http.StatusOK, response)
}

// setAuthCookies はアクセストークンとリフレッシュトークンをCookieに設定する
func setAuthCookies(ctx echo.Context, accessToken, refreshToken string, config *config.ServerConfig) {
	// アクセストークンをCookieに設定
	ctx.SetCookie(&http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		MaxAge:   int(config.JWTExpiration.Seconds()),
		HttpOnly: true,
		Secure:   config.CookieSecure,
		SameSite: http.SameSiteStrictMode,
	})

	// リフレッシュトークンをCookieに設定
	ctx.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		MaxAge:   int(config.RefreshTokenExpiration.Seconds()),
		HttpOnly: true,
		Secure:   config.CookieSecure,
		SameSite: http.SameSiteStrictMode,
	})
}

// setAccessTokenCookie はアクセストークンのみをCookieに設定する
func setAccessTokenCookie(ctx echo.Context, accessToken string, config *config.ServerConfig) {
	ctx.SetCookie(&http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		MaxAge:   int(config.JWTExpiration.Seconds()),
		HttpOnly: true,
		Secure:   config.CookieSecure,
		SameSite: http.SameSiteStrictMode,
	})
}

// Logout はユーザーをログアウトし、認証Cookieをクリアする
// @Summary ログアウト
// @Description ユーザーをログアウトし、認証Cookieをクリアします
// @Tags auth
// @Success 200 {object} map[string]string
// @Router /auth/logout [post]
func (c *AuthController) Logout(ctx echo.Context) error {
	// アクセストークンCookieをクリア
	ctx.SetCookie(&http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   c.serverConfig.CookieSecure,
		SameSite: http.SameSiteStrictMode,
	})

	// リフレッシュトークンCookieをクリア
	ctx.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   c.serverConfig.CookieSecure,
		SameSite: http.SameSiteStrictMode,
	})

	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "ログアウトしました",
	})
}
