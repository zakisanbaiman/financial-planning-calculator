'use client';

import React, { useState, useEffect } from 'react';
import InputField from './InputField';
import Button from './Button';
import type { Goal, GoalType } from '@/types/api';

export interface GoalFormProps {
  initialData?: Goal;
  userId: string;
  onSubmit: (goal: Goal) => Promise<void>;
  onCancel?: () => void;
  loading?: boolean;
}

interface FormErrors {
  title?: string;
  target_amount?: string;
  target_date?: string;
  monthly_contribution?: string;
}

const goalTypeLabels: Record<GoalType, string> = {
  savings: '貯蓄',
  retirement: '老後資金',
  emergency: '緊急資金',
  custom: 'カスタム',
};

const getDefaultTargetDate = () => {
  const today = new Date();
  today.setFullYear(today.getFullYear() + 1); // 今から1年後
  return today.toISOString().split('T')[0];
};

const GoalForm: React.FC<GoalFormProps> = ({
  initialData,
  userId,
  onSubmit,
  onCancel,
  loading = false,
}) => {
  const [type, setType] = useState<GoalType>(initialData?.goal_type || 'savings');
  const [title, setTitle] = useState(initialData?.title || '');
  const [targetAmount, setTargetAmount] = useState(initialData?.target_amount || 5000000); // デフォルト値を500万円に設定
  const [targetDate, setTargetDate] = useState(
    initialData?.target_date ? initialData.target_date.split('T')[0] : getDefaultTargetDate() // 1年後の日付をデフォルトに設定
  );
  const [currentAmount, setCurrentAmount] = useState(initialData?.current_amount || 1000000); // デフォルト値を100万円に設定
  const [monthlyContribution, setMonthlyContribution] = useState(
    initialData?.monthly_contribution || 50000 // デフォルト値を5万円に設定
  );
  const [isActive, setIsActive] = useState(initialData?.is_active ?? true);

  const [errors, setErrors] = useState<FormErrors>({});
  const [touched, setTouched] = useState<{ [key: string]: boolean }>({});

  // 目標タイプに応じたデフォルトタイトル設定
  useEffect(() => {
    if (!initialData && !title) {
      const defaultTitles: Record<GoalType, string> = {
        savings: '貯蓄目標',
        retirement: '老後資金準備',
        emergency: '緊急資金確保',
        custom: '',
      };
      setTitle(defaultTitles[type]);
    }
  }, [type, initialData, title]);

  const validateForm = (): boolean => {
    const newErrors: FormErrors = {};

    if (touched.title && !title.trim()) {
      newErrors.title = 'タイトルを入力してください';
    }

    if (touched.target_amount) {
      if (targetAmount <= 0) {
        newErrors.target_amount = '目標金額は0より大きい値を入力してください';
      } else if (targetAmount > 1000000000) {
        newErrors.target_amount = '目標金額が大きすぎます';
      }
    }

    if (touched.target_date) {
      if (!targetDate) {
        newErrors.target_date = '目標期日を入力してください';
      } else {
        const selectedDate = new Date(targetDate);
        const today = new Date();
        today.setHours(0, 0, 0, 0);
        if (selectedDate < today) {
          newErrors.target_date = '目標期日は今日以降の日付を選択してください';
        }
      }
    }

    if (touched.monthly_contribution && monthlyContribution < 0) {
      newErrors.monthly_contribution = '月間積立額は0以上の値を入力してください';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  useEffect(() => {
    validateForm();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [title, targetAmount, targetDate, monthlyContribution, touched]);

  const handleBlur = (field: string) => {
    setTouched((prev) => ({ ...prev, [field]: true }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    const allTouched = {
      title: true,
      target_amount: true,
      target_date: true,
      monthly_contribution: true,
    };
    setTouched(allTouched);

    if (!validateForm()) {
      return;
    }

    const goalData: Goal = {
      ...(initialData?.id && { id: initialData.id }),
      user_id: userId,
      goal_type: type, // プロパティ名を 'type' から 'goal_type' に変更
      title: title.trim(),
      target_amount: targetAmount,
      target_date: new Date(targetDate).toISOString(),
      current_amount: currentAmount,
      monthly_contribution: monthlyContribution,
      is_active: isActive,
    };

    await onSubmit(goalData);
  };

  // 計算された値
  const remainingAmount = Math.max(0, targetAmount - currentAmount);
  const progressRate = targetAmount > 0 ? (currentAmount / targetAmount) * 100 : 0;

  // 目標達成までの月数計算
  const calculateMonthsToGoal = (): number | null => {
    if (monthlyContribution <= 0 || remainingAmount <= 0) return null;
    return Math.ceil(remainingAmount / monthlyContribution);
  };

  const monthsToGoal = calculateMonthsToGoal();

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      {/* 目標タイプ */}
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">
          目標タイプ <span className="text-error-500">*</span>
        </label>
        <div className="grid grid-cols-2 gap-3">
          {(Object.keys(goalTypeLabels) as GoalType[]).map((goalType) => (
            <button
              key={goalType}
              type="button"
              onClick={() => setType(goalType)}
              className={`px-4 py-3 rounded-lg border-2 transition-colors ${
                type === goalType
                  ? 'border-primary-500 bg-primary-50 text-primary-700'
                  : 'border-gray-300 hover:border-gray-400 text-gray-700'
              }`}
            >
              {goalTypeLabels[goalType]}
            </button>
          ))}
        </div>
      </div>

      {/* タイトル */}
      <InputField
        type="text"
        label="目標タイトル"
        value={title}
        onChange={(e) => setTitle(e.target.value)}
        onBlur={() => handleBlur('title')}
        error={errors.title}
        placeholder="例: マイホーム購入資金"
        required
      />

      {/* 目標金額 */}
      <InputField
        type="number"
        label="目標金額"
        value={targetAmount || ''}
        onChange={(e) => setTargetAmount(Number(e.target.value))}
        onBlur={() => handleBlur('target_amount')}
        error={errors.target_amount}
        placeholder="5000000"
        required
        min="0"
        step="10000"
      />

      {/* 目標期日 */}
      <InputField
        type="date"
        label="目標期日"
        value={targetDate}
        onChange={(e) => setTargetDate(e.target.value)}
        onBlur={() => handleBlur('target_date')}
        error={errors.target_date}
        required
      />

      {/* 現在の積立額 */}
      <InputField
        type="number"
        label="現在の積立額"
        value={currentAmount || ''}
        onChange={(e) => setCurrentAmount(Number(e.target.value))}
        placeholder="1000000"
        helperText="この目標のために既に積み立てている金額"
        min="0"
        step="10000"
      />

      {/* 月間積立額 */}
      <InputField
        type="number"
        label="月間積立額"
        value={monthlyContribution || ''}
        onChange={(e) => setMonthlyContribution(Number(e.target.value))}
        onBlur={() => handleBlur('monthly_contribution')}
        error={errors.monthly_contribution}
        placeholder="50000"
        helperText="毎月この目標のために積み立てる金額"
        min="0"
        step="1000"
      />

      {/* 進捗状況表示 */}
      {targetAmount > 0 && (
        <div className="p-4 bg-gray-50 rounded-lg space-y-3">
          <div className="flex justify-between items-center text-sm">
            <span className="text-gray-600">進捗状況</span>
            <span className="font-semibold text-gray-900">
              {progressRate.toFixed(1)}%
            </span>
          </div>
          <div className="w-full bg-gray-200 rounded-full h-2">
            <div
              className="bg-primary-500 h-2 rounded-full transition-all"
              style={{ width: `${Math.min(progressRate, 100)}%` }}
            />
          </div>
          <div className="flex justify-between items-center text-sm">
            <span className="text-gray-600">残り</span>
            <span className="font-semibold text-gray-900">
              ¥{remainingAmount.toLocaleString()}
            </span>
          </div>
          {monthsToGoal !== null && (
            <div className="flex justify-between items-center text-sm">
              <span className="text-gray-600">達成まで（現在のペース）</span>
              <span className="font-semibold text-gray-900">
                約{monthsToGoal}ヶ月
              </span>
            </div>
          )}
        </div>
      )}

      {/* アクティブ状態 */}
      <div className="flex items-center">
        <input
          type="checkbox"
          id="is_active"
          checked={isActive}
          onChange={(e) => setIsActive(e.target.checked)}
          className="w-4 h-4 text-primary-500 border-gray-300 rounded focus:ring-primary-500"
        />
        <label htmlFor="is_active" className="ml-2 text-sm text-gray-700">
          この目標をアクティブにする
        </label>
      </div>

      {/* ボタン */}
      <div className="flex justify-end gap-3 pt-4">
        {onCancel && (
          <Button type="button" variant="outline" onClick={onCancel}>
            キャンセル
          </Button>
        )}
        <Button type="submit" loading={loading} disabled={loading}>
          {initialData ? '更新' : '作成'}
        </Button>
      </div>
    </form>
  );
};

export default GoalForm;
