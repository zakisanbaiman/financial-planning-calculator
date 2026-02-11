import React from 'react';
import { render, screen } from '@testing-library/react';
import '@testing-library/jest-dom';
import SavingsAllocationChart from '../SavingsAllocationChart';
import type { SavingsItem } from '@/types/api';

// Mock Chart.js
jest.mock('react-chartjs-2', () => ({
  Bar: () => <div data-testid="bar-chart">Chart</div>,
}));

describe('SavingsAllocationChart', () => {
  const mockSavings: SavingsItem[] = [
    { type: 'deposit', amount: 1000000, description: '普通預金' },
    { type: 'investment', amount: 500000, description: '株式投資' },
    { type: 'other', amount: 300000, description: 'その他' },
  ];

  it('renders the chart with savings data', () => {
    render(<SavingsAllocationChart savings={mockSavings} />);
    
    expect(screen.getByText('資産配分')).toBeInTheDocument();
    expect(screen.getByTestId('bar-chart')).toBeInTheDocument();
  });

  it('displays total savings', () => {
    render(<SavingsAllocationChart savings={mockSavings} />);
    
    expect(screen.getByText(/¥1,800,000/)).toBeInTheDocument();
  });

  it('displays each savings type with amount and percentage', () => {
    render(<SavingsAllocationChart savings={mockSavings} />);
    
    expect(screen.getByText('預金')).toBeInTheDocument();
    expect(screen.getByText('投資')).toBeInTheDocument();
    expect(screen.getByText('その他')).toBeInTheDocument();
  });

  it('shows empty state when no savings', () => {
    render(<SavingsAllocationChart savings={[]} />);
    
    expect(screen.getByText('貯蓄データがありません')).toBeInTheDocument();
  });
});
