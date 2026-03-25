'use client';

import { useState, FormEvent } from 'react';
import Link from 'next/link';
import InputField from '@/components/InputField';
import Button from '@/components/Button';
import { passwordResetAPI } from '@/lib/api-client';
import { APIError } from '@/lib/api-client';

export default function ForgotPasswordPage() {
  const [email, setEmail] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [isSubmitted, setIsSubmitted] = useState(false);
  const [localError, setLocalError] = useState('');

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setLocalError('');

    if (!email) {
      setLocalError('メールアドレスを入力してください');
      return;
    }

    if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
      setLocalError('有効なメールアドレスを入力してください');
      return;
    }

    setIsLoading(true);
    try {
      await passwordResetAPI.forgotPassword(email);
      setIsSubmitted(true);
    } catch (err) {
      if (err instanceof APIError && err.status !== 200) {
        setLocalError('送信中にエラーが発生しました。しばらく経ってから再度お試しください。');
      } else {
        // ユーザー列挙対策：エラーでも成功と同じ画面を表示
        setIsSubmitted(true);
      }
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        {/* ヘッダー */}
        <div className="text-center">
          <div className="flex justify-center mb-4">
            <span className="text-6xl">🔐</span>
          </div>
          <h2 className="text-3xl font-extrabold text-gray-900 dark:text-white">
            パスワードリセット
          </h2>
          {!isSubmitted && (
            <p className="mt-2 text-sm text-gray-600 dark:text-gray-400">
              登録したメールアドレスを入力してください
            </p>
          )}
        </div>

        {isSubmitted ? (
          /* 送信完了メッセージ */
          <div className="space-y-6">
            <div className="rounded-md bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-700 p-4">
              <div className="flex">
                <div className="flex-shrink-0">
                  <span className="text-green-500">✅</span>
                </div>
                <div className="ml-3">
                  <p className="text-sm text-green-800 dark:text-green-200">
                    パスワードリセット用のメールを送信しました。
                    <br />
                    メールをご確認いただき、記載のリンクからパスワードを再設定してください。
                  </p>
                  <p className="mt-2 text-sm text-green-700 dark:text-green-300">
                    ※ メールが届かない場合は、迷惑メールフォルダもご確認ください。
                  </p>
                </div>
              </div>
            </div>
            <div className="text-center">
              <Link href="/login">
                <Button variant="primary" fullWidth>
                  ログインページに戻る
                </Button>
              </Link>
            </div>
          </div>
        ) : (
          /* 入力フォーム */
          <form className="mt-8 space-y-6" onSubmit={handleSubmit}>
            <InputField
              label="メールアドレス"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="example@example.com"
              required
              disabled={isLoading}
              autoComplete="email"
            />

            {/* エラーメッセージ */}
            {localError && (
              <div className="rounded-md bg-error-50 dark:bg-error-900/20 p-4">
                <div className="flex">
                  <div className="flex-shrink-0">
                    <span className="text-error-500">⚠️</span>
                  </div>
                  <div className="ml-3">
                    <p className="text-sm text-error-800 dark:text-error-200">
                      {localError}
                    </p>
                  </div>
                </div>
              </div>
            )}

            <Button
              type="submit"
              fullWidth
              loading={isLoading}
              disabled={isLoading}
            >
              {isLoading ? '送信中...' : 'リセットメールを送信'}
            </Button>

            <div className="text-center">
              <Link
                href="/login"
                className="text-sm text-primary-600 hover:text-primary-500 dark:text-primary-400"
              >
                ログインページに戻る
              </Link>
            </div>
          </form>
        )}
      </div>
    </div>
  );
}
