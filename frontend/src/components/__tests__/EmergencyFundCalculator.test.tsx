import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import EmergencyFundCalculator from '../EmergencyFundCalculator';
import { calculationsAPI } from '@/lib/api-client';
import type { EmergencyFundResponse } from '@/types/api';

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

describe('EmergencyFundCalculator', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('フォーム表示', () => {
    it('緊急資金計算フォームが表示される', () => {
      render(<EmergencyFundCalculator userId="user-1" />);
      expect(screen.getByText('緊急資金計算')).toBeInTheDocument();
    });

    it('緊急資金の説明が表示される', () => {
      render(<EmergencyFundCalculator userId="user-1" />);
      expect(screen.getByText('💡 緊急資金とは？')).toBeInTheDocument();
    });

    it('確保したい期間フィールドが表示される', () => {
      render(<EmergencyFundCalculator userId="user-1" />);
      // InputFieldはhtmlFor未設定のため、デフォルト値でinputを特定する
      expect(screen.getByDisplayValue('6')).toBeInTheDocument();
    });

    it('計算ボタンが表示される', () => {
      render(<EmergencyFundCalculator userId="user-1" />);
      expect(screen.getByRole('button', { name: '計算する' })).toBeInTheDocument();
    });

    it('目標緊急資金額が表示される', () => {
      render(<EmergencyFundCalculator userId="user-1" />);
      expect(screen.getByText('目標緊急資金額')).toBeInTheDocument();
    });

    it('initialDataが反映される', () => {
      render(
        <EmergencyFundCalculator
          userId="user-1"
          initialData={{ monthly_expenses: 200000, target_months: 3 }}
        />
      );
      expect(screen.getByDisplayValue('3')).toBeInTheDocument();
    });
  });

  describe('API送信', () => {
    const mockResponse: EmergencyFundResponse = {
      required_amount: 1680000,
      current_amount: 600000,
      shortfall: 1080000,
      sufficiency_rate: 35.7,
      months_to_target: 18,
    };

    it('計算ボタンクリックでAPIが呼ばれる', async () => {
      const user = userEvent.setup();
      mockedCalcAPI.emergencyFund.mockResolvedValue(mockResponse);

      render(<EmergencyFundCalculator userId="user-1" />);
      await user.click(screen.getByRole('button', { name: '計算する' }));

      await waitFor(() => {
        expect(mockedCalcAPI.emergencyFund).toHaveBeenCalledWith(
          expect.objectContaining({ user_id: 'user-1' })
        );
      });
    });
  });

  describe('結果表示', () => {
    const mockResponseInsufficient: EmergencyFundResponse = {
      required_amount: 1680000,
      current_amount: 600000,
      shortfall: 1080000,
      sufficiency_rate: 35.7,
      months_to_target: 18,
    };

    it('充足率が表示される', async () => {
      const user = userEvent.setup();
      mockedCalcAPI.emergencyFund.mockResolvedValue(mockResponseInsufficient);

      render(<EmergencyFundCalculator userId="user-1" />);
      await user.click(screen.getByRole('button', { name: '計算する' }));

      await waitFor(() => {
        expect(screen.getByText('35.7%')).toBeInTheDocument();
      });
    });

    it('緊急資金充足状況セクションが表示される', async () => {
      const user = userEvent.setup();
      mockedCalcAPI.emergencyFund.mockResolvedValue(mockResponseInsufficient);

      render(<EmergencyFundCalculator userId="user-1" />);
      await user.click(screen.getByRole('button', { name: '計算する' }));

      await waitFor(() => {
        expect(screen.getByText('緊急資金充足状況')).toBeInTheDocument();
      });
    });

    it('不足時に警告メッセージが表示される', async () => {
      const user = userEvent.setup();
      mockedCalcAPI.emergencyFund.mockResolvedValue(mockResponseInsufficient);

      render(<EmergencyFundCalculator userId="user-1" />);
      await user.click(screen.getByRole('button', { name: '計算する' }));

      await waitFor(() => {
        expect(screen.getByText('⚠️ 緊急資金が不足しています')).toBeInTheDocument();
      });
    });

    it('充足時に成功メッセージが表示される', async () => {
      const user = userEvent.setup();
      mockedCalcAPI.emergencyFund.mockResolvedValue({
        ...mockResponseInsufficient,
        shortfall: 0,
        sufficiency_rate: 110.0,
      });

      render(<EmergencyFundCalculator userId="user-1" />);
      await user.click(screen.getByRole('button', { name: '計算する' }));

      await waitFor(() => {
        expect(
          screen.getByText('✅ 緊急資金は十分に確保されています')
        ).toBeInTheDocument();
      });
    });

    it('現在カバーできる期間と目標期間が表示される', async () => {
      const user = userEvent.setup();
      mockedCalcAPI.emergencyFund.mockResolvedValue(mockResponseInsufficient);

      render(<EmergencyFundCalculator userId="user-1" />);
      await user.click(screen.getByRole('button', { name: '計算する' }));

      await waitFor(() => {
        expect(screen.getByText('現在カバーできる期間')).toBeInTheDocument();
        expect(screen.getByText('目標期間')).toBeInTheDocument();
      });
    });
  });

  describe('エラーハンドリング', () => {
    it('APIエラー時にエラーメッセージが表示される', async () => {
      const user = userEvent.setup();
      mockedCalcAPI.emergencyFund.mockRejectedValue(new Error('計算に失敗'));

      render(<EmergencyFundCalculator userId="user-1" />);
      await user.click(screen.getByRole('button', { name: '計算する' }));

      await waitFor(() => {
        expect(screen.getByText('計算に失敗')).toBeInTheDocument();
      });
    });
  });
});
