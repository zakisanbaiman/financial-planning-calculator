import { test, expect } from '@playwright/test';
import { generateTestUserId, createFinancialData, addYearsToDate, API_BASE_URL } from './test-utils';

/**
 * E2E Test: Goals Management Scenarios
 * 
 * Comprehensive scenario tests for goal creation, management, and related operations
 */

test.describe('Goals Management Scenarios', () => {
  test('Scenario: Create financial data and then create a goal', async ({ request }) => {
    const userId = generateTestUserId();

    // Step 1: Create financial data first
    const financialData = await createFinancialData(request, userId);
    expect(financialData.user_id).toBe(userId);

    // Step 2: Create a savings goal
    const goalResponse = await request.post(`${API_BASE_URL}/api/goals`, {
      data: {
        user_id: userId,
        goal_type: 'savings',
        title: 'マイホーム購入資金',
        target_amount: 10000000,
        target_date: addYearsToDate(3),
        current_amount: 1000000,
        monthly_contribution: 150000,
        description: '3年後にマイホームを購入するための貯蓄',
      },
    });

    expect(goalResponse.status()).toBe(201);
    const goalData = await goalResponse.json();
    expect(goalData.user_id).toBe(userId);
    expect(goalData.goal_id).toBeDefined();
  });

  test('Scenario: Create multiple goals and retrieve them', async ({ request }) => {
    const userId = generateTestUserId();

    // Create financial data
    await createFinancialData(request, userId);

    // Create multiple goals
    const goals = [
      {
        goal_type: 'savings',
        title: 'マイホーム購入',
        target_amount: 10000000,
        current_amount: 1000000,
        monthly_contribution: 150000,
      },
      {
        goal_type: 'retirement',
        title: '老後資金',
        target_amount: 30000000,
        current_amount: 5000000,
        monthly_contribution: 100000,
      },
      {
        goal_type: 'emergency',
        title: '緊急資金',
        target_amount: 3000000,
        current_amount: 500000,
        monthly_contribution: 50000,
      },
    ];

    const createdGoals = [];
    for (const goal of goals) {
      const response = await request.post('${API_BASE_URL}/api/goals', {
        data: {
          user_id: userId,
          ...goal,
          target_date: addYearsToDate(5),
        },
      });
      expect(response.status()).toBe(201);
      createdGoals.push(await response.json());
    }

    // Retrieve all goals
    const getGoalsResponse = await request.get(`http://localhost:8080/api/goals?user_id=${userId}`);
    expect(getGoalsResponse.ok()).toBeTruthy();
    
    const goalsData = await getGoalsResponse.json();
    expect(goalsData.goals).toBeDefined();
    expect(goalsData.goals.length).toBe(3);
  });

  test('Scenario: Update a goal', async ({ request }) => {
    const userId = generateTestUserId();

    // Create financial data
    await createFinancialData(request, userId);

    // Create a goal
    const createResponse = await request.post('${API_BASE_URL}/api/goals', {
      data: {
        user_id: userId,
        goal_type: 'savings',
        title: '車購入資金',
        target_amount: 3000000,
        target_date: addYearsToDate(2),
        current_amount: 500000,
        monthly_contribution: 80000,
      },
    });

    expect(createResponse.status()).toBe(201);
    const goalData = await createResponse.json();

    // Update the goal
    const updateResponse = await request.put(
      `http://localhost:8080/api/goals/${goalData.goal_id}?user_id=${userId}`,
      {
        data: {
          title: '新車購入資金（更新）',
          target_amount: 3500000,
          monthly_contribution: 100000,
        },
      }
    );

    expect(updateResponse.ok()).toBeTruthy();
    const updatedData = await updateResponse.json();
    expect(updatedData.success).toBe(true);
  });

  test('Scenario: Update goal progress', async ({ request }) => {
    const userId = generateTestUserId();

    // Create financial data
    await createFinancialData(request, userId);

    // Create a goal
    const createResponse = await request.post('${API_BASE_URL}/api/goals', {
      data: {
        user_id: userId,
        goal_type: 'savings',
        title: '旅行資金',
        target_amount: 1000000,
        target_date: new Date(Date.now() + 365 * 24 * 60 * 60 * 1000).toISOString(),
        current_amount: 300000,
        monthly_contribution: 50000,
      },
    });

    const goalData = await createResponse.json();

    // Update progress
    const progressResponse = await request.put(
      `http://localhost:8080/api/goals/${goalData.goal_id}/progress?user_id=${userId}`,
      {
        data: {
          current_amount: 450000,
          note: '今月も順調に積立できました',
        },
      }
    );

    expect(progressResponse.ok()).toBeTruthy();
    const progressData = await progressResponse.json();
    expect(progressData.success).toBe(true);
  });

  test('Scenario: Get goal recommendations', async ({ request }) => {
    const userId = generateTestUserId();

    // Create financial data
    await createFinancialData(request, userId);

    // Create a goal
    const createResponse = await request.post('${API_BASE_URL}/api/goals', {
      data: {
        user_id: userId,
        goal_type: 'savings',
        title: '教育資金',
        target_amount: 5000000,
        target_date: addYearsToDate(10),
        current_amount: 500000,
        monthly_contribution: 40000,
      },
    });

    const goalData = await createResponse.json();

    // Get recommendations
    const recommendationsResponse = await request.get(
      `http://localhost:8080/api/goals/${goalData.goal_id}/recommendations?user_id=${userId}`
    );

    expect(recommendationsResponse.ok()).toBeTruthy();
    const recommendations = await recommendationsResponse.json();
    expect(recommendations.recommendations).toBeDefined();
  });

  test('Scenario: Analyze goal feasibility', async ({ request }) => {
    const userId = generateTestUserId();

    // Create financial data
    await createFinancialData(request, userId);

    // Create a goal
    const createResponse = await request.post('${API_BASE_URL}/api/goals', {
      data: {
        user_id: userId,
        goal_type: 'savings',
        title: '投資用不動産',
        target_amount: 20000000,
        target_date: addYearsToDate(5),
        current_amount: 2000000,
        monthly_contribution: 200000,
      },
    });

    const goalData = await createResponse.json();

    // Analyze feasibility
    const feasibilityResponse = await request.get(
      `http://localhost:8080/api/goals/${goalData.goal_id}/feasibility?user_id=${userId}`
    );

    expect(feasibilityResponse.ok()).toBeTruthy();
    const feasibility = await feasibilityResponse.json();
    expect(feasibility.achievable).toBeDefined();
    expect(feasibility.risk_level).toBeDefined();
  });

  test('Scenario: Delete a goal', async ({ request }) => {
    const userId = generateTestUserId();

    // Create financial data
    await createFinancialData(request, userId);

    // Create a goal
    const createResponse = await request.post('${API_BASE_URL}/api/goals', {
      data: {
        user_id: userId,
        goal_type: 'savings',
        title: 'テスト目標',
        target_amount: 1000000,
        target_date: new Date(Date.now() + 365 * 24 * 60 * 60 * 1000).toISOString(),
        current_amount: 0,
        monthly_contribution: 50000,
      },
    });

    const goalData = await createResponse.json();

    // Delete the goal
    const deleteResponse = await request.delete(
      `http://localhost:8080/api/goals/${goalData.goal_id}?user_id=${userId}`
    );

    expect(deleteResponse.status()).toBe(204);

    // Verify the goal is deleted
    const getResponse = await request.get(
      `http://localhost:8080/api/goals/${goalData.goal_id}?user_id=${userId}`
    );
    expect(getResponse.status()).toBe(404);
  });

  test('Scenario: Filter goals by type', async ({ request }) => {
    const userId = generateTestUserId();

    // Create financial data
    await createFinancialData(request, userId);

    // Create goals of different types
    const goalTypes = ['savings', 'retirement', 'emergency', 'savings'];
    for (const goalType of goalTypes) {
      await request.post('${API_BASE_URL}/api/goals', {
        data: {
          user_id: userId,
          goal_type: goalType,
          title: `${goalType} 目標`,
          target_amount: 1000000,
          target_date: new Date(Date.now() + 365 * 24 * 60 * 60 * 1000).toISOString(),
          current_amount: 100000,
          monthly_contribution: 50000,
        },
      });
    }

    // Get only savings goals
    const savingsResponse = await request.get(
      `http://localhost:8080/api/goals?user_id=${userId}&goal_type=savings`
    );
    expect(savingsResponse.ok()).toBeTruthy();
    const savingsData = await savingsResponse.json();
    expect(savingsData.goals.length).toBe(2);
  });

  test('Scenario: Error - Create goal without financial data', async ({ request }) => {
    const userId = generateTestUserId();

    // Try to create a goal without creating financial data first
    const goalResponse = await request.post('${API_BASE_URL}/api/goals', {
      data: {
        user_id: userId,
        goal_type: 'savings',
        title: 'テスト目標',
        target_amount: 1000000,
        target_date: new Date(Date.now() + 365 * 24 * 60 * 60 * 1000).toISOString(),
        current_amount: 0,
        monthly_contribution: 50000,
      },
    });

    // Should return 400 Bad Request due to missing financial data
    expect(goalResponse.status()).toBe(400);
  });

  test('Scenario: Error - Create goal with invalid data', async ({ request }) => {
    const userId = generateTestUserId();

    // Create financial data
    await createFinancialData(request, userId);

    // Try to create a goal with negative target amount
    const goalResponse = await request.post('${API_BASE_URL}/api/goals', {
      data: {
        user_id: userId,
        goal_type: 'savings',
        title: 'テスト目標',
        target_amount: -1000000, // Invalid: negative amount
        target_date: new Date(Date.now() + 365 * 24 * 60 * 60 * 1000).toISOString(),
        current_amount: 0,
        monthly_contribution: 50000,
      },
    });

    // Should return 400 Bad Request due to validation error
    expect(goalResponse.status()).toBe(400);
  });

  test('Scenario: Complete financial planning flow', async ({ request }) => {
    const userId = generateTestUserId();

    // Step 1: Create financial data
    const financialData = await createFinancialData(request, userId);
    expect(financialData.user_id).toBe(userId);

    // Step 2: Create multiple goals
    const savingsGoal = await request.post('${API_BASE_URL}/api/goals', {
      data: {
        user_id: userId,
        goal_type: 'savings',
        title: 'マイホーム購入',
        target_amount: 10000000,
        target_date: addYearsToDate(5),
        current_amount: 2000000,
        monthly_contribution: 150000,
      },
    });
    expect(savingsGoal.status()).toBe(201);

    const retirementGoal = await request.post('${API_BASE_URL}/api/goals', {
      data: {
        user_id: userId,
        goal_type: 'retirement',
        title: '老後資金',
        target_amount: 50000000,
        target_date: addYearsToDate(30),
        current_amount: 5000000,
        monthly_contribution: 100000,
      },
    });
    expect(retirementGoal.status()).toBe(201);

    // Step 3: Get all goals
    const allGoalsResponse = await request.get(`http://localhost:8080/api/goals?user_id=${userId}`);
    expect(allGoalsResponse.ok()).toBeTruthy();
    const allGoals = await allGoalsResponse.json();
    expect(allGoals.goals.length).toBeGreaterThanOrEqual(2);

    // Step 4: Calculate asset projection
    const projectionResponse = await request.post('${API_BASE_URL}/api/calculations/asset-projection', {
      data: {
        user_id: userId,
        years: 10,
      },
    });
    expect(projectionResponse.ok()).toBeTruthy();

    // Step 5: Update retirement data
    const retirementDataResponse = await request.put(`http://localhost:8080/api/financial-data/${userId}/retirement`, {
      data: {
        current_age: 35,
        retirement_age: 65,
        life_expectancy: 90,
        monthly_expenses_after_retirement: 250000,
        expected_pension: 150000,
      },
    });
    expect(retirementDataResponse.ok()).toBeTruthy();

    // Step 6: Calculate retirement projection
    const retirementProjectionResponse = await request.post('${API_BASE_URL}/api/calculations/retirement', {
      data: {
        user_id: userId,
      },
    });
    expect(retirementProjectionResponse.ok()).toBeTruthy();
  });
});
