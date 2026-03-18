import type { Metadata } from 'next';
import './globals.css';
import Navigation from '@/components/Navigation';
import Tutorial from '@/components/Tutorial';
import { AppProviders } from '@/lib/contexts/AppProviders';

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
    <html lang="ja" suppressHydrationWarning>
      <head>
        {/* Editorial fonts: Cormorant Garamond for display, Source Sans 3 for body, JetBrains Mono for numbers */}
        <link rel="preconnect" href="https://fonts.googleapis.com" />
        <link rel="preconnect" href="https://fonts.gstatic.com" crossOrigin="anonymous" />
        <link
          href="https://fonts.googleapis.com/css2?family=Cormorant+Garamond:ital,wght@0,400;0,500;0,600;0,700;1,400;1,500&family=Source+Sans+3:wght@300;400;500;600;700&family=JetBrains+Mono:wght@400;500&display=swap"
          rel="stylesheet"
        />
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
    </html>
  );
}
