package controllers

import (
	"net/http"

	"github.com/financial-planning-calculator/backend/application/usecases"
	"github.com/financial-planning-calculator/backend/domain/entities"
	"github.com/labstack/echo/v4"
)

// TwoFactorController は2段階認証関連のコントローラー
type TwoFactorController struct {
	authUseCase usecases.AuthUseCase
}

// NewTwoFactorController は新しいTwoFactorControllerを作成する
func NewTwoFactorController(authUseCase usecases.AuthUseCase) *TwoFactorController {
	return &TwoFactorController{
		authUseCase: authUseCase,
	}
}

// Setup2FAResponse は2FA設定開始のレスポンス
type Setup2FAResponse struct {
	Secret      string   `json:"secret"`
	QRCodeURL   string   `json:"qr_code_url"`
	BackupCodes []string `json:"backup_codes"`
}

// Enable2FARequest は2FA有効化のリクエスト
type Enable2FARequest struct {
	Code   string `json:"code" validate:"required,len=6"`
	Secret string `json:"secret" validate:"required"`
}

// Verify2FARequest は2FA検証のリクエスト
type Verify2FARequest struct {
	Code      string `json:"code" validate:"required"`
	UseBackup bool   `json:"use_backup"`
}

// Disable2FARequest は2FA無効化のリクエスト
type Disable2FARequest struct {
	Password string `json:"password" validate:"required"`
}

// RegenerateBackupCodesResponse はバックアップコード再生成のレスポンス
type RegenerateBackupCodesResponse struct {
	BackupCodes []string `json:"backup_codes"`
}

// getUserIDFromContext はコンテキストからユーザーIDを取得する
func getUserIDFromContext(ctx echo.Context) (string, error) {
	userID, ok := ctx.Get("user_id").(string)
	if !ok || userID == "" {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "ユーザー情報が取得できません")
	}

	// バリデーション
	_, err := entities.NewUserID(userID)
	if err != nil {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "無効なユーザーIDです")
	}

	return userID, nil
}

// Setup2FA は2段階認証のセットアップを開始する
// @Summary 2FA設定開始
// @Description 2段階認証のセットアップを開始し、QRコードとバックアップコードを生成します
// @Tags two-factor-auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} Setup2FAResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse "2段階認証は既に有効です"
// @Failure 500 {object} ErrorResponse
// @Router /auth/2fa/setup [post]
func (c *TwoFactorController) Setup2FA(ctx echo.Context) error {
	// JWTトークンからユーザーIDを取得
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, NewErrorResponse(ctx, ErrorCodeUnauthorized, "認証が必要です", err.Error()))
	}

	// 2FA設定を開始
	output, err := c.authUseCase.Setup2FA(ctx.Request().Context(), userID)
	if err != nil {
		if err.Error() == "2段階認証は既に有効です" {
			return ctx.JSON(http.StatusConflict, NewErrorResponse(ctx, ErrorCodeConflict, err.Error(), nil))
		}
		return ctx.JSON(http.StatusInternalServerError, NewErrorResponse(ctx, ErrorCodeInternalServer, "2FA設定の開始に失敗しました", err.Error()))
	}

	response := Setup2FAResponse{
		Secret:      output.Secret,
		QRCodeURL:   output.QRCodeURL,
		BackupCodes: output.BackupCodes,
	}

	return ctx.JSON(http.StatusOK, response)
}

// Enable2FA は2段階認証を有効化する
// @Summary 2FA有効化
// @Description 初回コード検証を行い、2段階認証を有効化します
// @Tags two-factor-auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body Enable2FARequest true "2FA有効化リクエスト"
// @Success 200 {object} map[string]string
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse "2段階認証は既に有効です"
// @Failure 500 {object} ErrorResponse
// @Router /auth/2fa/enable [post]
func (c *TwoFactorController) Enable2FA(ctx echo.Context) error {
	// JWTトークンからユーザーIDを取得
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, NewErrorResponse(ctx, ErrorCodeUnauthorized, "認証が必要です", err.Error()))
	}

	var req Enable2FARequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, NewErrorResponse(ctx, ErrorCodeBadRequest, "リクエストの解析に失敗しました", err.Error()))
	}

	if err := ctx.Validate(&req); err != nil {
		return err
	}

	// 2FAを有効化
	input := usecases.Enable2FAInput{
		UserID: userID,
		Code:   req.Code,
		Secret: req.Secret,
	}

	if err := c.authUseCase.Enable2FA(ctx.Request().Context(), input); err != nil {
		if err.Error() == "2段階認証は既に有効です" {
			return ctx.JSON(http.StatusConflict, NewErrorResponse(ctx, ErrorCodeConflict, err.Error(), nil))
		}
		if err.Error() == "認証コードが無効です" {
			return ctx.JSON(http.StatusBadRequest, NewErrorResponse(ctx, ErrorCodeValidation, err.Error(), nil))
		}
		return ctx.JSON(http.StatusInternalServerError, NewErrorResponse(ctx, ErrorCodeInternalServer, "2FAの有効化に失敗しました", err.Error()))
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "2段階認証が有効になりました",
	})
}

