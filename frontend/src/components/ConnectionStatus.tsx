'use client';

import React, { useEffect, useState } from 'react';
import { checkAPIHealth, checkAPIReadiness } from '@/lib/integration-utils';

interface ConnectionStatusProps {
  showWhenHealthy?: boolean;
  checkInterval?: number;
}

/**
 * ConnectionStatus - Displays API connection status
 * Shows a banner when API is unavailable
 */
export function ConnectionStatus({ 
  showWhenHealthy = false, 
  checkInterval = 30000 
}: ConnectionStatusProps) {
  const [isHealthy, setIsHealthy] = useState<boolean>(true);
  const [isChecking, setIsChecking] = useState<boolean>(false);
  const [lastCheck, setLastCheck] = useState<Date | null>(null);
  const [errorMessage, setErrorMessage] = useState<string>('');

  const checkConnection = async () => {
    setIsChecking(true);
    try {
      const health = await checkAPIHealth();
      setIsHealthy(health.healthy);
      setErrorMessage(health.message);
      setLastCheck(new Date());
    } catch (error) {
      setIsHealthy(false);
      setErrorMessage('APIサーバーに接続できません');
      setLastCheck(new Date());
    } finally {
      setIsChecking(false);
    }
  };

  useEffect(() => {
    // Initial check
    checkConnection();

    // Set up periodic checks
    const interval = setInterval(checkConnection, checkInterval);

    return () => clearInterval(interval);
  }, [checkInterval]);

  // Don't show anything if healthy and showWhenHealthy is false
  if (isHealthy && !showWhenHealthy) {
    return null;
  }

  return (
    <div
      className={`fixed top-0 left-0 right-0 z-50 ${
        isHealthy ? 'bg-green-500' : 'bg-red-500'
      } text-white px-4 py-2 shadow-lg`}
    >
      <div className="max-w-7xl mx-auto flex items-center justify-between">
        <div className="flex items-center gap-3">
          {isChecking ? (
            <svg
              className="animate-spin h-5 w-5"
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
            >
              <circle
                className="opacity-25"
                cx="12"
                cy="12"
                r="10"
                stroke="currentColor"
                strokeWidth="4"
              />
              <path
                className="opacity-75"
                fill="currentColor"
                d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
              />
            </svg>
          ) : isHealthy ? (
            <svg
              className="h-5 w-5"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
          ) : (
            <svg
              className="h-5 w-5"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
          )}
          
          <span className="text-sm font-medium">
            {isChecking
              ? '接続確認中...'
              : isHealthy
              ? 'APIサーバーは正常に動作しています'
              : errorMessage}
          </span>
          
          {lastCheck && (
            <span className="text-xs opacity-75">
              (最終確認: {lastCheck.toLocaleTimeString('ja-JP')})
            </span>
          )}
        </div>

        {!isHealthy && (
          <button
            onClick={checkConnection}
            disabled={isChecking}
            className="text-sm font-medium hover:underline disabled:opacity-50"
          >
            再確認
          </button>
        )}
      </div>
    </div>
  );
}

/**
 * InlineConnectionStatus - Compact connection status indicator
 */
export function InlineConnectionStatus() {
  const [isHealthy, setIsHealthy] = useState<boolean | null>(null);

  useEffect(() => {
    const checkStatus = async () => {
      const ready = await checkAPIReadiness();
      setIsHealthy(ready);
    };

    checkStatus();
    const interval = setInterval(checkStatus, 10000);

    return () => clearInterval(interval);
  }, []);

  if (isHealthy === null) {
    return (
      <div className="flex items-center gap-2 text-sm text-gray-500">
        <div className="w-2 h-2 rounded-full bg-gray-400 animate-pulse" />
        <span>確認中...</span>
      </div>
    );
  }

  return (
    <div className="flex items-center gap-2 text-sm">
      <div
        className={`w-2 h-2 rounded-full ${
          isHealthy ? 'bg-green-500' : 'bg-red-500'
        }`}
      />
      <span className={isHealthy ? 'text-green-700' : 'text-red-700'}>
        {isHealthy ? 'オンライン' : 'オフライン'}
      </span>
    </div>
  );
}
