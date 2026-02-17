#!/bin/bash

# Get the directory where the script is located
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
cd "$SCRIPT_DIR"

# Create generated directory if it doesn't exist
mkdir -p generated

# Path to schema is relative to this script: ../../schema
SCHEMA_PATH="../../schema"

echo "Generating schema files in $SCRIPT_DIR/generated using schema from $SCHEMA_PATH..."

# Generate JSON
echo "Generating schema.json..."
go run . -schema-path "$SCHEMA_PATH" > generated/schema.json

# Generate Markdown
echo "Generating schema.md..."
go run . -schema-path "$SCHEMA_PATH" -o markdown > generated/schema.md

# Generate Kuzu DDL
echo "Generating schema.ddl..."
go run . -schema-path "$SCHEMA_PATH" -o kuzu > generated/schema.ddl

echo "Generation complete."
ls -l generated
