import type { Metadata } from 'next';

export const metadata: Metadata = {
  title: '利用規約 — FinPlan',
};

export default function TermsPage() {
  return (
    <div className="container mx-auto px-4 py-12 max-w-3xl">
      <h1 className="font-display text-3xl font-semibold text-ink-900 dark:text-ink-100 mb-8">
        利用規約
      </h1>

      <div className="space-y-8 font-body text-ink-700 dark:text-ink-300">
        <section>
          <h2 className="text-xl font-semibold text-ink-900 dark:text-ink-100 mb-3">
            第1条（適用）
          </h2>
          <p className="leading-relaxed">
            本規約は、FinPlan（以下「本サービス」）の利用条件を定めるものです。
            本サービスをご利用になる場合は、本規約に同意したものとみなします。
          </p>
        </section>

        <section>
          <h2 className="text-xl font-semibold text-ink-900 dark:text-ink-100 mb-3">
            第2条（利用目的）
          </h2>
          <p className="leading-relaxed">
            本サービスは個人の財務計画・資産形成のシミュレーションを目的とした学習・参考ツールです。
            投資判断や金融商品の購入における意思決定の根拠として利用することはご遠慮ください。
          </p>
        </section>

        <section>
          <h2 className="text-xl font-semibold text-ink-900 dark:text-ink-100 mb-3">
            第3条（禁止事項）
          </h2>
          <p className="leading-relaxed mb-3">
            以下の行為を禁止します。
          </p>
          <ul className="list-disc list-inside space-y-2">
            <li>本サービスの運営を妨害する行為</li>
            <li>他のユーザーに不利益を与える行為</li>
            <li>法令または公序良俗に違反する行為</li>
            <li>本サービスを商業目的に無断で利用する行為</li>
          </ul>
        </section>

        <section>
          <h2 className="text-xl font-semibold text-ink-900 dark:text-ink-100 mb-3">
            第4条（免責事項）
          </h2>
          <p className="leading-relaxed">
            本サービスが提供する計算結果・シミュレーション結果は参考情報であり、
            投資・金融に関する専門的なアドバイスを意味するものではありません。
            本サービスの利用によって生じた損害については、一切の責任を負いません。
          </p>
        </section>

        <section>
          <h2 className="text-xl font-semibold text-ink-900 dark:text-ink-100 mb-3">
            第5条（規約の変更）
          </h2>
          <p className="leading-relaxed">
            本規約は予告なく変更される場合があります。変更後の規約は本ページに掲載した時点で効力を生じます。
          </p>
        </section>

        <p className="text-sm text-ink-400 dark:text-ink-500 pt-4 border-t border-ink-200 dark:border-ink-800">
          最終更新日：2026年4月7日
        </p>
      </div>
    </div>
  );
}
