#!/bin/bash

set -e

host="$DB_HOST"
port="$DB_PORT"
user="$DB_USER"
db="$DB_NAME"

echo "Waiting for PostgreSQL at $host:$port..."

until pg_isready -h "$host" -p "$port" -U "$user" -d "$db"; do
  echo "PostgreSQL is unavailable - retrying..."
  sleep 2
done

echo "PostgreSQL is ready - starting the application."
exec /app/gostosobookings-api
