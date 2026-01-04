'use client';

import React, { ReactNode } from 'react';
import { FinancialDataProvider } from './FinancialDataContext';
import { GoalsProvider } from './GoalsContext';
import { CalculationsProvider } from './CalculationsContext';
import { ThemeProvider } from './ThemeContext';

interface AppProvidersProps {
  children: ReactNode;
}

/**
 * アプリケーション全体のプロバイダーをまとめたコンポーネント
 * すべてのコンテキストプロバイダーをここで統合
 */
export function AppProviders({ children }: AppProvidersProps) {
  return (
    <ThemeProvider>
      <FinancialDataProvider>
        <GoalsProvider>
          <CalculationsProvider>
            {children}
          </CalculationsProvider>
        </GoalsProvider>
      </FinancialDataProvider>
    </ThemeProvider>
  );
}
