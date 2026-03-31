'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/lib/contexts/AuthContext';
import { useBot } from '@/lib/hooks/useBot';

export default function BotPage() {
  const router = useRouter();
  const { isAuthenticated, isLoading: authLoading } = useAuth();
  const [question, setQuestion] = useState('');
  const { answer, isLoading, error, ask, reset } = useBot();

  useEffect(() => {
    if (!authLoading && !isAuthenticated) {
      router.push('/login');
    }
  }, [isAuthenticated, authLoading, router]);

  if (authLoading || !isAuthenticated) {
    return null;
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!question.trim()) return;
    await ask(question);
  };

  const handleReset = () => {
    setQuestion('');
    reset();
  };

  return (
    <div className="container mx-auto px-4 py-8 max-w-3xl">
      <h1 className="text-2xl font-display font-semibold text-ink-900 dark:text-ink-100 mb-6">
        AIチャットボット
      </h1>
      <p className="text-sm text-ink-500 dark:text-ink-400 mb-6 font-body">
        財務計画に関する質問をどうぞ。FAQを参照しながら回答します。
      </p>

      <form onSubmit={handleSubmit} className="mb-6">
        <div className="flex gap-2">
          <input
            type="text"
            value={question}
            onChange={(e) => setQuestion(e.target.value)}
            placeholder="質問を入力してください（例: 積立NISAとは何ですか？）"
            className="flex-1 px-4 py-2 border border-ink-300 dark:border-ink-700 rounded bg-white dark:bg-ink-900 text-ink-900 dark:text-ink-100 text-sm font-body focus:outline-none focus:ring-2 focus:ring-ink-500"
            disabled={isLoading}
          />
          <button
            type="submit"
            disabled={isLoading || !question.trim()}
            className="px-4 py-2 bg-ink-900 dark:bg-ink-100 text-ink-50 dark:text-ink-900 text-sm font-body font-medium hover:bg-ink-700 dark:hover:bg-ink-300 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            {isLoading ? '回答中...' : '送信'}
          </button>
          {(answer || error) && (
            <button
              type="button"
              onClick={handleReset}
              className="px-4 py-2 border border-ink-300 dark:border-ink-700 text-ink-600 dark:text-ink-300 text-sm font-body hover:bg-ink-50 dark:hover:bg-ink-800 transition-colors"
            >
              リセット
            </button>
          )}
        </div>
      </form>

      {error && (
        <div className="mb-4 p-4 border border-red-300 dark:border-red-700 bg-red-50 dark:bg-red-950 text-red-700 dark:text-red-300 text-sm font-body rounded">
          {error}
        </div>
      )}

      {(answer || isLoading) && (
        <div className="p-4 border border-ink-200 dark:border-ink-800 bg-ink-50 dark:bg-ink-900 rounded">
          <h2 className="text-sm font-body font-medium text-ink-500 dark:text-ink-400 mb-2">
            回答
          </h2>
          <div className="text-sm font-body text-ink-900 dark:text-ink-100 whitespace-pre-wrap">
            {answer}
            {isLoading && (
              <span className="inline-block w-2 h-4 bg-ink-400 dark:bg-ink-500 animate-pulse ml-0.5" />
            )}
          </div>
        </div>
      )}
    </div>
  );
}
