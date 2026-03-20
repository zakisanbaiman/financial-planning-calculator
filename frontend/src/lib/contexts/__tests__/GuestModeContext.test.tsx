import React from 'react';
import { renderHook, act } from '@testing-library/react';
import { GuestModeProvider, useGuestMode } from '../GuestModeContext';
import type { FinancialData } from '@/types/api';

const wrapper = ({ children }: { children: React.ReactNode }) => (
  <GuestModeProvider>{children}</GuestModeProvider>
);

const mockFinancialData: FinancialData = {
  user_id: 'guest',
  profile: {
    monthly_income: 300000,
    monthly_expenses: [{ category: '生活費', amount: 100000 }],
    current_savings: [{ type: 'deposit', amount: 500000 }],
    investment_return: 5.0,
    inflation_rate: 2.0,
  },
};

describe('GuestModeContext', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    (localStorage.getItem as jest.Mock).mockReturnValue(null);
  });

  describe('初期状態', () => {
    it('デフォルトではゲストモードは無効', () => {
      const { result } = renderHook(() => useGuestMode(), { wrapper });
      expect(result.current.isGuestMode).toBe(false);
      expect(result.current.guestData).toBeNull();
    });
  });

  describe('ゲストモード開始', () => {
    it('startGuestModeでゲストモードが有効になる', () => {
      const { result } = renderHook(() => useGuestMode(), { wrapper });

      act(() => {
        result.current.startGuestMode();
      });

      expect(result.current.isGuestMode).toBe(true);
      expect(localStorage.setItem).toHaveBeenCalledWith('guest_mode', 'true');
    });
  });

  describe('ゲストモード終了', () => {
    it('exitGuestModeでゲストモードが無効になる', () => {
      const { result } = renderHook(() => useGuestMode(), { wrapper });

      act(() => {
        result.current.startGuestMode();
      });

      act(() => {
        result.current.exitGuestMode();
      });

      expect(result.current.isGuestMode).toBe(false);
      expect(result.current.guestData).toBeNull();
      expect(localStorage.removeItem).toHaveBeenCalledWith('guest_mode');
    });
  });

  describe('ゲストデータ管理', () => {
    it('setGuestDataでデータが保存される', () => {
      const { result } = renderHook(() => useGuestMode(), { wrapper });

      act(() => {
        result.current.setGuestData(mockFinancialData);
      });

      expect(result.current.guestData).toEqual(mockFinancialData);
      expect(localStorage.setItem).toHaveBeenCalledWith(
        'guest_financial_data',
        JSON.stringify(mockFinancialData)
      );
    });

    it('setGuestData(null)でデータが削除される', () => {
      const { result } = renderHook(() => useGuestMode(), { wrapper });

      act(() => {
        result.current.setGuestData(mockFinancialData);
      });
      act(() => {
        result.current.setGuestData(null);
      });

      expect(result.current.guestData).toBeNull();
      expect(localStorage.removeItem).toHaveBeenCalledWith('guest_financial_data');
    });

    it('clearGuestDataでデータがクリアされる', () => {
      const { result } = renderHook(() => useGuestMode(), { wrapper });

      act(() => {
        result.current.setGuestData(mockFinancialData);
      });
      act(() => {
        result.current.clearGuestData();
      });

      expect(result.current.guestData).toBeNull();
      expect(localStorage.removeItem).toHaveBeenCalledWith('guest_financial_data');
    });
  });

  describe('localStorage永続化', () => {
    it('ゲストモード状態がlocalStorageから復元される', () => {
      (localStorage.getItem as jest.Mock).mockImplementation((key: string) => {
        if (key === 'guest_mode') return 'true';
        if (key === 'guest_financial_data') return JSON.stringify(mockFinancialData);
        return null;
      });

      const { result } = renderHook(() => useGuestMode(), { wrapper });
      expect(result.current.isGuestMode).toBe(true);
      expect(result.current.guestData).toEqual(mockFinancialData);
    });

    it('不正なJSONデータの場合、データはクリアされる', () => {
      (localStorage.getItem as jest.Mock).mockImplementation((key: string) => {
        if (key === 'guest_mode') return 'true';
        if (key === 'guest_financial_data') return 'invalid json';
        return null;
      });

      const consoleSpy = jest.spyOn(console, 'error').mockImplementation(() => {});
      const { result } = renderHook(() => useGuestMode(), { wrapper });

      // パースエラーでデータがクリアされる
      expect(localStorage.removeItem).toHaveBeenCalledWith('guest_financial_data');
      consoleSpy.mockRestore();
    });

    it('user_idがないデータは無効として扱われる', () => {
      (localStorage.getItem as jest.Mock).mockImplementation((key: string) => {
        if (key === 'guest_mode') return 'true';
        if (key === 'guest_financial_data') return JSON.stringify({ name: 'invalid' });
        return null;
      });

      const consoleSpy = jest.spyOn(console, 'warn').mockImplementation(() => {});
      renderHook(() => useGuestMode(), { wrapper });

      expect(localStorage.removeItem).toHaveBeenCalledWith('guest_financial_data');
      consoleSpy.mockRestore();
    });
  });

  describe('Provider外でのフック使用', () => {
    it('Provider外で useGuestMode を使うとエラーが発生する', () => {
      const consoleSpy = jest.spyOn(console, 'error').mockImplementation(() => {});
      expect(() => {
        renderHook(() => useGuestMode());
      }).toThrow('useGuestMode must be used within a GuestModeProvider');
      consoleSpy.mockRestore();
    });
  });
});
