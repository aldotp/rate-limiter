#!/bin/bash

URL="http://localhost:8080/ping"
API_KEY="test-api-key"
TOTAL_REQUESTS=150
DELAY=0.0              

echo "Starting test to $URL with X-API-KEY: $API_KEY"
echo "Sending $TOTAL_REQUESTS requests with $DELAY sec delay"

for i in $(seq 1 $TOTAL_REQUESTS); do
    echo -n "[$i] "
    curl -s -o /dev/null -w "%{http_code}\n" -H "X-API-KEY: $API_KEY" "$URL"
    sleep $DELAY
done
