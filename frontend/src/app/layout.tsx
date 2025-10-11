import type { Metadata } from 'next';
import './globals.css';

export const metadata: Metadata = {
  title: '財務計画計算機',
  description: '将来の資産形成と老後の財務計画を可視化するアプリケーション',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="ja">
      <body className="bg-gray-50 min-h-screen">
        <div className="min-h-screen flex flex-col">
          {children}
        </div>
      </body>
    </html>
  );
}