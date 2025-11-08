import Link from 'next/link';

export default function DashboardPage() {
  return (
    <div className="container mx-auto px-4 py-8">
      {/* Header */}
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">ダッシュボード</h1>
        <p className="text-gray-600">財務状況の概要と主要な指標を確認できます</p>
      </div>

      {/* Quick Stats */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <div className="card">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">月間純貯蓄</p>
              <p className="text-2xl font-bold text-gray-900">¥120,000</p>
            </div>
            <div className="text-2xl">💰</div>
          </div>
          <p className="text-xs text-gray-500 mt-2">前月比 +5%</p>
        </div>

        <div className="card">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">総資産</p>
              <p className="text-2xl font-bold text-gray-900">¥1,500,000</p>
            </div>
            <div className="text-2xl">📈</div>
          </div>
          <p className="text-xs text-gray-500 mt-2">前月比 +8%</p>
        </div>

        <div className="card">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">老後資金充足率</p>
              <p className="text-2xl font-bold text-gray-900">65%</p>
            </div>
            <div className="text-2xl">🏖️</div>
          </div>
          <p className="text-xs text-gray-500 mt-2">目標まで35%</p>
        </div>

        <div className="card">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">緊急資金</p>
              <p className="text-2xl font-bold text-gray-900">6ヶ月分</p>
            </div>
            <div className="text-2xl">🚨</div>
          </div>
          <p className="text-xs text-success-600 mt-2">十分確保済み</p>
        </div>
      </div>

      {/* Main Content Grid */}
      <div className="grid lg:grid-cols-3 gap-8">
        {/* Left Column - Charts and Projections */}
        <div className="lg:col-span-2 space-y-6">
          {/* Asset Projection Chart Placeholder */}
          <div className="card">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-xl font-semibold text-gray-900">資産推移予測</h2>
              <Link href="/calculations" className="text-primary-600 hover:text-primary-700 text-sm font-medium">
                詳細計算 →
              </Link>
            </div>
            <div className="h-64 bg-gray-100 rounded-lg flex items-center justify-center">
              <div className="text-center text-gray-500">
                <div className="text-4xl mb-2">📊</div>
                <p>資産推移グラフ</p>
                <p className="text-sm">(Chart.js実装予定)</p>
              </div>
            </div>
          </div>

          {/* Monthly Breakdown */}
          <div className="card">
            <h2 className="text-xl font-semibold text-gray-900 mb-4">月間収支内訳</h2>
            <div className="space-y-3">
              <div className="flex items-center justify-between py-2 border-b border-gray-100">
                <span className="text-gray-600">月収</span>
                <span className="font-medium text-gray-900">¥400,000</span>
              </div>
              <div className="flex items-center justify-between py-2 border-b border-gray-100">
                <span className="text-gray-600">住居費</span>
                <span className="font-medium text-gray-900">¥120,000</span>
              </div>
              <div className="flex items-center justify-between py-2 border-b border-gray-100">
                <span className="text-gray-600">食費</span>
                <span className="font-medium text-gray-900">¥60,000</span>
              </div>
              <div className="flex items-center justify-between py-2 border-b border-gray-100">
                <span className="text-gray-600">その他支出</span>
                <span className="font-medium text-gray-900">¥100,000</span>
              </div>
              <div className="flex items-center justify-between py-2 font-semibold">
                <span className="text-gray-900">純貯蓄</span>
                <span className="text-success-600">¥120,000</span>
              </div>
            </div>
          </div>
        </div>

        {/* Right Column - Goals and Actions */}
        <div className="space-y-6">
          {/* Active Goals */}
          <div className="card">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-xl font-semibold text-gray-900">進行中の目標</h2>
              <Link href="/goals" className="text-primary-600 hover:text-primary-700 text-sm font-medium">
                管理 →
              </Link>
            </div>
            <div className="space-y-4">
              <div>
                <div className="flex items-center justify-between mb-2">
                  <span className="text-sm font-medium text-gray-900">緊急資金</span>
                  <span className="text-sm text-gray-600">100%</span>
                </div>
                <div className="w-full bg-gray-200 rounded-full h-2">
                  <div className="bg-success-500 h-2 rounded-full" style={{ width: '100%' }}></div>
                </div>
                <p className="text-xs text-gray-500 mt-1">¥600,000 / ¥600,000</p>
              </div>

              <div>
                <div className="flex items-center justify-between mb-2">
                  <span className="text-sm font-medium text-gray-900">老後資金</span>
                  <span className="text-sm text-gray-600">65%</span>
                </div>
                <div className="w-full bg-gray-200 rounded-full h-2">
                  <div className="bg-primary-500 h-2 rounded-full" style={{ width: '65%' }}></div>
                </div>
                <p className="text-xs text-gray-500 mt-1">¥19,500,000 / ¥30,000,000</p>
              </div>

              <div>
                <div className="flex items-center justify-between mb-2">
                  <span className="text-sm font-medium text-gray-900">マイホーム資金</span>
                  <span className="text-sm text-gray-600">25%</span>
                </div>
                <div className="w-full bg-gray-200 rounded-full h-2">
                  <div className="bg-warning-500 h-2 rounded-full" style={{ width: '25%' }}></div>
                </div>
                <p className="text-xs text-gray-500 mt-1">¥1,250,000 / ¥5,000,000</p>
              </div>
            </div>
          </div>

          {/* Quick Actions */}
          <div className="card">
            <h2 className="text-xl font-semibold text-gray-900 mb-4">クイックアクション</h2>
            <div className="space-y-3">
              <Link
                href="/financial-data"
                className="block w-full text-left p-3 rounded-lg border border-gray-200 hover:border-primary-300 hover:bg-primary-50 transition-colors"
              >
                <div className="flex items-center space-x-3">
                  <span className="text-xl">💰</span>
                  <div>
                    <p className="font-medium text-gray-900">財務データ更新</p>
                    <p className="text-sm text-gray-600">収入・支出を更新</p>
                  </div>
                </div>
              </Link>

              <Link
                href="/goals"
                className="block w-full text-left p-3 rounded-lg border border-gray-200 hover:border-primary-300 hover:bg-primary-50 transition-colors"
              >
                <div className="flex items-center space-x-3">
                  <span className="text-xl">🎯</span>
                  <div>
                    <p className="font-medium text-gray-900">新しい目標設定</p>
                    <p className="text-sm text-gray-600">財務目標を追加</p>
                  </div>
                </div>
              </Link>

              <Link
                href="/reports"
                className="block w-full text-left p-3 rounded-lg border border-gray-200 hover:border-primary-300 hover:bg-primary-50 transition-colors"
              >
                <div className="flex items-center space-x-3">
                  <span className="text-xl">📋</span>
                  <div>
                    <p className="font-medium text-gray-900">レポート生成</p>
                    <p className="text-sm text-gray-600">PDF形式で出力</p>
                  </div>
                </div>
              </Link>
            </div>
          </div>

          {/* Recommendations */}
          <div className="card">
            <h2 className="text-xl font-semibold text-gray-900 mb-4">推奨事項</h2>
            <div className="space-y-3">
              <div className="p-3 bg-success-50 border border-success-200 rounded-lg">
                <p className="text-sm font-medium text-success-800">✅ 緊急資金は十分確保されています</p>
              </div>
              <div className="p-3 bg-warning-50 border border-warning-200 rounded-lg">
                <p className="text-sm font-medium text-warning-800">⚠️ 老後資金の積立を月額¥50,000増やすことを推奨</p>
              </div>
              <div className="p-3 bg-primary-50 border border-primary-200 rounded-lg">
                <p className="text-sm font-medium text-primary-800">💡 投資利回りを5%→6%に改善すると目標達成が2年早まります</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}