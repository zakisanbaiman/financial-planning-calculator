import { APIRequestContext } from '@playwright/test';

/**
 * Shared test utilities for E2E tests
 */

// Configuration
export const API_BASE_URL = process.env.API_URL || 'http://localhost:8080';

export interface TestAuthCredentials {
  userId: string;
  token: string;
  email: string;
}

export const registerAndLoginTestUser = async (
  request: APIRequestContext
): Promise<TestAuthCredentials> => {
  const uniqueId = `${Date.now()}-${Math.random().toString(36).substring(7)}`;
  const email = `test-${uniqueId}@example.com`;
  const password = 'TestPass123!';
  const registerResponse = await request.post(`${API_BASE_URL}/api/auth/register`, {
    data: { email, password },
  });
  if (!registerResponse.ok()) {
    throw new Error(`Failed to register test user: ${registerResponse.status()}`);
  }
  const data = await registerResponse.json();
  return { userId: data.user_id, token: data.token, email: data.email };
};

export const authHeaders = (token: string) => ({
  Authorization: `Bearer ${token}`,
});

/**
 * Generate a unique test user ID
 */
export const generateTestUserId = () =>
  `test-user-${Date.now()}-${Math.random().toString(36).substring(7)}`;

/**
 * Add years to current date
 */
export const addYearsToDate = (years: number): string => {
  const date = new Date();
  date.setFullYear(date.getFullYear() + years);
  return date.toISOString();
};

/**
 * Create financial data for a test user
 */
export const createFinancialData = async (request: APIRequestContext, userId: string, token?: string) => {
  const response = await request.post(`${API_BASE_URL}/api/financial-data`, {
    headers: token ? { Authorization: `Bearer ${token}` } : undefined,
    data: {
      user_id: userId,
      monthly_income: 500000,
      monthly_expenses: [
        { category: '住居費', amount: 120000 },
        { category: '食費', amount: 60000 },
        { category: '光熱費', amount: 20000 },
      ],
      current_savings: [
        { type: 'deposit', amount: 2000000 },
        { type: 'investment', amount: 1000000 },
      ],
      investment_return: 5.0,
      inflation_rate: 2.0,
    },
  });

  if (!response.ok()) {
    throw new Error(`Failed to create financial data: ${response.status()}`);
  }

  return response.json();
};

/**
 * Setup complete financial profile including retirement data and emergency fund
 */
export const setupCompleteFinancialProfile = async (
  request: APIRequestContext,
  userId: string,
  token?: string
) => {
  const headers = token ? { Authorization: `Bearer ${token}` } : undefined;

  // Create financial data
  await request.post(`${API_BASE_URL}/api/financial-data`, {
    headers,
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
  await request.put(`${API_BASE_URL}/api/financial-data/${userId}/retirement`, {
    headers,
    data: {
      current_age: 35,
      retirement_age: 65,
      life_expectancy: 90,
      monthly_retirement_expenses: 250000,
      pension_amount: 150000,
    },
  });

  // Add emergency fund
  await request.put(`${API_BASE_URL}/api/financial-data/${userId}/emergency-fund`, {
    headers,
    data: {
      target_months: 6,
      current_amount: 800000,
      priority: 'high',
    },
  });
};
