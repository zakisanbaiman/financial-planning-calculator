'use client';

import React, { useState, useEffect, useCallback } from 'react';
import InputField from './InputField';
import Button from './Button';

export interface InvestmentSettings {
  investment_return: number;
  inflation_rate: number;
}

export interface InvestmentSettingsFormProps {
  initialData?: InvestmentSettings;
  onSubmit: (data: InvestmentSettings) => Promise<void>;
  loading?: boolean;
}

interface FormErrors {
  investment_return?: string;
  inflation_rate?: string;
}

const InvestmentSettingsForm: React.FC<InvestmentSettingsFormProps> = ({
  initialData,
  onSubmit,
  loading = false,
}) => {
  const [investmentReturn, setInvestmentReturn] = useState(
    initialData?.investment_return ?? 5.0
  );
  const [inflationRate, setInflationRate] = useState(initialData?.inflation_rate ?? 2.0);
  const [errors, setErrors] = useState<FormErrors>({});
  const [touched, setTouched] = useState<{ [key: string]: boolean }>({});

  const validateForm = useCallback((): boolean => {
    const newErrors: FormErrors = {};

    // 投資利回りバリデーション
    if (touched.investment_return) {
      if (investmentReturn < 0) {
        newErrors.investment_return = '投資利回りは0以上の値を入力してください';
      } else if (investmentReturn > 100) {
        newErrors.investment_return = '投資利回りは100%以下の値を入力してください';
      }
    }

    // インフレ率バリデーション
    if (touched.inflation_rate) {
      if (inflationRate < 0) {
        newErrors.inflation_rate = 'インフレ率は0以上の値を入力してください';
      } else if (inflationRate > 50) {
        newErrors.inflation_rate = 'インフレ率は50%以下の値を入力してください';
      }
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  }, [investmentReturn, inflationRate, touched]);

  // リアルタイムバリデーション
  useEffect(() => {
    validateForm();
  }, [validateForm]);

  const handleBlur = (field: string) => {
    setTouched((prev) => ({ ...prev, [field]: true }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    // すべてのフィールドをtouchedにする
    setTouched({
      investment_return: true,
      inflation_rate: true,
    });

    if (!validateForm()) {
      return;
    }

    const formData: InvestmentSettings = {
      investment_return: investmentReturn,
      inflation_rate: inflationRate,
    };

    await onSubmit(formData);
  };

  // 実質利回り計算
  const realReturn = investmentReturn - inflationRate;

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div className="card">
        <h3 className="text-lg font-semibold text-gray-900 mb-4">投資・インフレ設定</h3>
        
        <div className="space-y-4">
          {/* 投資利回り */}
          <div>
            <InputField
              type="number"
              label="期待投資利回り（年率）"
              value={investmentReturn || ''}
              onChange={(e) => setInvestmentReturn(Number(e.target.value))}
              onBlur={() => handleBlur('investment_return')}
              error={errors.investment_return}
              helperText="投資による年間の期待リターン率を入力してください"
              placeholder="5.0"
              required
              min="0"
              max="100"
              step="0.1"
            />
            <div className="mt-2 flex items-center gap-2">
              <div className="flex-1 bg-gray-100 rounded-lg p-3">
                <div className="text-xs text-gray-600 mb-1">保守的</div>
                <button
                  type="button"
                  onClick={() => setInvestmentReturn(3.0)}
                  className="text-sm text-primary-600 hover:text-primary-700 font-medium"
                >
                  3%
                </button>
              </div>
              <div className="flex-1 bg-gray-100 rounded-lg p-3">
                <div className="text-xs text-gray-600 mb-1">標準</div>
                <button
                  type="button"
                  onClick={() => setInvestmentReturn(5.0)}
                  className="text-sm text-primary-600 hover:text-primary-700 font-medium"
                >
                  5%
                </button>
              </div>
              <div className="flex-1 bg-gray-100 rounded-lg p-3">
                <div className="text-xs text-gray-600 mb-1">積極的</div>
                <button
                  type="button"
                  onClick={() => setInvestmentReturn(7.0)}
                  className="text-sm text-primary-600 hover:text-primary-700 font-medium"
                >
                  7%
                </button>
              </div>
            </div>
          </div>

          {/* インフレ率 */}
          <div>
            <InputField
              type="number"
              label="想定インフレ率（年率）"
              value={inflationRate || ''}
              onChange={(e) => setInflationRate(Number(e.target.value))}
              onBlur={() => handleBlur('inflation_rate')}
              error={errors.inflation_rate}
              helperText="将来の物価上昇率を入力してください"
              placeholder="2.0"
              required
              min="0"
              max="50"
              step="0.1"
            />
            <div className="mt-2 flex items-center gap-2">
              <div className="flex-1 bg-gray-100 rounded-lg p-3">
                <div className="text-xs text-gray-600 mb-1">低インフレ</div>
                <button
                  type="button"
                  onClick={() => setInflationRate(1.0)}
                  className="text-sm text-primary-600 hover:text-primary-700 font-medium"
                >
                  1%
                </button>
              </div>
              <div className="flex-1 bg-gray-100 rounded-lg p-3">
                <div className="text-xs text-gray-600 mb-1">標準</div>
                <button
                  type="button"
                  onClick={() => setInflationRate(2.0)}
                  className="text-sm text-primary-600 hover:text-primary-700 font-medium"
                >
                  2%
                </button>
              </div>
              <div className="flex-1 bg-gray-100 rounded-lg p-3">
                <div className="text-xs text-gray-600 mb-1">高インフレ</div>
                <button
                  type="button"
                  onClick={() => setInflationRate(3.0)}
                  className="text-sm text-primary-600 hover:text-primary-700 font-medium"
                >
                  3%
                </button>
              </div>
            </div>
          </div>

          {/* 実質利回り表示 */}
          <div className="bg-primary-50 border border-primary-200 rounded-lg p-4">
            <div className="flex justify-between items-center">
              <div>
                <div className="text-sm font-medium text-gray-700">実質利回り</div>
                <div className="text-xs text-gray-600 mt-1">
                  投資利回り - インフレ率
                </div>
              </div>
              <div className="text-2xl font-bold text-primary-600">
                {realReturn.toFixed(1)}%
              </div>
            </div>
            {realReturn < 0 && (
              <div className="mt-2 text-sm text-warning-600">
                ⚠️ 実質利回りがマイナスです。インフレにより資産価値が目減りする可能性があります。
              </div>
            )}
          </div>

          {/* 説明 */}
          <div className="bg-gray-50 rounded-lg p-4 text-sm text-gray-600">
            <h4 className="font-medium text-gray-900 mb-2">💡 設定のヒント</h4>
            <ul className="space-y-1 list-disc list-inside">
              <li>投資利回りは過去の実績や投資商品の特性を参考に設定してください</li>
              <li>インフレ率は日本銀行の目標値（2%）を基準に調整できます</li>
              <li>実質利回りは、インフレを考慮した実際の資産増加率を示します</li>
            </ul>
          </div>
        </div>
      </div>

      {/* 送信ボタン */}
      <div className="flex justify-end gap-3">
        <Button type="submit" loading={loading} disabled={loading}>
          設定を保存
        </Button>
      </div>
    </form>
  );
};

export default InvestmentSettingsForm;
