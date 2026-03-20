import React from 'react';
import { render, screen } from '@testing-library/react';
import LoadingSpinner from '../LoadingSpinner';

// SVG要素では className がSVGAnimatedStringになるため getAttribute('class') を使用
function getClassString(el: Element): string {
  return el.getAttribute('class') || '';
}

describe('LoadingSpinner', () => {
  describe('基本レンダリング', () => {
    it('デフォルトのスピナーが表示される', () => {
      render(<LoadingSpinner />);
      const spinner = screen.getByRole('status');
      expect(spinner).toBeInTheDocument();
      expect(spinner).toHaveAttribute('aria-label', 'Loading');
    });

    it('テキストなしの場合、テキストが表示されない', () => {
      const { container } = render(<LoadingSpinner />);
      expect(container.querySelector('p')).not.toBeInTheDocument();
    });
  });

  describe('テキスト表示', () => {
    it('text propが指定された場合、テキストが表示される', () => {
      render(<LoadingSpinner text="読み込み中..." />);
      expect(screen.getByText('読み込み中...')).toBeInTheDocument();
    });
  });

  describe('サイズバリアント', () => {
    it('size=sm の場合、h-4 w-4 クラスが適用される', () => {
      render(<LoadingSpinner size="sm" />);
      const spinner = screen.getByRole('status');
      const cls = getClassString(spinner);
      expect(cls).toContain('h-4');
      expect(cls).toContain('w-4');
    });

    it('size=lg の場合、h-12 w-12 クラスが適用される', () => {
      render(<LoadingSpinner size="lg" />);
      const spinner = screen.getByRole('status');
      const cls = getClassString(spinner);
      expect(cls).toContain('h-12');
      expect(cls).toContain('w-12');
    });

    it('size=xl の場合、h-16 w-16 クラスが適用される', () => {
      render(<LoadingSpinner size="xl" />);
      const spinner = screen.getByRole('status');
      const cls = getClassString(spinner);
      expect(cls).toContain('h-16');
      expect(cls).toContain('w-16');
    });
  });

  describe('色バリアント', () => {
    it('color=primary の場合、text-primary-500 クラスが適用される', () => {
      render(<LoadingSpinner color="primary" />);
      const cls = getClassString(screen.getByRole('status'));
      expect(cls).toContain('text-primary-500');
    });

    it('color=white の場合、text-white クラスが適用される', () => {
      render(<LoadingSpinner color="white" />);
      const cls = getClassString(screen.getByRole('status'));
      expect(cls).toContain('text-white');
    });

    it('color=gray の場合、text-gray-500 クラスが適用される', () => {
      render(<LoadingSpinner color="gray" />);
      const cls = getClassString(screen.getByRole('status'));
      expect(cls).toContain('text-gray-500');
    });
  });

  describe('fullScreenバリアント', () => {
    it('fullScreen=true の場合、固定位置のオーバーレイで表示される', () => {
      const { container } = render(<LoadingSpinner fullScreen />);
      const overlay = container.firstChild as HTMLElement;
      expect(overlay.className).toContain('fixed');
      expect(overlay.className).toContain('inset-0');
      expect(overlay.className).toContain('z-50');
    });

    it('fullScreen=false の場合、インラインで表示される', () => {
      const { container } = render(<LoadingSpinner />);
      const wrapper = container.firstChild as HTMLElement;
      expect(wrapper.className).not.toContain('fixed');
    });
  });
});
