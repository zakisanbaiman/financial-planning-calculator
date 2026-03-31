'use client';

import { useState, useCallback } from 'react';
import { streamBotAnswer } from '@/lib/api-client';

export interface UseBotReturn {
  answer: string;
  isLoading: boolean;
  error: string | null;
  ask: (question: string) => Promise<void>;
  reset: () => void;
}

export const useBot = (): UseBotReturn => {
  const [answer, setAnswer] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const ask = useCallback(async (question: string) => {
    if (!question.trim()) {
      setError('質問を入力してください');
      return;
    }

    setIsLoading(true);
    setError(null);
    setAnswer('');

    try {
      await streamBotAnswer(
        question,
        (token) => {
          setAnswer((prev) => prev + token);
        },
        () => {
          setIsLoading(false);
        },
        (message) => {
          setError(message);
          setIsLoading(false);
        },
      );
    } catch (err) {
      const message = err instanceof Error ? err.message : 'エラーが発生しました';
      setError(message);
      setIsLoading(false);
    }
  }, []);

  const reset = useCallback(() => {
    setAnswer('');
    setError(null);
    setIsLoading(false);
  }, []);

  return { answer, isLoading, error, ask, reset };
};
