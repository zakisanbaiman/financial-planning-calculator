#!/bin/sh

# Note: set -e is deferred until after database readiness check
# to allow retry logic to handle connection failures gracefully

# Function to wait for database to be ready
wait_for_db() {
    echo "Waiting for database to be ready..."
    # Default: 30 retries (60 seconds total with 2-second intervals)
    MAX_RETRIES=${DB_WAIT_MAX_RETRIES:-30}
    # Default: 2 seconds between retry attempts
    RETRY_INTERVAL=${DB_WAIT_INTERVAL:-2}
    RETRY_COUNT=0
    # Use unique log file with process ID to avoid conflicts
    LOG_FILE="/tmp/migrate_status_$$.log"
    
    while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
        # Try to connect using the migrate binary
        # Use explicit if-then to handle errors without toggling set -e
        if ./migrate -command status > "$LOG_FILE" 2>&1; then
            echo "Database is ready!"
            rm -f "$LOG_FILE"
            return 0
        fi
        
        RETRY_COUNT=$((RETRY_COUNT + 1))
        echo "Database not ready yet (attempt $RETRY_COUNT/$MAX_RETRIES). Waiting ${RETRY_INTERVAL} seconds..."
        sleep $RETRY_INTERVAL
    done
    
    echo "ERROR: Database did not become ready after $MAX_RETRIES attempts"
    echo "Last migration status check output:"
    if [ -f "$LOG_FILE" ]; then
        cat "$LOG_FILE"
        rm -f "$LOG_FILE"
    else
        echo "No migration log available"
    fi
    exit 1
}

# Wait for database to be ready
wait_for_db

# Enable strict error checking for the rest of the script
set -e

echo "Running database migrations..."
./migrate -command up

echo "Starting application..."
exec ./main
