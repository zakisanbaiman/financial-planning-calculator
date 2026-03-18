import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import InvestmentSettingsForm from '../InvestmentSettingsForm';

const mockOnSubmit = jest.fn();

describe('InvestmentSettingsForm', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('フォーム表示', () => {
    it('投資・インフレ設定のタイトルが表示される', () => {
      render(<InvestmentSettingsForm onSubmit={mockOnSubmit} />);
      expect(screen.getByText('投資・インフレ設定')).toBeInTheDocument();
    });

    it('期待投資利回りフィールドがデフォルト値5.0で表示される', () => {
      render(<InvestmentSettingsForm onSubmit={mockOnSubmit} />);
      expect(screen.getByDisplayValue('5')).toBeInTheDocument();
    });

    it('インフレ率フィールドがデフォルト値2.0で表示される', () => {
      render(<InvestmentSettingsForm onSubmit={mockOnSubmit} />);
      expect(screen.getByDisplayValue('2')).toBeInTheDocument();
    });

    it('設定を保存ボタンが表示される', () => {
      render(<InvestmentSettingsForm onSubmit={mockOnSubmit} />);
      expect(screen.getByRole('button', { name: '設定を保存' })).toBeInTheDocument();
    });

    it('実質利回りが表示される（5 - 2 = 3.0%）', () => {
      render(<InvestmentSettingsForm onSubmit={mockOnSubmit} />);
      expect(screen.getByText('3.0%')).toBeInTheDocument();
    });

    it('initialDataが反映される', () => {
      render(
        <InvestmentSettingsForm
          onSubmit={mockOnSubmit}
          initialData={{ investment_return: 7.0, inflation_rate: 1.5 }}
        />
      );
      expect(screen.getByDisplayValue('7')).toBeInTheDocument();
      expect(screen.getByDisplayValue('1.5')).toBeInTheDocument();
    });
  });

  describe('プリセットボタン', () => {
    it('投資利回りのプリセットボタン（3%/5%/7%）が表示される', () => {
      render(<InvestmentSettingsForm onSubmit={mockOnSubmit} />);
      // 3%は投資利回りとインフレ率の両方に存在するためAllByを使用
      const buttons3percent = screen.getAllByRole('button', { name: '3%' });
      expect(buttons3percent).toHaveLength(2);
      expect(screen.getByRole('button', { name: '5%' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: '7%' })).toBeInTheDocument();
    });

    it('インフレ率のプリセットボタン（1%/2%/3%）が表示される', () => {
      render(<InvestmentSettingsForm onSubmit={mockOnSubmit} />);
      expect(screen.getByRole('button', { name: '1%' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: '2%' })).toBeInTheDocument();
      // 3%は2つある
      expect(screen.getAllByRole('button', { name: '3%' })).toHaveLength(2);
    });

    it('積極的プリセット（7%）をクリックすると投資利回りが更新される', async () => {
      const user = userEvent.setup();
      render(<InvestmentSettingsForm onSubmit={mockOnSubmit} />);

      await user.click(screen.getByRole('button', { name: '7%' }));

      // 実質利回りが 7 - 2 = 5.0% になる
      await waitFor(() => {
        expect(screen.getByText('5.0%')).toBeInTheDocument();
      });
    });

    it('低インフレプリセット（1%）をクリックするとインフレ率が更新される', async () => {
      const user = userEvent.setup();
      render(<InvestmentSettingsForm onSubmit={mockOnSubmit} />);

      await user.click(screen.getByRole('button', { name: '1%' }));

      // 実質利回りが 5 - 1 = 4.0% になる
      await waitFor(() => {
        expect(screen.getByText('4.0%')).toBeInTheDocument();
      });
    });
  });

  describe('実質利回り計算', () => {
    it('実質利回りがマイナスの場合に警告メッセージが表示される', async () => {
      const user = userEvent.setup();
      render(
        <InvestmentSettingsForm
          onSubmit={mockOnSubmit}
          // 投資利回り1% < インフレ率2% → 実質利回りマイナス
          initialData={{ investment_return: 1.0, inflation_rate: 2.0 }}
        />
      );

      // 初期状態で既にマイナスになっているはず
      await waitFor(() => {
        expect(
          screen.getByText(/実質利回りがマイナスです/)
        ).toBeInTheDocument();
      });
    });

    it('高インフレプリセット（3%）クリックで実質利回りが変わる', async () => {
      const user = userEvent.setup();
      render(<InvestmentSettingsForm onSubmit={mockOnSubmit} />);

      // インフレ率を3%にする（2番目の3%ボタン）
      const buttons3percent = screen.getAllByRole('button', { name: '3%' });
      // 2番目が高インフレプリセット（インフレ率3%）
      await user.click(buttons3percent[1]);

      // 実質利回りが 5 - 3 = 2.0% になる
      await waitFor(() => {
        expect(screen.getByText('2.0%')).toBeInTheDocument();
      });
    });
  });

  describe('フォーム送信', () => {
    it('設定を保存ボタンクリックで onSubmit が呼ばれる', async () => {
      const user = userEvent.setup();
      mockOnSubmit.mockResolvedValue(undefined);

      render(<InvestmentSettingsForm onSubmit={mockOnSubmit} />);

      await user.click(screen.getByRole('button', { name: '設定を保存' }));

      await waitFor(() => {
        expect(mockOnSubmit).toHaveBeenCalledWith({
          investment_return: 5,
          inflation_rate: 2,
        });
      });
    });

    it('loading=trueのとき保存ボタンが無効になる', () => {
      render(<InvestmentSettingsForm onSubmit={mockOnSubmit} loading={true} />);
      expect(screen.getByRole('button', { name: '設定を保存' })).toBeDisabled();
    });
  });
});
