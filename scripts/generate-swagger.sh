#!/bin/bash

# Script to generate Swagger documentation for different environments
# Usage: ./generate-swagger.sh [environment]
# Environments: development, staging, production

ENVIRONMENT=${1:-development}

echo "Generating Swagger documentation for environment: $ENVIRONMENT"

# Set environment variables
export ENVIRONMENT=$ENVIRONMENT

case $ENVIRONMENT in
    "production")
        export PROD_HOST=${PROD_HOST:-"api.yourdomain.com"}
        export STAGING_HOST=${STAGING_HOST:-"staging-api.yourdomain.com"}
        export DEV_HOST=${DEV_HOST:-"localhost:8082"}
        echo "Production host: $PROD_HOST"
        ;;
    "staging")
        export PROD_HOST=${PROD_HOST:-"api.yourdomain.com"}
        export STAGING_HOST=${STAGING_HOST:-"staging.yourdomain.com"}
        export DEV_HOST=${DEV_HOST:-"localhost:8082"}
        if [ "$STAGING_HOST" = "staging.yourdomain.com" ]; then
            echo "⚠️  Warning: Using default staging host. Set STAGING_HOST env var for custom domain."
        fi
        echo "Staging host: $STAGING_HOST"
        ;;
    "development"|*)
        export PROD_HOST=${PROD_HOST:-"api.yourdomain.com"}
        export STAGING_HOST=${STAGING_HOST:-"staging-api.yourdomain.com"}
        export DEV_HOST=${DEV_HOST:-"localhost:8082"}
        echo "Development host: $DEV_HOST"
        ;;
esac

# Generate Swagger docs
echo "Running swag init..."
swag init -g cmd/main.go

echo "✅ Swagger documentation generated for $ENVIRONMENT environment"
echo "📁 Documentation files created in docs/ directory"
echo "🌐 Access Swagger UI at: http://$DEV_HOST/swagger/index.html (dev) or https://$PROD_HOST/swagger/index.html (prod)"
