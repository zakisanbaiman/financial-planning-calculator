'use client';

import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { calculationsAPI } from '@/lib/api-client';
import type { AssetProjectionRequest, AssetProjectionResponse } from '@/types/api';
import AssetProjectionChart from './AssetProjectionChart';
import { InputField, Button, LoadingSpinner } from './index';

// バリデーションスキーマ
const assetProjectionSchema = z.object({
  years: z.number().min(1, '1年以上を指定してください').max(100, '100年以内で指定してください'),
  monthly_income: z.number().min(0, '0以上の値を入力してください'),
  monthly_expenses: z.number().min(0, '0以上の値を入力してください'),
  current_savings: z.number().min(0, '0以上の値を入力してください'),
  investment_return: z.number().min(0, '0以上の値を入力してください').max(100, '100以下の値を入力してください'),
  inflation_rate: z.number().min(0, '0以上の値を入力してください').max(50, '50以下の値を入力してください'),
});

type AssetProjectionFormData = z.infer<typeof assetProjectionSchema>;

interface AssetProjectionCalculatorProps {
  userId: string;
  initialData?: Partial<AssetProjectionFormData>;
}

export default function AssetProjectionCalculator({
  userId,
  initialData,
}: AssetProjectionCalculatorProps) {
  const [result, setResult] = useState<AssetProjectionResponse | null>(null);
  const [isCalculating, setIsCalculating] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const {
    register,
    handleSubmit,
    formState: { errors },
    watch,
  } = useForm<AssetProjectionFormData>({
    resolver: zodResolver(assetProjectionSchema),
    defaultValues: {
      years: initialData?.years || 30,
      monthly_income: initialData?.monthly_income || 400000,
      monthly_expenses: initialData?.monthly_expenses || 280000,
      current_savings: initialData?.current_savings || 1500000,
      investment_return: initialData?.investment_return || 5.0,
      inflation_rate: initialData?.inflation_rate || 2.0,
    },
  });

  const monthlyIncome = watch('monthly_income');
  const monthlyExpenses = watch('monthly_expenses');
  const monthlySavings = monthlyIncome - monthlyExpenses;

  const onSubmit = async (data: AssetProjectionFormData) => {
    setIsCalculating(true);
    setError(null);

    try {
      const request: AssetProjectionRequest = {
        user_id: userId,
        ...data,
      };

      const response = await calculationsAPI.assetProjection(request);
      setResult(response);
    } catch (err) {
      setError(err instanceof Error ? err.message : '計算中にエラーが発生しました');
    } finally {
      setIsCalculating(false);
    }
  };

  return (
    <div className="space-y-6">
      {/* 計算フォーム */}
      <div className="card">
        <h2 className="text-xl font-semibold text-gray-900 mb-4">
          資産推移シミュレーション
        </h2>
        <p className="text-gray-600 mb-6">
          現在の貯蓄ペースで将来どれだけ資産が増えるかを計算します
        </p>

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div className="grid md:grid-cols-2 gap-4">
            <InputField
              label="シミュレーション期間（年）"
              type="number"
              {...register('years', { valueAsNumber: true })}
              error={errors.years?.message}
              placeholder="30"
            />

            <InputField
              label="現在の貯蓄額（円）"
              type="number"
              {...register('current_savings', { valueAsNumber: true })}
              error={errors.current_savings?.message}
              placeholder="1500000"
            />

            <InputField
              label="月収（円）"
              type="number"
              {...register('monthly_income', { valueAsNumber: true })}
              error={errors.monthly_income?.message}
              placeholder="400000"
            />

            <InputField
              label="月間支出（円）"
              type="number"
              {...register('monthly_expenses', { valueAsNumber: true })}
              error={errors.monthly_expenses?.message}
              placeholder="280000"
            />

            <InputField
              label="投資利回り（%）"
              type="number"
              step="0.1"
              {...register('investment_return', { valueAsNumber: true })}
              error={errors.investment_return?.message}
              placeholder="5.0"
            />

            <InputField
              label="インフレ率（%）"
              type="number"
              step="0.1"
              {...register('inflation_rate', { valueAsNumber: true })}
              error={errors.inflation_rate?.message}
              placeholder="2.0"
            />
          </div>

          {/* 月間貯蓄額表示 */}
          <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
            <div className="flex justify-between items-center">
              <span className="text-gray-700 font-medium">月間貯蓄額</span>
              <span
                className={`text-xl font-bold ${
                  monthlySavings >= 0 ? 'text-success-600' : 'text-error-600'
                }`}
              >
                {new Intl.NumberFormat('ja-JP', {
                  style: 'currency',
                  currency: 'JPY',
                  maximumFractionDigits: 0,
                }).format(monthlySavings)}
              </span>
            </div>
            {monthlySavings < 0 && (
              <p className="text-sm text-error-600 mt-2">
                ⚠️ 支出が収入を上回っています
              </p>
            )}
          </div>

          {error && (
            <div className="bg-error-50 border border-error-200 rounded-lg p-4">
              <p className="text-error-600">{error}</p>
            </div>
          )}

          <Button type="submit" disabled={isCalculating} className="w-full">
            {isCalculating ? <LoadingSpinner size="sm" /> : '計算する'}
          </Button>
        </form>
      </div>

      {/* 計算結果 */}
      {result && (
        <>
          {/* グラフ表示 */}
          <div className="card">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">
              資産推移グラフ
            </h3>
            <AssetProjectionChart
              projections={result.projections}
              showRealValue={true}
              showContributions={true}
              height={400}
            />
          </div>

          {/* サマリー */}
          <div className="card">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">
              計算結果サマリー
            </h3>
            <div className="grid md:grid-cols-3 gap-6">
              <div className="text-center">
                <p className="text-sm text-gray-600 mb-1">最終資産額</p>
                <p className="text-2xl font-bold text-blue-600">
                  {new Intl.NumberFormat('ja-JP', {
                    style: 'currency',
                    currency: 'JPY',
                    maximumFractionDigits: 0,
                  }).format(result.final_amount)}
                </p>
              </div>
              <div className="text-center">
                <p className="text-sm text-gray-600 mb-1">積立元本合計</p>
                <p className="text-2xl font-bold text-orange-600">
                  {new Intl.NumberFormat('ja-JP', {
                    style: 'currency',
                    currency: 'JPY',
                    maximumFractionDigits: 0,
                  }).format(result.total_contributions)}
                </p>
              </div>
              <div className="text-center">
                <p className="text-sm text-gray-600 mb-1">投資収益</p>
                <p className="text-2xl font-bold text-success-600">
                  {new Intl.NumberFormat('ja-JP', {
                    style: 'currency',
                    currency: 'JPY',
                    maximumFractionDigits: 0,
                  }).format(result.total_gains)}
                </p>
              </div>
            </div>

            {/* 複利効果の説明 */}
            <div className="mt-6 bg-green-50 border border-green-200 rounded-lg p-4">
              <p className="text-sm text-gray-700">
                <strong>複利効果：</strong>
                投資収益により、積立元本の
                <strong className="text-success-600">
                  {' '}
                  {((result.total_gains / result.total_contributions) * 100).toFixed(1)}%
                </strong>
                の資産増加が見込まれます
              </p>
            </div>
          </div>

          {/* 年次詳細 */}
          <div className="card">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">
              年次詳細（主要マイルストーン）
            </h3>
            <div className="overflow-x-auto">
              <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                      経過年数
                    </th>
                    <th className="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">
                      総資産
                    </th>
                    <th className="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">
                      実質価値
                    </th>
                    <th className="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">
                      積立元本
                    </th>
                    <th className="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">
                      投資収益
                    </th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {result.projections
                    .filter((p) => p.year % 5 === 0 || p.year === result.projections.length)
                    .map((projection) => (
                      <tr key={projection.year}>
                        <td className="px-4 py-3 text-sm text-gray-900">
                          {projection.year}年後
                        </td>
                        <td className="px-4 py-3 text-sm text-right text-gray-900 font-medium">
                          {new Intl.NumberFormat('ja-JP', {
                            style: 'currency',
                            currency: 'JPY',
                            maximumFractionDigits: 0,
                          }).format(projection.total_assets)}
                        </td>
                        <td className="px-4 py-3 text-sm text-right text-success-600">
                          {new Intl.NumberFormat('ja-JP', {
                            style: 'currency',
                            currency: 'JPY',
                            maximumFractionDigits: 0,
                          }).format(projection.real_value)}
                        </td>
                        <td className="px-4 py-3 text-sm text-right text-orange-600">
                          {new Intl.NumberFormat('ja-JP', {
                            style: 'currency',
                            currency: 'JPY',
                            maximumFractionDigits: 0,
                          }).format(projection.contributed_amount)}
                        </td>
                        <td className="px-4 py-3 text-sm text-right text-blue-600">
                          {new Intl.NumberFormat('ja-JP', {
                            style: 'currency',
                            currency: 'JPY',
                            maximumFractionDigits: 0,
                          }).format(projection.investment_gains)}
                        </td>
                      </tr>
                    ))}
                </tbody>
              </table>
            </div>
          </div>
        </>
      )}
    </div>
  );
}
