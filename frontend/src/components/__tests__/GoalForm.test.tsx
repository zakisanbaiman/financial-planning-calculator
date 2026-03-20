import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import GoalForm from '../GoalForm';
import type { Goal } from '@/types/api';

// requestAnimationFrameモック（CurrencyInput使用のため）
beforeAll(() => {
  jest.spyOn(window, 'requestAnimationFrame').mockImplementation((cb) => {
    cb(0);
    return 0;
  });
});

afterAll(() => {
  (window.requestAnimationFrame as jest.Mock).mockRestore();
});

const mockOnSubmit = jest.fn();
const mockOnCancel = jest.fn();

const mockGoal: Goal = {
  id: 'goal-1',
  user_id: 'user-1',
  goal_type: 'savings',
  title: '既存の目標',
  target_amount: 3000000,
  target_date: '2027-12-31T00:00:00.000Z',
  current_amount: 500000,
  monthly_contribution: 30000,
  is_active: true,
  created_at: '2024-01-01T00:00:00.000Z',
  updated_at: '2024-06-01T00:00:00.000Z',
};

describe('GoalForm', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('フォーム表示', () => {
    it('目標タイプのボタンが表示される', () => {
      render(
        <GoalForm userId="user-1" onSubmit={mockOnSubmit} />
      );
      expect(screen.getByRole('button', { name: '貯蓄' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: '老後資金' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: '緊急資金' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'カスタム' })).toBeInTheDocument();
    });

    it('作成ボタンが表示される', () => {
      render(
        <GoalForm userId="user-1" onSubmit={mockOnSubmit} />
      );
      expect(screen.getByRole('button', { name: '作成' })).toBeInTheDocument();
    });

    it('キャンセルボタンが表示される（onCancel提供時）', () => {
      render(
        <GoalForm userId="user-1" onSubmit={mockOnSubmit} onCancel={mockOnCancel} />
      );
      expect(screen.getByRole('button', { name: 'キャンセル' })).toBeInTheDocument();
    });

    it('initialDataがある場合は更新ボタンが表示される', () => {
      render(
        <GoalForm userId="user-1" initialData={mockGoal} onSubmit={mockOnSubmit} />
      );
      expect(screen.getByRole('button', { name: '更新' })).toBeInTheDocument();
    });

    it('initialDataのタイトルが表示される', () => {
      render(
        <GoalForm userId="user-1" initialData={mockGoal} onSubmit={mockOnSubmit} />
      );
      expect(screen.getByDisplayValue('既存の目標')).toBeInTheDocument();
    });

    it('デフォルトで貯蓄タイプのデフォルトタイトルが設定される', async () => {
      render(
        <GoalForm userId="user-1" onSubmit={mockOnSubmit} />
      );
      // useEffectによりデフォルトタイトルが設定される
      await waitFor(() => {
        expect(screen.getByDisplayValue('貯蓄目標')).toBeInTheDocument();
      });
    });
  });

  describe('目標タイプ選択', () => {
    it('老後資金ボタンをクリックするとそのタイプが選択状態になる', async () => {
      const user = userEvent.setup();
      render(
        <GoalForm userId="user-1" onSubmit={mockOnSubmit} />
      );

      await user.click(screen.getByRole('button', { name: '老後資金' }));

      // 老後資金ボタンがアクティブクラスを持つ
      const btn = screen.getByRole('button', { name: '老後資金' });
      expect(btn.className).toContain('border-primary-500');
    });

    it('緊急資金ボタンをクリックするとそのタイプが選択状態になる', async () => {
      const user = userEvent.setup();
      render(
        <GoalForm userId="user-1" onSubmit={mockOnSubmit} />
      );

      await user.click(screen.getByRole('button', { name: '緊急資金' }));

      const btn = screen.getByRole('button', { name: '緊急資金' });
      expect(btn.className).toContain('border-primary-500');
    });

    it('カスタムボタンをクリックするとそのタイプが選択状態になる', async () => {
      const user = userEvent.setup();
      render(
        <GoalForm userId="user-1" onSubmit={mockOnSubmit} />
      );

      await user.click(screen.getByRole('button', { name: 'カスタム' }));

      const btn = screen.getByRole('button', { name: 'カスタム' });
      expect(btn.className).toContain('border-primary-500');
    });
  });

  describe('フォーム送信', () => {
    it('有効なデータでフォーム送信されると onSubmit が呼ばれる', async () => {
      const user = userEvent.setup();
      mockOnSubmit.mockResolvedValue(undefined);

      render(
        <GoalForm userId="user-1" onSubmit={mockOnSubmit} />
      );

      await user.click(screen.getByRole('button', { name: '作成' }));

      await waitFor(() => {
        expect(mockOnSubmit).toHaveBeenCalledWith(
          expect.objectContaining({
            user_id: 'user-1',
            goal_type: 'savings',
          })
        );
      });
    });

    it('キャンセルボタンクリックで onCancel が呼ばれる', async () => {
      const user = userEvent.setup();
      render(
        <GoalForm userId="user-1" onSubmit={mockOnSubmit} onCancel={mockOnCancel} />
      );

      await user.click(screen.getByRole('button', { name: 'キャンセル' }));

      expect(mockOnCancel).toHaveBeenCalled();
    });
  });

  describe('進捗表示', () => {
    it('initialDataがある場合に進捗状況が表示される', () => {
      render(
        <GoalForm userId="user-1" initialData={mockGoal} onSubmit={mockOnSubmit} />
      );
      expect(screen.getByText('進捗状況')).toBeInTheDocument();
    });

    it('「この目標をアクティブにする」チェックボックスが表示される', () => {
      render(
        <GoalForm userId="user-1" onSubmit={mockOnSubmit} />
      );
      expect(screen.getByLabelText('この目標をアクティブにする')).toBeInTheDocument();
    });
  });
});
