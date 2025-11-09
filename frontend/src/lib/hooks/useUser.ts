'use client';

import { useState, useEffect } from 'react';

/**
 * ユーザーセッション管理フック
 * 本番環境では認証システムと統合する必要があります
 * 現在はローカルストレージを使用した簡易実装
 */
export function useUser() {
  const [userId, setUserId] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // ローカルストレージからユーザーIDを取得
    const storedUserId = localStorage.getItem('userId');
    
    if (storedUserId) {
      setUserId(storedUserId);
    } else {
      // ユーザーIDが存在しない場合は新規作成（デモ用）
      const newUserId = `user_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
      localStorage.setItem('userId', newUserId);
      setUserId(newUserId);
    }
    
    setLoading(false);
  }, []);

  const clearUser = () => {
    localStorage.removeItem('userId');
    setUserId(null);
  };

  return {
    userId,
    loading,
    clearUser,
  };
}
