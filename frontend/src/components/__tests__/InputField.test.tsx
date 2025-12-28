import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import InputField from '../InputField';

describe('InputField', () => {
  describe('基本表示', () => {
    it('ラベルが正しく表示される', () => {
      render(<InputField label="テストラベル" />);
      expect(screen.getByText('テストラベル')).toBeInTheDocument();
    });

    it('ラベルなしでもレンダリングできる', () => {
      render(<InputField placeholder="テスト入力" />);
      expect(screen.getByPlaceholderText('テスト入力')).toBeInTheDocument();
    });

    it('必須フィールドに * マークが表示される', () => {
      render(<InputField label="必須項目" required />);
      expect(screen.getByText('*')).toBeInTheDocument();
    });

    it('プレースホルダーが正しく表示される', () => {
      render(<InputField placeholder="ここに入力" />);
      expect(screen.getByPlaceholderText('ここに入力')).toBeInTheDocument();
    });
  });

  describe('エラー表示', () => {
    it('エラーメッセージが表示される', () => {
      render(<InputField label="入力" error="入力エラーです" />);
      expect(screen.getByText('入力エラーです')).toBeInTheDocument();
    });

    it('エラー時にエラースタイルが適用される', () => {
      render(<InputField label="入力" error="エラー" data-testid="input" />);
      const input = screen.getByTestId('input');
      expect(input.className).toContain('border-error');
    });

    it('エラーがない時はヘルパーテキストが表示される', () => {
      render(<InputField label="入力" helperText="補足説明" />);
      expect(screen.getByText('補足説明')).toBeInTheDocument();
    });

    it('エラーがある時はヘルパーテキストが表示されない', () => {
      render(<InputField label="入力" error="エラー" helperText="補足説明" />);
      expect(screen.queryByText('補足説明')).not.toBeInTheDocument();
      expect(screen.getByText('エラー')).toBeInTheDocument();
    });
  });

  describe('入力操作', () => {
    it('テキスト入力が正しく動作する', async () => {
      const user = userEvent.setup();
      render(<InputField label="名前" data-testid="name-input" />);
      
      const input = screen.getByTestId('name-input');
      await user.type(input, 'テスト太郎');
      
      expect(input).toHaveValue('テスト太郎');
    });

    it('数値入力が正しく動作する', async () => {
      const user = userEvent.setup();
      render(<InputField type="number" label="金額" data-testid="amount-input" />);
      
      const input = screen.getByTestId('amount-input');
      await user.type(input, '10000');
      
      expect(input).toHaveValue(10000);
    });

    it('onChange ハンドラが呼ばれる', async () => {
      const handleChange = jest.fn();
      const user = userEvent.setup();
      
      render(
        <InputField 
          label="入力" 
          data-testid="input" 
          onChange={handleChange} 
        />
      );
      
      const input = screen.getByTestId('input');
      await user.type(input, 'a');
      
      expect(handleChange).toHaveBeenCalled();
    });
  });

  describe('無効化状態', () => {
    it('disabled 属性が正しく適用される', () => {
      render(<InputField label="入力" disabled data-testid="input" />);
      expect(screen.getByTestId('input')).toBeDisabled();
    });

    it('disabled 時に入力できない', async () => {
      const user = userEvent.setup();
      render(<InputField label="入力" disabled data-testid="input" />);
      
      const input = screen.getByTestId('input');
      await user.type(input, 'テスト');
      
      expect(input).toHaveValue('');
    });
  });

  describe('スタイル', () => {
    it('fullWidth がデフォルトで適用される', () => {
      render(<InputField label="入力" data-testid="input" />);
      const input = screen.getByTestId('input');
      expect(input.className).toContain('w-full');
    });

    it('fullWidth=false でフル幅にならない', () => {
      render(<InputField label="入力" fullWidth={false} data-testid="input" />);
      const input = screen.getByTestId('input');
      expect(input.className).not.toContain('w-full');
    });

    it('カスタムクラス名が適用できる', () => {
      render(<InputField label="入力" className="custom-class" data-testid="input" />);
      const input = screen.getByTestId('input');
      expect(input.className).toContain('custom-class');
    });
  });

  describe('ref フォワーディング', () => {
    it('ref が正しくフォワードされる', () => {
      const ref = React.createRef<HTMLInputElement>();
      render(<InputField label="入力" ref={ref} />);
      
      expect(ref.current).toBeInstanceOf(HTMLInputElement);
    });

    it('ref を使ってフォーカスできる', () => {
      const ref = React.createRef<HTMLInputElement>();
      render(<InputField label="入力" ref={ref} data-testid="input" />);
      
      ref.current?.focus();
      
      expect(screen.getByTestId('input')).toHaveFocus();
    });
  });
});
