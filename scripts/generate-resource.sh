#!/bin/bash

# Latens Resource Generator
# Usage: ./scripts/generate-resource.sh <ResourceName>
# Example: ./scripts/generate-resource.sh Task

set -e

if [ $# -eq 0 ]; then
    echo "🔥 Latens Resource Generator"
    echo ""
    echo "Usage: $0 <ResourceName>"
    echo "Example: $0 Task"
    echo ""
    echo "This will generate:"
    echo "  ✅ Repository interface and implementation"
    echo "  ✅ Service layer with CRUD methods"
    echo "  ✅ Controller with HTTP handlers"
    echo "  ✅ Routes snippet for setup"
    exit 1
fi

RESOURCE_NAME=$1
SCRIPT_DIR=$(dirname "$0")
ROOT_DIR="$SCRIPT_DIR/.."

echo "🚀 Generating resource: $RESOURCE_NAME"
echo ""

# Run the Go generator
cd "$ROOT_DIR"
go run scripts/generate-resource.go "$RESOURCE_NAME"

echo ""
echo "🎯 Next steps:"
echo "1. Add your $RESOURCE_NAME model to internal/database/models/models.go"
echo "2. Copy the routes snippet to cmd/api/routes.go"
echo "3. Run 'go mod tidy' if needed"
echo ""
echo "Happy coding! 🎉"