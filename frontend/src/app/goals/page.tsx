'use client';

import React, { useEffect, useState } from 'react';
import { useGoals } from '@/lib/contexts/GoalsContext';
import { useUser } from '@/lib/hooks/useUser';
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

const goalTypeColors: Record<GoalType, string> = {
  savings: 'bg-blue-100 text-blue-800',
  retirement: 'bg-purple-100 text-purple-800',
  emergency: 'bg-orange-100 text-orange-800',
  custom: 'bg-gray-100 text-gray-800',
};

export default function GoalsPage() {
  const { userId } = useUser();
  const { goals, loading, error, fetchGoals, createGoal, updateGoal, deleteGoal, clearError } =
    useGoals();

  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const [editingGoal, setEditingGoal] = useState<Goal | null>(null);
  const [deletingGoalId, setDeletingGoalId] = useState<string | null>(null);

  useEffect(() => {
    if (userId) {
      fetchGoals(userId);
    }
  }, [userId, fetchGoals]);

  const handleCreateGoal = async (goal: Goal) => {
    try {
      await createGoal(goal);
      setIsCreateModalOpen(false);
    } catch (err) {
      console.error('Failed to create goal:', err);
    }
  };

  const handleUpdateGoal = async (goal: Goal) => {
    if (!editingGoal?.id) return;
    try {
      await updateGoal(editingGoal.id, userId, goal);
      setEditingGoal(null);
    } catch (err) {
      console.error('Failed to update goal:', err);
    }
  };

  const handleDeleteGoal = async (goalId: string) => {
    try {
      await deleteGoal(goalId, userId);
      setDeletingGoalId(null);
    } catch (err) {
      console.error('Failed to delete goal:', err);
    }
  };

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

  const activeGoals = goals.filter((g) => g.is_active);
  const inactiveGoals = goals.filter((g) => !g.is_active);

  if (loading && goals.length === 0) {
    return (
      <div className="flex justify-center items-center min-h-screen">
        <LoadingSpinner size="lg" />
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8 max-w-6xl">
      {/* ヘッダー */}
      <div className="flex justify-between items-center mb-8">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">目標管理</h1>
          <p className="text-gray-600 mt-2">財務目標を設定して進捗を追跡しましょう</p>
        </div>
        <Button onClick={() => setIsCreateModalOpen(true)}>+ 新しい目標</Button>
      </div>

      {/* エラー表示 */}
      {error && (
        <div className="mb-6 p-4 bg-error-50 border border-error-200 rounded-lg flex justify-between items-center">
          <p className="text-error-700">{error}</p>
          <button
            onClick={clearError}
            className="text-error-500 hover:text-error-700"
            aria-label="エラーを閉じる"
          >
            ✕
          </button>
        </div>
      )}

      {/* アクティブな目標 */}
      {activeGoals.length > 0 && (
        <div className="mb-8">
          <h2 className="text-xl font-semibold text-gray-900 mb-4">アクティブな目標</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {activeGoals.map((goal) => {
              const progress = calculateProgress(goal);
              const daysRemaining = calculateDaysRemaining(goal.target_date);
              const remainingAmount = Math.max(0, goal.target_amount - goal.current_amount);

              return (
                <div key={goal.id} className="card hover:shadow-lg transition-shadow">
                  {/* ヘッダー */}
                  <div className="flex justify-between items-start mb-4">
                    <div className="flex-1">
                      <div className="flex items-center gap-2 mb-2">
                        <span
                          className={`px-2 py-1 rounded text-xs font-medium ${
                            goalTypeColors[goal.type]
                          }`}
                        >
                          {goalTypeLabels[goal.type]}
                        </span>
                      </div>
                      <h3 className="text-lg font-semibold text-gray-900">{goal.title}</h3>
                    </div>
                    <div className="flex gap-2">
                      <button
                        onClick={() => setEditingGoal(goal)}
                        className="text-gray-400 hover:text-primary-500 transition-colors"
                        aria-label="編集"
                      >
                        <svg
                          className="w-5 h-5"
                          fill="none"
                          stroke="currentColor"
                          viewBox="0 0 24 24"
                        >
                          <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth={2}
                            d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
                          />
                        </svg>
                      </button>
                      <button
                        onClick={() => setDeletingGoalId(goal.id!)}
                        className="text-gray-400 hover:text-error-500 transition-colors"
                        aria-label="削除"
                      >
                        <svg
                          className="w-5 h-5"
                          fill="none"
                          stroke="currentColor"
                          viewBox="0 0 24 24"
                        >
                          <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth={2}
                            d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                          />
                        </svg>
                      </button>
                    </div>
                  </div>

                  {/* 進捗バー */}
                  <div className="mb-4">
                    <div className="flex justify-between items-center text-sm mb-2">
                      <span className="text-gray-600">進捗</span>
                      <span className="font-semibold text-gray-900">{progress.toFixed(1)}%</span>
                    </div>
                    <div className="w-full bg-gray-200 rounded-full h-3">
                      <div
                        className={`h-3 rounded-full transition-all ${
                          progress >= 100
                            ? 'bg-success-500'
                            : progress >= 75
                            ? 'bg-primary-500'
                            : progress >= 50
                            ? 'bg-warning-500'
                            : 'bg-orange-500'
                        }`}
                        style={{ width: `${progress}%` }}
                      />
                    </div>
                  </div>

                  {/* 金額情報 */}
                  <div className="space-y-2 mb-4">
                    <div className="flex justify-between text-sm">
                      <span className="text-gray-600">現在の積立額</span>
                      <span className="font-medium text-gray-900">
                        ¥{goal.current_amount.toLocaleString()}
                      </span>
                    </div>
                    <div className="flex justify-between text-sm">
                      <span className="text-gray-600">目標金額</span>
                      <span className="font-medium text-gray-900">
                        ¥{goal.target_amount.toLocaleString()}
                      </span>
                    </div>
                    <div className="flex justify-between text-sm">
                      <span className="text-gray-600">残り</span>
                      <span className="font-semibold text-primary-600">
                        ¥{remainingAmount.toLocaleString()}
                      </span>
                    </div>
                  </div>

                  {/* 期日情報 */}
                  <div className="pt-4 border-t border-gray-200">
                    <div className="flex justify-between items-center text-sm">
                      <span className="text-gray-600">目標期日</span>
                      <span className="font-medium text-gray-900">
                        {new Date(goal.target_date).toLocaleDateString('ja-JP')}
                      </span>
                    </div>
                    <div className="flex justify-between items-center text-sm mt-1">
                      <span className="text-gray-600">残り日数</span>
                      <span
                        className={`font-medium ${
                          daysRemaining < 30
                            ? 'text-error-600'
                            : daysRemaining < 90
                            ? 'text-warning-600'
                            : 'text-gray-900'
                        }`}
                      >
                        {daysRemaining > 0 ? `${daysRemaining}日` : '期限切れ'}
                      </span>
                    </div>
                    {goal.monthly_contribution > 0 && (
                      <div className="flex justify-between items-center text-sm mt-1">
                        <span className="text-gray-600">月間積立額</span>
                        <span className="font-medium text-gray-900">
                          ¥{goal.monthly_contribution.toLocaleString()}
                        </span>
                      </div>
                    )}
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      )}

      {/* 非アクティブな目標 */}
      {inactiveGoals.length > 0 && (
        <div>
          <h2 className="text-xl font-semibold text-gray-900 mb-4">非アクティブな目標</h2>
          <div className="space-y-3">
            {inactiveGoals.map((goal) => {
              const progress = calculateProgress(goal);

              return (
                <div
                  key={goal.id}
                  className="card bg-gray-50 hover:bg-gray-100 transition-colors"
                >
                  <div className="flex justify-between items-center">
                    <div className="flex-1">
                      <div className="flex items-center gap-3">
                        <span
                          className={`px-2 py-1 rounded text-xs font-medium ${
                            goalTypeColors[goal.type]
                          }`}
                        >
                          {goalTypeLabels[goal.type]}
                        </span>
                        <h3 className="text-base font-medium text-gray-700">{goal.title}</h3>
                        <span className="text-sm text-gray-500">
                          {progress.toFixed(0)}% 達成
                        </span>
                      </div>
                    </div>
                    <div className="flex gap-2">
                      <button
                        onClick={() => setEditingGoal(goal)}
                        className="text-gray-400 hover:text-primary-500 transition-colors"
                        aria-label="編集"
                      >
                        <svg
                          className="w-5 h-5"
                          fill="none"
                          stroke="currentColor"
                          viewBox="0 0 24 24"
                        >
                          <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth={2}
                            d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
                          />
                        </svg>
                      </button>
                      <button
                        onClick={() => setDeletingGoalId(goal.id!)}
                        className="text-gray-400 hover:text-error-500 transition-colors"
                        aria-label="削除"
                      >
                        <svg
                          className="w-5 h-5"
                          fill="none"
                          stroke="currentColor"
                          viewBox="0 0 24 24"
                        >
                          <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth={2}
                            d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                          />
                        </svg>
                      </button>
                    </div>
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      )}

      {/* 目標がない場合 */}
      {goals.length === 0 && !loading && (
        <div className="text-center py-12">
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
          <h3 className="mt-2 text-sm font-medium text-gray-900">目標がありません</h3>
          <p className="mt-1 text-sm text-gray-500">
            新しい財務目標を作成して、進捗を追跡しましょう
          </p>
          <div className="mt-6">
            <Button onClick={() => setIsCreateModalOpen(true)}>+ 最初の目標を作成</Button>
          </div>
        </div>
      )}

      {/* 作成モーダル */}
      <Modal
        isOpen={isCreateModalOpen}
        onClose={() => setIsCreateModalOpen(false)}
        title="新しい目標を作成"
        size="lg"
      >
        <GoalForm
          userId={userId}
          onSubmit={handleCreateGoal}
          onCancel={() => setIsCreateModalOpen(false)}
          loading={loading}
        />
      </Modal>

      {/* 編集モーダル */}
      <Modal
        isOpen={!!editingGoal}
        onClose={() => setEditingGoal(null)}
        title="目標を編集"
        size="lg"
      >
        {editingGoal && (
          <GoalForm
            initialData={editingGoal}
            userId={userId}
            onSubmit={handleUpdateGoal}
            onCancel={() => setEditingGoal(null)}
            loading={loading}
          />
        )}
      </Modal>

      {/* 削除確認モーダル */}
      <Modal
        isOpen={!!deletingGoalId}
        onClose={() => setDeletingGoalId(null)}
        title="目標を削除"
        size="sm"
      >
        <div className="space-y-4">
          <p className="text-gray-700">この目標を削除してもよろしいですか？</p>
          <p className="text-sm text-gray-500">この操作は取り消せません。</p>
          <div className="flex justify-end gap-3">
            <Button variant="outline" onClick={() => setDeletingGoalId(null)}>
              キャンセル
            </Button>
            <Button
              variant="error"
              onClick={() => deletingGoalId && handleDeleteGoal(deletingGoalId)}
              loading={loading}
            >
              削除
            </Button>
          </div>
        </div>
      </Modal>
    </div>
  );
}
