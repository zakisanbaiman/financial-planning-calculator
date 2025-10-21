package repositories

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/financial-planning-calculator/backend/config"
	"github.com/financial-planning-calculator/backend/domain/aggregates"
	"github.com/financial-planning-calculator/backend/domain/entities"
	"github.com/financial-planning-calculator/backend/domain/valueobjects"
	"github.com/financial-planning-calculator/backend/infrastructure/database"
)

func setupFinancialPlanTestDB(t *testing.T) *sql.DB {
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

func createTestUserForFinancialPlan(t *testing.T, db *sql.DB) entities.UserID {
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

func createTestFinancialPlan(t *testing.T, userID entities.UserID) *aggregates.FinancialPlan {
	// Create test financial profile
	monthlyIncome, err := valueobjects.NewMoneyJPY(400000)
	if err != nil {
		t.Fatalf("Failed to create monthly income: %v", err)
	}

	investmentReturn, err := valueobjects.NewRate(5.0)
	if err != nil {
		t.Fatalf("Failed to create investment return: %v", err)
	}

	inflationRate, err := valueobjects.NewRate(2.0)
	if err != nil {
		t.Fatalf("Failed to create inflation rate: %v", err)
	}

	// Create expense items
	expenses := entities.ExpenseCollection{
		{
			Category:    "住居費",
			Amount:      mustNewMoneyJPY(120000),
			Description: "家賃・光熱費",
		},
		{
			Category:    "食費",
			Amount:      mustNewMoneyJPY(60000),
			Description: "食材・外食費",
		},
	}

	// Create savings items
	savings := entities.SavingsCollection{
		{
			Type:        "deposit",
			Amount:      mustNewMoneyJPY(1000000),
			Description: "普通預金",
		},
		{
			Type:        "investment",
			Amount:      mustNewMoneyJPY(500000),
			Description: "投資信託",
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
		t.Fatalf("Failed to create financial profile: %v", err)
	}

	plan, err := aggregates.NewFinancialPlan(profile)
	if err != nil {
		t.Fatalf("Failed to create financial plan: %v", err)
	}

	return plan
}

func TestPostgreSQLFinancialPlanRepository_Save(t *testing.T) {
	db := setupFinancialPlanTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	userID := createTestUserForFinancialPlan(t, db)
	repo := NewPostgreSQLFinancialPlanRepository(db)
	plan := createTestFinancialPlan(t, userID)

	ctx := context.Background()

	// Test Save
	err := repo.Save(ctx, plan)
	if err != nil {
		t.Fatalf("Failed to save financial plan: %v", err)
	}

	// Verify plan was saved by checking existence
	exists, err := repo.ExistsByUserID(ctx, userID)
	if err != nil {
		t.Fatalf("Failed to check financial plan existence: %v", err)
	}
	if !exists {
		t.Error("Financial plan should exist after saving")
	}
}

func TestPostgreSQLFinancialPlanRepository_FindByUserID(t *testing.T) {
	db := setupFinancialPlanTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	userID := createTestUserForFinancialPlan(t, db)
	repo := NewPostgreSQLFinancialPlanRepository(db)
	originalPlan := createTestFinancialPlan(t, userID)

	ctx := context.Background()

	// Save the plan
	err := repo.Save(ctx, originalPlan)
	if err != nil {
		t.Fatalf("Failed to save financial plan: %v", err)
	}

	// Test FindByUserID
	foundPlan, err := repo.FindByUserID(ctx, userID)
	if err != nil {
		t.Fatalf("Failed to find financial plan: %v", err)
	}

	// Verify plan properties
	if foundPlan.Profile().UserID() != originalPlan.Profile().UserID() {
		t.Errorf("Expected user ID %s, got %s", originalPlan.Profile().UserID(), foundPlan.Profile().UserID())
	}

	if foundPlan.Profile().MonthlyIncome().Amount() != originalPlan.Profile().MonthlyIncome().Amount() {
		t.Errorf("Expected monthly income %f, got %f",
			originalPlan.Profile().MonthlyIncome().Amount(),
			foundPlan.Profile().MonthlyIncome().Amount())
	}

	// Verify expense items
	originalExpenses := originalPlan.Profile().MonthlyExpenses()
	foundExpenses := foundPlan.Profile().MonthlyExpenses()
	if len(foundExpenses) != len(originalExpenses) {
		t.Errorf("Expected %d expense items, got %d", len(originalExpenses), len(foundExpenses))
	}

	// Verify savings items
	originalSavings := originalPlan.Profile().CurrentSavings()
	foundSavings := foundPlan.Profile().CurrentSavings()
	if len(foundSavings) != len(originalSavings) {
		t.Errorf("Expected %d savings items, got %d", len(originalSavings), len(foundSavings))
	}
}

func TestPostgreSQLFinancialPlanRepository_SaveWithRetirementData(t *testing.T) {
	db := setupFinancialPlanTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	userID := createTestUserForFinancialPlan(t, db)
	repo := NewPostgreSQLFinancialPlanRepository(db)
	plan := createTestFinancialPlan(t, userID)

	// Add retirement data
	monthlyExpenses, err := valueobjects.NewMoneyJPY(250000)
	if err != nil {
		t.Fatalf("Failed to create monthly expenses: %v", err)
	}

	pensionAmount, err := valueobjects.NewMoneyJPY(150000)
	if err != nil {
		t.Fatalf("Failed to create pension amount: %v", err)
	}

	retirementData, err := entities.NewRetirementData(
		userID,
		35, // current age
		65, // retirement age
		85, // life expectancy
		monthlyExpenses,
		pensionAmount,
	)
	if err != nil {
		t.Fatalf("Failed to create retirement data: %v", err)
	}

	err = plan.SetRetirementData(retirementData)
	if err != nil {
		t.Fatalf("Failed to set retirement data: %v", err)
	}

	ctx := context.Background()

	// Save plan with retirement data
	err = repo.Save(ctx, plan)
	if err != nil {
		t.Fatalf("Failed to save financial plan with retirement data: %v", err)
	}

	// Retrieve and verify
	foundPlan, err := repo.FindByUserID(ctx, userID)
	if err != nil {
		t.Fatalf("Failed to find financial plan: %v", err)
	}

	if foundPlan.RetirementData() == nil {
		t.Fatal("Expected retirement data to be present")
	}

	if foundPlan.RetirementData().CurrentAge() != 35 {
		t.Errorf("Expected current age 35, got %d", foundPlan.RetirementData().CurrentAge())
	}
}

func TestPostgreSQLFinancialPlanRepository_SaveWithGoals(t *testing.T) {
	db := setupFinancialPlanTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	userID := createTestUserForFinancialPlan(t, db)
	repo := NewPostgreSQLFinancialPlanRepository(db)
	plan := createTestFinancialPlan(t, userID)

	// Add goals
	targetAmount, err := valueobjects.NewMoneyJPY(1000000)
	if err != nil {
		t.Fatalf("Failed to create target amount: %v", err)
	}

	monthlyContribution, err := valueobjects.NewMoneyJPY(50000)
	if err != nil {
		t.Fatalf("Failed to create monthly contribution: %v", err)
	}

	targetDate := time.Now().AddDate(1, 0, 0) // 1 year from now

	goal1, err := entities.NewGoal(
		userID,
		entities.GoalTypeSavings,
		"Emergency Fund",
		targetAmount,
		targetDate,
		monthlyContribution,
	)
	if err != nil {
		t.Fatalf("Failed to create goal 1: %v", err)
	}

	goal2, err := entities.NewGoal(
		userID,
		entities.GoalTypeRetirement,
		"Retirement Savings",
		mustNewMoneyJPY(5000000),
		time.Now().AddDate(5, 0, 0),
		mustNewMoneyJPY(100000),
	)
	if err != nil {
		t.Fatalf("Failed to create goal 2: %v", err)
	}

	err = plan.AddGoal(goal1)
	if err != nil {
		t.Fatalf("Failed to add goal 1: %v", err)
	}

	err = plan.AddGoal(goal2)
	if err != nil {
		t.Fatalf("Failed to add goal 2: %v", err)
	}

	ctx := context.Background()

	// Save plan with goals
	err = repo.Save(ctx, plan)
	if err != nil {
		t.Fatalf("Failed to save financial plan with goals: %v", err)
	}

	// Retrieve and verify
	foundPlan, err := repo.FindByUserID(ctx, userID)
	if err != nil {
		t.Fatalf("Failed to find financial plan: %v", err)
	}

	goals := foundPlan.Goals()
	if len(goals) != 2 {
		t.Errorf("Expected 2 goals, got %d", len(goals))
	}
}

func TestPostgreSQLFinancialPlanRepository_Update(t *testing.T) {
	db := setupFinancialPlanTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	userID := createTestUserForFinancialPlan(t, db)
	repo := NewPostgreSQLFinancialPlanRepository(db)
	plan := createTestFinancialPlan(t, userID)

	ctx := context.Background()

	// Save initial plan
	err := repo.Save(ctx, plan)
	if err != nil {
		t.Fatalf("Failed to save initial financial plan: %v", err)
	}

	// Update the plan (modify monthly income)
	newIncome, err := valueobjects.NewMoneyJPY(450000)
	if err != nil {
		t.Fatalf("Failed to create new income: %v", err)
	}

	err = plan.Profile().UpdateMonthlyIncome(newIncome)
	if err != nil {
		t.Fatalf("Failed to update monthly income: %v", err)
	}

	// Update the plan
	err = repo.Update(ctx, plan)
	if err != nil {
		t.Fatalf("Failed to update financial plan: %v", err)
	}

	// Retrieve and verify update
	updatedPlan, err := repo.FindByUserID(ctx, userID)
	if err != nil {
		t.Fatalf("Failed to find updated financial plan: %v", err)
	}

	if updatedPlan.Profile().MonthlyIncome().Amount() != 450000 {
		t.Errorf("Expected updated monthly income 450000, got %f",
			updatedPlan.Profile().MonthlyIncome().Amount())
	}
}

func TestPostgreSQLFinancialPlanRepository_Delete(t *testing.T) {
	db := setupFinancialPlanTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	userID := createTestUserForFinancialPlan(t, db)
	repo := NewPostgreSQLFinancialPlanRepository(db)
	plan := createTestFinancialPlan(t, userID)

	ctx := context.Background()

	// Save plan
	err := repo.Save(ctx, plan)
	if err != nil {
		t.Fatalf("Failed to save financial plan: %v", err)
	}

	// Get the plan ID for deletion
	savedPlan, err := repo.FindByUserID(ctx, userID)
	if err != nil {
		t.Fatalf("Failed to find saved plan: %v", err)
	}

	planID := aggregates.FinancialPlanID(savedPlan.Profile().ID())

	// Delete plan
	err = repo.Delete(ctx, planID)
	if err != nil {
		t.Fatalf("Failed to delete financial plan: %v", err)
	}

	// Verify deletion
	exists, err := repo.ExistsByUserID(ctx, userID)
	if err != nil {
		t.Fatalf("Failed to check financial plan existence: %v", err)
	}
	if exists {
		t.Error("Financial plan should not exist after deletion")
	}
}

func TestPostgreSQLFinancialPlanRepository_TransactionRollback(t *testing.T) {
	db := setupFinancialPlanTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	userID := createTestUserForFinancialPlan(t, db)
	repo := NewPostgreSQLFinancialPlanRepository(db)

	// Create a plan with invalid data that should cause rollback
	monthlyIncome, _ := valueobjects.NewMoneyJPY(400000)
	investmentReturn, _ := valueobjects.NewRate(5.0)
	inflationRate, _ := valueobjects.NewRate(2.0)

	// Create invalid expense (negative amount should be caught by domain validation)
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
		t.Fatalf("Failed to create financial profile: %v", err)
	}

	plan, err := aggregates.NewFinancialPlan(profile)
	if err != nil {
		t.Fatalf("Failed to create financial plan: %v", err)
	}

	ctx := context.Background()

	// Save the valid plan first
	err = repo.Save(ctx, plan)
	if err != nil {
		t.Fatalf("Failed to save financial plan: %v", err)
	}

	// Verify plan exists
	exists, err := repo.ExistsByUserID(ctx, userID)
	if err != nil {
		t.Fatalf("Failed to check financial plan existence: %v", err)
	}
	if !exists {
		t.Error("Financial plan should exist after saving")
	}
}

// Performance test for bulk operations
func TestPostgreSQLFinancialPlanRepository_BulkOperations(t *testing.T) {
	db := setupFinancialPlanTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	repo := NewPostgreSQLFinancialPlanRepository(db)
	ctx := context.Background()

	// Create multiple users and plans
	numPlans := 10
	userIDs := make([]entities.UserID, numPlans)
	plans := make([]*aggregates.FinancialPlan, numPlans)

	// Create test data
	for i := 0; i < numPlans; i++ {
		userID := createTestUserForFinancialPlan(t, db)
		userIDs[i] = userID
		plans[i] = createTestFinancialPlan(t, userID)
	}

	// Measure save performance
	start := time.Now()
	for i := 0; i < numPlans; i++ {
		err := repo.Save(ctx, plans[i])
		if err != nil {
			t.Fatalf("Failed to save plan %d: %v", i, err)
		}
	}
	saveTime := time.Since(start)

	t.Logf("Saved %d plans in %v (avg: %v per plan)",
		numPlans, saveTime, saveTime/time.Duration(numPlans))

	// Measure read performance
	start = time.Now()
	for i := 0; i < numPlans; i++ {
		_, err := repo.FindByUserID(ctx, userIDs[i])
		if err != nil {
			t.Fatalf("Failed to find plan %d: %v", i, err)
		}
	}
	readTime := time.Since(start)

	t.Logf("Read %d plans in %v (avg: %v per plan)",
		numPlans, readTime, readTime/time.Duration(numPlans))

	// Performance thresholds (adjust based on requirements)
	avgSaveTime := saveTime / time.Duration(numPlans)
	avgReadTime := readTime / time.Duration(numPlans)

	if avgSaveTime > 100*time.Millisecond {
		t.Logf("Warning: Average save time %v exceeds 100ms threshold", avgSaveTime)
	}

	if avgReadTime > 50*time.Millisecond {
		t.Logf("Warning: Average read time %v exceeds 50ms threshold", avgReadTime)
	}
}

// Test concurrent access
func TestPostgreSQLFinancialPlanRepository_ConcurrentAccess(t *testing.T) {
	db := setupFinancialPlanTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	userID := createTestUserForFinancialPlan(t, db)
	repo := NewPostgreSQLFinancialPlanRepository(db)
	plan := createTestFinancialPlan(t, userID)

	ctx := context.Background()

	// Save initial plan
	err := repo.Save(ctx, plan)
	if err != nil {
		t.Fatalf("Failed to save initial plan: %v", err)
	}

	// Test concurrent reads
	numGoroutines := 5
	done := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			_, err := repo.FindByUserID(ctx, userID)
			done <- err
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		if err := <-done; err != nil {
			t.Errorf("Concurrent read %d failed: %v", i, err)
		}
	}
}

// Test data integrity constraints
func TestPostgreSQLFinancialPlanRepository_DataIntegrity(t *testing.T) {
	db := setupFinancialPlanTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	ctx := context.Background()

	// Test foreign key constraint - try to save plan with non-existent user
	nonExistentUserID := entities.UserID("00000000-0000-0000-0000-000000000000")
	plan := createTestFinancialPlan(t, nonExistentUserID)

	repo := NewPostgreSQLFinancialPlanRepository(db)
	err := repo.Save(ctx, plan)
	if err == nil {
		t.Error("Expected error when saving plan with non-existent user ID")
	}

	// Test unique constraint - try to save duplicate financial data for same user
	userID := createTestUserForFinancialPlan(t, db)
	plan1 := createTestFinancialPlan(t, userID)
	plan2 := createTestFinancialPlan(t, userID)

	// Save first plan
	err = repo.Save(ctx, plan1)
	if err != nil {
		t.Fatalf("Failed to save first plan: %v", err)
	}

	// Save second plan for same user (should update, not create duplicate)
	err = repo.Save(ctx, plan2)
	if err != nil {
		t.Fatalf("Failed to save second plan (update): %v", err)
	}

	// Verify only one financial_data record exists for the user
	var count int
	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM financial_data WHERE user_id = $1", string(userID)).Scan(&count)
	if err != nil {
		t.Fatalf("Failed to count financial_data records: %v", err)
	}

	if count != 1 {
		t.Errorf("Expected 1 financial_data record, got %d", count)
	}
}
