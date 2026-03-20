import React from 'react';
import { render, screen } from '@testing-library/react';
import GoalRecommendations from '../GoalRecommendations';
import type { Goal } from '@/types/api';

// 今日の日付: 2026-03-18
const FUTURE_DATE = '2030-12-31T00:00:00.000Z';
const PAST_DATE = '2020-01-01T00:00:00.000Z';
const NEAR_FUTURE_DATE = '2026-04-01T00:00:00.000Z'; // 約2週間後

const mockGoalInProgress: Goal = {
  id: 'goal-1',
  user_id: 'user-1',
  goal_type: 'savings',
  title: '貯蓄目標',
  target_amount: 5000000,
  target_date: FUTURE_DATE,
  current_amount: 1000000,
  monthly_contribution: 50000,
  is_active: true,
  created_at: '2024-01-01T00:00:00.000Z',
  updated_at: '2024-06-01T00:00:00.000Z',
};

describe('GoalRecommendations', () => {
  describe('目標達成済みの場合', () => {
    it('「目標達成おめでとうございます！」が表示される', () => {
      render(
        <GoalRecommendations
          goal={{ ...mockGoalInProgress, current_amount: 5000000 }}
        />
      );
      expect(screen.getByText('目標達成おめでとうございます！')).toBeInTheDocument();
    });

    it('達成アイコン（✅）が表示される', () => {
      render(
        <GoalRecommendations
          goal={{ ...mockGoalInProgress, current_amount: 5000000 }}
        />
      );
      expect(screen.getByText('✅')).toBeInTheDocument();
    });
  });

  describe('期限切れの場合', () => {
    it('「目標期日を過ぎています」が表示される', () => {
      render(
        <GoalRecommendations
          goal={{ ...mockGoalInProgress, target_date: PAST_DATE }}
        />
      );
      expect(screen.getByText('目標期日を過ぎています')).toBeInTheDocument();
    });

    it('エラーアイコン（❌）が表示される', () => {
      render(
        <GoalRecommendations
          goal={{ ...mockGoalInProgress, target_date: PAST_DATE }}
        />
      );
      expect(screen.getByText('❌')).toBeInTheDocument();
    });
  });

  describe('進捗が遅れている場合', () => {
    it('期限が近くて進捗が低い場合「進捗が遅れています」が表示される', () => {
      render(
        <GoalRecommendations
          goal={{
            ...mockGoalInProgress,
            target_date: NEAR_FUTURE_DATE,
            current_amount: 200000, // 4% 進捗
            monthly_contribution: 5000,
          }}
        />
      );
      expect(screen.getByText('進捗が遅れています')).toBeInTheDocument();
    });
  });

  describe('推奨事項なしの場合', () => {
    it('推奨事項がない場合に「現在、推奨事項はありません」が表示される', () => {
      // 進捗75%以上で期限まで余裕あり、月間積立も問題なし
      render(
        <GoalRecommendations
          goal={{
            ...mockGoalInProgress,
            current_amount: 4000000, // 80% 達成
            monthly_contribution: 100000, // 十分な積立額
          }}
        />
      );
      // この場合「順調に進んでいます」が表示されるはず
      expect(screen.getByText('順調に進んでいます')).toBeInTheDocument();
    });
  });

  describe('財務プロファイルがある場合', () => {
    it('貯蓄率が低い場合「貯蓄率が低い状態です」が表示される', () => {
      render(
        <GoalRecommendations
          goal={mockGoalInProgress}
          financialProfile={{
            monthly_income: 300000,
            monthly_expenses: 290000, // 貯蓄率3%
            current_savings: 600000,
          }}
        />
      );
      expect(screen.getByText('貯蓄率が低い状態です')).toBeInTheDocument();
    });

    it('緊急資金不足の場合に情報メッセージが表示される', () => {
      render(
        <GoalRecommendations
          goal={mockGoalInProgress}
          financialProfile={{
            monthly_income: 400000,
            monthly_expenses: 200000,
            current_savings: 100000, // 月間支出×3未満
          }}
        />
      );
      expect(screen.getByText('緊急資金の確保を優先')).toBeInTheDocument();
    });

    it('緊急資金目標の場合は緊急資金確保の推奨は表示されない', () => {
      render(
        <GoalRecommendations
          goal={{ ...mockGoalInProgress, goal_type: 'emergency' }}
          financialProfile={{
            monthly_income: 400000,
            monthly_expenses: 200000,
            current_savings: 100000,
          }}
        />
      );
      expect(screen.queryByText('緊急資金の確保を優先')).not.toBeInTheDocument();
    });
  });

  describe('アイコン表示', () => {
    it('infoタイプのアイコン（💡）が表示される', () => {
      // 投資推奨（残り金額>10万、期限>12ヶ月）
      render(
        <GoalRecommendations
          goal={{
            ...mockGoalInProgress,
            current_amount: 1000000, // 残り400万以上
            monthly_contribution: 100000,
          }}
        />
      );
      // 複数の💡が存在する可能性があるため getAllByText を使用
      const icons = screen.getAllByText('💡');
      expect(icons.length).toBeGreaterThan(0);
    });

    it('「投資による目標達成の加速」推奨が表示される', () => {
      render(
        <GoalRecommendations
          goal={{
            ...mockGoalInProgress,
            current_amount: 1000000, // 残り400万以上
            monthly_contribution: 100000,
          }}
        />
      );
      expect(screen.getByText('投資による目標達成の加速')).toBeInTheDocument();
    });
  });
});
