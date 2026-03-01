import { renderHook } from '@testing-library/react';
import { useUser } from '../useUser';

// useAuth と useGuestMode をモック
jest.mock('../../contexts/AuthContext', () => ({
  useAuth: jest.fn(),
}));
jest.mock('../../contexts/GuestModeContext', () => ({
  useGuestMode: jest.fn(),
}));

import { useAuth } from '../../contexts/AuthContext';
import { useGuestMode } from '../../contexts/GuestModeContext';

const mockUseAuth = useAuth as jest.MockedFunction<typeof useAuth>;
const mockUseGuestMode = useGuestMode as jest.MockedFunction<typeof useGuestMode>;

describe('useUser', () => {
  const mockLogout = jest.fn();

  beforeEach(() => {
    jest.clearAllMocks();
    mockUseGuestMode.mockReturnValue({
      isGuestMode: false,
      guestData: null,
      setGuestData: jest.fn(),
      clearGuestData: jest.fn(),
      startGuestMode: jest.fn(),
      exitGuestMode: jest.fn(),
    });
  });

  describe('認証済みユーザー', () => {
    it('認証済みユーザーの userId を返す', () => {
      mockUseAuth.mockReturnValue({
        user: { userId: 'existing_user_123', email: 'test@example.com' },
        isAuthenticated: true,
        isLoading: false,
        login: jest.fn(),
        register: jest.fn(),
        logout: mockLogout,
        error: null,
        clearError: jest.fn(),
        setAuthData: jest.fn(),
      });

      const { result } = renderHook(() => useUser());

      expect(result.current.userId).toBe('existing_user_123');
      expect(result.current.email).toBe('test@example.com');
      expect(result.current.loading).toBe(false);
      expect(result.current.isGuest).toBe(false);
    });

    it('ロード中は loading が true になる', () => {
      mockUseAuth.mockReturnValue({
        user: null,
        isAuthenticated: false,
        isLoading: true,
        login: jest.fn(),
        register: jest.fn(),
        logout: mockLogout,
        error: null,
        clearError: jest.fn(),
        setAuthData: jest.fn(),
      });

      const { result } = renderHook(() => useUser());

      expect(result.current.loading).toBe(true);
      expect(result.current.userId).toBeNull();
    });

    it('未ログイン時は userId が null になる', () => {
      mockUseAuth.mockReturnValue({
        user: null,
        isAuthenticated: false,
        isLoading: false,
        login: jest.fn(),
        register: jest.fn(),
        logout: mockLogout,
        error: null,
        clearError: jest.fn(),
        setAuthData: jest.fn(),
      });

      const { result } = renderHook(() => useUser());

      expect(result.current.userId).toBeNull();
      expect(result.current.email).toBeNull();
    });
  });

  describe('ゲストモード', () => {
    it('ゲストモード時は userId が "guest" になる', () => {
      mockUseAuth.mockReturnValue({
        user: null,
        isAuthenticated: false,
        isLoading: false,
        login: jest.fn(),
        register: jest.fn(),
        logout: mockLogout,
        error: null,
        clearError: jest.fn(),
        setAuthData: jest.fn(),
      });
      mockUseGuestMode.mockReturnValue({
        isGuestMode: true,
        guestData: null,
        setGuestData: jest.fn(),
        clearGuestData: jest.fn(),
        startGuestMode: jest.fn(),
        exitGuestMode: jest.fn(),
      });

      const { result } = renderHook(() => useUser());

      expect(result.current.userId).toBe('guest');
      expect(result.current.email).toBeNull();
      expect(result.current.isGuest).toBe(true);
    });
  });

  describe('返り値の型', () => {
    it('userId, email, loading, clearUser, isGuest を返す', () => {
      mockUseAuth.mockReturnValue({
        user: { userId: 'user_123', email: 'user@example.com' },
        isAuthenticated: true,
        isLoading: false,
        login: jest.fn(),
        register: jest.fn(),
        logout: mockLogout,
        error: null,
        clearError: jest.fn(),
        setAuthData: jest.fn(),
      });

      const { result } = renderHook(() => useUser());

      expect(result.current).toHaveProperty('userId');
      expect(result.current).toHaveProperty('email');
      expect(result.current).toHaveProperty('loading');
      expect(result.current).toHaveProperty('clearUser');
      expect(result.current).toHaveProperty('isGuest');
      expect(typeof result.current.clearUser).toBe('function');
    });

    it('clearUser を呼ぶと logout が実行される', () => {
      mockUseAuth.mockReturnValue({
        user: { userId: 'user_123', email: 'user@example.com' },
        isAuthenticated: true,
        isLoading: false,
        login: jest.fn(),
        register: jest.fn(),
        logout: mockLogout,
        error: null,
        clearError: jest.fn(),
        setAuthData: jest.fn(),
      });

      const { result } = renderHook(() => useUser());

      result.current.clearUser();

      expect(mockLogout).toHaveBeenCalledTimes(1);
    });
  });
});
