#!/bin/bash

# Generate Production Swagger Documentation
# Usage: ./scripts/generate-production-swagger.sh

set -e

echo "🚀 Generating Production Swagger Documentation..."

# Check if swag is installed
if ! command -v swag &> /dev/null; then
    echo "❌ Error: swag command not found"
    echo "Please install swag first: go install github.com/swaggo/swag/cmd/swag@latest"
    exit 1
fi

# Set production environment variables
export ENVIRONMENT=production
export PROD_HOST=${1:-"independent-carlene-tushar27x-a3461680.koyeb.app"}

echo "📝 Using PROD_HOST: $PROD_HOST"

# Generate Swagger docs
echo "🔧 Running swag init..."
swag init -g cmd/main.go --parseDependency --parseInternal

echo "✅ Production Swagger documentation generated successfully!"
echo "📚 Swagger UI will now use: $PROD_HOST"
echo ""
echo "🔄 Next steps:"
echo "1. Commit and push the updated docs/"
echo "2. Redeploy to Koyeb"
echo "3. Your Swagger UI will now send requests to $PROD_HOST"
