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
      title: '財務プロフィール',
      description: '収入・支出・貯蓄を入力して、正確なシミュレーションの土台を構築します',
      href: '/financial-data',
      label: '詳細を見る',
    },
    {
      title: '資産推移予測',
      description: '現在の貯蓄率をもとに、将来の資産がどう増えるかを可視化します',
      href: '/calculations',
      label: '計算する',
    },
    {
      title: '老後計画',
      description: '年金受給額や想定生活費を考慮して、必要な老後資金を算出します',
      href: '/calculations',
      label: '計画する',
    },
    {
      title: '緊急資金',
      description: '万一の事態に備えた緊急資金の必要額を計算・確認します',
      href: '/calculations',
      label: '確認する',
    },
    {
      title: '目標管理',
      description: '具体的な財務目標を設定し、達成に向けた進捗をトラッキングします',
      href: '/goals',
      label: '目標を設定',
    },
    {
      title: 'レポート',
      description: '財務状況の総合レポートをPDF形式で生成・共有できます',
      href: '/reports',
      label: 'レポートを作成',
    },
  ];

  return (
    <div className="min-h-screen">
      {/* Hero Section */}
      <section className="container mx-auto px-4 pt-20 pb-24">
        <div className="max-w-3xl">
          <p className="font-body text-sm font-semibold tracking-editorial uppercase text-accent-600 dark:text-accent-400 mb-4">
            スマートな財務計画
          </p>
          <h1 className="font-display text-6xl md:text-7xl lg:text-8xl font-semibold text-ink-900 dark:text-ink-100 leading-[0.95] mb-8">
            FinPlan
          </h1>
          <p className="font-body text-xl md:text-2xl text-ink-500 dark:text-ink-400 leading-relaxed max-w-2xl mb-6 font-light">
            あなたの未来のための財務計画ツール
          </p>
          <p className="font-body text-base md:text-lg text-ink-400 dark:text-ink-500 mb-12 max-w-2xl">
            将来の資産を可視化し、老後に備え、自信を持って目標を達成しましょう
          </p>
          <div className="flex flex-col sm:flex-row items-start gap-4">
            {isAuthenticated ? (
              <Link
                href="/dashboard"
                className="btn-primary inline-flex items-center text-base px-8 py-3"
              >
                ダッシュボードを開く
              </Link>
            ) : (
              <>
                <button
                  onClick={handleGuestStart}
                  className="btn-primary inline-flex items-center text-base px-8 py-3"
                >
                  ゲストとして試す
                </button>
                <Link
                  href="/login"
                  className="btn-secondary inline-flex items-center text-base px-8 py-3"
                >
                  ログイン / 新規登録
                </Link>
              </>
            )}
          </div>
          {!isAuthenticated && (
            <p className="text-sm text-ink-400 dark:text-ink-500 mt-6 font-body">
              ゲストモード対応 ― 登録不要で全機能をお試しいただけます（データはローカルに保存）
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
            機能一覧
          </p>
          <h2 className="font-display text-4xl md:text-5xl font-semibold text-ink-900 dark:text-ink-100">
            できること
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
              はじめ方
            </p>
            <h2 className="font-display text-4xl md:text-5xl font-semibold text-ink-900 dark:text-ink-100">
              使い方
            </h2>
          </div>

          <div className="space-y-12">
            <div className="flex items-start gap-8">
              <span className="font-mono text-4xl font-medium text-ink-300 dark:text-ink-700 shrink-0 leading-none pt-1">
                01
              </span>
              <div className="border-t border-ink-200 dark:border-ink-800 pt-4 flex-1">
                <h3 className="font-display text-2xl font-semibold text-ink-900 dark:text-ink-100 mb-2">
                  財務データを入力する
                </h3>
                <p className="font-body text-ink-500 dark:text-ink-400 leading-relaxed">
                  現在の収入・支出・貯蓄額を入力してください
                </p>
              </div>
            </div>
            <div className="flex items-start gap-8">
              <span className="font-mono text-4xl font-medium text-ink-300 dark:text-ink-700 shrink-0 leading-none pt-1">
                02
              </span>
              <div className="border-t border-ink-200 dark:border-ink-800 pt-4 flex-1">
                <h3 className="font-display text-2xl font-semibold text-ink-900 dark:text-ink-100 mb-2">
                  目標を設定する
                </h3>
                <p className="font-body text-ink-500 dark:text-ink-400 leading-relaxed">
                  達成したい財務目標を定義してください
                </p>
              </div>
            </div>
            <div className="flex items-start gap-8">
              <span className="font-mono text-4xl font-medium text-ink-300 dark:text-ink-700 shrink-0 leading-none pt-1">
                03
              </span>
              <div className="border-t border-ink-200 dark:border-ink-800 pt-4 flex-1">
                <h3 className="font-display text-2xl font-semibold text-ink-900 dark:text-ink-100 mb-2">
                  計算して可視化する
                </h3>
                <p className="font-body text-ink-500 dark:text-ink-400 leading-relaxed">
                  将来の資産推移と必要資金をグラフで確認できます
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
