import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import Tutorial from '../Tutorial';
import { TutorialProvider } from '@/lib/contexts/TutorialContext';

// next/navigation モック
const mockPush = jest.fn();
const mockUsePathname = jest.fn();
jest.mock('next/navigation', () => ({
  usePathname: () => mockUsePathname(),
  useRouter: () => ({ push: mockPush }),
}));

function renderTutorial() {
  return render(
    <TutorialProvider>
      <Tutorial />
    </TutorialProvider>
  );
}

describe('Tutorial', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    mockUsePathname.mockReturnValue('/');
    // チュートリアル未完了状態（自動開始される）
    (localStorage.getItem as jest.Mock).mockReturnValue(null);
  });

  describe('自動表示', () => {
    it('チュートリアル未完了の場合、自動的に表示される', () => {
      renderTutorial();
      expect(screen.getByText('財務計画計算機へようこそ！')).toBeInTheDocument();
    });

    it('チュートリアル完了済みの場合、表示されない', () => {
      (localStorage.getItem as jest.Mock).mockImplementation((key: string) => {
        if (key === 'financial-calculator-tutorial-completed') return 'completed';
        return null;
      });
      renderTutorial();
      expect(screen.queryByText('財務計画計算機へようこそ！')).not.toBeInTheDocument();
    });
  });

  describe('ステップ表示', () => {
    it('ステップ番号が表示される', () => {
      renderTutorial();
      expect(screen.getByText('ステップ 1 / 9')).toBeInTheDocument();
    });

    it('進捗バーが表示される', () => {
      const { container } = renderTutorial();
      const progressBar = container.querySelector('.bg-primary-500.h-2');
      expect(progressBar).toBeInTheDocument();
    });
  });

  describe('ステップ進行', () => {
    it('「次へ」ボタンで次のステップに進む', async () => {
      const user = userEvent.setup();
      renderTutorial();

      expect(screen.getByText('財務計画計算機へようこそ！')).toBeInTheDocument();

      await user.click(screen.getByText('次へ'));
      expect(screen.getByText('ダッシュボード')).toBeInTheDocument();
      expect(screen.getByText('ステップ 2 / 9')).toBeInTheDocument();
    });

    it('「前へ」ボタンで前のステップに戻る', async () => {
      const user = userEvent.setup();
      renderTutorial();

      // ステップ2に進む
      await user.click(screen.getByText('次へ'));
      expect(screen.getByText('ダッシュボード')).toBeInTheDocument();

      // ステップ1に戻る
      await user.click(screen.getByText('← 前へ'));
      expect(screen.getByText('財務計画計算機へようこそ！')).toBeInTheDocument();
    });

    it('最初のステップでは「前へ」ボタンが表示されない', () => {
      renderTutorial();
      expect(screen.queryByText('← 前へ')).not.toBeInTheDocument();
    });
  });

  describe('スキップ', () => {
    it('「スキップ」ボタンでチュートリアルが閉じる', async () => {
      const user = userEvent.setup();
      renderTutorial();

      await user.click(screen.getByText('スキップ'));
      expect(screen.queryByText('財務計画計算機へようこそ！')).not.toBeInTheDocument();
    });

    it('スキップ時にlocalStorageに保存される', async () => {
      const user = userEvent.setup();
      renderTutorial();

      await user.click(screen.getByText('スキップ'));
      expect(localStorage.setItem).toHaveBeenCalledWith(
        'financial-calculator-tutorial-completed',
        'skipped'
      );
    });
  });

  describe('完了', () => {
    it('最後のステップで「完了」ボタンが表示される', async () => {
      const user = userEvent.setup();
      renderTutorial();

      // 全ステップを進む (9ステップ)
      for (let i = 0; i < 8; i++) {
        await user.click(screen.getByText('次へ'));
      }

      expect(screen.getByText('完了')).toBeInTheDocument();
      expect(screen.getByText('チュートリアル完了！')).toBeInTheDocument();
    });

    it('「完了」ボタンでチュートリアルが閉じる', async () => {
      const user = userEvent.setup();
      renderTutorial();

      // 全ステップを進む
      for (let i = 0; i < 8; i++) {
        await user.click(screen.getByText('次へ'));
      }

      await user.click(screen.getByText('完了'));
      expect(screen.queryByText('チュートリアル完了！')).not.toBeInTheDocument();
    });
  });

  describe('スキップボタン（Xアイコン）', () => {
    it('Xボタンでチュートリアルをスキップできる', async () => {
      const user = userEvent.setup();
      renderTutorial();

      await user.click(screen.getByLabelText('チュートリアルをスキップ'));
      expect(screen.queryByText('財務計画計算機へようこそ！')).not.toBeInTheDocument();
    });
  });
});
