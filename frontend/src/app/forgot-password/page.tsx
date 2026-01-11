'use client';

import Link from 'next/link';
import Button from '@/components/Button';

export default function ForgotPasswordPage() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        <div className="text-center">
          <div className="flex justify-center mb-4">
            <span className="text-6xl">🔐</span>
          </div>
          <h2 className="text-3xl font-extrabold text-gray-900 dark:text-white">
            パスワードリセット
          </h2>
          <p className="mt-2 text-sm text-gray-600 dark:text-gray-400">
            この機能は現在開発中です
          </p>
        </div>

        <div className="mt-8 bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-700 rounded-lg p-4">
          <p className="text-sm text-yellow-800 dark:text-yellow-200">
            パスワードリセット機能は将来のリリースで実装予定です。<br />
            現在は新しいアカウントを作成するか、管理者にお問い合わせください。
          </p>
        </div>

        <div className="text-center">
          <Link href="/login">
            <Button variant="primary" fullWidth>
              ログインページに戻る
            </Button>
          </Link>
        </div>
      </div>
    </div>
  );
}
