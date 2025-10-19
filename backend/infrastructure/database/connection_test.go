package database

import (
	"testing"

	"github.com/financial-planning-calculator/backend/config"
)

func TestDatabaseConnection(t *testing.T) {
	// Skip if no database is available
	if testing.Short() {
		t.Skip("Skipping database connection test in short mode")
	}

	// Load test database configuration
	dbConfig := config.NewDatabaseConfig()

	// Try to connect to database
	db, err := config.NewDatabaseConnection(dbConfig)
	if err != nil {
		t.Skipf("Database connection failed (expected in CI): %v", err)
		return
	}
	defer db.Close()

	// Test basic query
	var result int
	err = db.QueryRow("SELECT 1").Scan(&result)
	if err != nil {
		t.Fatalf("Basic query failed: %v", err)
	}

	if result != 1 {
		t.Fatalf("Expected 1, got %d", result)
	}
}

func TestMigratorCreation(t *testing.T) {
	// This test doesn't require a real database connection
	// We're just testing that the migrator can be created
	migrator := NewMigrator(nil)
	if migrator == nil {
		t.Fatal("Failed to create migrator")
	}
}

func TestSeederCreation(t *testing.T) {
	// This test doesn't require a real database connection
	// We're just testing that the seeder can be created
	seeder := NewSeeder(nil)
	if seeder == nil {
		t.Fatal("Failed to create seeder")
	}
}
