import { test, expect } from '@playwright/test';

/**
 * E2E Test: Basic Health Checks
 * 
 * Minimal tests to verify the application is running
 */
test.describe('Health Checks', () => {
  test('should check backend API health', async ({ request }) => {
    const response = await request.get('http://localhost:8080/health');

    expect(response.ok()).toBeTruthy();

    const data = await response.json();
    expect(data.status).toBe('ok');
  });

  test('should load frontend homepage', async ({ page }) => {
    await page.goto('/');

    // Wait for page to be ready
    await page.waitForLoadState('networkidle');

    // Check page title
    await expect(page).toHaveTitle(/財務計画計算機/);
  });
});
