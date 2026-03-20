import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { ErrorBoundary, APIErrorDisplay } from '../ErrorBoundary';

// エラーを投げるテスト用コンポーネント
const ThrowError = ({ shouldThrow }: { shouldThrow: boolean }) => {
  if (shouldThrow) {
    throw new Error('テストエラー');
  }
  return <div>正常なコンテンツ</div>;
};

describe('ErrorBoundary', () => {
  // console.errorを抑制（React内部のエラーログ）
  const originalConsoleError = console.error;
  beforeAll(() => {
    console.error = jest.fn();
  });
  afterAll(() => {
    console.error = originalConsoleError;
  });

  describe('正常時', () => {
    it('子コンポーネントが正常にレンダリングされる', () => {
      render(
        <ErrorBoundary>
          <ThrowError shouldThrow={false} />
        </ErrorBoundary>
      );
      expect(screen.getByText('正常なコンテンツ')).toBeInTheDocument();
    });
  });

  describe('エラーキャッチ', () => {
    it('子コンポーネントがエラーを投げた場合、デフォルトフォールバックが表示される', () => {
      render(
        <ErrorBoundary>
          <ThrowError shouldThrow={true} />
        </ErrorBoundary>
      );
      expect(screen.getByText('エラーが発生しました')).toBeInTheDocument();
      expect(
        screen.getByText(/申し訳ございません。予期しないエラーが発生しました/)
      ).toBeInTheDocument();
    });

    it('カスタムfallbackが指定された場合、それが表示される', () => {
      render(
        <ErrorBoundary fallback={<div>カスタムエラー表示</div>}>
          <ThrowError shouldThrow={true} />
        </ErrorBoundary>
      );
      expect(screen.getByText('カスタムエラー表示')).toBeInTheDocument();
    });

    it('onErrorコールバックが呼ばれる', () => {
      const onError = jest.fn();
      render(
        <ErrorBoundary onError={onError}>
          <ThrowError shouldThrow={true} />
        </ErrorBoundary>
      );
      expect(onError).toHaveBeenCalledWith(
        expect.any(Error),
        expect.objectContaining({ componentStack: expect.any(String) })
      );
    });
  });

  describe('リセット機能', () => {
    it('「もう一度試す」ボタンでエラー状態がリセットされる', async () => {
      const user = userEvent.setup();

      // 最初はエラーを投げ、リセット後は正常にレンダリング
      let shouldThrow = true;
      const TestComponent = () => {
        if (shouldThrow) throw new Error('テストエラー');
        return <div>復旧しました</div>;
      };

      render(
        <ErrorBoundary>
          <TestComponent />
        </ErrorBoundary>
      );

      expect(screen.getByText('エラーが発生しました')).toBeInTheDocument();

      // リセット後はエラーを投げないように設定
      shouldThrow = false;
      await user.click(screen.getByText('もう一度試す'));

      expect(screen.getByText('復旧しました')).toBeInTheDocument();
    });
  });

  describe('ホームに戻るボタン', () => {
    it('「ホームに戻る」ボタンが表示される', () => {
      render(
        <ErrorBoundary>
          <ThrowError shouldThrow={true} />
        </ErrorBoundary>
      );
      expect(screen.getByText('ホームに戻る')).toBeInTheDocument();
    });
  });
});

describe('APIErrorDisplay', () => {
  it('エラーメッセージが表示される', () => {
    const error = new Error('API接続エラー');
    render(<APIErrorDisplay error={error} />);
    expect(screen.getByText('API接続エラー')).toBeInTheDocument();
    expect(screen.getByText('エラーが発生しました')).toBeInTheDocument();
  });

  it('ネットワークエラーの場合、ネットワークエラーというタイトルが表示される', () => {
    const error = new Error('ネットワークに接続できません');
    render(<APIErrorDisplay error={error} />);
    expect(screen.getByText('ネットワークエラー')).toBeInTheDocument();
  });

  it('再試行ボタンが表示され、クリックするとonRetryが呼ばれる', async () => {
    const user = userEvent.setup();
    const onRetry = jest.fn();
    render(<APIErrorDisplay error={new Error('エラー')} onRetry={onRetry} />);
    await user.click(screen.getByText('再試行'));
    expect(onRetry).toHaveBeenCalledTimes(1);
  });

  it('閉じるボタンが表示され、クリックするとonDismissが呼ばれる', async () => {
    const user = userEvent.setup();
    const onDismiss = jest.fn();
    render(<APIErrorDisplay error={new Error('エラー')} onDismiss={onDismiss} />);
    await user.click(screen.getByText('閉じる'));
    expect(onDismiss).toHaveBeenCalledTimes(1);
  });

  it('onRetry/onDismissが未指定の場合、ボタンが表示されない', () => {
    render(<APIErrorDisplay error={new Error('エラー')} />);
    expect(screen.queryByText('再試行')).not.toBeInTheDocument();
    expect(screen.queryByText('閉じる')).not.toBeInTheDocument();
  });
});
