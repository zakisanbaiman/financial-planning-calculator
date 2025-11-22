// Performance monitoring and optimization utilities

/**
 * Measure component render time
 */
export function measureRenderTime(componentName: string, callback: () => void) {
  if (typeof window === 'undefined' || process.env.NODE_ENV !== 'development') {
    callback();
    return;
  }

  const startTime = performance.now();
  callback();
  const endTime = performance.now();
  
  console.log(`[Performance] ${componentName} rendered in ${(endTime - startTime).toFixed(2)}ms`);
}

/**
 * Report Web Vitals
 */
export function reportWebVitals(metric: any) {
  if (process.env.NODE_ENV === 'development') {
    console.log('[Web Vitals]', metric);
  }
  
  // In production, you could send this to an analytics service
  // Example: sendToAnalytics(metric);
}

/**
 * Lazy load component with loading state
 */
export function lazyLoadComponent<T extends React.ComponentType<any>>(
  importFunc: () => Promise<{ default: T }>,
  fallback?: React.ReactNode
) {
  const React = require('react');
  return React.lazy(importFunc);
}

/**
 * Debounce function for performance optimization
 */
export function debounce<T extends (...args: any[]) => any>(
  func: T,
  wait: number
): (...args: Parameters<T>) => void {
  let timeout: NodeJS.Timeout | null = null;

  return function executedFunction(...args: Parameters<T>) {
    const later = () => {
      timeout = null;
      func(...args);
    };

    if (timeout) {
      clearTimeout(timeout);
    }
    timeout = setTimeout(later, wait);
  };
}

/**
 * Throttle function for performance optimization
 */
export function throttle<T extends (...args: any[]) => any>(
  func: T,
  limit: number
): (...args: Parameters<T>) => void {
  let inThrottle: boolean;

  return function executedFunction(...args: Parameters<T>) {
    if (!inThrottle) {
      func(...args);
      inThrottle = true;
      setTimeout(() => (inThrottle = false), limit);
    }
  };
}

/**
 * Memoize expensive calculations
 */
export function memoize<T extends (...args: any[]) => any>(fn: T): T {
  const cache = new Map();

  return ((...args: Parameters<T>) => {
    const key = JSON.stringify(args);
    
    if (cache.has(key)) {
      return cache.get(key);
    }

    const result = fn(...args);
    cache.set(key, result);
    
    return result;
  }) as T;
}

/**
 * Optimize large list rendering
 */
export function chunkArray<T>(array: T[], chunkSize: number): T[][] {
  const chunks: T[][] = [];
  for (let i = 0; i < array.length; i += chunkSize) {
    chunks.push(array.slice(i, i + chunkSize));
  }
  return chunks;
}

/**
 * Request idle callback wrapper
 */
export function runWhenIdle(callback: () => void, options?: IdleRequestOptions) {
  if (typeof window === 'undefined') {
    callback();
    return;
  }

  if ('requestIdleCallback' in window) {
    window.requestIdleCallback(callback, options);
  } else {
    setTimeout(callback, 1);
  }
}

/**
 * Preload critical resources
 */
export function preloadResource(href: string, as: string) {
  if (typeof document === 'undefined') return;

  const link = document.createElement('link');
  link.rel = 'preload';
  link.href = href;
  link.as = as;
  document.head.appendChild(link);
}

/**
 * Monitor long tasks
 */
export function monitorLongTasks() {
  if (typeof window === 'undefined' || !('PerformanceObserver' in window)) {
    return;
  }

  try {
    const observer = new PerformanceObserver((list) => {
      for (const entry of list.getEntries()) {
        if (entry.duration > 50) {
          console.warn('[Performance] Long task detected:', {
            duration: entry.duration,
            startTime: entry.startTime,
          });
        }
      }
    });

    observer.observe({ entryTypes: ['longtask'] });
  } catch (e) {
    // PerformanceObserver not supported
  }
}

/**
 * Get performance metrics
 */
export function getPerformanceMetrics() {
  if (typeof window === 'undefined' || !window.performance) {
    return null;
  }

  const navigation = performance.getEntriesByType('navigation')[0] as PerformanceNavigationTiming;
  
  if (!navigation) return null;

  return {
    // Page load metrics
    domContentLoaded: navigation.domContentLoadedEventEnd - navigation.domContentLoadedEventStart,
    loadComplete: navigation.loadEventEnd - navigation.loadEventStart,
    
    // Network metrics
    dns: navigation.domainLookupEnd - navigation.domainLookupStart,
    tcp: navigation.connectEnd - navigation.connectStart,
    request: navigation.responseStart - navigation.requestStart,
    response: navigation.responseEnd - navigation.responseStart,
    
    // Rendering metrics
    domInteractive: navigation.domInteractive - navigation.fetchStart,
    domComplete: navigation.domComplete - navigation.fetchStart,
    
    // Total time
    totalTime: navigation.loadEventEnd - navigation.fetchStart,
  };
}

/**
 * Log performance metrics
 */
export function logPerformanceMetrics() {
  if (process.env.NODE_ENV !== 'development') return;

  runWhenIdle(() => {
    const metrics = getPerformanceMetrics();
    if (metrics) {
      console.table(metrics);
    }
  });
}

/**
 * Optimize chart data for rendering
 * Reduces data points for better performance while maintaining visual accuracy
 */
export function optimizeChartData<T extends { x: number; y: number }>(
  data: T[],
  maxPoints: number = 100
): T[] {
  if (data.length <= maxPoints) {
    return data;
  }

  const step = Math.ceil(data.length / maxPoints);
  const optimized: T[] = [];

  for (let i = 0; i < data.length; i += step) {
    optimized.push(data[i]);
  }

  // Always include the last point
  if (optimized[optimized.length - 1] !== data[data.length - 1]) {
    optimized.push(data[data.length - 1]);
  }

  return optimized;
}

/**
 * Format large numbers efficiently
 */
export function formatLargeNumber(num: number): string {
  if (num >= 100000000) {
    return `${(num / 100000000).toFixed(1)}億`;
  }
  if (num >= 10000) {
    return `${(num / 10000).toFixed(1)}万`;
  }
  return num.toLocaleString('ja-JP');
}

/**
 * Batch state updates
 */
export function batchUpdates(updates: Array<() => void>) {
  if (typeof window === 'undefined') {
    updates.forEach(update => update());
    return;
  }

  requestAnimationFrame(() => {
    updates.forEach(update => update());
  });
}
