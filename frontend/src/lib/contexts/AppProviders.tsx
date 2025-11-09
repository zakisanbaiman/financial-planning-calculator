'use client';

import React, { ReactNode } from 'react';
import { FinancialDataProvider } from './FinancialDataContext';
import { GoalsProvider } from './GoalsContext';
import { CalculationsProvider } from './CalculationsContext';

interface AppProvidersProps {
  children: ReactNode;
}

/**
 * アプリケーション全体のプロバイダーをまとめたコンポーネント
 * すべてのコンテキストプロバイダーをここで統合
 */
export function AppProviders({ children }: AppProvidersProps) {
  return (
    <FinancialDataProvider>
      <GoalsProvider>
        <CalculationsProvider>
          {children}
        </CalculationsProvider>
      </GoalsProvider>
    </FinancialDataProvider>
  );
}
