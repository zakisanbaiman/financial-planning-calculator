package web

import (
	"github.com/financial-planning-calculator/backend/config"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

// GitHubOAuthMiddleware はGitHub OAuth設定をコンテキストに注入するミドルウェア
func GitHubOAuthMiddleware(cfg *config.ServerConfig) echo.MiddlewareFunc {
	githubOAuthConfig := &oauth2.Config{
		ClientID:     cfg.GitHubClientID,
		ClientSecret: cfg.GitHubClientSecret,
		RedirectURL:  cfg.GitHubCallbackURL,
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("github_oauth_config", githubOAuthConfig)
			c.Set("oauth_success_redirect", cfg.OAuthSuccessRedirect)
			c.Set("oauth_failure_redirect", cfg.OAuthFailureRedirect)
			return next(c)
		}
	}
}
