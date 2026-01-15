package controllers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/financial-planning-calculator/backend/application/usecases"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

// GitHubLogin はGitHubでのログイン開始（OAuth認証画面へリダイレクト）
// @Summary GitHub OAuth認証開始
// @Description GitHubの認証画面にリダイレクトします（Issue: #67）
// @Tags auth
// @Success 302 "GitHubの認証画面へリダイレクト"
// @Router /auth/github [get]
func (c *AuthController) GitHubLogin(ctx echo.Context) error {
	// OAuthコンフィグを取得
	oauthConfig := ctx.Get("github_oauth_config")
	if oauthConfig == nil {
		return ctx.JSON(http.StatusInternalServerError, NewErrorResponse(ctx, ErrorCodeInternalServer, "OAuth設定が見つかりません", nil))
	}

	// ステートパラメータ（CSRF対策）を生成
	state := generateRandomState()

	// セッションにステートを保存（本番ではセッションストアを使用）
	// 簡易実装のため、クッキーに保存
	ctx.SetCookie(&http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		MaxAge:   300, // 5分間有効
		HttpOnly: true,
		Secure:   false, // 開発環境ではfalse、本番ではtrue
		SameSite: http.SameSiteLaxMode,
	})

	// GitHubの認証画面にリダイレクト
	authURL := oauthConfig.(*oauth2.Config).AuthCodeURL(state, oauth2.AccessTypeOffline)
	return ctx.Redirect(http.StatusTemporaryRedirect, authURL)
}

// GitHubCallback はGitHub認証後のコールバック処理
// @Summary GitHub OAuth コールバック
// @Description GitHub認証後のコールバックを処理し、ユーザー情報を取得してログイン/登録します（Issue: #67）
// @Tags auth
// @Param code query string true "GitHub認証コード"
// @Param state query string true "CSRFトークン"
// @Success 302 "ダッシュボードへリダイレクト"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/github/callback [get]
func (c *AuthController) GitHubCallback(ctx echo.Context) error {
	// ステート検証（CSRF対策）
	stateCookie, err := ctx.Cookie("oauth_state")
	if err != nil || stateCookie.Value != ctx.QueryParam("state") {
		return ctx.Redirect(http.StatusTemporaryRedirect, getOAuthFailureRedirect(ctx)+"?error=invalid_state")
	}

	// OAuth認証コードを取得
	code := ctx.QueryParam("code")
	if code == "" {
		return ctx.Redirect(http.StatusTemporaryRedirect, getOAuthFailureRedirect(ctx)+"?error=no_code")
	}

	// OAuthコンフィグを取得
	oauthConfig := ctx.Get("github_oauth_config").(*oauth2.Config)

	// 認証コードをトークンに交換
	token, err := oauthConfig.Exchange(ctx.Request().Context(), code)
	if err != nil {
		return ctx.Redirect(http.StatusTemporaryRedirect, getOAuthFailureRedirect(ctx)+"?error=token_exchange_failed")
	}

	// GitHubからユーザー情報を取得
	client := oauthConfig.Client(ctx.Request().Context(), token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return ctx.Redirect(http.StatusTemporaryRedirect, getOAuthFailureRedirect(ctx)+"?error=user_fetch_failed")
	}
	defer resp.Body.Close()

	var githubUser struct {
		ID        int64  `json:"id"`
		Login     string `json:"login"`
		Email     string `json:"email"`
		Name      string `json:"name"`
		AvatarURL string `json:"avatar_url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&githubUser); err != nil {
		return ctx.Redirect(http.StatusTemporaryRedirect, getOAuthFailureRedirect(ctx)+"?error=user_parse_failed")
	}

	// メールアドレスが空の場合、GitHub APIから取得
	if githubUser.Email == "" {
		emailResp, err := client.Get("https://api.github.com/user/emails")
		if err == nil {
			defer emailResp.Body.Close()
			var emails []struct {
				Email    string `json:"email"`
				Primary  bool   `json:"primary"`
				Verified bool   `json:"verified"`
			}
			if json.NewDecoder(emailResp.Body).Decode(&emails) == nil {
				for _, email := range emails {
					if email.Primary && email.Verified {
						githubUser.Email = email.Email
						break
					}
				}
			}
		}
	}

	// メールアドレスが取得できない場合はエラー
	if githubUser.Email == "" {
		return ctx.Redirect(http.StatusTemporaryRedirect, getOAuthFailureRedirect(ctx)+"?error=no_email")
	}

	// ユーザー名が空の場合はログインIDを使用
	if githubUser.Name == "" {
		githubUser.Name = githubUser.Login
	}

	// ユースケースを使用してログイン/登録
	output, err := c.authUseCase.GitHubOAuthLogin(ctx.Request().Context(), usecases.GitHubOAuthInput{
		GitHubUserID: fmt.Sprintf("%d", githubUser.ID),
		Email:        githubUser.Email,
		Name:         githubUser.Name,
		AvatarURL:    githubUser.AvatarURL,
	})
	if err != nil {
		return ctx.Redirect(http.StatusTemporaryRedirect, getOAuthFailureRedirect(ctx)+"?error=login_failed")
	}

	// 成功時、トークンをクッキーに保存してフロントエンドにリダイレクト
	// 本来はフロントエンドが適切にトークンを受け取る仕組みが必要
	successURL := getOAuthSuccessRedirect(ctx) +
		"?token=" + output.Token +
		"&refresh_token=" + output.RefreshToken +
		"&user_id=" + output.UserID +
		"&email=" + output.Email

	return ctx.Redirect(http.StatusTemporaryRedirect, successURL)
}

// Helper functions

func generateRandomState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func getOAuthSuccessRedirect(ctx echo.Context) string {
	if url := ctx.Get("oauth_success_redirect"); url != nil {
		return url.(string)
	}
	return "http://localhost:3000/auth/callback"
}

func getOAuthFailureRedirect(ctx echo.Context) string {
	if url := ctx.Get("oauth_failure_redirect"); url != nil {
		return url.(string)
	}
	return "http://localhost:3000/login"
}
