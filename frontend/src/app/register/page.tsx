'use client';

import { useState, FormEvent } from 'react';
import { useAuth } from '@/lib/contexts/AuthContext';
import InputField from '@/components/InputField';
import Button from '@/components/Button';
import Link from 'next/link';

export default function RegisterPage() {
  const { register, error, clearError, isLoading } = useAuth();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [localError, setLocalError] = useState('');

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setLocalError('');
    clearError();

    // バリデーション
    if (!email || !password || !confirmPassword) {
      setLocalError('すべてのフィールドを入力してください');
      return;
    }

    if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
      setLocalError('有効なメールアドレスを入力してください');
      return;
    }

    if (password.length < 8) {
      setLocalError('パスワードは8文字以上である必要があります');
      return;
    }

    if (password !== confirmPassword) {
      setLocalError('パスワードが一致しません');
      return;
    }

    try {
      await register(email, password);
      // 成功時はAuthContextでリダイレクト処理される
    } catch (err) {
      // エラーはAuthContextで管理される
      console.error('Registration failed:', err);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        {/* ヘッダー */}
        <div className="text-center">
          <div className="flex justify-center mb-4">
            <span className="text-6xl">💼</span>
          </div>
          <h2 className="text-3xl font-extrabold text-gray-900 dark:text-white">
            財務計画計算機
          </h2>
          <p className="mt-2 text-sm text-gray-600 dark:text-gray-400">
            新規アカウント登録
          </p>
        </div>

        {/* フォーム */}
        <form className="mt-8 space-y-6" onSubmit={handleSubmit}>
          <div className="rounded-md shadow-sm space-y-4">
            <InputField
              label="メールアドレス"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="example@example.com"
              required
              disabled={isLoading}
              autoComplete="email"
              helperText="ログイン時に使用するメールアドレス"
            />

            <InputField
              label="パスワード"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="••••••••"
              required
              disabled={isLoading}
              autoComplete="new-password"
              helperText="8文字以上で入力してください"
            />

            <InputField
              label="パスワード（確認）"
              type="password"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              placeholder="••••••••"
              required
              disabled={isLoading}
              autoComplete="new-password"
              helperText="確認のため再度入力してください"
            />
          </div>

          {/* エラーメッセージ */}
          {(error || localError) && (
            <div className="rounded-md bg-error-50 dark:bg-error-900/20 p-4">
              <div className="flex">
                <div className="flex-shrink-0">
                  <span className="text-error-500">⚠️</span>
                </div>
                <div className="ml-3">
                  <p className="text-sm text-error-800 dark:text-error-200">
                    {error || localError}
                  </p>
                </div>
              </div>
            </div>
          )}

          {/* 登録ボタン */}
          <Button
            type="submit"
            fullWidth
            loading={isLoading}
            disabled={isLoading}
          >
            {isLoading ? '登録中...' : 'アカウントを作成'}
          </Button>

          {/* ログインリンク */}
          <div className="text-center">
            <p className="text-sm text-gray-600 dark:text-gray-400">
              既にアカウントをお持ちですか？{' '}
              <Link
                href="/login"
                className="font-medium text-primary-600 hover:text-primary-500 dark:text-primary-400"
              >
                ログイン
              </Link>
            </p>
          </div>

          {/* 利用規約（将来実装） */}
          <div className="text-center">
            <p className="text-xs text-gray-500 dark:text-gray-400">
              登録することで、
              <Link href="/terms" className="underline hover:text-gray-700 dark:hover:text-gray-300">
                利用規約
              </Link>
              と
              <Link href="/privacy" className="underline hover:text-gray-700 dark:hover:text-gray-300">
                プライバシーポリシー
              </Link>
              に同意したものとみなされます
            </p>
          </div>
        </form>
      </div>
    </div>
  );
}
