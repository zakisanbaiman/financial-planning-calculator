package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/financial-planning-calculator/backend/config"
	"github.com/financial-planning-calculator/backend/domain/aggregates"
	"github.com/financial-planning-calculator/backend/domain/entities"
	"github.com/financial-planning-calculator/backend/domain/repositories"
	"github.com/financial-planning-calculator/backend/domain/valueobjects"
	"github.com/financial-planning-calculator/backend/infrastructure/database"
)

// Helper function to create Money JPY without error handling (for test data)
func mustNewMoneyJPY(amount float64) valueobjects.Money {
	money, err := valueobjects.NewMoneyJPY(amount)
	if err != nil {
		panic(fmt.Sprintf("Failed to create money: %v", err))
	}
	return money
}

// BenchmarkGoalRepository_Save benchmarks goal saving performance
func BenchmarkGoalRepository_Save(b *testing.B) {
	db := setupBenchmarkDB(b)
	if db == nil {
		return
	}
	defer db.Close()

	userID := createBenchmarkUser(b, db)
	repo := NewPostgreSQLGoalRepository(db)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		goal := createBenchmarkGoal(b, userID, i)
		err := repo.Save(ctx, goal)
		if err != nil {
			b.Fatalf("Failed to save goal: %v", err)
		}
	}
}

// BenchmarkGoalRepository_FindByID benchmarks goal retrieval by ID
func BenchmarkGoalRepository_FindByID(b *testing.B) {
	db := setupBenchmarkDB(b)
	if db == nil {
		return
	}
	defer db.Close()

	userID := createBenchmarkUser(b, db)
	repo := NewPostgreSQLGoalRepository(db)
	ctx := context.Background()

	// Pre-create goals for benchmarking
	numGoals := 1000
	goalIDs := make([]entities.GoalID, numGoals)
	for i := 0; i < numGoals; i++ {
		goal := createBenchmarkGoal(b, userID, i)
		err := repo.Save(ctx, goal)
		if err != nil {
			b.Fatalf("Failed to save goal %d: %v", i, err)
		}
		goalIDs[i] = goal.ID()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		goalID := goalIDs[i%numGoals]
		_, err := repo.FindByID(ctx, goalID)
		if err != nil {
			b.Fatalf("Failed to find goal: %v", err)
		}
	}
}

