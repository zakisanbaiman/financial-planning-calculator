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

const mockLogout = jest.fn();

describe('useUser', () => {
  const defaultGuestMode = {
    isGuestMode: false,
    guestData: null,
    setGuestData: jest.fn(),
    clearGuestData: jest.fn(),
    startGuestMode: jest.fn(),
    exitGuestMode: jest.fn(),
  };

  beforeEach(() => {
    jest.clearAllMocks();
    mockUseGuestMode.mockReturnValue(defaultGuestMode);
  });

  describe('認証済みユーザー', () => {
    it('ログイン済みの場合、userId を返す', () => {
      mockUseAuth.mockReturnValue({
        user: { userId: 'user-001', email: 'test@example.com' },
        isLoading: false,
        isAuthenticated: true,
        login: jest.fn(),
        register: jest.fn(),
        logout: mockLogout,
        error: null,
        clearError: jest.fn(),
        setAuthData: jest.fn(),
      });

      const { result } = renderHook(() => useUser());

      expect(result.current.userId).toBe('user-001');
      expect(result.current.email).toBe('test@example.com');
      expect(result.current.loading).toBe(false);
      expect(result.current.isGuest).toBe(false);
    });

    it('未ログインの場合、userId が null になる', () => {
      mockUseAuth.mockReturnValue({
        user: null,
        isLoading: false,
        isAuthenticated: false,
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

    it('ローディング中は loading が true になる', () => {
      mockUseAuth.mockReturnValue({
        user: null,
        isLoading: true,
        isAuthenticated: false,
        login: jest.fn(),
        register: jest.fn(),
        logout: mockLogout,
        error: null,
        clearError: jest.fn(),
        setAuthData: jest.fn(),
      });

      const { result } = renderHook(() => useUser());

      expect(result.current.loading).toBe(true);
    });
  });

  describe('ゲストモード', () => {
    it('ゲストモードの場合、userId が "guest" になる', () => {
      mockUseAuth.mockReturnValue({
        user: null,
        isLoading: false,
        isAuthenticated: false,
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
      expect(result.current.isGuest).toBe(true);
    });
  });

  describe('clearUser', () => {
    it('clearUser を呼ぶと logout が実行される', () => {
      mockUseAuth.mockReturnValue({
        user: { userId: 'user-001', email: 'test@example.com' },
        isLoading: false,
        isAuthenticated: true,
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

  describe('返り値の型', () => {
    it('userId, loading, clearUser, email, isGuest を返す', () => {
      mockUseAuth.mockReturnValue({
        user: { userId: 'user-001', email: 'test@example.com' },
        isLoading: false,
        isAuthenticated: true,
        login: jest.fn(),
        register: jest.fn(),
        logout: mockLogout,
        error: null,
        clearError: jest.fn(),
        setAuthData: jest.fn(),
      });

      const { result } = renderHook(() => useUser());

      expect(result.current).toHaveProperty('userId');
      expect(result.current).toHaveProperty('loading');
      expect(result.current).toHaveProperty('clearUser');
      expect(result.current).toHaveProperty('email');
      expect(result.current).toHaveProperty('isGuest');
      expect(typeof result.current.clearUser).toBe('function');
    });
  });
});
