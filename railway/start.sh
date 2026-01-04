#!/bin/sh

# Start backend in background
cd /app/backend
./main &

# Start frontend in foreground
cd /app/frontend
npm run start