// BenchmarkGoalRepository_FindByUserID benchmarks goal retrieval by user ID
func BenchmarkGoalRepository_FindByUserID(b *testing.B) {
	db := setupBenchmarkDB(b)
	if db == nil {
		return
	}
	defer db.Close()

	userID := createBenchmarkUser(b, db)
	repo := NewPostgreSQLGoalRepository(db)
	ctx := context.Background()

	// Pre-create goals
	numGoals := 50
	for i := 0; i < numGoals; i++ {
		goal := createBenchmarkGoal(b, userID, i)
		err := repo.Save(ctx, goal)
		if err != nil {
			b.Fatalf("Failed to save goal %d: %v", i, err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := repo.FindByUserID(ctx, userID)
		if err != nil {
			b.Fatalf("Failed to find goals by user ID: %v", err)
		}
	}
}

// BenchmarkFinancialPlanRepository_Save benchmarks financial plan saving
func BenchmarkFinancialPlanRepository_Save(b *testing.B) {
	db := setupBenchmarkDB(b)
	if db == nil {
		return
	}
	defer db.Close()

	repo := NewPostgreSQLFinancialPlanRepository(db)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		userID := createBenchmarkUser(b, db)
		plan := createBenchmarkFinancialPlan(b, userID)
		err := repo.Save(ctx, plan)
		if err != nil {
			b.Fatalf("Failed to save financial plan: %v", err)
		}
	}
}

// TestDatabaseStressTest performs stress testing on the database
func TestDatabaseStressTest(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stress test in short mode")
	}

	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	// Configure connection pool for stress testing
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(time.Minute * 5)

	factory := NewRepositoryFactory(db)
	goalRepo := factory.NewGoalRepository()
	financialPlanRepo := factory.NewFinancialPlanRepository()

	ctx := context.Background()

	// Stress test parameters
	numUsers := 100
	numGoalsPerUser := 20
	numConcurrentWorkers := 20

	t.Logf("Starting stress test: %d users, %d goals per user, %d concurrent workers",
		numUsers, numGoalsPerUser, numConcurrentWorkers)

	start := time.Now()

	// Create users
	userIDs := make([]entities.UserID, numUsers)
	for i := 0; i < numUsers; i++ {
		userIDs[i] = createTestUser(t, db)
	}

	// Channel for work distribution
	type workItem struct {
		userID    entities.UserID
		goalIndex int
	}
	workChan := make(chan workItem, numUsers*numGoalsPerUser)
	errorChan := make(chan error, numUsers*numGoalsPerUser)

	// Populate work channel
	for i, userID := range userIDs {
		// Create financial plan for each user
		plan := createTestFinancialPlan(t, userID)
		err := financialPlanRepo.Save(ctx, plan)
		if err != nil {
			t.Fatalf("Failed to save financial plan for user %d: %v", i, err)
		}

		// Add goal creation work
		for j := 0; j < numGoalsPerUser; j++ {
			workChan <- workItem{userID: userID, goalIndex: j}
		}
	}
	close(workChan)

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < numConcurrentWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for work := range workChan {
				err := createAndSaveGoal(ctx, goalRepo, work.userID, work.goalIndex)
				errorChan <- err
			}
		}(i)
	}

	// Wait for completion
	wg.Wait()
	close(errorChan)

	// Check for errors
	errorCount := 0
	for err := range errorChan {
		if err != nil {
			t.Logf("Worker error: %v", err)
			errorCount++
		}
	}

	duration := time.Since(start)
	totalOperations := numUsers * numGoalsPerUser

	t.Logf("Stress test completed in %v", duration)
	t.Logf("Total operations: %d", totalOperations)
	t.Logf("Operations per second: %.2f", float64(totalOperations)/duration.Seconds())
	t.Logf("Errors: %d (%.2f%%)", errorCount, float64(errorCount)/float64(totalOperations)*100)

	if errorCount > totalOperations/10 { // Allow up to 10% error rate
		t.Errorf("Too many errors: %d out of %d operations failed", errorCount, totalOperations)
	}

	// Verify data integrity
	t.Log("Verifying data integrity...")
	for i, userID := range userIDs {
		goals, err := goalRepo.FindByUserID(ctx, userID)
		if err != nil {
			t.Errorf("Failed to retrieve goals for user %d: %v", i, err)
			continue
		}

		expectedGoals := numGoalsPerUser
		if len(goals) != expectedGoals {
			t.Errorf("User %d: expected %d goals, got %d", i, expectedGoals, len(goals))
		}
	}

	// Check database connection stats
	stats := db.Stats()
	t.Logf("Final DB stats - Open: %d, InUse: %d, Idle: %d, WaitCount: %d, WaitDuration: %v",
		stats.OpenConnections, stats.InUse, stats.Idle, stats.WaitCount, stats.WaitDuration)
}

// TestDatabaseConnectionLeaks tests for connection leaks
func TestDatabaseConnectionLeaks(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	// Set a small connection pool to detect leaks easily
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(2)

	factory := NewRepositoryFactory(db)
	goalRepo := factory.NewGoalRepository()
	userID := createTestUser(t, db)

	ctx := context.Background()

	// Record initial stats
	initialStats := db.Stats()
	t.Logf("Initial DB stats - Open: %d, InUse: %d, Idle: %d",
		initialStats.OpenConnections, initialStats.InUse, initialStats.Idle)

	// Perform many operations
	numOperations := 100
	for i := 0; i < numOperations; i++ {
		goal := createBenchmarkGoal(t, userID, i)

		err := goalRepo.Save(ctx, goal)
		if err != nil {
			t.Fatalf("Failed to save goal %d: %v", i, err)
		}

		_, err = goalRepo.FindByID(ctx, goal.ID())
		if err != nil {
			t.Fatalf("Failed to find goal %d: %v", i, err)
		}

		// Check stats periodically
		if i%20 == 0 {
			stats := db.Stats()
			t.Logf("Operation %d - Open: %d, InUse: %d, Idle: %d",
				i, stats.OpenConnections, stats.InUse, stats.Idle)
		}
	}

	// Final stats check
	finalStats := db.Stats()
	t.Logf("Final DB stats - Open: %d, InUse: %d, Idle: %d",
		finalStats.OpenConnections, finalStats.InUse, finalStats.Idle)

	// Verify no connection leaks
	if finalStats.OpenConnections > 5 {
		t.Errorf("Possible connection leak: %d open connections (max: 5)", finalStats.OpenConnections)
	}

	if finalStats.InUse > 0 {
		t.Errorf("Connections still in use: %d", finalStats.InUse)
	}
}

