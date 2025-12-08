import { test, expect } from '@playwright/test';

/**
 * E2E Test: Calculation Scenarios with Goals
 * 
 * Tests calculation endpoints in conjunction with goals and financial data
 */

// Test data helpers
const generateTestUserId = () => `test-user-${Date.now()}-${Math.random().toString(36).substring(7)}`;

const setupCompleteFinancialProfile = async (request: any, userId: string) => {
  // Create financial data
  await request.post('http://localhost:8080/api/financial-data', {
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
        { type: 'investment', amount: 2000000 },
      ],
      investment_return: 5.0,
      inflation_rate: 2.0,
    },
  });

  // Add retirement data
  await request.put(`http://localhost:8080/api/financial-data/${userId}/retirement`, {
    data: {
      current_age: 35,
      retirement_age: 65,
      life_expectancy: 90,
      monthly_expenses_after_retirement: 250000,
      expected_pension: 150000,
    },
  });

  // Add emergency fund
  await request.put(`http://localhost:8080/api/financial-data/${userId}/emergency-fund`, {
    data: {
      target_months: 6,
      current_amount: 800000,
      priority: 'high',
    },
  });
};

test.describe('Calculation Scenarios with Goals', () => {
  test('Scenario: Calculate asset projection', async ({ request }) => {
    const userId = generateTestUserId();
    await setupCompleteFinancialProfile(request, userId);

    const response = await request.post('http://localhost:8080/api/calculations/asset-projection', {
      data: {
        user_id: userId,
        years: 10,
      },
    });

    expect(response.ok()).toBeTruthy();
    const data = await response.json();
    expect(data.projections).toBeDefined();
    expect(data.summary).toBeDefined();
  });

  test('Scenario: Calculate retirement projection', async ({ request }) => {
    const userId = generateTestUserId();
    await setupCompleteFinancialProfile(request, userId);

    const response = await request.post('http://localhost:8080/api/calculations/retirement', {
      data: {
        user_id: userId,
      },
    });

    expect(response.ok()).toBeTruthy();
    const data = await response.json();
    expect(data.calculation).toBeDefined();
    expect(data.sufficiency_level).toBeDefined();
  });

  test('Scenario: Calculate emergency fund projection', async ({ request }) => {
    const userId = generateTestUserId();
    await setupCompleteFinancialProfile(request, userId);

    const response = await request.post('http://localhost:8080/api/calculations/emergency-fund', {
      data: {
        user_id: userId,
      },
    });

    expect(response.ok()).toBeTruthy();
    const data = await response.json();
    expect(data.status).toBeDefined();
    expect(data.priority).toBeDefined();
  });

  test('Scenario: Calculate comprehensive projection', async ({ request }) => {
    const userId = generateTestUserId();
    await setupCompleteFinancialProfile(request, userId);

    const response = await request.post('http://localhost:8080/api/calculations/comprehensive', {
      data: {
        user_id: userId,
        years: 15,
      },
    });

    expect(response.ok()).toBeTruthy();
    const data = await response.json();
    expect(data.asset_projection).toBeDefined();
    expect(data.retirement_analysis).toBeDefined();
    expect(data.emergency_fund_status).toBeDefined();
  });

  test('Scenario: Calculate goal projection', async ({ request }) => {
    const userId = generateTestUserId();
    await setupCompleteFinancialProfile(request, userId);

    // Create a goal
    const goalResponse = await request.post('http://localhost:8080/api/goals', {
      data: {
        user_id: userId,
        goal_type: 'savings',
        title: 'マイホーム購入',
        target_amount: 10000000,
        target_date: new Date(Date.now() + 365 * 24 * 60 * 60 * 1000 * 5).toISOString(),
        current_amount: 2000000,
        monthly_contribution: 150000,
      },
    });
    const goalData = await goalResponse.json();

    // Calculate goal projection
    const projectionResponse = await request.post('http://localhost:8080/api/calculations/goal-projection', {
      data: {
        user_id: userId,
        goal_id: goalData.goal_id,
      },
    });

    expect(projectionResponse.ok()).toBeTruthy();
    const projection = await projectionResponse.json();
    expect(projection.projection).toBeDefined();
  });

  test('Scenario: Multiple goals with comprehensive calculation', async ({ request }) => {
    const userId = generateTestUserId();
    await setupCompleteFinancialProfile(request, userId);

    // Create multiple goals
    await request.post('http://localhost:8080/api/goals', {
      data: {
        user_id: userId,
        goal_type: 'savings',
        title: 'マイホーム購入',
        target_amount: 10000000,
        target_date: new Date(Date.now() + 365 * 24 * 60 * 60 * 1000 * 5).toISOString(),
        current_amount: 2000000,
        monthly_contribution: 150000,
      },
    });

    await request.post('http://localhost:8080/api/goals', {
      data: {
        user_id: userId,
        goal_type: 'retirement',
        title: '老後資金',
        target_amount: 50000000,
        target_date: new Date(Date.now() + 365 * 24 * 60 * 60 * 1000 * 30).toISOString(),
        current_amount: 5000000,
        monthly_contribution: 100000,
      },
    });

    await request.post('http://localhost:8080/api/goals', {
      data: {
        user_id: userId,
        goal_type: 'emergency',
        title: '緊急資金',
        target_amount: 3000000,
        target_date: new Date(Date.now() + 365 * 24 * 60 * 60 * 1000 * 2).toISOString(),
        current_amount: 800000,
        monthly_contribution: 80000,
      },
    });

    // Calculate comprehensive projection with multiple goals
    const comprehensiveResponse = await request.post('http://localhost:8080/api/calculations/comprehensive', {
      data: {
        user_id: userId,
        years: 20,
      },
    });

    expect(comprehensiveResponse.ok()).toBeTruthy();
    const data = await comprehensiveResponse.json();
    expect(data.asset_projection).toBeDefined();
    expect(data.retirement_analysis).toBeDefined();
  });

  test('Scenario: Validate calculation with different timeframes', async ({ request }) => {
    const userId = generateTestUserId();
    await setupCompleteFinancialProfile(request, userId);

    // Test multiple timeframes
    const timeframes = [5, 10, 20, 30];
    for (const years of timeframes) {
      const response = await request.post('http://localhost:8080/api/calculations/asset-projection', {
        data: {
          user_id: userId,
          years: years,
        },
      });

      expect(response.ok()).toBeTruthy();
      const data = await response.json();
      expect(data.projections).toBeDefined();
    }
  });

  test('Scenario: Error - Calculate without financial data', async ({ request }) => {
    const userId = 'non-existent-user-12345';

    const response = await request.post('http://localhost:8080/api/calculations/asset-projection', {
      data: {
        user_id: userId,
        years: 10,
      },
    });

    // Should return an error status
    expect(response.status()).toBeGreaterThanOrEqual(400);
  });

  test('Scenario: Error - Invalid calculation parameters', async ({ request }) => {
    const userId = generateTestUserId();
    await setupCompleteFinancialProfile(request, userId);

    // Try with negative years
    const response = await request.post('http://localhost:8080/api/calculations/asset-projection', {
      data: {
        user_id: userId,
        years: -5,
      },
    });

    expect(response.status()).toBe(400);
  });

  test('Scenario: Complete planning with calculations and reports', async ({ request }) => {
    const userId = generateTestUserId();
    await setupCompleteFinancialProfile(request, userId);

    // Create goals
    await request.post('http://localhost:8080/api/goals', {
      data: {
        user_id: userId,
        goal_type: 'savings',
        title: 'マイホーム購入',
        target_amount: 10000000,
        target_date: new Date(Date.now() + 365 * 24 * 60 * 60 * 1000 * 5).toISOString(),
        current_amount: 2000000,
        monthly_contribution: 150000,
      },
    });

    // Calculate asset projection
    const assetProjection = await request.post('http://localhost:8080/api/calculations/asset-projection', {
      data: {
        user_id: userId,
        years: 10,
      },
    });
    expect(assetProjection.ok()).toBeTruthy();

    // Calculate retirement projection
    const retirementProjection = await request.post('http://localhost:8080/api/calculations/retirement', {
      data: {
        user_id: userId,
      },
    });
    expect(retirementProjection.ok()).toBeTruthy();

    // Calculate emergency fund
    const emergencyFund = await request.post('http://localhost:8080/api/calculations/emergency-fund', {
      data: {
        user_id: userId,
      },
    });
    expect(emergencyFund.ok()).toBeTruthy();

    // Generate financial summary report
    const summaryReport = await request.post('http://localhost:8080/api/reports/financial-summary', {
      data: {
        user_id: userId,
      },
    });
    expect(summaryReport.ok()).toBeTruthy();

    // Generate asset projection report
    const assetReport = await request.post('http://localhost:8080/api/reports/asset-projection', {
      data: {
        user_id: userId,
        years: 10,
      },
    });
    expect(assetReport.ok()).toBeTruthy();
  });
});
