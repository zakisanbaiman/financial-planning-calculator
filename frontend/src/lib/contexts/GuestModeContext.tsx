'use client';

import React, { createContext, useContext, useState, useEffect, useCallback, ReactNode } from 'react';
import type { FinancialData } from '@/types/api';

// ゲストモードコンテキストの型
interface GuestModeContextType {
  isGuestMode: boolean;
  guestData: FinancialData | null;
  setGuestData: (data: FinancialData | null) => void;
  clearGuestData: () => void;
  startGuestMode: () => void;
  exitGuestMode: () => void;
}

const GuestModeContext = createContext<GuestModeContextType | undefined>(undefined);

// ローカルストレージのキー
const GUEST_MODE_KEY = 'guest_mode';
const GUEST_DATA_KEY = 'guest_financial_data';

interface GuestModeProviderProps {
  children: ReactNode;
}

export function GuestModeProvider({ children }: GuestModeProviderProps) {
  const [isGuestMode, setIsGuestMode] = useState(false);
  const [guestData, setGuestDataState] = useState<FinancialData | null>(null);

  // 初期化時にローカルストレージからゲストモードの状態を復元
  useEffect(() => {
    try {
      const storedMode = localStorage.getItem(GUEST_MODE_KEY);
      const storedData = localStorage.getItem(GUEST_DATA_KEY);

      if (storedMode === 'true') {
        setIsGuestMode(true);
        if (storedData) {
          setGuestDataState(JSON.parse(storedData));
        }
      }
    } catch (e) {
      console.error('Failed to restore guest mode state:', e);
    }
  }, []);

  // ゲストデータを保存
  const setGuestData = useCallback((data: FinancialData | null) => {
    setGuestDataState(data);
    if (data) {
      localStorage.setItem(GUEST_DATA_KEY, JSON.stringify(data));
    } else {
      localStorage.removeItem(GUEST_DATA_KEY);
    }
  }, []);

  // ゲストデータをクリア
  const clearGuestData = useCallback(() => {
    setGuestDataState(null);
    localStorage.removeItem(GUEST_DATA_KEY);
  }, []);

  // ゲストモードを開始
  const startGuestMode = useCallback(() => {
    setIsGuestMode(true);
    localStorage.setItem(GUEST_MODE_KEY, 'true');
  }, []);

  // ゲストモードを終了
  const exitGuestMode = useCallback(() => {
    setIsGuestMode(false);
    clearGuestData();
    localStorage.removeItem(GUEST_MODE_KEY);
  }, [clearGuestData]);

  const value = {
    isGuestMode,
    guestData,
    setGuestData,
    clearGuestData,
    startGuestMode,
    exitGuestMode,
  };

  return <GuestModeContext.Provider value={value}>{children}</GuestModeContext.Provider>;
}

// ゲストモードコンテキストを使用するフック
export function useGuestMode() {
  const context = useContext(GuestModeContext);
  if (context === undefined) {
    throw new Error('useGuestMode must be used within a GuestModeProvider');
  }
  return context;
}
