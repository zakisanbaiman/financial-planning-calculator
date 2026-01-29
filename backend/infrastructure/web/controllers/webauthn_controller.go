package controllers

import (
	"net/http"

	"github.com/financial-planning-calculator/backend/application/usecases"
	"github.com/labstack/echo/v4"
)

// WebAuthnController はパスキー（WebAuthn）認証関連のコントローラー
type WebAuthnController struct {
	webAuthnUseCase usecases.WebAuthnUseCase
}

// NewWebAuthnController は新しいWebAuthnControllerを作成する
func NewWebAuthnController(webAuthnUseCase usecases.WebAuthnUseCase) *WebAuthnController {
	return &WebAuthnController{
		webAuthnUseCase: webAuthnUseCase,
	}
}

// BeginRegistrationResponse はパスキー登録開始のレスポンス
type BeginRegistrationResponse struct {
	PublicKey   string `json:"publicKey"`
	SessionData string `json:"sessionData"`
}

// FinishRegistrationRequest はパスキー登録完了のリクエスト
type FinishRegistrationRequest struct {
	CredentialName string `json:"credential_name" validate:"required"`
	Response       string `json:"response" validate:"required"`
	SessionData    string `json:"session_data" validate:"required"`
}

// BeginLoginRequest はパスキーログイン開始のリクエスト
type BeginLoginRequest struct {
	Email string `json:"email"`
}

// BeginLoginResponse はパスキーログイン開始のレスポンス
type BeginLoginResponse struct {
	PublicKey   string `json:"publicKey"`
	SessionData string `json:"sessionData"`
}

// FinishLoginRequest はパスキーログイン完了のリクエスト
type FinishLoginRequest struct {
	Response    string `json:"response" validate:"required"`
	SessionData string `json:"session_data" validate:"required"`
}

// CredentialInfoResponse はパスキー情報のレスポンス
type CredentialInfoResponse struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	CreatedAt  string  `json:"created_at"`
	LastUsedAt *string `json:"last_used_at,omitempty"`
}

// RenameCredentialRequest はパスキー名変更のリクエスト
type RenameCredentialRequest struct {
	Name string `json:"name" validate:"required"`
}

// BeginRegistration はパスキー登録を開始する
// @Summary パスキー登録開始
// @Description パスキー（WebAuthn）の登録を開始します。認証済みユーザーのみ利用可能です。
// @Tags passkey
// @Security BearerAuth
// @Produce json
// @Success 200 {object} BeginRegistrationResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure 503 {object} ErrorResponse "パスキー機能は利用できません"
// @Router /auth/passkey/register/begin [post]
func (c *WebAuthnController) BeginRegistration(ctx echo.Context) error {
	// WebAuthn機能の利用可否をチェック
	if c.webAuthnUseCase == nil {
		return ctx.JSON(http.StatusServiceUnavailable, NewErrorResponse(ctx, ErrorCodeInternalServer, "パスキー機能は現在利用できません", nil))
	}

	// ユーザーIDを取得
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, NewErrorResponse(ctx, ErrorCodeUnauthorized, "認証が必要です", err.Error()))
	}

	// パスキー登録を開始
	output, err := c.webAuthnUseCase.BeginRegistration(ctx.Request().Context(), userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, NewErrorResponse(ctx, ErrorCodeInternalServer, "パスキー登録の開始に失敗しました", err.Error()))
	}

	response := BeginRegistrationResponse{
		PublicKey:   output.PublicKeyOptions,
		SessionData: output.SessionData,
	}

	return ctx.JSON(http.StatusOK, response)
}

