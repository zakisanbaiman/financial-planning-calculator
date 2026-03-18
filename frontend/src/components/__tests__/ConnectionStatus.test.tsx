import React from 'react';
import { render, screen, act, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { ConnectionStatus, InlineConnectionStatus } from '../ConnectionStatus';

// integration-utils モック
jest.mock('@/lib/integration-utils', () => ({
  checkAPIHealth: jest.fn(),
  checkAPIReadiness: jest.fn(),
}));

import { checkAPIHealth, checkAPIReadiness } from '@/lib/integration-utils';
const mockedCheckAPIHealth = checkAPIHealth as jest.MockedFunction<typeof checkAPIHealth>;
const mockedCheckAPIReadiness = checkAPIReadiness as jest.MockedFunction<typeof checkAPIReadiness>;

describe('ConnectionStatus', () => {
  beforeEach(() => {
    jest.useFakeTimers();
    jest.clearAllMocks();
  });

  afterEach(() => {
    jest.useRealTimers();
  });

  describe('正常な接続', () => {
    it('showWhenHealthy=false の場合、正常時は何も表示しない', async () => {
      mockedCheckAPIHealth.mockResolvedValue({
        healthy: true,
        message: 'APIサーバーは正常に動作しています',
      });

      const { container } = render(<ConnectionStatus />);

      await act(async () => {
        await Promise.resolve();
      });

      expect(container.firstChild).toBeNull();
    });

    it('showWhenHealthy=true の場合、正常時にバナーが表示される', async () => {
      mockedCheckAPIHealth.mockResolvedValue({
        healthy: true,
        message: 'APIサーバーは正常に動作しています',
      });

      render(<ConnectionStatus showWhenHealthy />);

      await waitFor(() => {
        expect(screen.getByText('APIサーバーは正常に動作しています')).toBeInTheDocument();
      });
    });
  });

  describe('異常な接続', () => {
    it('APIが不健全な場合、エラーバナーが表示される', async () => {
      mockedCheckAPIHealth.mockResolvedValue({
        healthy: false,
        message: 'APIサーバーに接続できません',
      });

      render(<ConnectionStatus />);

      await waitFor(() => {
        expect(screen.getByText('APIサーバーに接続できません')).toBeInTheDocument();
      });
    });

    it('再確認ボタンが表示される', async () => {
      mockedCheckAPIHealth.mockResolvedValue({
        healthy: false,
        message: 'APIサーバーに接続できません',
      });

      render(<ConnectionStatus />);

      await waitFor(() => {
        expect(screen.getByText('再確認')).toBeInTheDocument();
      });
    });

    it('再確認ボタンをクリックすると再チェックが実行される', async () => {
      const user = userEvent.setup({ advanceTimers: jest.advanceTimersByTime });
      mockedCheckAPIHealth.mockResolvedValue({
        healthy: false,
        message: 'APIサーバーに接続できません',
      });

      render(<ConnectionStatus />);

      await waitFor(() => {
        expect(screen.getByText('再確認')).toBeInTheDocument();
      });

      // 2回目の呼び出しで正常に
      mockedCheckAPIHealth.mockResolvedValue({
        healthy: true,
        message: 'APIサーバーは正常に動作しています',
      });

      await user.click(screen.getByText('再確認'));

      await waitFor(() => {
        expect(mockedCheckAPIHealth).toHaveBeenCalledTimes(2);
      });
    });
  });

  describe('例外ハンドリング', () => {
    it('checkAPIHealthが例外を投げた場合、エラー表示になる', async () => {
      mockedCheckAPIHealth.mockRejectedValue(new Error('Network error'));

      render(<ConnectionStatus />);

      await waitFor(() => {
        expect(screen.getByText('APIサーバーに接続できません')).toBeInTheDocument();
      });
    });
  });

  describe('定期チェック', () => {
    it('指定間隔でヘルスチェックが実行される', async () => {
      mockedCheckAPIHealth.mockResolvedValue({
        healthy: true,
        message: 'OK',
      });

      render(<ConnectionStatus checkInterval={5000} showWhenHealthy />);

      await act(async () => {
        await Promise.resolve();
      });

      expect(mockedCheckAPIHealth).toHaveBeenCalledTimes(1);

      await act(async () => {
        jest.advanceTimersByTime(5000);
        await Promise.resolve();
      });

      expect(mockedCheckAPIHealth).toHaveBeenCalledTimes(2);
    });
  });
});

describe('InlineConnectionStatus', () => {
  beforeEach(() => {
    jest.useFakeTimers();
    jest.clearAllMocks();
  });

  afterEach(() => {
    jest.useRealTimers();
  });

  it('チェック中は「確認中...」が表示される', () => {
    mockedCheckAPIReadiness.mockImplementation(() => new Promise(() => {}));
    render(<InlineConnectionStatus />);
    expect(screen.getByText('確認中...')).toBeInTheDocument();
  });

  it('正常時は「オンライン」が表示される', async () => {
    mockedCheckAPIReadiness.mockResolvedValue(true);
    render(<InlineConnectionStatus />);

    await waitFor(() => {
      expect(screen.getByText('オンライン')).toBeInTheDocument();
    });
  });

  it('異常時は「オフライン」が表示される', async () => {
    mockedCheckAPIReadiness.mockResolvedValue(false);
    render(<InlineConnectionStatus />);

    await waitFor(() => {
      expect(screen.getByText('オフライン')).toBeInTheDocument();
    });
  });
});
