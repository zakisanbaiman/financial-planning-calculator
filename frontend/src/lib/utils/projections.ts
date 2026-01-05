import type { AssetProjectionPoint } from '@/types/api';

/**
 * Generate asset projection data with monthly compounding
 * @param years - Number of years to project
 * @param initialAssets - Initial asset amount
 * @param monthlyContribution - Monthly contribution amount
 * @param investmentReturn - Annual investment return rate (as decimal, e.g., 0.05 for 5%)
 * @param inflationRate - Annual inflation rate (as decimal, e.g., 0.02 for 2%)
 * @returns Array of asset projection points
 */
export function generateAssetProjections(
  years: number,
  initialAssets: number,
  monthlyContribution: number,
  investmentReturn: number = 0.05,
  inflationRate: number = 0.02
): AssetProjectionPoint[] {
  const projections: AssetProjectionPoint[] = [];
  const monthlyRate = investmentReturn / 12;
  
  for (let year = 0; year <= years; year++) {
    const months = year * 12;
    
    // Calculate total assets using monthly compounding
    // Initial assets compound, and each monthly contribution compounds from its deposit date
    let totalAssets = initialAssets * Math.pow(1 + monthlyRate, months);
    
    // Add compounded monthly contributions using future value of annuity formula
    if (months > 0 && monthlyRate > 0) {
      totalAssets += monthlyContribution * ((Math.pow(1 + monthlyRate, months) - 1) / monthlyRate);
    }
    
    // Calculate total contributed amount (principal only)
    const contributedAmount = initialAssets + (monthlyContribution * months);
    
    // Calculate real value (adjusted for inflation)
    const realValue = totalAssets / Math.pow(1 + inflationRate, year);
    
    // Calculate investment gains
    const investmentGains = totalAssets - contributedAmount;
    
    projections.push({
      year,
      total_assets: Math.round(totalAssets),
      real_value: Math.round(realValue),
      contributed_amount: Math.round(contributedAmount),
      investment_gains: Math.round(investmentGains),
    });
  }
  
  return projections;
}
