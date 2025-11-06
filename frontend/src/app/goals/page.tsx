export default function GoalsPage() {
  return (
    <div className="container mx-auto px-4 py-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">目標設定・進捗管理</h1>
        <p className="text-gray-600">具体的な財務目標を設定し、進捗を追跡してモチベーションを維持しましょう</p>
      </div>

      {/* Add New Goal Button */}
      <div className="mb-6">
        <button className="btn-primary">
          <span className="mr-2">➕</span>
          新しい目標を追加
        </button>
      </div>

      {/* Goals List */}
      <div className="space-y-6">
        {/* Emergency Fund Goal */}
        <div className="card">
          <div className="flex items-start justify-between mb-4">
            <div className="flex items-center space-x-3">
              <div className="text-2xl">🚨</div>
              <div>
                <h3 className="text-lg font-semibold text-gray-900">緊急資金</h3>
                <p className="text-sm text-gray-600">6ヶ月分の生活費を確保</p>
              </div>
            </div>
            <div className="flex items-center space-x-2">
              <span className="px-2 py-1 bg-success-100 text-success-800 text-xs font-medium rounded-full">
                達成済み
              </span>
              <button className="text-gray-400 hover:text-gray-600">⋯</button>
            </div>
          </div>
          
          <div className="mb-4">
            <div className="flex justify-between items-center mb-2">
              <span className="text-sm font-medium text-gray-700">進捗</span>
              <span className="text-sm font-medium text-gray-900">¥600,000 / ¥600,000</span>
            </div>
            <div className="w-full bg-gray-200 rounded-full h-3">
              <div className="bg-success-500 h-3 rounded-full" style={{ width: '100%' }}></div>
            </div>
            <div className="flex justify-between items-center mt-2">
              <span className="text-xs text-gray-500">目標達成日: 2024年3月</span>
              <span className="text-xs font-medium text-success-600">100%</span>
            </div>
          </div>
        </div>

        {/* Retirement Fund Goal */}
        <div className="card">
          <div className="flex items-start justify-between mb-4">
            <div className="flex items-center space-x-3">
              <div className="text-2xl">🏖️</div>
              <div>
                <h3 className="text-lg font-semibold text-gray-900">老後資金</h3>
                <p className="text-sm text-gray-600">65歳までに3,000万円を準備</p>
              </div>
            </div>
            <div className="flex items-center space-x-2">
              <span className="px-2 py-1 bg-primary-100 text-primary-800 text-xs font-medium rounded-full">
                進行中
              </span>
              <button className="text-gray-400 hover:text-gray-600">⋯</button>
            </div>
          </div>
          
          <div className="mb-4">
            <div className="flex justify-between items-center mb-2">
              <span className="text-sm font-medium text-gray-700">進捗</span>
              <span className="text-sm font-medium text-gray-900">¥19,500,000 / ¥30,000,000</span>
            </div>
            <div className="w-full bg-gray-200 rounded-full h-3">
              <div className="bg-primary-500 h-3 rounded-full" style={{ width: '65%' }}></div>
            </div>
            <div className="flex justify-between items-center mt-2">
              <span className="text-xs text-gray-500">予想達成日: 2040年8月</span>
              <span className="text-xs font-medium text-primary-600">65%</span>
            </div>
          </div>

          <div className="bg-primary-50 border border-primary-200 rounded-lg p-3">
            <p className="text-sm text-primary-800">
              💡 月額貯蓄を¥50,000増やすと、目標達成が2年早まります
            </p>
          </div>
        </div>

        {/* House Purchase Goal */}
        <div className="card">
          <div className="flex items-start justify-between mb-4">
            <div className="flex items-center space-x-3">
              <div className="text-2xl">🏠</div>
              <div>
                <h3 className="text-lg font-semibold text-gray-900">マイホーム資金</h3>
                <p className="text-sm text-gray-600">頭金500万円を5年以内に準備</p>
              </div>
            </div>
            <div className="flex items-center space-x-2">
              <span className="px-2 py-1 bg-warning-100 text-warning-800 text-xs font-medium rounded-full">
                要注意
              </span>
              <button className="text-gray-400 hover:text-gray-600">⋯</button>
            </div>
          </div>
          
          <div className="mb-4">
            <div className="flex justify-between items-center mb-2">
              <span className="text-sm font-medium text-gray-700">進捗</span>
              <span className="text-sm font-medium text-gray-900">¥1,250,000 / ¥5,000,000</span>
            </div>
            <div className="w-full bg-gray-200 rounded-full h-3">
              <div className="bg-warning-500 h-3 rounded-full" style={{ width: '25%' }}></div>
            </div>
            <div className="flex justify-between items-center mt-2">
              <span className="text-xs text-gray-500">目標達成日: 2029年12月</span>
              <span className="text-xs font-medium text-warning-600">25%</span>
            </div>
          </div>

          <div className="bg-warning-50 border border-warning-200 rounded-lg p-3">
            <p className="text-sm text-warning-800">
              ⚠️ 現在のペースでは目標達成が困難です。月額貯蓄の見直しを検討してください
            </p>
          </div>
        </div>
      </div>

      {/* Goal Setting Form Placeholder */}
      <div className="mt-8">
        <div className="card">
          <h2 className="text-xl font-semibold text-gray-900 mb-4">目標設定フォーム</h2>
          <div className="text-center py-8 text-gray-500">
            <div className="text-4xl mb-2">🎯</div>
            <p>目標設定・編集フォーム</p>
            <p className="text-sm">(タスク9.1で実装予定)</p>
          </div>
        </div>
      </div>
    </div>
  );
}