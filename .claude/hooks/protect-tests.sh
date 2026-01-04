#!/bin/bash
# Block edits/writes to test files without confirmation

input=$(cat)
file_path=$(echo "$input" | grep -o '"file_path":"[^"]*"' | cut -d'"' -f4)

# Check if it's a test file
if [[ "$file_path" == *"_test.go"* ]] || [[ "$file_path" == *".test.ts"* ]] || [[ "$file_path" == *".test.js"* ]] || [[ "$file_path" == *".spec.ts"* ]] || [[ "$file_path" == *".spec.js"* ]]; then
  echo "BLOCKED: Cannot modify test file '$file_path' without user approval. Ask the user first." >&2
  exit 2
fi

exit 0
