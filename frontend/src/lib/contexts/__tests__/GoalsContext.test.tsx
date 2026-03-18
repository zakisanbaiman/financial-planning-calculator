import React from 'react';
import { renderHook, act } from '@testing-library/react';
import { GoalsProvider, useGoals } from '../GoalsContext';
import { GuestModeProvider } from '../GuestModeContext';
import { goalsAPI } from '@/lib/api-client';
import type { Goal } from '@/types/api';

jest.mock('@/lib/api-client');
const mockedGoalsAPI = goalsAPI as jest.Mocked<typeof goalsAPI>;

const wrapper = ({ children }: { children: React.ReactNode }) => (
  <GuestModeProvider>
    <GoalsProvider>{children}</GoalsProvider>
  </GuestModeProvider>
);

const mockGoal: Goal = {
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
};

describe('GoalsContext', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    (localStorage.getItem as jest.Mock).mockReturnValue(null);
  });

  describe('初期状態', () => {
    it('初期状態は空のgoals配列', () => {
      const { result } = renderHook(() => useGoals(), { wrapper });
      expect(result.current.goals).toEqual([]);
      expect(result.current.loading).toBe(false);
      expect(result.current.error).toBeNull();
    });
  });

  describe('目標取得 (APIモード)', () => {
    it('fetchGoalsで目標一覧を取得できる', async () => {
      mockedGoalsAPI.list.mockResolvedValue([mockGoal]);

      const { result } = renderHook(() => useGoals(), { wrapper });

      await act(async () => {
        await result.current.fetchGoals('user-1');
      });

      expect(result.current.goals).toEqual([mockGoal]);
      expect(mockedGoalsAPI.list).toHaveBeenCalledWith('user-1');
    });

    it('fetchGoalsでエラーが発生した場合、errorが設定される', async () => {
      mockedGoalsAPI.list.mockRejectedValue(new Error('取得に失敗'));

      const { result } = renderHook(() => useGoals(), { wrapper });

      let thrownError: Error | undefined;
      await act(async () => {
        try {
          await result.current.fetchGoals('user-1');
        } catch (e) {
          thrownError = e as Error;
        }
      });

      expect(thrownError).toBeDefined();
      expect(result.current.error).toBe('取得に失敗');
    });
  });

  describe('目標作成 (APIモード)', () => {
    it('createGoalで目標を作成できる', async () => {
      mockedGoalsAPI.create.mockResolvedValue(mockGoal);

      const { result } = renderHook(() => useGoals(), { wrapper });

      await act(async () => {
        await result.current.createGoal(mockGoal);
      });

      expect(result.current.goals).toContainEqual(mockGoal);
      expect(mockedGoalsAPI.create).toHaveBeenCalledWith(mockGoal);
    });
  });

  describe('目標更新 (APIモード)', () => {
    it('updateGoalで目標を更新できる', async () => {
      const updatedGoal = { ...mockGoal, title: '更新後の目標' };
      mockedGoalsAPI.list.mockResolvedValue([mockGoal]);
      mockedGoalsAPI.update.mockResolvedValue(updatedGoal);

      const { result } = renderHook(() => useGoals(), { wrapper });

      await act(async () => {
        await result.current.fetchGoals('user-1');
      });

      await act(async () => {
        await result.current.updateGoal('goal-1', 'user-1', { title: '更新後の目標' });
      });

      expect(result.current.goals[0].title).toBe('更新後の目標');
    });
  });

  describe('目標進捗更新 (APIモード)', () => {
    it('updateGoalProgressで進捗を更新できる', async () => {
      const updatedGoal = { ...mockGoal, current_amount: 2000000 };
      mockedGoalsAPI.list.mockResolvedValue([mockGoal]);
      mockedGoalsAPI.updateProgress.mockResolvedValue(updatedGoal);

      const { result } = renderHook(() => useGoals(), { wrapper });

      await act(async () => {
        await result.current.fetchGoals('user-1');
      });

      await act(async () => {
        await result.current.updateGoalProgress('goal-1', 'user-1', 2000000);
      });

      expect(result.current.goals[0].current_amount).toBe(2000000);
    });
  });

  describe('目標削除 (APIモード)', () => {
    it('deleteGoalで目標を削除できる', async () => {
      mockedGoalsAPI.list.mockResolvedValue([mockGoal]);
      mockedGoalsAPI.delete.mockResolvedValue(undefined);

      const { result } = renderHook(() => useGoals(), { wrapper });

      await act(async () => {
        await result.current.fetchGoals('user-1');
      });

      await act(async () => {
        await result.current.deleteGoal('goal-1', 'user-1');
      });

      expect(result.current.goals).toEqual([]);
    });
  });

  describe('エラー管理', () => {
    it('clearErrorでエラーがクリアされる', async () => {
      mockedGoalsAPI.list.mockRejectedValue(new Error('エラー'));

      const { result } = renderHook(() => useGoals(), { wrapper });

      await act(async () => {
        try {
          await result.current.fetchGoals('user-1');
        } catch {
          // expected
        }
      });

      expect(result.current.error).toBe('エラー');

      act(() => {
        result.current.clearError();
      });

      expect(result.current.error).toBeNull();
    });
  });

  describe('ゲストモード', () => {
    it('ゲストモードではlocalStorageから目標を取得する', async () => {
      (localStorage.getItem as jest.Mock).mockImplementation((key: string) => {
        if (key === 'guest_mode') return 'true';
        if (key === 'guest_goals') return JSON.stringify([mockGoal]);
        return null;
      });

      const { result } = renderHook(() => useGoals(), { wrapper });

      await act(async () => {
        await result.current.fetchGoals('guest');
      });

      expect(result.current.goals).toEqual([mockGoal]);
      expect(mockedGoalsAPI.list).not.toHaveBeenCalled();
    });
  });

  describe('Provider外でのフック使用', () => {
    it('Provider外で useGoals を使うとエラーが発生する', () => {
      const consoleSpy = jest.spyOn(console, 'error').mockImplementation(() => {});
      expect(() => {
        renderHook(() => useGoals());
      }).toThrow('useGoals must be used within a GoalsProvider');
      consoleSpy.mockRestore();
    });
  });
});
