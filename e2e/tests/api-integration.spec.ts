import { test, expect } from '@playwright/test';

/**
 * E2E Test: API Integration
 * 
 * Tests the integration between frontend and backend APIs
 */
test.describe('API Integration', () => {
  const apiBaseUrl = process.env.API_URL || 'http://localhost:8080/api';

  test('should check API health', async ({ request }) => {
    const response = await request.get('http://localhost:8080/health');

    expect(response.ok()).toBeTruthy();

    const data = await response.json();
    if (!response.ok()) {
      console.log('API Error Response (Emergency Fund):', await response.text());
    }
    expect(response.ok()).toBeTruthy();
  });

  test('should calculate asset projection via API', async ({ request }) => {
    const response = await request.post(`${apiBaseUrl}/calculations/asset-projection`, {
      data: {
        user_id: 'test-user',
        monthly_income: 400000,
        monthly_expenses: 280000,
        current_savings: 1500000,
        investment_return: 5.0,
        inflation_rate: 2.0,
        projection_years: 30,
      },
    });

    if (!response.ok()) {
      console.log('API Error Response:', await response.text());
    }
    expect(response.ok()).toBeTruthy();

    const data = await response.json();
    expect(data.projections).toBeDefined();
    expect(Array.isArray(data.projections)).toBeTruthy();
    expect(data.projections.length).toBeGreaterThan(0);
  });

  test('should calculate retirement needs via API', async ({ request }) => {
    const response = await request.post(`${apiBaseUrl}/calculations/retirement`, {
      data: {
        user_id: 'test-user',
        current_age: 35,
        retirement_age: 65,
        life_expectancy: 90,
        monthly_retirement_expenses: 250000,
        pension_amount: 150000,
        current_savings: 1500000,
        monthly_savings: 120000,
        investment_return: 5.0,
        inflation_rate: 2.0,
      },
    });

    if (!response.ok()) {
      console.log('API Error Response (Retirement):', await response.text());
    }
    expect(response.ok()).toBeTruthy();

    const data = await response.json();
    expect(data.required_amount).toBeDefined();
    expect(data.projected_amount).toBeDefined();
    expect(data.sufficiency_rate).toBeDefined();
  });

  test('should calculate emergency fund via API', async ({ request }) => {
    const response = await request.post(`${apiBaseUrl}/calculations/emergency-fund`, {
      data: {
        user_id: 'test-user',
        monthly_expenses: 280000,
        emergency_months: 6,
        current_savings: 1500000,
      },
    });

    expect(response.ok()).toBeTruthy();

    const data = await response.json();
    expect(data.required_amount).toBeDefined();
    expect(data.current_amount).toBeDefined();
    expect(data.shortfall).toBeDefined();
  });

  test('should handle API validation errors', async ({ request }) => {
    const response = await request.post(`${apiBaseUrl}/calculations/asset-projection`, {
      data: {
        user_id: 'test-user',
        monthly_income: -1000, // Invalid negative value
        monthly_expenses: 280000,
        current_savings: 1500000,
        investment_return: 5.0,
        inflation_rate: 2.0,
        projection_years: 30,
      },
    });

    expect(response.status()).toBe(400);

    const data = await response.json();
    expect(data.error).toBeDefined();
  });



  test('should handle CORS correctly', async ({ request }) => {
    const response = await request.get(`${apiBaseUrl}/`, {
      headers: {
        'Origin': 'http://localhost:3000',
      },
    });

    const headers = response.headers();
    expect(headers['access-control-allow-origin']).toBeDefined();
  });

  test('should return consistent response format', async ({ request }) => {
    const response = await request.post(`${apiBaseUrl}/calculations/asset-projection`, {
      data: {
        user_id: 'test-user',
        monthly_income: 400000,
        monthly_expenses: 280000,
        current_savings: 1500000,
        investment_return: 5.0,
        inflation_rate: 2.0,
        projection_years: 30,
      },
    });

    expect(response.ok()).toBeTruthy();

    const data = await response.json();

    // Check response structure
    expect(data).toHaveProperty('projections');
    expect(data).toHaveProperty('summary');

    // Check projection structure
    const projection = data.projections[0];
    expect(projection).toHaveProperty('year');
    expect(projection).toHaveProperty('total_assets');
    expect(projection).toHaveProperty('contributed_amount');
    expect(projection).toHaveProperty('investment_gains');
  });




});
