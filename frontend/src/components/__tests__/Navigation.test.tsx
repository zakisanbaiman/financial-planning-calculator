import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import Navigation from '../Navigation';
import { ThemeProvider } from '@/lib/contexts/ThemeContext';
import { TutorialProvider } from '@/lib/contexts/TutorialContext';
import { AuthProvider } from '@/lib/contexts/AuthContext';
import { GuestModeProvider } from '@/lib/contexts/GuestModeContext';

// next/navigation モック
const mockUsePathname = jest.fn();
const mockPush = jest.fn();
jest.mock('next/navigation', () => ({
  usePathname: () => mockUsePathname(),
  useRouter: () => ({ push: mockPush }),
}));

function renderNavigation() {
  return render(
    <AuthProvider>
      <GuestModeProvider>
        <ThemeProvider>
          <TutorialProvider>
            <Navigation />
          </TutorialProvider>
        </ThemeProvider>
      </GuestModeProvider>
    </AuthProvider>
  );
}

describe('Navigation', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    mockUsePathname.mockReturnValue('/');
    // localStorage をデフォルト状態にリセット
    (localStorage.getItem as jest.Mock).mockReturnValue(null);
  });

  describe('基本レンダリング', () => {
    it('ロゴが表示される', () => {
      renderNavigation();
      expect(screen.getByText('FinPlan')).toBeInTheDocument();
    });

    it('ナビゲーションリンクが表示される', () => {
      renderNavigation();
      expect(screen.getAllByText('Home').length).toBeGreaterThanOrEqual(1);
      expect(screen.getAllByText('Dashboard').length).toBeGreaterThanOrEqual(1);
      expect(screen.getAllByText('Profile').length).toBeGreaterThanOrEqual(1);
      expect(screen.getAllByText('Goals').length).toBeGreaterThanOrEqual(1);
      expect(screen.getAllByText('Calculator').length).toBeGreaterThanOrEqual(1);
      expect(screen.getAllByText('Reports').length).toBeGreaterThanOrEqual(1);
    });

    it('ヘルプボタンが表示される', () => {
      renderNavigation();
      const helpButtons = screen.getAllByText('ヘルプ');
      expect(helpButtons.length).toBeGreaterThanOrEqual(1);
    });
  });

  describe('未認証状態', () => {
    it('ログインリンクが表示される', () => {
      renderNavigation();
      const loginLinks = screen.getAllByText('ログイン');
      expect(loginLinks.length).toBeGreaterThanOrEqual(1);
    });

    it('ログアウトボタンは表示されない', () => {
      renderNavigation();
      expect(screen.queryByText('ログアウト')).not.toBeInTheDocument();
    });
  });

  describe('ゲストモード', () => {
    it('ゲストモードではゲストモード表示と登録リンクが表示される', () => {
      (localStorage.getItem as jest.Mock).mockImplementation((key: string) => {
        if (key === 'guest_mode') return 'true';
        return null;
      });
      renderNavigation();
      const guestLabels = screen.getAllByText('ゲストモード');
      expect(guestLabels.length).toBeGreaterThanOrEqual(1);
      const registerLinks = screen.getAllByText('登録してデータを保存');
      expect(registerLinks.length).toBeGreaterThanOrEqual(1);
    });
  });

  describe('認証済み状態', () => {
    it('認証済みではメールアドレスとログアウトボタンが表示される', () => {
      (localStorage.getItem as jest.Mock).mockImplementation((key: string) => {
        if (key === 'auth_user') {
          return JSON.stringify({ userId: 'user-1', email: 'test@example.com' });
        }
        return null;
      });
      renderNavigation();
      const emails = screen.getAllByText('test@example.com');
      expect(emails.length).toBeGreaterThanOrEqual(1);
      const logoutButtons = screen.getAllByText('ログアウト');
      expect(logoutButtons.length).toBeGreaterThanOrEqual(1);
    });
  });

  describe('モバイルメニュー', () => {
    it('モバイルメニューボタンが表示される', () => {
      renderNavigation();
      expect(screen.getByLabelText('メニューを開く')).toBeInTheDocument();
    });
  });
});
