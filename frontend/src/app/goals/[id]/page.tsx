'use client';

import React, { useEffect, useState } from 'react';
import { useRouter, useParams } from 'next/navigation';
import { useGoals } from '@/lib/contexts/GoalsContext';
import { useFinancialData } from '@/lib/contexts/FinancialDataContext';
import { useUser } from '@/lib/hooks/useUser';
import GoalRecommendations from '@/components/GoalRecommendations';
import GoalForm from '@/components/GoalForm';
import Button from '@/components/Button';
import Modal from '@/components/Modal';
import LoadingSpinner from '@/components/LoadingSpinner';
import type { Goal, GoalType } from '@/types/api';

const goalTypeLabels: Record<GoalType, string> = {
  savings: '貯蓄',
  retirement: '老後資金',
  emergency: '緊急資金',
  custom: 'カスタム',
};

export default function GoalDetailPage() {
  const router = useRouter();
  const params = useParams();
  const goalId = typeof params.id === 'string' ? params.id : '';
  const { userId } = useUser();
  const { goals, loading: goalsLoading, updateGoal, updateGoalProgress } = useGoals();
  const { financialData, loading: financialLoading } = useFinancialData();

  const [isEditModalOpen, setIsEditModalOpen] = useState(false);
  const [isProgressModalOpen, setIsProgressModalOpen] = useState(false);
  const [newProgress, setNewProgress] = useState(0);

  const goal = goals.find((g) => g.id === goalId);

  useEffect(() => {
    if (goal) {
      setNewProgress(goal.current_amount);
    }
  }, [goal]);

  if (goalsLoading || financialLoading) {
    return (
      <div className="flex justify-center items-center min-h-screen">
        <LoadingSpinner size="lg" />
      </div>
    );
  }

  if (!goal) {
    return (
      <div className="container mx-auto px-4 py-8 max-w-4xl">
        <div className="text-center py-12">
          <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-4">目標が見つかりません</h2>
          <Button onClick={() => router.push('/goals')}>目標一覧に戻る</Button>
        </div>
      </div>
    );
  }

  const handleUpdateGoal = async (updatedGoal: Goal) => {
    const id = goal?.id || goalId;
    if (!id || id === '' || !userId) return;
    try {
      await updateGoal(id, userId, updatedGoal);
      setIsEditModalOpen(false);
    } catch (err) {
      console.error('Failed to update goal:', err);
    }
  };

  const handleUpdateProgress = async () => {
    const id = goal?.id || goalId;
    if (!id || id === '' || !userId) return;
    try {
      await updateGoalProgress(id, userId, newProgress);
      setIsProgressModalOpen(false);
    } catch (err) {
      console.error('Failed to update progress:', err);
    }
  };

  const progress = goal.target_amount > 0 ? (goal.current_amount / goal.target_amount) * 100 : 0;
  const remainingAmount = Math.max(0, goal.target_amount - goal.current_amount);
  const targetDate = new Date(goal.target_date);
  const today = new Date();
  const daysRemaining = Math.ceil((targetDate.getTime() - today.getTime()) / (1000 * 60 * 60 * 24));
  const monthsRemaining = Math.max(
    0,
    (targetDate.getFullYear() - today.getFullYear()) * 12 +
      (targetDate.getMonth() - today.getMonth())
  );

  const financialProfile = financialData?.profile
    ? {
        monthly_income: financialData.profile.monthly_income || 0,
        monthly_expenses: (financialData.profile.monthly_expenses || []).reduce(
          (sum, e) => sum + e.amount,
          0
        ),
        current_savings: (financialData.profile.current_savings || []).reduce((sum, s) => sum + s.amount, 0),
      }
    : undefined;

  return (
    <div className="container mx-auto px-4 py-8 max-w-6xl">
      {/* ヘッダー */}
      <div className="mb-6">
        <button
          onClick={() => router.push('/goals')}
          className="text-primary-600 hover:text-primary-700 mb-4 inline-flex items-center"
        >
          ← 目標一覧に戻る
        </button>
        <div className="flex justify-between items-start">
          <div>
            <div className="flex items-center gap-3 mb-2">
              <span className="px-3 py-1 bg-primary-100 text-primary-800 rounded text-sm font-medium">
                {goalTypeLabels[goal.goal_type]}
              </span>
              {!goal.is_active && (
                <span className="px-3 py-1 bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300 rounded text-sm font-medium">
                  非アクティブ
                </span>
              )}
            </div>
            <h1 className="text-3xl font-bold text-gray-900 dark:text-white">{goal.title}</h1>
          </div>
          <div className="flex gap-3">
            <Button variant="outline" onClick={() => setIsEditModalOpen(true)}>
              編集
            </Button>
            <Button onClick={() => setIsProgressModalOpen(true)}>進捗を更新</Button>
          </div>
        </div>
      </div>

      <div className="grid lg:grid-cols-3 gap-8">
        {/* 左カラム - 進捗情報 */}
        <div className="lg:col-span-2 space-y-6">
          {/* 進捗カード */}
          <div className="card">
            <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-4">進捗状況</h2>
            <div className="space-y-4">
              {/* 進捗率 */}
              <div>
                <div className="flex justify-between items-center mb-2">
                  <span className="text-sm font-medium text-gray-700">達成率</span>
                  <span className="text-2xl font-bold text-gray-900 dark:text-white">{progress.toFixed(1)}%</span>
                </div>
                <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-4">
                  <div
                    className={`h-4 rounded-full transition-all ${
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
              </div>

              {/* 金額情報 */}
              <div className="grid grid-cols-2 gap-4 pt-4">
                <div className="p-4 bg-gray-50 rounded-lg">
                  <div className="text-sm text-gray-600 dark:text-gray-300 mb-1">現在の積立額</div>
                  <div className="text-2xl font-bold text-gray-900 dark:text-white">
                    ¥{goal.current_amount.toLocaleString()}
                  </div>
                </div>
                <div className="p-4 bg-gray-50 rounded-lg">
                  <div className="text-sm text-gray-600 dark:text-gray-300 mb-1">目標金額</div>
                  <div className="text-2xl font-bold text-gray-900 dark:text-white">
                    ¥{goal.target_amount.toLocaleString()}
                  </div>
                </div>
                <div className="p-4 bg-primary-50 rounded-lg">
                  <div className="text-sm text-primary-600 mb-1">残り金額</div>
                  <div className="text-2xl font-bold text-primary-900">
                    ¥{remainingAmount.toLocaleString()}
                  </div>
                </div>
                <div className="p-4 bg-orange-50 rounded-lg">
                  <div className="text-sm text-orange-600 mb-1">月間積立額</div>
                  <div className="text-2xl font-bold text-orange-900">
                    ¥{goal.monthly_contribution.toLocaleString()}
                  </div>
                </div>
              </div>
            </div>
          </div>

          {/* 期日情報 */}
          <div className="card">
            <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-4">期日情報</h2>
            <div className="grid grid-cols-2 gap-4">
              <div>
                <div className="text-sm text-gray-600 dark:text-gray-300 mb-1">目標期日</div>
                <div className="text-lg font-semibold text-gray-900 dark:text-white">
                  {targetDate.toLocaleDateString('ja-JP', {
                    year: 'numeric',
                    month: 'long',
                    day: 'numeric',
                  })}
                </div>
              </div>
              <div>
                <div className="text-sm text-gray-600 dark:text-gray-300 mb-1">残り日数</div>
                <div
                  className={`text-lg font-semibold ${
                    daysRemaining < 0
                      ? 'text-error-600'
                      : daysRemaining < 30
                      ? 'text-warning-600'
                      : 'text-gray-900 dark:text-white'
                  }`}
                >
                  {daysRemaining > 0 ? `${daysRemaining}日` : '期限切れ'}
                </div>
              </div>
              {monthsRemaining > 0 && remainingAmount > 0 && goal.monthly_contribution > 0 && (
                <>
                  <div>
                    <div className="text-sm text-gray-600 dark:text-gray-300 mb-1">達成まで（現在のペース）</div>
                    <div className="text-lg font-semibold text-gray-900 dark:text-white">
                      約{Math.ceil(remainingAmount / goal.monthly_contribution)}ヶ月
                    </div>
                  </div>
                  <div>
                    <div className="text-sm text-gray-600 dark:text-gray-300 mb-1">推奨月間積立額</div>
                    <div className="text-lg font-semibold text-primary-600">
                      ¥{Math.ceil(remainingAmount / monthsRemaining).toLocaleString()}
                    </div>
                  </div>
                </>
              )}
            </div>
          </div>

          {/* 推奨事項 */}
          <div className="card">
            <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-4">推奨事項とアドバイス</h2>
            <GoalRecommendations goal={goal} financialProfile={financialProfile} />
          </div>
        </div>

        {/* 右カラム - 履歴とアクション */}
        <div className="space-y-6">
          {/* クイックアクション */}
          <div className="card">
            <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">クイックアクション</h3>
            <div className="space-y-3">
              <Button
                fullWidth
                onClick={() => setIsProgressModalOpen(true)}
              >
                進捗を更新
              </Button>
              <Button
                fullWidth
                variant="outline"
                onClick={() => setIsEditModalOpen(true)}
              >
                目標を編集
              </Button>
            </div>
          </div>

          {/* 目標情報 */}
          <div className="card">
            <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">目標情報</h3>
            <div className="space-y-3 text-sm">
              <div className="flex justify-between">
                <span className="text-gray-600 dark:text-gray-300">作成日</span>
                <span className="font-medium text-gray-900 dark:text-white">
                  {goal.created_at
                    ? new Date(goal.created_at).toLocaleDateString('ja-JP')
                    : '-'}
                </span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600 dark:text-gray-300">最終更新</span>
                <span className="font-medium text-gray-900 dark:text-white">
                  {goal.updated_at
                    ? new Date(goal.updated_at).toLocaleDateString('ja-JP')
                    : '-'}
                </span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600 dark:text-gray-300">ステータス</span>
                <span
                  className={`font-medium ${
                    goal.is_active ? 'text-success-600' : 'text-gray-500 dark:text-gray-400'
                  }`}
                >
                  {goal.is_active ? 'アクティブ' : '非アクティブ'}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* 編集モーダル */}
      <Modal
        isOpen={isEditModalOpen}
        onClose={() => setIsEditModalOpen(false)}
        title="目標を編集"
        size="lg"
      >
        {goal && userId && (
          <GoalForm
            initialData={goal}
            userId={userId}
            onSubmit={handleUpdateGoal}
            onCancel={() => setIsEditModalOpen(false)}
            loading={goalsLoading}
          />
        )}
      </Modal>

      {/* 進捗更新モーダル */}
      <Modal
        isOpen={isProgressModalOpen}
        onClose={() => setIsProgressModalOpen(false)}
        title="進捗を更新"
        size="md"
      >
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              現在の積立額
            </label>
            <input
              type="number"
              value={newProgress}
              onChange={(e) => setNewProgress(Number(e.target.value))}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500"
              min="0"
              step="1000"
            />
            <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
              目標金額: ¥{goal.target_amount.toLocaleString()}
            </p>
          </div>
          <div className="flex justify-end gap-3">
            <Button variant="outline" onClick={() => setIsProgressModalOpen(false)}>
              キャンセル
            </Button>
            <Button onClick={handleUpdateProgress} loading={goalsLoading}>
              更新
            </Button>
          </div>
        </div>
      </Modal>
    </div>
  );
}
