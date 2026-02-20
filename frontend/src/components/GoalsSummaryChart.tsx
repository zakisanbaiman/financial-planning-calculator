'use client';

import React, { useEffect, useRef } from 'react';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
  Tooltip,
  Legend,
  ArcElement,
} from 'chart.js';
import { Bar, Doughnut } from 'react-chartjs-2';
import type { Goal } from '@/types/api';

ChartJS.register(CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend, ArcElement);

export interface GoalsSummaryChartProps {
  goals: Goal[];
  chartType?: 'bar' | 'doughnut';
}

const GoalsSummaryChart: React.FC<GoalsSummaryChartProps> = ({ goals, chartType = 'bar' }) => {
  const activeGoals = goals.filter((g) => g.is_active);

  if (activeGoals.length === 0) {
    return (
      <div className="card text-center py-8">
        <p className="text-gray-500">表示する目標がありません</p>
      </div>
    );
  }

  // バーチャート用データ
  const barChartData = {
    labels: activeGoals.map((g) => g.title),
    datasets: [
      {
        label: '現在の積立額',
        data: activeGoals.map((g) => g.current_amount),
        backgroundColor: 'rgba(59, 130, 246, 0.8)',
        borderColor: 'rgba(59, 130, 246, 1)',
        borderWidth: 1,
      },
      {
        label: '目標金額',
        data: activeGoals.map((g) => g.target_amount),
        backgroundColor: 'rgba(209, 213, 219, 0.8)',
        borderColor: 'rgba(209, 213, 219, 1)',
        borderWidth: 1,
      },
    ],
  };

  const barChartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        position: 'top' as const,
      },
      title: {
        display: true,
        text: '目標別進捗状況',
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
    },
  };

  // ドーナツチャート用データ
  const totalTarget = activeGoals.reduce((sum, g) => sum + g.target_amount, 0);
  const totalCurrent = activeGoals.reduce((sum, g) => sum + g.current_amount, 0);
  const totalRemaining = Math.max(0, totalTarget - totalCurrent);

  const doughnutChartData = {
    labels: ['達成済み', '未達成'],
    datasets: [
      {
        data: [totalCurrent, totalRemaining],
        backgroundColor: ['rgba(34, 197, 94, 0.8)', 'rgba(209, 213, 219, 0.8)'],
        borderColor: ['rgba(34, 197, 94, 1)', 'rgba(209, 213, 219, 1)'],
        borderWidth: 1,
      },
    ],
  };

  const doughnutChartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        position: 'bottom' as const,
      },
      title: {
        display: true,
        text: '全体の達成状況',
        font: {
          size: 16,
          weight: 'bold' as const,
        },
      },
      tooltip: {
        callbacks: {
          label: function (context: any) {
            let label = context.label || '';
            if (label) {
              label += ': ';
            }
            if (context.parsed !== null) {
              label += '¥' + context.parsed.toLocaleString();
              const percentage = ((context.parsed / totalTarget) * 100).toFixed(1);
              label += ` (${percentage}%)`;
            }
            return label;
          },
        },
      },
    },
  };

  const overallProgress = totalTarget > 0 ? (totalCurrent / totalTarget) * 100 : 0;

  return (
    <div className="space-y-6">
      {/* 全体サマリー */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div className="card bg-primary-50 dark:bg-primary-900/30">
          <div className="text-sm text-primary-600 dark:text-primary-300 font-medium mb-1">総目標金額</div>
          <div className="text-2xl font-bold text-primary-900 dark:text-primary-200">
            ¥{totalTarget.toLocaleString()}
          </div>
        </div>
        <div className="card bg-success-50 dark:bg-success-900/30">
          <div className="text-sm text-success-600 dark:text-success-300 font-medium mb-1">現在の積立額</div>
          <div className="text-2xl font-bold text-success-900 dark:text-success-200">
            ¥{totalCurrent.toLocaleString()}
          </div>
        </div>
        <div className="card bg-orange-50 dark:bg-orange-900/30">
          <div className="text-sm text-orange-600 dark:text-orange-300 font-medium mb-1">残り金額</div>
          <div className="text-2xl font-bold text-orange-900 dark:text-orange-200">
            ¥{totalRemaining.toLocaleString()}
          </div>
        </div>
      </div>

      {/* 全体進捗バー */}
      <div className="card">
        <div className="flex justify-between items-center mb-2">
          <span className="text-sm font-medium text-gray-700 dark:text-gray-200">全体達成率</span>
          <span className="text-lg font-bold text-gray-900 dark:text-white">{overallProgress.toFixed(1)}%</span>
        </div>
        <div className="w-full bg-gray-200 dark:bg-gray-600 rounded-full h-4">
          <div
            className="bg-gradient-to-r from-primary-500 to-success-500 h-4 rounded-full transition-all"
            style={{ width: `${Math.min(overallProgress, 100)}%` }}
          />
        </div>
      </div>

      {/* チャート */}
      <div className="card">
        <div className="h-80">
          {chartType === 'bar' ? (
            <Bar data={barChartData} options={barChartOptions} />
          ) : (
            <Doughnut data={doughnutChartData} options={doughnutChartOptions} />
          )}
        </div>
      </div>

      {/* 目標別詳細 */}
      <div className="card">
        <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">目標別詳細</h3>
        <div className="space-y-3">
          {activeGoals.map((goal) => {
            const progress = (goal.current_amount / goal.target_amount) * 100;
            return (
              <div key={goal.id} className="border-b border-gray-200 dark:border-gray-600 pb-3 last:border-b-0">
                <div className="flex justify-between items-center mb-2">
                  <span className="font-medium text-gray-900 dark:text-white">{goal.title}</span>
                  <span className="text-sm font-semibold text-gray-700 dark:text-gray-200">
                    {progress.toFixed(1)}%
                  </span>
                </div>
                <div className="w-full bg-gray-200 dark:bg-gray-600 rounded-full h-2">
                  <div
                    className={`h-2 rounded-full transition-all ${
                      progress >= 100
                        ? 'bg-success-500'
                        : progress >= 75
                        ? 'bg-primary-500'
                        : progress >= 50
                        ? 'bg-warning-500'
                        : 'bg-orange-500'
                    }`}
                    style={{ width: `${Math.min(progress, 100)}%` }}
                  />
                </div>
                <div className="flex justify-between items-center mt-1 text-xs text-gray-500 dark:text-gray-400">
                  <span>¥{goal.current_amount.toLocaleString()}</span>
                  <span>¥{goal.target_amount.toLocaleString()}</span>
                </div>
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
};

export default GoalsSummaryChart;
