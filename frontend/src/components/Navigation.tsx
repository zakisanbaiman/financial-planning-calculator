'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { useTutorial } from '@/lib/contexts/TutorialContext';
import { useAuth } from '@/lib/contexts/AuthContext';
import { useGuestMode } from '@/lib/contexts/GuestModeContext';
import ThemeToggle from './ThemeToggle';

const Navigation = () => {
  const pathname = usePathname();
  const { startTutorial } = useTutorial();
  const { isAuthenticated, user, logout } = useAuth();
  const { isGuestMode, exitGuestMode } = useGuestMode();

  const navItems = [
    { href: '/', label: 'Home' },
    { href: '/dashboard', label: 'Dashboard' },
    { href: '/financial-data', label: 'Profile' },
    { href: '/goals', label: 'Goals' },
    { href: '/calculations', label: 'Calculator' },
    { href: '/reports', label: 'Reports' },
    ...(isAuthenticated ? [{ href: '/bot', label: 'Bot' }] : []),
  ];

  return (
    <nav className="border-b border-ink-200 dark:border-ink-800 bg-ink-50 dark:bg-ink-950 transition-colors">
      <div className="container mx-auto px-4">
        <div className="flex items-center justify-between h-14">
          {/* Logo */}
          <Link href="/" className="flex items-center space-x-2 group">
            <span className="font-display text-2xl font-semibold text-ink-900 dark:text-ink-100 tracking-tight">
              FinPlan
            </span>
          </Link>

          {/* Navigation Links */}
          <div className="hidden md:flex items-center space-x-1">
            {navItems.map((item) => (
              <Link
                key={item.href}
                href={item.href}
                className={`px-3 py-1.5 text-sm font-body font-medium transition-colors ${
                  pathname === item.href
                    ? 'text-ink-900 dark:text-ink-100 border-b-2 border-ink-900 dark:border-ink-100'
                    : 'text-ink-500 hover:text-ink-800 dark:text-ink-400 dark:hover:text-ink-200'
                }`}
              >
                {item.label}
              </Link>
            ))}
            <button
              onClick={startTutorial}
              className="px-3 py-1.5 text-sm font-body font-medium transition-colors text-ink-500 hover:text-ink-800 dark:text-ink-400 dark:hover:text-ink-200"
              title="チュートリアルを表示"
            >
              ヘルプ
            </button>
            <ThemeToggle />

            {/* 認証状態に応じた表示 */}
            {isAuthenticated ? (
              <div className="flex items-center space-x-2 ml-3 pl-3 border-l border-ink-200 dark:border-ink-700">
                <span className="text-sm text-ink-500 dark:text-ink-400 font-body">
                  {user?.email}
                </span>
                <Link
                  href="/settings/security"
                  className={`px-3 py-1.5 text-sm font-body font-medium transition-colors ${
                    pathname === '/settings/security'
                      ? 'text-ink-900 dark:text-ink-100'
                      : 'text-ink-500 hover:text-ink-800 dark:text-ink-400 dark:hover:text-ink-200'
                  }`}
                  title="セキュリティ設定"
                >
                  セキュリティ
                </Link>
                <button
                  onClick={logout}
                  className="px-3 py-1.5 text-sm font-body font-medium transition-colors text-ink-500 hover:text-ink-800 dark:text-ink-400 dark:hover:text-ink-200"
                  title="ログアウト"
                >
                  ログアウト
                </button>
              </div>
            ) : isGuestMode ? (
              <div className="flex items-center space-x-2 ml-3 pl-3 border-l border-ink-200 dark:border-ink-700">
                <span className="text-sm text-accent-600 dark:text-accent-400 font-body">
                  ゲストモード
                </span>
                <Link
                  href="/register"
                  className="px-3 py-1.5 text-sm font-body font-semibold text-ink-900 dark:text-ink-100 border border-ink-900 dark:border-ink-100 hover:bg-ink-900 hover:text-ink-50 dark:hover:bg-ink-100 dark:hover:text-ink-900 transition-colors"
                  title="データを保存するには登録が必要です"
                >
                  登録してデータを保存
                </Link>
              </div>
            ) : (
              <div className="flex items-center space-x-2 ml-3 pl-3 border-l border-ink-200 dark:border-ink-700">
                <Link
                  href="/login"
                  className="px-3 py-1.5 text-sm font-body font-medium transition-colors text-ink-500 hover:text-ink-800 dark:text-ink-400 dark:hover:text-ink-200"
                >
                  ログイン
                </Link>
              </div>
            )}
          </div>

          {/* Mobile Menu Button */}
          <div className="md:hidden flex items-center space-x-2">
            <ThemeToggle />
            <button
              type="button"
              className="text-ink-600 dark:text-ink-300 hover:text-ink-900 dark:hover:text-ink-100 focus:outline-none"
              aria-label="メニューを開く"
            >
              <svg className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
                <path strokeLinecap="round" strokeLinejoin="round" d="M4 6h16M4 12h16M4 18h16" />
              </svg>
            </button>
          </div>
        </div>

        {/* Mobile Navigation */}
        <div className="md:hidden border-t border-ink-200 dark:border-ink-800 py-2">
          {navItems.map((item) => (
            <Link
              key={item.href}
              href={item.href}
              className={`block px-3 py-2 text-sm font-body font-medium transition-colors ${
                pathname === item.href
                  ? 'text-ink-900 dark:text-ink-100'
                  : 'text-ink-500 hover:text-ink-800 dark:text-ink-400 dark:hover:text-ink-200'
              }`}
            >
              {item.label}
            </Link>
          ))}
          <button
            onClick={startTutorial}
            className="block w-full text-left px-3 py-2 text-sm font-body font-medium transition-colors text-ink-500 hover:text-ink-800 dark:text-ink-400 dark:hover:text-ink-200"
          >
            ヘルプ
          </button>

          {/* モバイル認証メニュー */}
          {isAuthenticated ? (
            <>
              <div className="px-3 py-2 text-sm text-ink-500 dark:text-ink-400 border-t border-ink-200 dark:border-ink-800 mt-2 pt-2 font-body">
                {user?.email}
              </div>
              <Link
                href="/settings/security"
                className={`block px-3 py-2 text-sm font-body font-medium transition-colors ${
                  pathname === '/settings/security'
                    ? 'text-ink-900 dark:text-ink-100'
                    : 'text-ink-500 hover:text-ink-800 dark:text-ink-400 dark:hover:text-ink-200'
                }`}
              >
                セキュリティ設定
              </Link>
              <button
                onClick={logout}
                className="block w-full text-left px-3 py-2 text-sm font-body font-medium transition-colors text-ink-500 hover:text-ink-800 dark:text-ink-400 dark:hover:text-ink-200"
              >
                ログアウト
              </button>
            </>
          ) : isGuestMode ? (
            <>
              <div className="px-3 py-2 text-sm text-accent-600 dark:text-accent-400 border-t border-ink-200 dark:border-ink-800 mt-2 pt-2 font-body">
                ゲストモード
              </div>
              <Link
                href="/register"
                className="block mx-3 mt-1 px-3 py-2 text-sm font-body font-semibold text-center text-ink-900 dark:text-ink-100 border border-ink-900 dark:border-ink-100 hover:bg-ink-900 hover:text-ink-50 dark:hover:bg-ink-100 dark:hover:text-ink-900 transition-colors"
              >
                登録してデータを保存
              </Link>
            </>
          ) : (
            <Link
              href="/login"
              className="block px-3 py-2 text-sm font-body font-medium transition-colors text-ink-500 hover:text-ink-800 dark:text-ink-400 dark:hover:text-ink-200 border-t border-ink-200 dark:border-ink-800 mt-2 pt-2"
            >
              ログイン
            </Link>
          )}
        </div>
      </div>
    </nav>
  );
};

export default Navigation;
