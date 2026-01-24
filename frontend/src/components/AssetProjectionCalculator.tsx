'use client';

import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { calculationsAPI } from '@/lib/api-client';
import type { AssetProjectionRequest, AssetProjectionResponse } from '@/types/api';
import AssetProjectionChart from './AssetProjectionChart';
import { InputField, Button, LoadingSpinner } from './index';
import CurrencyInputWithPresets from './CurrencyInputWithPresets';

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
    setValue,
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
  const currentSavings = watch('current_savings');
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
        <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-4">
          資産推移シミュレーション
        </h2>
        <p className="text-gray-600 dark:text-gray-300 mb-6">
          現在の貯蓄ペースで将来どれだけ資産が増えるかを計算します
        </p>

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <InputField
              label="シミュレーション期間（年）"
              type="number"
              {...register('years', { valueAsNumber: true })}
              error={errors.years?.message}
              placeholder="30"
              className="text-base py-3"
            />

            <CurrencyInputWithPresets
              label="現在の貯蓄額"
              value={currentSavings}
              onChange={(value) => setValue('current_savings', value)}
              error={errors.current_savings?.message}
              presets={[
                { label: '50万', value: 500000 },
                { label: '100万', value: 1000000 },
                { label: '300万', value: 3000000 },
                { label: '500万', value: 5000000 },
              ]}
            />

            <CurrencyInputWithPresets
              label="月収"
              value={monthlyIncome}
              onChange={(value) => setValue('monthly_income', value)}
              error={errors.monthly_income?.message}
              presets={[
                { label: '20万', value: 200000 },
                { label: '30万', value: 300000 },
                { label: '40万', value: 400000 },
                { label: '50万', value: 500000 },
              ]}
            />

            <CurrencyInputWithPresets
              label="月間支出"
              value={monthlyExpenses}
              onChange={(value) => setValue('monthly_expenses', value)}
              error={errors.monthly_expenses?.message}
              presets={[
                { label: '15万', value: 150000 },
                { label: '20万', value: 200000 },
                { label: '30万', value: 300000 },
                { label: '40万', value: 400000 },
              ]}
            />

            <InputField
              label="投資利回り（%）"
              type="number"
              step="0.1"
              {...register('investment_return', { valueAsNumber: true })}
              error={errors.investment_return?.message}
              placeholder="5.0"
              className="text-base py-3"
            />

            <InputField
              label="インフレ率（%）"
              type="number"
              step="0.1"
              {...register('inflation_rate', { valueAsNumber: true })}
              error={errors.inflation_rate?.message}
              placeholder="2.0"
              className="text-base py-3"
            />
          </div>

          {/* 月間貯蓄額表示 */}
          <div className="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-4">
            <div className="flex flex-col sm:flex-row sm:justify-between sm:items-center gap-2">
              <span className="text-gray-700 dark:text-gray-300 font-medium">月間貯蓄額</span>
              <span
                className={`text-2xl font-bold ${
                  monthlySavings >= 0 ? 'text-success-600 dark:text-success-400' : 'text-error-600 dark:text-error-400'
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
              <p className="text-sm text-error-600 dark:text-error-400 mt-2">
                ⚠️ 支出が収入を上回っています
              </p>
            )}
          </div>

          {error && (
            <div className="bg-error-50 border border-error-200 rounded-lg p-4">
              <p className="text-error-600">{error}</p>
            </div>
          )}

          <Button type="submit" disabled={isCalculating} className="w-full py-3 text-lg min-h-[48px]">
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
            <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">
              計算結果サマリー
            </h3>
            <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 sm:gap-6">
              <div className="text-center p-4 bg-blue-50 dark:bg-blue-900/20 rounded-lg">
                <p className="text-sm text-gray-600 dark:text-gray-400 mb-1">最終資産額</p>
                <p className="text-xl sm:text-2xl font-bold text-blue-600 dark:text-blue-400 break-words">
                  {new Intl.NumberFormat('ja-JP', {
                    style: 'currency',
                    currency: 'JPY',
                    maximumFractionDigits: 0,
                  }).format(result.final_amount)}
                </p>
              </div>
              <div className="text-center p-4 bg-orange-50 dark:bg-orange-900/20 rounded-lg">
                <p className="text-sm text-gray-600 dark:text-gray-400 mb-1">積立元本合計</p>
                <p className="text-xl sm:text-2xl font-bold text-orange-600 dark:text-orange-400 break-words">
                  {new Intl.NumberFormat('ja-JP', {
                    style: 'currency',
                    currency: 'JPY',
                    maximumFractionDigits: 0,
                  }).format(result.total_contributions)}
                </p>
              </div>
              <div className="text-center p-4 bg-green-50 dark:bg-green-900/20 rounded-lg">
                <p className="text-sm text-gray-600 dark:text-gray-400 mb-1">投資収益</p>
                <p className="text-xl sm:text-2xl font-bold text-success-600 dark:text-success-400 break-words">
                  {new Intl.NumberFormat('ja-JP', {
                    style: 'currency',
                    currency: 'JPY',
                    maximumFractionDigits: 0,
                  }).format(result.total_gains)}
                </p>
              </div>
            </div>

            {/* 複利効果の説明 */}
            <div className="mt-6 bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-lg p-4">
              <p className="text-sm text-gray-700 dark:text-gray-300">
                <strong>複利効果：</strong>
                投資収益により、積立元本の
                <strong className="text-success-600 dark:text-success-400">
                  {' '}
                  {((result.total_gains / result.total_contributions) * 100).toFixed(1)}%
                </strong>
                の資産増加が見込まれます
              </p>
            </div>
          </div>

          {/* 年次詳細 */}
          <div className="card">
            <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">
              年次詳細（主要マイルストーン）
            </h3>
            <div className="overflow-x-auto -mx-4 sm:mx-0">
              <div className="inline-block min-w-full align-middle">
                <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
                  <thead className="bg-gray-50 dark:bg-gray-800">
                    <tr>
                      <th className="px-3 sm:px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase whitespace-nowrap">
                        経過年数
                      </th>
                      <th className="px-3 sm:px-4 py-3 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase whitespace-nowrap">
                        総資産
                      </th>
                      <th className="px-3 sm:px-4 py-3 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase whitespace-nowrap">
                        実質価値
                      </th>
                      <th className="px-3 sm:px-4 py-3 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase whitespace-nowrap">
                        積立元本
                      </th>
                      <th className="px-3 sm:px-4 py-3 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase whitespace-nowrap">
                        投資収益
                      </th>
                    </tr>
                  </thead>
                  <tbody className="bg-white dark:bg-gray-900 divide-y divide-gray-200 dark:divide-gray-700">
                    {result.projections
                      .filter((p) => p.year % 5 === 0 || p.year === result.projections.length)
                      .map((projection) => (
                        <tr key={projection.year}>
                          <td className="px-3 sm:px-4 py-3 text-sm text-gray-900 dark:text-white whitespace-nowrap">
                            {projection.year}年後
                          </td>
                          <td className="px-3 sm:px-4 py-3 text-sm text-right text-gray-900 dark:text-white font-medium whitespace-nowrap">
                            {new Intl.NumberFormat('ja-JP', {
                              style: 'currency',
                              currency: 'JPY',
                              maximumFractionDigits: 0,
                            }).format(projection.total_assets)}
                          </td>
                          <td className="px-3 sm:px-4 py-3 text-sm text-right text-success-600 dark:text-success-400 whitespace-nowrap">
                            {new Intl.NumberFormat('ja-JP', {
                              style: 'currency',
                              currency: 'JPY',
                              maximumFractionDigits: 0,
                            }).format(projection.real_value)}
                          </td>
                          <td className="px-3 sm:px-4 py-3 text-sm text-right text-orange-600 dark:text-orange-400 whitespace-nowrap">
                            {new Intl.NumberFormat('ja-JP', {
                              style: 'currency',
                              currency: 'JPY',
                              maximumFractionDigits: 0,
                            }).format(projection.contributed_amount)}
                          </td>
                          <td className="px-3 sm:px-4 py-3 text-sm text-right text-blue-600 dark:text-blue-400 whitespace-nowrap">
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
          </div>
        </>
      )}
    </div>
  );
}
