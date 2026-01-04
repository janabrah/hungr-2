#!/bin/bash
# Block edits/writes to test files unless approved

input=$(cat)
file_path=$(echo "$input" | grep -o '"file_path":"[^"]*"' | cut -d'"' -f4)

# Check if it's a test file
if [[ "$file_path" == *"_test.go"* ]] || [[ "$file_path" == *".test.ts"* ]] || [[ "$file_path" == *".test.js"* ]] || [[ "$file_path" == *".spec.ts"* ]] || [[ "$file_path" == *".spec.js"* ]]; then
  # Allow if approval file exists
  if [[ -f "$CLAUDE_PROJECT_DIR/.claude/approve-test-edits" ]]; then
    exit 0
  fi
  echo "BLOCKED: Cannot modify test file '$file_path'. Run 'touch .claude/approve-test-edits' to allow, or delete it when done." >&2
  exit 2
fi

exit 0
