'use client';

import { useEffect } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { useAuth } from '@/lib/contexts/AuthContext';

/**
 * OAuthコールバックページ（Issue: #67）
 * GitHub OAuth認証後のリダイレクト先
 * URLパラメータからトークンを取得してAuthContextに保存
 */
export default function AuthCallbackPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { setAuthData } = useAuth();

  useEffect(() => {
    const handleCallback = async () => {
      // URLパラメータからトークン情報を取得
      const token = searchParams.get('token');
      const refreshToken = searchParams.get('refresh_token');
      const userId = searchParams.get('user_id');
      const email = searchParams.get('email');
      const error = searchParams.get('error');

      // エラーがある場合
      if (error) {
        console.error('OAuth authentication failed:', error);
        router.push(`/login?error=${error}`);
        return;
      }

      // トークンが存在する場合
      if (token && refreshToken && userId && email) {
        // AuthContextに認証情報を保存
        setAuthData({
          token,
          refreshToken,
          user: {
            userId,
            email,
          },
        });

        // ダッシュボードにリダイレクト
        router.push('/dashboard');
      } else {
        // トークン情報が不足している場合
        console.error('Missing authentication data in callback');
        router.push('/login?error=missing_data');
      }
    };

    handleCallback();
  }, [searchParams, router, setAuthData]);

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
