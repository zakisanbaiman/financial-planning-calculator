package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/financial-planning-calculator/backend/docs"
)

// @title 財務計画計算機 API
// @version 1.0
// @description 将来の資産形成と老後の財務計画を可視化するアプリケーションのAPI
// @host localhost:8080
// @BasePath /api
func main() {
	e := echo.New()

	// ミドルウェア設定
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Swagger UI
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// ヘルスチェック
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status":  "ok",
			"message": "財務計画計算機 API サーバーが正常に動作しています",
		})
	})

	// API ルートグループ
	api := e.Group("/api")

	// 基本的なエンドポイント（後で実装）
	api.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "財務計画計算機 API v1.0",
			"docs":    "/swagger/index.html",
		})
	})

	log.Println("サーバーを開始します: http://localhost:8080")
	log.Println("Swagger UI: http://localhost:8080/swagger/index.html")

	e.Logger.Fatal(e.Start(":8080"))
}
