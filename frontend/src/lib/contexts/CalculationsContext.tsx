'use client';

import React, { createContext, useContext, useState, useCallback, ReactNode } from 'react';
import { calculationsAPI } from '@/lib/api-client';
import type {
  AssetProjectionRequest,
  AssetProjectionResponse,
  RetirementCalculationRequest,
  RetirementCalculationResponse,
  EmergencyFundRequest,
  EmergencyFundResponse,
  GoalProjectionRequest,
  GoalProjectionResponse,
} from '@/types/api';

// コンテキスト型定義
interface CalculationsContextType {
  assetProjection: AssetProjectionResponse | null;
  retirementCalculation: RetirementCalculationResponse | null;
  emergencyFund: EmergencyFundResponse | null;
  goalProjection: GoalProjectionResponse | null;
  loading: boolean;
  error: string | null;
  calculateAssetProjection: (data: AssetProjectionRequest) => Promise<void>;
  calculateRetirement: (data: RetirementCalculationRequest) => Promise<void>;
  calculateEmergencyFund: (data: EmergencyFundRequest) => Promise<void>;
  calculateGoalProjection: (data: GoalProjectionRequest) => Promise<void>;
  clearCalculations: () => void;
  clearError: () => void;
}

// コンテキスト作成
const CalculationsContext = createContext<CalculationsContextType | undefined>(
  undefined
);

// プロバイダープロパティ
interface CalculationsProviderProps {
  children: ReactNode;
}

// プロバイダーコンポーネント
export function CalculationsProvider({ children }: CalculationsProviderProps) {
  const [assetProjection, setAssetProjection] = useState<AssetProjectionResponse | null>(null);
  const [retirementCalculation, setRetirementCalculation] = useState<RetirementCalculationResponse | null>(null);
  const [emergencyFund, setEmergencyFund] = useState<EmergencyFundResponse | null>(null);
  const [goalProjection, setGoalProjection] = useState<GoalProjectionResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // エラークリア
  const clearError = useCallback(() => {
    setError(null);
  }, []);

  // 計算結果クリア
  const clearCalculations = useCallback(() => {
    setAssetProjection(null);
    setRetirementCalculation(null);
    setEmergencyFund(null);
    setGoalProjection(null);
    setError(null);
  }, []);

  // 資産推移計算
  const calculateAssetProjection = useCallback(
    async (data: AssetProjectionRequest) => {
      setLoading(true);
      setError(null);
      try {
        const result = await calculationsAPI.assetProjection(data);
        setAssetProjection(result);
      } catch (err) {
        const errorMessage = err instanceof Error ? err.message : '資産推移の計算に失敗しました';
        setError(errorMessage);
        throw err;
      } finally {
        setLoading(false);
      }
    },
    []
  );

  // 老後資金計算
  const calculateRetirement = useCallback(
    async (data: RetirementCalculationRequest) => {
      setLoading(true);
      setError(null);
      try {
        const result = await calculationsAPI.retirement(data);
        setRetirementCalculation(result);
      } catch (err) {
        const errorMessage = err instanceof Error ? err.message : '老後資金の計算に失敗しました';
        setError(errorMessage);
        throw err;
      } finally {
        setLoading(false);
      }
    },
    []
  );

  // 緊急資金計算
  const calculateEmergencyFund = useCallback(
    async (data: EmergencyFundRequest) => {
      setLoading(true);
      setError(null);
      try {
        const result = await calculationsAPI.emergencyFund(data);
        setEmergencyFund(result);
      } catch (err) {
        const errorMessage = err instanceof Error ? err.message : '緊急資金の計算に失敗しました';
        setError(errorMessage);
        throw err;
      } finally {
        setLoading(false);
      }
    },
    []
  );

  // 目標達成計算
  const calculateGoalProjection = useCallback(
    async (data: GoalProjectionRequest) => {
      setLoading(true);
      setError(null);
      try {
        const result = await calculationsAPI.goalProjection(data);
        setGoalProjection(result);
      } catch (err) {
        const errorMessage = err instanceof Error ? err.message : '目標達成の計算に失敗しました';
        setError(errorMessage);
        throw err;
      } finally {
        setLoading(false);
      }
    },
    []
  );

  const value: CalculationsContextType = {
    assetProjection,
    retirementCalculation,
    emergencyFund,
    goalProjection,
    loading,
    error,
    calculateAssetProjection,
    calculateRetirement,
    calculateEmergencyFund,
    calculateGoalProjection,
    clearCalculations,
    clearError,
  };

  return (
    <CalculationsContext.Provider value={value}>
      {children}
    </CalculationsContext.Provider>
  );
}

// カスタムフック
export function useCalculations() {
  const context = useContext(CalculationsContext);
  if (context === undefined) {
    throw new Error('useCalculations must be used within a CalculationsProvider');
  }
  return context;
}
