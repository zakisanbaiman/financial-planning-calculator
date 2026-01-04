'use client';

import React, { useEffect } from 'react';
import { usePathname, useRouter } from 'next/navigation';
import { useTutorial } from '@/lib/contexts/TutorialContext';

export default function Tutorial() {
  const pathname = usePathname();
  const router = useRouter();
  const {
    isActive,
    currentStep,
    totalSteps,
    currentStepData,
    nextStep,
    previousStep,
    skipTutorial,
    completeTutorial,
  } = useTutorial();

  // ページ遷移が必要な場合の処理
  useEffect(() => {
    if (isActive && currentStepData && currentStepData.page !== pathname) {
      router.push(currentStepData.page);
    }
  }, [isActive, currentStepData, pathname, router]);

  if (!isActive || !currentStepData) {
    return null;
  }

  const isFirstStep = currentStep === 0;
  const isLastStep = currentStep === totalSteps - 1;
  const progress = ((currentStep + 1) / totalSteps) * 100;

  return (
    <>
      {/* Overlay */}
      <div className="fixed inset-0 bg-black/50 z-[100] transition-opacity" />

      {/* Tutorial Card */}
      <div className="fixed inset-0 z-[101] flex items-center justify-center p-4">
        <div className="bg-white dark:bg-gray-800 rounded-lg shadow-2xl max-w-2xl w-full max-h-[90vh] overflow-y-auto">
          {/* Header */}
          <div className="border-b border-gray-200 dark:border-gray-700 p-6">
            <div className="flex items-start justify-between mb-4">
              <div className="flex items-center space-x-3">
                <div className="w-10 h-10 bg-primary-500 rounded-full flex items-center justify-center text-white font-bold">
                  {currentStep + 1}
                </div>
                <div>
                  <h2 className="text-2xl font-bold text-gray-900 dark:text-white">
                    {currentStepData.title}
                  </h2>
                  <p className="text-sm text-gray-500 dark:text-gray-400">
                    ステップ {currentStep + 1} / {totalSteps}
                  </p>
                </div>
              </div>
              <button
                onClick={skipTutorial}
                className="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 transition-colors"
                aria-label="チュートリアルをスキップ"
              >
                <svg
                  className="w-6 h-6"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M6 18L18 6M6 6l12 12"
                  />
                </svg>
              </button>
            </div>

            {/* Progress Bar */}
            <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2">
              <div
                className="bg-primary-500 h-2 rounded-full transition-all duration-300"
                style={{ width: `${progress}%` }}
              />
            </div>
          </div>

          {/* Content */}
          <div className="p-6">
            <div className="prose dark:prose-invert max-w-none">
              <p className="text-gray-700 dark:text-gray-300 text-lg leading-relaxed whitespace-pre-line">
                {currentStepData.description}
              </p>
            </div>

            {/* Icon for specific steps */}
            {currentStep === 0 && (
              <div className="mt-6 text-center">
                <div className="text-6xl mb-4">🎓</div>
              </div>
            )}
            {currentStepData.id === 'dashboard-intro' && (
              <div className="mt-6 text-center">
                <div className="text-6xl mb-4">📊</div>
              </div>
            )}
            {currentStepData.id === 'financial-data' && (
              <div className="mt-6 text-center">
                <div className="text-6xl mb-4">💰</div>
              </div>
            )}
            {currentStepData.id === 'calculations' && (
              <div className="mt-6 text-center">
                <div className="text-6xl mb-4">📈</div>
              </div>
            )}
            {currentStepData.id === 'goals' && (
              <div className="mt-6 text-center">
                <div className="text-6xl mb-4">🎯</div>
              </div>
            )}
            {currentStepData.id === 'reports' && (
              <div className="mt-6 text-center">
                <div className="text-6xl mb-4">📋</div>
              </div>
            )}
            {isLastStep && (
              <div className="mt-6 text-center">
                <div className="text-6xl mb-4">🎉</div>
              </div>
            )}
          </div>

          {/* Footer */}
          <div className="border-t border-gray-200 dark:border-gray-700 p-6">
            <div className="flex items-center justify-between">
              <div className="flex items-center space-x-3">
                {!isFirstStep && (
                  <button
                    onClick={previousStep}
                    className="px-4 py-2 text-gray-600 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white font-medium transition-colors"
                  >
                    ← 前へ
                  </button>
                )}
                <button
                  onClick={skipTutorial}
                  className="px-4 py-2 text-gray-600 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white font-medium transition-colors"
                >
                  スキップ
                </button>
              </div>

              <button
                onClick={isLastStep ? completeTutorial : nextStep}
                className="btn-primary flex items-center space-x-2"
              >
                <span>{isLastStep ? '完了' : '次へ'}</span>
                {!isLastStep && (
                  <svg
                    className="w-5 h-5"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M9 5l7 7-7 7"
                    />
                  </svg>
                )}
                {isLastStep && (
                  <svg
                    className="w-5 h-5"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M5 13l4 4L19 7"
                    />
                  </svg>
                )}
              </button>
            </div>
          </div>
        </div>
      </div>
    </>
  );
}
