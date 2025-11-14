'use client';

import React from 'react';
import type { Goal, GoalType } from '@/types/api';

export interface GoalProgressTrackerProps {
  goals: Goal[];
  onGoalClick?: (goal: Goal) => void;
}

const goalTypeLabels: Record<GoalType, string> = {
  savings: '貯蓄',
  retirement: '老後資金',
  emergency: '緊急資金',
  custom: 'カスタム',
};

const goalTypeColors: Record<GoalType, { bg: string; text: string; progress: string }> = {
  savings: { bg: 'bg-blue-50', text: 'text-blue-700', progress: 'bg-blue-500' },
  retirement: { bg: 'bg-purple-50', text: 'text-purple-700', progress: 'bg-purple-500' },
  emergency: { bg: 'bg-orange-50', text: 'text-orange-700', progress: 'bg-orange-500' },
  custom: { bg: 'bg-gray-50', text: 'text-gray-700', progress: 'bg-gray-500' },
};

const GoalProgressTracker: React.FC<GoalProgressTrackerProps> = ({ goals, onGoalClick }) => {
  const activeGoals = goals.filter((g) => g.is_active);

  const calculateProgress = (goal: Goal): number => {
    if (goal.target_amount <= 0) return 0;
    return Math.min((goal.current_amount / goal.target_amount) * 100, 100);
  };

  const calculateDaysRemaining = (targetDate: string): number => {
    const target = new Date(targetDate);
    const today = new Date();
    const diffTime = target.getTime() - today.getTime();
    return Math.ceil(diffTime / (1000 * 60 * 60 * 24));
  };

  const getProgressStatus = (progress: number, daysRemaining: number): string => {
    if (progress >= 100) return 'completed';
    if (daysRemaining < 0) return 'overdue';
    if (daysRemaining < 30) return 'urgent';
    if (progress >= 75) return 'on-track';
    if (progress >= 50) return 'moderate';
    return 'behind';
  };

  const statusConfig: Record<string, { label: string; color: string; icon: string }> = {
    completed: { label: '達成', color: 'text-success-600', icon: '✓' },
    overdue: { label: '期限切れ', color: 'text-error-600', icon: '!' },
    urgent: { label: '緊急', color: 'text-warning-600', icon: '⚠' },
    'on-track': { label: '順調', color: 'text-success-600', icon: '↑' },
    moderate: { label: '進行中', color: 'text-primary-600', icon: '→' },
    behind: { label: '遅延', color: 'text-orange-600', icon: '↓' },
  };

  if (activeGoals.length === 0) {
    return (
      <div className="card text-center py-8">
        <svg
          className="mx-auto h-12 w-12 text-gray-400"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"
          />
        </svg>
        <h3 className="mt-2 text-sm font-medium text-gray-900">アクティブな目標がありません</h3>
        <p className="mt-1 text-sm text-gray-500">目標を作成して進捗を追跡しましょう</p>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {activeGoals.map((goal) => {
        const progress = calculateProgress(goal);
        const daysRemaining = calculateDaysRemaining(goal.target_date);
        const remainingAmount = Math.max(0, goal.target_amount - goal.current_amount);
        const status = getProgressStatus(progress, daysRemaining);
        const statusInfo = statusConfig[status];
        const colors = goalTypeColors[goal.type];

        return (
          <div
            key={goal.id}
            className={`card hover:shadow-md transition-all cursor-pointer ${colors.bg}`}
            onClick={() => onGoalClick?.(goal)}
          >
            {/* ヘッダー */}
            <div className="flex justify-between items-start mb-3">
              <div className="flex-1">
                <div className="flex items-center gap-2 mb-1">
                  <span className={`text-xs font-medium ${colors.text}`}>
                    {goalTypeLabels[goal.type]}
                  </span>
                  <span className={`text-xs font-medium ${statusInfo.color}`}>
                    {statusInfo.icon} {statusInfo.label}
                  </span>
                </div>
                <h3 className="text-base font-semibold text-gray-900">{goal.title}</h3>
              </div>
              <div className="text-right">
                <div className="text-2xl font-bold text-gray-900">{progress.toFixed(0)}%</div>
                <div className="text-xs text-gray-500">達成率</div>
              </div>
            </div>

            {/* 進捗バー */}
            <div className="mb-3">
              <div className="w-full bg-gray-200 rounded-full h-2.5">
                <div
                  className={`h-2.5 rounded-full transition-all ${colors.progress}`}
                  style={{ width: `${progress}%` }}
                />
              </div>
            </div>

            {/* 詳細情報 */}
            <div className="grid grid-cols-2 gap-3 text-sm">
              <div>
                <div className="text-gray-600 text-xs">現在 / 目標</div>
                <div className="font-semibold text-gray-900">
                  ¥{goal.current_amount.toLocaleString()} / ¥{goal.target_amount.toLocaleString()}
                </div>
              </div>
              <div>
                <div className="text-gray-600 text-xs">残り金額</div>
                <div className="font-semibold text-primary-600">
                  ¥{remainingAmount.toLocaleString()}
                </div>
              </div>
              <div>
                <div className="text-gray-600 text-xs">目標期日</div>
                <div className="font-medium text-gray-900">
                  {new Date(goal.target_date).toLocaleDateString('ja-JP', {
                    year: 'numeric',
                    month: 'short',
                    day: 'numeric',
                  })}
                </div>
              </div>
              <div>
                <div className="text-gray-600 text-xs">残り日数</div>
                <div
                  className={`font-medium ${
                    daysRemaining < 0
                      ? 'text-error-600'
                      : daysRemaining < 30
                      ? 'text-warning-600'
                      : 'text-gray-900'
                  }`}
                >
                  {daysRemaining > 0 ? `${daysRemaining}日` : '期限切れ'}
                </div>
              </div>
            </div>

            {/* 月間積立情報 */}
            {goal.monthly_contribution > 0 && (
              <div className="mt-3 pt-3 border-t border-gray-200">
                <div className="flex justify-between items-center text-sm">
                  <span className="text-gray-600">月間積立額</span>
                  <span className="font-semibold text-gray-900">
                    ¥{goal.monthly_contribution.toLocaleString()}
                  </span>
                </div>
                {remainingAmount > 0 && goal.monthly_contribution > 0 && (
                  <div className="flex justify-between items-center text-xs text-gray-500 mt-1">
                    <span>達成まで（現在のペース）</span>
                    <span>約{Math.ceil(remainingAmount / goal.monthly_contribution)}ヶ月</span>
                  </div>
                )}
              </div>
            )}
          </div>
        );
      })}
    </div>
  );
};

export default GoalProgressTracker;
