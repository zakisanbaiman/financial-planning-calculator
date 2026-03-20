import React from 'react';
import { render, screen } from '@testing-library/react';
import AssetProjectionChart from '../AssetProjectionChart';
import type { AssetProjectionPoint } from '@/types/api';

// react-chartjs-2モック
jest.mock('react-chartjs-2', () => ({
  Line: (props: any) => (
    <div data-testid="mock-line-chart" data-chart-data={JSON.stringify(props.data)}>
      Line Chart
    </div>
  ),
}));

// chart.jsモック
jest.mock('chart.js', () => ({
  Chart: { register: jest.fn() },
  CategoryScale: jest.fn(),
  LinearScale: jest.fn(),
  PointElement: jest.fn(),
  LineElement: jest.fn(),
  Title: jest.fn(),
  Tooltip: jest.fn(),
  Legend: jest.fn(),
  Filler: jest.fn(),
}));

const mockProjections: AssetProjectionPoint[] = [
  { year: 1, total_assets: 4800000, real_value: 4700000, contributed_amount: 4800000, investment_gains: 0 },
  { year: 5, total_assets: 10000000, real_value: 9500000, contributed_amount: 9000000, investment_gains: 1000000 },
  { year: 10, total_assets: 20000000, real_value: 18000000, contributed_amount: 15000000, investment_gains: 5000000 },
];

describe('AssetProjectionChart', () => {
  describe('空データの表示', () => {
    it('dataがない場合に「データがありません」と表示される', () => {
      render(<AssetProjectionChart projections={[]} />);
      expect(screen.getByText('データがありません')).toBeInTheDocument();
    });

    it('nullが渡された場合に「データがありません」と表示される', () => {
      // TypeScript的にはnullは不正だがランタイムエラーを防ぐため
      render(<AssetProjectionChart projections={null as any} />);
      expect(screen.getByText('データがありません')).toBeInTheDocument();
    });
  });

  describe('チャート表示', () => {
    it('データがある場合にLineチャートが表示される', () => {
      render(<AssetProjectionChart projections={mockProjections} />);
      expect(screen.getByTestId('mock-line-chart')).toBeInTheDocument();
    });

    it('チャートのdata-chart-dataにラベルが含まれる', () => {
      render(<AssetProjectionChart projections={mockProjections} />);
      const chart = screen.getByTestId('mock-line-chart');
      const chartData = JSON.parse(chart.getAttribute('data-chart-data') || '{}');
      expect(chartData.labels).toContain('1年後');
      expect(chartData.labels).toContain('5年後');
      expect(chartData.labels).toContain('10年後');
    });

    it('デフォルトでデータセットに「総資産」が含まれる', () => {
      render(<AssetProjectionChart projections={mockProjections} />);
      const chart = screen.getByTestId('mock-line-chart');
      const chartData = JSON.parse(chart.getAttribute('data-chart-data') || '{}');
      const labels = chartData.datasets.map((d: any) => d.label);
      expect(labels).toContain('総資産');
    });

    it('showRealValue=trueのとき「実質価値」データセットが含まれる', () => {
      render(<AssetProjectionChart projections={mockProjections} showRealValue={true} />);
      const chart = screen.getByTestId('mock-line-chart');
      const chartData = JSON.parse(chart.getAttribute('data-chart-data') || '{}');
      const labels = chartData.datasets.map((d: any) => d.label);
      expect(labels).toContain('実質価値（インフレ調整後）');
    });

    it('showRealValue=falseのとき「実質価値」データセットが含まれない', () => {
      render(<AssetProjectionChart projections={mockProjections} showRealValue={false} />);
      const chart = screen.getByTestId('mock-line-chart');
      const chartData = JSON.parse(chart.getAttribute('data-chart-data') || '{}');
      const labels = chartData.datasets.map((d: any) => d.label);
      expect(labels).not.toContain('実質価値（インフレ調整後）');
    });

    it('showContributions=trueのとき「積立元本」データセットが含まれる', () => {
      render(<AssetProjectionChart projections={mockProjections} showContributions={true} />);
      const chart = screen.getByTestId('mock-line-chart');
      const chartData = JSON.parse(chart.getAttribute('data-chart-data') || '{}');
      const labels = chartData.datasets.map((d: any) => d.label);
      expect(labels).toContain('積立元本');
    });

    it('showContributions=falseのとき「積立元本」データセットが含まれない', () => {
      render(<AssetProjectionChart projections={mockProjections} showContributions={false} />);
      const chart = screen.getByTestId('mock-line-chart');
      const chartData = JSON.parse(chart.getAttribute('data-chart-data') || '{}');
      const labels = chartData.datasets.map((d: any) => d.label);
      expect(labels).not.toContain('積立元本');
    });
  });
});
