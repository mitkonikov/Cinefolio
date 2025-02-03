#!/bin/bash
set -e

DB_PATH="/data/db"

# Check if MongoDB was not shut down properly
if [ -f "$DB_PATH/mongod.lock" ]; then
    echo "Detected mongod.lock, running repair..."
    rm -f "$DB_PATH/mongod.lock"
    mongod --repair --dbpath "$DB_PATH"
fi

# Start MongoDB
exec mongod --config /etc/mongod.conf --bind_ip_all
