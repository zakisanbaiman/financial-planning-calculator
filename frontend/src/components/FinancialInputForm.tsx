'use client';

import React, { useState, useEffect, useCallback } from 'react';
import InputField from './InputField';
import Button from './Button';
import type { ExpenseItem, SavingsItem, FinancialProfile } from '@/types/api';

export interface FinancialInputFormProps {
  initialData?: FinancialProfile;
  onSubmit: (data: FinancialProfile) => Promise<void>;
  loading?: boolean;
}

interface FormErrors {
  monthly_income?: string;
  expenses?: { [key: number]: string };
  savings?: { [key: number]: string };
  investment_return?: string;
  inflation_rate?: string;
}

const FinancialInputForm: React.FC<FinancialInputFormProps> = ({
  initialData,
  onSubmit,
  loading = false,
}) => {
  // フォーム状態
  const [monthlyIncome, setMonthlyIncome] = useState(initialData?.monthly_income || 0);
  const [expenses, setExpenses] = useState<ExpenseItem[]>(
    initialData?.monthly_expenses || [
      { category: '住居費', amount: 0 },
      { category: '食費', amount: 0 },
      { category: '交通費', amount: 0 },
      { category: 'その他', amount: 0 },
    ]
  );
  const [savings, setSavings] = useState<SavingsItem[]>(
    initialData?.current_savings || [
      { type: 'deposit', amount: 0 },
      { type: 'investment', amount: 0 },
    ]
  );
  const [investmentReturn, setInvestmentReturn] = useState(
    initialData?.investment_return || 5.0
  );
  const [inflationRate, setInflationRate] = useState(initialData?.inflation_rate || 2.0);

  // エラー状態
  const [errors, setErrors] = useState<FormErrors>({});
  const [touched, setTouched] = useState<{ [key: string]: boolean }>({});

  const validateForm = useCallback((): boolean => {
    const newErrors: FormErrors = {};

    // 月収バリデーション
    if (touched.monthly_income) {
      if (monthlyIncome <= 0) {
        newErrors.monthly_income = '月収は0より大きい値を入力してください';
      } else if (monthlyIncome > 100000000) {
        newErrors.monthly_income = '月収が大きすぎます';
      }
    }

    // 支出バリデーション
    newErrors.expenses = {};
    expenses.forEach((expense, index) => {
      if (touched[`expense_${index}`]) {
        if (expense.amount < 0) {
          newErrors.expenses![index] = '支出額は0以上の値を入力してください';
        } else if (expense.amount > monthlyIncome) {
          newErrors.expenses![index] = '支出額が月収を超えています';
        }
      }
    });

    // 貯蓄バリデーション
    newErrors.savings = {};
    savings.forEach((saving, index) => {
      if (touched[`saving_${index}`]) {
        if (saving.amount < 0) {
          newErrors.savings![index] = '貯蓄額は0以上の値を入力してください';
        }
      }
    });

    // 投資利回りバリデーション
    if (touched.investment_return) {
      if (investmentReturn < 0 || investmentReturn > 100) {
        newErrors.investment_return = '投資利回りは0〜100%の範囲で入力してください';
      }
    }

    // インフレ率バリデーション
    if (touched.inflation_rate) {
      if (inflationRate < 0 || inflationRate > 50) {
        newErrors.inflation_rate = 'インフレ率は0〜50%の範囲で入力してください';
      }
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0 || 
           (Object.keys(newErrors.expenses || {}).length === 0 && 
            Object.keys(newErrors.savings || {}).length === 0 &&
            !newErrors.monthly_income &&
            !newErrors.investment_return &&
            !newErrors.inflation_rate);
  }, [monthlyIncome, expenses, savings, investmentReturn, inflationRate, touched]);

  // リアルタイムバリデーション
  useEffect(() => {
    validateForm();
  }, [validateForm]);

  const handleBlur = (field: string) => {
    setTouched((prev) => ({ ...prev, [field]: true }));
  };

  const handleExpenseChange = (index: number, field: keyof ExpenseItem, value: any) => {
    const newExpenses = [...expenses];
    newExpenses[index] = { ...newExpenses[index], [field]: value };
    setExpenses(newExpenses);
  };

  const handleSavingChange = (index: number, field: keyof SavingsItem, value: any) => {
    const newSavings = [...savings];
    newSavings[index] = { ...newSavings[index], [field]: value };
    setSavings(newSavings);
  };

  const addExpense = () => {
    setExpenses([...expenses, { category: '', amount: 0 }]);
  };

  const removeExpense = (index: number) => {
    if (expenses.length > 1) {
      setExpenses(expenses.filter((_, i) => i !== index));
    }
  };

  const addSaving = () => {
    setSavings([...savings, { type: 'other', amount: 0 }]);
  };

  const removeSaving = (index: number) => {
    if (savings.length > 1) {
      setSavings(savings.filter((_, i) => i !== index));
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    // すべてのフィールドをtouchedにする
    const allTouched: { [key: string]: boolean } = {
      monthly_income: true,
      investment_return: true,
      inflation_rate: true,
    };
    expenses.forEach((_, index) => {
      allTouched[`expense_${index}`] = true;
    });
    savings.forEach((_, index) => {
      allTouched[`saving_${index}`] = true;
    });
    setTouched(allTouched);

    if (!validateForm()) {
      return;
    }

    const formData: FinancialProfile = {
      monthly_income: monthlyIncome,
      monthly_expenses: expenses.filter(e => e.category && e.amount >= 0),
      current_savings: savings.filter(s => s.amount >= 0),
      investment_return: investmentReturn,
      inflation_rate: inflationRate,
    };

    await onSubmit(formData);
  };

  // 計算された値
  const totalExpenses = expenses.reduce((sum, e) => sum + (e.amount || 0), 0);
  const totalSavings = savings.reduce((sum, s) => sum + (s.amount || 0), 0);
  const netSavings = monthlyIncome - totalExpenses;

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      {/* 月収入力 */}
      <div className="card">
        <h3 className="text-lg font-semibold text-gray-900 mb-4">月収</h3>
        <InputField
          type="number"
          label="月収（税込）"
          value={monthlyIncome || ''}
          onChange={(e) => setMonthlyIncome(Number(e.target.value))}
          onBlur={() => handleBlur('monthly_income')}
          error={errors.monthly_income}
          placeholder="400000"
          required
          min="0"
          step="1000"
        />
      </div>

      {/* 月間支出入力 */}
      <div className="card">
        <div className="flex justify-between items-center mb-4">
          <h3 className="text-lg font-semibold text-gray-900">月間支出</h3>
          <Button type="button" variant="outline" size="sm" onClick={addExpense}>
            + 項目追加
          </Button>
        </div>
        <div className="space-y-4">
          {expenses.map((expense, index) => (
            <div key={index} className="flex gap-3 items-start">
              <div className="flex-1">
                <InputField
                  type="text"
                  label="カテゴリー"
                  value={expense.category}
                  onChange={(e) => handleExpenseChange(index, 'category', e.target.value)}
                  placeholder="例: 住居費"
                  required
                />
              </div>
              <div className="flex-1">
                <InputField
                  type="number"
                  label="金額"
                  value={expense.amount || ''}
                  onChange={(e) => handleExpenseChange(index, 'amount', Number(e.target.value))}
                  onBlur={() => handleBlur(`expense_${index}`)}
                  error={errors.expenses?.[index]}
                  placeholder="120000"
                  required
                  min="0"
                  step="1000"
                />
              </div>
              {expenses.length > 1 && (
                <button
                  type="button"
                  onClick={() => removeExpense(index)}
                  className="mt-7 text-error-500 hover:text-error-600 p-2"
                  aria-label="削除"
                >
                  ✕
                </button>
              )}
            </div>
          ))}
        </div>
        <div className="mt-4 pt-4 border-t border-gray-200">
          <div className="flex justify-between items-center">
            <span className="font-medium text-gray-700">合計支出</span>
            <span className="text-lg font-semibold text-gray-900">
              ¥{totalExpenses.toLocaleString()}
            </span>
          </div>
          <div className="flex justify-between items-center mt-2">
            <span className="font-medium text-gray-700">月間純貯蓄</span>
            <span
              className={`text-lg font-semibold ${
                netSavings >= 0 ? 'text-success-600' : 'text-error-600'
              }`}
            >
              ¥{netSavings.toLocaleString()}
            </span>
          </div>
          {netSavings < 0 && (
            <p className="mt-2 text-sm text-error-600">
              ⚠️ 支出が収入を上回っています
            </p>
          )}
        </div>
      </div>

      {/* 現在の貯蓄入力 */}
      <div className="card">
        <div className="flex justify-between items-center mb-4">
          <h3 className="text-lg font-semibold text-gray-900">現在の貯蓄</h3>
          <Button type="button" variant="outline" size="sm" onClick={addSaving}>
            + 項目追加
          </Button>
        </div>
        <div className="space-y-4">
          {savings.map((saving, index) => (
            <div key={index} className="flex gap-3 items-start">
              <div className="flex-1">
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  種類 <span className="text-error-500">*</span>
                </label>
                <select
                  value={saving.type}
                  onChange={(e) =>
                    handleSavingChange(index, 'type', e.target.value as SavingsItem['type'])
                  }
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500"
                  required
                >
                  <option value="deposit">預金</option>
                  <option value="investment">投資</option>
                  <option value="other">その他</option>
                </select>
              </div>
              <div className="flex-1">
                <InputField
                  type="number"
                  label="金額"
                  value={saving.amount || ''}
                  onChange={(e) => handleSavingChange(index, 'amount', Number(e.target.value))}
                  onBlur={() => handleBlur(`saving_${index}`)}
                  error={errors.savings?.[index]}
                  placeholder="1000000"
                  required
                  min="0"
                  step="10000"
                />
              </div>
              {savings.length > 1 && (
                <button
                  type="button"
                  onClick={() => removeSaving(index)}
                  className="mt-7 text-error-500 hover:text-error-600 p-2"
                  aria-label="削除"
                >
                  ✕
                </button>
              )}
            </div>
          ))}
        </div>
        <div className="mt-4 pt-4 border-t border-gray-200">
          <div className="flex justify-between items-center">
            <span className="font-medium text-gray-700">総資産</span>
            <span className="text-lg font-semibold text-gray-900">
              ¥{totalSavings.toLocaleString()}
            </span>
          </div>
        </div>
      </div>

      {/* 送信ボタン */}
      <div className="flex justify-end gap-3">
        <Button type="submit" loading={loading} disabled={loading}>
          保存
        </Button>
      </div>
    </form>
  );
};

export default FinancialInputForm;
