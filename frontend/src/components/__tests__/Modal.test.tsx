import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import Modal from '../Modal';

describe('Modal', () => {
  const defaultProps = {
    isOpen: true,
    onClose: jest.fn(),
    children: <p>モーダルコンテンツ</p>,
  };

  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('開閉制御', () => {
    it('isOpen=true の場合、モーダルが表示される', () => {
      render(<Modal {...defaultProps} />);
      expect(screen.getByText('モーダルコンテンツ')).toBeInTheDocument();
    });

    it('isOpen=false の場合、モーダルが表示されない', () => {
      render(<Modal {...defaultProps} isOpen={false} />);
      expect(screen.queryByText('モーダルコンテンツ')).not.toBeInTheDocument();
    });
  });

  describe('タイトル', () => {
    it('title propが指定された場合、タイトルが表示される', () => {
      render(<Modal {...defaultProps} title="テストモーダル" />);
      expect(screen.getByText('テストモーダル')).toBeInTheDocument();
    });

    it('タイトルにaria-labelledbyが設定される', () => {
      render(<Modal {...defaultProps} title="テストモーダル" />);
      const dialog = screen.getByRole('dialog');
      expect(dialog).toHaveAttribute('aria-labelledby', 'modal-title');
    });
  });

  describe('閉じるボタン', () => {
    it('デフォルトで閉じるボタンが表示される', () => {
      render(<Modal {...defaultProps} title="テスト" />);
      expect(screen.getByLabelText('Close modal')).toBeInTheDocument();
    });

    it('閉じるボタンをクリックするとonCloseが呼ばれる', async () => {
      const user = userEvent.setup();
      render(<Modal {...defaultProps} title="テスト" />);
      await user.click(screen.getByLabelText('Close modal'));
      expect(defaultProps.onClose).toHaveBeenCalledTimes(1);
    });

    it('showCloseButton=false の場合、閉じるボタンが非表示', () => {
      render(<Modal {...defaultProps} title="テスト" showCloseButton={false} />);
      expect(screen.queryByLabelText('Close modal')).not.toBeInTheDocument();
    });
  });

  describe('Escapeキー', () => {
    it('Escapeキーを押すとonCloseが呼ばれる', async () => {
      const user = userEvent.setup();
      render(<Modal {...defaultProps} />);
      await user.keyboard('{Escape}');
      expect(defaultProps.onClose).toHaveBeenCalledTimes(1);
    });
  });

  describe('バックドロップクリック', () => {
    it('バックドロップをクリックするとonCloseが呼ばれる', async () => {
      const user = userEvent.setup();
      render(<Modal {...defaultProps} />);
      // バックドロップは最外側のdiv
      const backdrop = screen.getByRole('dialog').parentElement!;
      await user.click(backdrop);
      expect(defaultProps.onClose).toHaveBeenCalledTimes(1);
    });

    it('モーダルコンテンツをクリックしてもonCloseは呼ばれない', async () => {
      const user = userEvent.setup();
      render(<Modal {...defaultProps} />);
      await user.click(screen.getByText('モーダルコンテンツ'));
      expect(defaultProps.onClose).not.toHaveBeenCalled();
    });
  });

  describe('サイズバリアント', () => {
    it.each([
      ['sm', 'max-w-md'],
      ['md', 'max-w-lg'],
      ['lg', 'max-w-2xl'],
      ['xl', 'max-w-4xl'],
    ] as const)('size=%s の場合、%s クラスが適用される', (size, expectedClass) => {
      render(<Modal {...defaultProps} size={size} />);
      const dialog = screen.getByRole('dialog');
      expect(dialog.className).toContain(expectedClass);
    });
  });

  describe('body overflow制御', () => {
    it('モーダルが開くとbody overflowがhiddenになる', () => {
      render(<Modal {...defaultProps} />);
      expect(document.body.style.overflow).toBe('hidden');
    });

    it('モーダルが閉じるとbody overflowがunsetに戻る', () => {
      const { unmount } = render(<Modal {...defaultProps} />);
      unmount();
      expect(document.body.style.overflow).toBe('unset');
    });
  });

  describe('アクセシビリティ', () => {
    it('aria-modal="true" が設定される', () => {
      render(<Modal {...defaultProps} />);
      const dialog = screen.getByRole('dialog');
      expect(dialog).toHaveAttribute('aria-modal', 'true');
    });
  });
});
