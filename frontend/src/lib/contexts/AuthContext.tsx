'use client';

import React, { createContext, useContext, useState, useEffect, useCallback, ReactNode } from 'react';
import { useRouter } from 'next/navigation';

// 認証ユーザー情報
export interface AuthUser {
  userId: string;
  email: string;
}

// 認証コンテキストの型
interface AuthContextType {
  user: AuthUser | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (email: string, password: string) => Promise<void>;
  register: (email: string, password: string) => Promise<void>;
  logout: () => void;
  error: string | null;
  clearError: () => void;
  setAuthData: (data: { user: AuthUser }) => void; // Issue: #67 - トークン情報を削除
}

// 認証レスポンスの型
interface AuthResponse {
  user_id: string;
  email: string;
  token: string;
  refresh_token?: string; // 2FA有効時は空
  expires_at?: string;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

// ローカルストレージのキー（ユーザー情報のみ保存）
const USER_KEY = 'auth_user';
// 2FA用の一時トークンキー
const TEMP_TOKEN_KEY = 'auth_token';

// API ベースURL
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api';

interface AuthProviderProps {
  children: ReactNode;
}

export function AuthProvider({ children }: AuthProviderProps) {
  const [user, setUser] = useState<AuthUser | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const router = useRouter();

  // 初期化時にローカルストレージからユーザー情報を復元
  // トークンはCookieで管理されるため、ここでは復元しない
  useEffect(() => {
    const initAuth = () => {
      try {
        const storedUser = localStorage.getItem(USER_KEY);

        if (storedUser) {
          setUser(JSON.parse(storedUser));
        }
      } catch (e) {
        console.error('Failed to restore auth state:', e);
      } finally {
        setIsLoading(false);
      }
    };

    initAuth();
  }, []);

  // 認証情報を保存（ユーザー情報のみ）
  const saveAuthData = useCallback((response: AuthResponse) => {
    const authUser: AuthUser = {
      userId: response.user_id,
      email: response.email,
    };

    localStorage.setItem(USER_KEY, JSON.stringify(authUser));
    setUser(authUser);
  }, []);

  // 認証情報をクリア
  const clearAuthData = useCallback(() => {
    localStorage.removeItem(USER_KEY);
    setUser(null);
  }, []);

  // ログイン
  const login = useCallback(async (email: string, password: string) => {
    setIsLoading(true);
    setError(null);

    try {
      const response = await fetch(`${API_BASE_URL}/auth/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email, password }),
        credentials: 'include', // Cookieを送受信
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        if (response.status === 401) {
          throw new Error('メールアドレスまたはパスワードが正しくありません');
        }
        throw new Error(errorData.error || errorData.message || 'ログインに失敗しました');
      }

      const data: AuthResponse = await response.json();
      
      // リフレッシュトークンが空の場合は2FA検証が必要
      // 仮トークンはCookieで自動的に設定されるため、localStorageへの保存は不要
      if (!data.refresh_token) {
        // 2FA検証ページにリダイレクト
        router.push('/auth/2fa-verify');
      } else {
        // 通常のログイン
        saveAuthData(data);
        router.push('/dashboard');
      }
    } catch (e) {
      const message = e instanceof Error ? e.message : 'ログインに失敗しました';
      setError(message);
      throw e;
    } finally {
      setIsLoading(false);
    }
  }, [router, saveAuthData]);

  // ユーザー登録
  const register = useCallback(async (email: string, password: string) => {
    setIsLoading(true);
    setError(null);

    try {
      const response = await fetch(`${API_BASE_URL}/auth/register`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email, password }),
        credentials: 'include', // Cookieを送受信
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        if (response.status === 409) {
          throw new Error('このメールアドレスは既に登録されています');
        }
        throw new Error(errorData.error || errorData.message || '登録に失敗しました');
      }

      const data: AuthResponse = await response.json();
      saveAuthData(data);
      router.push('/dashboard');
    } catch (e) {
      const message = e instanceof Error ? e.message : '登録に失敗しました';
      setError(message);
      throw e;
    } finally {
      setIsLoading(false);
    }
  }, [router, saveAuthData]);

  // ログアウト
  const logout = useCallback(async () => {
    try {
      // バックエンドのログアウトエンドポイントを呼び出してCookieをクリア
      await fetch(`${API_BASE_URL}/auth/logout`, {
        method: 'POST',
        credentials: 'include',
      });
    } catch (e) {
      console.error('Logout API call failed:', e);
      // エラーが発生してもローカルの状態はクリアする
    }
    
    clearAuthData();
    router.push('/login');
  }, [clearAuthData, router]);

  // エラーをクリア
  const clearError = useCallback(() => {
    setError(null);
  }, []);

  // OAuth用：外部からユーザーデータを設定（Issue: #67）
  // トークンはCookieで管理されるため、ユーザー情報のみを受け取る
  const setAuthData = useCallback((data: { user: AuthUser }) => {
    localStorage.setItem(USER_KEY, JSON.stringify(data.user));
    setUser(data.user);
    setError(null);
  }, []);

  const value = {
    user,
    isAuthenticated: !!user,
    isLoading,
    login,
    register,
    logout,
    error,
    clearError,
    setAuthData, // Issue: #67
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

// 認証コンテキストを使用するフック
export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}
