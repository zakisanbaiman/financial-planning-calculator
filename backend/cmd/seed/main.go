package main

import (
	"log"

	"github.com/financial-planning-calculator/backend/config"
	"github.com/financial-planning-calculator/backend/infrastructure/database"
)

func main() {
	// Load database configuration
	dbConfig := config.NewDatabaseConfig()

	// Connect to database
	db, err := config.NewDatabaseConnection(dbConfig)
	if err != nil {
		log.Fatalf("データベース接続に失敗しました: %v", err)
	}
	defer db.Close()

	// Create seeder
	seeder := database.NewSeeder(db)

	// Execute seeding
	if err := seeder.Seed(); err != nil {
		log.Fatalf("シードデータの投入に失敗しました: %v", err)
	}

	log.Println("シードデータの投入が完了しました")
}