// FinishRegistration はパスキー登録を完了する
// @Summary パスキー登録完了
// @Description パスキー（WebAuthn）の登録を完了します
// @Tags passkey
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body FinishRegistrationRequest true "パスキー登録完了リクエスト"
// @Success 200 {object} map[string]string
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/passkey/register/finish [post]
func (c *WebAuthnController) FinishRegistration(ctx echo.Context) error {
	// WebAuthn機能の利用可否をチェック
	if c.webAuthnUseCase == nil {
		return ctx.JSON(http.StatusServiceUnavailable, NewErrorResponse(ctx, ErrorCodeInternalServer, "パスキー機能は現在利用できません", nil))
	}

	// ユーザーIDを取得
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, NewErrorResponse(ctx, ErrorCodeUnauthorized, "認証が必要です", err.Error()))
	}

	var req FinishRegistrationRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, NewErrorResponse(ctx, ErrorCodeBadRequest, "リクエストの解析に失敗しました", err.Error()))
	}

	if err := ctx.Validate(&req); err != nil {
		return err // Validator already returns proper error response
	}

	// パスキー登録を完了
	input := usecases.FinishRegistrationInput{
		UserID:         userID,
		CredentialName: req.CredentialName,
		Response:       req.Response,
		SessionData:    req.SessionData,
	}

	if err := c.webAuthnUseCase.FinishRegistration(ctx.Request().Context(), input); err != nil {
		return ctx.JSON(http.StatusInternalServerError, NewErrorResponse(ctx, ErrorCodeInternalServer, "パスキー登録の完了に失敗しました", err.Error()))
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "パスキーを登録しました",
	})
}

// BeginLogin はパスキーログインを開始する
// @Summary パスキーログイン開始
// @Description パスキー（WebAuthn）でのログインを開始します
// @Tags passkey
// @Accept json
// @Produce json
// @Param request body BeginLoginRequest false "パスキーログイン開始リクエスト"
// @Success 200 {object} BeginLoginResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/passkey/login/begin [post]
func (c *WebAuthnController) BeginLogin(ctx echo.Context) error {
	// WebAuthn機能の利用可否をチェック
	if c.webAuthnUseCase == nil {
		return ctx.JSON(http.StatusServiceUnavailable, NewErrorResponse(ctx, ErrorCodeInternalServer, "パスキー機能は現在利用できません", nil))
	}

	var req BeginLoginRequest
	if err := ctx.Bind(&req); err != nil {
		// リクエストボディがない場合もOK（ユーザーレス認証）
		req = BeginLoginRequest{}
	}

	// パスキーログインを開始
	input := usecases.BeginLoginInput{
		Email: req.Email,
	}

	output, err := c.webAuthnUseCase.BeginLogin(ctx.Request().Context(), input)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, NewErrorResponse(ctx, ErrorCodeInternalServer, "パスキーログインの開始に失敗しました", err.Error()))
	}

	response := BeginLoginResponse{
		PublicKey:   output.PublicKeyOptions,
		SessionData: output.SessionData,
	}

	return ctx.JSON(http.StatusOK, response)
}

// FinishLogin はパスキーログインを完了する
// @Summary パスキーログイン完了
// @Description パスキー（WebAuthn）でのログインを完了し、JWTトークンを発行します
// @Tags passkey
// @Accept json
// @Produce json
// @Param request body FinishLoginRequest true "パスキーログイン完了リクエスト"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/passkey/login/finish [post]
func (c *WebAuthnController) FinishLogin(ctx echo.Context) error {
	// WebAuthn機能の利用可否をチェック
	if c.webAuthnUseCase == nil {
		return ctx.JSON(http.StatusServiceUnavailable, NewErrorResponse(ctx, ErrorCodeInternalServer, "パスキー機能は現在利用できません", nil))
	}

	var req FinishLoginRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, NewErrorResponse(ctx, ErrorCodeBadRequest, "リクエストの解析に失敗しました", err.Error()))
	}

	if err := ctx.Validate(&req); err != nil {
		return err // Validator already returns proper error response
	}

	// パスキーログインを完了
	input := usecases.FinishLoginInput{
		Response:    req.Response,
		SessionData: req.SessionData,
	}

	output, err := c.webAuthnUseCase.FinishLogin(ctx.Request().Context(), input)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, NewErrorResponse(ctx, ErrorCodeUnauthorized, "パスキー認証に失敗しました", err.Error()))
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

