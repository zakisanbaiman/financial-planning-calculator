'use client';

import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useTutorial } from '@/lib/contexts/TutorialContext';
import { useGuestMode } from '@/lib/contexts/GuestModeContext';
import { useAuth } from '@/lib/contexts/AuthContext';

export default function HomePage() {
  const router = useRouter();
  const { startTutorial } = useTutorial();
  const { startGuestMode } = useGuestMode();
  const { isAuthenticated } = useAuth();

  const handleGuestStart = () => {
    startGuestMode();
    router.push('/dashboard');
  };

  return (
    <div className="container mx-auto px-4 py-8 gradient-bg min-h-screen">
      {/* Hero Section */}
      <div className="text-center mb-16 py-12">
        <div className="inline-block mb-6">
          <div className="w-20 h-20 bg-gradient-to-br from-primary-500 to-primary-600 rounded-3xl flex items-center justify-center shadow-2xl shadow-primary-500/40 mx-auto">
            <span className="text-5xl">📊</span>
          </div>
        </div>
        <h1 className="text-5xl md:text-6xl font-bold mb-4 bg-gradient-to-r from-primary-600 via-primary-700 to-primary-800 bg-clip-text text-transparent dark:from-primary-400 dark:via-primary-500 dark:to-primary-600">
          FinPlan
        </h1>
        <p className="text-xl md:text-2xl text-gray-600 dark:text-gray-300 mb-8 max-w-3xl mx-auto font-light">
          Smart Financial Planning for Your Future
        </p>
        <p className="text-base md:text-lg text-gray-500 dark:text-gray-400 mb-10 max-w-2xl mx-auto">
          Visualize your financial future, plan for retirement, and achieve your goals with confidence
        </p>
        <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
          {isAuthenticated ? (
            <Link
              href="/dashboard"
              className="btn-primary inline-flex items-center space-x-2 text-lg px-10 py-4"
            >
              <span>📊</span>
              <span>Open Dashboard</span>
            </Link>
          ) : (
            <>
              <button
                onClick={handleGuestStart}
                className="btn-primary inline-flex items-center space-x-2 text-lg px-10 py-4"
              >
                <span>✨</span>
                <span>Try as Guest</span>
              </button>
              <Link
                href="/login"
                className="btn-secondary inline-flex items-center space-x-2 text-lg px-10 py-4"
              >
                <span>🔐</span>
                <span>Sign In / Sign Up</span>
              </Link>
            </>
          )}
        </div>
        {!isAuthenticated && (
          <p className="text-sm text-gray-500 dark:text-gray-400 mt-6 flex items-center justify-center gap-2">
            <span>💡</span>
            <span>Guest mode available - Try all features without registration (data saved locally)</span>
          </p>
        )}
      </div>

      {/* Features Grid */}
      <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6 mb-16 max-w-7xl mx-auto">
        <div className="card hover:scale-105 transition-transform duration-200">
          <div className="w-14 h-14 bg-gradient-to-br from-primary-500 to-primary-600 rounded-2xl flex items-center justify-center mb-4 shadow-lg shadow-primary-500/30">
            <span className="text-3xl">💼</span>
          </div>
          <h3 className="text-lg font-bold text-gray-900 dark:text-white mb-2">Financial Profile</h3>
          <p className="text-gray-600 dark:text-gray-300 mb-4 text-sm leading-relaxed">
            Manage your income, expenses, and savings to build a solid foundation for accurate projections
          </p>
          <Link href="/financial-data" className="text-primary-600 hover:text-primary-700 dark:text-primary-400 dark:hover:text-primary-300 font-semibold text-sm inline-flex items-center group">
            <span>Learn more</span>
            <span className="ml-1 transform group-hover:translate-x-1 transition-transform">→</span>
          </Link>
        </div>

        <div className="card hover:scale-105 transition-transform duration-200">
          <div className="w-14 h-14 bg-gradient-to-br from-success-500 to-success-600 rounded-2xl flex items-center justify-center mb-4 shadow-lg shadow-success-500/30">
            <span className="text-3xl">📈</span>
          </div>
          <h3 className="text-lg font-bold text-gray-900 dark:text-white mb-2">Asset Projection</h3>
          <p className="text-gray-600 dark:text-gray-300 mb-4 text-sm leading-relaxed">
            Visualize how your assets will grow over time based on your current savings rate
          </p>
          <Link href="/calculations" className="text-primary-600 hover:text-primary-700 dark:text-primary-400 dark:hover:text-primary-300 font-semibold text-sm inline-flex items-center group">
            <span>Calculate now</span>
            <span className="ml-1 transform group-hover:translate-x-1 transition-transform">→</span>
          </Link>
        </div>

        <div className="card hover:scale-105 transition-transform duration-200">
          <div className="w-14 h-14 bg-gradient-to-br from-warning-500 to-warning-600 rounded-2xl flex items-center justify-center mb-4 shadow-lg shadow-warning-500/30">
            <span className="text-3xl">🏖️</span>
          </div>
          <h3 className="text-lg font-bold text-gray-900 dark:text-white mb-2">Retirement Planning</h3>
          <p className="text-gray-600 dark:text-gray-300 mb-4 text-sm leading-relaxed">
            Calculate required retirement funds considering pension and expected lifestyle costs
          </p>
          <Link href="/calculations" className="text-primary-600 hover:text-primary-700 dark:text-primary-400 dark:hover:text-primary-300 font-semibold text-sm inline-flex items-center group">
            <span>Plan now</span>
            <span className="ml-1 transform group-hover:translate-x-1 transition-transform">→</span>
          </Link>
        </div>

        <div className="card hover:scale-105 transition-transform duration-200">
          <div className="w-14 h-14 bg-gradient-to-br from-error-500 to-error-600 rounded-2xl flex items-center justify-center mb-4 shadow-lg shadow-error-500/30">
            <span className="text-3xl">🚨</span>
          </div>
          <h3 className="text-lg font-bold text-gray-900 dark:text-white mb-2">Emergency Fund</h3>
          <p className="text-gray-600 dark:text-gray-300 mb-4 text-sm leading-relaxed">
            Calculate and verify the emergency funds needed for unexpected situations
          </p>
          <Link href="/calculations" className="text-primary-600 hover:text-primary-700 dark:text-primary-400 dark:hover:text-primary-300 font-semibold text-sm inline-flex items-center group">
            <span>Check now</span>
            <span className="ml-1 transform group-hover:translate-x-1 transition-transform">→</span>
          </Link>
        </div>

        <div className="card hover:scale-105 transition-transform duration-200">
          <div className="w-14 h-14 bg-gradient-to-br from-primary-400 to-primary-500 rounded-2xl flex items-center justify-center mb-4 shadow-lg shadow-primary-400/30">
            <span className="text-3xl">🎯</span>
          </div>
          <h3 className="text-lg font-bold text-gray-900 dark:text-white mb-2">Goal Tracking</h3>
          <p className="text-gray-600 dark:text-gray-300 mb-4 text-sm leading-relaxed">
            Set specific financial goals and track progress to stay motivated
          </p>
          <Link href="/goals" className="text-primary-600 hover:text-primary-700 dark:text-primary-400 dark:hover:text-primary-300 font-semibold text-sm inline-flex items-center group">
            <span>Set goals</span>
            <span className="ml-1 transform group-hover:translate-x-1 transition-transform">→</span>
          </Link>
        </div>

        <div className="card hover:scale-105 transition-transform duration-200">
          <div className="w-14 h-14 bg-gradient-to-br from-gray-600 to-gray-700 rounded-2xl flex items-center justify-center mb-4 shadow-lg shadow-gray-600/30">
            <span className="text-3xl">📋</span>
          </div>
          <h3 className="text-lg font-bold text-gray-900 dark:text-white mb-2">Reports</h3>
          <p className="text-gray-600 dark:text-gray-300 mb-4 text-sm leading-relaxed">
            Generate comprehensive financial reports in PDF format for easy sharing
          </p>
          <Link href="/reports" className="text-primary-600 hover:text-primary-700 dark:text-primary-400 dark:hover:text-primary-300 font-semibold text-sm inline-flex items-center group">
            <span>Create report</span>
            <span className="ml-1 transform group-hover:translate-x-1 transition-transform">→</span>
          </Link>
        </div>
      </div>

      {/* Getting Started Section */}
      <div className="card max-w-3xl mx-auto">
        <h2 className="text-3xl font-bold text-gray-900 dark:text-white mb-8 text-center bg-gradient-to-r from-primary-600 to-primary-800 bg-clip-text text-transparent dark:from-primary-400 dark:to-primary-600">Getting Started</h2>
        <div className="space-y-6">
          <div className="flex items-start space-x-4">
            <div className="flex-shrink-0 w-12 h-12 bg-gradient-to-br from-primary-500 to-primary-600 text-white rounded-2xl flex items-center justify-center text-lg font-bold shadow-lg shadow-primary-500/30">1</div>
            <div className="flex-1">
              <h4 className="font-bold text-gray-900 dark:text-white text-lg mb-1">Enter Your Financial Data</h4>
              <p className="text-gray-600 dark:text-gray-300 leading-relaxed">Input your current income, expenses, and savings amount</p>
            </div>
          </div>
          <div className="flex items-start space-x-4">
            <div className="flex-shrink-0 w-12 h-12 bg-gradient-to-br from-primary-500 to-primary-600 text-white rounded-2xl flex items-center justify-center text-lg font-bold shadow-lg shadow-primary-500/30">2</div>
            <div className="flex-1">
              <h4 className="font-bold text-gray-900 dark:text-white text-lg mb-1">Set Your Goals</h4>
              <p className="text-gray-600 dark:text-gray-300 leading-relaxed">Define the financial goals you want to achieve</p>
            </div>
          </div>
          <div className="flex items-start space-x-4">
            <div className="flex-shrink-0 w-12 h-12 bg-gradient-to-br from-primary-500 to-primary-600 text-white rounded-2xl flex items-center justify-center text-lg font-bold shadow-lg shadow-primary-500/30">3</div>
            <div className="flex-1">
              <h4 className="font-bold text-gray-900 dark:text-white text-lg mb-1">Calculate & Visualize</h4>
              <p className="text-gray-600 dark:text-gray-300 leading-relaxed">Review your future asset projections and required funds</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}