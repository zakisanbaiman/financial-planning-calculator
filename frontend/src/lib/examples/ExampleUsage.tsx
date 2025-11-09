'use client';

import React, { useEffect } from 'react';
import { useFinancialData, useGoals, useCalculations } from '@/lib/contexts';
import { useUser } from '@/lib/hooks';

/**
 * API クライアントと状態管理の使用例
 * このコンポーネントは実装の参考例として提供されています
 */
export function ExampleUsage() {
  const { userId, loading: userLoading } = useUser();
  const {
    financialData,
    loading: financialLoading,
    error: financialError,
    fetchFinancialData,
    updateProfile,
  } = useFinancialData();
  const {
    goals,
    loading: goalsLoading,
    error: goalsError,
    fetchGoals,
    createGoal,
  } = useGoals();
  const {
    assetProjection,
    loading: calculationLoading,
    error: calculationError,
    calculateAssetProjection,
  } = useCalculations();

  // ユーザーIDが取得できたら財務データと目標を取得
  useEffect(() => {
    if (userId) {
      fetchFinancialData(userId).catch(err => {
        console.error('財務データの取得に失敗:', err);
      });
      fetchGoals(userId).catch(err => {
        console.error('目標の取得に失敗:', err);
      });
    }
  }, [userId, fetchFinancialData, fetchGoals]);

  // 財務プロファイル更新の例
  const handleUpdateProfile = async () => {
    if (!userId) return;

    try {
      await updateProfile(userId, {
        monthly_income: 400000,
        monthly_expenses: [
          { category: '住居費', amount: 120000 },
          { category: '食費', amount: 60000 },
          { category: '交通費', amount: 20000 },
          { category: 'その他', amount: 80000 },
        ],
        current_savings: [
          { type: 'deposit', amount: 1000000 },
          { type: 'investment', amount: 500000 },
        ],
        investment_return: 5.0,
        inflation_rate: 2.0,
      });
      console.log('プロファイルを更新しました');
    } catch (error) {
      console.error('プロファイルの更新に失敗:', error);
    }
  };

  // 目標作成の例
  const handleCreateGoal = async () => {
    if (!userId) return;

    try {
      await createGoal({
        user_id: userId,
        type: 'savings',
        title: '新車購入',
        target_amount: 3000000,
        target_date: '2025-12-31',
        current_amount: 500000,
        monthly_contribution: 100000,
        is_active: true,
      });
      console.log('目標を作成しました');
    } catch (error) {
      console.error('目標の作成に失敗:', error);
    }
  };

  // 資産推移計算の例
  const handleCalculateProjection = async () => {
    if (!userId) return;

    try {
      await calculateAssetProjection({
        user_id: userId,
        years: 30,
        monthly_income: 400000,
        monthly_expenses: 280000,
        current_savings: 1500000,
        investment_return: 5.0,
        inflation_rate: 2.0,
      });
      console.log('資産推移を計算しました');
    } catch (error) {
      console.error('資産推移の計算に失敗:', error);
    }
  };

  if (userLoading) {
    return <div className="p-4">ユーザー情報を読み込み中...</div>;
  }

  return (
    <div className="p-4 space-y-6">
      <h1 className="text-2xl font-bold">API クライアントと状態管理の使用例</h1>

      {/* ユーザー情報 */}
      <section className="border p-4 rounded">
        <h2 className="text-xl font-semibold mb-2">ユーザー情報</h2>
        <p>ユーザーID: {userId}</p>
      </section>

      {/* 財務データ */}
      <section className="border p-4 rounded">
        <h2 className="text-xl font-semibold mb-2">財務データ</h2>
        {financialLoading && <p>読み込み中...</p>}
        {financialError && <p className="text-red-600">エラー: {financialError}</p>}
        {financialData && (
          <div>
            <p>月収: ¥{financialData.profile?.monthly_income?.toLocaleString()}</p>
            <p>投資利回り: {financialData.profile?.investment_return}%</p>
          </div>
        )}
        <button
          onClick={handleUpdateProfile}
          className="mt-2 px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
        >
          プロファイルを更新
        </button>
      </section>

      {/* 目標 */}
      <section className="border p-4 rounded">
        <h2 className="text-xl font-semibold mb-2">目標</h2>
        {goalsLoading && <p>読み込み中...</p>}
        {goalsError && <p className="text-red-600">エラー: {goalsError}</p>}
        {goals.length > 0 ? (
          <ul className="space-y-2">
            {goals.map((goal) => (
              <li key={goal.id} className="border-l-4 border-blue-500 pl-2">
                <p className="font-semibold">{goal.title}</p>
                <p>目標額: ¥{goal.target_amount.toLocaleString()}</p>
                <p>現在額: ¥{goal.current_amount.toLocaleString()}</p>
              </li>
            ))}
          </ul>
        ) : (
          <p>目標がありません</p>
        )}
        <button
          onClick={handleCreateGoal}
          className="mt-2 px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700"
        >
          目標を作成
        </button>
      </section>

      {/* 計算 */}
      <section className="border p-4 rounded">
        <h2 className="text-xl font-semibold mb-2">資産推移計算</h2>
        {calculationLoading && <p>計算中...</p>}
        {calculationError && <p className="text-red-600">エラー: {calculationError}</p>}
        {assetProjection && (
          <div>
            <p>最終資産額: ¥{assetProjection.final_amount.toLocaleString()}</p>
            <p>総積立額: ¥{assetProjection.total_contributions.toLocaleString()}</p>
            <p>投資収益: ¥{assetProjection.total_gains.toLocaleString()}</p>
            <p>データポイント数: {assetProjection.projections.length}</p>
          </div>
        )}
        <button
          onClick={handleCalculateProjection}
          className="mt-2 px-4 py-2 bg-purple-600 text-white rounded hover:bg-purple-700"
        >
          資産推移を計算
        </button>
      </section>
    </div>
  );
}
