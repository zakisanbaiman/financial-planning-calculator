'use client';

import { Line } from 'react-chartjs-2';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler,
} from 'chart.js';

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler
);

export interface MonthlyCashFlowChartProps {
  monthlyIncome: number;
  monthlyExpenses: number;
  monthlySavings: number;
  months?: number;
  height?: number;
}

const MonthlyCashFlowChart: React.FC<MonthlyCashFlowChartProps> = ({
  monthlyIncome,
  monthlyExpenses,
  monthlySavings,
  months = 12,
  height = 300,
}) => {
  if (monthlyIncome === 0 && monthlyExpenses === 0) {
    return (
      <div className="h-64 bg-gray-100 dark:bg-gray-800 rounded-lg flex items-center justify-center">
        <div className="text-center text-gray-500 dark:text-gray-400">
          <div className="text-4xl mb-2">📈</div>
          <p>収支データがありません</p>
        </div>
      </div>
    );
  }

  // 過去から現在までの月を生成
  const labels = Array.from({ length: months }, (_, i) => {
    const monthsAgo = months - 1 - i;
    if (monthsAgo === 0) return '今月';
    return `${monthsAgo}ヶ月前`;
  });

  // 簡易的な変動を加えたデータ生成（実際のアプリでは履歴データを使用）
  const generateTrendData = (baseValue: number, variance: number = 0.1) => {
    return Array.from({ length: months }, (_, i) => {
      // 最後の月は実際の値を使用
      if (i === months - 1) return baseValue;
      // それ以前は±10%の範囲でランダムな値
      const randomFactor = 1 + (Math.random() - 0.5) * variance;
      return Math.round(baseValue * randomFactor);
    });
  };

  const incomeData = generateTrendData(monthlyIncome, 0.08);
  const expenseData = generateTrendData(monthlyExpenses, 0.12);
  const savingsData = incomeData.map((income, i) => income - expenseData[i]);

  const chartData = {
    labels,
    datasets: [
      {
        label: '収入',
        data: incomeData,
        borderColor: 'rgb(34, 197, 94)',
        backgroundColor: 'rgba(34, 197, 94, 0.1)',
        fill: false,
        tension: 0.4,
        pointRadius: 4,
        pointHoverRadius: 6,
      },
      {
        label: '支出',
        data: expenseData,
        borderColor: 'rgb(239, 68, 68)',
        backgroundColor: 'rgba(239, 68, 68, 0.1)',
        fill: false,
        tension: 0.4,
        pointRadius: 4,
        pointHoverRadius: 6,
      },
      {
        label: '純貯蓄',
        data: savingsData,
        borderColor: 'rgb(59, 130, 246)',
        backgroundColor: 'rgba(59, 130, 246, 0.1)',
        fill: true,
        tension: 0.4,
        pointRadius: 4,
        pointHoverRadius: 6,
      },
    ],
  };

  const options = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        position: 'top' as const,
        labels: {
          padding: 15,
          usePointStyle: true,
        },
      },
      title: {
        display: true,
        text: '月間収支推移',
        font: {
          size: 16,
          weight: 'bold' as const,
        },
      },
      tooltip: {
        callbacks: {
          label: function (context: any) {
            let label = context.dataset.label || '';
            if (label) {
              label += ': ';
            }
            if (context.parsed.y !== null) {
              label += '¥' + context.parsed.y.toLocaleString();
            }
            return label;
          },
        },
      },
    },
    scales: {
      y: {
        beginAtZero: true,
        ticks: {
          callback: function (value: any) {
            return '¥' + value.toLocaleString();
          },
        },
      },
      x: {
        grid: {
          display: false,
        },
      },
    },
    interaction: {
      mode: 'index' as const,
      intersect: false,
    },
  };

  // 平均値の計算
  const avgIncome = incomeData.reduce((sum, val) => sum + val, 0) / incomeData.length;
  const avgExpense = expenseData.reduce((sum, val) => sum + val, 0) / expenseData.length;
  const avgSavings = savingsData.reduce((sum, val) => sum + val, 0) / savingsData.length;

  return (
    <div className="space-y-4">
      <div style={{ height: `${height}px` }}>
        <Line data={chartData} options={options} />
      </div>
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
        <div className="p-3 bg-success-50 dark:bg-success-900/20 rounded-lg">
          <div className="text-sm text-success-700 dark:text-success-300 mb-1">平均収入</div>
          <div className="text-xl font-bold text-success-900 dark:text-success-100">
            ¥{Math.round(avgIncome).toLocaleString()}
          </div>
        </div>
        <div className="p-3 bg-red-50 dark:bg-red-900/20 rounded-lg">
          <div className="text-sm text-red-700 dark:text-red-300 mb-1">平均支出</div>
          <div className="text-xl font-bold text-red-900 dark:text-red-100">
            ¥{Math.round(avgExpense).toLocaleString()}
          </div>
        </div>
        <div className="p-3 bg-primary-50 dark:bg-primary-900/20 rounded-lg">
          <div className="text-sm text-primary-700 dark:text-primary-300 mb-1">平均貯蓄</div>
          <div className={`text-xl font-bold ${
            avgSavings >= 0 
              ? 'text-primary-900 dark:text-primary-100' 
              : 'text-red-900 dark:text-red-100'
          }`}>
            ¥{Math.round(avgSavings).toLocaleString()}
          </div>
        </div>
      </div>
    </div>
  );
};

export default MonthlyCashFlowChart;
