import React, { ReactNode } from 'react';
import { renderHook, act, waitFor } from '@testing-library/react';
import { FinancialDataProvider, useFinancialData } from '../FinancialDataContext';
import { financialDataAPI } from '@/lib/api-client';
import type { FinancialData, FinancialProfile, RetirementData, EmergencyFund } from '@/types/api';

// API モック
jest.mock('@/lib/api-client', () => ({
  financialDataAPI: {
    get: jest.fn(),
    create: jest.fn(),
    updateProfile: jest.fn(),
    updateRetirement: jest.fn(),
    updateEmergencyFund: jest.fn(),
    delete: jest.fn(),
  },
}));

const mockedAPI = financialDataAPI as jest.Mocked<typeof financialDataAPI>;

// テスト用ラッパー
const wrapper = ({ children }: { children: ReactNode }) => (
  <FinancialDataProvider>{children}</FinancialDataProvider>
);

// テスト用データ
const mockFinancialData: FinancialData = {
  user_id: 'user_123',
  profile: {
    monthly_income: 500000,
    monthly_expenses: [
      { category: '生活費', amount: 200000 },
      { category: '住居費', amount: 100000 },
    ],
    current_savings: [
      { type: 'deposit', amount: 2000000 },
    ],
    investment_return: 5,
    inflation_rate: 2,
  } as FinancialProfile,
  retirement: {
    current_age: 30,
    retirement_age: 65,
    life_expectancy: 90,
    monthly_retirement_expenses: 250000,
    pension_amount: 150000,
  } as RetirementData,
  emergency_fund: {
    target_months: 6,
    monthly_expenses: 300000,
    current_amount: 1000000,
  } as EmergencyFund,
};

