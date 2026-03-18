import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import RetirementCalculator from '../RetirementCalculator';
import { calculationsAPI } from '@/lib/api-client';
import type { RetirementCalculationResponse } from '@/types/api';

jest.mock('@/lib/api-client');
const mockedCalcAPI = calculationsAPI as jest.Mocked<typeof calculationsAPI>;

// requestAnimationFrameモック
beforeAll(() => {
  jest.spyOn(window, 'requestAnimationFrame').mockImplementation((cb) => {
    cb(0);
    return 0;
  });
});

afterAll(() => {
  (window.requestAnimationFrame as jest.Mock).mockRestore();
});

describe('RetirementCalculator', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('フォーム表示', () => {
    it('老後資金計算フォームが表示される', () => {
      render(<RetirementCalculator userId="user-1" />);
      expect(screen.getByText('老後資金計算')).toBeInTheDocument();
    });

    it('年齢フィールドが表示される', () => {
      render(<RetirementCalculator userId="user-1" />);
      // InputFieldはhtmlFor未設定のため、デフォルト値でinputを特定する
      expect(screen.getByDisplayValue('35')).toBeInTheDocument();
      expect(screen.getByDisplayValue('65')).toBeInTheDocument();
      expect(screen.getByDisplayValue('90')).toBeInTheDocument();
    });

    it('デフォルト値が設定されている', () => {
      render(<RetirementCalculator userId="user-1" />);
      expect(screen.getByDisplayValue('35')).toBeInTheDocument();
      expect(screen.getByDisplayValue('65')).toBeInTheDocument();
      expect(screen.getByDisplayValue('90')).toBeInTheDocument();
    });

    it('退職までの期間が表示される', () => {
      render(<RetirementCalculator userId="user-1" />);
      expect(screen.getByText('退職までの期間')).toBeInTheDocument();
      expect(screen.getByText('30年')).toBeInTheDocument();
    });

    it('退職後の期間が表示される', () => {
      render(<RetirementCalculator userId="user-1" />);
      expect(screen.getByText('退職後の期間')).toBeInTheDocument();
      expect(screen.getByText('25年')).toBeInTheDocument();
    });

    it('計算ボタンが表示される', () => {
      render(<RetirementCalculator userId="user-1" />);
      expect(screen.getByRole('button', { name: '計算する' })).toBeInTheDocument();
    });

    it('initialDataが反映される', () => {
      render(
        <RetirementCalculator
          userId="user-1"
          initialData={{ current_age: 40, retirement_age: 60 }}
        />
      );
      expect(screen.getByDisplayValue('40')).toBeInTheDocument();
      expect(screen.getByDisplayValue('60')).toBeInTheDocument();
    });
  });

  describe('API送信', () => {
    const mockResponse: RetirementCalculationResponse = {
      required_amount: 30000000,
      projected_amount: 35000000,
      shortfall: -5000000,
      sufficiency_rate: 116.7,
      recommended_monthly_savings: 80000,
      years_until_retirement: 30,
    };

    it('計算ボタンクリックでAPIが呼ばれる', async () => {
      const user = userEvent.setup();
      mockedCalcAPI.retirement.mockResolvedValue(mockResponse);

      render(<RetirementCalculator userId="user-1" />);
      await user.click(screen.getByRole('button', { name: '計算する' }));

      await waitFor(() => {
        expect(mockedCalcAPI.retirement).toHaveBeenCalledWith(
          expect.objectContaining({ user_id: 'user-1' })
        );
      });
    });

    it('計算中はLoadingSpinnerが表示される', async () => {
      const user = userEvent.setup();
      let resolveCalc: ((val: RetirementCalculationResponse) => void) | undefined;
      mockedCalcAPI.retirement.mockImplementation(
        () => new Promise((resolve) => { resolveCalc = resolve; })
      );

      render(<RetirementCalculator userId="user-1" />);
      await user.click(screen.getByRole('button', { name: '計算する' }));

      expect(screen.getByRole('status')).toBeInTheDocument();
      resolveCalc!(mockResponse);
    });
  });

  describe('結果表示', () => {
    const mockResponseSufficient: RetirementCalculationResponse = {
      required_amount: 30000000,
      projected_amount: 35000000,
      shortfall: -5000000,
      sufficiency_rate: 116.7,
      recommended_monthly_savings: 80000,
      years_until_retirement: 30,
    };

    it('充足率が表示される', async () => {
      const user = userEvent.setup();
      mockedCalcAPI.retirement.mockResolvedValue(mockResponseSufficient);

      render(<RetirementCalculator userId="user-1" />);
      await user.click(screen.getByRole('button', { name: '計算する' }));

      await waitFor(() => {
        expect(screen.getByText('116.7%')).toBeInTheDocument();
      });
    });

    it('老後資金計算結果セクションが表示される', async () => {
      const user = userEvent.setup();
      mockedCalcAPI.retirement.mockResolvedValue(mockResponseSufficient);

      render(<RetirementCalculator userId="user-1" />);
      await user.click(screen.getByRole('button', { name: '計算する' }));

      await waitFor(() => {
        expect(screen.getByText('老後資金計算結果')).toBeInTheDocument();
      });
    });

    it('資金充足時に良好な状態メッセージが表示される', async () => {
      const user = userEvent.setup();
      mockedCalcAPI.retirement.mockResolvedValue(mockResponseSufficient);

      render(<RetirementCalculator userId="user-1" />);
      await user.click(screen.getByRole('button', { name: '計算する' }));

      await waitFor(() => {
        expect(screen.getByText('✅ 良好な状態です')).toBeInTheDocument();
      });
    });

    it('資金不足時に推奨事項が表示される', async () => {
      const user = userEvent.setup();
      mockedCalcAPI.retirement.mockResolvedValue({
        ...mockResponseSufficient,
        shortfall: 5000000,
        sufficiency_rate: 85.0,
      });

      render(<RetirementCalculator userId="user-1" />);
      await user.click(screen.getByRole('button', { name: '計算する' }));

      await waitFor(() => {
        expect(screen.getByText('⚠️ 推奨事項')).toBeInTheDocument();
      });
    });
  });

  describe('エラーハンドリング', () => {
    it('APIエラー時にエラーメッセージが表示される', async () => {
      const user = userEvent.setup();
      mockedCalcAPI.retirement.mockRejectedValue(new Error('計算に失敗しました'));

      render(<RetirementCalculator userId="user-1" />);
      await user.click(screen.getByRole('button', { name: '計算する' }));

      await waitFor(() => {
        expect(screen.getByText('計算に失敗しました')).toBeInTheDocument();
      });
    });
  });
});
