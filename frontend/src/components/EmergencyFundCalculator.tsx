'use client';

import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { calculationsAPI } from '@/lib/api-client';
import type { EmergencyFundRequest, EmergencyFundResponse } from '@/types/api';
import { InputField, Button, LoadingSpinner } from './index';
import CurrencyInputWithPresets from './CurrencyInputWithPresets';

// バリデーションスキーマ
const emergencyFundSchema = z.object({
  monthly_expenses: z.number().min(0, '0以上の値を入力してください'),
  target_months: z
    .number()
    .min(1, '1ヶ月以上を指定してください')
    .max(24, '24ヶ月以内で指定してください'),
  current_savings: z.number().min(0, '0以上の値を入力してください'),
});

type EmergencyFundFormData = z.infer<typeof emergencyFundSchema>;

interface EmergencyFundCalculatorProps {
  userId: string;
  initialData?: Partial<EmergencyFundFormData>;
}

export default function EmergencyFundCalculator({
  userId,
  initialData,
}: EmergencyFundCalculatorProps) {
  const [result, setResult] = useState<EmergencyFundResponse | null>(null);
  const [isCalculating, setIsCalculating] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const {
    register,
    handleSubmit,
    formState: { errors },
    watch,
    setValue,
  } = useForm<EmergencyFundFormData>({
    resolver: zodResolver(emergencyFundSchema),
    defaultValues: {
      monthly_expenses: initialData?.monthly_expenses || 280000,
      target_months: initialData?.target_months || 6,
      current_savings: initialData?.current_savings || 600000,
    },
  });

  const monthlyExpenses = watch('monthly_expenses');
  const targetMonths = watch('target_months');
  const currentSavings = watch('current_savings');
  const targetAmount = monthlyExpenses * targetMonths;

  const onSubmit = async (data: EmergencyFundFormData) => {
    setIsCalculating(true);
    setError(null);

    try {
      const request: EmergencyFundRequest = {
        user_id: userId,
        ...data,
      };

      const response = await calculationsAPI.emergencyFund(request);
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
        <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-4">緊急資金計算</h2>
        <p className="text-gray-600 dark:text-gray-300 mb-6">
          万が一の時（失業、病気など）に必要な緊急資金を計算します
        </p>

        {/* 緊急資金の説明 */}
        <div className="mb-6 bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-4">
          <h3 className="text-sm font-semibold text-gray-900 dark:text-white mb-2">
            💡 緊急資金とは？
          </h3>
          <p className="text-sm text-gray-700 dark:text-gray-300 mb-2">
            予期せぬ出来事（失業、病気、事故など）に備えて、すぐに使える形で
            確保しておくべき資金です。
          </p>
          <ul className="text-sm text-gray-700 dark:text-gray-300 space-y-1 ml-4">
            <li>• 一般的には生活費の3〜6ヶ月分が推奨されます</li>
            <li>• 自営業や収入が不安定な場合は6〜12ヶ月分が理想的です</li>
            <li>• 預金など、すぐに引き出せる形で保管することが重要です</li>
          </ul>
        </div>

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <CurrencyInputWithPresets
              label="月間生活費"
              value={monthlyExpenses}
              onChange={(value) => setValue('monthly_expenses', value)}
              error={errors.monthly_expenses?.message}
              helperText="家賃、食費、光熱費など必要最低限の支出"
              presets={[
                { label: '15万', value: 150000 },
                { label: '20万', value: 200000 },
                { label: '30万', value: 300000 },
                { label: '40万', value: 400000 },
              ]}
            />

            <InputField
              label="確保したい期間（ヶ月）"
              type="number"
              {...register('target_months', { valueAsNumber: true })}
              error={errors.target_months?.message}
              placeholder="6"
              helperText="3〜6ヶ月が一般的"
              className="text-base py-3"
            />

            <CurrencyInputWithPresets
              label="現在の緊急資金"
              value={currentSavings}
              onChange={(value) => setValue('current_savings', value)}
              error={errors.current_savings?.message}
              helperText="すぐに引き出せる預金額"
              presets={[
                { label: '50万', value: 500000 },
                { label: '100万', value: 1000000 },
                { label: '200万', value: 2000000 },
                { label: '300万', value: 3000000 },
              ]}
            />
          </div>

          {/* 目標額表示 */}
          <div className="bg-purple-50 dark:bg-purple-900/20 border border-purple-200 dark:border-purple-800 rounded-lg p-4">
            <div className="flex flex-col sm:flex-row sm:justify-between sm:items-center gap-2">
              <span className="text-gray-700 dark:text-gray-300 font-medium">目標緊急資金額</span>
              <span className="text-2xl font-bold text-purple-600 dark:text-purple-400">
                {new Intl.NumberFormat('ja-JP', {
                  style: 'currency',
                  currency: 'JPY',
                  maximumFractionDigits: 0,
                }).format(targetAmount)}
              </span>
            </div>
            <p className="text-sm text-gray-600 dark:text-gray-400 mt-2">
              月間生活費 {new Intl.NumberFormat('ja-JP').format(monthlyExpenses)}円 ×{' '}
              {targetMonths}ヶ月
            </p>
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
          {/* 充足状況 */}
          <div className="card">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">
              緊急資金充足状況
            </h3>

            {/* 充足率表示 */}
            <div className="mb-6">
              <div className="flex justify-between items-center mb-2">
                <span className="text-sm text-gray-600">充足率</span>
                <span
                  className={`text-2xl font-bold ${
                    result.sufficiency_rate >= 100
                      ? 'text-success-600'
                      : result.sufficiency_rate >= 50
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
                      : result.sufficiency_rate >= 50
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
                <span className="text-gray-600">推奨緊急資金額</span>
                <span className="text-lg font-semibold text-gray-900">
                  {new Intl.NumberFormat('ja-JP', {
                    style: 'currency',
                    currency: 'JPY',
                    maximumFractionDigits: 0,
                  }).format(result.required_amount)}
                </span>
              </div>

              <div className="flex justify-between items-center py-3 border-b border-gray-200">
                <span className="text-gray-600">現在の緊急資金</span>
                <span className="text-lg font-semibold text-gray-900">
                  {new Intl.NumberFormat('ja-JP', {
                    style: 'currency',
                    currency: 'JPY',
                    maximumFractionDigits: 0,
                  }).format(result.current_amount)}
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

              {result.shortfall > 0 && result.months_to_target > 0 && (
                <div className="flex justify-between items-center py-3">
                  <span className="text-gray-600">目標達成までの期間</span>
                  <span className="text-lg font-semibold text-gray-900">
                    {result.months_to_target}ヶ月
                  </span>
                </div>
              )}
            </div>
          </div>

          {/* 可視化 */}
          <div className="card">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">
              緊急資金の内訳
            </h3>
            <div className="space-y-4">
              {/* 現在の資金 */}
              <div>
                <div className="flex justify-between items-center mb-2">
                  <span className="text-sm text-gray-600">現在の緊急資金</span>
                  <span className="text-sm font-medium text-gray-900">
                    {new Intl.NumberFormat('ja-JP', {
                      style: 'currency',
                      currency: 'JPY',
                      maximumFractionDigits: 0,
                    }).format(result.current_amount)}
                  </span>
                </div>
                <div className="w-full bg-gray-200 rounded-full h-8">
                  <div
                    className="bg-success-600 h-8 rounded-full flex items-center justify-end pr-3"
                    style={{
                      width: `${Math.min(
                        100,
                        (result.current_amount / result.required_amount) * 100
                      )}%`,
                    }}
                  >
                    <span className="text-xs text-white font-medium">
                      {result.sufficiency_rate.toFixed(0)}%
                    </span>
                  </div>
                </div>
              </div>

              {/* 不足分 */}
              {result.shortfall > 0 && (
                <div>
                  <div className="flex justify-between items-center mb-2">
                    <span className="text-sm text-gray-600">不足分</span>
                    <span className="text-sm font-medium text-error-600">
                      {new Intl.NumberFormat('ja-JP', {
                        style: 'currency',
                        currency: 'JPY',
                        maximumFractionDigits: 0,
                      }).format(result.shortfall)}
                    </span>
                  </div>
                  <div className="w-full bg-gray-200 rounded-full h-8">
                    <div
                      className="bg-error-600 h-8 rounded-full flex items-center justify-end pr-3"
                      style={{
                        width: `${(result.shortfall / result.required_amount) * 100}%`,
                      }}
                    >
                      <span className="text-xs text-white font-medium">
                        {((result.shortfall / result.required_amount) * 100).toFixed(0)}%
                      </span>
                    </div>
                  </div>
                </div>
              )}
            </div>

            {/* 月数表示 */}
            <div className="mt-6 grid grid-cols-2 gap-4">
              <div className="bg-success-50 border border-success-200 rounded-lg p-4 text-center">
                <p className="text-sm text-gray-600 mb-1">現在カバーできる期間</p>
                <p className="text-2xl font-bold text-success-600">
                  {(result.current_amount / monthlyExpenses).toFixed(1)}ヶ月
                </p>
              </div>
              <div className="bg-purple-50 border border-purple-200 rounded-lg p-4 text-center">
                <p className="text-sm text-gray-600 mb-1">目標期間</p>
                <p className="text-2xl font-bold text-purple-600">{targetMonths}ヶ月</p>
              </div>
            </div>
          </div>

          {/* アドバイス */}
          {result.shortfall > 0 ? (
            <div className="card bg-warning-50 border-warning-200">
              <h3 className="text-lg font-semibold text-gray-900 mb-4">
                ⚠️ 緊急資金が不足しています
              </h3>
              <p className="text-gray-700 mb-4">
                予期せぬ出来事に備えて、緊急資金を増やすことをお勧めします。
              </p>

              <div className="space-y-3">
                <div className="bg-white rounded-lg p-4">
                  <p className="text-sm text-gray-600 mb-1">
                    月々の貯蓄で目標達成する場合
                  </p>
                  {result.months_to_target > 0 && (
                    <>
                      <p className="text-lg font-bold text-warning-600">
                        月額{' '}
                        {new Intl.NumberFormat('ja-JP', {
                          style: 'currency',
                          currency: 'JPY',
                          maximumFractionDigits: 0,
                        }).format(result.shortfall / result.months_to_target)}
                      </p>
                      <p className="text-sm text-gray-600 mt-1">
                        {result.months_to_target}ヶ月で目標達成
                      </p>
                    </>
                  )}
                </div>

                <div className="bg-white rounded-lg p-4">
                  <p className="text-sm font-semibold text-gray-900 mb-2">
                    緊急資金を増やす方法
                  </p>
                  <ul className="space-y-2 text-sm text-gray-700">
                    <li className="flex items-start gap-2">
                      <span className="text-success-600">✓</span>
                      <span>毎月の収入から一定額を自動的に緊急資金口座に移す</span>
                    </li>
                    <li className="flex items-start gap-2">
                      <span className="text-success-600">✓</span>
                      <span>ボーナスや臨時収入の一部を緊急資金に充てる</span>
                    </li>
                    <li className="flex items-start gap-2">
                      <span className="text-success-600">✓</span>
                      <span>不要な支出を見直して貯蓄に回す</span>
                    </li>
                    <li className="flex items-start gap-2">
                      <span className="text-success-600">✓</span>
                      <span>
                        普通預金や定期預金など、すぐに引き出せる形で保管する
                      </span>
                    </li>
                  </ul>
                </div>

                <div className="bg-white rounded-lg p-4">
                  <p className="text-sm font-semibold text-gray-900 mb-2">
                    緊急資金の重要性
                  </p>
                  <ul className="space-y-2 text-sm text-gray-700">
                    <li className="flex items-start gap-2">
                      <span className="text-blue-600">•</span>
                      <span>失業時の生活費をカバー</span>
                    </li>
                    <li className="flex items-start gap-2">
                      <span className="text-blue-600">•</span>
                      <span>急な医療費や修理費に対応</span>
                    </li>
                    <li className="flex items-start gap-2">
                      <span className="text-blue-600">•</span>
                      <span>投資資産を緊急時に売却せずに済む</span>
                    </li>
                    <li className="flex items-start gap-2">
                      <span className="text-blue-600">•</span>
                      <span>精神的な安心感を得られる</span>
                    </li>
                  </ul>
                </div>
              </div>
            </div>
          ) : (
            <div className="card bg-success-50 border-success-200">
              <h3 className="text-lg font-semibold text-gray-900 mb-4">
                ✅ 緊急資金は十分に確保されています
              </h3>
              <p className="text-gray-700 mb-4">
                予期せぬ出来事にも対応できる十分な緊急資金が確保されています。
              </p>
              <div className="bg-white rounded-lg p-4">
                <p className="text-sm font-semibold text-gray-900 mb-2">
                  今後の注意点
                </p>
                <ul className="space-y-2 text-sm text-gray-700">
                  <li className="flex items-start gap-2">
                    <span className="text-success-600">✓</span>
                    <span>生活費が変わったら緊急資金の目標額も見直しましょう</span>
                  </li>
                  <li className="flex items-start gap-2">
                    <span className="text-success-600">✓</span>
                    <span>緊急資金は投資に回さず、流動性の高い形で保管</span>
                  </li>
                  <li className="flex items-start gap-2">
                    <span className="text-success-600">✓</span>
                    <span>余裕資金は他の目標（老後資金など）に活用できます</span>
                  </li>
                  <li className="flex items-start gap-2">
                    <span className="text-success-600">✓</span>
                    <span>定期的に緊急資金の状況を確認しましょう</span>
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
                <p className="text-gray-600 mb-1">月間生活費</p>
                <p className="font-medium text-gray-900">
                  {new Intl.NumberFormat('ja-JP', {
                    style: 'currency',
                    currency: 'JPY',
                    maximumFractionDigits: 0,
                  }).format(monthlyExpenses)}
                </p>
              </div>
              <div>
                <p className="text-gray-600 mb-1">確保期間</p>
                <p className="font-medium text-gray-900">{targetMonths}ヶ月分</p>
              </div>
              <div>
                <p className="text-gray-600 mb-1">推奨額の根拠</p>
                <p className="font-medium text-gray-900">
                  生活費 × 期間
                </p>
              </div>
              <div>
                <p className="text-gray-600 mb-1">資金の性質</p>
                <p className="font-medium text-gray-900">流動性重視</p>
              </div>
            </div>
            <div className="mt-4 p-3 bg-blue-50 rounded-lg">
              <p className="text-xs text-gray-600">
                ※ 緊急資金は投資に回さず、普通預金や定期預金など、
                すぐに引き出せる形で保管することが重要です。
              </p>
            </div>
          </div>
        </>
      )}
    </div>
  );
}
