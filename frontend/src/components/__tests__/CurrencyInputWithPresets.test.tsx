import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import CurrencyInputWithPresets from '../CurrencyInputWithPresets';

describe('CurrencyInputWithPresets', () => {
  const presets = [
    { label: '10万', value: 100000 },
    { label: '30万', value: 300000 },
    { label: '50万', value: 500000 },
  ];

  const defaultProps = {
    value: 0,
    onChange: jest.fn(),
    presets,
  };

  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('基本レンダリング', () => {
    it('通貨シンボル¥が表示される', () => {
      render(<CurrencyInputWithPresets {...defaultProps} />);
      expect(screen.getByText('¥')).toBeInTheDocument();
    });

    it('単位が表示される（デフォルト: 円）', () => {
      render(<CurrencyInputWithPresets {...defaultProps} />);
      expect(screen.getByText('円')).toBeInTheDocument();
    });

    it('カスタム単位が表示される', () => {
      render(<CurrencyInputWithPresets {...defaultProps} unit="万円" />);
      expect(screen.getByText('万円')).toBeInTheDocument();
    });

    it('ラベルが表示される', () => {
      render(<CurrencyInputWithPresets {...defaultProps} label="月収" />);
      expect(screen.getByText('月収')).toBeInTheDocument();
    });

    it('required=true の場合、必須マークが表示される', () => {
      render(<CurrencyInputWithPresets {...defaultProps} label="月収" required />);
      expect(screen.getByText('*')).toBeInTheDocument();
    });
  });

  describe('プリセットボタン', () => {
    it('プリセットボタンが表示される', () => {
      render(<CurrencyInputWithPresets {...defaultProps} />);
      expect(screen.getByText('10万')).toBeInTheDocument();
      expect(screen.getByText('30万')).toBeInTheDocument();
      expect(screen.getByText('50万')).toBeInTheDocument();
    });

    it('プリセットボタンをクリックするとonChangeが呼ばれる', async () => {
      const user = userEvent.setup();
      const onChange = jest.fn();
      render(<CurrencyInputWithPresets {...defaultProps} onChange={onChange} />);

      await user.click(screen.getByText('30万'));
      expect(onChange).toHaveBeenCalledWith(300000);
    });

    it('選択中のプリセットがハイライトされる', () => {
      render(<CurrencyInputWithPresets {...defaultProps} value={300000} />);
      const activeButton = screen.getByText('30万');
      expect(activeButton.className).toContain('bg-primary-500');
      expect(activeButton.className).toContain('text-white');
    });

    it('未選択のプリセットはデフォルトスタイル', () => {
      render(<CurrencyInputWithPresets {...defaultProps} value={300000} />);
      const inactiveButton = screen.getByText('10万');
      expect(inactiveButton.className).toContain('bg-gray-100');
    });

    it('プリセットがない場合、ボタンエリアが表示されない', () => {
      render(
        <CurrencyInputWithPresets {...defaultProps} presets={[]} />
      );
      expect(screen.queryByText('10万')).not.toBeInTheDocument();
    });
  });

  describe('disabled状態', () => {
    it('disabled=true の場合、入力が無効になる', () => {
      render(<CurrencyInputWithPresets {...defaultProps} disabled />);
      const input = screen.getByRole('textbox');
      expect(input).toBeDisabled();
    });

    it('disabled=true の場合、プリセットボタンも無効になる', () => {
      render(<CurrencyInputWithPresets {...defaultProps} disabled />);
      const buttons = screen.getAllByRole('button');
      buttons.forEach((button) => {
        expect(button).toBeDisabled();
      });
    });

    it('disabled=true の場合、プリセットクリックでonChangeが呼ばれない', async () => {
      const user = userEvent.setup();
      const onChange = jest.fn();
      render(
        <CurrencyInputWithPresets {...defaultProps} onChange={onChange} disabled />
      );

      await user.click(screen.getByText('30万'));
      expect(onChange).not.toHaveBeenCalled();
    });
  });

  describe('エラー・ヘルパーテキスト', () => {
    it('エラーメッセージが表示される', () => {
      render(
        <CurrencyInputWithPresets {...defaultProps} error="入力が必要です" />
      );
      expect(screen.getByText('入力が必要です')).toBeInTheDocument();
    });

    it('ヘルパーテキストが表示される', () => {
      render(
        <CurrencyInputWithPresets {...defaultProps} helperText="推奨: 30万円" />
      );
      expect(screen.getByText('推奨: 30万円')).toBeInTheDocument();
    });

    it('エラーがある場合、ヘルパーテキストは非表示', () => {
      render(
        <CurrencyInputWithPresets
          {...defaultProps}
          error="エラー"
          helperText="推奨: 30万円"
        />
      );
      expect(screen.getByText('エラー')).toBeInTheDocument();
      expect(screen.queryByText('推奨: 30万円')).not.toBeInTheDocument();
    });
  });

  describe('値の表示', () => {
    it('値がフォーマットされて表示される', () => {
      render(<CurrencyInputWithPresets {...defaultProps} value={1500000} />);
      const input = screen.getByRole('textbox');
      expect(input).toHaveValue('1,500,000');
    });

    it('値が0の場合、空で表示される', () => {
      render(<CurrencyInputWithPresets {...defaultProps} value={0} />);
      const input = screen.getByRole('textbox');
      expect(input).toHaveValue('');
    });
  });
});
