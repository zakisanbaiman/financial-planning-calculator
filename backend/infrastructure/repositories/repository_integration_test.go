package repositories

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/financial-planning-calculator/backend/domain/aggregates"
	"github.com/financial-planning-calculator/backend/domain/entities"
	"github.com/financial-planning-calculator/backend/domain/valueobjects"
)

// TestRepositoryFactory tests the repository factory functionality
func TestRepositoryFactory(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	factory := NewRepositoryFactory(db)

	// Test financial plan repository creation
	financialPlanRepo := factory.NewFinancialPlanRepository()
	if financialPlanRepo == nil {
		t.Error("Financial plan repository should not be nil")
	}

	// Test goal repository creation
	goalRepo := factory.NewGoalRepository()
	if goalRepo == nil {
		t.Error("Goal repository should not be nil")
	}

	// Test that repositories are properly typed
	_, ok := financialPlanRepo.(*PostgreSQLFinancialPlanRepository)
	if !ok {
		t.Error("Financial plan repository should be PostgreSQL implementation")
	}

	_, ok = goalRepo.(*PostgreSQLGoalRepository)
	if !ok {
		t.Error("Goal repository should be PostgreSQL implementation")
	}
}

// TestCrossRepositoryOperations tests operations that span multiple repositories
func TestCrossRepositoryOperations(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	userID := createTestUser(t, db)
	factory := NewRepositoryFactory(db)
	financialPlanRepo := factory.NewFinancialPlanRepository()
	goalRepo := factory.NewGoalRepository()

	ctx := context.Background()

	// Create financial plan
	plan := createTestFinancialPlan(t, userID)
	err := financialPlanRepo.Save(ctx, plan)
	if err != nil {
		t.Fatalf("Failed to save financial plan: %v", err)
	}

	// Create goals using goal repository
	targetAmount, _ := valueobjects.NewMoneyJPY(1000000)
	monthlyContribution, _ := valueobjects.NewMoneyJPY(50000)
	targetDate := time.Now().AddDate(1, 0, 0)

	goal1, err := entities.NewGoal(
		userID,
		entities.GoalTypeSavings,
		"Savings Goal",
		targetAmount,
		targetDate,
		monthlyContribution,
	)
	if err != nil {
		t.Fatalf("Failed to create goal 1: %v", err)
	}

	goal2, err := entities.NewGoal(
		userID,
		entities.GoalTypeEmergency,
		"Emergency Fund",
		targetAmount,
		targetDate,
		monthlyContribution,
	)
	if err != nil {
		t.Fatalf("Failed to create goal 2: %v", err)
	}

	// Save goals using goal repository
	err = goalRepo.Save(ctx, goal1)
	if err != nil {
		t.Fatalf("Failed to save goal 1: %v", err)
	}

	err = goalRepo.Save(ctx, goal2)
	if err != nil {
		t.Fatalf("Failed to save goal 2: %v", err)
	}

	// Retrieve financial plan and verify goals are included
	retrievedPlan, err := financialPlanRepo.FindByUserID(ctx, userID)
	if err != nil {
		t.Fatalf("Failed to retrieve financial plan: %v", err)
	}

	goals := retrievedPlan.Goals()
	if len(goals) != 2 {
		t.Errorf("Expected 2 goals in financial plan, got %d", len(goals))
	}

	// Verify goals can be retrieved independently
	independentGoals, err := goalRepo.FindByUserID(ctx, userID)
	if err != nil {
		t.Fatalf("Failed to retrieve goals independently: %v", err)
	}

	if len(independentGoals) != 2 {
		t.Errorf("Expected 2 independent goals, got %d", len(independentGoals))
	}

	// Test cascading delete - delete financial plan should not affect independently saved goals
	planID := aggregates.FinancialPlanID(retrievedPlan.Profile().ID())
	err = financialPlanRepo.Delete(ctx, planID)
	if err != nil {
		t.Fatalf("Failed to delete financial plan: %v", err)
	}

	// Goals should still exist
	remainingGoals, err := goalRepo.FindByUserID(ctx, userID)
	if err != nil {
		t.Fatalf("Failed to retrieve remaining goals: %v", err)
	}

	if len(remainingGoals) != 2 {
		t.Errorf("Expected 2 remaining goals after plan deletion, got %d", len(remainingGoals))
	}
}

