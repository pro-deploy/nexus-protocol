#!/bin/bash

# Test script for Nexus Protocol version compatibility

echo "🧪 Testing Nexus Protocol version compatibility..."

# Start server in background (assuming it's already built)
# ./server/go/nexus-server &

SERVER_URL="http://localhost:8080"

# Test cases
test_cases=(
    # Compatible versions
    "2.0.0:200:Compatible version"
    "2.0.1:200:Compatible patch version"

    # Incompatible versions (client newer than server)
    "2.1.0:400:Client minor version too new"
    "1.0.0:400:Major version mismatch"
    "3.0.0:400:Client major version too new"
)

echo ""
echo "Testing protocol version validation..."

for test_case in "${test_cases[@]}"; do
    IFS=':' read -r version expected_status description <<< "$test_case"

    echo "Testing version $version - $description"

    # Test with query parameter
    response=$(curl -s -w "%{http_code}" -o /dev/null \
        -H "Content-Type: application/json" \
        "$SERVER_URL/api/v1/health?protocol_version=$version")

    if [ "$response" -eq "$expected_status" ]; then
        echo "✅ Query param: $version -> $response (expected $expected_status)"
    else
        echo "❌ Query param: $version -> $response (expected $expected_status)"
    fi

    # Test with header
    response=$(curl -s -w "%{http_code}" -o /dev/null \
        -H "Content-Type: application/json" \
        -H "X-Protocol-Version: $version" \
        "$SERVER_URL/api/v1/health")

    if [ "$response" -eq "$expected_status" ]; then
        echo "✅ Header: $version -> $response (expected $expected_status)"
    else
        echo "❌ Header: $version -> $response (expected $expected_status)"
    fi

    echo ""
done

echo "🎉 Protocol version testing completed!"
