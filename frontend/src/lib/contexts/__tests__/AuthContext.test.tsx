import React, { ReactNode } from 'react';
import { renderHook, act, waitFor } from '@testing-library/react';
import { AuthProvider, useAuth } from '../AuthContext';

// next/navigation モック
const mockPush = jest.fn();
jest.mock('next/navigation', () => ({
  useRouter: () => ({
    push: mockPush,
  }),
}));

// テスト用ラッパー
const wrapper = ({ children }: { children: ReactNode }) => (
  <AuthProvider>{children}</AuthProvider>
);

// ログイン成功レスポンス
const mockLoginSuccessResponse = {
  user_id: 'user_123',
  email: 'test@example.com',
  token: 'mock-token',
  refresh_token: 'mock-refresh-token',
  expires_at: '2026-12-31T00:00:00Z',
};

// initAuth の /api/auth/me fetch に対してデフォルトで 401 を返すモックを設定する
const mockInitAuthFailure = () => {
  (global.fetch as jest.Mock).mockResolvedValueOnce({
    ok: false,
    status: 401,
    json: async () => ({ error: 'Unauthorized' }),
  });
};

describe('AuthContext - login 関数', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    (global.fetch as jest.Mock).mockReset();
    // localStorage.getItem は初期化時に null を返す（未ログイン状態）
    (window.localStorage.getItem as jest.Mock).mockReturnValue(null);
  });

  describe('正常系', () => {
    it('ログイン成功時に isLoading が最終的に false になる', async () => {
      // initAuth の /api/auth/me fetch 用モック（1回目）
      mockInitAuthFailure();
      // login の /api/auth/login fetch 用モック（2回目）
      (global.fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => mockLoginSuccessResponse,
      });

      const { result } = renderHook(() => useAuth(), { wrapper });

      // 初期化完了を待つ（useEffect で setIsLoading(false) が呼ばれる）
      await waitFor(() => {
        expect(result.current.isLoading).toBe(false);
      });

      await act(async () => {
        await result.current.login('test@example.com', 'password123');
      });

      expect(result.current.isLoading).toBe(false);
      expect(result.current.error).toBeNull();
    });

    it('ログイン成功時に router.push が呼ばれる', async () => {
      // initAuth の /api/auth/me fetch 用モック（1回目）
      mockInitAuthFailure();
      // login の /api/auth/login fetch 用モック（2回目）
      (global.fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => mockLoginSuccessResponse,
      });

      const { result } = renderHook(() => useAuth(), { wrapper });

      await waitFor(() => {
        expect(result.current.isLoading).toBe(false);
      });

      await act(async () => {
        await result.current.login('test@example.com', 'password123');
      });

      expect(mockPush).toHaveBeenCalledWith('/dashboard');
    });
  });

  describe('エラー系', () => {
    it('fetch が失敗した場合に isLoading が false になる', async () => {
      // initAuth の /api/auth/me fetch 用モック（1回目）
      mockInitAuthFailure();
      // login の /api/auth/login fetch 用モック（2回目）
      (global.fetch as jest.Mock).mockRejectedValueOnce(
        new Error('Network error')
      );

      const { result } = renderHook(() => useAuth(), { wrapper });

      await waitFor(() => {
        expect(result.current.isLoading).toBe(false);
      });

      await act(async () => {
        try {
          await result.current.login('test@example.com', 'password123');
        } catch {
          // エラーがスローされることを期待
        }
      });

      expect(result.current.isLoading).toBe(false);
    });

    it('fetch が失敗した場合にエラーメッセージが設定される', async () => {
      // initAuth の /api/auth/me fetch 用モック（1回目）
      mockInitAuthFailure();
      // login の /api/auth/login fetch 用モック（2回目）
      (global.fetch as jest.Mock).mockRejectedValueOnce(
        new Error('Network error')
      );

      const { result } = renderHook(() => useAuth(), { wrapper });

      await waitFor(() => {
        expect(result.current.isLoading).toBe(false);
      });

      await act(async () => {
        try {
          await result.current.login('test@example.com', 'password123');
        } catch {
          // エラーがスローされることを期待
        }
      });

      expect(result.current.error).toBe('Network error');
    });

    it('401レスポンス時に isLoading が false になり、適切なエラーメッセージが設定される', async () => {
      // initAuth の /api/auth/me fetch 用モック（1回目）
      mockInitAuthFailure();
      // login の /api/auth/login fetch 用モック（2回目）
      (global.fetch as jest.Mock).mockResolvedValueOnce({
        ok: false,
        status: 401,
        json: async () => ({ error: 'Unauthorized' }),
      });

      const { result } = renderHook(() => useAuth(), { wrapper });

      await waitFor(() => {
        expect(result.current.isLoading).toBe(false);
      });

      await act(async () => {
        try {
          await result.current.login('test@example.com', 'wrong-password');
        } catch {
          // エラーがスローされることを期待
        }
      });

      expect(result.current.isLoading).toBe(false);
      expect(result.current.error).toBe(
        'メールアドレスまたはパスワードが正しくありません'
      );
    });
  });

  describe('タイムアウト系', () => {
    it('fetch がタイムアウト（AbortError）した場合に isLoading が false になる', async () => {
      const abortError = new DOMException('The operation was aborted.', 'AbortError');
      // initAuth の /api/auth/me fetch 用モック（1回目）
      mockInitAuthFailure();
      // login の /api/auth/login fetch 用モック（2回目）
      (global.fetch as jest.Mock).mockRejectedValueOnce(abortError);

      const { result } = renderHook(() => useAuth(), { wrapper });

      await waitFor(() => {
        expect(result.current.isLoading).toBe(false);
      });

      await act(async () => {
        try {
          await result.current.login('test@example.com', 'password123');
        } catch {
          // タイムアウトエラーがスローされることを期待
        }
      });

      expect(result.current.isLoading).toBe(false);
    });

    it('fetch がタイムアウト（AbortError）した場合にタイムアウトエラーメッセージが設定される', async () => {
      const abortError = new DOMException('The operation was aborted.', 'AbortError');
      // initAuth の /api/auth/me fetch 用モック（1回目）
      mockInitAuthFailure();
      // login の /api/auth/login fetch 用モック（2回目）
      (global.fetch as jest.Mock).mockRejectedValueOnce(abortError);

      const { result } = renderHook(() => useAuth(), { wrapper });

      await waitFor(() => {
        expect(result.current.isLoading).toBe(false);
      });

      await act(async () => {
        try {
          await result.current.login('test@example.com', 'password123');
        } catch {
          // タイムアウトエラーがスローされることを期待
        }
      });

      // タイムアウト時はタイムアウト関連のエラーメッセージが設定される（修正後の実装を期待）
      expect(result.current.error).toBeTruthy();
      expect(result.current.error).toMatch(/タイムアウト|接続がタイムアウト|リクエストがタイムアウト/);
    });
  });

  describe('API_BASE_URL 確認', () => {
    it('login 時の fetch が相対パス /api/auth/login を使っている（絶対 URL を使っていない）', async () => {
      // initAuth の /api/auth/me fetch 用モック（1回目）
      mockInitAuthFailure();
      // login の /api/auth/login fetch 用モック（2回目）
      (global.fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => mockLoginSuccessResponse,
      });

      const { result } = renderHook(() => useAuth(), { wrapper });

      await waitFor(() => {
        expect(result.current.isLoading).toBe(false);
      });

      await act(async () => {
        await result.current.login('test@example.com', 'password123');
      });

      // fetch が呼ばれた URL を確認（2回目の呼び出しが /api/auth/login）
      expect(global.fetch).toHaveBeenCalled();
      const fetchCalls = (global.fetch as jest.Mock).mock.calls;
      const loginFetchUrl = fetchCalls[fetchCalls.length - 1][0] as string;

      // 相対パスであること（絶対 URL ではないこと）
      expect(loginFetchUrl).toBe('/api/auth/login');
      expect(loginFetchUrl).not.toMatch(/^https?:\/\//);
    });
  });
});
