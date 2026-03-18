import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import FinancialInputForm from '../FinancialInputForm';

// requestAnimationFrameモック（CurrencyInput使用のため）
beforeAll(() => {
  jest.spyOn(window, 'requestAnimationFrame').mockImplementation((cb) => {
    cb(0);
    return 0;
  });
});

afterAll(() => {
  (window.requestAnimationFrame as jest.Mock).mockRestore();
});

const mockOnSubmit = jest.fn();

describe('FinancialInputForm', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('フォーム表示', () => {
    it('月収セクションが表示される', () => {
      render(<FinancialInputForm onSubmit={mockOnSubmit} />);
      expect(screen.getByText('月収')).toBeInTheDocument();
    });

    it('月間支出セクションが表示される', () => {
      render(<FinancialInputForm onSubmit={mockOnSubmit} />);
      expect(screen.getByText('月間支出')).toBeInTheDocument();
    });

    it('現在の貯蓄セクションが表示される', () => {
      render(<FinancialInputForm onSubmit={mockOnSubmit} />);
      expect(screen.getByText('現在の貯蓄')).toBeInTheDocument();
    });

    it('保存ボタンが表示される', () => {
      render(<FinancialInputForm onSubmit={mockOnSubmit} />);
      expect(screen.getByRole('button', { name: '保存' })).toBeInTheDocument();
    });

    it('デフォルトで生活費カテゴリーが表示される', () => {
      render(<FinancialInputForm onSubmit={mockOnSubmit} />);
      expect(screen.getByDisplayValue('生活費')).toBeInTheDocument();
    });

    it('合計支出が表示される', () => {
      render(<FinancialInputForm onSubmit={mockOnSubmit} />);
      expect(screen.getByText('合計支出')).toBeInTheDocument();
    });

    it('総資産が表示される', () => {
      render(<FinancialInputForm onSubmit={mockOnSubmit} />);
      expect(screen.getByText('総資産')).toBeInTheDocument();
    });

    it('月間純貯蓄が表示される', () => {
      render(<FinancialInputForm onSubmit={mockOnSubmit} />);
      expect(screen.getByText('月間純貯蓄')).toBeInTheDocument();
    });
  });

  describe('動的項目追加', () => {
    it('支出の「+ 項目追加」ボタンが表示される', () => {
      render(<FinancialInputForm onSubmit={mockOnSubmit} />);
      // 月間支出と現在の貯蓄の両方に「+ 項目追加」ボタンがある
      const addButtons = screen.getAllByRole('button', { name: '+ 項目追加' });
      expect(addButtons).toHaveLength(2);
    });

    it('支出の項目追加ボタンをクリックするとフィールドが増える', async () => {
      const user = userEvent.setup();
      render(<FinancialInputForm onSubmit={mockOnSubmit} />);

      // 初期状態はカテゴリーが1つ
      expect(screen.getAllByPlaceholderText('例: 住居費')).toHaveLength(1);

      const addButtons = screen.getAllByRole('button', { name: '+ 項目追加' });
      await user.click(addButtons[0]); // 最初の項目追加ボタン（支出）

      // カテゴリーが2つになる
      await waitFor(() => {
        expect(screen.getAllByPlaceholderText('例: 住居費')).toHaveLength(2);
      });
    });

    it('支出が2つ以上のとき削除ボタンが表示される', async () => {
      const user = userEvent.setup();
      render(<FinancialInputForm onSubmit={mockOnSubmit} />);

      const addButtons = screen.getAllByRole('button', { name: '+ 項目追加' });
      await user.click(addButtons[0]);

      await waitFor(() => {
        // aria-label="削除"ボタンが表示される
        const deleteButtons = screen.getAllByRole('button', { name: '削除' });
        expect(deleteButtons.length).toBeGreaterThan(0);
      });
    });

    it('削除ボタンをクリックすると支出項目が減る', async () => {
      const user = userEvent.setup();
      render(<FinancialInputForm onSubmit={mockOnSubmit} />);

      const addButtons = screen.getAllByRole('button', { name: '+ 項目追加' });
      await user.click(addButtons[0]);

      await waitFor(() => {
        expect(screen.getAllByPlaceholderText('例: 住居費')).toHaveLength(2);
      });

      const deleteButtons = screen.getAllByRole('button', { name: '削除' });
      await user.click(deleteButtons[0]);

      await waitFor(() => {
        expect(screen.getAllByPlaceholderText('例: 住居費')).toHaveLength(1);
      });
    });
  });

  describe('フォーム送信', () => {
    it('保存ボタンクリックで onSubmit が呼ばれる', async () => {
      const user = userEvent.setup();
      mockOnSubmit.mockResolvedValue(undefined);

      render(<FinancialInputForm onSubmit={mockOnSubmit} />);

      await user.click(screen.getByRole('button', { name: '保存' }));

      await waitFor(() => {
        expect(mockOnSubmit).toHaveBeenCalledWith(
          expect.objectContaining({
            monthly_income: expect.any(Number),
            monthly_expenses: expect.any(Array),
            current_savings: expect.any(Array),
          })
        );
      });
    });

    it('loading=trueのとき保存ボタンが無効になる', () => {
      render(<FinancialInputForm onSubmit={mockOnSubmit} loading={true} />);
      expect(screen.getByRole('button', { name: '保存' })).toBeDisabled();
    });
  });

  describe('initialData反映', () => {
    it('initialDataの月収が反映される', () => {
      render(
        <FinancialInputForm
          onSubmit={mockOnSubmit}
          initialData={{
            monthly_income: 500000,
            monthly_expenses: [{ category: '家賃', amount: 80000 }],
            current_savings: [{ type: 'deposit', amount: 2000000 }],
            investment_return: 5.0,
            inflation_rate: 2.0,
          }}
        />
      );
      expect(screen.getByDisplayValue('家賃')).toBeInTheDocument();
    });
  });
});
