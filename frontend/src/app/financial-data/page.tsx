'use client';

import { useState, useEffect, useRef } from 'react';
import { useFinancialData } from '@/lib/contexts/FinancialDataContext';
import { useUser } from '@/lib/hooks/useUser';
import FinancialInputForm from '@/components/FinancialInputForm';
import InvestmentSettingsForm from '@/components/InvestmentSettingsForm';
import LoadingSpinner from '@/components/LoadingSpinner';
import { financialDataAPI } from '@/lib/api-client';
import { generateCSVFromProfile, downloadCSVLocally } from '@/lib/utils/csvExport';
import type { FinancialProfile } from '@/types/api';
import type { InvestmentSettings } from '@/components/InvestmentSettingsForm';

export default function FinancialDataPage() {
  const { userId } = useUser();
  const { financialData, loading, error, fetchFinancialData, updateProfile } = useFinancialData();
  const [activeTab, setActiveTab] = useState<'basic' | 'investment'>('basic');
  const [successMessage, setSuccessMessage] = useState<string | null>(null);
  const [csvLoading, setCsvLoading] = useState(false);
  const [csvError, setCsvError] = useState<string | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

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

  // バックエンド経由でCSVをダウンロードする（GETで直接 text/csv を返す）
  const handleCSVDownloadFromBackend = async () => {
    if (!userId) return;
    setCsvLoading(true);
    setCsvError(null);
    try {
      const blob = await financialDataAPI.downloadCSV(userId);
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = 'financial_data.csv';
      a.style.display = 'none';
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      URL.revokeObjectURL(url);
    } catch (err) {
      setCsvError(err instanceof Error ? err.message : 'CSVダウンロードに失敗しました');
    } finally {
      setCsvLoading(false);
    }
  };

  // フロントエンドのみでCSVを生成してダウンロードする（Blob API）
  const handleCSVDownloadLocal = () => {
    if (!profile) return;
    const csvContent = generateCSVFromProfile(profile);
    downloadCSVLocally(csvContent, 'financial_data_local.csv');
  };

  // CSVファイルをアップロードして財務データを更新する
  const handleCSVUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file || !userId) return;

    setCsvLoading(true);
    setCsvError(null);
    try {
      await financialDataAPI.importCSV(userId, file);
      await fetchFinancialData(userId);
      setSuccessMessage('CSVから財務データをインポートしました');
      setTimeout(() => setSuccessMessage(null), 3000);
    } catch (err) {
      setCsvError(err instanceof Error ? err.message : 'CSVインポートに失敗しました');
    } finally {
      setCsvLoading(false);
      // ファイル入力をリセットして同じファイルの再選択を可能にする
      if (fileInputRef.current) fileInputRef.current.value = '';
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
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-2">財務データ管理</h1>
        <p className="text-gray-600 dark:text-gray-300">収入、支出、貯蓄状況を入力・管理して、正確な将来予測の基盤を作成します</p>
      </div>

      {/* Success Message */}
      {successMessage && (
        <div className="mb-6 bg-success-50 border border-success-200 text-success-700 px-4 py-3 rounded-lg">
          ✓ {successMessage}
        </div>
      )}

      {/* Error Message */}
      {error && (
        <div className="mb-6 bg-blue-50 border border-blue-200 text-blue-700 px-4 py-3 rounded-lg">
          <div className="flex items-start gap-3">
            <span className="text-xl">ℹ️</span>
            <div>
              <p className="font-medium">{error}</p>
              {error.includes('まだ作成されていません') && (
                <p className="text-sm text-blue-600 mt-1">下のフォームから最初のデータを入力してください。</p>
              )}
            </div>
          </div>
        </div>
      )}

      <div className="grid lg:grid-cols-2 gap-8">
        {/* Current Data Display */}
        <div className="card">
          <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-4">現在の財務状況</h2>
          {profile ? (
            <div className="space-y-4">
              <div className="flex justify-between items-center py-2 border-b border-gray-100 dark:border-gray-700">
                <span className="text-gray-600 dark:text-gray-300">月収</span>
                <span className="font-medium text-gray-900 dark:text-white">
                  ¥{(profile.monthly_income || 0).toLocaleString()}
                </span>
              </div>
              <div className="flex justify-between items-center py-2 border-b border-gray-100 dark:border-gray-700">
                <span className="text-gray-600 dark:text-gray-300">月間支出</span>
                <span className="font-medium text-gray-900 dark:text-white">¥{totalExpenses.toLocaleString()}</span>
              </div>
              <div className="flex justify-between items-center py-2 border-b border-gray-100 dark:border-gray-700">
                <span className="text-gray-600 dark:text-gray-300">月間純貯蓄</span>
                <span className={`font-medium ${netSavings >= 0 ? 'text-success-600' : 'text-error-600'}`}>
                  ¥{netSavings.toLocaleString()}
                </span>
              </div>
              <div className="flex justify-between items-center py-2">
                <span className="text-gray-600 dark:text-gray-300">総資産</span>
                <span className="font-medium text-gray-900 dark:text-white">¥{totalSavings.toLocaleString()}</span>
              </div>
              <div className="pt-4 border-t border-gray-200 dark:border-gray-700">
                <div className="flex justify-between items-center py-2">
                  <span className="text-gray-600 dark:text-gray-300">投資利回り</span>
                  <span className="font-medium text-gray-900 dark:text-white">{profile.investment_return}%</span>
                </div>
                <div className="flex justify-between items-center py-2">
                  <span className="text-gray-600 dark:text-gray-300">インフレ率</span>
                  <span className="font-medium text-gray-900 dark:text-white">{profile.inflation_rate}%</span>
                </div>
              </div>
            </div>
          ) : (
            <div className="text-center py-8 text-gray-500 dark:text-gray-400">
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
                  : 'bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300 hover:bg-gray-200 dark:bg-gray-700'
              }`}
            >
              基本情報
            </button>
            <button
              onClick={() => setActiveTab('investment')}
              className={`flex-1 px-4 py-2 rounded-lg font-medium transition-colors ${
                activeTab === 'investment'
                  ? 'bg-primary-500 text-white'
                  : 'bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300 hover:bg-gray-200 dark:bg-gray-700'
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
      {profile && profile.monthly_expenses && profile.monthly_expenses.length > 0 && (
        <div className="mt-8">
          <div className="card">
            <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-4">支出内訳</h2>
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
                      <span className="text-gray-600 dark:text-gray-300">{expense.category}</span>
                      <span className="font-medium text-gray-900 dark:text-white">
                        ¥{expense.amount.toLocaleString()} ({percentage.toFixed(0)}%)
                      </span>
                    </div>
                    <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2">
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

      {/* CSV Import/Export */}
      <div className="mt-8">
        <div className="card">
          <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-2">CSVエクスポート・インポート</h2>
          <p className="text-sm text-gray-500 dark:text-gray-400 mb-4">
            財務データをCSVファイルとして書き出し・読み込みできます。<br />
            「バックエンド」はGoサーバーでCSVを生成、「ローカル生成」はブラウザのみで生成します。
          </p>

          <div className="flex flex-wrap gap-3">
            {/* バックエンド経由ダウンロード */}
            <button
              onClick={handleCSVDownloadFromBackend}
              disabled={csvLoading || !userId || !profile}
              className="btn-primary px-4 py-2 rounded-lg font-medium disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {csvLoading ? 'ダウンロード中...' : 'CSVダウンロード（バックエンド）'}
            </button>

            {/* フロントエンド直接生成（Blob API 学習用） */}
            <button
              onClick={handleCSVDownloadLocal}
              disabled={!profile}
              className="px-4 py-2 rounded-lg font-medium bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-200 hover:bg-gray-200 dark:hover:bg-gray-600 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              CSVダウンロード（ローカル生成）
            </button>

            {/* アップロード */}
            <label className={`px-4 py-2 rounded-lg font-medium cursor-pointer ${csvLoading ? 'opacity-50 cursor-not-allowed' : 'bg-success-100 dark:bg-success-900 text-success-700 dark:text-success-300 hover:bg-success-200 dark:hover:bg-success-800'}`}>
              {csvLoading ? 'インポート中...' : 'CSVアップロード'}
              <input
                ref={fileInputRef}
                type="file"
                accept=".csv"
                onChange={handleCSVUpload}
                disabled={csvLoading || !userId}
                className="hidden"
              />
            </label>
          </div>

          {csvError && (
            <p className="mt-3 text-sm text-error-600 dark:text-error-400">{csvError}</p>
          )}

          <div className="mt-4 p-3 bg-gray-50 dark:bg-gray-800 rounded-lg text-xs text-gray-500 dark:text-gray-400">
            <p className="font-medium mb-1">CSVフォーマット（マルチセクション形式）</p>
            <pre className="whitespace-pre-wrap font-mono">{`# SECTION: PROFILE\nfield,value\nmonthly_income,300000\n...\n\n# SECTION: EXPENSES\ncategory,amount,description\n生活費,100000,\n...\n\n# SECTION: SAVINGS\ntype,amount,description\ndeposit,500000,普通預金`}</pre>
          </div>
        </div>
      </div>
    </div>
  );
}