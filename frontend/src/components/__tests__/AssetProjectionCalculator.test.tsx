import { z } from 'zod';

// Test the validation schema for asset projection
describe('AssetProjectionCalculator Validation', () => {
  // Extract the schema definition from the component
  const assetProjectionSchema = z.object({
    years: z.number().min(1, '1年以上を指定してください').max(100, '100年以内で指定してください'),
    monthly_income: z.number().min(0, '0以上の値を入力してください'),
    monthly_expenses: z.number().min(0, '0以上の値を入力してください'),
    current_savings: z.number().min(0, '0以上の値を入力してください'),
    investment_return: z.number().min(0, '0以上の値を入力してください').max(100, '100以下の値を入力してください'),
    inflation_rate: z.number().min(0, '0以上の値を入力してください').max(50, '50以下の値を入力してください'),
  });

  describe('years validation', () => {
    it('should accept 1 year', () => {
      const result = assetProjectionSchema.safeParse({
        years: 1,
        monthly_income: 400000,
        monthly_expenses: 280000,
        current_savings: 1500000,
        investment_return: 5.0,
        inflation_rate: 2.0,
      });
      expect(result.success).toBe(true);
    });

    it('should accept 100 years', () => {
      const result = assetProjectionSchema.safeParse({
        years: 100,
        monthly_income: 400000,
        monthly_expenses: 280000,
        current_savings: 1500000,
        investment_return: 5.0,
        inflation_rate: 2.0,
      });
      expect(result.success).toBe(true);
    });

    it('should accept 50 years (previously the max)', () => {
      const result = assetProjectionSchema.safeParse({
        years: 50,
        monthly_income: 400000,
        monthly_expenses: 280000,
        current_savings: 1500000,
        investment_return: 5.0,
        inflation_rate: 2.0,
      });
      expect(result.success).toBe(true);
    });

    it('should reject 0 years', () => {
      const result = assetProjectionSchema.safeParse({
        years: 0,
        monthly_income: 400000,
        monthly_expenses: 280000,
        current_savings: 1500000,
        investment_return: 5.0,
        inflation_rate: 2.0,
      });
      expect(result.success).toBe(false);
      if (!result.success) {
        expect(result.error.issues[0].message).toContain('1年以上');
      }
    });

    it('should reject 101 years', () => {
      const result = assetProjectionSchema.safeParse({
        years: 101,
        monthly_income: 400000,
        monthly_expenses: 280000,
        current_savings: 1500000,
        investment_return: 5.0,
        inflation_rate: 2.0,
      });
      expect(result.success).toBe(false);
      if (!result.success) {
        expect(result.error.issues[0].message).toContain('100年以内');
      }
    });

    it('should reject negative years', () => {
      const result = assetProjectionSchema.safeParse({
        years: -1,
        monthly_income: 400000,
        monthly_expenses: 280000,
        current_savings: 1500000,
        investment_return: 5.0,
        inflation_rate: 2.0,
      });
      expect(result.success).toBe(false);
    });
  });
});
