import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import GoalProgressTracker from '../GoalProgressTracker';
import type { Goal } from '@/types/api';

const mockGoalInProgress: Goal = {
  id: 'goal-1',
  user_id: 'user-1',
  goal_type: 'savings',
  title: '貯蓄目標',
  target_amount: 5000000,
  target_date: '2030-12-31T00:00:00.000Z', // 遠い未来
  current_amount: 1000000,
  monthly_contribution: 50000,
  is_active: true,
  created_at: '2024-01-01T00:00:00.000Z',
  updated_at: '2024-06-01T00:00:00.000Z',
};

const mockGoalCompleted: Goal = {
  ...mockGoalInProgress,
  id: 'goal-2',
  title: '達成済み目標',
  current_amount: 5000000, // target_amountと同じ → 100%
};

const mockGoalOverdue: Goal = {
  ...mockGoalInProgress,
  id: 'goal-3',
  title: '期限切れ目標',
  target_date: '2020-01-01T00:00:00.000Z', // 過去の日付
  current_amount: 1000000,
};

const mockInactiveGoal: Goal = {
  ...mockGoalInProgress,
  id: 'goal-4',
  is_active: false,
};

describe('GoalProgressTracker', () => {
  describe('空状態の表示', () => {
    it('goalsが空配列のとき空状態メッセージが表示される', () => {
      render(<GoalProgressTracker goals={[]} />);
      expect(screen.getByText('アクティブな目標がありません')).toBeInTheDocument();
    });

    it('アクティブでないgoalのみの場合に空状態メッセージが表示される', () => {
      render(<GoalProgressTracker goals={[mockInactiveGoal]} />);
      expect(screen.getByText('アクティブな目標がありません')).toBeInTheDocument();
    });
  });

  describe('進捗バーの表示', () => {
    it('アクティブな目標のタイトルが表示される', () => {
      render(<GoalProgressTracker goals={[mockGoalInProgress]} />);
      expect(screen.getByText('貯蓄目標')).toBeInTheDocument();
    });

    it('達成率が表示される', () => {
      render(<GoalProgressTracker goals={[mockGoalInProgress]} />);
      // 1000000 / 5000000 = 20%
      expect(screen.getByText('20%')).toBeInTheDocument();
    });

    it('目標タイプのラベルが表示される', () => {
      render(<GoalProgressTracker goals={[mockGoalInProgress]} />);
      expect(screen.getByText('貯蓄')).toBeInTheDocument();
    });

    it('現在/目標の金額が表示される', () => {
      render(<GoalProgressTracker goals={[mockGoalInProgress]} />);
      expect(screen.getByText(/現在 \/ 目標/)).toBeInTheDocument();
    });

    it('残り金額が表示される', () => {
      render(<GoalProgressTracker goals={[mockGoalInProgress]} />);
      expect(screen.getByText('残り金額')).toBeInTheDocument();
    });

    it('月間積立情報が表示される', () => {
      render(<GoalProgressTracker goals={[mockGoalInProgress]} />);
      expect(screen.getByText('月間積立額')).toBeInTheDocument();
    });
  });

  describe('ステータスバッジ', () => {
    it('達成済みの目標に「達成」バッジが表示される', () => {
      render(<GoalProgressTracker goals={[mockGoalCompleted]} />);
      expect(screen.getByText('✓ 達成')).toBeInTheDocument();
    });

    it('期限切れの目標に「期限切れ」バッジが表示される', () => {
      render(<GoalProgressTracker goals={[mockGoalOverdue]} />);
      expect(screen.getByText('! 期限切れ')).toBeInTheDocument();
    });
  });

  describe('クリックイベント', () => {
    it('目標カードをクリックすると onGoalClick が呼ばれる', async () => {
      const user = userEvent.setup();
      const mockOnClick = jest.fn();

      render(<GoalProgressTracker goals={[mockGoalInProgress]} onGoalClick={mockOnClick} />);

      await user.click(screen.getByText('貯蓄目標'));

      expect(mockOnClick).toHaveBeenCalledWith(mockGoalInProgress);
    });
  });

  describe('複数目標の表示', () => {
    it('複数のアクティブな目標が全て表示される', () => {
      render(<GoalProgressTracker goals={[mockGoalInProgress, mockGoalCompleted]} />);
      expect(screen.getByText('貯蓄目標')).toBeInTheDocument();
      expect(screen.getByText('達成済み目標')).toBeInTheDocument();
    });

    it('非アクティブな目標は表示されない', () => {
      render(
        <GoalProgressTracker goals={[mockGoalInProgress, mockInactiveGoal]} />
      );
      expect(screen.getByText('貯蓄目標')).toBeInTheDocument();
      // inactiveGoalはmockGoalInProgressと同じタイトルなので確認不可
      // 代わりに全カードの数を確認
      const cards = screen.getAllByText(/達成率/);
      expect(cards).toHaveLength(1); // activeGoalのみ
    });
  });
});
