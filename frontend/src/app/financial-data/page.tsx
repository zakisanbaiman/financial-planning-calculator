export default function FinancialDataPage() {
  return (
    <div className="container mx-auto px-4 py-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">財務データ管理</h1>
        <p className="text-gray-600">収入、支出、貯蓄状況を入力・管理して、正確な将来予測の基盤を作成します</p>
      </div>

      <div className="grid lg:grid-cols-2 gap-8">
        {/* Current Data Display */}
        <div className="card">
          <h2 className="text-xl font-semibold text-gray-900 mb-4">現在の財務状況</h2>
          <div className="space-y-4">
            <div className="flex justify-between items-center py-2 border-b border-gray-100">
              <span className="text-gray-600">月収</span>
              <span className="font-medium text-gray-900">¥400,000</span>
            </div>
            <div className="flex justify-between items-center py-2 border-b border-gray-100">
              <span className="text-gray-600">月間支出</span>
              <span className="font-medium text-gray-900">¥280,000</span>
            </div>
            <div className="flex justify-between items-center py-2 border-b border-gray-100">
              <span className="text-gray-600">月間純貯蓄</span>
              <span className="font-medium text-success-600">¥120,000</span>
            </div>
            <div className="flex justify-between items-center py-2">
              <span className="text-gray-600">総資産</span>
              <span className="font-medium text-gray-900">¥1,500,000</span>
            </div>
          </div>
        </div>

        {/* Input Form Placeholder */}
        <div className="card">
          <h2 className="text-xl font-semibold text-gray-900 mb-4">データ入力・更新</h2>
          <div className="space-y-4">
            <div className="text-center py-8 text-gray-500">
              <div className="text-4xl mb-2">📝</div>
              <p>財務データ入力フォーム</p>
              <p className="text-sm">(タスク7.1で実装予定)</p>
            </div>
          </div>
        </div>
      </div>

      {/* Expense Breakdown */}
      <div className="mt-8">
        <div className="card">
          <h2 className="text-xl font-semibold text-gray-900 mb-4">支出内訳</h2>
          <div className="grid md:grid-cols-2 gap-6">
            <div className="space-y-3">
              <div className="flex justify-between items-center">
                <span className="text-gray-600">住居費</span>
                <span className="font-medium text-gray-900">¥120,000 (43%)</span>
              </div>
              <div className="w-full bg-gray-200 rounded-full h-2">
                <div className="bg-primary-500 h-2 rounded-full" style={{ width: '43%' }}></div>
              </div>
            </div>
            <div className="space-y-3">
              <div className="flex justify-between items-center">
                <span className="text-gray-600">食費</span>
                <span className="font-medium text-gray-900">¥60,000 (21%)</span>
              </div>
              <div className="w-full bg-gray-200 rounded-full h-2">
                <div className="bg-success-500 h-2 rounded-full" style={{ width: '21%' }}></div>
              </div>
            </div>
            <div className="space-y-3">
              <div className="flex justify-between items-center">
                <span className="text-gray-600">交通費</span>
                <span className="font-medium text-gray-900">¥20,000 (7%)</span>
              </div>
              <div className="w-full bg-gray-200 rounded-full h-2">
                <div className="bg-warning-500 h-2 rounded-full" style={{ width: '7%' }}></div>
              </div>
            </div>
            <div className="space-y-3">
              <div className="flex justify-between items-center">
                <span className="text-gray-600">その他</span>
                <span className="font-medium text-gray-900">¥80,000 (29%)</span>
              </div>
              <div className="w-full bg-gray-200 rounded-full h-2">
                <div className="bg-error-500 h-2 rounded-full" style={{ width: '29%' }}></div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}