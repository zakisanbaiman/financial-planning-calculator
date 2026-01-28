'use client';

import React, { createContext, useContext, useState, useCallback, ReactNode } from 'react';
import { financialDataAPI } from '@/lib/api-client';
import { useGuestMode } from './GuestModeContext';
import type {
  FinancialData,
  FinancialProfile,
  RetirementData,
  EmergencyFund,
} from '@/types/api';

// コンテキスト型定義
interface FinancialDataContextType {
  financialData: FinancialData | null;
  loading: boolean;
  error: string | null;
  fetchFinancialData: (userId: string) => Promise<void>;
  createFinancialData: (data: FinancialData) => Promise<void>;
  updateProfile: (userId: string, profile: FinancialProfile) => Promise<void>;
  updateRetirement: (userId: string, retirement: RetirementData) => Promise<void>;
  updateEmergencyFund: (userId: string, emergencyFund: EmergencyFund) => Promise<void>;
  deleteFinancialData: (userId: string) => Promise<void>;
  clearError: () => void;
}

// コンテキスト作成
const FinancialDataContext = createContext<FinancialDataContextType | undefined>(
  undefined
);

// プロバイダープロパティ
interface FinancialDataProviderProps {
  children: ReactNode;
}

// プロバイダーコンポーネント
export function FinancialDataProvider({ children }: FinancialDataProviderProps) {
  const [financialData, setFinancialData] = useState<FinancialData | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const { isGuestMode, guestData, setGuestData } = useGuestMode();

  // エラークリア
  const clearError = useCallback(() => {
    setError(null);
  }, []);

  // 財務データ取得
  const fetchFinancialData = useCallback(async (userId: string) => {
    // ゲストモードの場合はローカルストレージから取得
    if (isGuestMode) {
      setLoading(true);
      try {
        if (guestData) {
          setFinancialData(guestData);
        } else {
          setError('財務データがまだ作成されていません。下のフォームから入力してください。');
        }
      } finally {
        setLoading(false);
      }
      return;
    }

    setLoading(true);
    setError(null);
    try {
      const data = await financialDataAPI.get(userId);
      setFinancialData(data);
    } catch (err) {
      // APIError の場合は status を確認
      let errorMessage = '財務データの取得に失敗しました';
      if (err instanceof Error) {
        // 404 の場合は、ユーザーが財務データをまだ入力していないことを示す
        if ((err as any).status === 404) {
          errorMessage = '財務データがまだ作成されていません。下のフォームから入力してください。';
        } else {
          errorMessage = err.message;
        }
      }
      setError(errorMessage);
      throw err;
    } finally {
      setLoading(false);
    }
  }, [isGuestMode, guestData]);

  // 財務データ作成
  const createFinancialData = useCallback(async (data: FinancialData) => {
    // ゲストモードの場合はローカルストレージに保存
    if (isGuestMode) {
      setLoading(true);
      try {
        setGuestData(data);
        setFinancialData(data);
      } finally {
        setLoading(false);
      }
      return;
    }

    setLoading(true);
    setError(null);
    try {
      const created = await financialDataAPI.create(data);
      setFinancialData(created);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : '財務データの作成に失敗しました';
      setError(errorMessage);
      throw err;
    } finally {
      setLoading(false);
    }
  }, [isGuestMode, setGuestData]);

  // 財務プロファイル更新
  const updateProfile = useCallback(
    async (userId: string, profile: FinancialProfile) => {
      // ゲストモードの場合はローカルストレージに保存
      if (isGuestMode) {
        setLoading(true);
        try {
          const updated = { ...financialData, user_id: userId, profile } as FinancialData;
          setGuestData(updated);
          setFinancialData(updated);
        } finally {
          setLoading(false);
        }
        return;
      }

      setLoading(true);
      setError(null);
      try {
        const updated = await financialDataAPI.updateProfile(userId, profile);
        setFinancialData(updated);
      } catch (err) {
        const errorMessage = err instanceof Error ? err.message : 'プロファイルの更新に失敗しました';
        setError(errorMessage);
        throw err;
      } finally {
        setLoading(false);
      }
    },
    [isGuestMode, financialData, setGuestData]
  );

  // 退職データ更新
  const updateRetirement = useCallback(
    async (userId: string, retirement: RetirementData) => {
      // ゲストモードの場合はローカルストレージに保存
      if (isGuestMode) {
        setLoading(true);
        try {
          const updated = { ...financialData, user_id: userId, retirement } as FinancialData;
          setGuestData(updated);
          setFinancialData(updated);
        } finally {
          setLoading(false);
        }
        return;
      }

      setLoading(true);
      setError(null);
      try {
        const updated = await financialDataAPI.updateRetirement(userId, retirement);
        setFinancialData(updated);
      } catch (err) {
        const errorMessage = err instanceof Error ? err.message : '退職データの更新に失敗しました';
        setError(errorMessage);
        throw err;
      } finally {
        setLoading(false);
      }
    },
    [isGuestMode, financialData, setGuestData]
  );

  // 緊急資金更新
  const updateEmergencyFund = useCallback(
    async (userId: string, emergencyFund: EmergencyFund) => {
      // ゲストモードの場合はローカルストレージに保存
      if (isGuestMode) {
        setLoading(true);
        try {
          const updated = { ...financialData, user_id: userId, emergency_fund: emergencyFund } as FinancialData;
          setGuestData(updated);
          setFinancialData(updated);
        } finally {
          setLoading(false);
        }
        return;
      }

      setLoading(true);
      setError(null);
      try {
        const updated = await financialDataAPI.updateEmergencyFund(userId, emergencyFund);
        setFinancialData(updated);
      } catch (err) {
        const errorMessage = err instanceof Error ? err.message : '緊急資金の更新に失敗しました';
        setError(errorMessage);
        throw err;
      } finally {
        setLoading(false);
      }
    },
    [isGuestMode, financialData, setGuestData]
  );

  // 財務データ削除
  const deleteFinancialData = useCallback(async (userId: string) => {
    // ゲストモードの場合はローカルストレージから削除
    if (isGuestMode) {
      setLoading(true);
      try {
        setGuestData(null);
        setFinancialData(null);
      } finally {
        setLoading(false);
      }
      return;
    }

    setLoading(true);
    setError(null);
    try {
      await financialDataAPI.delete(userId);
      setFinancialData(null);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : '財務データの削除に失敗しました';
      setError(errorMessage);
      throw err;
    } finally {
      setLoading(false);
    }
  }, [isGuestMode, setGuestData]);

  const value: FinancialDataContextType = {
    financialData,
    loading,
    error,
    fetchFinancialData,
    createFinancialData,
    updateProfile,
    updateRetirement,
    updateEmergencyFund,
    deleteFinancialData,
    clearError,
  };

  return (
    <FinancialDataContext.Provider value={value}>
      {children}
    </FinancialDataContext.Provider>
  );
}

// カスタムフック
export function useFinancialData() {
  const context = useContext(FinancialDataContext);
  if (context === undefined) {
    throw new Error('useFinancialData must be used within a FinancialDataProvider');
  }
  return context;
}
