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

  const features = [
    {
      title: 'Financial Profile',
      description: 'Manage your income, expenses, and savings to build a solid foundation for accurate projections',
      href: '/financial-data',
      label: 'Learn more',
    },
    {
      title: 'Asset Projection',
      description: 'Visualize how your assets will grow over time based on your current savings rate',
      href: '/calculations',
      label: 'Calculate now',
    },
    {
      title: 'Retirement Planning',
      description: 'Calculate required retirement funds considering pension and expected lifestyle costs',
      href: '/calculations',
      label: 'Plan now',
    },
    {
      title: 'Emergency Fund',
      description: 'Calculate and verify the emergency funds needed for unexpected situations',
      href: '/calculations',
      label: 'Check now',
    },
    {
      title: 'Goal Tracking',
      description: 'Set specific financial goals and track progress to stay motivated',
      href: '/goals',
      label: 'Set goals',
    },
    {
      title: 'Reports',
      description: 'Generate comprehensive financial reports in PDF format for easy sharing',
      href: '/reports',
      label: 'Create report',
    },
  ];

  return (
    <div className="min-h-screen">
      {/* Hero Section */}
      <section className="container mx-auto px-4 pt-20 pb-24">
        <div className="max-w-3xl">
          <p className="font-body text-sm font-semibold tracking-editorial uppercase text-accent-600 dark:text-accent-400 mb-4">
            Smart Financial Planning
          </p>
          <h1 className="font-display text-6xl md:text-7xl lg:text-8xl font-semibold text-ink-900 dark:text-ink-100 leading-[0.95] mb-8">
            FinPlan
          </h1>
          <p className="font-body text-xl md:text-2xl text-ink-500 dark:text-ink-400 leading-relaxed max-w-2xl mb-6 font-light">
            Smart Financial Planning for Your Future
          </p>
          <p className="font-body text-base md:text-lg text-ink-400 dark:text-ink-500 mb-12 max-w-2xl">
            Visualize your financial future, plan for retirement, and achieve your goals with confidence
          </p>
          <div className="flex flex-col sm:flex-row items-start gap-4">
            {isAuthenticated ? (
              <Link
                href="/dashboard"
                className="btn-primary inline-flex items-center text-base px-8 py-3"
              >
                Open Dashboard
              </Link>
            ) : (
              <>
                <button
                  onClick={handleGuestStart}
                  className="btn-primary inline-flex items-center text-base px-8 py-3"
                >
                  Try as Guest
                </button>
                <Link
                  href="/login"
                  className="btn-secondary inline-flex items-center text-base px-8 py-3"
                >
                  Sign In / Sign Up
                </Link>
              </>
            )}
          </div>
          {!isAuthenticated && (
            <p className="text-sm text-ink-400 dark:text-ink-500 mt-6 font-body">
              Guest mode available - Try all features without registration (data saved locally)
            </p>
          )}
        </div>
      </section>

      {/* Divider */}
      <div className="container mx-auto px-4">
        <hr className="border-ink-200 dark:border-ink-800" />
      </div>

      {/* Features Grid */}
      <section className="container mx-auto px-4 py-20">
        <div className="mb-12">
          <p className="font-body text-sm font-semibold tracking-editorial uppercase text-accent-600 dark:text-accent-400 mb-3">
            Features
          </p>
          <h2 className="font-display text-4xl md:text-5xl font-semibold text-ink-900 dark:text-ink-100">
            What You Can Do
          </h2>
        </div>

        <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-px bg-ink-200 dark:bg-ink-800 border border-ink-200 dark:border-ink-800">
          {features.map((feature, index) => (
            <div
              key={index}
              className="bg-ink-50 dark:bg-ink-950 p-8 group hover:bg-white dark:hover:bg-ink-900 transition-colors duration-200"
            >
              <h3 className="font-display text-2xl font-semibold text-ink-900 dark:text-ink-100 mb-3">
                {feature.title}
              </h3>
              <p className="font-body text-sm text-ink-500 dark:text-ink-400 leading-relaxed mb-6">
                {feature.description}
              </p>
              <Link
                href={feature.href}
                className="font-body text-sm font-semibold text-ink-900 dark:text-ink-100 inline-flex items-center group/link"
              >
                <span>{feature.label}</span>
                <span className="ml-2 transform group-hover/link:translate-x-1 transition-transform duration-150">&rarr;</span>
              </Link>
            </div>
          ))}
        </div>
      </section>

      {/* Divider */}
      <div className="container mx-auto px-4">
        <hr className="border-ink-200 dark:border-ink-800" />
      </div>

      {/* Getting Started Section */}
      <section className="container mx-auto px-4 py-20">
        <div className="max-w-2xl mx-auto">
          <div className="text-center mb-16">
            <p className="font-body text-sm font-semibold tracking-editorial uppercase text-accent-600 dark:text-accent-400 mb-3">
              Getting Started
            </p>
            <h2 className="font-display text-4xl md:text-5xl font-semibold text-ink-900 dark:text-ink-100">
              How It Works
            </h2>
          </div>

          <div className="space-y-12">
            <div className="flex items-start gap-8">
              <span className="font-mono text-4xl font-medium text-ink-300 dark:text-ink-700 shrink-0 leading-none pt-1">
                01
              </span>
              <div className="border-t border-ink-200 dark:border-ink-800 pt-4 flex-1">
                <h3 className="font-display text-2xl font-semibold text-ink-900 dark:text-ink-100 mb-2">
                  Enter Your Financial Data
                </h3>
                <p className="font-body text-ink-500 dark:text-ink-400 leading-relaxed">
                  Input your current income, expenses, and savings amount
                </p>
              </div>
            </div>
            <div className="flex items-start gap-8">
              <span className="font-mono text-4xl font-medium text-ink-300 dark:text-ink-700 shrink-0 leading-none pt-1">
                02
              </span>
              <div className="border-t border-ink-200 dark:border-ink-800 pt-4 flex-1">
                <h3 className="font-display text-2xl font-semibold text-ink-900 dark:text-ink-100 mb-2">
                  Set Your Goals
                </h3>
                <p className="font-body text-ink-500 dark:text-ink-400 leading-relaxed">
                  Define the financial goals you want to achieve
                </p>
              </div>
            </div>
            <div className="flex items-start gap-8">
              <span className="font-mono text-4xl font-medium text-ink-300 dark:text-ink-700 shrink-0 leading-none pt-1">
                03
              </span>
              <div className="border-t border-ink-200 dark:border-ink-800 pt-4 flex-1">
                <h3 className="font-display text-2xl font-semibold text-ink-900 dark:text-ink-100 mb-2">
                  Calculate & Visualize
                </h3>
                <p className="font-body text-ink-500 dark:text-ink-400 leading-relaxed">
                  Review your future asset projections and required funds
                </p>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Footer spacer */}
      <div className="h-20" />
    </div>
  );
}
