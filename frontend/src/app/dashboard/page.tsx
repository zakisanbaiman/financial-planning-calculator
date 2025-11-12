'use client';

import React, { useEffect, useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useGoals } from '@/lib/contexts/GoalsContext';
import { useUser } from '@/lib/hooks/useUser';
import GoalProgressTracker from '@/components/GoalProgressTracker';
import GoalsSummaryChart from '@/components/GoalsSummaryChart';
import LoadingSpinner from '@/components/LoadingSpinner';
import type { Goal } from '@/types/api';

export default function DashboardPage() {
  const router = useRouter();
  const { userId } = useUser();
  const { goals, loading, fetchGoals } = useGoals();
  const [chartType, setChartType] = useState<'bar' | 'doughnut'>('bar');

  useEffect(() => {
    if (userId) {
      fetchGoals(userId);
    }
  }, [userId, fetchGoals]);

  const handleGoalClick = (goal: Goal) => {
    router.push('/goals');
  };

  const activeGoals = goals.filter((g) => g.is_active);
  const totalTarget = activeGoals.reduce((sum, g) => sum + g.target_amount, 0);
  const totalCurrent = activeGoals.reduce((sum, g) => sum + g.current_amount, 0);
  const overallProgress = totalTarget > 0 ? (totalCurrent / totalTarget) * 100 : 0;
  return (
    <div className="container mx-auto px-4 py-8">
      {/* Header */}
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">ダッシュボード</h1>
        <p className="text-gray-600">財務状況の概要と主要な指標を確認できます</p>
      </div>

      {/* Quick Stats */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <div className="card">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">月間純貯蓄</p>
              <p className="text-2xl font-bold text-gray-900">¥120,000</p>
            </div>
            <div className="text-2xl">💰</div>
          </div>
          <p className="text-xs text-gray-500 mt-2">前月比 +5%</p>
        </div>

        <div className="card">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">総資産</p>
              <p className="text-2xl font-bold text-gray-900">¥1,500,000</p>
            </div>
            <div className="text-2xl">📈</div>
          </div>
          <p className="text-xs text-gray-500 mt-2">前月比 +8%</p>
        </div>

        <div className="card">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">老後資金充足率</p>
              <p className="text-2xl font-bold text-gray-900">65%</p>
            </div>
            <div className="text-2xl">🏖️</div>
          </div>
          <p className="text-xs text-gray-500 mt-2">目標まで35%</p>
        </div>

        <div className="card">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">緊急資金</p>
              <p className="text-2xl font-bold text-gray-900">6ヶ月分</p>
            </div>
            <div className="text-2xl">🚨</div>
          </div>
          <p className="text-xs text-success-600 mt-2">十分確保済み</p>
        </div>
      </div>

      {/* Main Content Grid */}
      <div className="grid lg:grid-cols-3 gap-8">
        {/* Left Column - Charts and Projections */}
        <div className="lg:col-span-2 space-y-6">
          {/* Asset Projection Chart Placeholder */}
          <div className="card">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-xl font-semibold text-gray-900">資産推移予測</h2>
              <Link href="/calculations" className="text-primary-600 hover:text-primary-700 text-sm font-medium">
                詳細計算 →
              </Link>
            </div>
            <div className="h-64 bg-gray-100 rounded-lg flex items-center justify-center">
              <div className="text-center text-gray-500">
                <div className="text-4xl mb-2">📊</div>
                <p>資産推移グラフ</p>
                <p className="text-sm">(Chart.js実装予定)</p>
              </div>
            </div>
          </div>

          {/* Monthly Breakdown */}
          <div className="card">
            <h2 className="text-xl font-semibold text-gray-900 mb-4">月間収支内訳</h2>
            <div className="space-y-3">
              <div className="flex items-center justify-between py-2 border-b border-gray-100">
                <span className="text-gray-600">月収</span>
                <span className="font-medium text-gray-900">¥400,000</span>
              </div>
              <div className="flex items-center justify-between py-2 border-b border-gray-100">
                <span className="text-gray-600">住居費</span>
                <span className="font-medium text-gray-900">¥120,000</span>
              </div>
              <div className="flex items-center justify-between py-2 border-b border-gray-100">
                <span className="text-gray-600">食費</span>
                <span className="font-medium text-gray-900">¥60,000</span>
              </div>
              <div className="flex items-center justify-between py-2 border-b border-gray-100">
                <span className="text-gray-600">その他支出</span>
                <span className="font-medium text-gray-900">¥100,000</span>
              </div>
              <div className="flex items-center justify-between py-2 font-semibold">
                <span className="text-gray-900">純貯蓄</span>
                <span className="text-success-600">¥120,000</span>
              </div>
            </div>
          </div>
        </div>

        {/* Right Column - Goals and Actions */}
        <div className="space-y-6">
          {/* Active Goals */}
          <div className="card">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-xl font-semibold text-gray-900">進行中の目標</h2>
              <Link href="/goals" className="text-primary-600 hover:text-primary-700 text-sm font-medium">
                管理 →
              </Link>
            </div>
            {loading ? (
              <div className="flex justify-center py-8">
                <LoadingSpinner />
              </div>
            ) : activeGoals.length > 0 ? (
              <div className="space-y-4">
                {activeGoals.slice(0, 3).map((goal) => {
                  const progress = (goal.current_amount / goal.target_amount) * 100;
                  return (
                    <div key={goal.id}>
                      <div className="flex items-center justify-between mb-2">
                        <span className="text-sm font-medium text-gray-900">{goal.title}</span>
                        <span className="text-sm text-gray-600">{progress.toFixed(0)}%</span>
                      </div>
                      <div className="w-full bg-gray-200 rounded-full h-2">
                        <div
                          className={`h-2 rounded-full ${
                            progress >= 100
                              ? 'bg-success-500'
                              : progress >= 75
                              ? 'bg-primary-500'
                              : progress >= 50
                              ? 'bg-warning-500'
                              : 'bg-orange-500'
                          }`}
                          style={{ width: `${Math.min(progress, 100)}%` }}
                        ></div>
                      </div>
                      <p className="text-xs text-gray-500 mt-1">
                        ¥{goal.current_amount.toLocaleString()} / ¥{goal.target_amount.toLocaleString()}
                      </p>
                    </div>
                  );
                })}
                {activeGoals.length > 3 && (
                  <Link
                    href="/goals"
                    className="block text-center text-sm text-primary-600 hover:text-primary-700 font-medium mt-3"
                  >
                    他{activeGoals.length - 3}件の目標を表示 →
                  </Link>
                )}
              </div>
            ) : (
              <div className="text-center py-6">
                <p className="text-gray-500 text-sm mb-3">目標が設定されていません</p>
                <Link
                  href="/goals"
                  className="text-primary-600 hover:text-primary-700 text-sm font-medium"
                >
                  最初の目標を作成 →
                </Link>
              </div>
            )}
          </div>

          {/* Quick Actions */}
          <div className="card">
            <h2 className="text-xl font-semibold text-gray-900 mb-4">クイックアクション</h2>
            <div className="space-y-3">
              <Link
                href="/financial-data"
                className="block w-full text-left p-3 rounded-lg border border-gray-200 hover:border-primary-300 hover:bg-primary-50 transition-colors"
              >
                <div className="flex items-center space-x-3">
                  <span className="text-xl">💰</span>
                  <div>
                    <p className="font-medium text-gray-900">財務データ更新</p>
                    <p className="text-sm text-gray-600">収入・支出を更新</p>
                  </div>
                </div>
              </Link>

              <Link
                href="/goals"
                className="block w-full text-left p-3 rounded-lg border border-gray-200 hover:border-primary-300 hover:bg-primary-50 transition-colors"
              >
                <div className="flex items-center space-x-3">
                  <span className="text-xl">🎯</span>
                  <div>
                    <p className="font-medium text-gray-900">新しい目標設定</p>
                    <p className="text-sm text-gray-600">財務目標を追加</p>
                  </div>
                </div>
              </Link>

              <Link
                href="/reports"
                className="block w-full text-left p-3 rounded-lg border border-gray-200 hover:border-primary-300 hover:bg-primary-50 transition-colors"
              >
                <div className="flex items-center space-x-3">
                  <span className="text-xl">📋</span>
                  <div>
                    <p className="font-medium text-gray-900">レポート生成</p>
                    <p className="text-sm text-gray-600">PDF形式で出力</p>
                  </div>
                </div>
              </Link>
            </div>
          </div>

          {/* Recommendations */}
          <div className="card">
            <h2 className="text-xl font-semibold text-gray-900 mb-4">推奨事項</h2>
            <div className="space-y-3">
              <div className="p-3 bg-success-50 border border-success-200 rounded-lg">
                <p className="text-sm font-medium text-success-800">✅ 緊急資金は十分確保されています</p>
              </div>
              <div className="p-3 bg-warning-50 border border-warning-200 rounded-lg">
                <p className="text-sm font-medium text-warning-800">⚠️ 老後資金の積立を月額¥50,000増やすことを推奨</p>
              </div>
              <div className="p-3 bg-primary-50 border border-primary-200 rounded-lg">
                <p className="text-sm font-medium text-primary-800">💡 投資利回りを5%→6%に改善すると目標達成が2年早まります</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Goals Dashboard Section */}
      {activeGoals.length > 0 && (
        <div className="mt-8">
          <div className="flex justify-between items-center mb-6">
            <h2 className="text-2xl font-bold text-gray-900">目標進捗ダッシュボード</h2>
            <div className="flex gap-2">
              <button
                onClick={() => setChartType('bar')}
                className={`px-3 py-1 rounded text-sm font-medium transition-colors ${
                  chartType === 'bar'
                    ? 'bg-primary-500 text-white'
                    : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
                }`}
              >
                棒グラフ
              </button>
              <button
                onClick={() => setChartType('doughnut')}
                className={`px-3 py-1 rounded text-sm font-medium transition-colors ${
                  chartType === 'doughnut'
                    ? 'bg-primary-500 text-white'
                    : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
                }`}
              >
                円グラフ
              </button>
            </div>
          </div>

          <div className="grid lg:grid-cols-3 gap-8">
            {/* Progress Tracker */}
            <div className="lg:col-span-1">
              <h3 className="text-lg font-semibold text-gray-900 mb-4">目標一覧</h3>
              <GoalProgressTracker goals={goals} onGoalClick={handleGoalClick} />
            </div>

            {/* Summary Chart */}
            <div className="lg:col-span-2">
              <GoalsSummaryChart goals={goals} chartType={chartType} />
            </div>
          </div>
        </div>
      )}
    </div>
  );
}