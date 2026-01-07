'use client';

import { useState } from 'react';
import {
  AssetProjectionCalculator,
  RetirementCalculator,
  EmergencyFundCalculator,
} from '@/components';
import AssetProjectionChart from '@/components/AssetProjectionChart';
import { generateAssetProjections } from '@/lib/utils/projections';
import { formatCurrency } from '@/lib/utils/currency';

type CalculatorView = 'menu' | 'asset-projection' | 'retirement' | 'emergency';

// Sample data constants for asset projection preview
const SAMPLE_INITIAL_ASSETS = 3000000; // ¥3,000,000
const SAMPLE_MONTHLY_CONTRIBUTION = 120000; // ¥120,000
const SAMPLE_INVESTMENT_RETURN = 0.05; // 5% annual
const SAMPLE_INFLATION_RATE = 0.02; // 2% annual
const SAMPLE_PROJECTION_YEARS = 30;

export default function CalculationsPage() {
  const [activeView, setActiveView] = useState<CalculatorView>('menu');
  const userId = 'user-001'; // TODO: 実際のユーザーIDを取得

  // サンプル資産推移データを生成（30年間）
  const sampleProjections = generateAssetProjections(
    SAMPLE_PROJECTION_YEARS,
    SAMPLE_INITIAL_ASSETS,
    SAMPLE_MONTHLY_CONTRIBUTION,
    SAMPLE_INVESTMENT_RETURN,
    SAMPLE_INFLATION_RATE
  );

  if (activeView === 'asset-projection') {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="mb-6">
          <button
            onClick={() => setActiveView('menu')}
            className="text-blue-600 dark:text-blue-400 hover:text-blue-700 dark:text-blue-300 flex items-center gap-2 mb-4"
          >
            ← メニューに戻る
          </button>
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-2">資産推移シミュレーション</h1>
          <p className="text-gray-600 dark:text-gray-300">現在の貯蓄ペースで将来どれだけ資産が増えるかを計算します</p>
        </div>
        <AssetProjectionCalculator userId={userId} />
      </div>
    );
  }

  if (activeView === 'retirement') {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="mb-6">
          <button
            onClick={() => setActiveView('menu')}
            className="text-blue-600 dark:text-blue-400 hover:text-blue-700 dark:text-blue-300 flex items-center gap-2 mb-4"
          >
            ← メニューに戻る
          </button>
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-2">老後資金計算</h1>
          <p className="text-gray-600 dark:text-gray-300">退職後に必要な資金と年金額を考慮した計算を行います</p>
        </div>
        <RetirementCalculator userId={userId} />
      </div>
    );
  }

  if (activeView === 'emergency') {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="mb-6">
          <button
            onClick={() => setActiveView('menu')}
            className="text-blue-600 dark:text-blue-400 hover:text-blue-700 dark:text-blue-300 flex items-center gap-2 mb-4"
          >
            ← メニューに戻る
          </button>
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-2">緊急資金計算</h1>
          <p className="text-gray-600 dark:text-gray-300">万が一の時に必要な緊急資金を計算します</p>
        </div>
        <EmergencyFundCalculator userId={userId} />
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-2">財務計算機</h1>
        <p className="text-gray-600 dark:text-gray-300">資産推移、老後資金、緊急資金の計算と可視化を行います</p>
      </div>

      {/* Calculation Categories */}
      <div className="grid md:grid-cols-3 gap-6 mb-8">
        <div className="card">
          <div className="text-center">
            <div className="text-4xl mb-4">📈</div>
            <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-2">資産推移シミュレーション</h3>
            <p className="text-gray-600 dark:text-gray-300 text-sm mb-4">
              現在の貯蓄ペースで将来どれだけ資産が増えるかを計算
            </p>
            <button
              onClick={() => setActiveView('asset-projection')}
              className="btn-primary w-full"
            >
              計算開始
            </button>
          </div>
        </div>

        <div className="card">
          <div className="text-center">
            <div className="text-4xl mb-4">🏖️</div>
            <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-2">老後資金計算</h3>
            <p className="text-gray-600 dark:text-gray-300 text-sm mb-4">
              退職後に必要な資金と年金額を考慮した計算
            </p>
            <button
              onClick={() => setActiveView('retirement')}
              className="btn-primary w-full"
            >
              計算開始
            </button>
          </div>
        </div>

        <div className="card">
          <div className="text-center">
            <div className="text-4xl mb-4">🚨</div>
            <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-2">緊急資金計算</h3>
            <p className="text-gray-600 dark:text-gray-300 text-sm mb-4">
              万が一の時に必要な緊急資金を計算
            </p>
            <button
              onClick={() => setActiveView('emergency')}
              className="btn-primary w-full"
            >
              計算開始
            </button>
          </div>
        </div>
      </div>

      {/* Sample Calculation Results */}
      <div className="grid lg:grid-cols-2 gap-8">
        {/* Asset Projection Chart */}
        <div className="card">
          <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-4">資産推移予測（30年間）</h2>
          <AssetProjectionChart
            projections={sampleProjections}
            showRealValue={true}
            showContributions={true}
            height={256}
          />
          <div className="grid grid-cols-3 gap-4 text-center mt-4">
            <div>
              <p className="text-sm text-gray-600 dark:text-gray-300">10年後</p>
              <p className="text-lg font-semibold text-gray-900 dark:text-white">
                {formatCurrency(sampleProjections[10]?.total_assets || 0)}
              </p>
            </div>
            <div>
              <p className="text-sm text-gray-600 dark:text-gray-300">20年後</p>
              <p className="text-lg font-semibold text-gray-900 dark:text-white">
                {formatCurrency(sampleProjections[20]?.total_assets || 0)}
              </p>
            </div>
            <div>
              <p className="text-sm text-gray-600 dark:text-gray-300">30年後</p>
              <p className="text-lg font-semibold text-gray-900 dark:text-white">
                {formatCurrency(sampleProjections[30]?.total_assets || 0)}
              </p>
            </div>
          </div>
        </div>

        {/* Calculation Summary */}
        <div className="space-y-6">
          {/* Retirement Calculation */}
          <div className="card">
            <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-3">老後資金計算結果</h3>
            <div className="space-y-3">
              <div className="flex justify-between">
                <span className="text-gray-600 dark:text-gray-300">必要老後資金</span>
                <span className="font-medium text-gray-900 dark:text-white">¥30,000,000</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600 dark:text-gray-300">予想達成額</span>
                <span className="font-medium text-gray-900 dark:text-white">¥45,600,000</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600 dark:text-gray-300">余裕額</span>
                <span className="font-medium text-success-600">+¥15,600,000</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600 dark:text-gray-300">充足率</span>
                <span className="font-medium text-success-600">152%</span>
              </div>
            </div>
          </div>

          {/* Emergency Fund Calculation */}
          <div className="card">
            <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-3">緊急資金計算結果</h3>
            <div className="space-y-3">
              <div className="flex justify-between">
                <span className="text-gray-600 dark:text-gray-300">推奨緊急資金</span>
                <span className="font-medium text-gray-900 dark:text-white">¥1,680,000</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600 dark:text-gray-300">現在の緊急資金</span>
                <span className="font-medium text-gray-900 dark:text-white">¥600,000</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600 dark:text-gray-300">不足額</span>
                <span className="font-medium text-warning-600">¥1,080,000</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600 dark:text-gray-300">達成までの期間</span>
                <span className="font-medium text-gray-900 dark:text-white">9ヶ月</span>
              </div>
            </div>
          </div>

          {/* Calculation Parameters */}
          <div className="card">
            <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-3">計算パラメータ</h3>
            <div className="space-y-2 text-sm">
              <div className="flex justify-between">
                <span className="text-gray-600 dark:text-gray-300">投資利回り</span>
                <span className="text-gray-900 dark:text-white">5.0%</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600 dark:text-gray-300">インフレ率</span>
                <span className="text-gray-900 dark:text-white">2.0%</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600 dark:text-gray-300">月間貯蓄額</span>
                <span className="text-gray-900 dark:text-white">¥120,000</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600 dark:text-gray-300">退職予定年齢</span>
                <span className="text-gray-900 dark:text-white">65歳</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Calculation Forms Placeholder */}
      <div className="mt-8">
        <div className="card">
          <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-4">詳細計算フォーム</h2>
          <div className="text-center py-8 text-gray-500 dark:text-gray-400">
            <div className="text-4xl mb-2">🧮</div>
            <p>計算パラメータ入力フォーム</p>
            <p className="text-sm">(タスク8.1-8.3で実装予定)</p>
          </div>
        </div>
      </div>
    </div>
  );
}