import type { FinancialProfile } from '@/types/api';

/**
 * 財務プロファイルをCSV文字列（マルチセクション形式）に変換する
 *
 * フォーマット:
 *   # SECTION: PROFILE
 *   field,value
 *   monthly_income,300000
 *   ...
 *
 *   # SECTION: EXPENSES
 *   category,amount,description
 *   生活費,100000,
 *   ...
 *
 * バックエンドの encoding/csv と同じフォーマットを生成することで
 * アップロード時にそのまま利用できる。
 */
export function generateCSVFromProfile(profile: FinancialProfile): string {
  const lines: string[] = [];

  // PROFILE セクション
  lines.push('# SECTION: PROFILE');
  lines.push('field,value');
  lines.push(`monthly_income,${profile.monthly_income ?? 0}`);
  lines.push(`investment_return,${profile.investment_return}`);
  lines.push(`inflation_rate,${profile.inflation_rate}`);

  // EXPENSES セクション
  lines.push('');
  lines.push('# SECTION: EXPENSES');
  lines.push('category,amount,description');
  for (const expense of profile.monthly_expenses ?? []) {
    const desc = expense.description ?? '';
    // カンマを含む文字列をクォートする
    const category = escapeCSVField(expense.category);
    lines.push(`${category},${expense.amount},${escapeCSVField(desc)}`);
  }

  // SAVINGS セクション
  lines.push('');
  lines.push('# SECTION: SAVINGS');
  lines.push('type,amount,description');
  for (const saving of profile.current_savings ?? []) {
    const desc = saving.description ?? '';
    lines.push(`${escapeCSVField(saving.type)},${saving.amount},${escapeCSVField(desc)}`);
  }

  return lines.join('\n');
}

/**
 * CSV文字列をダウンロードファイルとして保存する
 * Blob API を使ったクライアント側のみの実装（バックエンド不使用）
 */
export function downloadCSVLocally(csvContent: string, fileName: string): void {
  const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
  const url = URL.createObjectURL(blob);

  const a = document.createElement('a');
  a.href = url;
  a.download = fileName;
  a.style.display = 'none';
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);

  // blob: URL を解放してメモリリークを防ぐ
  URL.revokeObjectURL(url);
}

function escapeCSVField(field: string): string {
  if (field.includes(',') || field.includes('"') || field.includes('\n')) {
    return `"${field.replace(/"/g, '""')}"`;
  }
  return field;
}
