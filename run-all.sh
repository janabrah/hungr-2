#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

run_backend() {
  run_with_prefix "backend" bash -c "cd \"${ROOT_DIR}/backend\" && echo '==> make fmt' && make fmt"
  run_with_prefix "backend" bash -c "cd \"${ROOT_DIR}/backend\" && echo '==> make vet' && make vet"
  run_with_prefix "backend" bash -c "cd \"${ROOT_DIR}/backend\" && echo '==> make test' && make test"
}

run_frontend() {
  run_with_prefix "frontend" bash -c "cd \"${ROOT_DIR}/frontend\" && echo '==> npm run format' && npm run format"
  run_with_prefix "frontend" bash -c "cd \"${ROOT_DIR}/frontend\" && echo '==> npm run lint' && npm run lint"
  run_with_prefix "frontend" bash -c "cd \"${ROOT_DIR}/frontend\" && echo '==> npm run build' && npm run build"
  run_with_prefix "frontend" bash -c "cd \"${ROOT_DIR}/frontend\" && echo '==> npm run test' && npm run test"
#  run_with_prefix "frontend" bash -c "cd \"${ROOT_DIR}/frontend\" && echo '==> npx playwright test' && npx playwright test"
}

run_with_prefix() {
  local prefix=$1
  shift
  "$@" 2>&1 | sed "s/^/[$prefix] /"
}

run_backend &
backend_pid=$!

run_frontend &
frontend_pid=$!

wait "$backend_pid" "$frontend_pid"
