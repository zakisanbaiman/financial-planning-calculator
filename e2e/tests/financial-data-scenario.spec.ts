import { test, expect } from '@playwright/test';
import { generateTestUserId, addYearsToDate, API_BASE_URL } from './test-utils';

/**
 * E2E Test: Financial Data Flow Scenarios
 * 
 * Tests the complete flow of financial data creation, updates, and related operations
 */

test.describe('Financial Data Flow Scenarios', () => {
  test('Scenario: Create complete financial profile', async ({ request }) => {
    const userId = generateTestUserId();

    // Create financial data with complete profile
    const response = await request.post('${API_BASE_URL}/api/financial-data', {
      data: {
        user_id: userId,
        monthly_income: 600000,
        monthly_expenses: [
          { category: '住居費', amount: 150000 },
          { category: '食費', amount: 80000 },
          { category: '光熱費', amount: 25000 },
          { category: '通信費', amount: 15000 },
          { category: '交通費', amount: 20000 },
          { category: '保険', amount: 30000 },
          { category: 'その他', amount: 50000 },
        ],
        current_savings: [
          { type: 'deposit', amount: 3000000 },
          { type: 'investment', amount: 2000000 },
          { type: 'other', amount: 500000 },
        ],
        investment_return: 5.0,
        inflation_rate: 2.0,
      },
    });

    expect(response.status()).toBe(201);
    const data = await response.json();
    expect(data.user_id).toBe(userId);
    expect(data.plan_id).toBeDefined();
  });

  test('Scenario: Update financial profile', async ({ request }) => {
    const userId = generateTestUserId();

    // Create initial financial data
    const createResponse = await request.post('${API_BASE_URL}/api/financial-data', {
      data: {
        user_id: userId,
        monthly_income: 500000,
        monthly_expenses: [{ category: '住居費', amount: 120000 }],
        current_savings: [{ type: 'deposit', amount: 1000000 }],
        investment_return: 5.0,
        inflation_rate: 2.0,
      },
    });
    expect(createResponse.status()).toBe(201);

    // Update the profile
    const updateResponse = await request.put(`${API_BASE_URL}/api/financial-data/${userId}/profile`, {
      data: {
        monthly_income: 550000,
        monthly_expenses: [
          { category: '住居費', amount: 130000 },
          { category: '食費', amount: 60000 },
        ],
        current_savings: [
          { type: 'deposit', amount: 1200000 },
          { type: 'investment', amount: 500000 },
        ],
        investment_return: 6.0,
        inflation_rate: 2.5,
      },
    });

    expect(updateResponse.ok()).toBeTruthy();
    const updateData = await updateResponse.json();
    expect(updateData.success).toBe(true);
  });

  test('Scenario: Add and update retirement data', async ({ request }) => {
    const userId = generateTestUserId();

    // Create financial data
    await request.post('${API_BASE_URL}/api/financial-data', {
      data: {
        user_id: userId,
        monthly_income: 500000,
        monthly_expenses: [{ category: '住居費', amount: 120000 }],
        current_savings: [{ type: 'deposit', amount: 2000000 }],
        investment_return: 5.0,
        inflation_rate: 2.0,
      },
    });

    // Add retirement data
    const retirementResponse = await request.put(`${API_BASE_URL}/api/financial-data/${userId}/retirement`, {
      data: {
        current_age: 35,
        retirement_age: 65,
        life_expectancy: 90,
        monthly_expenses_after_retirement: 250000,
        expected_pension: 150000,
      },
    });

    expect(retirementResponse.ok()).toBeTruthy();
    const retirementData = await retirementResponse.json();
    expect(retirementData.success).toBe(true);
  });

  test('Scenario: Add emergency fund settings', async ({ request }) => {
    const userId = generateTestUserId();

    // Create financial data
    await request.post('${API_BASE_URL}/api/financial-data', {
      data: {
        user_id: userId,
        monthly_income: 500000,
        monthly_expenses: [{ category: '住居費', amount: 120000 }],
        current_savings: [{ type: 'deposit', amount: 2000000 }],
        investment_return: 5.0,
        inflation_rate: 2.0,
      },
    });

    // Add emergency fund
    const emergencyResponse = await request.put(`${API_BASE_URL}/api/financial-data/${userId}/emergency-fund`, {
      data: {
        target_months: 6,
        current_amount: 500000,
        priority: 'high',
      },
    });

    expect(emergencyResponse.ok()).toBeTruthy();
    const emergencyData = await emergencyResponse.json();
    expect(emergencyData.success).toBe(true);
  });

  test('Scenario: Retrieve financial data', async ({ request }) => {
    const userId = generateTestUserId();

    // Create financial data
    await request.post('${API_BASE_URL}/api/financial-data', {
      data: {
        user_id: userId,
        monthly_income: 500000,
        monthly_expenses: [{ category: '住居費', amount: 120000 }],
        current_savings: [{ type: 'deposit', amount: 2000000 }],
        investment_return: 5.0,
        inflation_rate: 2.0,
      },
    });

    // Retrieve the data
    const getResponse = await request.get(`${API_BASE_URL}/api/financial-data?user_id=${userId}`);
    expect(getResponse.ok()).toBeTruthy();
    
    const data = await getResponse.json();
    expect(data.plan).toBeDefined();
  });

  test('Scenario: Complete financial setup with goals', async ({ request }) => {
    const userId = generateTestUserId();

    // Step 1: Create financial data
    const financialResponse = await request.post('${API_BASE_URL}/api/financial-data', {
      data: {
        user_id: userId,
        monthly_income: 600000,
        monthly_expenses: [
          { category: '住居費', amount: 150000 },
          { category: '食費', amount: 70000 },
          { category: '光熱費', amount: 25000 },
        ],
        current_savings: [
          { type: 'deposit', amount: 3000000 },
          { type: 'investment', amount: 1500000 },
        ],
        investment_return: 5.5,
        inflation_rate: 2.0,
      },
    });
    expect(financialResponse.status()).toBe(201);

    // Step 2: Add retirement data
    const retirementResponse = await request.put(`${API_BASE_URL}/api/financial-data/${userId}/retirement`, {
      data: {
        current_age: 30,
        retirement_age: 60,
        life_expectancy: 85,
        monthly_expenses_after_retirement: 200000,
        expected_pension: 120000,
      },
    });
    expect(retirementResponse.ok()).toBeTruthy();

    // Step 3: Add emergency fund
    const emergencyResponse = await request.put(`${API_BASE_URL}/api/financial-data/${userId}/emergency-fund`, {
      data: {
        target_months: 6,
        current_amount: 1000000,
        priority: 'high',
      },
    });
    expect(emergencyResponse.ok()).toBeTruthy();

    // Step 4: Create goals
    const goal1 = await request.post('${API_BASE_URL}/api/goals', {
      data: {
        user_id: userId,
        goal_type: 'savings',
        title: 'マイホーム頭金',
        target_amount: 5000000,
        target_date: addYearsToDate(3),
        current_amount: 1000000,
        monthly_contribution: 100000,
      },
    });
    expect(goal1.status()).toBe(201);

    const goal2 = await request.post('${API_BASE_URL}/api/goals', {
      data: {
        user_id: userId,
        goal_type: 'retirement',
        title: '老後資金',
        target_amount: 30000000,
        target_date: addYearsToDate(30),
        current_amount: 3000000,
        monthly_contribution: 80000,
      },
    });
    expect(goal2.status()).toBe(201);

    // Step 5: Verify all data is retrievable
    const finalDataResponse = await request.get(`${API_BASE_URL}/api/financial-data?user_id=${userId}`);
    expect(finalDataResponse.ok()).toBeTruthy();

    const goalsResponse = await request.get(`${API_BASE_URL}/api/goals?user_id=${userId}`);
    expect(goalsResponse.ok()).toBeTruthy();
    const goalsData = await goalsResponse.json();
    expect(goalsData.goals.length).toBe(2);
  });

  test('Scenario: Delete financial data', async ({ request }) => {
    const userId = generateTestUserId();

    // Create financial data
    await request.post('${API_BASE_URL}/api/financial-data', {
      data: {
        user_id: userId,
        monthly_income: 500000,
        monthly_expenses: [{ category: '住居費', amount: 120000 }],
        current_savings: [{ type: 'deposit', amount: 1000000 }],
        investment_return: 5.0,
        inflation_rate: 2.0,
      },
    });

    // Delete the data
    const deleteResponse = await request.delete(`${API_BASE_URL}/api/financial-data/${userId}`);
    expect(deleteResponse.status()).toBe(204);

    // Verify deletion
    const getResponse = await request.get(`${API_BASE_URL}/api/financial-data?user_id=${userId}`);
    expect(getResponse.status()).toBe(404);
  });

  test('Scenario: Error - Get non-existent financial data', async ({ request }) => {
    const userId = 'non-existent-user-12345';

    const response = await request.get(`${API_BASE_URL}/api/financial-data?user_id=${userId}`);
    expect(response.status()).toBe(404);
  });

  test('Scenario: Error - Invalid financial data', async ({ request }) => {
    const userId = generateTestUserId();

    // Try to create with negative income
    const response = await request.post('${API_BASE_URL}/api/financial-data', {
      data: {
        user_id: userId,
        monthly_income: -100000, // Invalid: negative
        monthly_expenses: [{ category: '住居費', amount: 120000 }],
        current_savings: [{ type: 'deposit', amount: 1000000 }],
        investment_return: 5.0,
        inflation_rate: 2.0,
      },
    });

    expect(response.status()).toBe(400);
  });
});
