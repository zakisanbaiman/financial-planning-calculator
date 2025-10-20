package repositories

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/financial-planning-calculator/backend/config"
	"github.com/financial-planning-calculator/backend/domain/entities"
	"github.com/financial-planning-calculator/backend/domain/valueobjects"
	"github.com/financial-planning-calculator/backend/infrastructure/database"
)

func setupTestDB(t *testing.T) *sql.DB {
	// Skip if no database is available
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	dbConfig := config.NewDatabaseConfig()
	db, err := config.NewDatabaseConnection(dbConfig)
	if err != nil {
		t.Skipf("Database connection failed (expected in CI): %v", err)
		return nil
	}

	// Run migrations
	migrator := database.NewMigrator(db)
	if err := migrator.Up(); err != nil {
		t.Fatalf("Migration failed: %v", err)
	}

	return db
}

func createTestUser(t *testing.T, db *sql.DB) entities.UserID {
	// Use uuid_generate_v4() to create a proper UUID
	var userID string
	err := db.QueryRow("SELECT uuid_generate_v4()::text").Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to generate UUID: %v", err)
	}

	// Insert test user
	_, err = db.Exec("INSERT INTO users (id, email) VALUES ($1, $2)", userID, userID+"@test.com")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	return entities.UserID(userID)
}

func TestPostgreSQLGoalRepository_Save(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	userID := createTestUser(t, db)
	repo := NewPostgreSQLGoalRepository(db)

	// Create test goal
	targetAmount, err := valueobjects.NewMoneyJPY(1000000)
	if err != nil {
		t.Fatalf("Failed to create target amount: %v", err)
	}

	monthlyContribution, err := valueobjects.NewMoneyJPY(50000)
	if err != nil {
		t.Fatalf("Failed to create monthly contribution: %v", err)
	}

	targetDate := time.Now().AddDate(1, 0, 0) // 1 year from now

	goal, err := entities.NewGoal(
		userID,
		entities.GoalTypeSavings,
		"Test Savings Goal",
		targetAmount,
		targetDate,
		monthlyContribution,
	)
	if err != nil {
		t.Fatalf("Failed to create goal: %v", err)
	}

	// Test Save
	ctx := context.Background()
	err = repo.Save(ctx, goal)
	if err != nil {
		t.Fatalf("Failed to save goal: %v", err)
	}

	// Verify goal was saved
	exists, err := repo.Exists(ctx, goal.ID())
	if err != nil {
		t.Fatalf("Failed to check goal existence: %v", err)
	}
	if !exists {
		t.Error("Goal should exist after saving")
	}
}

func TestPostgreSQLGoalRepository_FindByID(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	userID := createTestUser(t, db)
	repo := NewPostgreSQLGoalRepository(db)

	// Create and save test goal
	targetAmount, _ := valueobjects.NewMoneyJPY(1000000)
	monthlyContribution, _ := valueobjects.NewMoneyJPY(50000)
	targetDate := time.Now().AddDate(1, 0, 0)

	originalGoal, err := entities.NewGoal(
		userID,
		entities.GoalTypeSavings,
		"Test Savings Goal",
		targetAmount,
		targetDate,
		monthlyContribution,
	)
	if err != nil {
		t.Fatalf("Failed to create goal: %v", err)
	}

	ctx := context.Background()
	err = repo.Save(ctx, originalGoal)
	if err != nil {
		t.Fatalf("Failed to save goal: %v", err)
	}

	// Test FindByID
	foundGoal, err := repo.FindByID(ctx, originalGoal.ID())
	if err != nil {
		t.Fatalf("Failed to find goal: %v", err)
	}

	// Verify goal properties
	if foundGoal.ID() != originalGoal.ID() {
		t.Errorf("Expected goal ID %s, got %s", originalGoal.ID(), foundGoal.ID())
	}
	if foundGoal.Title() != originalGoal.Title() {
		t.Errorf("Expected goal title %s, got %s", originalGoal.Title(), foundGoal.Title())
	}
	if foundGoal.GoalType() != originalGoal.GoalType() {
		t.Errorf("Expected goal type %s, got %s", originalGoal.GoalType(), foundGoal.GoalType())
	}
}

func TestPostgreSQLGoalRepository_FindByUserID(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	userID := createTestUser(t, db)
	repo := NewPostgreSQLGoalRepository(db)

	// Create and save multiple test goals
	targetAmount, _ := valueobjects.NewMoneyJPY(1000000)
	monthlyContribution, _ := valueobjects.NewMoneyJPY(50000)
	targetDate := time.Now().AddDate(1, 0, 0)

	goal1, _ := entities.NewGoal(userID, entities.GoalTypeSavings, "Goal 1", targetAmount, targetDate, monthlyContribution)
	goal2, _ := entities.NewGoal(userID, entities.GoalTypeEmergency, "Goal 2", targetAmount, targetDate, monthlyContribution)

	ctx := context.Background()
	repo.Save(ctx, goal1)
	repo.Save(ctx, goal2)

	// Test FindByUserID
	goals, err := repo.FindByUserID(ctx, userID)
	if err != nil {
		t.Fatalf("Failed to find goals by user ID: %v", err)
	}

	if len(goals) != 2 {
		t.Errorf("Expected 2 goals, got %d", len(goals))
	}
}

func TestPostgreSQLGoalRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	userID := createTestUser(t, db)
	repo := NewPostgreSQLGoalRepository(db)

	// Create and save test goal
	targetAmount, _ := valueobjects.NewMoneyJPY(1000000)
	monthlyContribution, _ := valueobjects.NewMoneyJPY(50000)
	targetDate := time.Now().AddDate(1, 0, 0)

	goal, err := entities.NewGoal(
		userID,
		entities.GoalTypeSavings,
		"Original Title",
		targetAmount,
		targetDate,
		monthlyContribution,
	)
	if err != nil {
		t.Fatalf("Failed to create goal: %v", err)
	}

	ctx := context.Background()
	err = repo.Save(ctx, goal)
	if err != nil {
		t.Fatalf("Failed to save goal: %v", err)
	}

	// Update goal
	err = goal.UpdateTitle("Updated Title")
	if err != nil {
		t.Fatalf("Failed to update goal title: %v", err)
	}

	// Test Update
	err = repo.Update(ctx, goal)
	if err != nil {
		t.Fatalf("Failed to update goal: %v", err)
	}

	// Verify update
	updatedGoal, err := repo.FindByID(ctx, goal.ID())
	if err != nil {
		t.Fatalf("Failed to find updated goal: %v", err)
	}

	if updatedGoal.Title() != "Updated Title" {
		t.Errorf("Expected updated title 'Updated Title', got '%s'", updatedGoal.Title())
	}
}

func TestPostgreSQLGoalRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	userID := createTestUser(t, db)
	repo := NewPostgreSQLGoalRepository(db)

	// Create and save test goal
	targetAmount, _ := valueobjects.NewMoneyJPY(1000000)
	monthlyContribution, _ := valueobjects.NewMoneyJPY(50000)
	targetDate := time.Now().AddDate(1, 0, 0)

	goal, err := entities.NewGoal(
		userID,
		entities.GoalTypeSavings,
		"Test Goal",
		targetAmount,
		targetDate,
		monthlyContribution,
	)
	if err != nil {
		t.Fatalf("Failed to create goal: %v", err)
	}

	ctx := context.Background()
	err = repo.Save(ctx, goal)
	if err != nil {
		t.Fatalf("Failed to save goal: %v", err)
	}

	// Test Delete
	err = repo.Delete(ctx, goal.ID())
	if err != nil {
		t.Fatalf("Failed to delete goal: %v", err)
	}

	// Verify deletion
	exists, err := repo.Exists(ctx, goal.ID())
	if err != nil {
		t.Fatalf("Failed to check goal existence: %v", err)
	}
	if exists {
		t.Error("Goal should not exist after deletion")
	}
}

func TestPostgreSQLGoalRepository_CountActiveGoalsByType(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	userID := createTestUser(t, db)
	repo := NewPostgreSQLGoalRepository(db)

	// Create and save test goals
	targetAmount, _ := valueobjects.NewMoneyJPY(1000000)
	monthlyContribution, _ := valueobjects.NewMoneyJPY(50000)
	targetDate := time.Now().AddDate(1, 0, 0)

	goal1, _ := entities.NewGoal(userID, entities.GoalTypeSavings, "Savings Goal 1", targetAmount, targetDate, monthlyContribution)
	goal2, _ := entities.NewGoal(userID, entities.GoalTypeSavings, "Savings Goal 2", targetAmount, targetDate, monthlyContribution)
	goal3, _ := entities.NewGoal(userID, entities.GoalTypeEmergency, "Emergency Goal", targetAmount, targetDate, monthlyContribution)

	// Deactivate one savings goal
	goal2.Deactivate()

	ctx := context.Background()
	repo.Save(ctx, goal1)
	repo.Save(ctx, goal2)
	repo.Save(ctx, goal3)

	// Test CountActiveGoalsByType
	count, err := repo.CountActiveGoalsByType(ctx, userID, entities.GoalTypeSavings)
	if err != nil {
		t.Fatalf("Failed to count active goals: %v", err)
	}

	if count != 1 {
		t.Errorf("Expected 1 active savings goal, got %d", count)
	}

	emergencyCount, err := repo.CountActiveGoalsByType(ctx, userID, entities.GoalTypeEmergency)
	if err != nil {
		t.Fatalf("Failed to count active emergency goals: %v", err)
	}

	if emergencyCount != 1 {
		t.Errorf("Expected 1 active emergency goal, got %d", emergencyCount)
	}
}
