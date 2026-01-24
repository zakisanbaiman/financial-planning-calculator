'use client';

import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { calculationsAPI } from '@/lib/api-client';
import type {
  RetirementCalculationRequest,
  RetirementCalculationResponse,
} from '@/types/api';
import { InputField, Button, LoadingSpinner } from './index';
import CurrencyInputWithPresets from './CurrencyInputWithPresets';

// バリデーションスキーマ
const retirementSchema = z.object({
  current_age: z
    .number()
    .min(18, '18歳以上を入力してください')
    .max(100, '100歳以下を入力してください'),
  retirement_age: z
    .number()
    .min(50, '50歳以上を入力してください')
    .max(100, '100歳以下を入力してください'),
  life_expectancy: z
    .number()
    .min(60, '60歳以上を入力してください')
    .max(120, '120歳以下を入力してください'),
  monthly_retirement_expenses: z
    .number()
    .min(0, '0以上の値を入力してください'),
  pension_amount: z.number().min(0, '0以上の値を入力してください'),
  current_savings: z.number().min(0, '0以上の値を入力してください'),
  monthly_savings: z.number().min(0, '0以上の値を入力してください'),
  investment_return: z
    .number()
    .min(0, '0以上の値を入力してください')
    .max(100, '100以下の値を入力してください'),
  inflation_rate: z
    .number()
    .min(0, '0以上の値を入力してください')
    .max(50, '50以下の値を入力してください'),
});

type RetirementFormData = z.infer<typeof retirementSchema>;

interface RetirementCalculatorProps {
  userId: string;
  initialData?: Partial<RetirementFormData>;
}

