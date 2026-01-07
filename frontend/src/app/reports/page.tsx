'use client';

import { useState } from 'react';
import { LoadingSpinner } from '@/components';
import AssetProjectionChart from '@/components/AssetProjectionChart';
import { generateAssetProjections } from '@/lib/utils/projections';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api';

// Sample data constants for report preview
const SAMPLE_INITIAL_ASSETS = 1500000; // ¥1,500,000
const SAMPLE_MONTHLY_CONTRIBUTION = 120000; // ¥120,000
const SAMPLE_INVESTMENT_RETURN = 0.05; // 5% annual
const SAMPLE_INFLATION_RATE = 0.02; // 2% annual

export default function ReportsPage() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [reportSettings, setReportSettings] = useState({
    years: 10,
    includeFinancialSummary: true,
    includeAssetProjection: true,
    includeGoalsProgress: true,
    includeDetails: false,
    includeRecommendations: true,
    format: 'pdf',
  });

  const sampleProjections = generateAssetProjections(
    reportSettings.years,
    SAMPLE_INITIAL_ASSETS,
    SAMPLE_MONTHLY_CONTRIBUTION,
    SAMPLE_INVESTMENT_RETURN,
    SAMPLE_INFLATION_RATE
  );

  const handleGenerateReport = async (reportType: string) => {
    setLoading(true);
    setError(null);

    try {
      const userId = 'user-001'; // 実際の実装ではログインユーザーIDを使用

      let endpoint = '';
      let requestBody: any = { user_id: userId };

      switch (reportType) {
        case 'comprehensive':
          endpoint = '/reports/comprehensive';
          requestBody.years = reportSettings.years;
          break;
        case 'financial_summary':
          endpoint = '/reports/financial-summary';
          break;
        case 'goals_progress':
          endpoint = '/reports/goals-progress';
          break;
        case 'asset_projection':
          endpoint = '/reports/asset-projection';
          requestBody.years = reportSettings.years;
          break;
        default:
          throw new Error('サポートされていないレポートタイプです');
      }

      const response = await fetch(`${API_BASE_URL}${endpoint}`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(requestBody),
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new Error(errorData.error || 'レポートの生成に失敗しました');
      }

      const data = await response.json();

      // レポートデータを取得したら、PDFエクスポートを実行
      if (data && data.report) {
        const exportResponse = await fetch(`${API_BASE_URL}/reports/export`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            user_id: userId,
            report_type: reportType,
            format: reportSettings.format,
            report_data: data.report,
          }),
        });

        if (!exportResponse.ok) {
          throw new Error('PDFエクスポートに失敗しました');
        }

        const exportData = await exportResponse.json();
        if (exportData && exportData.download_url) {
          // ダウンロードURLを開く（実際の実装では実ファイルをダウンロード）
          alert(`レポートが生成されました: ${exportData.file_name}\nダウンロードURL: ${exportData.download_url}`);
        }
      }
    } catch (err: any) {
      console.error('レポート生成エラー:', err);
      setError(err.message || 'レポートの生成に失敗しました');
    } finally {
      setLoading(false);
    }
  };

  const handleQuickDownload = async (reportType: string) => {
    setLoading(true);
    setError(null);

    try {
      const userId = 'user-001';
      const url = `${API_BASE_URL}/reports/pdf?user_id=${userId}&report_type=${reportType}&years=${reportSettings.years}`;
      
      const response = await fetch(url);
      
      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new Error(errorData.error || 'レポートのダウンロードに失敗しました');
      }

      const data = await response.json();
      if (data && data.download_url) {
        alert(`レポートが生成されました: ${data.file_name}\nダウンロードURL: ${data.download_url}`);
      }
    } catch (err: any) {
      console.error('レポートダウンロードエラー:', err);
      setError(err.message || 'レポートのダウンロードに失敗しました');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-2">レポート生成</h1>
        <p className="text-gray-600 dark:text-gray-300">財務状況をまとめたレポートをPDF形式で生成・印刷できます</p>
      </div>

      {error && (
        <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg">
          <p className="text-red-800">{error}</p>
        </div>
      )}

      {loading && (
        <div className="mb-6 flex items-center justify-center p-8">
          <LoadingSpinner />
          <span className="ml-3 text-gray-600 dark:text-gray-300">レポートを生成中...</span>
        </div>
      )}

      {/* Report Generation Options */}
      <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6 mb-8">
        <div className="card">
          <div className="text-center">
            <div className="text-4xl mb-4">📊</div>
            <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-2">総合財務レポート</h3>
            <p className="text-gray-600 dark:text-gray-300 text-sm mb-4">
              現在の財務状況と将来予測を含む包括的なレポート
            </p>
            <button 
              className="btn-primary w-full"
              onClick={() => handleGenerateReport('comprehensive')}
              disabled={loading}
            >
              {loading ? '生成中...' : 'PDF生成'}
            </button>
          </div>
        </div>

        <div className="card">
          <div className="text-center">
            <div className="text-4xl mb-4">🎯</div>
            <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-2">目標進捗レポート</h3>
            <p className="text-gray-600 dark:text-gray-300 text-sm mb-4">
              設定した目標の進捗状況と達成予測
            </p>
            <button 
              className="btn-primary w-full"
              onClick={() => handleGenerateReport('goals_progress')}
              disabled={loading}
            >
              {loading ? '生成中...' : 'PDF生成'}
            </button>
          </div>
        </div>

        <div className="card">
          <div className="text-center">
            <div className="text-4xl mb-4">📈</div>
            <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-2">資産推移レポート</h3>
            <p className="text-gray-600 dark:text-gray-300 text-sm mb-4">
              資産の推移予測とシナリオ分析
            </p>
            <button 
              className="btn-primary w-full"
              onClick={() => handleGenerateReport('asset_projection')}
              disabled={loading}
            >
              {loading ? '生成中...' : 'PDF生成'}
            </button>
          </div>
        </div>
      </div>

      {/* Report Preview */}
      <div className="grid lg:grid-cols-3 gap-8">
        {/* Report Content Preview */}
        <div className="lg:col-span-2">
          <div className="card">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-xl font-semibold text-gray-900 dark:text-white">レポートプレビュー</h2>
              <div className="flex space-x-2">
                <button 
                  className="btn-primary text-sm"
                  onClick={() => handleQuickDownload('comprehensive')}
                  disabled={loading}
                >
                  PDF出力
                </button>
              </div>
            </div>
            
            {/* Mock Report Content */}
            <div className="bg-white border border-gray-200 dark:border-gray-700 rounded-lg p-6 min-h-96">
              <div className="text-center mb-6">
                <h1 className="text-2xl font-bold text-gray-900 dark:text-white mb-2">財務計画レポート</h1>
                <p className="text-gray-600 dark:text-gray-300">作成日: 2024年11月7日</p>
              </div>

              <div className="space-y-6">
                <section>
                  <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-3">現在の財務状況</h2>
                  <div className="grid grid-cols-2 gap-4">
                    <div className="bg-gray-50 p-3 rounded">
                      <p className="text-sm text-gray-600 dark:text-gray-300">月収</p>
                      <p className="text-lg font-semibold">¥400,000</p>
                    </div>
                    <div className="bg-gray-50 p-3 rounded">
                      <p className="text-sm text-gray-600 dark:text-gray-300">月間支出</p>
                      <p className="text-lg font-semibold">¥280,000</p>
                    </div>
                    <div className="bg-gray-50 p-3 rounded">
                      <p className="text-sm text-gray-600 dark:text-gray-300">月間貯蓄</p>
                      <p className="text-lg font-semibold text-success-600">¥120,000</p>
                    </div>
                    <div className="bg-gray-50 p-3 rounded">
                      <p className="text-sm text-gray-600 dark:text-gray-300">総資産</p>
                      <p className="text-lg font-semibold">¥1,500,000</p>
                    </div>
                  </div>
                </section>

                <section>
                  <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-3">将来予測</h2>
                  <AssetProjectionChart
                    projections={sampleProjections}
                    showRealValue={false}
                    showContributions={false}
                    height={128}
                  />
                </section>

                <section>
                  <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-3">目標達成状況</h2>
                  <div className="space-y-2">
                    <div className="flex justify-between items-center">
                      <span>緊急資金</span>
                      <span className="text-success-600 font-medium">100%</span>
                    </div>
                    <div className="flex justify-between items-center">
                      <span>老後資金</span>
                      <span className="text-primary-600 font-medium">65%</span>
                    </div>
                    <div className="flex justify-between items-center">
                      <span>マイホーム資金</span>
                      <span className="text-warning-600 font-medium">25%</span>
                    </div>
                  </div>
                </section>
              </div>
            </div>
          </div>
        </div>

        {/* Report Settings */}
        <div className="space-y-6">
          <div className="card">
            <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">レポート設定</h3>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  予測期間（年）
                </label>
                <input
                  type="number"
                  className="input-field"
                  value={reportSettings.years}
                  onChange={(e) => setReportSettings({ ...reportSettings, years: parseInt(e.target.value) || 10 })}
                  min="1"
                  max="50"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  含める項目
                </label>
                <div className="space-y-2">
                  <label className="flex items-center">
                    <input 
                      type="checkbox" 
                      className="mr-2" 
                      checked={reportSettings.includeFinancialSummary}
                      onChange={(e) => setReportSettings({ ...reportSettings, includeFinancialSummary: e.target.checked })}
                    />
                    <span className="text-sm">現在の財務状況</span>
                  </label>
                  <label className="flex items-center">
                    <input 
                      type="checkbox" 
                      className="mr-2" 
                      checked={reportSettings.includeAssetProjection}
                      onChange={(e) => setReportSettings({ ...reportSettings, includeAssetProjection: e.target.checked })}
                    />
                    <span className="text-sm">資産推移予測</span>
                  </label>
                  <label className="flex items-center">
                    <input 
                      type="checkbox" 
                      className="mr-2" 
                      checked={reportSettings.includeGoalsProgress}
                      onChange={(e) => setReportSettings({ ...reportSettings, includeGoalsProgress: e.target.checked })}
                    />
                    <span className="text-sm">目標進捗状況</span>
                  </label>
                  <label className="flex items-center">
                    <input 
                      type="checkbox" 
                      className="mr-2"
                      checked={reportSettings.includeDetails}
                      onChange={(e) => setReportSettings({ ...reportSettings, includeDetails: e.target.checked })}
                    />
                    <span className="text-sm">詳細な計算過程</span>
                  </label>
                  <label className="flex items-center">
                    <input 
                      type="checkbox" 
                      className="mr-2"
                      checked={reportSettings.includeRecommendations}
                      onChange={(e) => setReportSettings({ ...reportSettings, includeRecommendations: e.target.checked })}
                    />
                    <span className="text-sm">推奨事項</span>
                  </label>
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  出力形式
                </label>
                <select 
                  className="input-field"
                  value={reportSettings.format}
                  onChange={(e) => setReportSettings({ ...reportSettings, format: e.target.value })}
                >
                  <option value="pdf">PDF (推奨)</option>
                  <option value="json">JSON</option>
                </select>
              </div>
            </div>
          </div>

          <div className="card">
            <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">クイックアクション</h3>
            <div className="space-y-3">
              <button
                className="w-full p-3 bg-gray-50 hover:bg-gray-100 dark:bg-gray-700 rounded text-left transition-colors"
                onClick={() => handleQuickDownload('comprehensive')}
                disabled={loading}
              >
                <p className="text-sm font-medium text-gray-900 dark:text-white">📊 総合レポート</p>
                <p className="text-xs text-gray-600 dark:text-gray-300">すべての情報を含む包括的レポート</p>
              </button>
              <button
                className="w-full p-3 bg-gray-50 hover:bg-gray-100 dark:bg-gray-700 rounded text-left transition-colors"
                onClick={() => handleQuickDownload('financial_summary')}
                disabled={loading}
              >
                <p className="text-sm font-medium text-gray-900 dark:text-white">💰 財務サマリー</p>
                <p className="text-xs text-gray-600 dark:text-gray-300">現在の財務状況の概要</p>
              </button>
              <button
                className="w-full p-3 bg-gray-50 hover:bg-gray-100 dark:bg-gray-700 rounded text-left transition-colors"
                onClick={() => handleGenerateReport('asset_projection')}
                disabled={loading}
              >
                <p className="text-sm font-medium text-gray-900 dark:text-white">📈 資産推移</p>
                <p className="text-xs text-gray-600 dark:text-gray-300">将来の資産予測レポート</p>
              </button>
            </div>
          </div>
        </div>
      </div>

      {/* Report Information */}
      <div className="mt-8">
        <div className="card">
          <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-4">レポートについて</h2>
          <div className="space-y-4 text-gray-600 dark:text-gray-300">
            <div>
              <h3 className="font-semibold text-gray-900 dark:text-white mb-2">📊 総合財務レポート</h3>
              <p className="text-sm">
                現在の財務状況、資産推移予測、目標進捗、アクションプランを含む包括的なレポートです。
                財務計画の全体像を把握するのに最適です。
              </p>
            </div>
            <div>
              <h3 className="font-semibold text-gray-900 dark:text-white mb-2">💰 財務サマリーレポート</h3>
              <p className="text-sm">
                現在の財務健全性スコア、主要指標、推奨事項をまとめたレポートです。
                定期的な財務状況のチェックに便利です。
              </p>
            </div>
            <div>
              <h3 className="font-semibold text-gray-900 dark:text-white mb-2">🎯 目標進捗レポート</h3>
              <p className="text-sm">
                設定した目標の進捗状況、達成予測、推奨アクションを確認できます。
                目標管理とモチベーション維持に役立ちます。
              </p>
            </div>
            <div>
              <h3 className="font-semibold text-gray-900 dark:text-white mb-2">📈 資産推移レポート</h3>
              <p className="text-sm">
                将来の資産推移予測、シナリオ分析、投資戦略の洞察を提供します。
                長期的な資産形成計画の策定に活用できます。
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}