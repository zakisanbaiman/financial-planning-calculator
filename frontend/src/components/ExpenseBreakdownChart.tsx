'use client';

import { Doughnut } from 'react-chartjs-2';
import {
  Chart as ChartJS,
  ArcElement,
  Tooltip,
  Legend,
} from 'chart.js';
import type { ExpenseItem } from '@/types/api';

ChartJS.register(ArcElement, Tooltip, Legend);

export interface ExpenseBreakdownChartProps {
  expenses: ExpenseItem[];
  height?: number;
}

const ExpenseBreakdownChart: React.FC<ExpenseBreakdownChartProps> = ({ 
  expenses, 
  height = 300 
}) => {
  if (!expenses || expenses.length === 0) {
    return (
      <div className="h-64 bg-gray-100 dark:bg-gray-800 rounded-lg flex items-center justify-center">
        <div className="text-center text-gray-500 dark:text-gray-400">
          <div className="text-4xl mb-2">📊</div>
          <p>支出データがありません</p>
        </div>
      </div>
    );
  }

  // カテゴリごとに集計
  const categoryMap = new Map<string, number>();
  expenses.forEach((expense) => {
    const current = categoryMap.get(expense.category) || 0;
    categoryMap.set(expense.category, current + expense.amount);
  });

  const categories = Array.from(categoryMap.keys());
  const amounts = Array.from(categoryMap.values());

  // カラーパレット（支出用の温色系）
  const colors = [
    'rgba(239, 68, 68, 0.8)',   // 赤
    'rgba(249, 115, 22, 0.8)',  // オレンジ
    'rgba(234, 179, 8, 0.8)',   // 黄色
    'rgba(132, 204, 22, 0.8)',  // ライム
    'rgba(34, 197, 94, 0.8)',   // 緑
    'rgba(20, 184, 166, 0.8)',  // ティール
    'rgba(59, 130, 246, 0.8)',  // 青
    'rgba(168, 85, 247, 0.8)',  // 紫
    'rgba(236, 72, 153, 0.8)',  // ピンク
  ];

  const borderColors = [
    'rgba(239, 68, 68, 1)',
    'rgba(249, 115, 22, 1)',
    'rgba(234, 179, 8, 1)',
    'rgba(132, 204, 22, 1)',
    'rgba(34, 197, 94, 1)',
    'rgba(20, 184, 166, 1)',
    'rgba(59, 130, 246, 1)',
    'rgba(168, 85, 247, 1)',
    'rgba(236, 72, 153, 1)',
  ];

  const chartData = {
    labels: categories,
    datasets: [
      {
        data: amounts,
        backgroundColor: colors.slice(0, categories.length),
        borderColor: borderColors.slice(0, categories.length),
        borderWidth: 2,
      },
    ],
  };

  const options = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        position: 'right' as const,
        labels: {
          padding: 15,
          font: {
            size: 12,
          },
        },
      },
      tooltip: {
        callbacks: {
          label: function (context: any) {
            const label = context.label || '';
            const value = context.parsed || 0;
            const total = amounts.reduce((sum, val) => sum + val, 0);
            const percentage = total > 0 ? ((value / total) * 100).toFixed(1) : '0.0';
            return `${label}: ¥${value.toLocaleString()} (${percentage}%)`;
          },
        },
      },
    },
  };

  const totalExpenses = amounts.reduce((sum, amount) => sum + amount, 0);

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center">
        <h3 className="text-lg font-semibold text-gray-900 dark:text-white">支出内訳</h3>
        <div className="text-sm font-medium text-gray-600 dark:text-gray-300">
          合計: <span className="text-red-600 dark:text-red-400">¥{totalExpenses.toLocaleString()}</span>
        </div>
      </div>
      <div style={{ height: `${height}px` }}>
        <Doughnut data={chartData} options={options} />
      </div>
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-2 text-sm">
        {categories.map((category, index) => {
          const amount = amounts[index];
          const percentage = totalExpenses > 0 ? ((amount / totalExpenses) * 100).toFixed(1) : '0.0';
          return (
            <div 
              key={category} 
              className="flex items-center justify-between p-2 bg-gray-50 dark:bg-gray-800 rounded"
            >
              <div className="flex items-center gap-2">
                <div 
                  className="w-3 h-3 rounded-full" 
                  style={{ backgroundColor: colors[index % colors.length] }}
                />
                <span className="text-gray-700 dark:text-gray-300">{category}</span>
              </div>
              <div className="text-right">
                <div className="font-medium text-gray-900 dark:text-white">
                  ¥{amount.toLocaleString()}
                </div>
                <div className="text-xs text-gray-500 dark:text-gray-400">
                  {percentage}%
                </div>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
};

export default ExpenseBreakdownChart;
