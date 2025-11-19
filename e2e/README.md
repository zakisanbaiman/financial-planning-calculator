# E2E Tests for Financial Planning Calculator

## Overview

This directory contains end-to-end (E2E) tests for the Financial Planning Calculator application using Playwright.

## Prerequisites

- Node.js 18+ installed
- Frontend server running on `http://localhost:3000`
- Backend server running on `http://localhost:8080`

## Installation

```bash
cd e2e
npm install
npm run install  # Install Playwright browsers
```

## Running Tests

### Run all tests

```bash
npm test
```

### Run tests in headed mode (see browser)

```bash
npm run test:headed
```

### Run tests in debug mode

```bash
npm run test:debug
```

### Run tests with UI mode

```bash
npm run test:ui
```

### Run tests for specific browser

```bash
npm run test:chromium
npm run test:firefox
npm run test:webkit
```

### Run mobile tests

```bash
npm run test:mobile
```

### Run specific test suites

```bash
npm run test:api          # API integration tests
npm run test:financial    # Financial data flow tests
npm run test:goals        # Goals management tests
```

## Test Structure

```
e2e/
├── tests/
│   ├── financial-data-flow.spec.ts  # Tests for financial data input and calculations
│   ├── goals-management.spec.ts     # Tests for goal creation and tracking
│   └── api-integration.spec.ts      # Tests for API endpoints
├── playwright.config.ts              # Playwright configuration
├── package.json                      # Dependencies and scripts
└── README.md                         # This file
```

## Test Scenarios

### Financial Data Flow Tests

- Homepage navigation
- Financial data input
- Asset projection calculation
- Retirement needs calculation
- Emergency fund calculation
- Form validation
- Error handling
- Data persistence

### Goals Management Tests

- Goal creation
- Goal progress tracking
- Goal recommendations
- Goal editing
- Goal deletion
- Goal filtering
- Goals summary visualization
- Concurrent updates

### API Integration Tests

- Health checks
- Calculation endpoints
- Validation errors
- Rate limiting
- CORS configuration
- Response format consistency
- Concurrent requests
- Response time measurement

## Configuration

### Environment Variables

- `BASE_URL`: Frontend URL (default: `http://localhost:3000`)
- `API_URL`: Backend API URL (default: `http://localhost:8080/api`)
- `CI`: Set to `true` in CI environment

### Browser Configuration

Tests run on multiple browsers by default:
- Desktop Chrome
- Desktop Firefox
- Desktop Safari
- Mobile Chrome (Pixel 5)
- Mobile Safari (iPhone 12)

## Test Reports

After running tests, view the HTML report:

```bash
npm run report
```

Reports are saved in:
- `test-results/html/` - HTML report
- `test-results/results.json` - JSON report

## Screenshots and Videos

- Screenshots are captured on test failure
- Videos are recorded for failed tests
- Traces are collected on first retry

Files are saved in `test-results/` directory.

## CI/CD Integration

### GitHub Actions Example

```yaml
name: E2E Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'
      
      - name: Install dependencies
        run: |
          cd e2e
          npm install
          npm run install
      
      - name: Start servers
        run: |
          cd frontend && npm run dev &
          cd backend && go run main.go &
          sleep 10
      
      - name: Run E2E tests
        run: |
          cd e2e
          npm test
      
      - name: Upload test results
        if: always()
        uses: actions/upload-artifact@v3
        with:
          name: playwright-report
          path: e2e/test-results/
```

## Best Practices

### Writing Tests

1. **Use data-testid attributes** for reliable selectors
2. **Wait for network idle** before assertions
3. **Use explicit waits** instead of arbitrary timeouts
4. **Clean up test data** after each test
5. **Make tests independent** - don't rely on test order

### Debugging

1. Use `--debug` flag to step through tests
2. Use `page.pause()` to pause execution
3. Check screenshots and videos in test-results
4. Use `--headed` to see browser actions

### Performance

1. Run tests in parallel when possible
2. Use `--project` to run specific browser tests
3. Reuse browser contexts when appropriate
4. Mock external API calls when needed

## Troubleshooting

### Tests fail with "Target closed"

- Ensure servers are running before tests
- Increase timeout in playwright.config.ts
- Check for JavaScript errors in browser console

### Tests are flaky

- Add explicit waits for elements
- Use `waitForLoadState('networkidle')`
- Increase timeout for slow operations
- Check for race conditions

### Cannot connect to servers

- Verify frontend is running on port 3000
- Verify backend is running on port 8080
- Check firewall settings
- Ensure no port conflicts

## Resources

- [Playwright Documentation](https://playwright.dev/)
- [Best Practices](https://playwright.dev/docs/best-practices)
- [API Reference](https://playwright.dev/docs/api/class-playwright)
- [Debugging Guide](https://playwright.dev/docs/debug)
