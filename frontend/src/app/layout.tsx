import type { Metadata } from 'next';
import './globals.css';
import Navigation from '@/components/Navigation';
import Tutorial from '@/components/Tutorial';
import { AppProviders } from '@/lib/contexts/AppProviders';

export const metadata: Metadata = {
  title: 'FinPlan - Smart Financial Planning',
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
      <body className="bg-gray-50 dark:bg-gray-900 min-h-screen transition-colors">
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