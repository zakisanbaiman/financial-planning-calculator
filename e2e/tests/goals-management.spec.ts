import { test, expect } from '@playwright/test';

/**
 * E2E Test: Goals Management Flow
 * 
 * Tests the complete flow of creating, updating, and tracking financial goals
 */
test.describe('Goals Management Flow', () => {
  const testUserId = 'e2e-test-user-' + Date.now();
  const testGoalTitle = 'E2E Test Goal ' + Date.now();

  test.beforeEach(async ({ page }) => {
    await page.goto('/');
    await page.waitForLoadState('networkidle');
  });

  test('should navigate to goals page', async ({ page }) => {
    // Click on goals link
    await page.getByRole('link', { name: /目標/ }).click();
    
    // Wait for navigation
    await page.waitForURL('**/goals');
    
    // Check page heading
    await expect(page.getByRole('heading', { name: /目標/ })).toBeVisible();
  });

  test('should create a new goal', async ({ page }) => {
    // Navigate to goals page
    await page.goto('/goals');
    
    // Click create goal button
    await page.getByRole('button', { name: /新しい目標/ }).click();
    
    // Fill in goal form
    await page.getByLabel(/目標名/).fill(testGoalTitle);
    await page.getByLabel(/目標金額/).fill('5000000');
    await page.getByLabel(/目標期日/).fill('2025-12-31');
    await page.getByLabel(/月間積立額/).fill('100000');
    
    // Select goal type
    await page.getByLabel(/目標種別/).selectOption('savings');
    
    // Submit form
    await page.getByRole('button', { name: /作成/ }).click();
    
    // Wait for success message
    await expect(page.getByText(/目標を作成しました/)).toBeVisible({ timeout: 5000 });
    
    // Check that goal appears in list
    await expect(page.getByText(testGoalTitle)).toBeVisible();
  });

  test('should display goal progress', async ({ page }) => {
    // Navigate to goals page
    await page.goto('/goals');
    
    // Wait for goals to load
    await page.waitForSelector('[data-testid="goal-item"]', { timeout: 10000 });
    
    // Check that progress bar is visible
    const progressBar = page.locator('[data-testid="goal-progress"]').first();
    await expect(progressBar).toBeVisible();
    
    // Check that progress percentage is displayed
    await expect(page.getByText(/%/)).toBeVisible();
  });

  test('should update goal progress', async ({ page }) => {
    // Navigate to goals page
    await page.goto('/goals');
    
    // Click on first goal
    await page.locator('[data-testid="goal-item"]').first().click();
    
    // Wait for goal detail page
    await page.waitForURL('**/goals/*');
    
    // Click update progress button
    await page.getByRole('button', { name: /進捗を更新/ }).click();
    
    // Fill in current amount
    await page.getByLabel(/現在の金額/).fill('1000000');
    
    // Submit
    await page.getByRole('button', { name: /更新/ }).click();
    
    // Wait for success message
    await expect(page.getByText(/進捗を更新しました/)).toBeVisible({ timeout: 5000 });
  });

  test('should view goal recommendations', async ({ page }) => {
    // Navigate to goals page
    await page.goto('/goals');
    
    // Click on first goal
    await page.locator('[data-testid="goal-item"]').first().click();
    
    // Wait for recommendations section
    await expect(page.getByText(/推奨事項/)).toBeVisible({ timeout: 10000 });
    
    // Check that recommendations are displayed
    await expect(page.getByText(/月間貯蓄額/)).toBeVisible();
  });

  test('should edit goal', async ({ page }) => {
    // Navigate to goals page
    await page.goto('/goals');
    
    // Click on first goal
    await page.locator('[data-testid="goal-item"]').first().click();
    
    // Click edit button
    await page.getByRole('button', { name: /編集/ }).click();
    
    // Update goal amount
    await page.getByLabel(/目標金額/).fill('6000000');
    
    // Submit
    await page.getByRole('button', { name: /保存/ }).click();
    
    // Wait for success message
    await expect(page.getByText(/目標を更新しました/)).toBeVisible({ timeout: 5000 });
  });

  test('should delete goal', async ({ page }) => {
    // Navigate to goals page
    await page.goto('/goals');
    
    // Get initial goal count
    const initialCount = await page.locator('[data-testid="goal-item"]').count();
    
    // Click on first goal
    await page.locator('[data-testid="goal-item"]').first().click();
    
    // Click delete button
    await page.getByRole('button', { name: /削除/ }).click();
    
    // Confirm deletion
    await page.getByRole('button', { name: /確認/ }).click();
    
    // Wait for success message
    await expect(page.getByText(/目標を削除しました/)).toBeVisible({ timeout: 5000 });
    
    // Check that goal count decreased
    await page.goto('/goals');
    const newCount = await page.locator('[data-testid="goal-item"]').count();
    expect(newCount).toBe(initialCount - 1);
  });

  test('should filter goals by status', async ({ page }) => {
    // Navigate to goals page
    await page.goto('/goals');
    
    // Click active filter
    await page.getByRole('button', { name: /アクティブ/ }).click();
    
    // Wait for filtered results
    await page.waitForTimeout(1000);
    
    // Check that only active goals are shown
    const goals = page.locator('[data-testid="goal-item"]');
    const count = await goals.count();
    
    for (let i = 0; i < count; i++) {
      const goal = goals.nth(i);
      await expect(goal.locator('[data-testid="goal-status"]')).toHaveText(/アクティブ/);
    }
  });

  test('should display goals summary chart', async ({ page }) => {
    // Navigate to goals page
    await page.goto('/goals');
    
    // Wait for chart to load
    await page.waitForSelector('[data-testid="goals-summary-chart"]', { timeout: 10000 });
    
    // Check that chart is visible
    const chart = page.locator('[data-testid="goals-summary-chart"]');
    await expect(chart).toBeVisible();
  });

  test('should validate goal form inputs', async ({ page }) => {
    // Navigate to goals page
    await page.goto('/goals');
    
    // Click create goal button
    await page.getByRole('button', { name: /新しい目標/ }).click();
    
    // Try to submit with invalid data
    await page.getByLabel(/目標金額/).fill('-1000');
    await page.getByRole('button', { name: /作成/ }).click();
    
    // Check for validation error
    await expect(page.getByText(/正の値を入力してください/)).toBeVisible();
  });

  test('should handle concurrent goal updates', async ({ page, context }) => {
    // Open two tabs
    const page2 = await context.newPage();
    
    // Navigate both to goals page
    await page.goto('/goals');
    await page2.goto('/goals');
    
    // Click on same goal in both tabs
    await page.locator('[data-testid="goal-item"]').first().click();
    await page2.locator('[data-testid="goal-item"]').first().click();
    
    // Update progress in first tab
    await page.getByRole('button', { name: /進捗を更新/ }).click();
    await page.getByLabel(/現在の金額/).fill('1000000');
    await page.getByRole('button', { name: /更新/ }).click();
    
    // Wait for update
    await expect(page.getByText(/進捗を更新しました/)).toBeVisible({ timeout: 5000 });
    
    // Refresh second tab and check updated value
    await page2.reload();
    await expect(page2.getByText(/1,000,000/)).toBeVisible({ timeout: 5000 });
    
    await page2.close();
  });
});
