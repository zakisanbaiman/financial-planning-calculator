# Database Integration Tests

This document describes the comprehensive database integration tests for the financial planning calculator backend.

## Overview

The integration tests verify the correct operation of the repository layer, ensuring data integrity, performance, and reliability under various conditions.

## Test Categories

### 1. Repository CRUD Operations

**Files:** `postgresql_financial_plan_repository_test.go`, `postgresql_goal_repository_test.go`

**Coverage:**
- ✅ Create operations (Save)
- ✅ Read operations (FindByID, FindByUserID)
- ✅ Update operations (Update)
- ✅ Delete operations (Delete)
- ✅ Existence checks (Exists, ExistsByUserID)

**Key Test Cases:**
- `TestPostgreSQLFinancialPlanRepository_Save` - Verifies financial plan creation
- `TestPostgreSQLFinancialPlanRepository_FindByUserID` - Tests plan retrieval with all components
- `TestPostgreSQLFinancialPlanRepository_SaveWithRetirementData` - Tests complex object persistence
- `TestPostgreSQLFinancialPlanRepository_SaveWithGoals` - Tests aggregate relationships
- `TestPostgreSQLGoalRepository_*` - Complete goal repository operations

### 2. Data Integrity and Constraints

**Coverage:**
- ✅ Foreign key constraints
- ✅ Unique constraints
- ✅ Check constraints (positive amounts, valid dates)
- ✅ Cascade operations
- ✅ Transaction rollback on errors

**Key Test Cases:**
- `TestPostgreSQLFinancialPlanRepository_DataIntegrity` - Foreign key and unique constraint validation
- `TestPostgreSQLFinancialPlanRepository_TransactionRollback` - Transaction consistency

### 3. Cross-Repository Operations

**File:** `repository_integration_test.go`

**Coverage:**
- ✅ Repository factory functionality
- ✅ Operations spanning multiple repositories
- ✅ Data consistency across repositories
- ✅ Cascading operations

**Key Test Cases:**
- `TestRepositoryFactory` - Factory pattern implementation
- `TestCrossRepositoryOperations` - Multi-repository workflows
- `TestTransactionConsistency` - ACID properties verification

### 4. Concurrent Access and Performance

**Files:** `repository_integration_test.go`, `database_performance_test.go`

**Coverage:**
- ✅ Concurrent read/write operations
- ✅ Connection pool management
- ✅ Deadlock detection and handling
- ✅ Performance under load
- ✅ Memory leak detection

**Key Test Cases:**
- `TestConcurrentRepositoryAccess` - Multi-threaded access patterns
- `TestDatabaseStressTest` - High-load scenarios
- `TestDatabaseConnectionLeaks` - Resource management
- `TestDatabaseDeadlockDetection` - Concurrent update scenarios

### 5. Performance Benchmarks

**File:** `database_performance_test.go`

**Benchmarks:**
- `BenchmarkGoalRepository_Save` - Goal creation performance
- `BenchmarkGoalRepository_FindByID` - Single goal retrieval
- `BenchmarkGoalRepository_FindByUserID` - Bulk goal retrieval
- `BenchmarkFinancialPlanRepository_Save` - Complex object persistence

## Running the Tests

### Prerequisites

1. **PostgreSQL Database**
   ```bash
   # Using Docker
   docker-compose up -d postgres
   
   # Or local PostgreSQL installation
   brew install postgresql  # macOS
   sudo apt-get install postgresql  # Ubuntu
   ```

2. **Environment Variables**
   ```bash
   export DB_HOST=localhost
   export DB_PORT=5432
   export DB_USER=postgres
   export DB_PASSWORD=password
   export DB_NAME=financial_planning_test
   export DB_SSLMODE=disable
   ```

### Running All Integration Tests

```bash
# Run the comprehensive test suite
./backend/scripts/run-integration-tests.sh

# Run with coverage report
./backend/scripts/run-integration-tests.sh --coverage
```

### Running Specific Test Categories

```bash
# Repository CRUD tests only
go test -v ./infrastructure/repositories -run "Test.*Repository.*"

# Performance and stress tests
go test -v ./infrastructure/repositories -run "TestDatabaseStressTest|TestDatabaseConnectionLeaks"

# Benchmarks
go test -v ./infrastructure/repositories -bench "Benchmark.*" -benchtime=5s

# Short tests (skip database-dependent tests)
go test -short ./infrastructure/repositories
```

## Test Data and Fixtures

### Test User Creation
```go
func createTestUser(t *testing.T, db *sql.DB) entities.UserID {
    var userID string
    err := db.QueryRow("SELECT uuid_generate_v4()::text").Scan(&userID)
    // ... insert user
    return entities.UserID(userID)
}
```

