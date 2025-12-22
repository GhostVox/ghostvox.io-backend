#!/bin/bash

set -eo pipefail

# --- Configuration ---
DIRECTION=${1:-up}
MIGRATIONS_DIR="/app/migrations"

# --- Database Migrations ---
echo "Running migrations in direction: $DIRECTION"
goose -dir "$MIGRATIONS_DIR" postgres "$DB_URL" "$DIRECTION"

# --- Start Application ---
echo "Starting application..."
exec /bin/server
