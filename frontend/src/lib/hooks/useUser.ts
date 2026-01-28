'use client';

import { useAuth } from '../contexts/AuthContext';
import { useGuestMode } from '../contexts/GuestModeContext';

/**
 * ユーザーセッション管理フック
 * AuthContextと統合され、JWT認証を使用します
 * ゲストモードの場合は特別なIDを返します
 */
export function useUser() {
  const { user, isLoading, logout } = useAuth();
  const { isGuestMode } = useGuestMode();

  // ゲストモードの場合は固定のIDを返す
  const userId = isGuestMode ? 'guest' : (user?.userId || null);
  const email = isGuestMode ? null : (user?.email || null);

  return {
    userId,
    email,
    loading: isLoading,
    clearUser: logout,
    isGuest: isGuestMode,
  };
}
