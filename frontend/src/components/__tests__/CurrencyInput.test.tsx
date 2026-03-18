import React from 'react';
import { render, screen, act } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import CurrencyInput from '../CurrencyInput';

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

describe('CurrencyInput', () => {
  const defaultProps = {
    value: 0,
    onChange: jest.fn(),
  };

  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('基本レンダリング', () => {
    it('通貨シンボル¥が表示される', () => {
      render(<CurrencyInput {...defaultProps} />);
      expect(screen.getByText('¥')).toBeInTheDocument();
    });

    it('inputMode="numeric" が設定される', () => {
      render(<CurrencyInput {...defaultProps} />);
      const input = screen.getByRole('textbox');
      expect(input).toHaveAttribute('inputMode', 'numeric');
    });

    it('プレースホルダーが表示される', () => {
      render(<CurrencyInput {...defaultProps} placeholder="1,000" />);
      expect(screen.getByPlaceholderText('1,000')).toBeInTheDocument();
    });
  });

  describe('ラベルとバリデーション表示', () => {
    it('ラベルが表示される', () => {
      render(<CurrencyInput {...defaultProps} label="月収" />);
      expect(screen.getByText('月収')).toBeInTheDocument();
    });

    it('required=true の場合、必須マークが表示される', () => {
      render(<CurrencyInput {...defaultProps} label="月収" required />);
      expect(screen.getByText('*')).toBeInTheDocument();
    });

    it('エラーメッセージが表示される', () => {
      render(<CurrencyInput {...defaultProps} error="入力が必要です" />);
      expect(screen.getByText('入力が必要です')).toBeInTheDocument();
    });

    it('ヘルパーテキストが表示される', () => {
      render(<CurrencyInput {...defaultProps} helperText="税込金額" />);
      expect(screen.getByText('税込金額')).toBeInTheDocument();
    });

    it('エラーがある場合、ヘルパーテキストは非表示', () => {
      render(
        <CurrencyInput {...defaultProps} error="エラー" helperText="税込金額" />
      );
      expect(screen.getByText('エラー')).toBeInTheDocument();
      expect(screen.queryByText('税込金額')).not.toBeInTheDocument();
    });
  });

  describe('値の表示', () => {
    it('値がフォーマットされて表示される', () => {
      render(<CurrencyInput {...defaultProps} value={1500000} />);
      const input = screen.getByRole('textbox');
      expect(input).toHaveValue('1,500,000');
    });

    it('値が0の場合、空で表示される', () => {
      render(<CurrencyInput {...defaultProps} value={0} />);
      const input = screen.getByRole('textbox');
      expect(input).toHaveValue('');
    });
  });

  describe('ユーザー入力', () => {
    it('数値を入力するとonChangeが呼ばれる', async () => {
      const user = userEvent.setup();
      const onChange = jest.fn();
      render(<CurrencyInput {...defaultProps} onChange={onChange} />);

      const input = screen.getByRole('textbox');
      await user.click(input);
      await user.type(input, '5000');

      expect(onChange).toHaveBeenCalled();
    });

    it('数値以外の入力は無視される', async () => {
      const user = userEvent.setup();
      const onChange = jest.fn();
      render(<CurrencyInput {...defaultProps} onChange={onChange} />);

      const input = screen.getByRole('textbox');
      await user.click(input);
      await user.type(input, 'abc');

      // 数値以外は onChange が呼ばれない
      const numericCalls = onChange.mock.calls.filter(
        ([val]: [number]) => val !== 0
      );
      expect(numericCalls).toHaveLength(0);
    });
  });

  describe('disabled状態', () => {
    it('disabled=true の場合、入力が無効になる', () => {
      render(<CurrencyInput {...defaultProps} disabled />);
      const input = screen.getByRole('textbox');
      expect(input).toBeDisabled();
    });
  });

  describe('onBlurコールバック', () => {
    it('フォーカスが外れるとonBlurが呼ばれる', async () => {
      const user = userEvent.setup();
      const onBlur = jest.fn();
      render(<CurrencyInput {...defaultProps} onBlur={onBlur} />);

      const input = screen.getByRole('textbox');
      await user.click(input);
      await user.tab();

      expect(onBlur).toHaveBeenCalledTimes(1);
    });
  });
});
