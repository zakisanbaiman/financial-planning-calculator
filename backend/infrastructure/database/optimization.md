# Database Performance Optimization Guide

## Overview

This document outlines database performance optimization strategies for the Financial Planning Calculator application.

## Indexing Strategy

### Primary Indexes

```sql
-- User-based queries (most common access pattern)
CREATE INDEX idx_financial_plans_user_id ON financial_plans(user_id);
CREATE INDEX idx_goals_user_id ON goals(user_id);

-- Status-based queries
CREATE INDEX idx_goals_user_status ON goals(user_id, is_active);

-- Date-based queries
CREATE INDEX idx_goals_target_date ON goals(target_date) WHERE is_active = true;
CREATE INDEX idx_financial_plans_updated ON financial_plans(updated_at);
```

### Composite Indexes

```sql
-- For common query patterns
CREATE INDEX idx_goals_user_type_active ON goals(user_id, goal_type, is_active);
CREATE INDEX idx_financial_plans_user_updated ON financial_plans(user_id, updated_at DESC);
```

## Query Optimization

### 1. Use Prepared Statements

Always use prepared statements to:
- Prevent SQL injection
- Enable query plan caching
- Improve performance for repeated queries

```go
// Good
stmt, err := db.Prepare("SELECT * FROM goals WHERE user_id = $1 AND is_active = $2")
defer stmt.Close()
rows, err := stmt.Query(userID, true)

// Avoid
query := fmt.Sprintf("SELECT * FROM goals WHERE user_id = '%s'", userID)
```

### 2. Limit Result Sets

Always use LIMIT when fetching multiple records:

```go
// Good
query := `SELECT * FROM goals WHERE user_id = $1 ORDER BY created_at DESC LIMIT 100`

// Avoid unbounded queries
query := `SELECT * FROM goals WHERE user_id = $1`
```

### 3. Select Only Required Columns

```go
// Good
query := `SELECT id, title, target_amount FROM goals WHERE user_id = $1`

// Avoid SELECT *
query := `SELECT * FROM goals WHERE user_id = $1`
```

### 4. Use Batch Operations

For multiple inserts/updates:

```go
// Good - Batch insert
query := `INSERT INTO goals (user_id, title, target_amount) VALUES ($1, $2, $3), ($4, $5, $6), ($7, $8, $9)`

// Avoid - Multiple single inserts
for _, goal := range goals {
    query := `INSERT INTO goals (user_id, title, target_amount) VALUES ($1, $2, $3)`
    db.Exec(query, goal.UserID, goal.Title, goal.TargetAmount)
}
```

## Connection Pooling

### Optimal Settings

```go
db.SetMaxOpenConns(25)        // Maximum number of open connections
db.SetMaxIdleConns(5)         // Maximum number of idle connections
db.SetConnMaxLifetime(5 * time.Minute)  // Maximum connection lifetime
db.SetConnMaxIdleTime(1 * time.Minute)  // Maximum idle time
```

### Rationale

- **MaxOpenConns (25)**: Prevents overwhelming the database while allowing concurrent requests
- **MaxIdleConns (5)**: Keeps connections ready for reuse without wasting resources
- **ConnMaxLifetime (5 min)**: Prevents stale connections
- **ConnMaxIdleTime (1 min)**: Closes unused connections quickly

## Caching Strategy

### 1. Application-Level Caching

Cache frequently accessed, rarely changing data:

```go
// Cache user financial profiles for 5 minutes
cache.Set(fmt.Sprintf("profile:%s", userID), profile, 5*time.Minute)

// Cache calculation results for 10 minutes
cache.Set(fmt.Sprintf("calc:%s:%s", userID, calcType), result, 10*time.Minute)
```

### 2. Query Result Caching

For expensive calculations:

```go
// Cache aggregated statistics
cache.Set("stats:goals:summary", summary, 15*time.Minute)
```

### 3. Cache Invalidation

Invalidate cache on data changes:

```go
func (r *Repository) UpdateGoal(ctx context.Context, goal Goal) error {
    // Update database
    err := r.db.Update(goal)
    
    // Invalidate cache
    cache.Delete(fmt.Sprintf("goal:%s", goal.ID))
    cache.Delete(fmt.Sprintf("goals:user:%s", goal.UserID))
    
    return err
}
```

