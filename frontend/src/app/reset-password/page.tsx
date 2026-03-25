'use client';

import { useState, FormEvent, Suspense } from 'react';
import { useSearchParams, useRouter } from 'next/navigation';
import Link from 'next/link';
import InputField from '@/components/InputField';
import Button from '@/components/Button';
import { passwordResetAPI, APIError } from '@/lib/api-client';

function ResetPasswordForm() {
  const searchParams = useSearchParams();
  const router = useRouter();
  const token = searchParams.get('token') ?? '';

  const [newPassword, setNewPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [isSuccess, setIsSuccess] = useState(false);
  const [localError, setLocalError] = useState('');

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setLocalError('');

    if (!token) {
      setLocalError('リセットトークンが見つかりません。メールのリンクから再度アクセスしてください。');
      return;
    }

    if (!newPassword || !confirmPassword) {
      setLocalError('パスワードを入力してください');
      return;
    }

    if (newPassword.length < 8) {
      setLocalError('パスワードは8文字以上で入力してください');
      return;
    }

    if (newPassword !== confirmPassword) {
      setLocalError('パスワードが一致しません');
      return;
    }

    setIsLoading(true);
    try {
      await passwordResetAPI.resetPassword(token, newPassword);
      setIsSuccess(true);
      // 3秒後にログインページへリダイレクト
      setTimeout(() => router.push('/login'), 3000);
    } catch (err) {
      if (err instanceof APIError) {
        if (err.status === 400) {
          setLocalError('リセットリンクが無効または期限切れです。パスワードリセットを再度お試しください。');
        } else {
          setLocalError('パスワードのリセットに失敗しました。しばらく経ってから再度お試しください。');
        }
      } else {
        setLocalError('ネットワークエラーが発生しました。');
      }
    } finally {
      setIsLoading(false);
    }
  };

  if (!token) {
    return (
      <div className="space-y-6">
        <div className="rounded-md bg-error-50 dark:bg-error-900/20 border border-red-200 dark:border-red-700 p-4">
          <div className="flex">
            <div className="flex-shrink-0">
              <span className="text-error-500">⚠️</span>
            </div>
            <div className="ml-3">
              <p className="text-sm text-error-800 dark:text-error-200">
                無効なリンクです。パスワードリセットメールのリンクから再度アクセスしてください。
              </p>
            </div>
          </div>
        </div>
        <div className="text-center">
          <Link href="/forgot-password">
            <Button variant="primary" fullWidth>
              パスワードリセットを再申請する
            </Button>
          </Link>
        </div>
      </div>
    );
  }

  if (isSuccess) {
    return (
      <div className="space-y-6">
        <div className="rounded-md bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-700 p-4">
          <div className="flex">
            <div className="flex-shrink-0">
              <span className="text-green-500">✅</span>
            </div>
            <div className="ml-3">
              <p className="text-sm text-green-800 dark:text-green-200">
                パスワードをリセットしました。
                <br />
                3秒後にログインページに移動します。
              </p>
            </div>
          </div>
        </div>
        <div className="text-center">
          <Link href="/login">
            <Button variant="primary" fullWidth>
              ログインページへ
            </Button>
          </Link>
        </div>
      </div>
    );
  }

  return (
    <form className="mt-8 space-y-6" onSubmit={handleSubmit}>
      <div className="space-y-4">
        <InputField
          label="新しいパスワード"
          type="password"
          value={newPassword}
          onChange={(e) => setNewPassword(e.target.value)}
          placeholder="••••••••"
          required
          disabled={isLoading}
          autoComplete="new-password"
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
        />
      </div>

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
        {isLoading ? 'リセット中...' : 'パスワードをリセット'}
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
  );
}

export default function ResetPasswordPage() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        {/* ヘッダー */}
        <div className="text-center">
          <div className="flex justify-center mb-4">
            <span className="text-6xl">🔑</span>
          </div>
          <h2 className="text-3xl font-extrabold text-gray-900 dark:text-white">
            新しいパスワードを設定
          </h2>
          <p className="mt-2 text-sm text-gray-600 dark:text-gray-400">
            8文字以上のパスワードを入力してください
          </p>
        </div>

        {/* useSearchParams を使うコンポーネントは Suspense でラップ必須 */}
        <Suspense fallback={<div className="text-center text-gray-500">読み込み中...</div>}>
          <ResetPasswordForm />
        </Suspense>
      </div>
    </div>
  );
}
