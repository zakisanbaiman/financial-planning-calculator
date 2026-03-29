import type { Metadata } from 'next';
import Link from 'next/link';

export const metadata: Metadata = {
  title: 'ページが見つかりません — FinPlan',
};

export default function NotFound() {
  return (
    <div className="flex flex-col items-center justify-center min-h-[60vh] px-4 text-center">
      <p className="font-mono text-ink-400 dark:text-ink-500 text-sm tracking-widest mb-4">404</p>
      <h1 className="font-display text-3xl text-ink-900 dark:text-ink-100 mb-3">
        ページが見つかりません
      </h1>
      <p className="text-ink-600 dark:text-ink-400 mb-8 max-w-md">
        お探しのページは存在しないか、移動・削除された可能性があります。
      </p>
      <Link
        href="/"
        className="px-6 py-2.5 bg-accent-600 hover:bg-accent-700 text-white rounded-lg transition-colors text-sm font-medium"
      >
        ホームへ戻る
      </Link>
    </div>
  );
}
