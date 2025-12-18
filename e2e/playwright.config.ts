import { defineConfig, devices } from '@playwright/test';

/**
 * E2E Test Configuration for Financial Planning Calculator
 * 
 * See https://playwright.dev/docs/test-configuration
 * 
 * CI mode: Only runs health-check tests on Chromium for faster feedback
 * Local mode: Full test suite available
 */
export default defineConfig({
  testDir: './tests',

  // Maximum time one test can run
  timeout: 30 * 1000,

  // Test execution settings
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 1 : 0,
  workers: process.env.CI ? 1 : undefined,

  // In CI, only run health-check tests
  testMatch: process.env.CI ? 'health-check.spec.ts' : '**/*.spec.ts',

  // Reporter configuration
  reporter: [
    ['html', { outputFolder: 'playwright-report' }],
    ['json', { outputFile: 'test-results/results.json' }],
    ['list'],
  ],

  // Shared settings for all tests
  use: {
    // Base URL for tests
    baseURL: process.env.BASE_URL || 'http://localhost:3000',

    // Collect trace on failure
    trace: 'on-first-retry',

    // Screenshot on failure
    screenshot: 'only-on-failure',

    // Video on failure - disabled in CI for speed
    video: process.env.CI ? 'off' : 'retain-on-failure',

    // API endpoint
    extraHTTPHeaders: {
      'Accept': 'application/json',
    },
  },

  // Configure projects - Chromium only (other browsers can be run manually if needed)
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
  ],

  // Run local dev server before starting tests
  webServer: [
    {
      command: 'cd ../frontend && npm run dev',
      url: 'http://localhost:3000',
      reuseExistingServer: true,
      timeout: 120 * 1000,
    },
    {
      command: 'cd ../backend && go run main.go',
      url: 'http://localhost:8080/health',
      reuseExistingServer: true,
      timeout: 120 * 1000,
    },
  ],
});
