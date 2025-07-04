#!/bin/bash

# Usage: ./run.bash
#
# Environment:
# - GATEWAY_URL
# - USER_ID
#
# This script tests getting a Passpoint profile through the gateway.
# It:
# - Authenticates with Keycloak
# - Sends an HTTP request to the gateway to get a temporary token
# - Sends an HTTP request to the gateway to get a Passpoint profile

# Load env
set -a
[ -f .env ] && . .env
set +a

# Login
cd ../keycloak_login
login_response=$(./run.bash)
cd -> /dev/null

# Get access token from login_response
access_token=$(echo "$login_response" | jq -r '.access_token')

if [ -z "$access_token" ] || [ "$access_token" == "null" ]; then
    echo "Failed to obtain access token from Keycloak" >&2
    echo "$login_response" >&2
    exit 1
fi

echo "Successfully obtained access token" >&2

# Get temporary token
temp_token_response=$(curl -s "${GATEWAY_URL}/idp/profiles/token" \
    -H "Authorization: Bearer ${access_token}")

echo "Temp token ${temp_token_response}" >&2

# Get profile
profile_response=$(curl -s "${GATEWAY_URL}/idp/profiles/${USER_ID}?token=${temp_token_response}&platform=android")

# Output the profile response for piping
echo "$profile_response"

# Print a success message to stderr
echo "Profile retrieval complete" >&2
