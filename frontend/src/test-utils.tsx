import React, { ReactElement, ReactNode } from 'react';
import { render, RenderOptions, renderHook, RenderHookOptions } from '@testing-library/react';
import { AuthProvider } from '@/lib/contexts/AuthContext';
import { GuestModeProvider } from '@/lib/contexts/GuestModeContext';
import { ThemeProvider } from '@/lib/contexts/ThemeContext';
import { TutorialProvider } from '@/lib/contexts/TutorialContext';
import { FinancialDataProvider } from '@/lib/contexts/FinancialDataContext';
import { GoalsProvider } from '@/lib/contexts/GoalsContext';
import { CalculationsProvider } from '@/lib/contexts/CalculationsContext';
import type { Goal, FinancialData, FinancialProfile } from '@/types/api';

// AppProviders のネスト順を再現するラッパー
function AllProviders({ children }: { children: ReactNode }) {
  return (
    <AuthProvider>
      <GuestModeProvider>
        <ThemeProvider>
          <TutorialProvider>
            <FinancialDataProvider>
              <GoalsProvider>
                <CalculationsProvider>
                  {children}
                </CalculationsProvider>
              </GoalsProvider>
            </FinancialDataProvider>
          </TutorialProvider>
        </ThemeProvider>
      </GuestModeProvider>
    </AuthProvider>
  );
}

// カスタム render: 全プロバイダー付き
function renderWithProviders(
  ui: ReactElement,
  options?: Omit<RenderOptions, 'wrapper'>
) {
  return render(ui, { wrapper: AllProviders, ...options });
}

// カスタム renderHook: 全プロバイダー付き
function renderHookWithProviders<Result, Props>(
  hook: (props: Props) => Result,
  options?: Omit<RenderHookOptions<Props>, 'wrapper'>
) {
  return renderHook(hook, { wrapper: AllProviders, ...options });
}

// ── モックデータファクトリ ──

function createMockGoal(overrides: Partial<Goal> = {}): Goal {
  return {
    id: 'goal-1',
    user_id: 'user-1',
    goal_type: 'savings',
    title: '貯蓄目標',
    target_amount: 5000000,
    target_date: '2027-12-31T00:00:00.000Z',
    current_amount: 1000000,
    monthly_contribution: 50000,
    is_active: true,
    created_at: '2024-01-01T00:00:00.000Z',
    updated_at: '2024-06-01T00:00:00.000Z',
    ...overrides,
  };
}

function createMockFinancialProfile(overrides: Partial<FinancialProfile> = {}): FinancialProfile {
  return {
    monthly_income: 400000,
    monthly_expenses: [
      { category: '住居費', amount: 100000 },
      { category: '食費', amount: 50000 },
    ],
    current_savings: [
      { type: 'deposit', amount: 2000000 },
      { type: 'investment', amount: 1000000 },
    ],
    investment_return: 5.0,
    inflation_rate: 2.0,
    ...overrides,
  };
}

function createMockFinancialData(overrides: Partial<FinancialData> = {}): FinancialData {
  return {
    id: 'fd-1',
    user_id: 'user-1',
    profile: createMockFinancialProfile(),
    retirement: {
      current_age: 35,
      retirement_age: 65,
      life_expectancy: 90,
      monthly_retirement_expenses: 250000,
      pension_amount: 150000,
    },
    emergency_fund: {
      target_months: 6,
      monthly_expenses: 250000,
      current_amount: 600000,
    },
    created_at: '2024-01-01T00:00:00.000Z',
    updated_at: '2024-06-01T00:00:00.000Z',
    ...overrides,
  };
}

export {
  renderWithProviders,
  renderHookWithProviders,
  AllProviders,
  createMockGoal,
  createMockFinancialProfile,
  createMockFinancialData,
};
