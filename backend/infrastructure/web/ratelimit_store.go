package web

import (
	"math"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// RateLimitInfo holds the current rate limit status for a given identifier.
type RateLimitInfo struct {
	Limit     int       `json:"limit"`
	Remaining int       `json:"remaining"`
	Reset     int64     `json:"reset"`     // Unix timestamp (seconds)
	ResetAt   time.Time `json:"reset_at"`  // Human-readable reset time
}

// rateLimiterEntry holds a per-IP limiter and the last access time (for expiry).
type rateLimiterEntry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// CustomRateLimiterStore is a thread-safe in-memory rate limiter store that
// implements echo middleware.RateLimiterStore and exposes per-IP limit info.
type CustomRateLimiterStore struct {
	rateLimit rate.Limit
	burst     int
	expiresIn time.Duration

	mu       sync.Mutex
	limiters map[string]*rateLimiterEntry
}

// NewCustomRateLimiterStore creates a new CustomRateLimiterStore.
func NewCustomRateLimiterStore(rps float64, burst int, expiresIn time.Duration) *CustomRateLimiterStore {
	store := &CustomRateLimiterStore{
		rateLimit: rate.Limit(rps),
		burst:     burst,
		expiresIn: expiresIn,
		limiters:  make(map[string]*rateLimiterEntry),
	}
	// Start background cleanup goroutine
	go store.cleanup()
	return store
}

// Allow implements echo middleware.RateLimiterStore.
// Returns true if the request is allowed, false if rate-limited.
func (s *CustomRateLimiterStore) Allow(identifier string) (bool, error) {
	entry := s.getOrCreate(identifier)
	return entry.limiter.Allow(), nil
}

// GetInfo returns the current rate limit status for the given identifier.
// This does NOT consume a token.
func (s *CustomRateLimiterStore) GetInfo(identifier string) RateLimitInfo {
	entry := s.getOrCreate(identifier)
	limiter := entry.limiter

	// Current available tokens (float64, may be fractional)
	tokens := limiter.Tokens()
	if tokens < 0 {
		tokens = 0
	}
	remaining := int(math.Floor(tokens))

	// Calculate reset time: time until the bucket is full (burst capacity)
	var resetAt time.Time
	if remaining < s.burst {
		// Tokens missing to reach burst
		missing := float64(s.burst) - tokens
		// Time in seconds to refill 'missing' tokens at rate rps
		secondsToFull := missing / float64(s.rateLimit)
		resetAt = time.Now().Add(time.Duration(secondsToFull * float64(time.Second)))
	} else {
		// Already full
		resetAt = time.Now()
	}

	return RateLimitInfo{
		Limit:     s.burst,
		Remaining: remaining,
		Reset:     resetAt.Unix(),
		ResetAt:   resetAt.UTC(),
	}
}

// getOrCreate retrieves an existing limiter or creates a new one for the identifier.
func (s *CustomRateLimiterStore) getOrCreate(identifier string) *rateLimiterEntry {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.limiters[identifier]
	if !ok {
		entry = &rateLimiterEntry{
			limiter: rate.NewLimiter(s.rateLimit, s.burst),
		}
		s.limiters[identifier] = entry
	}
	entry.lastSeen = time.Now()
	return entry
}

// cleanup removes expired entries periodically to prevent memory leaks.
func (s *CustomRateLimiterStore) cleanup() {
	// Run cleanup at half the expiry interval
	ticker := time.NewTicker(s.expiresIn / 2)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for id, entry := range s.limiters {
			if now.Sub(entry.lastSeen) > s.expiresIn {
				delete(s.limiters, id)
			}
		}
		s.mu.Unlock()
	}
}
