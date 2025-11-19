// Integration utilities for frontend-backend communication

import { APIError } from './api-client';

// Retry configuration
export interface RetryConfig {
  maxRetries: number;
  retryDelay: number;
  retryableStatuses: number[];
}

const DEFAULT_RETRY_CONFIG: RetryConfig = {
  maxRetries: 3,
  retryDelay: 1000,
  retryableStatuses: [408, 429, 500, 502, 503, 504],
};

// Retry wrapper for API calls
export async function withRetry<T>(
  fn: () => Promise<T>,
  config: Partial<RetryConfig> = {}
): Promise<T> {
  const { maxRetries, retryDelay, retryableStatuses } = {
    ...DEFAULT_RETRY_CONFIG,
    ...config,
  };

  let lastError: Error | null = null;

  for (let attempt = 0; attempt <= maxRetries; attempt++) {
    try {
      return await fn();
    } catch (error) {
      lastError = error as Error;

      // Check if error is retryable
      if (error instanceof APIError && error.status) {
        if (!retryableStatuses.includes(error.status)) {
          throw error;
        }
      }

      // Don't retry on last attempt
      if (attempt === maxRetries) {
        break;
      }

      // Wait before retrying with exponential backoff
      const delay = retryDelay * Math.pow(2, attempt);
      await new Promise((resolve) => setTimeout(resolve, delay));
    }
  }

  throw lastError;
}

// Connection health check
export async function checkAPIHealth(): Promise<{
  healthy: boolean;
  message: string;
  details?: any;
}> {
  try {
    const baseUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api';
    const healthUrl = baseUrl.replace('/api', '/health/detailed');
    
    const response = await fetch(healthUrl, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    if (response.ok) {
      const data = await response.json();
      return {
        healthy: data.status === 'ok',
        message: 'APIサーバーは正常に動作しています',
        details: data,
      };
    } else {
      return {
        healthy: false,
        message: 'APIサーバーに接続できません',
      };
    }
  } catch (error) {
    return {
      healthy: false,
      message: 'APIサーバーに接続できません',
      details: error instanceof Error ? error.message : 'Unknown error',
    };
  }
}

// Check API readiness
export async function checkAPIReadiness(): Promise<boolean> {
  try {
    const baseUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api';
    const readyUrl = baseUrl.replace('/api', '/ready');
    
    const response = await fetch(readyUrl);
    if (response.ok) {
      const data = await response.json();
      return data.ready === true;
    }
    return false;
  } catch {
    return false;
  }
}

// Error message formatter
export function formatAPIError(error: unknown): string {
  if (error instanceof APIError) {
    // Check if error has Japanese message
    if (error.data?.error) {
      return error.data.error;
    }
    if (error.data?.details) {
      return typeof error.data.details === 'string' 
        ? error.data.details 
        : JSON.stringify(error.data.details);
    }
    return error.message;
  }

  if (error instanceof Error) {
    return error.message;
  }

  return 'エラーが発生しました';
}

// Validation error formatter
export function formatValidationErrors(error: unknown): Record<string, string[]> {
  if (error instanceof APIError && error.data?.validation_errors) {
    return error.data.validation_errors;
  }
  return {};
}

// Network error detector
export function isNetworkError(error: unknown): boolean {
  if (error instanceof APIError) {
    return !error.status || error.status === 0;
  }
  
  if (error instanceof Error) {
    return (
      error.message.includes('ネットワーク') ||
      error.message.includes('network') ||
      error.message.includes('fetch')
    );
  }
  
  return false;
}

// Timeout error detector
export function isTimeoutError(error: unknown): boolean {
  if (error instanceof APIError) {
    return error.status === 408 || error.status === 504;
  }
  
  if (error instanceof Error) {
    return (
      error.message.includes('timeout') ||
      error.message.includes('タイムアウト')
    );
  }
  
  return false;
}

// Rate limit error detector
export function isRateLimitError(error: unknown): boolean {
  if (error instanceof APIError) {
    return error.status === 429;
  }
  return false;
}

// Server error detector
export function isServerError(error: unknown): boolean {
  if (error instanceof APIError && error.status) {
    return error.status >= 500 && error.status < 600;
  }
  return false;
}

// Client error detector
export function isClientError(error: unknown): boolean {
  if (error instanceof APIError && error.status) {
    return error.status >= 400 && error.status < 500;
  }
  return false;
}

// Get user-friendly error message
export function getUserFriendlyErrorMessage(error: unknown): string {
  if (isNetworkError(error)) {
    return 'ネットワーク接続を確認してください';
  }
  
  if (isTimeoutError(error)) {
    return 'リクエストがタイムアウトしました。もう一度お試しください';
  }
  
  if (isRateLimitError(error)) {
    return 'リクエストが多すぎます。しばらく待ってから再度お試しください';
  }
  
  if (isServerError(error)) {
    return 'サーバーエラーが発生しました。しばらく待ってから再度お試しください';
  }
  
  return formatAPIError(error);
}

// Debounce utility for API calls
export function debounce<T extends (...args: any[]) => any>(
  func: T,
  wait: number
): (...args: Parameters<T>) => void {
  let timeout: NodeJS.Timeout | null = null;

  return function executedFunction(...args: Parameters<T>) {
    const later = () => {
      timeout = null;
      func(...args);
    };

    if (timeout) {
      clearTimeout(timeout);
    }
    timeout = setTimeout(later, wait);
  };
}

// Throttle utility for API calls
export function throttle<T extends (...args: any[]) => any>(
  func: T,
  limit: number
): (...args: Parameters<T>) => void {
  let inThrottle: boolean;

  return function executedFunction(...args: Parameters<T>) {
    if (!inThrottle) {
      func(...args);
      inThrottle = true;
      setTimeout(() => (inThrottle = false), limit);
    }
  };
}

// Cache utility for API responses
class APICache {
  private cache: Map<string, { data: any; timestamp: number }> = new Map();
  private ttl: number;

  constructor(ttlSeconds: number = 300) {
    this.ttl = ttlSeconds * 1000;
  }

  get(key: string): any | null {
    const cached = this.cache.get(key);
    if (!cached) return null;

    const now = Date.now();
    if (now - cached.timestamp > this.ttl) {
      this.cache.delete(key);
      return null;
    }

    return cached.data;
  }

  set(key: string, data: any): void {
    this.cache.set(key, {
      data,
      timestamp: Date.now(),
    });
  }

  clear(): void {
    this.cache.clear();
  }

  delete(key: string): void {
    this.cache.delete(key);
  }
}

export const apiCache = new APICache(300); // 5 minutes default TTL

// Cached API call wrapper
export async function withCache<T>(
  key: string,
  fn: () => Promise<T>,
  useCache: boolean = true
): Promise<T> {
  if (useCache) {
    const cached = apiCache.get(key);
    if (cached !== null) {
      return cached as T;
    }
  }

  const result = await fn();
  if (useCache) {
    apiCache.set(key, result);
  }
  return result;
}
