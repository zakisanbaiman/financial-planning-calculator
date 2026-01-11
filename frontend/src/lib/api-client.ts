// API クライアント

import type {
  FinancialData,
  FinancialProfile,
  RetirementData,
  EmergencyFund,
  Goal,
  AssetProjectionRequest,
  AssetProjectionResponse,
  RetirementCalculationRequest,
  RetirementCalculationResponse,
  EmergencyFundRequest,
  EmergencyFundResponse,
  GoalProjectionRequest,
  GoalProjectionResponse,
  ReportRequest,
  FinancialSummaryReport,
  APIResponse,
} from '@/types/api';
import { getAuthToken, isAuthTokenValid } from './contexts/AuthContext';

// API ベースURL
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api';

// エラークラス
export class APIError extends Error {
  constructor(
    message: string,
    public status?: number,
    public data?: any
  ) {
    super(message);
    this.name = 'APIError';
  }
}

// 共通リクエストヘルパー
async function request<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const url = `${API_BASE_URL}${endpoint}`;
  
  // 認証トークンを取得
  const token = getAuthToken();
  
  const config: RequestInit = {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...(token && { 'Authorization': `Bearer ${token}` }),
      ...options.headers,
    },
  };

  try {
    const response = await fetch(url, config);
    
    // 401エラー（未認証）の場合、ログインページにリダイレクト
    if (response.status === 401) {
      // トークンが無効なのでクリア
      if (typeof window !== 'undefined') {
        localStorage.removeItem('auth_token');
        localStorage.removeItem('auth_user');
        localStorage.removeItem('auth_expires');
        window.location.href = '/login';
      }
      throw new APIError('認証が必要です。ログインしてください。', 401);
    }
    
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      throw new APIError(
        errorData.error || errorData.message || `HTTP ${response.status}`,
        response.status,
        errorData
      );
    }

    return await response.json();
  } catch (error) {
    if (error instanceof APIError) {
      throw error;
    }
    throw new APIError(
      error instanceof Error ? error.message : 'ネットワークエラーが発生しました'
    );
  }
}

// 財務データAPI
export const financialDataAPI = {
  // 財務データ作成
  create: async (data: FinancialData): Promise<FinancialData> => {
    return request<FinancialData>('/financial-data', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  },

  // 財務データ取得
  get: async (userId: string): Promise<FinancialData> => {
    return request<FinancialData>(`/financial-data?user_id=${userId}`);
  },

  // 財務プロファイル更新
  updateProfile: async (
    userId: string,
    profile: FinancialProfile
  ): Promise<FinancialData> => {
    return request<FinancialData>(`/financial-data/${userId}/profile`, {
      method: 'PUT',
      body: JSON.stringify(profile),
    });
  },

  // 退職データ更新
  updateRetirement: async (
    userId: string,
    retirement: RetirementData
  ): Promise<FinancialData> => {
    return request<FinancialData>(`/financial-data/${userId}/retirement`, {
      method: 'PUT',
      body: JSON.stringify(retirement),
    });
  },

  // 緊急資金更新
  updateEmergencyFund: async (
    userId: string,
    emergencyFund: EmergencyFund
  ): Promise<FinancialData> => {
    return request<FinancialData>(`/financial-data/${userId}/emergency-fund`, {
      method: 'PUT',
      body: JSON.stringify(emergencyFund),
    });
  },

  // 財務データ削除
  delete: async (userId: string): Promise<void> => {
    return request<void>(`/financial-data/${userId}`, {
      method: 'DELETE',
    });
  },
};

// 計算API
export const calculationsAPI = {
  // 資産推移計算
  assetProjection: async (
    data: AssetProjectionRequest
  ): Promise<AssetProjectionResponse> => {
    return request<AssetProjectionResponse>('/calculations/asset-projection', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  },

  // 老後資金計算
  retirement: async (
    data: RetirementCalculationRequest
  ): Promise<RetirementCalculationResponse> => {
    return request<RetirementCalculationResponse>('/calculations/retirement', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  },

  // 緊急資金計算
  emergencyFund: async (
    data: EmergencyFundRequest
  ): Promise<EmergencyFundResponse> => {
    return request<EmergencyFundResponse>('/calculations/emergency-fund', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  },

  // 目標達成計算
  goalProjection: async (
    data: GoalProjectionRequest
  ): Promise<GoalProjectionResponse> => {
    return request<GoalProjectionResponse>('/calculations/goal-projection', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  },
};

// 目標API
export const goalsAPI = {
  // 目標作成
  create: async (goal: Goal): Promise<Goal> => {
    return request<Goal>('/goals', {
      method: 'POST',
      body: JSON.stringify(goal),
    });
  },

  // 目標一覧取得
  list: async (userId: string): Promise<Goal[]> => {
    return request<Goal[]>(`/goals?user_id=${userId}`);
  },

  // 目標取得
  get: async (id: string, userId: string): Promise<Goal> => {
    return request<Goal>(`/goals/${id}?user_id=${userId}`);
  },

  // 目標更新
  update: async (id: string, userId: string, goal: Partial<Goal>): Promise<Goal> => {
    return request<Goal>(`/goals/${id}?user_id=${userId}`, {
      method: 'PUT',
      body: JSON.stringify(goal),
    });
  },

  // 目標進捗更新
  updateProgress: async (
    id: string,
    userId: string,
    currentAmount: number
  ): Promise<Goal> => {
    return request<Goal>(`/goals/${id}/progress?user_id=${userId}`, {
      method: 'PUT',
      body: JSON.stringify({ current_amount: currentAmount }),
    });
  },

  // 目標削除
  delete: async (id: string, userId: string): Promise<void> => {
    return request<void>(`/goals/${id}?user_id=${userId}`, {
      method: 'DELETE',
    });
  },

  // 目標推奨事項取得
  getRecommendations: async (id: string, userId: string): Promise<any> => {
    return request<any>(`/goals/${id}/recommendations?user_id=${userId}`);
  },

  // 目標実現可能性分析
  analyzeFeasibility: async (id: string, userId: string): Promise<any> => {
    return request<any>(`/goals/${id}/feasibility?user_id=${userId}`);
  },
};

// レポートAPI
export const reportsAPI = {
  // 財務サマリーレポート生成
  financialSummary: async (
    reportRequest: ReportRequest
  ): Promise<FinancialSummaryReport> => {
    return request<FinancialSummaryReport>('/reports/financial-summary', {
      method: 'POST',
      body: JSON.stringify(reportRequest),
    });
  },

  // PDFレポート取得
  getPDF: async (userId: string): Promise<Blob> => {
    const url = `${API_BASE_URL}/reports/pdf?user_id=${userId}`;
    const response = await fetch(url);
    
    if (!response.ok) {
      throw new APIError(`HTTP ${response.status}`, response.status);
    }
    
    return await response.blob();
  },
};

// ヘルスチェック
export const healthCheck = async (): Promise<{ status: string; message: string }> => {
  const url = `${API_BASE_URL.replace('/api', '')}/health`;
  const response = await fetch(url);
  return await response.json();
};
