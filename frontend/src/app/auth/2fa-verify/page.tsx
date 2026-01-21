'use client';

import { useState, useEffect, Suspense } from 'react';
import { useRouter } from 'next/navigation';
import { twoFactorAPI } from '@/lib/api-client';
import { useAuth } from '@/lib/contexts/AuthContext';

function TwoFactorVerifyContent() {
  const router = useRouter();
  const { setAuthData } = useAuth();
  const [code, setCode] = useState('');
  const [useBackup, setUseBackup] = useState(false);
  const [useBackup, setUseBackup] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    // 仮トークンがない場合はログインページにリダイレクト
    const tempToken = typeof window !== 'undefined' ? localStorage.getItem('auth_token') : null;
    if (!tempToken) {
      router.push('/login');
    }
  }, [router]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!code) {
      setError(useBackup ? 'バックアップコードを入力してください' : '認証コードを入力してください');
      return;
    }

    try {
      setLoading(true);
      setError(null);
      
      // 2FA検証
      const response = await twoFactorAPI.verify(code, useBackup);
      
      // AuthContextを使ってトークンを保存（UIも更新される）
      setAuthData({
        token: response.token,
        refreshToken: response.refresh_token,
        user: {
          userId: response.user_id,
          email: response.email,
        },
      });
      
      // ダッシュボードにリダイレクト
      router.push('/dashboard');
    } catch (err: any) {
      setError(err.message || '認証コードが無効です');
    } finally {
      setLoading(false);
    }
  };

  const handleCodeChange = (value: string) => {
    if (useBackup) {
      // バックアップコードは8文字の英数字
      setCode(value.toUpperCase().slice(0, 8));
    } else {
      // TOTPコードは6桁の数字
      setCode(value.replace(/\D/g, '').slice(0, 6));
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full">
        <div className="bg-white shadow rounded-lg p-8">
          <div className="text-center mb-6">
            <h1 className="text-2xl font-bold text-gray-900">2段階認証</h1>
            <p className="mt-2 text-sm text-gray-600">
              {useBackup 
                ? 'バックアップコードを入力してください' 
                : '認証アプリに表示されている6桁のコードを入力してください'}
            </p>
          </div>

          {error && (
            <div className="mb-4 bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
              {error}
            </div>
          )}

          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label htmlFor="code" className="block text-sm font-medium text-gray-700 mb-1">
                {useBackup ? 'バックアップコード' : '認証コード'}
              </label>
              <input
                id="code"
                type="text"
                value={code}
                onChange={(e) => handleCodeChange(e.target.value)}
                placeholder={useBackup ? 'XXXXXXXX' : '000000'}
                maxLength={useBackup ? 8 : 6}
                className="w-full px-4 py-3 text-center text-2xl tracking-widest border border-gray-300 rounded focus:ring-blue-500 focus:border-blue-500 font-mono"
                autoComplete="off"
                autoFocus
                required
              />
              <p className="mt-1 text-xs text-gray-500">
                {useBackup ? '8文字の英数字' : '6桁の数字'}
              </p>
            </div>

            <button
              type="submit"
              disabled={loading || (useBackup ? code.length !== 8 : code.length !== 6)}
              className="w-full bg-blue-600 text-white py-3 rounded font-semibold hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {loading ? '確認中...' : '確認'}
            </button>
          </form>

          <div className="mt-6 text-center">
            <button
              onClick={() => {
                setUseBackup(!useBackup);
                setCode('');
                setError(null);
              }}
              className="text-sm text-blue-600 hover:text-blue-800"
            >
              {useBackup 
                ? '← 認証コードを使用' 
                : 'バックアップコードを使用 →'}
            </button>
          </div>

          <div className="mt-4 text-center">
            <button
              onClick={() => {
                if (typeof window !== 'undefined') {
                  localStorage.removeItem('auth_token');
                  localStorage.removeItem('refresh_token');
                  localStorage.removeItem('auth_expires');
                  localStorage.removeItem('auth_user');
                }
                router.push('/login');
              }}
              className="text-sm text-gray-600 hover:text-gray-800"
            >
              ログアウト
            </button>
          </div>
        </div>

        <div className="mt-6 bg-blue-50 border border-blue-200 rounded p-4">
          <p className="text-sm text-blue-800">
            <strong>💡 ヒント:</strong> 認証アプリにアクセスできない場合は、バックアップコードを使用してログインできます。
          </p>
        </div>
      </div>
    </div>
  );
}

export default function TwoFactorVerifyPage() {
  return (
    <Suspense fallback={<div className="min-h-screen flex items-center justify-center">読み込み中...</div>}>
      <TwoFactorVerifyContent />
    </Suspense>
  );
}
