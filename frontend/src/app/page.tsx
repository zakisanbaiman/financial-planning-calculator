import Link from 'next/link';

export default function HomePage() {
  return (
    <div className="container mx-auto px-4 py-8">
      {/* Hero Section */}
      <div className="text-center mb-12">
        <h1 className="text-4xl font-bold text-gray-900 dark:text-white mb-4">
          財務計画計算機
        </h1>
        <p className="text-xl text-gray-600 dark:text-gray-300 mb-8 max-w-2xl mx-auto">
          将来の資産形成と老後の財務計画を可視化し、安心できる財務計画を立てましょう
        </p>
        <Link
          href="/dashboard"
          className="btn-primary inline-flex items-center space-x-2 text-lg px-8 py-3"
        >
          <span>📊</span>
          <span>ダッシュボードを開く</span>
        </Link>
      </div>

      {/* Features Grid */}
      <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6 mb-12">
        <div className="card">
          <div className="text-3xl mb-4">💰</div>
          <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-2">財務データ管理</h3>
          <p className="text-gray-600 dark:text-gray-300 mb-4">
            収入、支出、貯蓄状況を入力・管理して、正確な将来予測の基盤を作成
          </p>
          <Link href="/financial-data" className="text-primary-600 hover:text-primary-700 dark:text-primary-400 dark:hover:text-primary-300 font-medium">
            詳細を見る →
          </Link>
        </div>

        <div className="card">
          <div className="text-3xl mb-4">📈</div>
          <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-2">資産推移シミュレーション</h3>
          <p className="text-gray-600 dark:text-gray-300 mb-4">
            現在の貯蓄ペースで将来どれだけ資産が増えるかを可視化
          </p>
          <Link href="/calculations" className="text-primary-600 hover:text-primary-700 dark:text-primary-400 dark:hover:text-primary-300 font-medium">
            計算してみる →
          </Link>
        </div>

        <div className="card">
          <div className="text-3xl mb-4">🏖️</div>
          <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-2">老後資金計算</h3>
          <p className="text-gray-600 dark:text-gray-300 mb-4">
            退職後に必要な資金と年金額を考慮した老後の生活設計
          </p>
          <Link href="/calculations" className="text-primary-600 hover:text-primary-700 dark:text-primary-400 dark:hover:text-primary-300 font-medium">
            計算してみる →
          </Link>
        </div>

        <div className="card">
          <div className="text-3xl mb-4">🚨</div>
          <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-2">緊急資金計算</h3>
          <p className="text-gray-600 dark:text-gray-300 mb-4">
            万が一の時に必要な緊急資金を計算し、適切な備えを確認
          </p>
          <Link href="/calculations" className="text-primary-600 hover:text-primary-700 dark:text-primary-400 dark:hover:text-primary-300 font-medium">
            計算してみる →
          </Link>
        </div>

        <div className="card">
          <div className="text-3xl mb-4">🎯</div>
          <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-2">目標設定・進捗管理</h3>
          <p className="text-gray-600 dark:text-gray-300 mb-4">
            具体的な財務目標を設定し、進捗を追跡してモチベーションを維持
          </p>
          <Link href="/goals" className="text-primary-600 hover:text-primary-700 dark:text-primary-400 dark:hover:text-primary-300 font-medium">
            目標を設定 →
          </Link>
        </div>

        <div className="card">
          <div className="text-3xl mb-4">📋</div>
          <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-2">レポート生成</h3>
          <p className="text-gray-600 dark:text-gray-300 mb-4">
            財務状況をまとめたレポートをPDF形式で生成・印刷
          </p>
          <Link href="/reports" className="text-primary-600 hover:text-primary-700 dark:text-primary-400 dark:hover:text-primary-300 font-medium">
            レポート作成 →
          </Link>
        </div>
      </div>

      {/* Getting Started Section */}
      <div className="card max-w-2xl mx-auto text-center">
        <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-4">はじめ方</h2>
        <div className="space-y-4 text-left">
          <div className="flex items-start space-x-3">
            <div className="flex-shrink-0 w-6 h-6 bg-primary-500 text-white rounded-full flex items-center justify-center text-sm font-bold">1</div>
            <div>
              <h4 className="font-medium text-gray-900 dark:text-white">財務データを入力</h4>
              <p className="text-gray-600 dark:text-gray-300 text-sm">現在の収入、支出、貯蓄額を入力してください</p>
            </div>
          </div>
          <div className="flex items-start space-x-3">
            <div className="flex-shrink-0 w-6 h-6 bg-primary-500 text-white rounded-full flex items-center justify-center text-sm font-bold">2</div>
            <div>
              <h4 className="font-medium text-gray-900 dark:text-white">目標を設定</h4>
              <p className="text-gray-600 dark:text-gray-300 text-sm">達成したい財務目標を設定しましょう</p>
            </div>
          </div>
          <div className="flex items-start space-x-3">
            <div className="flex-shrink-0 w-6 h-6 bg-primary-500 text-white rounded-full flex items-center justify-center text-sm font-bold">3</div>
            <div>
              <h4 className="font-medium text-gray-900 dark:text-white">計算・可視化</h4>
              <p className="text-gray-600 dark:text-gray-300 text-sm">将来の資産推移や必要資金を確認してください</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}