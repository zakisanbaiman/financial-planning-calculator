package web

import (
	"net/http"
	"strings"

	"github.com/financial-planning-calculator/backend/application/usecases"
	"github.com/financial-planning-calculator/backend/domain/entities"
	"github.com/labstack/echo/v4"
)

// JWTAuthMiddleware はJWT認証ミドルウェア
// Cookieからトークンを取得し、なければAuthorizationヘッダーから取得（後方互換性のため）
// 2FA仮トークン（Requires2FA: true）は2FA検証エンドポイントのみで許可される
func JWTAuthMiddleware(authUseCase usecases.AuthUseCase) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var tokenString string

			// まずCookieからトークンを取得
			cookie, err := c.Cookie("access_token")
			if err == nil && cookie.Value != "" {
				tokenString = cookie.Value
			} else {
				// Cookieにトークンがない場合、Authorizationヘッダーから取得（後方互換性のため）
				authHeader := c.Request().Header.Get("Authorization")
				if authHeader == "" {
					return echo.NewHTTPError(http.StatusUnauthorized, "認証トークンが必要です")
				}

				// "Bearer "プレフィックスを確認
				const bearerPrefix = "Bearer "
				if !strings.HasPrefix(authHeader, bearerPrefix) {
					return echo.NewHTTPError(http.StatusUnauthorized, "無効な認証形式です")
				}

				// トークンを抽出
				tokenString = strings.TrimPrefix(authHeader, bearerPrefix)
				if tokenString == "" {
					return echo.NewHTTPError(http.StatusUnauthorized, "認証トークンが必要です")
				}
			}

			// トークンを検証
			claims, err := authUseCase.VerifyToken(c.Request().Context(), tokenString)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "無効または期限切れの認証トークンです")
			}

			// 2FA仮トークンの場合、2FA検証エンドポイントのみ許可
			if claims.Requires2FA || claims.TwoFactorVerify {
				path := c.Request().URL.Path
				// /api/auth/2fa/verify のみ許可
				if !strings.HasSuffix(path, "/auth/2fa/verify") {
					return echo.NewHTTPError(http.StatusUnauthorized, "2段階認証の検証が必要です")
				}
			}

			// ユーザー情報をコンテキストに保存
			c.Set("user_id", claims.UserID)
			c.Set("email", claims.Email)

			return next(c)
		}
	}
}

// GetUserIDFromContext はコンテキストからユーザーIDを取得する
func GetUserIDFromContext(c echo.Context) (entities.UserID, error) {
	userID, ok := c.Get("user_id").(string)
	if !ok || userID == "" {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "ユーザー情報が取得できません")
	}

	return entities.NewUserID(userID)
}

// GetEmailFromContext はコンテキストからメールアドレスを取得する
func GetEmailFromContext(c echo.Context) string {
	email, _ := c.Get("email").(string)
	return email
}