// Verify2FA はログイン時の2FAコード検証を行う
// @Summary 2FA検証
// @Description ログイン後の2段階認証コードを検証し、本トークンを発行します
// @Tags two-factor-auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body Verify2FARequest true "2FA検証リクエスト"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/2fa/verify [post]
func (c *TwoFactorController) Verify2FA(ctx echo.Context) error {
	// JWTトークンからユーザーIDを取得（仮トークン）
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, NewErrorResponse(ctx, ErrorCodeUnauthorized, "認証が必要です", err.Error()))
	}

	var req Verify2FARequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, NewErrorResponse(ctx, ErrorCodeBadRequest, "リクエストの解析に失敗しました", err.Error()))
	}

	if err := ctx.Validate(&req); err != nil {
		return err
	}

	// 2FAコードを検証
	input := usecases.Verify2FAInput{
		UserID:    userID,
		Code:      req.Code,
		UseBackup: req.UseBackup,
	}

	output, err := c.authUseCase.Verify2FA(ctx.Request().Context(), input)
	if err != nil {
		if err.Error() == "認証コードが無効です" {
			return ctx.JSON(http.StatusBadRequest, NewErrorResponse(ctx, ErrorCodeValidation, err.Error(), nil))
		}
		return ctx.JSON(http.StatusInternalServerError, NewErrorResponse(ctx, ErrorCodeInternalServer, "2FA検証に失敗しました", err.Error()))
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

// Disable2FA は2段階認証を無効化する
// @Summary 2FA無効化
// @Description 2段階認証を無効化します（パスワード検証が必要）
// @Tags two-factor-auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body Disable2FARequest true "2FA無効化リクエスト"
// @Success 200 {object} map[string]string
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/2fa [delete]
func (c *TwoFactorController) Disable2FA(ctx echo.Context) error {
	// JWTトークンからユーザーIDを取得
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, NewErrorResponse(ctx, ErrorCodeUnauthorized, "認証が必要です", err.Error()))
	}

	var req Disable2FARequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, NewErrorResponse(ctx, ErrorCodeBadRequest, "リクエストの解析に失敗しました", err.Error()))
	}

	if err := ctx.Validate(&req); err != nil {
		return err
	}

	// 2FAを無効化
	input := usecases.Disable2FAInput{
		UserID:   userID,
		Password: req.Password,
	}

	if err := c.authUseCase.Disable2FA(ctx.Request().Context(), input); err != nil {
		if err.Error() == "パスワードが正しくありません" {
			return ctx.JSON(http.StatusBadRequest, NewErrorResponse(ctx, ErrorCodeValidation, err.Error(), nil))
		}
		if err.Error() == "2段階認証は有効になっていません" {
			return ctx.JSON(http.StatusBadRequest, NewErrorResponse(ctx, ErrorCodeValidation, err.Error(), nil))
		}
		return ctx.JSON(http.StatusInternalServerError, NewErrorResponse(ctx, ErrorCodeInternalServer, "2FAの無効化に失敗しました", err.Error()))
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "2段階認証が無効になりました",
	})
}

// RegenerateBackupCodes はバックアップコードを再生成する
// @Summary バックアップコード再生成
// @Description 2段階認証のバックアップコードを再生成します
// @Tags two-factor-auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} RegenerateBackupCodesResponse
// @Failure 401 {object} ErrorResponse
// @Failure 400 {object} ErrorResponse "2段階認証が有効になっていません"
// @Failure 500 {object} ErrorResponse
// @Router /auth/2fa/backup-codes [post]
func (c *TwoFactorController) RegenerateBackupCodes(ctx echo.Context) error {
	// JWTトークンからユーザーIDを取得
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, NewErrorResponse(ctx, ErrorCodeUnauthorized, "認証が必要です", err.Error()))
	}

	// バックアップコードを再生成
	output, err := c.authUseCase.RegenerateBackupCodes(ctx.Request().Context(), userID)
	if err != nil {
		if err.Error() == "2段階認証が有効になっていません" {
			return ctx.JSON(http.StatusBadRequest, NewErrorResponse(ctx, ErrorCodeValidation, err.Error(), nil))
		}
		return ctx.JSON(http.StatusInternalServerError, NewErrorResponse(ctx, ErrorCodeInternalServer, "バックアップコードの再生成に失敗しました", err.Error()))
	}

	response := RegenerateBackupCodesResponse{
		BackupCodes: output.BackupCodes,
	}

	return ctx.JSON(http.StatusOK, response)
}

// Get2FAStatusResponse は2FAステータス取得のレスポンス
type Get2FAStatusResponse struct {
	Enabled bool `json:"enabled"`
}

// Get2FAStatus は2FAの有効状態を取得する
// @Summary 2FAステータス取得
// @Description 2段階認証の有効状態を取得します
// @Tags two-factor-auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} Get2FAStatusResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/2fa/status [get]
func (c *TwoFactorController) Get2FAStatus(ctx echo.Context) error {
	// JWTトークンからユーザーIDを取得
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, NewErrorResponse(ctx, ErrorCodeUnauthorized, "認証が必要です", err.Error()))
	}

	// 2FAステータスを取得
	output, err := c.authUseCase.Get2FAStatus(ctx.Request().Context(), userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, NewErrorResponse(ctx, ErrorCodeInternalServer, "2FAステータスの取得に失敗しました", err.Error()))
	}

	response := Get2FAStatusResponse{
		Enabled: output.Enabled,
	}

	return ctx.JSON(http.StatusOK, response)
}
