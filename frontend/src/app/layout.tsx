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