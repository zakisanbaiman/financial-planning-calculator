import { test, expect } from '@playwright/test';

/**
 * E2E Test: Basic Health Checks
 * 
 * Minimal tests to verify the application is running
 * These are the only tests run in CI for fast feedback
 */
test.describe('Health Checks', () => {
  test('should check backend API health', async ({ request }) => {
    const apiUrl = process.env.API_URL || 'http://localhost:8080';
    const response = await request.get(`${apiUrl}/health`);

    expect(response.ok()).toBeTruthy();

    const data = await response.json();
    expect(data.status).toBe('ok');
  });

  test('should load frontend homepage', async ({ page }) => {
    await page.goto('/');

    // Wait for page to be ready
    await page.waitForLoadState('domcontentloaded');

    // Check that the page loaded (more flexible check)
    const title = await page.title();
    expect(title).toBeTruthy();
  });
});
