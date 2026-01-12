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
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (email: string, password: string) => Promise<void>;
  register: (email: string, password: string) => Promise<void>;
  logout: () => void;
  error: string | null;
  clearError: () => void;
}

// 認証レスポンスの型
interface AuthResponse {
  user_id: string;
  email: string;
  token: string;
  refresh_token: string;
  expires_at: string;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

// ローカルストレージのキー
const TOKEN_KEY = 'auth_token';
const USER_KEY = 'auth_user';
const EXPIRES_KEY = 'auth_expires';
const REFRESH_TOKEN_KEY = 'refresh_token';

// API ベースURL
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api';

interface AuthProviderProps {
  children: ReactNode;
}

export function AuthProvider({ children }: AuthProviderProps) {
  const [user, setUser] = useState<AuthUser | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const router = useRouter();

  // トークンの有効期限チェック
  const isTokenExpired = useCallback((expiresAt: string | null): boolean => {
    if (!expiresAt) return true;
    return new Date(expiresAt) <= new Date();
  }, []);

  // 初期化時にローカルストレージから認証情報を復元
  useEffect(() => {
    const initAuth = () => {
      try {
        const storedToken = localStorage.getItem(TOKEN_KEY);
        const storedUser = localStorage.getItem(USER_KEY);
        const storedExpires = localStorage.getItem(EXPIRES_KEY);

        if (storedToken && storedUser && !isTokenExpired(storedExpires)) {
          setToken(storedToken);
          setUser(JSON.parse(storedUser));
        } else {
          // トークンが無効な場合はクリア
          localStorage.removeItem(TOKEN_KEY);
          localStorage.removeItem(USER_KEY);
          localStorage.removeItem(EXPIRES_KEY);
        }
      } catch (e) {
        console.error('Failed to restore auth state:', e);
      } finally {
        setIsLoading(false);
      }
    };

    initAuth();
  }, [isTokenExpired]);

  // 認証情報を保存
  const saveAuthData = useCallback((response: AuthResponse) => {
    const authUser: AuthUser = {
      userId: response.user_id,
      email: response.email,
    };

    localStorage.setItem(TOKEN_KEY, response.token);
    localStorage.setItem(USER_KEY, JSON.stringify(authUser));
    localStorage.setItem(EXPIRES_KEY, response.expires_at);
    localStorage.setItem(REFRESH_TOKEN_KEY, response.refresh_token);

    setToken(response.token);
    setUser(authUser);
  }, []);

  // 認証情報をクリア
  const clearAuthData = useCallback(() => {
    localStorage.removeItem(TOKEN_KEY);
    localStorage.removeItem(USER_KEY);
    localStorage.removeItem(EXPIRES_KEY);
    localStorage.removeItem(REFRESH_TOKEN_KEY);
    setToken(null);
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
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        if (response.status === 401) {
          throw new Error('メールアドレスまたはパスワードが正しくありません');
        }
        throw new Error(errorData.error || errorData.message || 'ログインに失敗しました');
      }

      const data: AuthResponse = await response.json();
      saveAuthData(data);
      router.push('/dashboard');
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
  const logout = useCallback(() => {
    clearAuthData();
    router.push('/login');
  }, [clearAuthData, router]);

  // エラーをクリア
  const clearError = useCallback(() => {
    setError(null);
  }, []);

  const value: AuthContextType = {
    user,
    token,
    isAuthenticated: !!token && !!user,
    isLoading,
    login,
    register,
    logout,
    error,
    clearError,
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

// トークンを取得するヘルパー関数（APIクライアント用）
export function getAuthToken(): string | null {
  if (typeof window === 'undefined') return null;
  return localStorage.getItem(TOKEN_KEY);
}

// トークンの有効期限をチェックするヘルパー関数
export function isAuthTokenValid(): boolean {
  if (typeof window === 'undefined') return false;
  const expiresAt = localStorage.getItem(EXPIRES_KEY);
  if (!expiresAt) return false;
  return new Date(expiresAt) > new Date();
}
