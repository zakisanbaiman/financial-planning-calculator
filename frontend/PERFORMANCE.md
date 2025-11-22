# Frontend Performance Optimization Guide

## Overview

This document outlines performance optimization strategies implemented in the Financial Planning Calculator frontend.

## Code Splitting

### Route-Based Splitting

Next.js automatically splits code by route. Each page is loaded on demand:

```typescript
// Automatic code splitting by Next.js
// app/dashboard/page.tsx - Only loaded when visiting /dashboard
// app/calculations/page.tsx - Only loaded when visiting /calculations
```

### Component-Based Splitting

For large components, use dynamic imports:

```typescript
import dynamic from 'next/dynamic';

// Lazy load heavy components
const AssetProjectionChart = dynamic(
  () => import('@/components/AssetProjectionChart'),
  {
    loading: () => <LoadingSpinner />,
    ssr: false, // Disable SSR for client-only components
  }
);
```

### Library Splitting

Import only what you need:

```typescript
// Bad
import _ from 'lodash';

// Good
import debounce from 'lodash/debounce';
```

## React Optimization

### 1. Memoization

Use React.memo for expensive components:

```typescript
import { memo } from 'react';

const ExpensiveComponent = memo(({ data }) => {
  // Expensive rendering logic
  return <div>{/* ... */}</div>;
});
```

### 2. useMemo for Expensive Calculations

```typescript
import { useMemo } from 'react';

function CalculationComponent({ data }) {
  const result = useMemo(() => {
    return expensiveCalculation(data);
  }, [data]);

  return <div>{result}</div>;
}
```

### 3. useCallback for Event Handlers

```typescript
import { useCallback } from 'react';

function FormComponent({ onSubmit }) {
  const handleSubmit = useCallback((e) => {
    e.preventDefault();
    onSubmit(formData);
  }, [formData, onSubmit]);

  return <form onSubmit={handleSubmit}>{/* ... */}</form>;
}
```

### 4. Virtualization for Long Lists

For lists with many items:

```typescript
import { FixedSizeList } from 'react-window';

function GoalsList({ goals }) {
  return (
    <FixedSizeList
      height={600}
      itemCount={goals.length}
      itemSize={80}
      width="100%"
    >
      {({ index, style }) => (
        <div style={style}>
          <GoalItem goal={goals[index]} />
        </div>
      )}
    </FixedSizeList>
  );
}
```

## Data Fetching Optimization

### 1. Request Deduplication

Use SWR or React Query for automatic deduplication:

```typescript
import useSWR from 'swr';

function useFinancialData(userId: string) {
  const { data, error } = useSWR(
    `/api/financial-data?user_id=${userId}`,
    fetcher,
    {
      revalidateOnFocus: false,
      dedupingInterval: 5000, // Dedupe requests within 5 seconds
    }
  );

  return { data, error, isLoading: !error && !data };
}
```

### 2. Prefetching

Prefetch data for likely next pages:

```typescript
import { useRouter } from 'next/navigation';

function Navigation() {
  const router = useRouter();

  return (
    <Link
      href="/dashboard"
      onMouseEnter={() => router.prefetch('/dashboard')}
    >
      Dashboard
    </Link>
  );
}
```

### 3. Caching

Implement client-side caching:

```typescript
import { apiCache, withCache } from '@/lib/integration-utils';

async function fetchGoals(userId: string) {
  return withCache(
    `goals:${userId}`,
    () => goalsAPI.list(userId),
    true // Use cache
  );
}
```

## Chart Optimization

### 1. Data Point Reduction

Reduce data points for large datasets:

```typescript
import { optimizeChartData } from '@/lib/performance';

function AssetChart({ data }) {
  const optimizedData = useMemo(
    () => optimizeChartData(data, 100), // Max 100 points
    [data]
  );

  return <LineChart data={optimizedData} />;
}
```

### 2. Debounced Updates

Debounce chart updates during user input:

```typescript
import { debounce } from '@/lib/performance';

function InteractiveChart({ onDataChange }) {
  const debouncedUpdate = useMemo(
    () => debounce(onDataChange, 300),
    [onDataChange]
  );

  return <Chart onChange={debouncedUpdate} />;
}
```

### 3. Canvas vs SVG

Use Canvas for large datasets (>1000 points):

```typescript
// For small datasets (< 1000 points)
<LineChart renderer="svg" />

// For large datasets (> 1000 points)
<LineChart renderer="canvas" />
```

## Image Optimization

### 1. Next.js Image Component

Always use Next.js Image component:

```typescript
import Image from 'next/image';

function Logo() {
  return (
    <Image
      src="/logo.png"
      alt="Logo"
      width={200}
      height={50}
      priority // Load immediately for above-the-fold images
    />
  );
}
```

### 2. Lazy Loading

Lazy load below-the-fold images:

```typescript
<Image
  src="/chart.png"
  alt="Chart"
  width={800}
  height={600}
  loading="lazy"
/>
```

## Bundle Size Optimization

### 1. Analyze Bundle

