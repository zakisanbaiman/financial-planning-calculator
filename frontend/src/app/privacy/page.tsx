import type { Metadata } from 'next';

export const metadata: Metadata = {
  title: 'プライバシーポリシー — FinPlan',
};

export default function PrivacyPage() {
  return (
    <div className="container mx-auto px-4 py-12 max-w-3xl">
      <h1 className="font-display text-3xl font-semibold text-ink-900 dark:text-ink-100 mb-8">
        プライバシーポリシー
      </h1>

      <div className="space-y-8 font-body text-ink-700 dark:text-ink-300">
        <section>
          <h2 className="text-xl font-semibold text-ink-900 dark:text-ink-100 mb-3">
            収集する情報
          </h2>
          <p className="leading-relaxed mb-3">
            本サービスでは以下の情報を収集します。
          </p>
          <ul className="list-disc list-inside space-y-2">
            <li>メールアドレス（アカウント登録時）</li>
            <li>財務プロフィール情報（収入・支出・資産など）</li>
            <li>目標設定情報</li>
            <li>アクセスログ（IPアドレス、ブラウザ情報など）</li>
          </ul>
        </section>

        <section>
          <h2 className="text-xl font-semibold text-ink-900 dark:text-ink-100 mb-3">
            情報の利用目的
          </h2>
          <p className="leading-relaxed mb-3">
            収集した情報は以下の目的で利用します。
          </p>
          <ul className="list-disc list-inside space-y-2">
            <li>サービスの提供・改善</li>
            <li>財務計算・シミュレーション機能の提供</li>
            <li>不正利用の防止</li>
          </ul>
        </section>

        <section>
          <h2 className="text-xl font-semibold text-ink-900 dark:text-ink-100 mb-3">
            第三者への提供
          </h2>
          <p className="leading-relaxed">
            法令に基づく場合を除き、ユーザーの個人情報を第三者に提供することはありません。
          </p>
        </section>

        <section>
          <h2 className="text-xl font-semibold text-ink-900 dark:text-ink-100 mb-3">
            Cookieの利用
          </h2>
          <p className="leading-relaxed">
            本サービスでは認証状態の維持にCookieを利用しています。
            また、Google Analytics（アクセス解析）を使用しており、匿名化されたアクセスデータを収集しています。
          </p>
        </section>

        <section>
          <h2 className="text-xl font-semibold text-ink-900 dark:text-ink-100 mb-3">
            データの削除
          </h2>
          <p className="leading-relaxed">
            アカウントを削除することで、保存された財務データはすべて削除されます。
            削除に関するお問い合わせは GitHub リポジトリの Issues よりご連絡ください。
          </p>
        </section>

        <p className="text-sm text-ink-400 dark:text-ink-500 pt-4 border-t border-ink-200 dark:border-ink-800">
          最終更新日：2026年4月7日
        </p>
      </div>
    </div>
  );
}