describe('FinancialDataContext', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('useFinancialData フック', () => {
    it('Provider の外で使用するとエラーがスローされる', () => {
      // コンソールエラーを抑制
      const consoleSpy = jest.spyOn(console, 'error').mockImplementation(() => {});
      
      expect(() => {
        renderHook(() => useFinancialData());
      }).toThrow('useFinancialData must be used within a FinancialDataProvider');
      
      consoleSpy.mockRestore();
    });

    it('初期状態が正しい', () => {
      const { result } = renderHook(() => useFinancialData(), { wrapper });
      
      expect(result.current.financialData).toBeNull();
      expect(result.current.loading).toBe(false);
      expect(result.current.error).toBeNull();
    });
  });

  describe('fetchFinancialData', () => {
    it('正常にデータを取得できる', async () => {
      mockedAPI.get.mockResolvedValue(mockFinancialData);
      
      const { result } = renderHook(() => useFinancialData(), { wrapper });
      
      await act(async () => {
        await result.current.fetchFinancialData('user_123');
      });
      
      expect(mockedAPI.get).toHaveBeenCalledWith('user_123');
      expect(result.current.financialData).toEqual(mockFinancialData);
      expect(result.current.loading).toBe(false);
      expect(result.current.error).toBeNull();
    });

    it('取得中は loading が true になる', async () => {
      let resolvePromise: (value: FinancialData) => void;
      mockedAPI.get.mockImplementation(() => 
        new Promise(resolve => { resolvePromise = resolve; })
      );
      
      const { result } = renderHook(() => useFinancialData(), { wrapper });
      
      act(() => {
        result.current.fetchFinancialData('user_123');
      });
      
      expect(result.current.loading).toBe(true);
      
      await act(async () => {
        resolvePromise!(mockFinancialData);
      });
      
      expect(result.current.loading).toBe(false);
    });

    it('404エラー時に適切なエラーメッセージが表示される', async () => {
      const error = new Error('Not found');
      (error as any).status = 404;
      mockedAPI.get.mockRejectedValue(error);
      
      const { result } = renderHook(() => useFinancialData(), { wrapper });
      
      let caughtError: Error | null = null;
      await act(async () => {
        try {
          await result.current.fetchFinancialData('user_123');
        } catch (e) {
          caughtError = e as Error;
        }
      });
      
      expect(caughtError).not.toBeNull();
      expect(result.current.error).toContain('まだ作成されていません');
    });

    it('一般エラー時にエラーメッセージが設定される', async () => {
      mockedAPI.get.mockRejectedValue(new Error('Network error'));
      
      const { result } = renderHook(() => useFinancialData(), { wrapper });
      
      let caughtError: Error | null = null;
      await act(async () => {
        try {
          await result.current.fetchFinancialData('user_123');
        } catch (e) {
          caughtError = e as Error;
        }
      });
      
      expect(caughtError).not.toBeNull();
      expect(result.current.error).toBe('Network error');
      expect(result.current.loading).toBe(false);
    });
  });

  describe('createFinancialData', () => {
    it('正常にデータを作成できる', async () => {
      mockedAPI.create.mockResolvedValue(mockFinancialData);
      
      const { result } = renderHook(() => useFinancialData(), { wrapper });
      
      await act(async () => {
        await result.current.createFinancialData(mockFinancialData);
      });
      
      expect(mockedAPI.create).toHaveBeenCalledWith(mockFinancialData);
      expect(result.current.financialData).toEqual(mockFinancialData);
    });

    it('作成失敗時にエラーが設定される', async () => {
      mockedAPI.create.mockRejectedValue(new Error('Creation failed'));
      
      const { result } = renderHook(() => useFinancialData(), { wrapper });
      
      let caughtError: Error | null = null;
      await act(async () => {
        try {
          await result.current.createFinancialData(mockFinancialData);
        } catch (e) {
          caughtError = e as Error;
        }
      });
      
      expect(caughtError).not.toBeNull();
      expect(result.current.error).toBe('Creation failed');
    });
  });

  describe('updateProfile', () => {
    it('正常にプロファイルを更新できる', async () => {
      const updatedProfile: FinancialProfile = {
        monthly_income: 550000,
        monthly_expenses: [
          { category: '生活費', amount: 200000 },
        ],
        current_savings: [
          { type: 'deposit', amount: 2500000 },
        ],
        investment_return: 6,
        inflation_rate: 2,
      };
      const updatedData = { ...mockFinancialData, profile: updatedProfile };
      mockedAPI.updateProfile.mockResolvedValue(updatedData);
      
      const { result } = renderHook(() => useFinancialData(), { wrapper });
      
      await act(async () => {
        await result.current.updateProfile('user_123', updatedProfile);
      });
      
      expect(mockedAPI.updateProfile).toHaveBeenCalledWith('user_123', updatedProfile);
      expect(result.current.financialData?.profile?.monthly_income).toBe(550000);
    });
  });

  describe('updateRetirement', () => {
    it('正常に退職データを更新できる', async () => {
      const updatedRetirement: RetirementData = {
        current_age: 30,
        retirement_age: 60,
        life_expectancy: 90,
        monthly_retirement_expenses: 250000,
        pension_amount: 150000,
      };
      const updatedData = { ...mockFinancialData, retirement: updatedRetirement };
      mockedAPI.updateRetirement.mockResolvedValue(updatedData);
      
      const { result } = renderHook(() => useFinancialData(), { wrapper });
      
      await act(async () => {
        await result.current.updateRetirement('user_123', updatedRetirement);
      });
      
      expect(mockedAPI.updateRetirement).toHaveBeenCalledWith('user_123', updatedRetirement);
      expect(result.current.financialData?.retirement?.retirement_age).toBe(60);
    });
  });

  describe('updateEmergencyFund', () => {
    it('正常に緊急資金を更新できる', async () => {
      const updatedEmergencyFund: EmergencyFund = {
        target_months: 12,
        monthly_expenses: 300000,
        current_amount: 1500000,
      };
      const updatedData = { ...mockFinancialData, emergency_fund: updatedEmergencyFund };
      mockedAPI.updateEmergencyFund.mockResolvedValue(updatedData);
      
      const { result } = renderHook(() => useFinancialData(), { wrapper });
      
      await act(async () => {
        await result.current.updateEmergencyFund('user_123', updatedEmergencyFund);
      });
      
      expect(mockedAPI.updateEmergencyFund).toHaveBeenCalledWith('user_123', updatedEmergencyFund);
      expect(result.current.financialData?.emergency_fund?.target_months).toBe(12);
    });
  });

  describe('deleteFinancialData', () => {
    it('正常にデータを削除できる', async () => {
      mockedAPI.delete.mockResolvedValue(undefined);
      mockedAPI.get.mockResolvedValue(mockFinancialData);
      
      const { result } = renderHook(() => useFinancialData(), { wrapper });
      
      // 先にデータを取得
      await act(async () => {
        await result.current.fetchFinancialData('user_123');
      });
      
      expect(result.current.financialData).not.toBeNull();
      
      // データを削除
      await act(async () => {
        await result.current.deleteFinancialData('user_123');
      });
      
      expect(mockedAPI.delete).toHaveBeenCalledWith('user_123');
      expect(result.current.financialData).toBeNull();
    });
  });

  describe('clearError', () => {
    it('エラーをクリアできる', async () => {
      mockedAPI.get.mockRejectedValue(new Error('Test error'));
      
      const { result } = renderHook(() => useFinancialData(), { wrapper });
      
      await act(async () => {
        try {
          await result.current.fetchFinancialData('user_123');
        } catch (e) {
          // エラーがスローされることを期待
        }
      });
      
      expect(result.current.error).toBe('Test error');
      
      act(() => {
        result.current.clearError();
      });
      
      expect(result.current.error).toBeNull();
    });
  });
});
