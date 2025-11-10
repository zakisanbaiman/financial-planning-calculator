'use client';

import { useState, useEffect } from 'react';
import { useFinancialData } from '@/lib/contexts/FinancialDataContext';
import { useUser } from '@/lib/hooks/useUser';
import FinancialInputForm from '@/components/FinancialInputForm';
import InvestmentSettingsForm from '@/components/InvestmentSettingsForm';
import LoadingSpinner from '@/components/LoadingSpinner';
import type { FinancialProfile } from '@/types/api';
import type { InvestmentSettings } from '@/components/InvestmentSettingsForm';

export default function FinancialDataPage() {
  const { userId } = useUser();
  const { financialData, loading, error, fetchFinancialData, updateProfile } = useFinancialData();
  const [activeTab, setActiveTab] = useState<'basic' | 'investment'>('basic');
  const [successMessage, setSuccessMessage] = useState<string | null>(null);

  useEffect(() => {
    if (userId) {
      fetchFinancialData(userId).catch(() => {
        // エラーは context で処理済み
      });
    }
  }, [userId, fetchFinancialData]);

  const handleFinancialSubmit = async (data: FinancialProfile) => {
    if (!userId) return;

    try {
      await updateProfile(userId, data);
      setSuccessMessage('財務データを保存しました');
      setTimeout(() => setSuccessMessage(null), 3000);
    } catch (err) {
      // エラーは context で処理済み
    }
  };

  const handleInvestmentSubmit = async (data: InvestmentSettings) => {
    if (!userId || !financialData?.profile) return;

    try {
      const updatedProfile: FinancialProfile = {
        ...financialData.profile,
        investment_return: data.investment_return,
        inflation_rate: data.inflation_rate,
      };
      await updateProfile(userId, updatedProfile);
      setSuccessMessage('投資設定を保存しました');
      setTimeout(() => setSuccessMessage(null), 3000);
    } catch (err) {
      // エラーは context で処理済み
    }
  };

  if (loading && !financialData) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="flex justify-center items-center h-64">
          <LoadingSpinner size="lg" />
        </div>
      </div>
    );
  }

  const profile = financialData?.profile;
  const totalExpenses = profile?.monthly_expenses?.reduce((sum, e) => sum + e.amount, 0) || 0;
  const totalSavings = profile?.current_savings?.reduce((sum, s) => sum + s.amount, 0) || 0;
  const netSavings = (profile?.monthly_income || 0) - totalExpenses;

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">財務データ管理</h1>
        <p className="text-gray-600">収入、支出、貯蓄状況を入力・管理して、正確な将来予測の基盤を作成します</p>
      </div>

      {/* Success Message */}
      {successMessage && (
        <div className="mb-6 bg-success-50 border border-success-200 text-success-700 px-4 py-3 rounded-lg">
          ✓ {successMessage}
        </div>
      )}

      {/* Error Message */}
      {error && (
        <div className="mb-6 bg-error-50 border border-error-200 text-error-700 px-4 py-3 rounded-lg">
          ✕ {error}
        </div>
      )}

      <div className="grid lg:grid-cols-2 gap-8">
        {/* Current Data Display */}
        <div className="card">
          <h2 className="text-xl font-semibold text-gray-900 mb-4">現在の財務状況</h2>
          {profile ? (
            <div className="space-y-4">
              <div className="flex justify-between items-center py-2 border-b border-gray-100">
                <span className="text-gray-600">月収</span>
                <span className="font-medium text-gray-900">
                  ¥{profile.monthly_income.toLocaleString()}
                </span>
              </div>
              <div className="flex justify-between items-center py-2 border-b border-gray-100">
                <span className="text-gray-600">月間支出</span>
                <span className="font-medium text-gray-900">¥{totalExpenses.toLocaleString()}</span>
              </div>
              <div className="flex justify-between items-center py-2 border-b border-gray-100">
                <span className="text-gray-600">月間純貯蓄</span>
                <span className={`font-medium ${netSavings >= 0 ? 'text-success-600' : 'text-error-600'}`}>
                  ¥{netSavings.toLocaleString()}
                </span>
              </div>
              <div className="flex justify-between items-center py-2">
                <span className="text-gray-600">総資産</span>
                <span className="font-medium text-gray-900">¥{totalSavings.toLocaleString()}</span>
              </div>
              <div className="pt-4 border-t border-gray-200">
                <div className="flex justify-between items-center py-2">
                  <span className="text-gray-600">投資利回り</span>
                  <span className="font-medium text-gray-900">{profile.investment_return}%</span>
                </div>
                <div className="flex justify-between items-center py-2">
                  <span className="text-gray-600">インフレ率</span>
                  <span className="font-medium text-gray-900">{profile.inflation_rate}%</span>
                </div>
              </div>
            </div>
          ) : (
            <div className="text-center py-8 text-gray-500">
              <div className="text-4xl mb-2">📊</div>
              <p>データがありません</p>
              <p className="text-sm">右側のフォームから入力してください</p>
            </div>
          )}
        </div>

        {/* Input Forms */}
        <div>
          {/* Tab Navigation */}
          <div className="flex gap-2 mb-4">
            <button
              onClick={() => setActiveTab('basic')}
              className={`flex-1 px-4 py-2 rounded-lg font-medium transition-colors ${
                activeTab === 'basic'
                  ? 'bg-primary-500 text-white'
                  : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
              }`}
            >
              基本情報
            </button>
            <button
              onClick={() => setActiveTab('investment')}
              className={`flex-1 px-4 py-2 rounded-lg font-medium transition-colors ${
                activeTab === 'investment'
                  ? 'bg-primary-500 text-white'
                  : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
              }`}
            >
              投資設定
            </button>
          </div>

          {/* Forms */}
          {activeTab === 'basic' ? (
            <FinancialInputForm
              initialData={profile}
              onSubmit={handleFinancialSubmit}
              loading={loading}
            />
          ) : (
            <InvestmentSettingsForm
              initialData={
                profile
                  ? {
                      investment_return: profile.investment_return,
                      inflation_rate: profile.inflation_rate,
                    }
                  : undefined
              }
              onSubmit={handleInvestmentSubmit}
              loading={loading}
            />
          )}
        </div>
      </div>

      {/* Expense Breakdown */}
      {profile && profile.monthly_expenses.length > 0 && (
        <div className="mt-8">
          <div className="card">
            <h2 className="text-xl font-semibold text-gray-900 mb-4">支出内訳</h2>
            <div className="grid md:grid-cols-2 gap-6">
              {profile.monthly_expenses.map((expense, index) => {
                const percentage = totalExpenses > 0 ? (expense.amount / totalExpenses) * 100 : 0;
                const colors = [
                  'bg-primary-500',
                  'bg-success-500',
                  'bg-warning-500',
                  'bg-error-500',
                  'bg-purple-500',
                  'bg-indigo-500',
                ];
                const colorClass = colors[index % colors.length];

                return (
                  <div key={index} className="space-y-3">
                    <div className="flex justify-between items-center">
                      <span className="text-gray-600">{expense.category}</span>
                      <span className="font-medium text-gray-900">
                        ¥{expense.amount.toLocaleString()} ({percentage.toFixed(0)}%)
                      </span>
                    </div>
                    <div className="w-full bg-gray-200 rounded-full h-2">
                      <div
                        className={`${colorClass} h-2 rounded-full`}
                        style={{ width: `${percentage}%` }}
                      ></div>
                    </div>
                  </div>
                );
              })}
            </div>
          </div>
        </div>
      )}
    </div>
  );
}