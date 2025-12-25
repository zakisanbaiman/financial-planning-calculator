import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import Button from '../Button';

describe('Button', () => {
  describe('基本表示', () => {
    it('子要素（テキスト）が正しく表示される', () => {
      render(<Button>クリック</Button>);
      expect(screen.getByRole('button', { name: 'クリック' })).toBeInTheDocument();
    });

    it('子要素（React ノード）が正しく表示される', () => {
      render(
        <Button>
          <span data-testid="icon">🔥</span>
          送信
        </Button>
      );
      expect(screen.getByTestId('icon')).toBeInTheDocument();
      expect(screen.getByRole('button')).toHaveTextContent('送信');
    });
  });

  describe('バリアント', () => {
    it('primary バリアントがデフォルトで適用される', () => {
      render(<Button>ボタン</Button>);
      const button = screen.getByRole('button');
      expect(button.className).toContain('bg-primary');
    });

    it('secondary バリアントが適用される', () => {
      render(<Button variant="secondary">ボタン</Button>);
      const button = screen.getByRole('button');
      expect(button.className).toContain('bg-gray');
    });

    it('success バリアントが適用される', () => {
      render(<Button variant="success">ボタン</Button>);
      const button = screen.getByRole('button');
      expect(button.className).toContain('bg-success');
    });

    it('warning バリアントが適用される', () => {
      render(<Button variant="warning">ボタン</Button>);
      const button = screen.getByRole('button');
      expect(button.className).toContain('bg-warning');
    });

    it('error バリアントが適用される', () => {
      render(<Button variant="error">ボタン</Button>);
      const button = screen.getByRole('button');
      expect(button.className).toContain('bg-error');
    });

    it('outline バリアントが適用される', () => {
      render(<Button variant="outline">ボタン</Button>);
      const button = screen.getByRole('button');
      expect(button.className).toContain('border-primary');
    });
  });

  describe('サイズ', () => {
    it('md サイズがデフォルトで適用される', () => {
      render(<Button>ボタン</Button>);
      const button = screen.getByRole('button');
      expect(button.className).toContain('py-2');
    });

    it('sm サイズが適用される', () => {
      render(<Button size="sm">ボタン</Button>);
      const button = screen.getByRole('button');
      expect(button.className).toContain('py-1.5');
    });

    it('lg サイズが適用される', () => {
      render(<Button size="lg">ボタン</Button>);
      const button = screen.getByRole('button');
      expect(button.className).toContain('py-3');
    });
  });

  describe('フル幅', () => {
    it('fullWidth でフル幅になる', () => {
      render(<Button fullWidth>ボタン</Button>);
      const button = screen.getByRole('button');
      expect(button.className).toContain('w-full');
    });

    it('fullWidth=false ではフル幅にならない', () => {
      render(<Button fullWidth={false}>ボタン</Button>);
      const button = screen.getByRole('button');
      expect(button.className).not.toContain('w-full');
    });
  });

  describe('クリックイベント', () => {
    it('onClick ハンドラが呼ばれる', async () => {
      const handleClick = jest.fn();
      const user = userEvent.setup();
      
      render(<Button onClick={handleClick}>クリック</Button>);
      
      await user.click(screen.getByRole('button'));
      
      expect(handleClick).toHaveBeenCalledTimes(1);
    });

    it('disabled 時は onClick が呼ばれない', async () => {
      const handleClick = jest.fn();
      const user = userEvent.setup();
      
      render(<Button onClick={handleClick} disabled>クリック</Button>);
      
      await user.click(screen.getByRole('button'));
      
      expect(handleClick).not.toHaveBeenCalled();
    });
  });

  describe('無効化状態', () => {
    it('disabled 属性が正しく適用される', () => {
      render(<Button disabled>ボタン</Button>);
      expect(screen.getByRole('button')).toBeDisabled();
    });

    it('disabled 時にスタイルが変更される', () => {
      render(<Button disabled>ボタン</Button>);
      const button = screen.getByRole('button');
      expect(button.className).toContain('disabled:opacity-50');
    });
  });

  describe('ローディング状態', () => {
    it('loading 時にスピナーが表示される', () => {
      render(<Button loading>送信中</Button>);
      const button = screen.getByRole('button');
      // SVG スピナーが存在することを確認
      expect(button.querySelector('svg')).toBeInTheDocument();
    });

    it('loading 時はボタンが無効化される', () => {
      render(<Button loading>送信中</Button>);
      expect(screen.getByRole('button')).toBeDisabled();
    });

    it('loading 時に onClick が呼ばれない', async () => {
      const handleClick = jest.fn();
      const user = userEvent.setup();
      
      render(<Button loading onClick={handleClick}>送信中</Button>);
      
      await user.click(screen.getByRole('button'));
      
      expect(handleClick).not.toHaveBeenCalled();
    });

    it('loading=false ではスピナーが表示されない', () => {
      render(<Button loading={false}>送信</Button>);
      const button = screen.getByRole('button');
      expect(button.querySelector('svg')).not.toBeInTheDocument();
    });
  });

  describe('カスタム属性', () => {
    it('type 属性が正しく適用される', () => {
      render(<Button type="submit">送信</Button>);
      expect(screen.getByRole('button')).toHaveAttribute('type', 'submit');
    });

    it('カスタムクラス名が適用される', () => {
      render(<Button className="custom-class">ボタン</Button>);
      const button = screen.getByRole('button');
      expect(button.className).toContain('custom-class');
    });

    it('data-testid が適用される', () => {
      render(<Button data-testid="submit-btn">送信</Button>);
      expect(screen.getByTestId('submit-btn')).toBeInTheDocument();
    });
  });
});
