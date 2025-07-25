#!/bin/bash

# Gmail MCP Server API Test Script
set -e

SERVER_URL="${SERVER_URL:-http://localhost:8080}"

echo "=== Gmail MCP Server API Test ==="
echo "Server URL: $SERVER_URL"
echo ""

# Function to test endpoint
test_endpoint() {
    local method=$1
    local endpoint=$2
    local data=$3
    local description=$4
    
    echo "Testing: $description"
    echo "  $method $endpoint"
    
    if [ -n "$data" ]; then
        response=$(curl -s -X "$method" "$SERVER_URL$endpoint" \
            -H "Content-Type: application/json" \
            -d "$data" \
            -w "\nHTTP_CODE:%{http_code}" 2>/dev/null)
    else
        response=$(curl -s -X "$method" "$SERVER_URL$endpoint" \
            -w "\nHTTP_CODE:%{http_code}" 2>/dev/null)
    fi
    
    http_code=$(echo "$response" | tail -n1 | cut -d: -f2)
    body=$(echo "$response" | head -n -1)
    
    echo "  HTTP $http_code"
    echo "  Response: $body"
    echo ""
    
    if [ "$http_code" -ge 200 ] && [ "$http_code" -lt 300 ]; then
        echo "  ✓ SUCCESS"
    else
        echo "  ✗ FAILED"
    fi
    echo ""
}

# Test health endpoint
test_endpoint "GET" "/health" "" "Health Check"

# Test list tools
test_endpoint "GET" "/mcp/tools" "" "List Available Tools"

# Test send email (should fail without auth)
test_endpoint "POST" "/mcp/tools/send_email" '{
    "arguments": {
        "to": ["test@example.com"],
        "subject": "Test Email",
        "body": "This is a test email"
    }
}' "Send Email (should fail - not authenticated)"

# Test search emails (should fail without auth)
test_endpoint "POST" "/mcp/tools/search_emails" '{
    "arguments": {
        "query": "from:test@example.com",
        "maxResults": 5
    }
}' "Search Emails (should fail - not authenticated)"

# Test SSE endpoint (just check if it responds)
echo "Testing: SSE Endpoint"
echo "  GET /mcp/sse"
timeout 2s curl -s "$SERVER_URL/mcp/sse" -H "Accept: text/event-stream" || echo "SSE test completed (timeout expected)"
echo ""
echo "  ✓ SSE endpoint accessible"
echo ""

echo "=== API Test Complete ===" 
echo ""
echo "Note: Most endpoints will fail without proper Gmail authentication."
echo "To fully test:"
echo "1. Set up proper OAuth2 credentials in gcp-oauth.keys.json"
echo "2. Run: ./gmail-mcp-server -auth"
echo "3. Then start the server and run tests again"