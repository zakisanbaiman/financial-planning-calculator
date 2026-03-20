import React from 'react';
import { renderHook, act } from '@testing-library/react';
import { ThemeProvider, useTheme } from '../ThemeContext';

const wrapper = ({ children }: { children: React.ReactNode }) => (
  <ThemeProvider>{children}</ThemeProvider>
);

describe('ThemeContext', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    document.documentElement.classList.remove('dark');
  });

  describe('初期状態', () => {
    it('デフォルトテーマはlightである', () => {
      const { result } = renderHook(() => useTheme(), { wrapper });
      expect(result.current.theme).toBe('light');
    });
  });

  describe('テーマ切替', () => {
    it('toggleThemeでlight→darkに切り替わる', () => {
      const { result } = renderHook(() => useTheme(), { wrapper });

      act(() => {
        result.current.toggleTheme();
      });

      expect(result.current.theme).toBe('dark');
    });

    it('toggleThemeでdark→lightに切り替わる', () => {
      const { result } = renderHook(() => useTheme(), { wrapper });

      act(() => {
        result.current.toggleTheme();
      });
      act(() => {
        result.current.toggleTheme();
      });

      expect(result.current.theme).toBe('light');
    });
  });

  describe('localStorage永続化', () => {
    it('テーマ切替時にlocalStorageに保存される', () => {
      const { result } = renderHook(() => useTheme(), { wrapper });

      act(() => {
        result.current.toggleTheme();
      });

      expect(localStorage.setItem).toHaveBeenCalledWith('theme', 'dark');
    });

    it('保存されたテーマが起動時に復元される', () => {
      (localStorage.getItem as jest.Mock).mockImplementation((key: string) => {
        if (key === 'theme') return 'dark';
        return null;
      });

      const { result } = renderHook(() => useTheme(), { wrapper });
      expect(result.current.theme).toBe('dark');
    });
  });

  describe('prefers-color-scheme', () => {
    it('localStorageにテーマがない場合、システム設定が使用される', () => {
      (localStorage.getItem as jest.Mock).mockReturnValue(null);
      // matchMediaモックでdarkを返す
      (window.matchMedia as jest.Mock).mockImplementation((query: string) => ({
        matches: query === '(prefers-color-scheme: dark)',
        media: query,
        onchange: null,
        addListener: jest.fn(),
        removeListener: jest.fn(),
        addEventListener: jest.fn(),
        removeEventListener: jest.fn(),
        dispatchEvent: jest.fn(),
      }));

      const { result } = renderHook(() => useTheme(), { wrapper });
      expect(result.current.theme).toBe('dark');
    });
  });

  describe('localStorage保存', () => {
    it('テーマ切替時にlocalStorageに保存される', () => {
      // matchMediaをリセットして確実にlightスタート
      (window.matchMedia as jest.Mock).mockImplementation((query: string) => ({
        matches: false,
        media: query,
        onchange: null,
        addListener: jest.fn(),
        removeListener: jest.fn(),
        addEventListener: jest.fn(),
        removeEventListener: jest.fn(),
        dispatchEvent: jest.fn(),
      }));
      (localStorage.getItem as jest.Mock).mockReturnValue(null);

      const { result } = renderHook(() => useTheme(), { wrapper });

      // 初期状態がlightであることを確認
      expect(result.current.theme).toBe('light');

      act(() => {
        result.current.toggleTheme();
      });

      expect(result.current.theme).toBe('dark');
      expect(localStorage.setItem).toHaveBeenCalledWith('theme', 'dark');
    });
  });

  describe('Provider外でのフック使用', () => {
    it('Provider外で useTheme を使うとエラーが発生する', () => {
      // console.errorを抑制
      const consoleSpy = jest.spyOn(console, 'error').mockImplementation(() => {});
      expect(() => {
        renderHook(() => useTheme());
      }).toThrow('useTheme must be used within a ThemeProvider');
      consoleSpy.mockRestore();
    });
  });
});