export default function RetirementCalculator({
  userId,
  initialData,
}: RetirementCalculatorProps) {
  const [result, setResult] = useState<RetirementCalculationResponse | null>(null);
  const [isCalculating, setIsCalculating] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const {
    register,
    handleSubmit,
    formState: { errors },
    watch,
    setValue,
  } = useForm<RetirementFormData>({
    resolver: zodResolver(retirementSchema),
    defaultValues: {
      current_age: initialData?.current_age || 35,
      retirement_age: initialData?.retirement_age || 65,
      life_expectancy: initialData?.life_expectancy || 90,
      monthly_retirement_expenses: initialData?.monthly_retirement_expenses || 250000,
      pension_amount: initialData?.pension_amount || 150000,
      current_savings: initialData?.current_savings || 1500000,
      monthly_savings: initialData?.monthly_savings || 120000,
      investment_return: initialData?.investment_return || 5.0,
      inflation_rate: initialData?.inflation_rate || 2.0,
    },
  });

  const currentAge = watch('current_age');
  const retirementAge = watch('retirement_age');
  const lifeExpectancy = watch('life_expectancy');
  const monthlyRetirementExpenses = watch('monthly_retirement_expenses');
  const pensionAmount = watch('pension_amount');
  const currentSavings = watch('current_savings');
  const monthlySavings = watch('monthly_savings');

  const yearsUntilRetirement = Math.max(0, retirementAge - currentAge);
  const yearsInRetirement = Math.max(0, lifeExpectancy - retirementAge);
  const monthlyShortfall = Math.max(0, monthlyRetirementExpenses - pensionAmount);

  const onSubmit = async (data: RetirementFormData) => {
    setIsCalculating(true);
    setError(null);

    try {
      const request: RetirementCalculationRequest = {
        user_id: userId,
        ...data,
      };

      const response = await calculationsAPI.retirement(request);
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
        <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-4">老後資金計算</h2>
        <p className="text-gray-600 dark:text-gray-300 mb-6">
          退職後に必要な資金と年金額を考慮して、老後の財務計画を立てます
        </p>

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
          {/* 年齢設定 */}
          <div>
            <h3 className="text-md font-semibold text-gray-900 dark:text-white mb-3">年齢設定</h3>
            <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
              <InputField
                label="現在の年齢"
                type="number"
                {...register('current_age', { valueAsNumber: true })}
                error={errors.current_age?.message}
                placeholder="35"
                className="text-base py-3"
              />

              <InputField
                label="退職予定年齢"
                type="number"
                {...register('retirement_age', { valueAsNumber: true })}
                error={errors.retirement_age?.message}
                placeholder="65"
                className="text-base py-3"
              />

              <InputField
                label="想定寿命"
                type="number"
                {...register('life_expectancy', { valueAsNumber: true })}
                error={errors.life_expectancy?.message}
                placeholder="90"
                className="text-base py-3"
              />
            </div>

            {/* 期間表示 */}
            <div className="mt-3 grid grid-cols-2 gap-4">
              <div className="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-3">
                <p className="text-sm text-gray-600 dark:text-gray-400">退職までの期間</p>
                <p className="text-xl font-bold text-blue-600 dark:text-blue-400">
                  {yearsUntilRetirement}年
                </p>
              </div>
              <div className="bg-purple-50 dark:bg-purple-900/20 border border-purple-200 dark:border-purple-800 rounded-lg p-3">
                <p className="text-sm text-gray-600 dark:text-gray-400">退職後の期間</p>
                <p className="text-xl font-bold text-purple-600 dark:text-purple-400">
                  {yearsInRetirement}年
                </p>
              </div>
            </div>
          </div>

          {/* 老後の生活費 */}
          <div>
            <h3 className="text-md font-semibold text-gray-900 dark:text-white mb-3">
              老後の生活費と年金
            </h3>
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              <CurrencyInputWithPresets
                label="月間生活費"
                value={monthlyRetirementExpenses}
                onChange={(value) => setValue('monthly_retirement_expenses', value)}
                error={errors.monthly_retirement_expenses?.message}
                presets={[
                  { label: '15万', value: 150000 },
                  { label: '20万', value: 200000 },
                  { label: '25万', value: 250000 },
                  { label: '30万', value: 300000 },
                ]}
              />

              <CurrencyInputWithPresets
                label="年金受給額（月額）"
                value={pensionAmount}
                onChange={(value) => setValue('pension_amount', value)}
                error={errors.pension_amount?.message}
                presets={[
                  { label: '10万', value: 100000 },
                  { label: '15万', value: 150000 },
                  { label: '20万', value: 200000 },
                  { label: '25万', value: 250000 },
                ]}
              />
            </div>

            {/* 不足額表示 */}
            <div className="mt-3 bg-orange-50 dark:bg-orange-900/20 border border-orange-200 dark:border-orange-800 rounded-lg p-3">
              <div className="flex flex-col sm:flex-row sm:justify-between sm:items-center gap-2">
                <span className="text-sm text-gray-600 dark:text-gray-400">
                  年金だけでは不足する月額
                </span>
                <span className="text-xl font-bold text-orange-600 dark:text-orange-400">
                  {new Intl.NumberFormat('ja-JP', {
                    style: 'currency',
                    currency: 'JPY',
                    maximumFractionDigits: 0,
                  }).format(monthlyShortfall)}
                </span>
              </div>
            </div>
          </div>

          {/* 現在の資産と貯蓄 */}
          <div>
            <h3 className="text-md font-semibold text-gray-900 dark:text-white mb-3">
              現在の資産と貯蓄
            </h3>
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              <CurrencyInputWithPresets
                label="現在の貯蓄額"
                value={currentSavings}
                onChange={(value) => setValue('current_savings', value)}
                error={errors.current_savings?.message}
                presets={[
                  { label: '100万', value: 1000000 },
                  { label: '300万', value: 3000000 },
                  { label: '500万', value: 5000000 },
                  { label: '1000万', value: 10000000 },
                ]}
              />

              <CurrencyInputWithPresets
                label="月間貯蓄額"
                value={monthlySavings}
                onChange={(value) => setValue('monthly_savings', value)}
                error={errors.monthly_savings?.message}
                presets={[
                  { label: '5万', value: 50000 },
                  { label: '10万', value: 100000 },
                  { label: '15万', value: 150000 },
                  { label: '20万', value: 200000 },
                ]}
              />
            </div>
          </div>

          {/* 投資設定 */}
          <div>
            <h3 className="text-md font-semibold text-gray-900 dark:text-white mb-3">投資設定</h3>
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
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
          {/* 結果サマリー */}
          <div className="card">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">
              老後資金計算結果
            </h3>

            {/* 充足率表示 */}
            <div className="mb-6">
              <div className="flex justify-between items-center mb-2">
                <span className="text-sm text-gray-600">老後資金充足率</span>
                <span
                  className={`text-2xl font-bold ${
                    result.sufficiency_rate >= 100
                      ? 'text-success-600'
                      : result.sufficiency_rate >= 70
                      ? 'text-warning-600'
                      : 'text-error-600'
                  }`}
                >
                  {result.sufficiency_rate.toFixed(1)}%
                </span>
              </div>
              <div className="w-full bg-gray-200 rounded-full h-4">
                <div
                  className={`h-4 rounded-full transition-all ${
                    result.sufficiency_rate >= 100
                      ? 'bg-success-600'
                      : result.sufficiency_rate >= 70
                      ? 'bg-warning-600'
                      : 'bg-error-600'
                  }`}
                  style={{ width: `${Math.min(100, result.sufficiency_rate)}%` }}
                />
              </div>
            </div>

            {/* 詳細数値 */}
            <div className="space-y-4">
              <div className="flex justify-between items-center py-3 border-b border-gray-200">
                <span className="text-gray-600">必要老後資金</span>
                <span className="text-lg font-semibold text-gray-900">
                  {new Intl.NumberFormat('ja-JP', {
                    style: 'currency',
                    currency: 'JPY',
                    maximumFractionDigits: 0,
                  }).format(result.required_amount)}
                </span>
              </div>

              <div className="flex justify-between items-center py-3 border-b border-gray-200">
                <span className="text-gray-600">予想達成額</span>
                <span className="text-lg font-semibold text-gray-900">
                  {new Intl.NumberFormat('ja-JP', {
                    style: 'currency',
                    currency: 'JPY',
                    maximumFractionDigits: 0,
                  }).format(result.projected_amount)}
                </span>
              </div>

              <div className="flex justify-between items-center py-3 border-b border-gray-200">
                <span className="text-gray-600">
                  {result.shortfall > 0 ? '不足額' : '余裕額'}
                </span>
                <span
                  className={`text-lg font-semibold ${
                    result.shortfall > 0 ? 'text-error-600' : 'text-success-600'
                  }`}
                >
                  {result.shortfall > 0 ? '' : '+'}
                  {new Intl.NumberFormat('ja-JP', {
                    style: 'currency',
                    currency: 'JPY',
                    maximumFractionDigits: 0,
                  }).format(Math.abs(result.shortfall))}
                </span>
              </div>

              <div className="flex justify-between items-center py-3">
                <span className="text-gray-600">退職までの期間</span>
                <span className="text-lg font-semibold text-gray-900">
                  {result.years_until_retirement}年
                </span>
              </div>
            </div>
          </div>

          {/* 推奨事項 */}
          {result.shortfall > 0 && (
            <div className="card bg-warning-50 border-warning-200">
              <h3 className="text-lg font-semibold text-gray-900 mb-4">
                ⚠️ 推奨事項
              </h3>
              <p className="text-gray-700 mb-4">
                現在の貯蓄ペースでは老後資金が不足する可能性があります。
                以下の対策をご検討ください：
              </p>
              <div className="space-y-3">
                <div className="bg-white rounded-lg p-4">
                  <p className="text-sm text-gray-600 mb-1">推奨月間貯蓄額</p>
                  <p className="text-2xl font-bold text-warning-600">
                    {new Intl.NumberFormat('ja-JP', {
                      style: 'currency',
                      currency: 'JPY',
                      maximumFractionDigits: 0,
                    }).format(result.recommended_monthly_savings)}
                  </p>
                  <p className="text-sm text-gray-600 mt-2">
                    現在より
                    <strong className="text-warning-600">
                      {' '}
                      {new Intl.NumberFormat('ja-JP', {
                        style: 'currency',
                        currency: 'JPY',
                        maximumFractionDigits: 0,
                      }).format(
                        result.recommended_monthly_savings - monthlyShortfall
                      )}
                    </strong>
                    の追加貯蓄が必要です
                  </p>
                </div>

                <div className="bg-white rounded-lg p-4">
                  <ul className="space-y-2 text-sm text-gray-700">
                    <li className="flex items-start gap-2">
                      <span className="text-success-600">✓</span>
                      <span>支出を見直して貯蓄額を増やす</span>
                    </li>
                    <li className="flex items-start gap-2">
                      <span className="text-success-600">✓</span>
                      <span>投資利回りを改善する（リスク許容度に応じて）</span>
                    </li>
                    <li className="flex items-start gap-2">
                      <span className="text-success-600">✓</span>
                      <span>退職年齢を延ばすことを検討する</span>
                    </li>
                    <li className="flex items-start gap-2">
                      <span className="text-success-600">✓</span>
                      <span>老後の生活費を見直す</span>
                    </li>
                  </ul>
                </div>
              </div>
            </div>
          )}

          {result.shortfall <= 0 && (
            <div className="card bg-success-50 border-success-200">
              <h3 className="text-lg font-semibold text-gray-900 mb-4">
                ✅ 良好な状態です
              </h3>
              <p className="text-gray-700 mb-4">
                現在の貯蓄ペースを維持すれば、老後資金は十分に確保できる見込みです。
              </p>
              <div className="bg-white rounded-lg p-4">
                <ul className="space-y-2 text-sm text-gray-700">
                  <li className="flex items-start gap-2">
                    <span className="text-success-600">✓</span>
                    <span>定期的に計画を見直しましょう</span>
                  </li>
                  <li className="flex items-start gap-2">
                    <span className="text-success-600">✓</span>
                    <span>余裕資金で他の目標達成を検討できます</span>
                  </li>
                  <li className="flex items-start gap-2">
                    <span className="text-success-600">✓</span>
                    <span>インフレ率の変動に注意しましょう</span>
                  </li>
                </ul>
              </div>
            </div>
          )}

          {/* 計算の前提 */}
          <div className="card">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">
              計算の前提条件
            </h3>
            <div className="grid md:grid-cols-2 gap-4 text-sm">
              <div>
                <p className="text-gray-600 mb-1">退職までの積立期間</p>
                <p className="font-medium text-gray-900">
                  {result.years_until_retirement}年間
                </p>
              </div>
              <div>
                <p className="text-gray-600 mb-1">退職後の生活期間</p>
                <p className="font-medium text-gray-900">{yearsInRetirement}年間</p>
              </div>
              <div>
                <p className="text-gray-600 mb-1">年金不足額（月額）</p>
                <p className="font-medium text-gray-900">
                  {new Intl.NumberFormat('ja-JP', {
                    style: 'currency',
                    currency: 'JPY',
                    maximumFractionDigits: 0,
                  }).format(monthlyShortfall)}
                </p>
              </div>
              <div>
                <p className="text-gray-600 mb-1">複利効果</p>
                <p className="font-medium text-gray-900">
                  投資利回り考慮済み
                </p>
              </div>
            </div>
            <div className="mt-4 p-3 bg-blue-50 rounded-lg">
              <p className="text-xs text-gray-600">
                ※ この計算は簡易的なシミュレーションです。実際の老後資金計画には、
                医療費、介護費用、住宅費用なども考慮する必要があります。
              </p>
            </div>
          </div>
        </>
      )}
    </div>
  );
}
