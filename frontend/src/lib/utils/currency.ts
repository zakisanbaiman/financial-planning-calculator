/**
 * Format a number as Japanese Yen currency
 * @param amount - The amount to format
 * @returns Formatted currency string (e.g., "¥1,234,567")
 */
export function formatCurrency(amount: number): string {
  return new Intl.NumberFormat('ja-JP', {
    style: 'currency',
    currency: 'JPY',
    maximumFractionDigits: 0,
  }).format(amount);
}
