'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { twoFactorAPI } from '@/lib/api-client';
import { QRCodeSVG } from 'qrcode.react';

export default function SecurityPage() {
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const [initialLoading, setInitialLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  
  // 2FA設定状態
  const [is2FAEnabled, setIs2FAEnabled] = useState(false);
  const [showSetup, setShowSetup] = useState(false);
  const [setupData, setSetupData] = useState<{
    secret: string;
    qr_code_url: string;
    backup_codes: string[];
  } | null>(null);
  const [verificationCode, setVerificationCode] = useState('');
  const [disablePassword, setDisablePassword] = useState('');
  const [showBackupCodes, setShowBackupCodes] = useState(false);
  const [newBackupCodes, setNewBackupCodes] = useState<string[]>([]);

  // ページロード時に2FAステータスを取得
  useEffect(() => {
    const fetch2FAStatus = async () => {
      try {
        const data = await twoFactorAPI.getStatus();
        setIs2FAEnabled(data.enabled);
      } catch (err: any) {
        // 認証エラーの場合はログインページにリダイレクト
        if (err.status === 401) {
          router.push('/login');
          return;
        }
        setError(err.message || '2FAステータスの取得に失敗しました');
      } finally {
        setInitialLoading(false);
      }
    };

    fetch2FAStatus();
  }, [router]);

  // 2FA設定開始
  const handleSetup2FA = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await twoFactorAPI.setup();
      setSetupData(data);
      setShowSetup(true);
    } catch (err: any) {
      setError(err.message || '2FA設定の開始に失敗しました');
    } finally {
      setLoading(false);
    }
  };

  // 2FA有効化
  const handleEnable2FA = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!setupData || !verificationCode) {
      setError('認証コードを入力してください');
      return;
    }

    try {
      setLoading(true);
      setError(null);
      setSuccess(null);
      console.log('2FA有効化リクエスト送信中...', { code: verificationCode, secretLength: setupData.secret.length });
      await twoFactorAPI.enable(verificationCode, setupData.secret);
      setSuccess('2段階認証が有効になりました');
      setIs2FAEnabled(true);
      setShowSetup(false);
      setVerificationCode('');
      setSetupData(null);
    } catch (err: any) {
      console.error('2FA有効化エラー:', err);
      const errorMessage = err.data?.error || err.message || '2FAの有効化に失敗しました';
      setError(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  // 2FA無効化
  const handleDisable2FA = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!disablePassword) {
      setError('パスワードを入力してください');
      return;
    }

    if (!confirm('2段階認証を無効にしますか？セキュリティが低下します。')) {
      return;
    }

    try {
      setLoading(true);
      setError(null);
      await twoFactorAPI.disable(disablePassword);
      setSuccess('2段階認証が無効になりました');
      setIs2FAEnabled(false);
      setDisablePassword('');
    } catch (err: any) {
      setError(err.message || '2FAの無効化に失敗しました');
    } finally {
      setLoading(false);
    }
  };

  // バックアップコード再生成
  const handleRegenerateBackupCodes = async () => {
    if (!confirm('現在のバックアップコードは無効になります。続けますか？')) {
      return;
    }

    try {
      setLoading(true);
      setError(null);
      const data = await twoFactorAPI.regenerateBackupCodes();
      setNewBackupCodes(data.backup_codes);
      setShowBackupCodes(true);
      setSuccess('バックアップコードを再生成しました');
    } catch (err: any) {
      setError(err.message || 'バックアップコードの再生成に失敗しました');
    } finally {
      setLoading(false);
    }
  };

  // バックアップコードのダウンロード
  const downloadBackupCodes = (codes: string[]) => {
    const text = codes.join('\n');
    const blob = new Blob([text], { type: 'text/plain' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'backup-codes.txt';
    a.click();
    URL.revokeObjectURL(url);
  };

  // 初期ロード中の表示
  if (initialLoading) {
    return (
      <div className="min-h-screen bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
        <div className="max-w-3xl mx-auto">
          <div className="bg-white shadow rounded-lg p-6">
            <div className="flex justify-center items-center py-12">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
              <span className="ml-3 text-gray-600">読み込み中...</span>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-3xl mx-auto">
        <div className="bg-white shadow rounded-lg">
          <div className="px-6 py-4 border-b border-gray-200">
            <h1 className="text-2xl font-bold text-gray-900">セキュリティ設定</h1>
            <p className="mt-1 text-sm text-gray-600">
              アカウントのセキュリティを強化するための設定
            </p>
          </div>

          <div className="p-6 space-y-6">
            {/* エラー・成功メッセージ */}
            {error && (
              <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
                {error}
              </div>
            )}
            {success && (
              <div className="bg-green-50 border border-green-200 text-green-700 px-4 py-3 rounded">
                {success}
              </div>
            )}

            {/* 2段階認証セクション */}
            <div className="border-b border-gray-200 pb-6">
              <h2 className="text-lg font-semibold text-gray-900 mb-2">
                2段階認証（2FA）
              </h2>
              <p className="text-sm text-gray-600 mb-4">
                Google Authenticator などの認証アプリを使用して、ログイン時に追加のセキュリティコードを要求します。
              </p>

              {!is2FAEnabled && !showSetup && (
                <button
                  onClick={handleSetup2FA}
                  disabled={loading}
                  className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700 disabled:opacity-50"
                >
                  {loading ? '処理中...' : '2段階認証を有効にする'}
                </button>
              )}

              {/* 2FA設定フロー */}
              {showSetup && setupData && (
                <div className="space-y-4">
                  <div className="bg-gray-50 p-4 rounded">
                    <h3 className="font-semibold mb-2">ステップ1: QRコードをスキャン</h3>
                    <p className="text-sm text-gray-600 mb-4">
                      Google Authenticator または Authy などの認証アプリでこのQRコードをスキャンしてください。
                    </p>
                    <div className="flex justify-center my-4">
                      <QRCodeSVG value={setupData.qr_code_url} size={200} />
                    </div>
                    <p className="text-xs text-gray-500 text-center">
                      手動入力: {setupData.secret}
                    </p>
                  </div>

                  <div className="bg-gray-50 p-4 rounded">
                    <h3 className="font-semibold mb-2">ステップ2: バックアップコード</h3>
                    <p className="text-sm text-gray-600 mb-2">
                      認証アプリにアクセスできなくなった場合に使用できるバックアップコードです。
                      安全な場所に保管してください。
                    </p>
                    <div className="bg-white p-3 rounded border border-gray-300 font-mono text-sm space-y-1">
                      {setupData.backup_codes.map((code, index) => (
                        <div key={index}>{code}</div>
                      ))}
                    </div>
                    <button
                      onClick={() => downloadBackupCodes(setupData.backup_codes)}
                      className="mt-2 text-sm text-blue-600 hover:text-blue-800"
                    >
                      📥 ダウンロード
                    </button>
                  </div>

                  <form onSubmit={handleEnable2FA} className="space-y-4">
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">
                        ステップ3: 認証コードを入力
                      </label>
                      <input
                        type="text"
                        inputMode="numeric"
                        pattern="[0-9]*"
                        autoComplete="one-time-code"
                        value={verificationCode}
                        onChange={(e) => {
                          const value = e.target.value.replace(/[^0-9]/g, '').slice(0, 6);
                          setVerificationCode(value);
                        }}
                        placeholder="6桁のコード"
                        maxLength={6}
                        className="w-full px-3 py-2 border border-gray-300 rounded focus:ring-blue-500 focus:border-blue-500 text-lg tracking-widest"
                        required
                      />
                    </div>
                    <div className="flex gap-2">
                      <button
                        type="submit"
                        disabled={loading || verificationCode.length !== 6}
                        className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700 disabled:opacity-50"
                      >
                        {loading ? '処理中...' : '有効化'}
                      </button>
                      <button
                        type="button"
                        onClick={() => {
                          setShowSetup(false);
                          setSetupData(null);
                          setVerificationCode('');
                        }}
                        className="bg-gray-300 text-gray-700 px-4 py-2 rounded hover:bg-gray-400"
                      >
                        キャンセル
                      </button>
                    </div>
                  </form>
                </div>
              )}

              {/* 2FA有効時の管理 */}
              {is2FAEnabled && (
                <div className="space-y-4">
                  <div className="bg-green-50 border border-green-200 p-4 rounded">
                    <p className="text-green-800 font-semibold">✓ 2段階認証が有効です</p>
                    <p className="text-sm text-green-700 mt-1">
                      ログイン時に認証コードの入力が必要になります。
                    </p>
                  </div>

                  <div>
                    <button
                      onClick={handleRegenerateBackupCodes}
                      disabled={loading}
                      className="bg-gray-600 text-white px-4 py-2 rounded hover:bg-gray-700 disabled:opacity-50 mr-2"
                    >
                      バックアップコードを再生成
                    </button>
                  </div>

                  {showBackupCodes && newBackupCodes.length > 0 && (
                    <div className="bg-yellow-50 border border-yellow-200 p-4 rounded">
                      <h4 className="font-semibold mb-2">新しいバックアップコード</h4>
                      <div className="bg-white p-3 rounded border border-gray-300 font-mono text-sm space-y-1">
                        {newBackupCodes.map((code, index) => (
                          <div key={index}>{code}</div>
                        ))}
                      </div>
                      <button
                        onClick={() => downloadBackupCodes(newBackupCodes)}
                        className="mt-2 text-sm text-blue-600 hover:text-blue-800"
                      >
                        📥 ダウンロード
                      </button>
                    </div>
                  )}

                  <form onSubmit={handleDisable2FA} className="space-y-4 border-t pt-4">
                    <h3 className="font-semibold text-red-700">2段階認証を無効にする</h3>
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">
                        パスワードを入力
                      </label>
                      <input
                        type="password"
                        value={disablePassword}
                        onChange={(e) => setDisablePassword(e.target.value)}
                        placeholder="パスワード"
                        className="w-full px-3 py-2 border border-gray-300 rounded focus:ring-red-500 focus:border-red-500"
                        required
                      />
                    </div>
                    <button
                      type="submit"
                      disabled={loading}
                      className="bg-red-600 text-white px-4 py-2 rounded hover:bg-red-700 disabled:opacity-50"
                    >
                      {loading ? '処理中...' : '無効にする'}
                    </button>
                  </form>
                </div>
              )}
            </div>

            {/* ダッシュボードに戻るボタン */}
            <div>
              <button
                onClick={() => router.push('/dashboard')}
                className="text-blue-600 hover:text-blue-800"
              >
                ← ダッシュボードに戻る
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
