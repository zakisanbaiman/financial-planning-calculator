'use client';

import React, { useEffect, useState, useMemo } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useGoals } from '@/lib/contexts/GoalsContext';
import { useFinancialData } from '@/lib/contexts/FinancialDataContext';
import { useUser } from '@/lib/hooks/useUser';
import GoalProgressTracker from '@/components/GoalProgressTracker';
import GoalsSummaryChart from '@/components/GoalsSummaryChart';
import AssetProjectionChart from '@/components/AssetProjectionChart';
import LoadingSpinner from '@/components/LoadingSpinner';
import type { Goal, AssetProjectionPoint } from '@/types/api';

export default function DashboardPage() {
  const router = useRouter();
  const { userId } = useUser();
  const { goals, loading: goalsLoading, fetchGoals } = useGoals();
  const { financialData, loading: financialLoading, fetchFinancialData } = useFinancialData();
  const [chartType, setChartType] = useState<'bar' | 'doughnut'>('bar');
  const [projectionYears, setProjectionYears] = useState<number>(30);

  useEffect(() => {
    if (userId) {
      fetchGoals(userId);
      fetchFinancialData(userId).catch(() => {
        // 財務データがまだない場合はエラーを無視
      });
    }
  }, [userId, fetchGoals, fetchFinancialData]);

  // 財務データから値を計算
  const financialStats = useMemo(() => {
    const profile = financialData?.profile;
    const retirement = financialData?.retirement;
    const emergencyFund = financialData?.emergency_fund;

    // 月収（配列の場合は合計）
    const monthlyIncome = profile?.monthly_income || 0;

    // 月間支出の合計
    const monthlyExpenses = profile?.monthly_expenses?.reduce((sum, item) => sum + item.amount, 0) || 0;

    // 月間純貯蓄
    const monthlySavings = monthlyIncome - monthlyExpenses;

    // 総資産（貯蓄の合計）
    const totalAssets = profile?.current_savings?.reduce((sum, item) => sum + item.amount, 0) || 0;

    // 老後資金充足率の計算
    let retirementSufficiency = 0;
    if (retirement) {
      const yearsInRetirement = retirement.life_expectancy - retirement.retirement_age;
      const monthsInRetirement = yearsInRetirement * 12;
      const requiredAmount = (retirement.monthly_retirement_expenses - retirement.pension_amount) * monthsInRetirement;
      if (requiredAmount > 0) {
        retirementSufficiency = Math.min((totalAssets / requiredAmount) * 100, 100);
      }
    }

    // 緊急資金の月数計算
    let emergencyMonths = 0;
    if (emergencyFund && emergencyFund.monthly_expenses > 0) {
      emergencyMonths = emergencyFund.current_amount / emergencyFund.monthly_expenses;
    } else if (monthlyExpenses > 0) {
      emergencyMonths = totalAssets / monthlyExpenses;
    }

    // 投資利回り・インフレ率
    const investmentReturn = (profile?.investment_return || 5) / 100;
    const inflationRate = (profile?.inflation_rate || 2) / 100;

    return {
      monthlyIncome,
      monthlyExpenses,
      monthlySavings,
      totalAssets,
      retirementSufficiency,
      emergencyMonths,
      investmentReturn,
      inflationRate,
      hasData: !!financialData,
    };
  }, [financialData]);

  // 支出の内訳を取得
  const expenseBreakdown = useMemo(() => {
    const expenses = financialData?.profile?.monthly_expenses || [];
    return expenses.map(item => ({
      category: item.category,
      amount: item.amount,
    }));
  }, [financialData]);

  const handleGoalClick = (goal: Goal) => {
    router.push('/goals');
  };

  const activeGoals = goals.filter((g) => g.is_active);
  const totalTarget = activeGoals.reduce((sum, g) => sum + g.target_amount, 0);
  const totalCurrent = activeGoals.reduce((sum, g) => sum + g.current_amount, 0);
  const overallProgress = totalTarget > 0 ? (totalCurrent / totalTarget) * 100 : 0;

  // 資産推移データを動的に生成（財務データを使用）
  const generateProjections = (years: number): AssetProjectionPoint[] => {
    const projections: AssetProjectionPoint[] = [];
    const initialAssets = financialStats.totalAssets || 3000000;
    const monthlyContribution = financialStats.monthlySavings > 0 ? financialStats.monthlySavings : 120000;
    const investmentReturn = financialStats.investmentReturn || 0.05;
    const inflationRate = financialStats.inflationRate || 0.02;
    
    for (let year = 0; year <= years; year++) {
      const contributedAmount = initialAssets + (monthlyContribution * 12 * year);
      const totalAssets = contributedAmount * Math.pow(1 + investmentReturn, year);
      const realValue = totalAssets / Math.pow(1 + inflationRate, year);
      const investmentGains = totalAssets - contributedAmount;
      
      projections.push({
        year,
        total_assets: Math.round(totalAssets),
        real_value: Math.round(realValue),
        contributed_amount: Math.round(contributedAmount),
        investment_gains: Math.round(investmentGains),
      });
    }
    
    return projections;
  };

  const sampleProjections = generateProjections(projectionYears);

  const loading = goalsLoading || financialLoading;

  return (
    <div className="container mx-auto px-4 py-8">
      {/* Header */}
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-2">ダッシュボード</h1>
        <p className="text-gray-600 dark:text-gray-300">財務状況の概要と主要な指標を確認できます</p>
      </div>

      {/* Quick Stats */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <div className="card">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600 dark:text-gray-300">月間純貯蓄</p>
              <p className="text-2xl font-bold text-gray-900 dark:text-white">
                {financialStats.hasData ? `¥${financialStats.monthlySavings.toLocaleString()}` : '---'}
              </p>
            </div>
            <div className="text-2xl">💰</div>
          </div>
          <p className="text-xs text-gray-500 dark:text-gray-400 mt-2">
            {financialStats.hasData ? `収入 ¥${financialStats.monthlyIncome.toLocaleString()} - 支出 ¥${financialStats.monthlyExpenses.toLocaleString()}` : '財務データを入力してください'}
          </p>
        </div>

        <div className="card">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600 dark:text-gray-300">総資産</p>
              <p className="text-2xl font-bold text-gray-900 dark:text-white">
                {financialStats.hasData ? `¥${financialStats.totalAssets.toLocaleString()}` : '---'}
              </p>
            </div>
            <div className="text-2xl">📈</div>
          </div>
          <p className="text-xs text-gray-500 dark:text-gray-400 mt-2">
            {financialStats.hasData ? '現在の貯蓄合計' : '財務データを入力してください'}
          </p>
        </div>

        <div className="card">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600 dark:text-gray-300">老後資金充足率</p>
              <p className="text-2xl font-bold text-gray-900 dark:text-white">
                {financialStats.hasData && financialStats.retirementSufficiency > 0 
                  ? `${financialStats.retirementSufficiency.toFixed(0)}%` 
                  : '---'}
              </p>
            </div>
            <div className="text-2xl">🏖️</div>
          </div>
          <p className="text-xs text-gray-500 dark:text-gray-400 mt-2">
            {financialStats.hasData && financialStats.retirementSufficiency > 0
              ? `目標まで${(100 - financialStats.retirementSufficiency).toFixed(0)}%`
              : '退職データを入力してください'}
          </p>
        </div>

        <div className="card">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600 dark:text-gray-300">緊急資金</p>
              <p className="text-2xl font-bold text-gray-900 dark:text-white">
                {financialStats.hasData && financialStats.emergencyMonths > 0
                  ? `${financialStats.emergencyMonths.toFixed(1)}ヶ月分`
                  : '---'}
              </p>
            </div>
            <div className="text-2xl">🚨</div>
          </div>
          <p className={`text-xs mt-2 ${
            financialStats.emergencyMonths >= 6 
              ? 'text-success-600' 
              : financialStats.emergencyMonths >= 3 
                ? 'text-warning-600' 
                : 'text-gray-500 dark:text-gray-400'
          }`}>
            {financialStats.hasData && financialStats.emergencyMonths > 0
              ? financialStats.emergencyMonths >= 6 
                ? '十分確保済み' 
                : financialStats.emergencyMonths >= 3
                  ? '最低限確保'
                  : '積み増しを推奨'
              : '財務データを入力してください'}
          </p>
        </div>
      </div>

      {/* Main Content Grid */}
      <div className="grid lg:grid-cols-3 gap-8">
        {/* Left Column - Charts and Projections */}
        <div className="lg:col-span-2 space-y-6">
          {/* Asset Projection Chart */}
          <div className="card">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-xl font-semibold text-gray-900 dark:text-white">資産推移予測</h2>
              <div className="flex items-center gap-3">
                <div className="flex items-center gap-2">
                  <label htmlFor="projection-years" className="text-sm text-gray-600 dark:text-gray-300">
                    期間:
                  </label>
                  <select
                    id="projection-years"
                    value={projectionYears}
                    onChange={(e) => setProjectionYears(Number(e.target.value))}
                    className="px-2 py-1 border border-gray-300 dark:border-gray-600 rounded text-sm bg-white dark:bg-gray-800 text-gray-900 dark:text-white"
                  >
                    <option value={5}>5年</option>
                    <option value={10}>10年</option>
                    <option value={20}>20年</option>
                    <option value={30}>30年</option>
                    <option value={50}>50年</option>
                    <option value={75}>75年</option>
                    <option value={100}>100年</option>
                  </select>
                </div>
                <Link href="/calculations" className="text-primary-600 hover:text-primary-700 text-sm font-medium">
                  詳細計算 →
                </Link>
              </div>
            </div>
            <AssetProjectionChart
              projections={sampleProjections}
              showRealValue={true}
              showContributions={true}
              height={256}
            />
            <div className="mt-3 grid grid-cols-3 gap-4 text-center text-sm">
              <div>
                <p className="text-gray-600 dark:text-gray-400">最終資産額</p>
                <p className="font-bold text-primary-600 dark:text-primary-400">
                  ¥{sampleProjections[sampleProjections.length - 1]?.total_assets.toLocaleString() || 0}
                </p>
              </div>
              <div>
                <p className="text-gray-600 dark:text-gray-400">積立元本</p>
                <p className="font-bold text-orange-600 dark:text-orange-400">
                  ¥{sampleProjections[sampleProjections.length - 1]?.contributed_amount.toLocaleString() || 0}
                </p>
              </div>
              <div>
                <p className="text-gray-600 dark:text-gray-400">投資収益</p>
                <p className="font-bold text-success-600 dark:text-success-400">
                  ¥{sampleProjections[sampleProjections.length - 1]?.investment_gains.toLocaleString() || 0}
                </p>
              </div>
            </div>
          </div>

          {/* Monthly Breakdown */}
          <div className="card">
            <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-4">月間収支内訳</h2>
            {financialStats.hasData ? (
              <div className="space-y-3">
                <div className="flex items-center justify-between py-2 border-b border-gray-100 dark:border-gray-700">
                  <span className="text-gray-600 dark:text-gray-300">月収</span>
                  <span className="font-medium text-gray-900 dark:text-white">¥{financialStats.monthlyIncome.toLocaleString()}</span>
                </div>
                {expenseBreakdown.map((expense, index) => (
                  <div key={index} className="flex items-center justify-between py-2 border-b border-gray-100 dark:border-gray-700">
                    <span className="text-gray-600 dark:text-gray-300">{expense.category}</span>
                    <span className="font-medium text-gray-900 dark:text-white">¥{expense.amount.toLocaleString()}</span>
                  </div>
                ))}
                <div className="flex items-center justify-between py-2 font-semibold">
                  <span className="text-gray-900 dark:text-white">純貯蓄</span>
                  <span className={financialStats.monthlySavings >= 0 ? 'text-success-600' : 'text-red-600'}>
                    ¥{financialStats.monthlySavings.toLocaleString()}
                  </span>
                </div>
              </div>
            ) : (
              <div className="text-center py-6">
                <p className="text-gray-500 dark:text-gray-400 text-sm mb-3">財務データが入力されていません</p>
                <Link
                  href="/financial-data"
                  className="text-primary-600 hover:text-primary-700 text-sm font-medium"
                >
                  財務データを入力 →
                </Link>
              </div>
            )}
          </div>
        </div>

        {/* Right Column - Goals and Actions */}
        <div className="space-y-6">
          {/* Active Goals */}
          <div className="card">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-xl font-semibold text-gray-900 dark:text-white">進行中の目標</h2>
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
                        <span className="text-sm font-medium text-gray-900 dark:text-white">{goal.title}</span>
                        <span className="text-sm text-gray-600 dark:text-gray-300">{progress.toFixed(0)}%</span>
                      </div>
                      <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2">
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
                      <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
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
                <p className="text-gray-500 dark:text-gray-400 text-sm mb-3">目標が設定されていません</p>
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
            <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-4">クイックアクション</h2>
            <div className="space-y-3">
              <Link
                href="/financial-data"
                className="block w-full text-left p-3 rounded-lg border border-gray-200 dark:border-gray-700 hover:border-primary-300 hover:bg-primary-50 dark:hover:bg-primary-900/30 transition-colors"
              >
                <div className="flex items-center space-x-3">
                  <span className="text-xl">💰</span>
                  <div>
                    <p className="font-medium text-gray-900 dark:text-white">財務データ更新</p>
                    <p className="text-sm text-gray-600 dark:text-gray-300">収入・支出を更新</p>
                  </div>
                </div>
              </Link>

              <Link
                href="/goals"
                className="block w-full text-left p-3 rounded-lg border border-gray-200 dark:border-gray-700 hover:border-primary-300 hover:bg-primary-50 dark:hover:bg-primary-900/30 transition-colors"
              >
                <div className="flex items-center space-x-3">
                  <span className="text-xl">🎯</span>
                  <div>
                    <p className="font-medium text-gray-900 dark:text-white">新しい目標設定</p>
                    <p className="text-sm text-gray-600 dark:text-gray-300">財務目標を追加</p>
                  </div>
                </div>
              </Link>

              <Link
                href="/reports"
                className="block w-full text-left p-3 rounded-lg border border-gray-200 dark:border-gray-700 hover:border-primary-300 hover:bg-primary-50 dark:hover:bg-primary-900/30 transition-colors"
              >
                <div className="flex items-center space-x-3">
                  <span className="text-xl">📋</span>
                  <div>
                    <p className="font-medium text-gray-900 dark:text-white">レポート生成</p>
                    <p className="text-sm text-gray-600 dark:text-gray-300">PDF形式で出力</p>
                  </div>
                </div>
              </Link>
            </div>
          </div>

          {/* Recommendations */}
          <div className="card">
            <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-4">推奨事項</h2>
            <div className="space-y-3">
              {!financialStats.hasData ? (
                <div className="p-3 bg-primary-50 dark:bg-primary-900/30 border border-primary-200 dark:border-primary-700 rounded-lg">
                  <p className="text-sm font-medium text-primary-800 dark:text-primary-200">💡 財務データを入力すると、パーソナライズされた推奨事項が表示されます</p>
                </div>
              ) : (
                <>
                  {financialStats.emergencyMonths >= 6 ? (
                    <div className="p-3 bg-success-50 dark:bg-success-900/30 border border-success-200 dark:border-success-700 rounded-lg">
                      <p className="text-sm font-medium text-success-800 dark:text-success-200">✅ 絊急資金は十分確保されています（{financialStats.emergencyMonths.toFixed(1)}ケ月分）</p>
                    </div>
                  ) : financialStats.emergencyMonths >= 3 ? (
                    <div className="p-3 bg-warning-50 dark:bg-warning-900/30 border border-warning-200 dark:border-warning-700 rounded-lg">
                      <p className="text-sm font-medium text-warning-800 dark:text-warning-200">⚠️ 絊急資金を6ケ月分まで増やすことを推奨（現在{financialStats.emergencyMonths.toFixed(1)}ケ月分）</p>
                    </div>
                  ) : (
                    <div className="p-3 bg-red-50 dark:bg-red-900/30 border border-red-200 dark:border-red-700 rounded-lg">
                      <p className="text-sm font-medium text-red-800 dark:text-red-200">🚨 絊急資金が不足しています。最低3ケ月分の確保を優先してください</p>
                    </div>
                  )}

                  {financialStats.retirementSufficiency > 0 && financialStats.retirementSufficiency < 100 && (
                    <div className="p-3 bg-warning-50 dark:bg-warning-900/30 border border-warning-200 dark:border-warning-700 rounded-lg">
                      <p className="text-sm font-medium text-warning-800 dark:text-warning-200">
                        ⚠️ 老後資金の充足率は{financialStats.retirementSufficiency.toFixed(0)}%です。積立額の増額を検討してください
                      </p>
                    </div>
                  )}

                  {financialStats.monthlySavings <= 0 ? (
                    <div className="p-3 bg-red-50 dark:bg-red-900/30 border border-red-200 dark:border-red-700 rounded-lg">
                      <p className="text-sm font-medium text-red-800 dark:text-red-200">🚨 支出が収入を上回っています。支出の見直しを検討してください</p>
                    </div>
                  ) : financialStats.monthlySavings < financialStats.monthlyIncome * 0.2 && (
                    <div className="p-3 bg-primary-50 dark:bg-primary-900/30 border border-primary-200 dark:border-primary-700 rounded-lg">
                      <p className="text-sm font-medium text-primary-800 dark:text-primary-200">💡 収入の20%以上を貯蓄に回すことで、目標達成が早まります</p>
                    </div>
                  )}
                </>
              )}
            </div>
          </div>
        </div>
      </div>

      {/* Goals Dashboard Section */}
      {activeGoals.length > 0 && (
        <div className="mt-8">
          <div className="flex justify-between items-center mb-6">
            <h2 className="text-2xl font-bold text-gray-900 dark:text-white">目標進捗ダッシュボード</h2>
            <div className="flex gap-2">
              <button
                onClick={() => setChartType('bar')}
                className={`px-3 py-1 rounded text-sm font-medium transition-colors ${
                  chartType === 'bar'
                    ? 'bg-primary-500 text-white'
                    : 'bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-200 hover:bg-gray-300 dark:hover:bg-gray-600'
                }`}
              >
                棒グラフ
              </button>
              <button
                onClick={() => setChartType('doughnut')}
                className={`px-3 py-1 rounded text-sm font-medium transition-colors ${
                  chartType === 'doughnut'
                    ? 'bg-primary-500 text-white'
                    : 'bg-gray-200 dark:bg-gray-700 text-gray-700 hover:bg-gray-300'
                }`}
              >
                円グラフ
              </button>
            </div>
          </div>

          <div className="grid lg:grid-cols-3 gap-8">
            {/* Progress Tracker */}
            <div className="lg:col-span-1">
              <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">目標一覧</h3>
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