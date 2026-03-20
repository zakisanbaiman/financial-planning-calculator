import React from 'react';

export const Line = (props: any) => (
  <div data-testid="mock-line-chart" data-chart-data={JSON.stringify(props.data)}>
    Line Chart
  </div>
);

export const Bar = (props: any) => (
  <div data-testid="mock-bar-chart" data-chart-data={JSON.stringify(props.data)}>
    Bar Chart
  </div>
);

export const Doughnut = (props: any) => (
  <div data-testid="mock-doughnut-chart" data-chart-data={JSON.stringify(props.data)}>
    Doughnut Chart
  </div>
);
