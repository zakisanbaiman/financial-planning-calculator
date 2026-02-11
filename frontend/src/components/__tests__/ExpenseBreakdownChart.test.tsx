import React from 'react';
import { render, screen } from '@testing-library/react';
import '@testing-library/jest-dom';
import ExpenseBreakdownChart from '../ExpenseBreakdownChart';
import type { ExpenseItem } from '@/types/api';

// Mock Chart.js
jest.mock('react-chartjs-2', () => ({
  Doughnut: () => <div data-testid="doughnut-chart">Chart</div>,
}));

describe('ExpenseBreakdownChart', () => {
  const mockExpenses: ExpenseItem[] = [
    { category: '食費', amount: 50000 },
    { category: '住居費', amount: 80000 },
    { category: '交通費', amount: 20000 },
  ];

  it('renders the chart with expense data', () => {
    render(<ExpenseBreakdownChart expenses={mockExpenses} />);
    
    expect(screen.getByText('支出内訳')).toBeInTheDocument();
    expect(screen.getByTestId('doughnut-chart')).toBeInTheDocument();
  });

  it('displays total expenses', () => {
    render(<ExpenseBreakdownChart expenses={mockExpenses} />);
    
    const totalExpenses = mockExpenses.reduce((sum, e) => sum + e.amount, 0);
    expect(screen.getByText(/¥150,000/)).toBeInTheDocument();
  });

  it('displays each category with amount', () => {
    render(<ExpenseBreakdownChart expenses={mockExpenses} />);
    
    expect(screen.getByText('食費')).toBeInTheDocument();
    expect(screen.getByText('住居費')).toBeInTheDocument();
    expect(screen.getByText('交通費')).toBeInTheDocument();
    expect(screen.getByText(/¥50,000/)).toBeInTheDocument();
    expect(screen.getByText(/¥80,000/)).toBeInTheDocument();
    expect(screen.getByText(/¥20,000/)).toBeInTheDocument();
  });

  it('shows empty state when no expenses', () => {
    render(<ExpenseBreakdownChart expenses={[]} />);
    
    expect(screen.getByText('支出データがありません')).toBeInTheDocument();
  });
});