// Helper functions for benchmarks and stress tests

func setupBenchmarkDB(b *testing.B) *sql.DB {
	if testing.Short() {
		b.Skip("Skipping benchmark in short mode")
	}

	dbConfig := config.NewDatabaseConfig()
	db, err := config.NewDatabaseConnection(dbConfig)
	if err != nil {
		b.Skipf("Database connection failed (expected in CI): %v", err)
		return nil
	}

	migrator := database.NewMigrator(db)
	if err := migrator.Up(); err != nil {
		b.Fatalf("Migration failed: %v", err)
	}

	return db
}

func createBenchmarkUser(b testing.TB, db *sql.DB) entities.UserID {
	var userID string
	err := db.QueryRow("SELECT uuid_generate_v4()::text").Scan(&userID)
	if err != nil {
		b.Fatalf("Failed to generate UUID: %v", err)
	}

	_, err = db.Exec("INSERT INTO users (id, email) VALUES ($1, $2)", userID, userID+"@benchmark.com")
	if err != nil {
		b.Fatalf("Failed to create benchmark user: %v", err)
	}

	return entities.UserID(userID)
}

func createBenchmarkGoal(b testing.TB, userID entities.UserID, index int) *entities.Goal {
	targetAmount, err := valueobjects.NewMoneyJPY(float64((index + 1) * 100000))
	if err != nil {
		b.Fatalf("Failed to create target amount: %v", err)
	}

	monthlyContribution, err := valueobjects.NewMoneyJPY(float64((index + 1) * 1000))
	if err != nil {
		b.Fatalf("Failed to create monthly contribution: %v", err)
	}

	targetDate := time.Now().AddDate(1, index%12, 0)

	goalTypes := []entities.GoalType{
		entities.GoalTypeSavings,
		entities.GoalTypeEmergency,
		entities.GoalTypeRetirement,
		entities.GoalTypeCustom,
	}

	goal, err := entities.NewGoal(
		userID,
		goalTypes[index%len(goalTypes)],
		fmt.Sprintf("Benchmark Goal %d", index),
		targetAmount,
		targetDate,
		monthlyContribution,
	)
	if err != nil {
		b.Fatalf("Failed to create benchmark goal: %v", err)
	}

	return goal
}

func createBenchmarkFinancialPlan(b testing.TB, userID entities.UserID) *aggregates.FinancialPlan {
	monthlyIncome, err := valueobjects.NewMoneyJPY(400000)
	if err != nil {
		b.Fatalf("Failed to create monthly income: %v", err)
	}

	investmentReturn, err := valueobjects.NewRate(5.0)
	if err != nil {
		b.Fatalf("Failed to create investment return: %v", err)
	}

	inflationRate, err := valueobjects.NewRate(2.0)
	if err != nil {
		b.Fatalf("Failed to create inflation rate: %v", err)
	}

	expenses := entities.ExpenseCollection{
		{
			Category:    "住居費",
			Amount:      mustNewMoneyJPY(120000),
			Description: "家賃・光熱費",
		},
	}

	savings := entities.SavingsCollection{
		{
			Type:        "deposit",
			Amount:      mustNewMoneyJPY(1000000),
			Description: "普通預金",
		},
	}

	profile, err := entities.NewFinancialProfile(
		userID,
		monthlyIncome,
		expenses,
		savings,
		investmentReturn,
		inflationRate,
	)
	if err != nil {
		b.Fatalf("Failed to create financial profile: %v", err)
	}

	plan, err := aggregates.NewFinancialPlan(profile)
	if err != nil {
		b.Fatalf("Failed to create financial plan: %v", err)
	}

	return plan
}

func createAndSaveGoal(ctx context.Context, repo repositories.GoalRepository, userID entities.UserID, index int) error {
	targetAmount, err := valueobjects.NewMoneyJPY(float64((index + 1) * 50000))
	if err != nil {
		return err
	}

	monthlyContribution, err := valueobjects.NewMoneyJPY(float64((index + 1) * 500))
	if err != nil {
		return err
	}

	targetDate := time.Now().AddDate(0, index+1, 0)

	goal, err := entities.NewGoal(
		userID,
		entities.GoalTypeSavings,
		fmt.Sprintf("Stress Test Goal %d", index),
		targetAmount,
		targetDate,
		monthlyContribution,
	)
	if err != nil {
		return err
	}

	return repo.Save(ctx, goal)
}
