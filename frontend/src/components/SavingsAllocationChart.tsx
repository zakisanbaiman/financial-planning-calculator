'use client';

import { Bar } from 'react-chartjs-2';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
  Tooltip,
  Legend,
} from 'chart.js';
import type { SavingsItem } from '@/types/api';

ChartJS.register(
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
  Tooltip,
  Legend
);

export interface SavingsAllocationChartProps {
  savings: SavingsItem[];
  height?: number;
}

const SavingsAllocationChart: React.FC<SavingsAllocationChartProps> = ({ 
  savings, 
  height = 300 
}) => {
  if (!savings || savings.length === 0) {
    return (
      <div className="h-64 bg-gray-100 dark:bg-gray-800 rounded-lg flex items-center justify-center">
        <div className="text-center text-gray-500 dark:text-gray-400">
          <div className="text-4xl mb-2">💰</div>
          <p>貯蓄データがありません</p>
        </div>
      </div>
    );
  }

  // タイプごとに集計
  const typeMap = new Map<string, number>();
  const typeLabels: { [key: string]: string } = {
    deposit: '預金',
    investment: '投資',
    other: 'その他',
  };

  savings.forEach((item) => {
    const current = typeMap.get(item.type) || 0;
    typeMap.set(item.type, current + item.amount);
  });

  const types = Array.from(typeMap.keys());
  const amounts = Array.from(typeMap.values());
  const labels = types.map(type => typeLabels[type] || type);

  const chartData = {
    labels,
    datasets: [
      {
        label: '資産額',
        data: amounts,
        backgroundColor: [
          'rgba(59, 130, 246, 0.8)',   // 預金 - 青
          'rgba(34, 197, 94, 0.8)',    // 投資 - 緑
          'rgba(168, 85, 247, 0.8)',   // その他 - 紫
        ].slice(0, types.length),
        borderColor: [
          'rgba(59, 130, 246, 1)',
          'rgba(34, 197, 94, 1)',
          'rgba(168, 85, 247, 1)',
        ].slice(0, types.length),
        borderWidth: 2,
      },
    ],
  };

  const options = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        display: false,
      },
      title: {
        display: true,
        text: '資産配分',
        font: {
          size: 16,
          weight: 'bold' as const,
        },
      },
      tooltip: {
        callbacks: {
          label: function (context: any) {
            const value = context.parsed.y || 0;
            const total = amounts.reduce((sum, val) => sum + val, 0);
            const percentage = total > 0 ? ((value / total) * 100).toFixed(1) : '0.0';
            return `¥${value.toLocaleString()} (${percentage}%)`;
          },
        },
      },
    },
    scales: {
      y: {
        beginAtZero: true,
        ticks: {
          callback: function (value: any) {
            return '¥' + value.toLocaleString();
          },
        },
      },
    },
  };

  const totalSavings = amounts.reduce((sum, amount) => sum + amount, 0);

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center">
        <h3 className="text-lg font-semibold text-gray-900 dark:text-white">資産配分</h3>
        <div className="text-sm font-medium text-gray-600 dark:text-gray-300">
          合計: <span className="text-primary-600 dark:text-primary-400">¥{totalSavings.toLocaleString()}</span>
        </div>
      </div>
      <div style={{ height: `${height}px` }}>
        <Bar data={chartData} options={options} />
      </div>
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
        {labels.map((label, index) => {
          const amount = amounts[index];
          const percentage = totalSavings > 0 ? ((amount / totalSavings) * 100).toFixed(1) : '0.0';
          const colors = [
            'bg-primary-500',
            'bg-success-500',
            'bg-purple-500',
          ];
          return (
            <div 
              key={label} 
              className="p-3 bg-gray-50 dark:bg-gray-800 rounded-lg"
            >
              <div className="flex items-center gap-2 mb-2">
                <div className={`w-3 h-3 rounded-full ${colors[index % colors.length]}`} />
                <span className="text-sm font-medium text-gray-700 dark:text-gray-300">{label}</span>
              </div>
              <div className="text-xl font-bold text-gray-900 dark:text-white">
                ¥{amount.toLocaleString()}
              </div>
              <div className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                全体の{percentage}%
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
};

export default SavingsAllocationChart;
