import React from 'react';
import { render, screen } from '@testing-library/react';

// next/font/google モック（CSS変数を返す）
// jest.mock はホイストされるため、ファクトリ内で jest.fn を定義する
jest.mock('next/font/google', () => ({
  Cormorant_Garamond: jest.fn((_: unknown) => ({
    variable: '--font-display',
    className: 'cormorant-garamond',
  })),
  Source_Sans_3: jest.fn((_: unknown) => ({
    variable: '--font-body',
    className: 'source-sans-3',
  })),
  JetBrains_Mono: jest.fn((_: unknown) => ({
    variable: '--font-mono',
    className: 'jetbrains-mono',
  })),
}));

// 子コンポーネントをモック
jest.mock('@/components/Navigation', () => ({
  __esModule: true,
  default: () => <div data-testid="navigation" />,
}));

jest.mock('@/components/Tutorial', () => ({
  __esModule: true,
  default: () => <div data-testid="tutorial" />,
}));

jest.mock('@/lib/contexts/AppProviders', () => ({
  AppProviders: ({ children }: { children: React.ReactNode }) => (
    <div data-testid="app-providers">{children}</div>
  ),
}));

// next/navigation モック
jest.mock('next/navigation', () => ({
  usePathname: () => '/',
  useRouter: () => ({ push: jest.fn() }),
}));

import RootLayout from '../layout';
import { Cormorant_Garamond, Source_Sans_3, JetBrains_Mono } from 'next/font/google';

const mockCormorantGaramond = jest.mocked(Cormorant_Garamond);
const mockSourceSans3 = jest.mocked(Source_Sans_3);
const mockJetBrainsMono = jest.mocked(JetBrains_Mono);

describe('RootLayout', () => {
  describe('next/font/google の設定', () => {
    // フォント関数はモジュールレベルで一度だけ呼ばれる
    // jest.clearAllMocks()（beforeEach で実行）でリセットされる前に引数をキャプチャする
    let cormorantArgs: unknown;
    let sourceSans3Args: unknown;
    let jetbrainsMonoArgs: unknown;

    beforeAll(() => {
      cormorantArgs = mockCormorantGaramond.mock.calls[0]?.[0];
      sourceSans3Args = mockSourceSans3.mock.calls[0]?.[0];
      jetbrainsMonoArgs = mockJetBrainsMono.mock.calls[0]?.[0];
    });

    it('Cormorant Garamond を正しいオプションで初期化する', () => {
      expect(cormorantArgs).toMatchObject({
        subsets: ['latin'],
        weight: expect.arrayContaining(['400', '500', '600', '700']),
        style: expect.arrayContaining(['normal', 'italic']),
        variable: '--font-display',
        display: 'swap',
      });
    });

    it('Source Sans 3 を正しいオプションで初期化する', () => {
      expect(sourceSans3Args).toMatchObject({
        subsets: ['latin'],
        weight: expect.arrayContaining(['300', '400', '500', '600', '700']),
        variable: '--font-body',
        display: 'swap',
      });
    });

    it('JetBrains Mono を正しいオプションで初期化する', () => {
      expect(jetbrainsMonoArgs).toMatchObject({
        subsets: ['latin'],
        weight: expect.arrayContaining(['400', '500']),
        variable: '--font-mono',
        display: 'swap',
      });
    });
  });

  describe('Google Fonts の <link> タグ', () => {
    it('fonts.googleapis.com への <link> タグが存在しない', () => {
      const { container } = render(<RootLayout><div data-testid="child" /></RootLayout>);

      const googleFontLinks = container.querySelectorAll(
        'link[href*="fonts.googleapis.com"]'
      );
      expect(googleFontLinks).toHaveLength(0);
    });

    it('fonts.gstatic.com への preconnect <link> タグが存在しない', () => {
      const { container } = render(<RootLayout><div /></RootLayout>);

      const gstaticLinks = container.querySelectorAll(
        'link[href*="fonts.gstatic.com"]'
      );
      expect(gstaticLinks).toHaveLength(0);
    });
  });

  describe('子コンテンツのレンダリング', () => {
    it('children が描画される', () => {
      render(
        <RootLayout>
          <div data-testid="page-content">ページコンテンツ</div>
        </RootLayout>
      );

      expect(screen.getByTestId('page-content')).toBeInTheDocument();
      expect(screen.getByText('ページコンテンツ')).toBeInTheDocument();
    });

    it('Navigation が描画される', () => {
      render(<RootLayout><div /></RootLayout>);

      expect(screen.getByTestId('navigation')).toBeInTheDocument();
    });

    it('Tutorial が描画される', () => {
      render(<RootLayout><div /></RootLayout>);

      expect(screen.getByTestId('tutorial')).toBeInTheDocument();
    });
  });

  describe('CSS変数クラスの適用', () => {
    it('3つのフォントCSS変数クラスがレンダリング結果に含まれる', () => {
      const { container } = render(<RootLayout><div /></RootLayout>);

      // モックが返す variable 値（'--font-display' など）が className として使われる
      // jsdom の挙動に関わらず innerHTML で検証できる
      const html = container.innerHTML + document.documentElement.outerHTML;
      expect(html).toContain('--font-display');
      expect(html).toContain('--font-body');
      expect(html).toContain('--font-mono');
    });
  });
});
