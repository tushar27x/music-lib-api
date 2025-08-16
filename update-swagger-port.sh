#!/bin/bash

# Script to update Swagger port configuration
# Usage: ./update-swagger-port.sh <port_number>

if [ $# -eq 0 ]; then
    echo "Usage: $0 <port_number>"
    echo "Example: $0 8082"
    exit 1
fi

PORT=$1

echo "Updating Swagger configuration to use port $PORT..."

# Update main.go
sed -i "s/localhost:[0-9]*/localhost:$PORT/g" cmd/main.go

# Update README.md
sed -i "s/localhost:[0-9]*/localhost:$PORT/g" README.md

# Regenerate Swagger docs
echo "Regenerating Swagger documentation..."
swag init -g cmd/main.go

echo "✅ Swagger configuration updated to use port $PORT"
echo "🔄 Restart your server and access Swagger UI at: http://localhost:$PORT/swagger/index.html"
