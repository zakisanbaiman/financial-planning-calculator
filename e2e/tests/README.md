# E2E Scenario Tests

This directory contains end-to-end scenario tests for the Financial Planning Calculator application.

## Test Files

### 1. health-check.spec.ts
Basic health checks to verify the application is running.
- Backend API health check
- Frontend homepage load check

### 2. goals-scenario.spec.ts
Comprehensive scenario tests for goals management functionality.

**Test Scenarios:**
- Create financial data and then create a goal
- Create multiple goals and retrieve them
- Update a goal
- Update goal progress
- Get goal recommendations
- Analyze goal feasibility
- Delete a goal
- Filter goals by type
- Error handling: Create goal without financial data
- Error handling: Create goal with invalid data
- Complete financial planning flow (integration test)

### 3. financial-data-scenario.spec.ts
Tests for financial data creation, updates, and management.

**Test Scenarios:**
- Create complete financial profile
- Update financial profile
- Add and update retirement data
- Add emergency fund settings
- Retrieve financial data
- Complete financial setup with goals
- Delete financial data
- Error handling: Get non-existent financial data
- Error handling: Invalid financial data

### 4. calculation-scenario.spec.ts
Tests for calculation endpoints with goals and financial data.

**Test Scenarios:**
- Calculate asset projection
- Calculate retirement projection
- Calculate emergency fund projection
- Calculate comprehensive projection
- Calculate goal projection
- Multiple goals with comprehensive calculation
- Validate calculation with different timeframes
- Error handling: Calculate without financial data
- Error handling: Invalid calculation parameters
- Complete planning with calculations and reports

## Running the Tests

### Prerequisites
1. Install dependencies:
   ```bash
   npm install
   ```

2. Install Playwright browsers:
   ```bash
   npx playwright install --with-deps
   ```

3. Start the application (backend and frontend must be running):
   - Backend: `http://localhost:8080`
   - Frontend: `http://localhost:3000`

### Run Tests

Run all tests:
```bash
npm test
```

Run specific test file:
```bash
npm test goals-scenario.spec.ts
```

Run tests in headed mode (with browser UI):
```bash
npm run test:headed
```

Run tests in debug mode:
```bash
npm run test:debug
```

Run tests with UI mode:
```bash
npm run test:ui
```

Run tests for specific browser:
```bash
npm run test:chromium
npm run test:firefox
npm run test:webkit
```

Run mobile tests:
```bash
npm run test:mobile
```

## Test Structure

Each scenario test follows this pattern:

1. **Setup**: Generate unique test user ID and create necessary data
2. **Action**: Perform the test action (create, update, delete, etc.)
3. **Verification**: Assert expected results
4. **Cleanup**: Tests use unique user IDs to avoid conflicts

## Helper Functions

### `generateTestUserId()`
Generates a unique test user ID for each test to avoid conflicts.

### `createFinancialData(request, userId)`
Helper function to create financial data for a user. Used as a prerequisite for goal creation tests.

### `setupCompleteFinancialProfile(request, userId)`
Helper function to create a complete financial profile including:
- Financial data
- Retirement data
- Emergency fund settings

## Test Data

Tests use realistic Japanese financial data:
- Monthly income: 500,000 - 600,000 JPY
- Housing costs: 120,000 - 150,000 JPY
- Food expenses: 60,000 - 80,000 JPY
- Savings: 1,000,000 - 3,000,000 JPY
- Investment returns: 5% - 6%
- Inflation rate: 2% - 2.5%

## Coverage

These scenario tests cover:
- ✅ Goal creation with financial data
- ✅ Goal updates and progress tracking
- ✅ Goal recommendations and feasibility analysis
- ✅ Multiple goal management
- ✅ Financial data CRUD operations
- ✅ Retirement planning data
- ✅ Emergency fund settings
- ✅ Asset projection calculations
- ✅ Retirement projection calculations
- ✅ Comprehensive financial planning flow
- ✅ Error handling and validation
- ✅ Integration between financial data, goals, and calculations

## CI/CD Integration

These tests can be run in CI/CD pipelines:

```bash
# Run tests in CI mode
CI=true npm test
```

The tests are configured with:
- Automatic retry on failure (2 retries in CI)
- Screenshot on failure
- Video recording on failure
- HTML report generation
- JSON report output

## Troubleshooting

### Tests fail with connection errors
Ensure backend and frontend are running:
- Backend: `http://localhost:8080/health` should return 200
- Frontend: `http://localhost:3000` should load

### Tests timeout
Increase timeout in `playwright.config.ts` if needed:
```typescript
timeout: 60 * 1000, // 60 seconds
```

### Database conflicts
Tests use unique user IDs to avoid conflicts. If you see data conflicts, ensure the database is clean or use a test database.

## Future Enhancements

Potential improvements:
- [ ] Add UI-based scenario tests (beyond API tests)
- [ ] Add performance benchmarking tests
- [ ] Add accessibility tests
- [ ] Add visual regression tests
- [ ] Add load testing scenarios
