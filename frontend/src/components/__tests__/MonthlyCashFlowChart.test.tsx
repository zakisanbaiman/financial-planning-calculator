import React from 'react';
import { render, screen } from '@testing-library/react';
import '@testing-library/jest-dom';
import MonthlyCashFlowChart from '../MonthlyCashFlowChart';

// Mock Chart.js
jest.mock('react-chartjs-2', () => ({
  Line: () => <div data-testid="line-chart">Chart</div>,
}));

describe('MonthlyCashFlowChart', () => {
  const mockProps = {
    monthlyIncome: 400000,
    monthlyExpenses: 280000,
    monthlySavings: 120000,
  };

  it('renders the chart with cash flow data', () => {
    render(<MonthlyCashFlowChart {...mockProps} />);
    
    expect(screen.getByTestId('line-chart')).toBeInTheDocument();
  });

  it('displays average income, expense, and savings', () => {
    render(<MonthlyCashFlowChart {...mockProps} />);
    
    expect(screen.getByText('平均収入')).toBeInTheDocument();
    expect(screen.getByText('平均支出')).toBeInTheDocument();
    expect(screen.getByText('平均貯蓄')).toBeInTheDocument();
  });

  it('shows empty state when no data', () => {
    render(<MonthlyCashFlowChart monthlyIncome={0} monthlyExpenses={0} monthlySavings={0} />);
    
    expect(screen.getByText('収支データがありません')).toBeInTheDocument();
  });

  it('renders with custom month count', () => {
    render(<MonthlyCashFlowChart {...mockProps} months={6} />);
    
    expect(screen.getByTestId('line-chart')).toBeInTheDocument();
  });
});
