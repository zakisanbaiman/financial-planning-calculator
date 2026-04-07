import type { Metadata } from 'next';
import { Cormorant_Garamond, Source_Sans_3, JetBrains_Mono } from 'next/font/google';
import { GoogleAnalytics } from '@next/third-parties/google';
import './globals.css';
import Navigation from '@/components/Navigation';
import Tutorial from '@/components/Tutorial';
import { AppProviders } from '@/lib/contexts/AppProviders';

const cormorantGaramond = Cormorant_Garamond({
  subsets: ['latin'],
  weight: ['400', '500', '600', '700'],
  style: ['normal', 'italic'],
  variable: '--font-display',
  display: 'swap',
});

const sourceSans3 = Source_Sans_3({
  subsets: ['latin'],
  weight: ['300', '400', '500', '600', '700'],
  variable: '--font-body',
  display: 'swap',
});

const jetbrainsMono = JetBrains_Mono({
  subsets: ['latin'],
  weight: ['400', '500'],
  variable: '--font-mono',
  display: 'swap',
});

export const metadata: Metadata = {
  metadataBase: new URL('https://financial-planning-frontend-production.up.railway.app'),
  title: { default: 'FinPlan - Smart Financial Planning', template: '%s | FinPlan' },
  description: '将来の資産推移を可視化し、老後資金・緊急資金を計画するツール',
  openGraph: {
    title: 'FinPlan - Smart Financial Planning',
    description: '将来の資産推移を可視化し、老後資金・緊急資金を計画するツール',
    url: 'https://financial-planning-frontend-production.up.railway.app',
    siteName: 'FinPlan',
    images: [{ url: '/api/og', width: 1200, height: 630 }],
    locale: 'ja_JP',
    type: 'website',
  },
  twitter: {
    card: 'summary_large_image',
    title: 'FinPlan - Smart Financial Planning',
    description: '将来の資産推移を可視化し、老後資金・緊急資金を計画するツール',
    images: ['/api/og'],
  },
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html
      lang="ja"
      suppressHydrationWarning
      className={`${cormorantGaramond.variable} ${sourceSans3.variable} ${jetbrainsMono.variable}`}
    >
      <head>
        {/* テーマの初期値を素早く読み込んで、白飛びを防止する */}
        <script
          dangerouslySetInnerHTML={{
            __html: `
              (function() {
                try {
                  const storedTheme = localStorage.getItem('theme');
                  const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
                  const theme = storedTheme || (prefersDark ? 'dark' : 'light');

                  if (theme === 'dark') {
                    document.documentElement.classList.add('dark');
                  } else {
                    document.documentElement.classList.remove('dark');
                  }
                } catch (error) {
                  console.warn('Failed to load theme:', error);
                }
              })();
            `,
          }}
        />
      </head>
      <body className="bg-ink-50 dark:bg-ink-950 min-h-screen transition-colors">
        <AppProviders>
          <div className="min-h-screen flex flex-col">
            <Navigation />
            <main className="flex-1">
              {children}
            </main>
          </div>
          <Tutorial />
        </AppProviders>
      </body>
      {process.env.NEXT_PUBLIC_GA_ID && (
        <GoogleAnalytics gaId={process.env.NEXT_PUBLIC_GA_ID} />
      )}
    </html>
  );
}
