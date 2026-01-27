'use client';

import React, { ReactNode } from 'react';
import { AuthProvider } from './AuthContext';
import { GuestModeProvider } from './GuestModeContext';
import { FinancialDataProvider } from './FinancialDataContext';
import { GoalsProvider } from './GoalsContext';
import { CalculationsProvider } from './CalculationsContext';
import { ThemeProvider } from './ThemeContext';
import { TutorialProvider } from './TutorialContext';

interface AppProvidersProps {
  children: ReactNode;
}

/**
 * アプリケーション全体のプロバイダーをまとめたコンポーネント
 * すべてのコンテキストプロバイダーをここで統合
 */
export function AppProviders({ children }: AppProvidersProps) {
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