// TestTransactionConsistency tests transaction consistency across operations
func TestTransactionConsistency(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	userID := createTestUser(t, db)
	factory := NewRepositoryFactory(db)
	financialPlanRepo := factory.NewFinancialPlanRepository()

	ctx := context.Background()

	// Create a plan with multiple components
	plan := createTestFinancialPlan(t, userID)

	// Add retirement data
	monthlyExpenses, _ := valueobjects.NewMoneyJPY(250000)
	pensionAmount, _ := valueobjects.NewMoneyJPY(150000)
	retirementData, err := entities.NewRetirementData(
		userID, 35, 65, 85, monthlyExpenses, pensionAmount,
	)
	if err != nil {
		t.Fatalf("Failed to create retirement data: %v", err)
	}
	plan.SetRetirementData(retirementData)

	// Add multiple goals
	for i := 0; i < 3; i++ {
		targetAmount, _ := valueobjects.NewMoneyJPY(float64((i + 1) * 1000000))
		monthlyContribution, _ := valueobjects.NewMoneyJPY(50000)
		targetDate := time.Now().AddDate(i+1, 0, 0)

		goal, err := entities.NewGoal(
			userID,
			entities.GoalTypeSavings,
			fmt.Sprintf("Goal %d", i+1),
			targetAmount,
			targetDate,
			monthlyContribution,
		)
		if err != nil {
			t.Fatalf("Failed to create goal %d: %v", i+1, err)
		}

		err = plan.AddGoal(goal)
		if err != nil {
			t.Fatalf("Failed to add goal %d to plan: %v", i+1, err)
		}
	}

	// Save the complete plan
	err = financialPlanRepo.Save(ctx, plan)
	if err != nil {
		t.Fatalf("Failed to save complete financial plan: %v", err)
	}

	// Verify all components were saved consistently
	retrievedPlan, err := financialPlanRepo.FindByUserID(ctx, userID)
	if err != nil {
		t.Fatalf("Failed to retrieve financial plan: %v", err)
	}

	// Check financial profile
	if retrievedPlan.Profile() == nil {
		t.Error("Financial profile should not be nil")
	}

	// Check retirement data
	if retrievedPlan.RetirementData() == nil {
		t.Error("Retirement data should not be nil")
	}

	// Check goals
	goals := retrievedPlan.Goals()
	if len(goals) != 3 {
		t.Errorf("Expected 3 goals, got %d", len(goals))
	}

	// Verify data integrity by checking database directly
	var financialDataCount, retirementDataCount, goalsCount int

	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM financial_data WHERE user_id = $1", string(userID)).Scan(&financialDataCount)
	if err != nil {
		t.Fatalf("Failed to count financial_data: %v", err)
	}

	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM retirement_data WHERE user_id = $1", string(userID)).Scan(&retirementDataCount)
	if err != nil {
		t.Fatalf("Failed to count retirement_data: %v", err)
	}

	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM goals WHERE user_id = $1", string(userID)).Scan(&goalsCount)
	if err != nil {
		t.Fatalf("Failed to count goals: %v", err)
	}

	if financialDataCount != 1 {
		t.Errorf("Expected 1 financial_data record, got %d", financialDataCount)
	}

	if retirementDataCount != 1 {
		t.Errorf("Expected 1 retirement_data record, got %d", retirementDataCount)
	}

	if goalsCount != 3 {
		t.Errorf("Expected 3 goals records, got %d", goalsCount)
	}
}

// TestConcurrentRepositoryAccess tests concurrent access to repositories
func TestConcurrentRepositoryAccess(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	factory := NewRepositoryFactory(db)
	goalRepo := factory.NewGoalRepository()

	// Create test users and goals
	numUsers := 5
	numGoalsPerUser := 3
	userIDs := make([]entities.UserID, numUsers)

	for i := 0; i < numUsers; i++ {
		userIDs[i] = createTestUser(t, db)
	}

	ctx := context.Background()
	var wg sync.WaitGroup
	errors := make(chan error, numUsers*numGoalsPerUser)

	// Concurrent goal creation
	for i := 0; i < numUsers; i++ {
		for j := 0; j < numGoalsPerUser; j++ {
			wg.Add(1)
			go func(userIndex, goalIndex int) {
				defer wg.Done()

				targetAmount, _ := valueobjects.NewMoneyJPY(float64((goalIndex + 1) * 100000))
				monthlyContribution, _ := valueobjects.NewMoneyJPY(10000)
				targetDate := time.Now().AddDate(1, 0, 0)

				goal, err := entities.NewGoal(
					userIDs[userIndex],
					entities.GoalTypeSavings,
					fmt.Sprintf("User%d Goal%d", userIndex, goalIndex),
					targetAmount,
					targetDate,
					monthlyContribution,
				)
				if err != nil {
					errors <- err
					return
				}

				err = goalRepo.Save(ctx, goal)
				if err != nil {
					errors <- err
					return
				}

				errors <- nil
			}(i, j)
		}
	}

	wg.Wait()
	close(errors)

	// Check for errors
	errorCount := 0
	for err := range errors {
		if err != nil {
			t.Errorf("Concurrent operation failed: %v", err)
			errorCount++
		}
	}

	if errorCount > 0 {
		t.Errorf("Failed %d concurrent operations", errorCount)
	}

	// Verify all goals were created
	totalGoals := 0
	for i := 0; i < numUsers; i++ {
		goals, err := goalRepo.FindByUserID(ctx, userIDs[i])
		if err != nil {
			t.Errorf("Failed to retrieve goals for user %d: %v", i, err)
			continue
		}
		totalGoals += len(goals)
	}

	expectedTotal := numUsers * numGoalsPerUser
	if totalGoals != expectedTotal {
		t.Errorf("Expected %d total goals, got %d", expectedTotal, totalGoals)
	}
}

