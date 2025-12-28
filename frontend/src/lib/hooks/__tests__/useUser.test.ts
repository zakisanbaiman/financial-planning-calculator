import { renderHook, act, waitFor } from '@testing-library/react';
import { useUser } from '../useUser';

// localStorageモックを取得
const localStorageMock = window.localStorage as jest.Mocked<Storage>;

describe('useUser', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('初期化', () => {
    it('既存のユーザーIDがある場合、それを使用する', async () => {
      const existingUserId = 'existing_user_123';
      localStorageMock.getItem.mockReturnValue(existingUserId);
      
      const { result } = renderHook(() => useUser());
      
      await waitFor(() => {
        expect(result.current.loading).toBe(false);
      });
      
      expect(result.current.userId).toBe(existingUserId);
      expect(localStorageMock.setItem).not.toHaveBeenCalled();
    });

    it('ユーザーIDがない場合、新規作成して保存する', async () => {
      localStorageMock.getItem.mockReturnValue(null);
      
      const { result } = renderHook(() => useUser());
      
      await waitFor(() => {
        expect(result.current.loading).toBe(false);
      });
      
      expect(result.current.userId).not.toBeNull();
      expect(result.current.userId).toMatch(/^user_\d+_[a-z0-9]+$/);
      expect(localStorageMock.setItem).toHaveBeenCalledWith(
        'userId',
        expect.stringMatching(/^user_\d+_[a-z0-9]+$/)
      );
    });
  });

  describe('clearUser', () => {
    it('ユーザーをクリアすると userId が null になる', async () => {
      const existingUserId = 'existing_user_123';
      localStorageMock.getItem.mockReturnValue(existingUserId);
      
      const { result } = renderHook(() => useUser());
      
      await waitFor(() => {
        expect(result.current.userId).toBe(existingUserId);
      });
      
      act(() => {
        result.current.clearUser();
      });
      
      expect(result.current.userId).toBeNull();
      expect(localStorageMock.removeItem).toHaveBeenCalledWith('userId');
    });
  });

  describe('返り値の型', () => {
    it('userId, loading, clearUser を返す', async () => {
      localStorageMock.getItem.mockReturnValue('user_123');
      
      const { result } = renderHook(() => useUser());
      
      await waitFor(() => {
        expect(result.current.loading).toBe(false);
      });
      
      expect(result.current).toHaveProperty('userId');
      expect(result.current).toHaveProperty('loading');
      expect(result.current).toHaveProperty('clearUser');
      expect(typeof result.current.clearUser).toBe('function');
    });
  });
});
