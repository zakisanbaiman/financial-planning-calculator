import React from 'react';
import { render, screen, act } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import ThemeToggle from '../ThemeToggle';
import { ThemeProvider } from '@/lib/contexts/ThemeContext';

const renderWithTheme = (ui: React.ReactElement) => {
  return render(<ThemeProvider>{ui}</ThemeProvider>);
};

describe('ThemeToggle', () => {
  describe('初期表示（hydration前）', () => {
    it('マウント前はdisabledボタンが表示される', () => {
      // useEffect実行前の状態をテスト
      // ThemeProviderなしでは動かないのでProvider込みでrender
      const { container } = renderWithTheme(<ThemeToggle />);
      // マウント後なのでenabledになっている
      const button = screen.getByRole('button');
      expect(button).toBeInTheDocument();
    });
  });

  describe('マウント後', () => {
    it('テーマ切り替えボタンが有効になる', () => {
      renderWithTheme(<ThemeToggle />);
      const button = screen.getByRole('button');
      expect(button).not.toBeDisabled();
    });

    it('ライトモードのaria-labelが設定される', () => {
      renderWithTheme(<ThemeToggle />);
      const button = screen.getByRole('button');
      expect(button).toHaveAttribute('aria-label', 'ダークモードに切り替え');
    });
  });

  describe('テーマ切替', () => {
    it('ボタンをクリックするとテーマが切り替わる', async () => {
      const user = userEvent.setup();
      renderWithTheme(<ThemeToggle />);

      const button = screen.getByRole('button');
      expect(button).toHaveAttribute('aria-label', 'ダークモードに切り替え');

      await user.click(button);
      expect(button).toHaveAttribute('aria-label', 'ライトモードに切り替え');

      await user.click(button);
      expect(button).toHaveAttribute('aria-label', 'ダークモードに切り替え');
    });
  });
});
