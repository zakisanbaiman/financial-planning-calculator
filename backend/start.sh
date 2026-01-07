#!/bin/sh

# Function to wait for database to be ready
wait_for_db() {
    echo "Waiting for database to be ready..."
    MAX_RETRIES=${DB_WAIT_MAX_RETRIES:-30}
    RETRY_INTERVAL=${DB_WAIT_INTERVAL:-2}
    RETRY_COUNT=0
    
    while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
        # Try to connect using the migrate binary which will fail if DB is not ready
        # Temporarily disable exit on error for this check
        set +e
        ./migrate -command=status > ./migrate_status.log 2>&1
        MIGRATE_EXIT_CODE=$?
        set -e
        
        if [ $MIGRATE_EXIT_CODE -eq 0 ]; then
            echo "Database is ready!"
            rm -f ./migrate_status.log
            return 0
        fi
        
        RETRY_COUNT=$((RETRY_COUNT + 1))
        echo "Database not ready yet (attempt $RETRY_COUNT/$MAX_RETRIES). Waiting ${RETRY_INTERVAL} seconds..."
        sleep $RETRY_INTERVAL
    done
    
    echo "ERROR: Database did not become ready after $MAX_RETRIES attempts"
    echo "Last migration status check output:"
    cat ./migrate_status.log
    exit 1
}

# Wait for database to be ready
wait_for_db

# Now enable strict error checking for the rest of the script
set -e

echo "Running database migrations..."
./migrate

echo "Starting application..."
exec ./main
