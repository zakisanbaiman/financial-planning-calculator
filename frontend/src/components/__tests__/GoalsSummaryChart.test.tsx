import React from 'react';
import { render, screen } from '@testing-library/react';
import GoalsSummaryChart from '../GoalsSummaryChart';
import type { Goal } from '@/types/api';

// react-chartjs-2モック
jest.mock('react-chartjs-2', () => ({
  Bar: (props: any) => (
    <div data-testid="mock-bar-chart" data-chart-data={JSON.stringify(props.data)}>
      Bar Chart
    </div>
  ),
  Doughnut: (props: any) => (
    <div data-testid="mock-doughnut-chart" data-chart-data={JSON.stringify(props.data)}>
      Doughnut Chart
    </div>
  ),
}));

// chart.jsモック
jest.mock('chart.js', () => ({
  Chart: { register: jest.fn() },
  CategoryScale: jest.fn(),
  LinearScale: jest.fn(),
  BarElement: jest.fn(),
  Title: jest.fn(),
  Tooltip: jest.fn(),
  Legend: jest.fn(),
  ArcElement: jest.fn(),
}));

const mockActiveGoal: Goal = {
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

const mockInactiveGoal: Goal = {
  ...mockActiveGoal,
  id: 'goal-2',
  is_active: false,
};

describe('GoalsSummaryChart', () => {
  describe('空状態の表示', () => {
    it('goalsが空配列のとき「表示する目標がありません」と表示される', () => {
      render(<GoalsSummaryChart goals={[]} />);
      expect(screen.getByText('表示する目標がありません')).toBeInTheDocument();
    });

    it('アクティブでないgoalのみの場合「表示する目標がありません」と表示される', () => {
      render(<GoalsSummaryChart goals={[mockInactiveGoal]} />);
      expect(screen.getByText('表示する目標がありません')).toBeInTheDocument();
    });
  });

  describe('チャート表示', () => {
    it('デフォルト（bar）でBarチャートが表示される', () => {
      render(<GoalsSummaryChart goals={[mockActiveGoal]} />);
      expect(screen.getByTestId('mock-bar-chart')).toBeInTheDocument();
    });

    it('chartType=doughnutでDoughnutチャートが表示される', () => {
      render(<GoalsSummaryChart goals={[mockActiveGoal]} chartType="doughnut" />);
      expect(screen.getByTestId('mock-doughnut-chart')).toBeInTheDocument();
    });

    it('chartType=barのときDoughnutチャートは表示されない', () => {
      render(<GoalsSummaryChart goals={[mockActiveGoal]} chartType="bar" />);
      expect(screen.queryByTestId('mock-doughnut-chart')).not.toBeInTheDocument();
    });

    it('chartType=doughnutのときBarチャートは表示されない', () => {
      render(<GoalsSummaryChart goals={[mockActiveGoal]} chartType="doughnut" />);
      expect(screen.queryByTestId('mock-bar-chart')).not.toBeInTheDocument();
    });
  });

  describe('サマリー表示', () => {
    it('総目標金額が表示される', () => {
      render(<GoalsSummaryChart goals={[mockActiveGoal]} />);
      expect(screen.getByText('総目標金額')).toBeInTheDocument();
    });

    it('現在の積立額が表示される', () => {
      render(<GoalsSummaryChart goals={[mockActiveGoal]} />);
      expect(screen.getByText('現在の積立額')).toBeInTheDocument();
    });

    it('残り金額が表示される', () => {
      render(<GoalsSummaryChart goals={[mockActiveGoal]} />);
      expect(screen.getByText('残り金額')).toBeInTheDocument();
    });

    it('全体達成率が表示される', () => {
      render(<GoalsSummaryChart goals={[mockActiveGoal]} />);
      expect(screen.getByText('全体達成率')).toBeInTheDocument();
    });

    it('目標別詳細が表示される', () => {
      render(<GoalsSummaryChart goals={[mockActiveGoal]} />);
      expect(screen.getByText('目標別詳細')).toBeInTheDocument();
    });

    it('アクティブな目標のタイトルが表示される', () => {
      render(<GoalsSummaryChart goals={[mockActiveGoal]} />);
      expect(screen.getByText('貯蓄目標')).toBeInTheDocument();
    });
  });

  describe('Barチャートのデータ', () => {
    it('barチャートに目標タイトルがラベルとして含まれる', () => {
      render(<GoalsSummaryChart goals={[mockActiveGoal]} />);
      const chart = screen.getByTestId('mock-bar-chart');
      const chartData = JSON.parse(chart.getAttribute('data-chart-data') || '{}');
      expect(chartData.labels).toContain('貯蓄目標');
    });
  });
});
