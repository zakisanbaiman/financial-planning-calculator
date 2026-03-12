import { render, screen } from '@testing-library/react';
import CalculationsPage from '../page';

// useUser フックをモック
jest.mock('@/lib/hooks/useUser', () => ({
  useUser: jest.fn(),
}));

// 重いコンポーネントをモック
jest.mock('@/components', () => ({
  AssetProjectionCalculator: () => <div data-testid="asset-projection-calculator" />,
  RetirementCalculator: () => <div data-testid="retirement-calculator" />,
  EmergencyFundCalculator: () => <div data-testid="emergency-fund-calculator" />,
}));

jest.mock('@/components/AssetProjectionChart', () => ({
  __esModule: true,
  default: () => <div data-testid="asset-projection-chart" />,
}));

jest.mock('@/lib/utils/projections', () => ({
  generateAssetProjections: jest.fn(() => Array(31).fill({ total_assets: 0 })),
}));

jest.mock('@/lib/utils/currency', () => ({
  formatCurrency: (v: number) => `¥${v.toLocaleString()}`,
}));

import { useUser } from '@/lib/hooks/useUser';

const mockUseUser = useUser as jest.MockedFunction<typeof useUser>;

describe('CalculationsPage', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('ローディング中はローディングメッセージを表示する', () => {
    mockUseUser.mockReturnValue({
      userId: null,
      email: null,
      loading: true,
      clearUser: jest.fn(),
      isGuest: false,
    });

    render(<CalculationsPage />);

    expect(screen.getByText('読み込み中...')).toBeInTheDocument();
  });

  it('未ログイン時はエラーメッセージを表示する', () => {
    mockUseUser.mockReturnValue({
      userId: null,
      email: null,
      loading: false,
      clearUser: jest.fn(),
      isGuest: false,
    });

    render(<CalculationsPage />);

    expect(screen.getByText('ログインが必要です')).toBeInTheDocument();
    expect(screen.getByText('計算機能を使用するにはログインしてください。')).toBeInTheDocument();
  });

  it('ログイン済みの場合は計算メニューを表示する', () => {
    mockUseUser.mockReturnValue({
      userId: 'test-user-id',
      email: 'test@example.com',
      loading: false,
      clearUser: jest.fn(),
      isGuest: false,
    });

    render(<CalculationsPage />);

    expect(screen.getByText('財務計算機')).toBeInTheDocument();
    expect(screen.getByText('資産推移シミュレーション')).toBeInTheDocument();
    expect(screen.getByText('老後資金計算')).toBeInTheDocument();
    expect(screen.getByText('緊急資金計算')).toBeInTheDocument();
  });

  it('ゲストモードの場合も計算メニューを表示する', () => {
    mockUseUser.mockReturnValue({
      userId: 'guest',
      email: null,
      loading: false,
      clearUser: jest.fn(),
      isGuest: true,
    });

    render(<CalculationsPage />);

    expect(screen.getByText('財務計算機')).toBeInTheDocument();
  });
});
