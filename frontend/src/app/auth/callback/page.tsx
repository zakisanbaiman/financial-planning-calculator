'use client';

import { Suspense, useEffect, useRef } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { useAuth } from '@/lib/contexts/AuthContext';

/**
 * OAuthコールバック処理コンポーネント（Issue: #67）
 * useSearchParams()を使用するため、Suspenseでラップ
 * トークンはhttpOnly Cookieで管理されるため、ユーザー情報のみを処理
 */
function CallbackHandler() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { setAuthData, isAuthenticated, isLoading } = useAuth();
  const hasProcessed = useRef(false);

  // トークン情報をAuthContextに保存
  useEffect(() => {
    if (hasProcessed.current) return;

    const handleCallback = async () => {
      // URLパラメータからユーザー情報を取得（トークンはCookieで管理）
      const userId = searchParams.get('user_id');
      const email = searchParams.get('email');
      const error = searchParams.get('error');

      // エラーがある場合
      if (error) {
        console.error('OAuth authentication failed:', error);
        router.push(`/login?error=${error}`);
        return;
      }

      // ユーザー情報が存在する場合
      if (userId && email) {
        // URLからユーザー情報をクリア（ブラウザ履歴に残さない）
        window.history.replaceState({}, '', '/auth/callback');

        // AuthContextにユーザー情報を保存（トークンはCookieで管理されている）
        setAuthData({
          user: {
            userId,
            email,
          },
        });

        hasProcessed.current = true;
      } else {
        // ユーザー情報が不足している場合
        console.error('Missing user data in callback');
        router.push('/login?error=missing_data');
      }
    };

    handleCallback();
  }, [searchParams, router, setAuthData]);

  // 認証状態が更新されたらダッシュボードにリダイレクト
  useEffect(() => {
    if (hasProcessed.current && isAuthenticated && !isLoading) {
      router.push('/dashboard');
    }
  }, [isAuthenticated, isLoading, router]);

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900">
      <div className="text-center">
        <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
        <p className="mt-4 text-gray-600 dark:text-gray-400">
          認証中...
        </p>
      </div>
    </div>
  );
}

/**
 * OAuthコールバックページ（Issue: #67）
 * GitHub OAuth認証後のリダイレクト先
 * URLパラメータからユーザー情報を取得してAuthContextに保存
 * トークンはhttpOnly Cookieで管理されるため、URLには含まれない
 */
export default function AuthCallbackPage() {
  return (
    <Suspense
      fallback={
        <div className="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900">
          <div className="text-center">
            <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
            <p className="mt-4 text-gray-600 dark:text-gray-400">
              読み込み中...
            </p>
          </div>
        </div>
      }
    >
      <CallbackHandler />
    </Suspense>
  );
}
