#!/bin/bash
# Smoke test for the Hungr API
# Usage: ./smoke-test.sh [BASE_URL]

BASE_URL="${1:-http://localhost:8080}"
TEST_EMAIL="smoke-test@example.com"
FAILED=0

echo "Running smoke tests against $BASE_URL"
echo "========================================"

# Helper function
check_response() {
    local name="$1"
    local expected_code="$2"
    local actual_code="$3"
    local response="$4"

    if [ "$actual_code" -eq "$expected_code" ]; then
        echo "✓ $name (HTTP $actual_code)"
    else
        echo "✗ $name - Expected HTTP $expected_code, got $actual_code"
        echo "  Response: $response"
        FAILED=1
    fi
}

# Test 1: Health check
echo ""
echo "1. Health check"
RESPONSE=$(curl -s -w "\n%{http_code}" "$BASE_URL/health")
HTTP_CODE=$(echo "$RESPONSE" | tail -1)
BODY=$(echo "$RESPONSE" | head -n -1)
check_response "GET /health" 200 "$HTTP_CODE" "$BODY"

# Test 2: Create user (for testing)
echo ""
echo "2. Create test user"
RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/api/users" \
    -H "Content-Type: application/json" \
    -d "{\"email\": \"$TEST_EMAIL\", \"name\": \"Smoke Test User\"}")
HTTP_CODE=$(echo "$RESPONSE" | tail -1)
BODY=$(echo "$RESPONSE" | head -n -1)
# 200 = created, 409 = already exists (both OK)
if [ "$HTTP_CODE" -eq 200 ] || [ "$HTTP_CODE" -eq 409 ]; then
    echo "✓ POST /api/users (HTTP $HTTP_CODE)"
else
    echo "✗ POST /api/users - Expected HTTP 200 or 409, got $HTTP_CODE"
    echo "  Response: $BODY"
    FAILED=1
fi

# Test 3: Get recipes with email parameter (the bug we fixed)
echo ""
echo "3. Get recipes with email parameter"
RESPONSE=$(curl -s -w "\n%{http_code}" "$BASE_URL/api/recipes?email=$TEST_EMAIL")
HTTP_CODE=$(echo "$RESPONSE" | tail -1)
BODY=$(echo "$RESPONSE" | head -n -1)
check_response "GET /api/recipes?email=..." 200 "$HTTP_CODE" "$BODY"

# Verify response has expected structure
if echo "$BODY" | grep -q '"recipeData"'; then
    echo "  ✓ Response contains recipeData"
else
    echo "  ✗ Response missing recipeData"
    FAILED=1
fi

# Test 4: Get recipes with old user_uuid parameter should fail
echo ""
echo "4. Verify old user_uuid parameter is rejected"
RESPONSE=$(curl -s -w "\n%{http_code}" "$BASE_URL/api/recipes?user_uuid=some-uuid")
HTTP_CODE=$(echo "$RESPONSE" | tail -1)
BODY=$(echo "$RESPONSE" | head -n -1)
check_response "GET /api/recipes?user_uuid=... (should fail)" 400 "$HTTP_CODE" "$BODY"

# Test 5: Get user by email
echo ""
echo "5. Get user by email"
RESPONSE=$(curl -s -w "\n%{http_code}" "$BASE_URL/api/users?email=$TEST_EMAIL")
HTTP_CODE=$(echo "$RESPONSE" | tail -1)
BODY=$(echo "$RESPONSE" | head -n -1)
check_response "GET /api/users?email=..." 200 "$HTTP_CODE" "$BODY"

echo ""
echo "========================================"
if [ "$FAILED" -eq 0 ]; then
    echo "All smoke tests passed! ✓"
    exit 0
else
    echo "Some smoke tests failed! ✗"
    exit 1
fi
