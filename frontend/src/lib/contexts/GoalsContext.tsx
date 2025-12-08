'use client';

import React, { createContext, useContext, useState, useCallback, ReactNode } from 'react';
import { goalsAPI } from '@/lib/api-client';
import type { Goal } from '@/types/api';

// コンテキスト型定義
interface GoalsContextType {
  goals: Goal[];
  loading: boolean;
  error: string | null;
  fetchGoals: (userId: string) => Promise<void>;
  createGoal: (goal: Goal) => Promise<void>;
  updateGoal: (id: string, userId: string, goal: Partial<Goal>) => Promise<void>;
  updateGoalProgress: (id: string, userId: string, currentAmount: number) => Promise<void>;
  deleteGoal: (id: string, userId: string) => Promise<void>;
  clearError: () => void;
}

// コンテキスト作成
const GoalsContext = createContext<GoalsContextType | undefined>(undefined);

// プロバイダープロパティ
interface GoalsProviderProps {
  children: ReactNode;
}

// プロバイダーコンポーネント
export function GoalsProvider({ children }: GoalsProviderProps) {
  const [goals, setGoals] = useState<Goal[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // エラークリア
  const clearError = useCallback(() => {
    setError(null);
  }, []);

  // 目標一覧取得
  const fetchGoals = useCallback(async (userId: string) => {
    setLoading(true);
    setError(null);
    try {
      const data = await goalsAPI.list(userId);

      // API may return either an array of goals or an object like { goals: [...], summary: {...} }
      if (Array.isArray(data)) {
        setGoals(data);
      } else if (data && Array.isArray((data as any).goals)) {
        // If each item is a wrapper with a nested `goal` field (GoalWithStatus), extract it.
        const extracted = (data as any).goals.map((g: any) => (g.goal ? g.goal : g));
        setGoals(extracted);
      } else {
        // Fallback: set to empty array to avoid runtime errors
        setGoals([]);
      }
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : '目標の取得に失敗しました';
      setError(errorMessage);
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  // 目標作成
  const createGoal = useCallback(async (goal: Goal) => {
    setLoading(true);
    setError(null);
    try {
      const created = await goalsAPI.create(goal);
      setGoals((prev) => [...prev, created]);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : '目標の作成に失敗しました';
      setError(errorMessage);
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  // 目標更新
  const updateGoal = useCallback(
    async (id: string, userId: string, goal: Partial<Goal>) => {
      setLoading(true);
      setError(null);
      try {
        const updated = await goalsAPI.update(id, userId, goal);
        setGoals((prev) =>
          prev.map((g) => (g.id === id ? updated : g))
        );
      } catch (err) {
        const errorMessage = err instanceof Error ? err.message : '目標の更新に失敗しました';
        setError(errorMessage);
        throw err;
      } finally {
        setLoading(false);
      }
    },
    []
  );

  // 目標進捗更新
  const updateGoalProgress = useCallback(
    async (id: string, userId: string, currentAmount: number) => {
      setLoading(true);
      setError(null);
      try {
        const updated = await goalsAPI.updateProgress(id, userId, currentAmount);
        setGoals((prev) =>
          prev.map((g) => (g.id === id ? updated : g))
        );
      } catch (err) {
        const errorMessage = err instanceof Error ? err.message : '進捗の更新に失敗しました';
        setError(errorMessage);
        throw err;
      } finally {
        setLoading(false);
      }
    },
    []
  );

  // 目標削除
  const deleteGoal = useCallback(async (id: string, userId: string) => {
    setLoading(true);
    setError(null);
    try {
      await goalsAPI.delete(id, userId);
      setGoals((prev) => prev.filter((g) => g.id !== id));
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : '目標の削除に失敗しました';
      setError(errorMessage);
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  const value: GoalsContextType = {
    goals,
    loading,
    error,
    fetchGoals,
    createGoal,
    updateGoal,
    updateGoalProgress,
    deleteGoal,
    clearError,
  };

  return <GoalsContext.Provider value={value}>{children}</GoalsContext.Provider>;
}

// カスタムフック
export function useGoals() {
  const context = useContext(GoalsContext);
  if (context === undefined) {
    throw new Error('useGoals must be used within a GoalsProvider');
  }
  return context;
}
