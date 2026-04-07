import type { Metadata } from 'next';
import { Cormorant_Garamond, Source_Sans_3, JetBrains_Mono } from 'next/font/google';
import { GoogleAnalytics } from '@next/third-parties/google';
import './globals.css';
import Navigation from '@/components/Navigation';
import Footer from '@/components/Footer';
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
  title: 'FinPlan — Financial Planning',
  description: 'Visualize your financial future and plan for retirement with confidence',
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
            <Footer />
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