// TestRepositoryPerformanceUnderLoad tests repository performance under load
func TestRepositoryPerformanceUnderLoad(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	factory := NewRepositoryFactory(db)
	goalRepo := factory.NewGoalRepository()
	userID := createTestUser(t, db)

	ctx := context.Background()

	// Create a large number of goals
	numGoals := 100
	goals := make([]*entities.Goal, numGoals)

	// Measure creation time
	start := time.Now()
	for i := 0; i < numGoals; i++ {
		targetAmount, _ := valueobjects.NewMoneyJPY(float64((i + 1) * 10000))
		monthlyContribution, _ := valueobjects.NewMoneyJPY(1000)
		targetDate := time.Now().AddDate(0, i+1, 0)

		goal, err := entities.NewGoal(
			userID,
			entities.GoalTypeSavings,
			fmt.Sprintf("Performance Goal %d", i+1),
			targetAmount,
			targetDate,
			monthlyContribution,
		)
		if err != nil {
			t.Fatalf("Failed to create goal %d: %v", i+1, err)
		}

		err = goalRepo.Save(ctx, goal)
		if err != nil {
			t.Fatalf("Failed to save goal %d: %v", i+1, err)
		}

		goals[i] = goal
	}
	creationTime := time.Since(start)

	t.Logf("Created %d goals in %v (avg: %v per goal)",
		numGoals, creationTime, creationTime/time.Duration(numGoals))

	// Measure retrieval time
	start = time.Now()
	retrievedGoals, err := goalRepo.FindByUserID(ctx, userID)
	if err != nil {
		t.Fatalf("Failed to retrieve goals: %v", err)
	}
	retrievalTime := time.Since(start)

	t.Logf("Retrieved %d goals in %v", len(retrievedGoals), retrievalTime)

	if len(retrievedGoals) != numGoals {
		t.Errorf("Expected %d goals, got %d", numGoals, len(retrievedGoals))
	}

	// Measure individual goal retrieval time
	start = time.Now()
	for i := 0; i < 10; i++ { // Test first 10 goals
		_, err := goalRepo.FindByID(ctx, goals[i].ID())
		if err != nil {
			t.Errorf("Failed to retrieve goal %d: %v", i, err)
		}
	}
	individualRetrievalTime := time.Since(start)

	t.Logf("Retrieved 10 individual goals in %v (avg: %v per goal)",
		individualRetrievalTime, individualRetrievalTime/10)

	// Performance assertions (adjust thresholds as needed)
	avgCreationTime := creationTime / time.Duration(numGoals)
	if avgCreationTime > 10*time.Millisecond {
		t.Logf("Warning: Average goal creation time %v exceeds 10ms threshold", avgCreationTime)
	}

	if retrievalTime > 100*time.Millisecond {
		t.Logf("Warning: Bulk retrieval time %v exceeds 100ms threshold", retrievalTime)
	}
}

// TestDatabaseConnectionPooling tests database connection pooling behavior
func TestDatabaseConnectionPooling(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	// Configure connection pool
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	factory := NewRepositoryFactory(db)
	goalRepo := factory.NewGoalRepository()
	userID := createTestUser(t, db)

	ctx := context.Background()

	// Test concurrent operations that would require multiple connections
	numConcurrentOps := 20
	var wg sync.WaitGroup
	errors := make(chan error, numConcurrentOps)

	for i := 0; i < numConcurrentOps; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			targetAmount, _ := valueobjects.NewMoneyJPY(float64((index + 1) * 10000))
			monthlyContribution, _ := valueobjects.NewMoneyJPY(1000)
			targetDate := time.Now().AddDate(0, index+1, 0)

			goal, err := entities.NewGoal(
				userID,
				entities.GoalTypeSavings,
				fmt.Sprintf("Pool Test Goal %d", index+1),
				targetAmount,
				targetDate,
				monthlyContribution,
			)
			if err != nil {
				errors <- err
				return
			}

			// Save goal
			err = goalRepo.Save(ctx, goal)
			if err != nil {
				errors <- err
				return
			}

			// Immediately retrieve it
			_, err = goalRepo.FindByID(ctx, goal.ID())
			if err != nil {
				errors <- err
				return
			}

			errors <- nil
		}(i)
	}

	wg.Wait()
	close(errors)

	// Check for errors
	for err := range errors {
		if err != nil {
			t.Errorf("Connection pool operation failed: %v", err)
		}
	}

	// Verify database stats
	stats := db.Stats()
	t.Logf("Database connection stats - Open: %d, InUse: %d, Idle: %d",
		stats.OpenConnections, stats.InUse, stats.Idle)

	if stats.OpenConnections > 10 {
		t.Errorf("Expected max 10 open connections, got %d", stats.OpenConnections)
	}
}
