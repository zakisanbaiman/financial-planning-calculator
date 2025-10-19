package main

import (
	"flag"
	"log"

	"github.com/financial-planning-calculator/backend/config"
	"github.com/financial-planning-calculator/backend/infrastructure/database"
)

func main() {
	var command string
	flag.StringVar(&command, "command", "up", "Migration command: up, down, status")
	flag.Parse()

	// Load database configuration
	dbConfig := config.NewDatabaseConfig()

	// Connect to database
	db, err := config.NewDatabaseConnection(dbConfig)
	if err != nil {
		log.Fatalf("データベース接続に失敗しました: %v", err)
	}
	defer db.Close()

	// Create migrator
	migrator := database.NewMigrator(db)

	// Execute command
	switch command {
	case "up":
		if err := migrator.Up(); err != nil {
			log.Fatalf("マイグレーションの実行に失敗しました: %v", err)
		}
	case "down":
		if err := migrator.Down(); err != nil {
			log.Fatalf("マイグレーションのロールバックに失敗しました: %v", err)
		}
	case "status":
		if err := migrator.Status(); err != nil {
			log.Fatalf("マイグレーション状況の取得に失敗しました: %v", err)
		}
	default:
		log.Fatalf("無効なコマンドです: %s (使用可能: up, down, status)", command)
	}
}
