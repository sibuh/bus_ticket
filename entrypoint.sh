#!/bin/sh

# Start CockroachDB in the background
cockroach start-single-node --insecure --http-port=8080 --store=cockroach-data &

# Wait for CockroachDB to be fully started
sleep 3

# Create the database
cockroach sql --insecure --host=localhost:26257 -e "CREATE DATABASE IF NOT EXISTS tickets;"

# Keep the container running
wait
