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
  AuthResponse,
  Setup2FAResponse,
  RegenerateBackupCodesResponse,
} from '@/types/api';
import { getAuthToken, isAuthTokenValid } from './contexts/AuthContext';

// API ベースURL
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api';

// リフレッシュ中フラグ（複数のリクエストが同時にリフレッシュするのを防ぐ）
let isRefreshing = false;
let refreshPromise: Promise<boolean> | null = null;

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

// リフレッシュトークンで新しいアクセストークンを取得
async function refreshAccessToken(): Promise<boolean> {
  const refreshToken = typeof window !== 'undefined' 
    ? localStorage.getItem('refresh_token')
    : null;

  if (!refreshToken) {
    return false;
  }

  try {
    const response = await fetch(`${API_BASE_URL}/auth/refresh`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ refresh_token: refreshToken }),
    });

    if (!response.ok) {
      return false;
    }

    const data = await response.json();
    
    if (typeof window !== 'undefined') {
      localStorage.setItem('auth_token', data.token);
      localStorage.setItem('auth_expires', data.expires_at);
    }

    return true;
  } catch (error) {
    console.error('トークンリフレッシュに失敗しました:', error);
    return false;
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
    
    // 401エラー（未認証）の場合、リフレッシュトークンで自動更新を試みる
    if (response.status === 401) {
      // リフレッシュ中でない場合、リフレッシュを開始
      if (!isRefreshing) {
        isRefreshing = true;
        refreshPromise = refreshAccessToken().finally(() => {
          isRefreshing = false;
          refreshPromise = null;
        });
      }

      // リフレッシュ完了を待つ
      const refreshed = await refreshPromise;

      if (refreshed) {
        // 新しいトークンでリトライ
        const newToken = getAuthToken();
        const retryConfig: RequestInit = {
          ...config,
          headers: {
            ...config.headers,
            'Authorization': `Bearer ${newToken}`,
          },
        };
        
        const retryResponse = await fetch(url, retryConfig);
        
        if (retryResponse.ok) {
          return await retryResponse.json();
        }
      }

      // リフレッシュ失敗またはリトライ失敗の場合、ログインページにリダイレクト
      if (typeof window !== 'undefined') {
        localStorage.removeItem('auth_token');
        localStorage.removeItem('auth_user');
        localStorage.removeItem('auth_expires');
        localStorage.removeItem('refresh_token');
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

// 2段階認証API
export const twoFactorAPI = {
  // 2FAステータス取得
  getStatus: async (): Promise<{ enabled: boolean }> => {
    return request<{ enabled: boolean }>('/auth/2fa/status', {
      method: 'GET',
    });
  },

  // 2FA設定開始（QRコード取得）
  setup: async (): Promise<Setup2FAResponse> => {
    return request<Setup2FAResponse>('/auth/2fa/setup', {
      method: 'POST',
    });
  },

  // 2FA有効化
  enable: async (code: string, secret: string): Promise<{ message: string }> => {
    return request<{ message: string }>('/auth/2fa/enable', {
      method: 'POST',
      body: JSON.stringify({ code, secret }),
    });
  },

  // 2FA検証（ログイン時）
  verify: async (code: string, useBackup: boolean = false): Promise<AuthResponse> => {
    return request<AuthResponse>('/auth/2fa/verify', {
      method: 'POST',
      body: JSON.stringify({ code, use_backup: useBackup }),
    });
  },

  // 2FA無効化
  disable: async (password: string): Promise<{ message: string }> => {
    return request<{ message: string }>('/auth/2fa', {
      method: 'DELETE',
      body: JSON.stringify({ password }),
    });
  },

  // バックアップコード再生成
  regenerateBackupCodes: async (): Promise<RegenerateBackupCodesResponse> => {
    return request<RegenerateBackupCodesResponse>('/auth/2fa/backup-codes', {
      method: 'POST',
    });
  },
};

// ヘルスチェック
export const healthCheck = async (): Promise<{ status: string; message: string }> => {
  const url = `${API_BASE_URL.replace('/api', '')}/health`;
  const response = await fetch(url);
  return await response.json();
};
