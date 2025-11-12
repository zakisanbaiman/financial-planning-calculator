'use client';

import { Line } from 'react-chartjs-2';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler,
} from 'chart.js';
import type { AssetProjectionPoint } from '@/types/api';

// Chart.js登録
ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler
);

interface AssetProjectionChartProps {
  projections: AssetProjectionPoint[];
  showRealValue?: boolean;
  showContributions?: boolean;
  height?: number;
}

export default function AssetProjectionChart({
  projections,
  showRealValue = true,
  showContributions = true,
  height = 300,
}: AssetProjectionChartProps) {
  if (!projections || projections.length === 0) {
    return (
      <div className="h-64 bg-gray-100 rounded-lg flex items-center justify-center">
        <div className="text-center text-gray-500">
          <div className="text-4xl mb-2">📊</div>
          <p>データがありません</p>
        </div>
      </div>
    );
  }

  const labels = projections.map((p) => `${p.year}年後`);

  const datasets = [
    {
      label: '総資産',
      data: projections.map((p) => p.total_assets),
      borderColor: 'rgb(59, 130, 246)',
      backgroundColor: 'rgba(59, 130, 246, 0.1)',
      fill: true,
      tension: 0.4,
    },
  ];

  if (showRealValue) {
    datasets.push({
      label: '実質価値（インフレ調整後）',
      data: projections.map((p) => p.real_value),
      borderColor: 'rgb(16, 185, 129)',
      backgroundColor: 'rgba(16, 185, 129, 0.1)',
      fill: true,
      tension: 0.4,
    });
  }

  if (showContributions) {
    datasets.push({
      label: '積立元本',
      data: projections.map((p) => p.contributed_amount),
      borderColor: 'rgb(251, 146, 60)',
      backgroundColor: 'rgba(251, 146, 60, 0.1)',
      fill: true,
      tension: 0.4,
    });
  }

  const options = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        position: 'top' as const,
      },
      tooltip: {
        callbacks: {
          label: function (context: any) {
            let label = context.dataset.label || '';
            if (label) {
              label += ': ';
            }
            if (context.parsed.y !== null) {
              label += new Intl.NumberFormat('ja-JP', {
                style: 'currency',
                currency: 'JPY',
                maximumFractionDigits: 0,
              }).format(context.parsed.y);
            }
            return label;
          },
        },
      },
    },
    scales: {
      y: {
        beginAtZero: true,
        ticks: {
          callback: function (value: any) {
            return new Intl.NumberFormat('ja-JP', {
              style: 'currency',
              currency: 'JPY',
              maximumFractionDigits: 0,
              notation: 'compact',
            }).format(value);
          },
        },
      },
    },
  };

  return (
    <div style={{ height: `${height}px` }}>
      <Line data={{ labels, datasets }} options={options} />
    </div>
  );
}
