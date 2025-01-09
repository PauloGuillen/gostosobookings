#!/bin/bash
# wait-for-db.sh
echo "Waiting for database to be ready..."
until pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER"; do
  echo "Waiting for PostgreSQL to become available..."
  sleep 2
done
echo "Database is ready!"
exec "$@"  # This will start the Go application
