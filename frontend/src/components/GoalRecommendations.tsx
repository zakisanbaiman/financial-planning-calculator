'use client';

import React from 'react';
import type { Goal } from '@/types/api';

export interface GoalRecommendationsProps {
  goal: Goal;
  financialProfile?: {
    monthly_income: number;
    monthly_expenses: number;
    current_savings: number;
  };
}

interface Recommendation {
  type: 'success' | 'warning' | 'info' | 'error';
  title: string;
  message: string;
  action?: string;
}

const GoalRecommendations: React.FC<GoalRecommendationsProps> = ({ goal, financialProfile }) => {
  const calculateRecommendations = (): Recommendation[] => {
    const recommendations: Recommendation[] = [];
    const remainingAmount = Math.max(0, goal.target_amount - goal.current_amount);
    const targetDate = new Date(goal.target_date);
    const today = new Date();
    const monthsRemaining = Math.max(
      0,
      (targetDate.getFullYear() - today.getFullYear()) * 12 +
        (targetDate.getMonth() - today.getMonth())
    );
    const progress = goal.target_amount > 0 ? (goal.current_amount / goal.target_amount) * 100 : 0;

    // 目標達成済み
    if (progress >= 100) {
      recommendations.push({
        type: 'success',
        title: '目標達成おめでとうございます！',
        message: `${goal.title}の目標金額を達成しました。`,
        action: '次の目標を設定しましょう',
      });
      return recommendations;
    }

    // 期限切れ
    if (monthsRemaining <= 0) {
      recommendations.push({
        type: 'error',
        title: '目標期日を過ぎています',
        message: `目標期日を延長するか、目標金額を見直すことをお勧めします。`,
        action: '目標を編集',
      });
    }

    // 推奨月間積立額の計算
    const recommendedMonthly = monthsRemaining > 0 ? remainingAmount / monthsRemaining : 0;

    // 現在の積立額が不足している場合
    if (goal.monthly_contribution > 0 && goal.monthly_contribution < recommendedMonthly) {
      const shortfall = recommendedMonthly - goal.monthly_contribution;
      recommendations.push({
        type: 'warning',
        title: '月間積立額の増額を推奨',
        message: `目標達成には月額¥${Math.ceil(recommendedMonthly).toLocaleString()}の積立が必要です。現在より¥${Math.ceil(
          shortfall
        ).toLocaleString()}の増額をお勧めします。`,
        action: `月額¥${Math.ceil(shortfall).toLocaleString()}増額`,
      });
    }

    // 順調な進捗
    if (progress >= 75 && monthsRemaining > 3) {
      recommendations.push({
        type: 'success',
        title: '順調に進んでいます',
        message: `現在のペースを維持すれば、目標期日までに達成できる見込みです。`,
      });
    }

    // 進捗が遅い場合
    if (progress < 50 && monthsRemaining < 12) {
      recommendations.push({
        type: 'warning',
        title: '進捗が遅れています',
        message: `目標達成が困難な状況です。積立額の増額または目標の見直しを検討してください。`,
        action: '代替案を確認',
      });
    }

    // 財務プロファイルがある場合の追加推奨事項
    if (financialProfile) {
      const netSavings = financialProfile.monthly_income - financialProfile.monthly_expenses;
      const savingsRate = financialProfile.monthly_income > 0 
        ? (netSavings / financialProfile.monthly_income) * 100 
        : 0;

      // 貯蓄率が低い場合
      if (savingsRate < 10) {
        recommendations.push({
          type: 'warning',
          title: '貯蓄率が低い状態です',
          message: `現在の貯蓄率は${savingsRate.toFixed(1)}%です。支出を見直して貯蓄率を向上させることをお勧めします。`,
          action: '支出を見直す',
        });
      }

      // 推奨積立額が純貯蓄を超える場合
      if (recommendedMonthly > netSavings) {
        const deficit = recommendedMonthly - netSavings;
        recommendations.push({
          type: 'error',
          title: '収支の改善が必要です',
          message: `目標達成には月額¥${Math.ceil(recommendedMonthly).toLocaleString()}必要ですが、現在の純貯蓄は¥${netSavings.toLocaleString()}です。収入を増やすか支出を¥${Math.ceil(
            deficit
          ).toLocaleString()}削減する必要があります。`,
          action: '財務計画を見直す',
        });
      }

      // 緊急資金の確認
      if (goal.goal_type !== 'emergency' && financialProfile.current_savings < financialProfile.monthly_expenses * 3) {
        recommendations.push({
          type: 'info',
          title: '緊急資金の確保を優先',
          message: `この目標の前に、まず3〜6ヶ月分の生活費を緊急資金として確保することをお勧めします。`,
          action: '緊急資金目標を作成',
        });
      }
    }

    // 代替案の提案
    if (monthsRemaining > 0 && remainingAmount > 0) {
      const alternativeScenarios = [
        {
          months: monthsRemaining + 6,
          monthly: remainingAmount / (monthsRemaining + 6),
        },
        {
          months: monthsRemaining + 12,
          monthly: remainingAmount / (monthsRemaining + 12),
        },
      ];

      const easierScenario = alternativeScenarios.find(
        (s) => s.monthly < goal.monthly_contribution * 0.8
      );

      if (easierScenario) {
        const newDate = new Date(today);
        newDate.setMonth(newDate.getMonth() + easierScenario.months);
        recommendations.push({
          type: 'info',
          title: '代替案の提案',
          message: `目標期日を${newDate.toLocaleDateString('ja-JP')}に延長すると、月額¥${Math.ceil(
            easierScenario.monthly
          ).toLocaleString()}で達成可能です。`,
          action: '期日を延長',
        });
      }
    }

    // 投資による加速
    if (remainingAmount > 100000 && monthsRemaining > 12) {
      const withInvestment = remainingAmount * Math.pow(1.05, monthsRemaining / 12);
      const savingsFromInvestment = withInvestment - remainingAmount;
      if (savingsFromInvestment > 10000) {
        recommendations.push({
          type: 'info',
          title: '投資による目標達成の加速',
          message: `年利5%で運用すると、約¥${Math.ceil(savingsFromInvestment).toLocaleString()}の追加収益が見込めます。`,
          action: '投資プランを確認',
        });
      }
    }

    return recommendations;
  };

  const recommendations = calculateRecommendations();

  const iconMap = {
    success: '✅',
    warning: '⚠️',
    info: '💡',
    error: '❌',
  };

  const colorMap = {
    success: {
      bg: 'bg-success-50',
      border: 'border-success-200',
      text: 'text-success-800',
      title: 'text-success-900',
    },
    warning: {
      bg: 'bg-warning-50',
      border: 'border-warning-200',
      text: 'text-warning-800',
      title: 'text-warning-900',
    },
    info: {
      bg: 'bg-primary-50',
      border: 'border-primary-200',
      text: 'text-primary-800',
      title: 'text-primary-900',
    },
    error: {
      bg: 'bg-error-50',
      border: 'border-error-200',
      text: 'text-error-800',
      title: 'text-error-900',
    },
  };

  if (recommendations.length === 0) {
    return (
      <div className="card text-center py-8">
        <p className="text-gray-500">現在、推奨事項はありません</p>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {recommendations.map((rec, index) => {
        const colors = colorMap[rec.type];
        return (
          <div
            key={index}
            className={`p-4 rounded-lg border ${colors.bg} ${colors.border}`}
          >
            <div className="flex items-start gap-3">
              <span className="text-2xl flex-shrink-0">{iconMap[rec.type]}</span>
              <div className="flex-1">
                <h4 className={`font-semibold mb-1 ${colors.title}`}>{rec.title}</h4>
                <p className={`text-sm ${colors.text}`}>{rec.message}</p>
                {rec.action && (
                  <button
                    className={`mt-2 text-sm font-medium ${colors.text} hover:underline`}
                  >
                    {rec.action} →
                  </button>
                )}
              </div>
            </div>
          </div>
        );
      })}
    </div>
  );
};

export default GoalRecommendations;
