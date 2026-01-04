'use client';

import React, { createContext, useContext, useState, useCallback, useEffect } from 'react';

export interface TutorialStep {
  id: string;
  page: string;
  title: string;
  description: string;
  elementId?: string;
  position?: 'top' | 'bottom' | 'left' | 'right' | 'center';
}

interface TutorialContextType {
  isActive: boolean;
  currentStep: number;
  totalSteps: number;
  currentStepData: TutorialStep | null;
  startTutorial: () => void;
  nextStep: () => void;
  previousStep: () => void;
  skipTutorial: () => void;
  completeTutorial: () => void;
  resetTutorial: () => void;
}

const TutorialContext = createContext<TutorialContextType | undefined>(undefined);

const TUTORIAL_STORAGE_KEY = 'financial-calculator-tutorial-completed';

// チュートリアルステップの定義
const tutorialSteps: TutorialStep[] = [
  {
    id: 'welcome',
    page: '/',
    title: '財務計画計算機へようこそ！',
    description: 'このアプリケーションでは、将来の資産形成と老後の財務計画を簡単に可視化できます。まずは主要な機能をご案内します。',
    position: 'center',
  },
  {
    id: 'dashboard-intro',
    page: '/dashboard',
    title: 'ダッシュボード',
    description: 'ここでは財務状況の概要と主要な指標を一目で確認できます。月間純貯蓄、総資産、老後資金充足率などが表示されます。',
    position: 'center',
  },
  {
    id: 'financial-data',
    page: '/financial-data',
    title: '財務データの入力',
    description: '最初に、現在の収入、支出、貯蓄状況を入力しましょう。これが正確な将来予測の基盤となります。',
    position: 'center',
  },
  {
    id: 'financial-data-form',
    page: '/financial-data',
    title: '基本情報の入力',
    description: '月収、月間支出（住居費、食費など）、現在の貯蓄額を入力してください。投資利回りとインフレ率の設定も可能です。',
    position: 'center',
  },
  {
    id: 'calculations',
    page: '/calculations',
    title: '財務計算機',
    description: 'ここでは3つの重要な計算ができます：\n• 資産推移シミュレーション\n• 老後資金計算\n• 緊急資金計算',
    position: 'center',
  },
  {
    id: 'asset-projection',
    page: '/calculations',
    title: '資産推移シミュレーション',
    description: '現在の貯蓄ペースで将来どれだけ資産が増えるかをグラフで可視化します。10年後、20年後、30年後の資産額が一目瞭然です。',
    position: 'center',
  },
  {
    id: 'goals',
    page: '/goals',
    title: '目標設定・進捗管理',
    description: '具体的な財務目標（マイホーム購入、教育資金、老後資金など）を設定し、進捗を追跡できます。目標達成までの期間や必要な月間貯蓄額も計算されます。',
    position: 'center',
  },
  {
    id: 'reports',
    page: '/reports',
    title: 'レポート生成',
    description: '財務状況をまとめたレポートをPDF形式で生成・印刷できます。家族との共有や記録保存に便利です。',
    position: 'center',
  },
  {
    id: 'complete',
    page: '/',
    title: 'チュートリアル完了！',
    description: 'これで基本的な使い方は完了です。さっそく財務データを入力して、あなたの未来を計画しましょう！いつでもヘルプセクションから再度チュートリアルを表示できます。',
    position: 'center',
  },
];

export function TutorialProvider({ children }: { children: React.ReactNode }) {
  const [isActive, setIsActive] = useState(false);
  const [currentStep, setCurrentStep] = useState(0);
  const [hasCheckedStorage, setHasCheckedStorage] = useState(false);

  useEffect(() => {
    // 初回マウント時にチュートリアル完了状態をチェック
    if (typeof window !== 'undefined' && !hasCheckedStorage) {
      const completed = localStorage.getItem(TUTORIAL_STORAGE_KEY);
      // 完了していない場合は自動的にチュートリアルを開始
      if (!completed) {
        setIsActive(true);
      }
      setHasCheckedStorage(true);
    }
  }, [hasCheckedStorage]);

  const startTutorial = useCallback(() => {
    setCurrentStep(0);
    setIsActive(true);
  }, []);

  const completeTutorial = useCallback(() => {
    setIsActive(false);
    if (typeof window !== 'undefined') {
      localStorage.setItem(TUTORIAL_STORAGE_KEY, 'completed');
    }
  }, []);

  const nextStep = useCallback(() => {
    if (currentStep < tutorialSteps.length - 1) {
      setCurrentStep((prev) => prev + 1);
    } else {
      completeTutorial();
    }
  }, [currentStep, completeTutorial]);

  const previousStep = useCallback(() => {
    if (currentStep > 0) {
      setCurrentStep((prev) => prev - 1);
    }
  }, [currentStep]);

  const skipTutorial = useCallback(() => {
    setIsActive(false);
    if (typeof window !== 'undefined') {
      localStorage.setItem(TUTORIAL_STORAGE_KEY, 'skipped');
    }
  }, []);

  const resetTutorial = useCallback(() => {
    setCurrentStep(0);
    setIsActive(false);
    if (typeof window !== 'undefined') {
      localStorage.removeItem(TUTORIAL_STORAGE_KEY);
    }
  }, []);

  const currentStepData = tutorialSteps[currentStep] || null;

  return (
    <TutorialContext.Provider
      value={{
        isActive,
        currentStep,
        totalSteps: tutorialSteps.length,
        currentStepData,
        startTutorial,
        nextStep,
        previousStep,
        skipTutorial,
        completeTutorial,
        resetTutorial,
      }}
    >
      {children}
    </TutorialContext.Provider>
  );
}

export function useTutorial() {
  const context = useContext(TutorialContext);
  if (context === undefined) {
    throw new Error('useTutorial must be used within a TutorialProvider');
  }
  return context;
}
