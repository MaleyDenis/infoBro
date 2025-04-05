#!/bin/bash

# Exit on error
set -e

# Build backend
echo "Building backend..."
go build -o bin/infobro ./cmd/server

# Start MongoDB and Redis if they are not running
MONGO_RUNNING=$(docker ps -q -f name=mongodb)
REDIS_RUNNING=$(docker ps -q -f name=redis)

if [ -z "$MONGO_RUNNING" ]; then
  echo "Starting MongoDB..."
  docker run -d -p 27017:27017 --name mongodb mongo
else
  echo "MongoDB is already running"
fi

if [ -z "$REDIS_RUNNING" ]; then
  echo "Starting Redis..."
  docker run -d -p 6379:6379 --name redis redis
else
  echo "Redis is already running"
fi

# Start backend in background
echo "Starting backend..."
./bin/infobro --mongo-uri mongodb://localhost:27017 --redis-addr localhost:6379 &
BACKEND_PID=$!

# Change to web directory
cd web

# Install frontend dependencies if needed
if [ ! -d "node_modules" ]; then
  echo "Installing frontend dependencies..."
  npm install
fi

# Start frontend
echo "Starting frontend..."
npm run start &
FRONTEND_PID=$!

# Handle exit
function cleanup {
  echo "Shutting down..."
  kill $BACKEND_PID
  kill $FRONTEND_PID
}

trap cleanup EXIT

# Wait for user to exit
echo "Press Ctrl+C to exit"
wait