### Test Financial Plan Creation
```go
func createTestFinancialPlan(t *testing.T, userID entities.UserID) *aggregates.FinancialPlan {
    // Creates a complete financial plan with:
    // - Monthly income: ¥400,000
    // - Expenses: Housing ¥120,000, Food ¥60,000
    // - Savings: Deposit ¥1,000,000, Investment ¥500,000
    // - Investment return: 5%
    // - Inflation rate: 2%
}
```

## Performance Expectations

### Benchmarks Targets

| Operation | Target Performance | Notes |
|-----------|-------------------|-------|
| Goal Save | < 10ms per operation | Single goal creation |
| Goal FindByID | < 5ms per operation | Single goal retrieval |
| Goal FindByUserID | < 50ms per operation | Bulk retrieval (50 goals) |
| Financial Plan Save | < 100ms per operation | Complex object with relationships |

### Stress Test Metrics

| Metric | Target | Actual (Example) |
|--------|--------|------------------|
| Concurrent Users | 100 users | ✅ Supported |
| Goals per User | 20 goals | ✅ Supported |
| Concurrent Workers | 20 workers | ✅ Supported |
| Error Rate | < 10% | ✅ < 5% typical |
| Operations/Second | > 100 ops/sec | ✅ ~200 ops/sec |

## Database Schema Validation

The tests verify the following schema constraints:

### Financial Data Table
- ✅ Unique constraint on `user_id`
- ✅ Positive amount checks
- ✅ Rate percentage limits (0-100% for investment, 0-50% for inflation)

### Goals Table
- ✅ Valid goal types (savings, retirement, emergency, custom)
- ✅ Positive target amounts
- ✅ Future target dates
- ✅ Non-negative current amounts and contributions

### Retirement Data Table
- ✅ Age validations (current < retirement < life expectancy)
- ✅ Positive expense and pension amounts

## Error Handling Verification

### Database Errors
- ✅ Connection failures (graceful degradation)
- ✅ Constraint violations (proper error messages)
- ✅ Transaction rollbacks (data consistency)
- ✅ Deadlock detection (automatic retry logic)

### Application Errors
- ✅ Invalid domain objects (validation at entity level)
- ✅ Missing required fields (comprehensive validation)
- ✅ Business rule violations (domain-specific constraints)

## Continuous Integration

### CI Pipeline Integration
```yaml
# Example GitHub Actions configuration
- name: Run Integration Tests
  run: |
    docker-compose up -d postgres
    sleep 10  # Wait for PostgreSQL to be ready
    ./backend/scripts/run-integration-tests.sh
  env:
    DB_HOST: localhost
    DB_PASSWORD: test_password
```

### Test Isolation
- Each test creates its own test users
- Tests clean up after themselves
- Database is reset between test runs
- No shared state between tests

## Troubleshooting

### Common Issues

1. **Database Connection Failures**
   ```
   Error: Database connection failed
   Solution: Ensure PostgreSQL is running and credentials are correct
   ```

2. **Migration Failures**
   ```
   Error: Migration failed
   Solution: Check database permissions and schema conflicts
   ```

3. **Performance Test Timeouts**
   ```
   Error: Test timeout exceeded
   Solution: Increase timeout or reduce test load
   ```

### Debug Mode
```bash
# Enable verbose logging
export DB_LOG_LEVEL=debug
go test -v ./infrastructure/repositories -run TestSpecificTest
```

## Future Enhancements

### Planned Improvements
- [ ] Add Redis integration tests for caching layer
- [ ] Implement database migration rollback tests
- [ ] Add cross-database compatibility tests (MySQL, SQLite)
- [ ] Performance regression detection
- [ ] Automated performance baseline updates

### Monitoring Integration
- [ ] Prometheus metrics collection during tests
- [ ] Performance trend analysis
- [ ] Automated alerting for performance degradation

## Contributing

When adding new repository methods or modifying existing ones:

1. **Add corresponding integration tests**
2. **Update performance benchmarks if applicable**
3. **Verify data integrity constraints**
4. **Test concurrent access scenarios**
5. **Update this documentation**

### Test Naming Convention
```go
// Repository method tests
func TestPostgreSQLXxxRepository_MethodName(t *testing.T)

// Integration tests
func TestCrossRepository_FeatureName(t *testing.T)

// Performance tests
func TestDatabase_PerformanceAspect(t *testing.T)

// Benchmarks
func BenchmarkXxxRepository_MethodName(b *testing.B)
```