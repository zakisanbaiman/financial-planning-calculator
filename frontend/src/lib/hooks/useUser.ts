'use client';

import { useAuth } from '../contexts/AuthContext';

/**
 * ユーザーセッション管理フック
 * AuthContextと統合され、JWT認証を使用します
 */
export function useUser() {
  const { user, isLoading, logout } = useAuth();

  return {
    userId: user?.userId || null,
    email: user?.email || null,
    loading: isLoading,
    clearUser: logout,
  };
}
