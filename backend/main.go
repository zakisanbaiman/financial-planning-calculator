package main

import (
	"log"

	"github.com/financial-planning-calculator/backend/config"
	"github.com/financial-planning-calculator/backend/infrastructure/web"
	"github.com/labstack/echo/v4"

	_ "github.com/financial-planning-calculator/backend/docs"
)

// @title 財務計画計算機 API
// @version 1.0
// @description 将来の資産形成と老後の財務計画を可視化するアプリケーションのAPI
// @host localhost:8080
// @BasePath /api
func main() {
	// 設定読み込み
	cfg := config.LoadServerConfig()

	// Echo インスタンス作成
	e := echo.New()

	// サーバー設定
	e.HideBanner = true
	e.Debug = cfg.Debug

	// カスタムエラーハンドラー
	e.HTTPErrorHandler = web.CustomHTTPErrorHandler

	// ミドルウェア設定
	web.SetupMiddleware(e, cfg)

	// ルーティング設定
	web.SetupRoutes(e)

	// サーバー起動
	log.Printf("サーバーを開始します: http://localhost:%s", cfg.Port)
	log.Printf("Swagger UI: http://localhost:%s/swagger/index.html", cfg.Port)
	log.Printf("API Base URL: http://localhost:%s/api", cfg.Port)
	log.Printf("Debug モード: %v", cfg.Debug)
	log.Printf("許可されたオリジン: %v", cfg.AllowedOrigins)

	e.Logger.Fatal(e.Start(":" + cfg.Port))
}
