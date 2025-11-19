package web

import (
	"context"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

// ResponseCache provides simple in-memory caching for API responses
type ResponseCache struct {
	cache map[string]*CacheEntry
	mu    sync.RWMutex
	ttl   time.Duration
}

// CacheEntry represents a cached response
type CacheEntry struct {
	Data      interface{}
	Timestamp time.Time
}

// NewResponseCache creates a new response cache
func NewResponseCache(ttl time.Duration) *ResponseCache {
	cache := &ResponseCache{
		cache: make(map[string]*CacheEntry),
		ttl:   ttl,
	}

	// Start cleanup goroutine
	go cache.cleanupExpired()

	return cache
}

// Get retrieves a cached response
func (rc *ResponseCache) Get(key string) (interface{}, bool) {
	rc.mu.RLock()
	defer rc.mu.RUnlock()

	entry, exists := rc.cache[key]
	if !exists {
		return nil, false
	}

	// Check if expired
	if time.Since(entry.Timestamp) > rc.ttl {
		return nil, false
	}

	return entry.Data, true
}

// Set stores a response in cache
func (rc *ResponseCache) Set(key string, data interface{}) {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	rc.cache[key] = &CacheEntry{
		Data:      data,
		Timestamp: time.Now(),
	}
}

// Delete removes a cached response
func (rc *ResponseCache) Delete(key string) {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	delete(rc.cache, key)
}

// Clear removes all cached responses
func (rc *ResponseCache) Clear() {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	rc.cache = make(map[string]*CacheEntry)
}

// cleanupExpired removes expired entries periodically
func (rc *ResponseCache) cleanupExpired() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rc.mu.Lock()
		now := time.Now()
		for key, entry := range rc.cache {
			if now.Sub(entry.Timestamp) > rc.ttl {
				delete(rc.cache, key)
			}
		}
		rc.mu.Unlock()
	}
}

// CachingMiddleware provides response caching for GET requests
func CachingMiddleware(cache *ResponseCache) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Only cache GET requests
			if c.Request().Method != "GET" {
				return next(c)
			}

			// Generate cache key from URL and query params
			cacheKey := c.Request().URL.String()

			// Check cache
			if cached, found := cache.Get(cacheKey); found {
				c.Response().Header().Set("X-Cache", "HIT")
				return c.JSON(200, cached)
			}

			// Execute handler
			c.Response().Header().Set("X-Cache", "MISS")
			return next(c)
		}
	}
}

// CalculationCache provides caching specifically for calculation results
type CalculationCache struct {
	*ResponseCache
}

// NewCalculationCache creates a new calculation cache with 5 minute TTL
func NewCalculationCache() *CalculationCache {
	return &CalculationCache{
		ResponseCache: NewResponseCache(5 * time.Minute),
	}
}

// ConnectionPool manages database connection pooling
type ConnectionPool struct {
	maxConnections int
	idleTimeout    time.Duration
	maxLifetime    time.Duration
}

// NewConnectionPool creates optimized connection pool settings
func NewConnectionPool() *ConnectionPool {
	return &ConnectionPool{
		maxConnections: 25,
		idleTimeout:    5 * time.Minute,
		maxLifetime:    30 * time.Minute,
	}
}

// BatchProcessor processes items in batches for better performance
type BatchProcessor struct {
	batchSize int
	timeout   time.Duration
}

// NewBatchProcessor creates a new batch processor
func NewBatchProcessor(batchSize int, timeout time.Duration) *BatchProcessor {
	return &BatchProcessor{
		batchSize: batchSize,
		timeout:   timeout,
	}
}

// Process processes items in batches
func (bp *BatchProcessor) Process(ctx context.Context, items []interface{}, processFn func([]interface{}) error) error {
	for i := 0; i < len(items); i += bp.batchSize {
		end := i + bp.batchSize
		if end > len(items) {
			end = len(items)
		}

		batch := items[i:end]

		// Check context cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Process batch
		if err := processFn(batch); err != nil {
			return err
		}
	}

	return nil
}

// WorkerPool manages concurrent workers for parallel processing
type WorkerPool struct {
	workers int
	jobs    chan func()
	wg      sync.WaitGroup
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(workers int) *WorkerPool {
	pool := &WorkerPool{
		workers: workers,
		jobs:    make(chan func(), workers*2),
	}

	// Start workers
	for i := 0; i < workers; i++ {
		pool.wg.Add(1)
		go pool.worker()
	}

	return pool
}

// worker processes jobs from the queue
func (wp *WorkerPool) worker() {
	defer wp.wg.Done()

	for job := range wp.jobs {
		job()
	}
}

// Submit submits a job to the worker pool
func (wp *WorkerPool) Submit(job func()) {
	wp.jobs <- job
}

// Close closes the worker pool and waits for all jobs to complete
func (wp *WorkerPool) Close() {
	close(wp.jobs)
	wp.wg.Wait()
}

// PerformanceMetrics tracks API performance metrics
type PerformanceMetrics struct {
	mu              sync.RWMutex
	requestCount    int64
	totalDuration   time.Duration
	slowestRequest  time.Duration
	fastestRequest  time.Duration
	errorCount      int64
}

// NewPerformanceMetrics creates a new performance metrics tracker
func NewPerformanceMetrics() *PerformanceMetrics {
	return &PerformanceMetrics{
		fastestRequest: time.Hour, // Initialize with large value
	}
}

// RecordRequest records a request's performance
func (pm *PerformanceMetrics) RecordRequest(duration time.Duration, isError bool) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	pm.requestCount++
	pm.totalDuration += duration

	if duration > pm.slowestRequest {
		pm.slowestRequest = duration
	}

	if duration < pm.fastestRequest {
		pm.fastestRequest = duration
	}

	if isError {
		pm.errorCount++
	}
}

// GetMetrics returns current performance metrics
func (pm *PerformanceMetrics) GetMetrics() map[string]interface{} {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	avgDuration := time.Duration(0)
	if pm.requestCount > 0 {
		avgDuration = pm.totalDuration / time.Duration(pm.requestCount)
	}

	errorRate := float64(0)
	if pm.requestCount > 0 {
		errorRate = float64(pm.errorCount) / float64(pm.requestCount) * 100
	}

	return map[string]interface{}{
		"total_requests":  pm.requestCount,
		"average_latency": avgDuration.String(),
		"slowest_request": pm.slowestRequest.String(),
		"fastest_request": pm.fastestRequest.String(),
		"error_count":     pm.errorCount,
		"error_rate":      errorRate,
	}
}

// PerformanceMonitoringMiddleware tracks request performance
func PerformanceMonitoringMiddleware(metrics *PerformanceMetrics) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)

			duration := time.Since(start)
			isError := err != nil || c.Response().Status >= 400

			metrics.RecordRequest(duration, isError)

			// Add performance header
			c.Response().Header().Set("X-Response-Time", duration.String())

			return err
		}
	}
}

// CompressionConfig provides optimized compression settings
type CompressionConfig struct {
	Level     int
	MinLength int
}

// GetOptimalCompressionConfig returns optimal compression settings
func GetOptimalCompressionConfig() CompressionConfig {
	return CompressionConfig{
		Level:     5, // Balance between speed and compression ratio
		MinLength: 1024, // Only compress responses larger than 1KB
	}
}