// ListCredentials はユーザーの全パスキーを取得する
// @Summary パスキー一覧取得
// @Description 認証済みユーザーの登録済みパスキー一覧を取得します
// @Tags passkey
// @Security BearerAuth
// @Produce json
// @Success 200 {array} CredentialInfoResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/passkey/credentials [get]
func (c *WebAuthnController) ListCredentials(ctx echo.Context) error {
	// WebAuthn機能の利用可否をチェック
	if c.webAuthnUseCase == nil {
		return ctx.JSON(http.StatusServiceUnavailable, NewErrorResponse(ctx, ErrorCodeInternalServer, "パスキー機能は現在利用できません", nil))
	}

	// ユーザーIDを取得
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, NewErrorResponse(ctx, ErrorCodeUnauthorized, "認証が必要です", err.Error()))
	}

	// パスキー一覧を取得
	credentials, err := c.webAuthnUseCase.ListCredentials(ctx.Request().Context(), userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, NewErrorResponse(ctx, ErrorCodeInternalServer, "パスキー一覧の取得に失敗しました", err.Error()))
	}

	response := make([]CredentialInfoResponse, 0, len(credentials))
	for _, cred := range credentials {
		response = append(response, CredentialInfoResponse{
			ID:         cred.ID,
			Name:       cred.Name,
			CreatedAt:  cred.CreatedAt,
			LastUsedAt: cred.LastUsedAt,
		})
	}

	return ctx.JSON(http.StatusOK, response)
}

// DeleteCredential はパスキーを削除する
// @Summary パスキー削除
// @Description 認証済みユーザーの指定されたパスキーを削除します
// @Tags passkey
// @Security BearerAuth
// @Param credential_id path string true "クレデンシャルID"
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/passkey/credentials/{credential_id} [delete]
func (c *WebAuthnController) DeleteCredential(ctx echo.Context) error {
	// WebAuthn機能の利用可否をチェック
	if c.webAuthnUseCase == nil {
		return ctx.JSON(http.StatusServiceUnavailable, NewErrorResponse(ctx, ErrorCodeInternalServer, "パスキー機能は現在利用できません", nil))
	}

	// ユーザーIDを取得
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, NewErrorResponse(ctx, ErrorCodeUnauthorized, "認証が必要です", err.Error()))
	}

	credentialID := ctx.Param("credential_id")
	if credentialID == "" {
		return ctx.JSON(http.StatusBadRequest, NewErrorResponse(ctx, ErrorCodeBadRequest, "クレデンシャルIDが必要です", nil))
	}

	// パスキーを削除
	if err := c.webAuthnUseCase.DeleteCredential(ctx.Request().Context(), userID, credentialID); err != nil {
		return ctx.JSON(http.StatusInternalServerError, NewErrorResponse(ctx, ErrorCodeInternalServer, "パスキーの削除に失敗しました", err.Error()))
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "パスキーを削除しました",
	})
}

// RenameCredential はパスキーの名前を変更する
// @Summary パスキー名変更
// @Description 認証済みユーザーの指定されたパスキーの名前を変更します
// @Tags passkey
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param credential_id path string true "クレデンシャルID"
// @Param request body RenameCredentialRequest true "パスキー名変更リクエスト"
// @Success 200 {object} map[string]string
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/passkey/credentials/{credential_id} [put]
func (c *WebAuthnController) RenameCredential(ctx echo.Context) error {
	// WebAuthn機能の利用可否をチェック
	if c.webAuthnUseCase == nil {
		return ctx.JSON(http.StatusServiceUnavailable, NewErrorResponse(ctx, ErrorCodeInternalServer, "パスキー機能は現在利用できません", nil))
	}

	// ユーザーIDを取得
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, NewErrorResponse(ctx, ErrorCodeUnauthorized, "認証が必要です", err.Error()))
	}

	credentialID := ctx.Param("credential_id")
	if credentialID == "" {
		return ctx.JSON(http.StatusBadRequest, NewErrorResponse(ctx, ErrorCodeBadRequest, "クレデンシャルIDが必要です", nil))
	}

	var req RenameCredentialRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, NewErrorResponse(ctx, ErrorCodeBadRequest, "リクエストの解析に失敗しました", err.Error()))
	}

	if err := ctx.Validate(&req); err != nil {
		return err // Validator already returns proper error response
	}

	// パスキーの名前を変更
	if err := c.webAuthnUseCase.RenameCredential(ctx.Request().Context(), userID, credentialID, req.Name); err != nil {
		return ctx.JSON(http.StatusInternalServerError, NewErrorResponse(ctx, ErrorCodeInternalServer, "パスキー名の変更に失敗しました", err.Error()))
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "パスキー名を変更しました",
	})
}
