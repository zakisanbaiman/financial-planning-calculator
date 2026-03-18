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

  const formatCurrency = (value: number) => `¥${value.toLocaleString()}`;

  return (
    <div className="container mx-auto px-4 py-10">
      {/* Header */}
      <div className="mb-12">
        <p className="font-body text-sm font-semibold tracking-editorial uppercase text-accent-600 dark:text-accent-400 mb-2">
          Overview
        </p>
        <h1 className="font-display text-4xl md:text-5xl font-semibold text-ink-900 dark:text-ink-100">
          ダッシュボード
        </h1>
        <p className="font-body text-ink-500 dark:text-ink-400 mt-2">財務状況の概要と主要な指標を確認できます</p>
      </div>

      {/* Quick Stats — editorial table style */}
      <div className="grid grid-cols-2 lg:grid-cols-4 border border-ink-200 dark:border-ink-800 mb-12">
        {[
          {
            label: '月間純貯蓄',
            value: financialStats.hasData ? formatCurrency(financialStats.monthlySavings) : '---',
            sub: financialStats.hasData
              ? `収入 ${formatCurrency(financialStats.monthlyIncome)} - 支出 ${formatCurrency(financialStats.monthlyExpenses)}`
              : '財務データを入力してください',
          },
          {
            label: '総資産',
            value: financialStats.hasData ? formatCurrency(financialStats.totalAssets) : '---',
            sub: financialStats.hasData ? '現在の貯蓄合計' : '財務データを入力してください',
          },
          {
            label: '老後資金充足率',
            value: financialStats.hasData && financialStats.retirementSufficiency > 0
              ? `${financialStats.retirementSufficiency.toFixed(0)}%`
              : '---',
            sub: financialStats.hasData && financialStats.retirementSufficiency > 0
              ? `目標まで${(100 - financialStats.retirementSufficiency).toFixed(0)}%`
              : '退職データを入力してください',
          },
          {
            label: '緊急資金',
            value: financialStats.hasData && financialStats.emergencyMonths > 0
              ? `${financialStats.emergencyMonths.toFixed(1)}ヶ月分`
              : '---',
            sub: financialStats.hasData && financialStats.emergencyMonths > 0
              ? financialStats.emergencyMonths >= 6
                ? '十分確保済み'
                : financialStats.emergencyMonths >= 3
                  ? '最低限確保'
                  : '積み増しを推奨'
              : '財務データを入力してください',
            subColor: financialStats.emergencyMonths >= 6
              ? 'text-sage-600'
              : financialStats.emergencyMonths >= 3
                ? 'text-accent-600'
                : undefined,
          },
        ].map((stat, i) => (
          <div
            key={i}
            className={`p-6 bg-white dark:bg-ink-900 ${
              i < 3 ? 'border-r border-ink-200 dark:border-ink-800' : ''
            } ${i < 2 ? 'max-lg:border-b max-lg:border-ink-200 max-lg:dark:border-ink-800' : ''}`}
          >
            <p className="font-body text-xs font-semibold tracking-editorial uppercase text-ink-400 dark:text-ink-500 mb-2">
              {stat.label}
            </p>
            <p className="font-mono text-2xl md:text-3xl font-medium text-ink-900 dark:text-ink-100 mb-1">
              {stat.value}
            </p>
            <p className={`font-body text-xs ${stat.subColor || 'text-ink-400 dark:text-ink-500'}`}>
              {stat.sub}
            </p>
          </div>
        ))}
      </div>

      {/* Main Content Grid */}
      <div className="grid lg:grid-cols-3 gap-8">
        {/* Left Column - Charts and Projections */}
        <div className="lg:col-span-2 space-y-8">
          {/* Asset Projection Chart */}
          <div className="card">
            <div className="flex items-center justify-between mb-6">
              <h2 className="font-display text-2xl font-semibold text-ink-900 dark:text-ink-100">資産推移予測</h2>
              <div className="flex items-center gap-4">
                <div className="flex items-center gap-2">
                  <label htmlFor="projection-years" className="text-xs font-body text-ink-400 dark:text-ink-500 uppercase tracking-editorial">
                    期間
                  </label>
                  <select
                    id="projection-years"
                    value={projectionYears}
                    onChange={(e) => setProjectionYears(Number(e.target.value))}
                    className="px-2 py-1 border border-ink-300 dark:border-ink-700 text-sm bg-transparent font-body text-ink-800 dark:text-ink-200 focus:outline-none focus:border-ink-900 dark:focus:border-ink-200"
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
                <Link href="/calculations" className="text-sm font-body font-semibold text-ink-900 dark:text-ink-100 hover:text-accent-600 dark:hover:text-accent-400 transition-colors">
                  詳細計算 &rarr;
                </Link>
              </div>
            </div>
            <AssetProjectionChart
              projections={sampleProjections}
              showRealValue={true}
              showContributions={true}
              height={256}
            />
            <div className="mt-4 grid grid-cols-3 gap-4 border-t border-ink-200 dark:border-ink-800 pt-4">
              <div>
                <p className="font-body text-xs font-semibold tracking-editorial uppercase text-ink-400 dark:text-ink-500 mb-1">最終資産額</p>
                <p className="font-mono text-lg font-medium text-ink-900 dark:text-ink-100">
                  {formatCurrency(sampleProjections[sampleProjections.length - 1]?.total_assets || 0)}
                </p>
              </div>
              <div>
                <p className="font-body text-xs font-semibold tracking-editorial uppercase text-ink-400 dark:text-ink-500 mb-1">積立元本</p>
                <p className="font-mono text-lg font-medium text-ink-900 dark:text-ink-100">
                  {formatCurrency(sampleProjections[sampleProjections.length - 1]?.contributed_amount || 0)}
                </p>
              </div>
              <div>
                <p className="font-body text-xs font-semibold tracking-editorial uppercase text-ink-400 dark:text-ink-500 mb-1">投資収益</p>
                <p className="font-mono text-lg font-medium text-sage-700 dark:text-sage-400">
                  {formatCurrency(sampleProjections[sampleProjections.length - 1]?.investment_gains || 0)}
                </p>
              </div>
            </div>
          </div>

          {/* Monthly Breakdown */}
          <div className="card">
            <h2 className="font-display text-2xl font-semibold text-ink-900 dark:text-ink-100 mb-6">月間収支内訳</h2>
            {financialStats.hasData ? (
              <div>
                <div className="flex items-center justify-between py-3 border-b border-ink-200 dark:border-ink-800">
                  <span className="font-body text-sm text-ink-500 dark:text-ink-400">月収</span>
                  <span className="font-mono text-sm font-medium text-ink-900 dark:text-ink-100">
                    {formatCurrency(financialStats.monthlyIncome)}
                  </span>
                </div>
                {expenseBreakdown.map((expense, index) => (
                  <div key={index} className="flex items-center justify-between py-3 border-b border-ink-100 dark:border-ink-800/50">
                    <span className="font-body text-sm text-ink-500 dark:text-ink-400">{expense.category}</span>
                    <span className="font-mono text-sm text-ink-700 dark:text-ink-300">
                      {formatCurrency(expense.amount)}
                    </span>
                  </div>
                ))}
                <div className="flex items-center justify-between py-3 mt-1">
                  <span className="font-body text-sm font-semibold text-ink-900 dark:text-ink-100">純貯蓄</span>
                  <span className={`font-mono text-sm font-semibold ${
                    financialStats.monthlySavings >= 0 ? 'text-sage-700 dark:text-sage-400' : 'text-error-600 dark:text-error-400'
                  }`}>
                    {formatCurrency(financialStats.monthlySavings)}
                  </span>
                </div>
              </div>
            ) : (
              <div className="text-center py-8">
                <p className="font-body text-sm text-ink-400 dark:text-ink-500 mb-3">財務データが入力されていません</p>
                <Link
                  href="/financial-data"
                  className="font-body text-sm font-semibold text-ink-900 dark:text-ink-100 hover:text-accent-600 dark:hover:text-accent-400 transition-colors"
                >
                  財務データを入力 &rarr;
                </Link>
              </div>
            )}
          </div>
        </div>

        {/* Right Column - Goals and Actions */}
        <div className="space-y-8">
          {/* Active Goals */}
          <div className="card">
            <div className="flex items-center justify-between mb-6">
              <h2 className="font-display text-2xl font-semibold text-ink-900 dark:text-ink-100">進行中の目標</h2>
              <Link href="/goals" className="font-body text-sm font-semibold text-ink-900 dark:text-ink-100 hover:text-accent-600 dark:hover:text-accent-400 transition-colors">
                管理 &rarr;
              </Link>
            </div>
            {loading ? (
              <div className="flex justify-center py-8">
                <LoadingSpinner />
              </div>
            ) : activeGoals.length > 0 ? (
              <div className="space-y-5">
                {activeGoals.slice(0, 3).map((goal) => {
                  const progress = (goal.current_amount / goal.target_amount) * 100;
                  return (
                    <div key={goal.id}>
                      <div className="flex items-center justify-between mb-2">
                        <span className="font-body text-sm font-medium text-ink-800 dark:text-ink-200">{goal.title}</span>
                        <span className="font-mono text-xs text-ink-400 dark:text-ink-500">{progress.toFixed(0)}%</span>
                      </div>
                      <div className="w-full bg-ink-100 dark:bg-ink-800 h-1">
                        <div
                          className={`h-1 transition-all duration-500 ${
                            progress >= 100
                              ? 'bg-sage-500'
                              : progress >= 50
                              ? 'bg-ink-700 dark:bg-ink-300'
                              : 'bg-accent-500'
                          }`}
                          style={{ width: `${Math.min(progress, 100)}%` }}
                        />
                      </div>
                      <p className="font-mono text-xs text-ink-400 dark:text-ink-500 mt-1">
                        {formatCurrency(goal.current_amount)} / {formatCurrency(goal.target_amount)}
                      </p>
                    </div>
                  );
                })}
                {activeGoals.length > 3 && (
                  <Link
                    href="/goals"
                    className="block text-center font-body text-sm font-semibold text-ink-900 dark:text-ink-100 hover:text-accent-600 dark:hover:text-accent-400 transition-colors pt-2"
                  >
                    他{activeGoals.length - 3}件の目標を表示 &rarr;
                  </Link>
                )}
              </div>
            ) : (
              <div className="text-center py-8">
                <p className="font-body text-sm text-ink-400 dark:text-ink-500 mb-3">目標が設定されていません</p>
                <Link
                  href="/goals"
                  className="font-body text-sm font-semibold text-ink-900 dark:text-ink-100 hover:text-accent-600 dark:hover:text-accent-400 transition-colors"
                >
                  最初の目標を作成 &rarr;
                </Link>
              </div>
            )}
          </div>

          {/* Quick Actions */}
          <div className="card">
            <h2 className="font-display text-2xl font-semibold text-ink-900 dark:text-ink-100 mb-6">クイックアクション</h2>
            <div className="space-y-1">
              {[
                { href: '/financial-data', title: '財務データ更新', sub: '収入・支出を更新' },
                { href: '/goals', title: '新しい目標設定', sub: '財務目標を追加' },
                { href: '/reports', title: 'レポート生成', sub: 'PDF形式で出力' },
                { href: '/settings/security', title: 'セキュリティ設定', sub: '2段階認証の管理' },
              ].map((action) => (
                <Link
                  key={action.href}
                  href={action.href}
                  className="block p-3 -mx-1 hover:bg-ink-100 dark:hover:bg-ink-800 transition-colors group"
                >
                  <p className="font-body text-sm font-medium text-ink-800 dark:text-ink-200 group-hover:text-ink-900 dark:group-hover:text-ink-100">
                    {action.title}
                  </p>
                  <p className="font-body text-xs text-ink-400 dark:text-ink-500">
                    {action.sub}
                  </p>
                </Link>
              ))}
            </div>
          </div>

          {/* Recommendations */}
          <div className="card">
            <h2 className="font-display text-2xl font-semibold text-ink-900 dark:text-ink-100 mb-6">推奨事項</h2>
            <div className="space-y-3">
              {!financialStats.hasData ? (
                <div className="p-4 border-l-2 border-accent-500 bg-accent-50 dark:bg-accent-900/20">
                  <p className="font-body text-sm text-ink-700 dark:text-ink-300">💡 財務データを入力すると、パーソナライズされた推奨事項が表示されます</p>
                </div>
              ) : (
                <>
                  {financialStats.emergencyMonths >= 6 ? (
                    <div className="p-4 border-l-2 border-sage-500 bg-sage-50 dark:bg-sage-900/20">
                      <p className="font-body text-sm text-ink-700 dark:text-ink-300">✅ 緊急資金は十分確保されています（{financialStats.emergencyMonths.toFixed(1)}ケ月分）</p>
                    </div>
                  ) : financialStats.emergencyMonths >= 3 ? (
                    <div className="p-4 border-l-2 border-accent-500 bg-accent-50 dark:bg-accent-900/20">
                      <p className="font-body text-sm text-ink-700 dark:text-ink-300">⚠️ 緊急資金を6ケ月分まで増やすことを推奨（現在{financialStats.emergencyMonths.toFixed(1)}ケ月分）</p>
                    </div>
                  ) : (
                    <div className="p-4 border-l-2 border-error-500 bg-error-50 dark:bg-error-900/20">
                      <p className="font-body text-sm text-ink-700 dark:text-ink-300">🚨 緊急資金が不足しています。最低3ケ月分の確保を優先してください</p>
                    </div>
                  )}

                  {financialStats.retirementSufficiency > 0 && financialStats.retirementSufficiency < 100 && (
                    <div className="p-4 border-l-2 border-accent-500 bg-accent-50 dark:bg-accent-900/20">
                      <p className="font-body text-sm text-ink-700 dark:text-ink-300">
                        ⚠️ 老後資金の充足率は{financialStats.retirementSufficiency.toFixed(0)}%です。積立額の増額を検討してください
                      </p>
                    </div>
                  )}

                  {financialStats.monthlySavings <= 0 ? (
                    <div className="p-4 border-l-2 border-error-500 bg-error-50 dark:bg-error-900/20">
                      <p className="font-body text-sm text-ink-700 dark:text-ink-300">🚨 支出が収入を上回っています。支出の見直しを検討してください</p>
                    </div>
                  ) : financialStats.monthlySavings < financialStats.monthlyIncome * 0.2 && (
                    <div className="p-4 border-l-2 border-accent-500 bg-accent-50 dark:bg-accent-900/20">
                      <p className="font-body text-sm text-ink-700 dark:text-ink-300">💡 収入の20%以上を貯蓄に回すことで、目標達成が早まります</p>
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
        <div className="mt-16">
          <div className="flex justify-between items-end mb-8">
            <div>
              <p className="font-body text-sm font-semibold tracking-editorial uppercase text-accent-600 dark:text-accent-400 mb-2">
                Progress
              </p>
              <h2 className="font-display text-3xl font-semibold text-ink-900 dark:text-ink-100">目標進捗ダッシュボード</h2>
            </div>
            <div className="flex gap-1">
              <button
                onClick={() => setChartType('bar')}
                className={`px-3 py-1.5 text-xs font-body font-semibold tracking-editorial uppercase transition-colors ${
                  chartType === 'bar'
                    ? 'bg-ink-900 text-ink-50 dark:bg-ink-100 dark:text-ink-900'
                    : 'bg-transparent text-ink-500 hover:text-ink-800 dark:text-ink-400 dark:hover:text-ink-200 border border-ink-300 dark:border-ink-700'
                }`}
              >
                棒グラフ
              </button>
              <button
                onClick={() => setChartType('doughnut')}
                className={`px-3 py-1.5 text-xs font-body font-semibold tracking-editorial uppercase transition-colors ${
                  chartType === 'doughnut'
                    ? 'bg-ink-900 text-ink-50 dark:bg-ink-100 dark:text-ink-900'
                    : 'bg-transparent text-ink-500 hover:text-ink-800 dark:text-ink-400 dark:hover:text-ink-200 border border-ink-300 dark:border-ink-700'
                }`}
              >
                円グラフ
              </button>
            </div>
          </div>

          <div className="grid lg:grid-cols-3 gap-8">
            {/* Progress Tracker */}
            <div className="lg:col-span-1">
              <h3 className="font-display text-xl font-semibold text-ink-900 dark:text-ink-100 mb-4">目標一覧</h3>
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