```bash
# Analyze bundle size
npm run build
npx @next/bundle-analyzer
```

### 2. Tree Shaking

Ensure imports support tree shaking:

```typescript
// Good - Tree shakeable
import { Button } from '@/components';

// Bad - Imports entire library
import * as Components from '@/components';
```

### 3. Remove Unused Dependencies

```bash
# Find unused dependencies
npx depcheck

# Remove unused packages
npm uninstall unused-package
```

## CSS Optimization

### 1. Tailwind CSS Purging

Tailwind automatically purges unused CSS in production:

```javascript
// tailwind.config.js
module.exports = {
  content: [
    './src/**/*.{js,ts,jsx,tsx}',
  ],
  // Unused classes are automatically removed
};
```

### 2. Critical CSS

Next.js automatically inlines critical CSS.

### 3. CSS Modules

Use CSS modules for component-specific styles:

```typescript
import styles from './Component.module.css';

function Component() {
  return <div className={styles.container}>Content</div>;
}
```

## Runtime Performance

### 1. Avoid Inline Functions

```typescript
// Bad
<button onClick={() => handleClick(id)}>Click</button>

// Good
const handleButtonClick = useCallback(() => handleClick(id), [id]);
<button onClick={handleButtonClick}>Click</button>
```

### 2. Debounce Input Handlers

```typescript
import { debounce } from '@/lib/performance';

function SearchInput() {
  const debouncedSearch = useMemo(
    () => debounce((value) => performSearch(value), 300),
    []
  );

  return <input onChange={(e) => debouncedSearch(e.target.value)} />;
}
```

### 3. Throttle Scroll Handlers

```typescript
import { throttle } from '@/lib/performance';

function ScrollComponent() {
  const throttledScroll = useMemo(
    () => throttle(() => handleScroll(), 100),
    []
  );

  useEffect(() => {
    window.addEventListener('scroll', throttledScroll);
    return () => window.removeEventListener('scroll', throttledScroll);
  }, [throttledScroll]);
}
```

## Web Vitals Targets

### Core Web Vitals

- **LCP (Largest Contentful Paint)**: < 2.5s
- **FID (First Input Delay)**: < 100ms
- **CLS (Cumulative Layout Shift)**: < 0.1

### Additional Metrics

- **FCP (First Contentful Paint)**: < 1.8s
- **TTI (Time to Interactive)**: < 3.8s
- **TBT (Total Blocking Time)**: < 200ms

## Monitoring

### 1. Web Vitals Reporting

```typescript
// app/layout.tsx
import { reportWebVitals } from '@/lib/performance';

export function reportWebVitals(metric) {
  console.log(metric);
  // Send to analytics service
}
```

### 2. Performance Observer

```typescript
import { monitorLongTasks } from '@/lib/performance';

useEffect(() => {
  monitorLongTasks();
}, []);
```

### 3. Lighthouse CI

```bash
# Run Lighthouse in CI
npm install -g @lhci/cli
lhci autorun
```

## Optimization Checklist

### Build Time

- [ ] Code splitting is enabled
- [ ] Tree shaking is working
- [ ] Bundle size is analyzed
- [ ] Unused dependencies are removed
- [ ] Images are optimized

### Runtime

- [ ] Components are memoized where appropriate
- [ ] Expensive calculations use useMemo
- [ ] Event handlers use useCallback
- [ ] Long lists are virtualized
- [ ] Charts are optimized

### Network

- [ ] API responses are cached
- [ ] Requests are deduplicated
- [ ] Critical resources are prefetched
- [ ] Compression is enabled
- [ ] CDN is configured

### Rendering

- [ ] No unnecessary re-renders
- [ ] Layout shifts are minimized
- [ ] Critical CSS is inlined
- [ ] Fonts are optimized
- [ ] Images are lazy loaded

## Performance Budget

### JavaScript

- **Initial bundle**: < 200KB (gzipped)
- **Per route**: < 50KB (gzipped)
- **Total**: < 500KB (gzipped)

### CSS

- **Critical CSS**: < 14KB
- **Total CSS**: < 50KB (gzipped)

### Images

- **Hero images**: < 200KB
- **Thumbnails**: < 50KB
- **Icons**: Use SVG or icon fonts

### API Responses

- **Simple queries**: < 100KB
- **Calculations**: < 200KB
- **Reports**: < 500KB

## Tools

### Development

- **React DevTools**: Profile component renders
- **Chrome DevTools**: Performance profiling
- **Lighthouse**: Performance audits
- **Bundle Analyzer**: Analyze bundle size

### Production

- **Vercel Analytics**: Real user monitoring
- **Sentry**: Error tracking with performance data
- **Google Analytics**: User behavior tracking

## Resources

- [Next.js Performance](https://nextjs.org/docs/advanced-features/measuring-performance)
- [React Performance](https://react.dev/learn/render-and-commit)
- [Web Vitals](https://web.dev/vitals/)
- [Bundle Size Optimization](https://web.dev/reduce-javascript-payloads-with-code-splitting/)
