import React from 'react';
import { renderHook, act } from '@testing-library/react';
import { TutorialProvider, useTutorial } from '../TutorialContext';

const wrapper = ({ children }: { children: React.ReactNode }) => (
  <TutorialProvider>{children}</TutorialProvider>
);

describe('TutorialContext', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    (localStorage.getItem as jest.Mock).mockReturnValue(null);
  });

  describe('初期状態', () => {
    it('チュートリアル未完了の場合、自動的にアクティブになる', () => {
      const { result } = renderHook(() => useTutorial(), { wrapper });
      expect(result.current.isActive).toBe(true);
      expect(result.current.currentStep).toBe(0);
    });

    it('チュートリアル完了済みの場合、アクティブにならない', () => {
      (localStorage.getItem as jest.Mock).mockImplementation((key: string) => {
        if (key === 'financial-calculator-tutorial-completed') return 'completed';
        return null;
      });
      const { result } = renderHook(() => useTutorial(), { wrapper });
      expect(result.current.isActive).toBe(false);
    });

    it('totalStepsは9である', () => {
      const { result } = renderHook(() => useTutorial(), { wrapper });
      expect(result.current.totalSteps).toBe(9);
    });

    it('currentStepDataが最初のステップデータを返す', () => {
      const { result } = renderHook(() => useTutorial(), { wrapper });
      expect(result.current.currentStepData).toEqual(
        expect.objectContaining({
          id: 'welcome',
          title: '財務計画計算機へようこそ！',
        })
      );
    });
  });

  describe('ステップ進行', () => {
    it('nextStepで次のステップに進む', () => {
      const { result } = renderHook(() => useTutorial(), { wrapper });

      act(() => {
        result.current.nextStep();
      });

      expect(result.current.currentStep).toBe(1);
      expect(result.current.currentStepData?.id).toBe('dashboard-intro');
    });

    it('previousStepで前のステップに戻る', () => {
      const { result } = renderHook(() => useTutorial(), { wrapper });

      act(() => {
        result.current.nextStep();
      });
      act(() => {
        result.current.previousStep();
      });

      expect(result.current.currentStep).toBe(0);
    });

    it('最初のステップでpreviousStepを呼んでも0のまま', () => {
      const { result } = renderHook(() => useTutorial(), { wrapper });

      act(() => {
        result.current.previousStep();
      });

      expect(result.current.currentStep).toBe(0);
    });

    it('最後のステップでnextStepを呼ぶと完了する', () => {
      const { result } = renderHook(() => useTutorial(), { wrapper });

      // 最後のステップまで進む
      for (let i = 0; i < 8; i++) {
        act(() => {
          result.current.nextStep();
        });
      }

      expect(result.current.currentStep).toBe(8);

      // nextStepで完了
      act(() => {
        result.current.nextStep();
      });

      expect(result.current.isActive).toBe(false);
      expect(localStorage.setItem).toHaveBeenCalledWith(
        'financial-calculator-tutorial-completed',
        'completed'
      );
    });
  });

  describe('スキップ', () => {
    it('skipTutorialでチュートリアルが非アクティブになる', () => {
      const { result } = renderHook(() => useTutorial(), { wrapper });

      act(() => {
        result.current.skipTutorial();
      });

      expect(result.current.isActive).toBe(false);
      expect(localStorage.setItem).toHaveBeenCalledWith(
        'financial-calculator-tutorial-completed',
        'skipped'
      );
    });
  });

  describe('完了', () => {
    it('completeTutorialでチュートリアルが完了する', () => {
      const { result } = renderHook(() => useTutorial(), { wrapper });

      act(() => {
        result.current.completeTutorial();
      });

      expect(result.current.isActive).toBe(false);
      expect(localStorage.setItem).toHaveBeenCalledWith(
        'financial-calculator-tutorial-completed',
        'completed'
      );
    });
  });

  describe('リセット', () => {
    it('resetTutorialでステップがリセットされる', () => {
      const { result } = renderHook(() => useTutorial(), { wrapper });

      act(() => {
        result.current.nextStep();
        result.current.nextStep();
      });
      act(() => {
        result.current.resetTutorial();
      });

      expect(result.current.currentStep).toBe(0);
      expect(result.current.isActive).toBe(false);
      expect(localStorage.removeItem).toHaveBeenCalledWith(
        'financial-calculator-tutorial-completed'
      );
    });
  });

  describe('再開始', () => {
    it('startTutorialでチュートリアルを再開始できる', () => {
      (localStorage.getItem as jest.Mock).mockImplementation((key: string) => {
        if (key === 'financial-calculator-tutorial-completed') return 'completed';
        return null;
      });

      const { result } = renderHook(() => useTutorial(), { wrapper });
      expect(result.current.isActive).toBe(false);

      act(() => {
        result.current.startTutorial();
      });

      expect(result.current.isActive).toBe(true);
      expect(result.current.currentStep).toBe(0);
    });
  });

  describe('Provider外でのフック使用', () => {
    it('Provider外で useTutorial を使うとエラーが発生する', () => {
      const consoleSpy = jest.spyOn(console, 'error').mockImplementation(() => {});
      expect(() => {
        renderHook(() => useTutorial());
      }).toThrow('useTutorial must be used within a TutorialProvider');
      consoleSpy.mockRestore();
    });
  });
});
