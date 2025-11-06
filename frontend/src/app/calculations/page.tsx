export default function CalculationsPage() {
  return (
    <div className="container mx-auto px-4 py-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">財務計算機</h1>
        <p className="text-gray-600">資産推移、老後資金、緊急資金の計算と可視化を行います</p>
      </div>

      {/* Calculation Categories */}
      <div className="grid md:grid-cols-3 gap-6 mb-8">
        <div className="card">
          <div className="text-center">
            <div className="text-4xl mb-4">📈</div>
            <h3 className="text-lg font-semibold text-gray-900 mb-2">資産推移シミュレーション</h3>
            <p className="text-gray-600 text-sm mb-4">
              現在の貯蓄ペースで将来どれだけ資産が増えるかを計算
            </p>
            <button className="btn-primary w-full">計算開始</button>
          </div>
        </div>

        <div className="card">
          <div className="text-center">
            <div className="text-4xl mb-4">🏖️</div>
            <h3 className="text-lg font-semibold text-gray-900 mb-2">老後資金計算</h3>
            <p className="text-gray-600 text-sm mb-4">
              退職後に必要な資金と年金額を考慮した計算
            </p>
            <button className="btn-primary w-full">計算開始</button>
          </div>
        </div>

        <div className="card">
          <div className="text-center">
            <div className="text-4xl mb-4">🚨</div>
            <h3 className="text-lg font-semibold text-gray-900 mb-2">緊急資金計算</h3>
            <p className="text-gray-600 text-sm mb-4">
              万が一の時に必要な緊急資金を計算
            </p>
            <button className="btn-primary w-full">計算開始</button>
          </div>
        </div>
      </div>

      {/* Sample Calculation Results */}
      <div className="grid lg:grid-cols-2 gap-8">
        {/* Asset Projection Chart */}
        <div className="card">
          <h2 className="text-xl font-semibold text-gray-900 mb-4">資産推移予測（30年間）</h2>
          <div className="h-64 bg-gray-100 rounded-lg flex items-center justify-center mb-4">
            <div className="text-center text-gray-500">
              <div className="text-4xl mb-2">📊</div>
              <p>資産推移グラフ</p>
              <p className="text-sm">(タスク8.1で実装予定)</p>
            </div>
          </div>
          <div className="grid grid-cols-3 gap-4 text-center">
            <div>
              <p className="text-sm text-gray-600">10年後</p>
              <p className="text-lg font-semibold text-gray-900">¥16,200,000</p>
            </div>
            <div>
              <p className="text-sm text-gray-600">20年後</p>
              <p className="text-lg font-semibold text-gray-900">¥38,400,000</p>
            </div>
            <div>
              <p className="text-sm text-gray-600">30年後</p>
              <p className="text-lg font-semibold text-gray-900">¥69,800,000</p>
            </div>
          </div>
        </div>

        {/* Calculation Summary */}
        <div className="space-y-6">
          {/* Retirement Calculation */}
          <div className="card">
            <h3 className="text-lg font-semibold text-gray-900 mb-3">老後資金計算結果</h3>
            <div className="space-y-3">
              <div className="flex justify-between">
                <span className="text-gray-600">必要老後資金</span>
                <span className="font-medium text-gray-900">¥30,000,000</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600">予想達成額</span>
                <span className="font-medium text-gray-900">¥45,600,000</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600">余裕額</span>
                <span className="font-medium text-success-600">+¥15,600,000</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600">充足率</span>
                <span className="font-medium text-success-600">152%</span>
              </div>
            </div>
          </div>

          {/* Emergency Fund Calculation */}
          <div className="card">
            <h3 className="text-lg font-semibold text-gray-900 mb-3">緊急資金計算結果</h3>
            <div className="space-y-3">
              <div className="flex justify-between">
                <span className="text-gray-600">推奨緊急資金</span>
                <span className="font-medium text-gray-900">¥1,680,000</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600">現在の緊急資金</span>
                <span className="font-medium text-gray-900">¥600,000</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600">不足額</span>
                <span className="font-medium text-warning-600">¥1,080,000</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600">達成までの期間</span>
                <span className="font-medium text-gray-900">9ヶ月</span>
              </div>
            </div>
          </div>

          {/* Calculation Parameters */}
          <div className="card">
            <h3 className="text-lg font-semibold text-gray-900 mb-3">計算パラメータ</h3>
            <div className="space-y-2 text-sm">
              <div className="flex justify-between">
                <span className="text-gray-600">投資利回り</span>
                <span className="text-gray-900">5.0%</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600">インフレ率</span>
                <span className="text-gray-900">2.0%</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600">月間貯蓄額</span>
                <span className="text-gray-900">¥120,000</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600">退職予定年齢</span>
                <span className="text-gray-900">65歳</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Calculation Forms Placeholder */}
      <div className="mt-8">
        <div className="card">
          <h2 className="text-xl font-semibold text-gray-900 mb-4">詳細計算フォーム</h2>
          <div className="text-center py-8 text-gray-500">
            <div className="text-4xl mb-2">🧮</div>
            <p>計算パラメータ入力フォーム</p>
            <p className="text-sm">(タスク8.1-8.3で実装予定)</p>
          </div>
        </div>
      </div>
    </div>
  );
}