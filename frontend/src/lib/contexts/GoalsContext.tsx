'use client';

import React, { createContext, useContext, useState, useCallback, ReactNode } from 'react';
import { goalsAPI } from '@/lib/api-client';
import { useGuestMode } from './GuestModeContext';
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

// ローカルストレージのキー
const GUEST_GOALS_KEY = 'guest_goals';

// プロバイダーコンポーネント
export function GoalsProvider({ children }: GoalsProviderProps) {
  const [goals, setGoals] = useState<Goal[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const { isGuestMode } = useGuestMode();

  // エラークリア
  const clearError = useCallback(() => {
    setError(null);
  }, []);

  // 目標一覧取得
  const fetchGoals = useCallback(async (userId: string) => {
    // ゲストモードの場合はローカルストレージから取得
    if (isGuestMode) {
      setLoading(true);
      setError(null);
      try {
        const stored = localStorage.getItem(GUEST_GOALS_KEY);
        if (stored) {
          const data = JSON.parse(stored) as Goal[];
          setGoals(data);
        } else {
          setGoals([]);
        }
      } catch (err) {
        console.error('Failed to restore guest goals:', err);
        setGoals([]);
      } finally {
        setLoading(false);
      }
      return;
    }

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
  }, [isGuestMode]);

  // 目標作成
  const createGoal = useCallback(async (goal: Goal) => {
    // ゲストモードの場合はローカルストレージに保存
    if (isGuestMode) {
      setLoading(true);
      setError(null);
      try {
        const newGoal: Goal = {
          ...goal,
          id: goal.id || `guest-${Date.now()}`,
          created_at: goal.created_at || new Date().toISOString(),
          updated_at: new Date().toISOString(),
        };
        setGoals((prev) => {
          const updated = [...prev, newGoal];
          localStorage.setItem(GUEST_GOALS_KEY, JSON.stringify(updated));
          return updated;
        });
        return;
      } catch (err) {
        const errorMessage = err instanceof Error ? err.message : '目標の作成に失敗しました';
        setError(errorMessage);
        throw err;
      } finally {
        setLoading(false);
      }
    }

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
  }, [isGuestMode]);

  // 目標更新
  const updateGoal = useCallback(
    async (id: string, userId: string, goal: Partial<Goal>) => {
      // ゲストモードの場合はローカルストレージを更新
      if (isGuestMode) {
        setLoading(true);
        setError(null);
        try {
          setGoals((prev) => {
            const updated = prev.map((g) =>
              g.id === id ? { ...g, ...goal, updated_at: new Date().toISOString() } : g
            );
            localStorage.setItem(GUEST_GOALS_KEY, JSON.stringify(updated));
            return updated;
          });
          return;
        } catch (err) {
          const errorMessage = err instanceof Error ? err.message : '目標の更新に失敗しました';
          setError(errorMessage);
          throw err;
        } finally {
          setLoading(false);
        }
      }

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
    [isGuestMode]
  );

  // 目標進捗更新
  const updateGoalProgress = useCallback(
    async (id: string, userId: string, currentAmount: number) => {
      // ゲストモードの場合はローカルストレージを更新
      if (isGuestMode) {
        setLoading(true);
        setError(null);
        try {
          setGoals((prev) => {
            const updated = prev.map((g) =>
              g.id === id ? { ...g, current_amount: currentAmount, updated_at: new Date().toISOString() } : g
            );
            localStorage.setItem(GUEST_GOALS_KEY, JSON.stringify(updated));
            return updated;
          });
          return;
        } catch (err) {
          const errorMessage = err instanceof Error ? err.message : '進捗の更新に失敗しました';
          setError(errorMessage);
          throw err;
        } finally {
          setLoading(false);
        }
      }

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
    [isGuestMode]
  );

  // 目標削除
  const deleteGoal = useCallback(async (id: string, userId: string) => {
    // ゲストモードの場合はローカルストレージから削除
    if (isGuestMode) {
      setLoading(true);
      setError(null);
      try {
        setGoals((prev) => {
          const updated = prev.filter((g) => g.id !== id);
          localStorage.setItem(GUEST_GOALS_KEY, JSON.stringify(updated));
          return updated;
        });
        return;
      } catch (err) {
        const errorMessage = err instanceof Error ? err.message : '目標の削除に失敗しました';
        setError(errorMessage);
        throw err;
      } finally {
        setLoading(false);
      }
    }

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
  }, [isGuestMode]);

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
