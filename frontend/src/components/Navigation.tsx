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
    { href: '/', label: 'Home', icon: '🏠' },
    { href: '/dashboard', label: 'Dashboard', icon: '📊' },
    { href: '/financial-data', label: 'Profile', icon: '💼' },
    { href: '/goals', label: 'Goals', icon: '🎯' },
    { href: '/calculations', label: 'Calculator', icon: '🧮' },
    { href: '/reports', label: 'Reports', icon: '📋' },
  ];

  return (
    <nav className="bg-white dark:bg-gray-800 shadow-sm border-b border-gray-200 dark:border-gray-700 transition-colors">
      <div className="container mx-auto px-4">
        <div className="flex items-center justify-between h-16">
          {/* Logo */}
          <Link href="/" className="flex items-center space-x-2 group">
            <div className="w-10 h-10 bg-gradient-to-br from-primary-500 to-primary-600 rounded-xl flex items-center justify-center shadow-lg shadow-primary-500/30 group-hover:shadow-xl transition-all duration-200 group-hover:scale-105">
              <span className="text-xl">📊</span>
            </div>
            <span className="text-xl font-bold bg-gradient-to-r from-primary-600 to-primary-800 bg-clip-text text-transparent dark:from-primary-400 dark:to-primary-600">FinPlan</span>
          </Link>

          {/* Navigation Links */}
          <div className="hidden md:flex items-center space-x-1">
            {navItems.map((item) => (
              <Link
                key={item.href}
                href={item.href}
                className={`flex items-center space-x-1 px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
                  pathname === item.href
                    ? 'bg-primary-50 text-primary-700 dark:bg-primary-900 dark:text-primary-300'
                    : 'text-gray-600 hover:text-gray-900 hover:bg-gray-50 dark:text-gray-300 dark:hover:text-white dark:hover:bg-gray-700'
                }`}
              >
                <span>{item.icon}</span>
                <span>{item.label}</span>
              </Link>
            ))}
            <button
              onClick={startTutorial}
              className="flex items-center space-x-1 px-3 py-2 rounded-lg text-sm font-medium transition-colors text-gray-600 hover:text-gray-900 hover:bg-gray-50 dark:text-gray-300 dark:hover:text-white dark:hover:bg-gray-700"
              title="チュートリアルを表示"
            >
              <span>🎓</span>
              <span>ヘルプ</span>
            </button>
            <ThemeToggle />
            
            {/* 認証状態に応じた表示 */}
            {isAuthenticated ? (
              <div className="flex items-center space-x-2 ml-2 pl-2 border-l border-gray-300 dark:border-gray-600">
                <span className="text-sm text-gray-600 dark:text-gray-300">
                  {user?.email}
                </span>
                <button
                  onClick={logout}
                  className="flex items-center space-x-1 px-3 py-2 rounded-lg text-sm font-medium transition-colors text-gray-600 hover:text-gray-900 hover:bg-gray-50 dark:text-gray-300 dark:hover:text-white dark:hover:bg-gray-700"
                  title="ログアウト"
                >
                  <span>🚪</span>
                  <span>ログアウト</span>
                </button>
              </div>
            ) : isGuestMode ? (
              <div className="flex items-center space-x-2 ml-2 pl-2 border-l border-gray-300 dark:border-gray-600">
                <span className="text-sm text-warning-600 dark:text-warning-400 flex items-center space-x-1">
                  <span>✨</span>
                  <span>ゲストモード</span>
                </span>
                <Link
                  href="/register"
                  className="flex items-center space-x-1 px-3 py-2 rounded-lg text-sm font-medium transition-colors bg-primary-600 text-white hover:bg-primary-700"
                  title="データを保存するには登録が必要です"
                >
                  <span>💾</span>
                  <span>登録してデータを保存</span>
                </Link>
              </div>
            ) : (
              <div className="flex items-center space-x-2 ml-2 pl-2 border-l border-gray-300 dark:border-gray-600">
                <Link
                  href="/login"
                  className="flex items-center space-x-1 px-3 py-2 rounded-lg text-sm font-medium transition-colors text-gray-600 hover:text-gray-900 hover:bg-gray-50 dark:text-gray-300 dark:hover:text-white dark:hover:bg-gray-700"
                >
                  <span>🔐</span>
                  <span>ログイン</span>
                </Link>
              </div>
            )}
          </div>

          {/* Mobile Menu Button */}
          <div className="md:hidden flex items-center space-x-2">
            <ThemeToggle />
            <button
              type="button"
              className="text-gray-600 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white focus:outline-none focus:text-gray-900 dark:focus:text-white"
              aria-label="メニューを開く"
            >
              <svg className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
              </svg>
            </button>
          </div>
        </div>

        {/* Mobile Navigation (hidden by default, would need state management for toggle) */}
        <div className="md:hidden border-t border-gray-200 dark:border-gray-700 py-2">
          {navItems.map((item) => (
            <Link
              key={item.href}
              href={item.href}
              className={`flex items-center space-x-2 px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
                pathname === item.href
                  ? 'bg-primary-50 text-primary-700 dark:bg-primary-900 dark:text-primary-300'
                  : 'text-gray-600 hover:text-gray-900 hover:bg-gray-50 dark:text-gray-300 dark:hover:text-white dark:hover:bg-gray-700'
              }`}
            >
              <span>{item.icon}</span>
              <span>{item.label}</span>
            </Link>
          ))}
          <button
            onClick={startTutorial}
            className="flex items-center space-x-2 px-3 py-2 rounded-lg text-sm font-medium transition-colors text-gray-600 hover:text-gray-900 hover:bg-gray-50 dark:text-gray-300 dark:hover:text-white dark:hover:bg-gray-700 w-full"
          >
            <span>🎓</span>
            <span>ヘルプ</span>
          </button>
          
          {/* モバイル認証メニュー */}
          {isAuthenticated ? (
            <>
              <div className="px-3 py-2 text-sm text-gray-600 dark:text-gray-300 border-t border-gray-200 dark:border-gray-700 mt-2 pt-2">
                {user?.email}
              </div>
              <button
                onClick={logout}
                className="flex items-center space-x-2 px-3 py-2 rounded-lg text-sm font-medium transition-colors text-gray-600 hover:text-gray-900 hover:bg-gray-50 dark:text-gray-300 dark:hover:text-white dark:hover:bg-gray-700 w-full"
              >
                <span>🚪</span>
                <span>ログアウト</span>
              </button>
            </>
          ) : isGuestMode ? (
            <>
              <div className="px-3 py-2 text-sm text-warning-600 dark:text-warning-400 border-t border-gray-200 dark:border-gray-700 mt-2 pt-2 flex items-center space-x-1">
                <span>✨</span>
                <span>ゲストモード</span>
              </div>
              <Link
                href="/register"
                className="flex items-center space-x-2 px-3 py-2 rounded-lg text-sm font-medium transition-colors bg-primary-600 text-white hover:bg-primary-700"
              >
                <span>💾</span>
                <span>登録してデータを保存</span>
              </Link>
            </>
          ) : (
            <Link
              href="/login"
              className="flex items-center space-x-2 px-3 py-2 rounded-lg text-sm font-medium transition-colors text-gray-600 hover:text-gray-900 hover:bg-gray-50 dark:text-gray-300 dark:hover:text-white dark:hover:bg-gray-700 border-t border-gray-200 dark:border-gray-700 mt-2 pt-2"
            >
              <span>🔐</span>
              <span>ログイン</span>
            </Link>
          )}
        </div>
      </div>
    </nav>
  );
};

export default Navigation;