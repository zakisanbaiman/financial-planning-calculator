export default function ReportsPage() {
  return (
    <div className="container mx-auto px-4 py-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">レポート生成</h1>
        <p className="text-gray-600">財務状況をまとめたレポートをPDF形式で生成・印刷できます</p>
      </div>

      {/* Report Generation Options */}
      <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6 mb-8">
        <div className="card">
          <div className="text-center">
            <div className="text-4xl mb-4">📊</div>
            <h3 className="text-lg font-semibold text-gray-900 mb-2">総合財務レポート</h3>
            <p className="text-gray-600 text-sm mb-4">
              現在の財務状況と将来予測を含む包括的なレポート
            </p>
            <button className="btn-primary w-full">PDF生成</button>
          </div>
        </div>

        <div className="card">
          <div className="text-center">
            <div className="text-4xl mb-4">🎯</div>
            <h3 className="text-lg font-semibold text-gray-900 mb-2">目標進捗レポート</h3>
            <p className="text-gray-600 text-sm mb-4">
              設定した目標の進捗状況と達成予測
            </p>
            <button className="btn-primary w-full">PDF生成</button>
          </div>
        </div>

        <div className="card">
          <div className="text-center">
            <div className="text-4xl mb-4">📈</div>
            <h3 className="text-lg font-semibold text-gray-900 mb-2">資産推移レポート</h3>
            <p className="text-gray-600 text-sm mb-4">
              資産の推移予測とシナリオ分析
            </p>
            <button className="btn-primary w-full">PDF生成</button>
          </div>
        </div>
      </div>

      {/* Report Preview */}
      <div className="grid lg:grid-cols-3 gap-8">
        {/* Report Content Preview */}
        <div className="lg:col-span-2">
          <div className="card">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-xl font-semibold text-gray-900">レポートプレビュー</h2>
              <div className="flex space-x-2">
                <button className="btn-secondary text-sm">編集</button>
                <button className="btn-primary text-sm">PDF出力</button>
              </div>
            </div>
            
            {/* Mock Report Content */}
            <div className="bg-white border border-gray-200 rounded-lg p-6 min-h-96">
              <div className="text-center mb-6">
                <h1 className="text-2xl font-bold text-gray-900 mb-2">財務計画レポート</h1>
                <p className="text-gray-600">作成日: 2024年11月7日</p>
              </div>

              <div className="space-y-6">
                <section>
                  <h2 className="text-lg font-semibold text-gray-900 mb-3">現在の財務状況</h2>
                  <div className="grid grid-cols-2 gap-4">
                    <div className="bg-gray-50 p-3 rounded">
                      <p className="text-sm text-gray-600">月収</p>
                      <p className="text-lg font-semibold">¥400,000</p>
                    </div>
                    <div className="bg-gray-50 p-3 rounded">
                      <p className="text-sm text-gray-600">月間支出</p>
                      <p className="text-lg font-semibold">¥280,000</p>
                    </div>
                    <div className="bg-gray-50 p-3 rounded">
                      <p className="text-sm text-gray-600">月間貯蓄</p>
                      <p className="text-lg font-semibold text-success-600">¥120,000</p>
                    </div>
                    <div className="bg-gray-50 p-3 rounded">
                      <p className="text-sm text-gray-600">総資産</p>
                      <p className="text-lg font-semibold">¥1,500,000</p>
                    </div>
                  </div>
                </section>

                <section>
                  <h2 className="text-lg font-semibold text-gray-900 mb-3">将来予測</h2>
                  <div className="h-32 bg-gray-100 rounded flex items-center justify-center">
                    <span className="text-gray-500">資産推移グラフ</span>
                  </div>
                </section>

                <section>
                  <h2 className="text-lg font-semibold text-gray-900 mb-3">目標達成状況</h2>
                  <div className="space-y-2">
                    <div className="flex justify-between items-center">
                      <span>緊急資金</span>
                      <span className="text-success-600 font-medium">100%</span>
                    </div>
                    <div className="flex justify-between items-center">
                      <span>老後資金</span>
                      <span className="text-primary-600 font-medium">65%</span>
                    </div>
                    <div className="flex justify-between items-center">
                      <span>マイホーム資金</span>
                      <span className="text-warning-600 font-medium">25%</span>
                    </div>
                  </div>
                </section>
              </div>
            </div>
          </div>
        </div>

        {/* Report Settings */}
        <div className="space-y-6">
          <div className="card">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">レポート設定</h3>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  レポート期間
                </label>
                <select className="input-field">
                  <option>過去1年間</option>
                  <option>過去6ヶ月</option>
                  <option>過去3ヶ月</option>
                  <option>カスタム期間</option>
                </select>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  含める項目
                </label>
                <div className="space-y-2">
                  <label className="flex items-center">
                    <input type="checkbox" className="mr-2" defaultChecked />
                    <span className="text-sm">現在の財務状況</span>
                  </label>
                  <label className="flex items-center">
                    <input type="checkbox" className="mr-2" defaultChecked />
                    <span className="text-sm">資産推移予測</span>
                  </label>
                  <label className="flex items-center">
                    <input type="checkbox" className="mr-2" defaultChecked />
                    <span className="text-sm">目標進捗状況</span>
                  </label>
                  <label className="flex items-center">
                    <input type="checkbox" className="mr-2" />
                    <span className="text-sm">詳細な計算過程</span>
                  </label>
                  <label className="flex items-center">
                    <input type="checkbox" className="mr-2" />
                    <span className="text-sm">推奨事項</span>
                  </label>
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  出力形式
                </label>
                <select className="input-field">
                  <option>PDF (推奨)</option>
                  <option>Excel</option>
                  <option>CSV</option>
                </select>
              </div>
            </div>
          </div>

          <div className="card">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">レポート履歴</h3>
            <div className="space-y-3">
              <div className="flex items-center justify-between p-3 bg-gray-50 rounded">
                <div>
                  <p className="text-sm font-medium text-gray-900">総合財務レポート</p>
                  <p className="text-xs text-gray-600">2024/11/01</p>
                </div>
                <button className="text-primary-600 hover:text-primary-700 text-sm">
                  ダウンロード
                </button>
              </div>
              <div className="flex items-center justify-between p-3 bg-gray-50 rounded">
                <div>
                  <p className="text-sm font-medium text-gray-900">目標進捗レポート</p>
                  <p className="text-xs text-gray-600">2024/10/15</p>
                </div>
                <button className="text-primary-600 hover:text-primary-700 text-sm">
                  ダウンロード
                </button>
              </div>
              <div className="flex items-center justify-between p-3 bg-gray-50 rounded">
                <div>
                  <p className="text-sm font-medium text-gray-900">資産推移レポート</p>
                  <p className="text-xs text-gray-600">2024/10/01</p>
                </div>
                <button className="text-primary-600 hover:text-primary-700 text-sm">
                  ダウンロード
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Report Generation Placeholder */}
      <div className="mt-8">
        <div className="card">
          <h2 className="text-xl font-semibold text-gray-900 mb-4">PDF生成機能</h2>
          <div className="text-center py-8 text-gray-500">
            <div className="text-4xl mb-2">📋</div>
            <p>PDF生成とダウンロード機能</p>
            <p className="text-sm">(タスク10.1で実装予定)</p>
          </div>
        </div>
      </div>
    </div>
  );
}