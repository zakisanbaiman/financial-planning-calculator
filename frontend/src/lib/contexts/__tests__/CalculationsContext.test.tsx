import React from 'react';
import { renderHook, act } from '@testing-library/react';
import { CalculationsProvider, useCalculations } from '../CalculationsContext';
import { calculationsAPI } from '@/lib/api-client';
import type {
  AssetProjectionRequest,
  AssetProjectionResponse,
  RetirementCalculationRequest,
  RetirementCalculationResponse,
  EmergencyFundRequest,
  EmergencyFundResponse,
  GoalProjectionRequest,
  GoalProjectionResponse,
} from '@/types/api';

jest.mock('@/lib/api-client');
const mockedCalcAPI = calculationsAPI as jest.Mocked<typeof calculationsAPI>;

const wrapper = ({ children }: { children: React.ReactNode }) => (
  <CalculationsProvider>{children}</CalculationsProvider>
);

describe('CalculationsContext', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('初期状態', () => {
    it('全計算結果がnull、loading=false、error=null', () => {
      const { result } = renderHook(() => useCalculations(), { wrapper });
      expect(result.current.assetProjection).toBeNull();
      expect(result.current.retirementCalculation).toBeNull();
      expect(result.current.emergencyFund).toBeNull();
      expect(result.current.goalProjection).toBeNull();
      expect(result.current.loading).toBe(false);
      expect(result.current.error).toBeNull();
    });
  });

  describe('資産推移計算', () => {
    const mockRequest: AssetProjectionRequest = {
      user_id: 'user-1',
      years: 30,
      monthly_income: 400000,
      monthly_expenses: 250000,
      current_savings: 3000000,
      investment_return: 5.0,
      inflation_rate: 2.0,
    };

    const mockResponse: AssetProjectionResponse = {
      projections: [
        { year: 1, total_assets: 4800000, real_value: 4700000, contributed_amount: 4800000, investment_gains: 0 },
      ],
      final_amount: 50000000,
      total_contributions: 30000000,
      total_gains: 20000000,
    };

    it('calculateAssetProjectionで結果が保存される', async () => {
      mockedCalcAPI.assetProjection.mockResolvedValue(mockResponse);

      const { result } = renderHook(() => useCalculations(), { wrapper });

      await act(async () => {
        await result.current.calculateAssetProjection(mockRequest);
      });

      expect(result.current.assetProjection).toEqual(mockResponse);
    });

    it('エラー時にerrorが設定される', async () => {
      mockedCalcAPI.assetProjection.mockRejectedValue(new Error('計算失敗'));

      const { result } = renderHook(() => useCalculations(), { wrapper });

      await act(async () => {
        try {
          await result.current.calculateAssetProjection(mockRequest);
        } catch {
          // expected
        }
      });

      expect(result.current.error).toBe('計算失敗');
    });
  });

  describe('老後資金計算', () => {
    const mockRequest: RetirementCalculationRequest = {
      user_id: 'user-1',
      current_age: 35,
      retirement_age: 65,
      life_expectancy: 90,
      monthly_retirement_expenses: 250000,
      pension_amount: 150000,
      current_savings: 3000000,
      monthly_savings: 100000,
      investment_return: 5.0,
      inflation_rate: 2.0,
    };

    const mockResponse: RetirementCalculationResponse = {
      required_amount: 30000000,
      projected_amount: 35000000,
      shortfall: -5000000,
      sufficiency_rate: 116.7,
      recommended_monthly_savings: 80000,
      years_until_retirement: 30,
    };

    it('calculateRetirementで結果が保存される', async () => {
      mockedCalcAPI.retirement.mockResolvedValue(mockResponse);

      const { result } = renderHook(() => useCalculations(), { wrapper });

      await act(async () => {
        await result.current.calculateRetirement(mockRequest);
      });

      expect(result.current.retirementCalculation).toEqual(mockResponse);
    });
  });

  describe('緊急資金計算', () => {
    const mockRequest: EmergencyFundRequest = {
      user_id: 'user-1',
      monthly_expenses: 250000,
      target_months: 6,
      current_savings: 600000,
    };

    const mockResponse: EmergencyFundResponse = {
      required_amount: 1500000,
      current_amount: 600000,
      shortfall: 900000,
      sufficiency_rate: 40.0,
      months_to_target: 12,
    };

    it('calculateEmergencyFundで結果が保存される', async () => {
      mockedCalcAPI.emergencyFund.mockResolvedValue(mockResponse);

      const { result } = renderHook(() => useCalculations(), { wrapper });

      await act(async () => {
        await result.current.calculateEmergencyFund(mockRequest);
      });

      expect(result.current.emergencyFund).toEqual(mockResponse);
    });
  });

  describe('目標達成計算', () => {
    const mockRequest: GoalProjectionRequest = {
      user_id: 'user-1',
      goal_id: 'goal-1',
      target_amount: 5000000,
      target_date: '2027-12-31',
      current_amount: 1000000,
      monthly_contribution: 50000,
      investment_return: 5.0,
    };

    const mockResponse: GoalProjectionResponse = {
      goal_id: 'goal-1',
      is_achievable: true,
      projected_completion_date: '2027-06-15',
      shortfall: 0,
      recommended_monthly_contribution: 45000,
      progress_rate: 20.0,
    };

    it('calculateGoalProjectionで結果が保存される', async () => {
      mockedCalcAPI.goalProjection.mockResolvedValue(mockResponse);

      const { result } = renderHook(() => useCalculations(), { wrapper });

      await act(async () => {
        await result.current.calculateGoalProjection(mockRequest);
      });

      expect(result.current.goalProjection).toEqual(mockResponse);
    });
  });

  describe('クリア操作', () => {
    it('clearCalculationsで全計算結果がクリアされる', async () => {
      mockedCalcAPI.emergencyFund.mockResolvedValue({
        required_amount: 1500000,
        current_amount: 600000,
        shortfall: 900000,
        sufficiency_rate: 40.0,
        months_to_target: 12,
      });

      const { result } = renderHook(() => useCalculations(), { wrapper });

      await act(async () => {
        await result.current.calculateEmergencyFund({
          user_id: 'user-1',
          monthly_expenses: 250000,
          target_months: 6,
          current_savings: 600000,
        });
      });

      expect(result.current.emergencyFund).not.toBeNull();

      act(() => {
        result.current.clearCalculations();
      });

      expect(result.current.assetProjection).toBeNull();
      expect(result.current.retirementCalculation).toBeNull();
      expect(result.current.emergencyFund).toBeNull();
      expect(result.current.goalProjection).toBeNull();
      expect(result.current.error).toBeNull();
    });

    it('clearErrorでエラーのみクリアされる', async () => {
      mockedCalcAPI.assetProjection.mockRejectedValue(new Error('エラー'));

      const { result } = renderHook(() => useCalculations(), { wrapper });

      await act(async () => {
        try {
          await result.current.calculateAssetProjection({
            user_id: 'user-1',
            years: 30,
            monthly_income: 400000,
            monthly_expenses: 250000,
            current_savings: 3000000,
            investment_return: 5.0,
            inflation_rate: 2.0,
          });
        } catch {
          // expected
        }
      });

      act(() => {
        result.current.clearError();
      });

      expect(result.current.error).toBeNull();
    });
  });

  describe('Provider外でのフック使用', () => {
    it('Provider外で useCalculations を使うとエラーが発生する', () => {
      const consoleSpy = jest.spyOn(console, 'error').mockImplementation(() => {});
      expect(() => {
        renderHook(() => useCalculations());
      }).toThrow('useCalculations must be used within a CalculationsProvider');
      consoleSpy.mockRestore();
    });
  });
});