## Transaction Management

### 1. Keep Transactions Short

```go
// Good
tx, _ := db.Begin()
defer tx.Rollback()

// Quick operations only
tx.Exec("UPDATE goals SET current_amount = $1 WHERE id = $2", amount, id)

tx.Commit()
```

### 2. Use Read-Only Transactions

For queries that don't modify data:

```go
tx, _ := db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
defer tx.Rollback()

// Read operations
rows, _ := tx.Query("SELECT * FROM goals WHERE user_id = $1", userID)
```

### 3. Avoid Long-Running Transactions

```go
// Avoid
tx, _ := db.Begin()
// ... complex business logic ...
// ... external API calls ...
tx.Commit()

// Good - Keep transactions focused
result := performBusinessLogic()
tx, _ := db.Begin()
tx.Exec("INSERT INTO results VALUES ($1)", result)
tx.Commit()
```

## Monitoring and Analysis

### 1. Enable Query Logging

In development:

```sql
-- PostgreSQL
ALTER DATABASE financial_planning SET log_statement = 'all';
ALTER DATABASE financial_planning SET log_duration = on;
ALTER DATABASE financial_planning SET log_min_duration_statement = 100; -- Log queries > 100ms
```

### 2. Use EXPLAIN ANALYZE

Analyze slow queries:

```sql
EXPLAIN ANALYZE
SELECT g.* FROM goals g
WHERE g.user_id = 'user-123'
AND g.is_active = true
ORDER BY g.target_date;
```

### 3. Monitor Connection Pool

```go
stats := db.Stats()
log.Printf("Open connections: %d", stats.OpenConnections)
log.Printf("In use: %d", stats.InUse)
log.Printf("Idle: %d", stats.Idle)
log.Printf("Wait count: %d", stats.WaitCount)
log.Printf("Wait duration: %s", stats.WaitDuration)
```

## Performance Benchmarks

### Target Metrics

- **Simple queries**: < 10ms
- **Complex calculations**: < 100ms
- **Batch operations**: < 500ms
- **Report generation**: < 2s

### Load Testing

```bash
# Test concurrent requests
ab -n 1000 -c 10 http://localhost:8080/api/calculations/asset-projection

# Monitor database during load
watch -n 1 'psql -c "SELECT * FROM pg_stat_activity WHERE datname = '\''financial_planning'\'';"'
```

## Common Anti-Patterns to Avoid

### 1. N+1 Query Problem

```go
// Bad
for _, goal := range goals {
    user := fetchUser(goal.UserID) // N queries
}

// Good
userIDs := extractUserIDs(goals)
users := fetchUsersBatch(userIDs) // 1 query
```

### 2. Fetching Unnecessary Data

```go
// Bad
goals := fetchAllGoals() // Fetches all columns, all rows
filtered := filterInMemory(goals)

// Good
goals := fetchGoalsWithFilter(filter) // Filter in database
```

### 3. Missing Indexes

```sql
-- Bad - Full table scan
SELECT * FROM goals WHERE target_date > NOW();

-- Good - Use index
CREATE INDEX idx_goals_target_date ON goals(target_date);
SELECT * FROM goals WHERE target_date > NOW();
```

## Optimization Checklist

- [ ] All foreign keys have indexes
- [ ] Frequently queried columns are indexed
- [ ] Connection pool is properly configured
- [ ] Prepared statements are used for repeated queries
- [ ] Transactions are kept short
- [ ] Query results are cached when appropriate
- [ ] Slow query logging is enabled
- [ ] Regular VACUUM and ANALYZE are scheduled
- [ ] Database statistics are up to date

## Maintenance Tasks

### Daily

```sql
-- Update statistics
ANALYZE;
```

### Weekly

```sql
-- Vacuum to reclaim space
VACUUM ANALYZE;
```

### Monthly

```sql
-- Full vacuum (requires downtime)
VACUUM FULL ANALYZE;

-- Reindex
REINDEX DATABASE financial_planning;
```

## Resources

- [PostgreSQL Performance Tuning](https://wiki.postgresql.org/wiki/Performance_Optimization)
- [Go database/sql Best Practices](https://go.dev/doc/database/manage-connections)
- [Indexing Strategies](https://use-the-index-luke.com/)
