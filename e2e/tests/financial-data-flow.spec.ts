import { test, expect } from '@playwright/test';

/**
 * E2E Test: Financial Data Input Flow
 * 
 * Tests the complete flow of entering financial data and viewing calculations
 */
test.describe('Financial Data Input Flow', () => {
  const testUserId = 'e2e-test-user-' + Date.now();

  test.beforeEach(async ({ page }) => {
    // Navigate to the application
    await page.goto('/');
    
    // Wait for page to be ready
    await page.waitForLoadState('networkidle');
  });

  test('should display homepage with navigation', async ({ page }) => {
    // Check page title
    await expect(page).toHaveTitle(/財務計画計算機/);
    
    // Check navigation is present
    const nav = page.locator('nav');
    await expect(nav).toBeVisible();
    
    // Check main navigation links
    await expect(page.getByRole('link', { name: /ダッシュボード/ })).toBeVisible();
    await expect(page.getByRole('link', { name: /財務データ/ })).toBeVisible();
    await expect(page.getByRole('link', { name: /計算/ })).toBeVisible();
    await expect(page.getByRole('link', { name: /目標/ })).toBeVisible();
  });

  test('should navigate to financial data page', async ({ page }) => {
    // Click on financial data link
    await page.getByRole('link', { name: /財務データ/ }).click();
    
    // Wait for navigation
    await page.waitForURL('**/financial-data');
    
    // Check page heading
    await expect(page.getByRole('heading', { name: /財務情報/ })).toBeVisible();
  });

  test('should input financial profile data', async ({ page }) => {
    // Navigate to financial data page
    await page.goto('/financial-data');
    
    // Fill in monthly income
    const incomeInput = page.getByLabel(/月収/);
    await incomeInput.fill('400000');
    
    // Fill in monthly expenses
    const expensesInput = page.getByLabel(/月間支出/);
    await expensesInput.fill('280000');
    
    // Fill in current savings
    const savingsInput = page.getByLabel(/現在の貯蓄/);
    await savingsInput.fill('1500000');
    
    // Fill in investment return
    const returnInput = page.getByLabel(/投資利回り/);
    await returnInput.fill('5');
    
    // Fill in inflation rate
    const inflationInput = page.getByLabel(/インフレ率/);
    await inflationInput.fill('2');
    
    // Submit form
    await page.getByRole('button', { name: /保存/ }).click();
    
    // Wait for success message
    await expect(page.getByText(/保存しました/)).toBeVisible({ timeout: 5000 });
  });

  test('should calculate asset projection', async ({ page }) => {
    // Navigate to calculations page
    await page.goto('/calculations');
    
    // Wait for page to load
    await page.waitForLoadState('networkidle');
    
    // Fill in calculation parameters
    await page.getByLabel(/月収/).fill('400000');
    await page.getByLabel(/月間支出/).fill('280000');
    await page.getByLabel(/現在の貯蓄/).fill('1500000');
    await page.getByLabel(/投資利回り/).fill('5');
    await page.getByLabel(/予測期間/).fill('30');
    
    // Click calculate button
    await page.getByRole('button', { name: /計算/ }).click();
    
    // Wait for results
    await page.waitForSelector('[data-testid="asset-projection-chart"]', { timeout: 10000 });
    
    // Check that chart is displayed
    const chart = page.locator('[data-testid="asset-projection-chart"]');
    await expect(chart).toBeVisible();
    
    // Check that results summary is displayed
    await expect(page.getByText(/年後の予想資産/)).toBeVisible();
  });

  test('should calculate retirement needs', async ({ page }) => {
    // Navigate to calculations page
    await page.goto('/calculations');
    
    // Switch to retirement calculator tab
    await page.getByRole('tab', { name: /老後資金/ }).click();
    
    // Fill in retirement data
    await page.getByLabel(/現在の年齢/).fill('35');
    await page.getByLabel(/退職年齢/).fill('65');
    await page.getByLabel(/平均寿命/).fill('90');
    await page.getByLabel(/老後の月間生活費/).fill('250000');
    await page.getByLabel(/年金受給額/).fill('150000');
    
    // Calculate
    await page.getByRole('button', { name: /計算/ }).click();
    
    // Wait for results
    await expect(page.getByText(/必要老後資金/)).toBeVisible({ timeout: 10000 });
    await expect(page.getByText(/充足率/)).toBeVisible();
  });

  test('should calculate emergency fund', async ({ page }) => {
    // Navigate to calculations page
    await page.goto('/calculations');
    
    // Switch to emergency fund calculator tab
    await page.getByRole('tab', { name: /緊急資金/ }).click();
    
    // Fill in emergency fund data
    await page.getByLabel(/月間支出/).fill('280000');
    await page.getByLabel(/緊急時期間/).fill('6');
    await page.getByLabel(/現在の貯蓄/).fill('1500000');
    
    // Calculate
    await page.getByRole('button', { name: /計算/ }).click();
    
    // Wait for results
    await expect(page.getByText(/必要緊急資金/)).toBeVisible({ timeout: 10000 });
    await expect(page.getByText(/充足状況/)).toBeVisible();
  });

  test('should handle validation errors', async ({ page }) => {
    // Navigate to financial data page
    await page.goto('/financial-data');
    
    // Try to submit with invalid data
    await page.getByLabel(/月収/).fill('-1000');
    await page.getByRole('button', { name: /保存/ }).click();
    
    // Check for error message
    await expect(page.getByText(/正の値を入力してください/)).toBeVisible();
  });

  test('should handle API errors gracefully', async ({ page }) => {
    // Intercept API call and return error
    await page.route('**/api/financial-data', route => {
      route.fulfill({
        status: 500,
        body: JSON.stringify({ error: 'Internal Server Error' }),
      });
    });
    
    // Navigate to financial data page
    await page.goto('/financial-data');
    
    // Try to submit form
    await page.getByLabel(/月収/).fill('400000');
    await page.getByRole('button', { name: /保存/ }).click();
    
    // Check for error message
    await expect(page.getByText(/エラーが発生しました/)).toBeVisible({ timeout: 5000 });
  });

  test('should persist data across page reloads', async ({ page }) => {
    // Navigate to financial data page
    await page.goto('/financial-data');
    
    // Fill in data
    await page.getByLabel(/月収/).fill('400000');
    await page.getByRole('button', { name: /保存/ }).click();
    
    // Wait for save
    await expect(page.getByText(/保存しました/)).toBeVisible({ timeout: 5000 });
    
    // Reload page
    await page.reload();
    
    // Check that data is still there
    const incomeInput = page.getByLabel(/月収/);
    await expect(incomeInput).toHaveValue('400000');
  });
});
