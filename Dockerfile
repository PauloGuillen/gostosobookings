# Build stage
FROM golang:1.22 AS build

WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code
COPY . .

# Build the main binary
RUN go build -o gostosobookings-api ./cmd/gostosobookings-api

# Final stage
FROM golang:1.22

WORKDIR /app

# Install PostgreSQL client and curl to use pg_isready
RUN apt-get update && \
    apt-get install -y postgresql-client curl && \
    rm -rf /var/lib/apt/lists/*

# Copy the binary and the initialization script
COPY --from=build /app/gostosobookings-api /app/gostosobookings-api
COPY wait-for-db.sh /app/wait-for-db.sh

# Grant execute permissions to the initialization script
RUN chmod +x /app/wait-for-db.sh

# Default environment variables (optional)
ENV DB_HOST=postgres \
    DB_PORT=5432 \
    DB_USER=gostoso_user \
    DB_PASSWORD=gostoso_password \
    DB_NAME=gostosobookings_db \
    SERVER_PORT=8080

# Set the initialization script as the entrypoint
ENTRYPOINT ["/app/wait-for-db.sh", "&&", "/app/gostosobookings-api"]

# Expose the application port
EXPOSE 8080
