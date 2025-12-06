// API型定義

// 共通型
export interface APIResponse<T> {
  data?: T;
  error?: string;
  message?: string;
}

// 財務データ型
export interface ExpenseItem {
  category: string;
  amount: number;
  description?: string;
}

export interface SavingsItem {
  type: 'deposit' | 'investment' | 'other';
  amount: number;
  description?: string;
}

export interface FinancialProfile {
  monthly_income: number;
  monthly_expenses: ExpenseItem[];
  current_savings: SavingsItem[];
  investment_return: number;
  inflation_rate: number;
}

export interface RetirementData {
  current_age: number;
  retirement_age: number;
  life_expectancy: number;
  monthly_retirement_expenses: number;
  pension_amount: number;
}

export interface EmergencyFund {
  target_months: number;
  monthly_expenses: number;
  current_amount: number;
}

export interface FinancialData {
  id?: string;
  user_id: string;
  profile?: FinancialProfile;
  retirement?: RetirementData;
  emergency_fund?: EmergencyFund;
  created_at?: string;
  updated_at?: string;
}

// 目標型
export type GoalType = 'savings' | 'retirement' | 'emergency' | 'custom';

export interface Goal {
  id?: string;
  user_id: string;
  goal_type: GoalType; // プロパティ名を 'goal_type' に変更
  title: string;
  target_amount: number;
  target_date: string;
  current_amount: number;
  monthly_contribution: number;
  is_active: boolean;
  created_at?: string;
  updated_at?: string;
}

// 計算リクエスト型
export interface AssetProjectionRequest {
  user_id: string;
  years: number;
  monthly_income: number;
  monthly_expenses: number;
  current_savings: number;
  investment_return: number;
  inflation_rate: number;
}

export interface RetirementCalculationRequest {
  user_id: string;
  current_age: number;
  retirement_age: number;
  life_expectancy: number;
  monthly_retirement_expenses: number;
  pension_amount: number;
  current_savings: number;
  monthly_savings: number;
  investment_return: number;
  inflation_rate: number;
}

export interface EmergencyFundRequest {
  user_id: string;
  monthly_expenses: number;
  target_months: number;
  current_savings: number;
}

export interface GoalProjectionRequest {
  user_id: string;
  goal_id: string;
  target_amount: number;
  target_date: string;
  current_amount: number;
  monthly_contribution: number;
  investment_return: number;
}

// 計算レスポンス型
export interface AssetProjectionPoint {
  year: number;
  total_assets: number;
  real_value: number;
  contributed_amount: number;
  investment_gains: number;
}

export interface AssetProjectionResponse {
  projections: AssetProjectionPoint[];
  final_amount: number;
  total_contributions: number;
  total_gains: number;
}

export interface RetirementCalculationResponse {
  required_amount: number;
  projected_amount: number;
  shortfall: number;
  sufficiency_rate: number;
  recommended_monthly_savings: number;
  years_until_retirement: number;
}

export interface EmergencyFundResponse {
  required_amount: number;
  current_amount: number;
  shortfall: number;
  sufficiency_rate: number;
  months_to_target: number;
}

export interface GoalProjectionResponse {
  goal_id: string;
  is_achievable: boolean;
  projected_completion_date: string;
  shortfall: number;
  recommended_monthly_contribution: number;
  progress_rate: number;
}

// レポート型
export interface ReportRequest {
  user_id: string;
  include_charts?: boolean;
  format?: 'json' | 'pdf';
}

export interface FinancialSummaryReport {
  user_id: string;
  summary: {
    total_assets: number;
    monthly_net_savings: number;
    emergency_fund_status: string;
    retirement_readiness: string;
  };
  generated_at: string;
}
