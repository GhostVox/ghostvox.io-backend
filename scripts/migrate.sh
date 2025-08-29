#!/bin/sh
direction=${1:-up}
echo "Running migrations in direction: $direction"
 goose -dir /app/migrations postgres ${DB_URL} $direction

echo "Starting application..."
exec /bin/server
