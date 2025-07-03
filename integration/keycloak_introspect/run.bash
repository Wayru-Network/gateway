#!/bin/bash

# Usage: ./kc-introspect.sh <token>
#
# Environment:
# - KEYCLOAK_URL
# - KEYCLOAK_REALM
# - KEYCLOAK_CLIENT_ID
# - KEYCLOAK_CLIENT_SECRET

# Load env
set -a
[ -f .env ] && . .env
set +a

token=$1

introspect_response=$(curl -s -X POST "$KEYCLOAK_URL/realms/$KEYCLOAK_REALM/protocol/openid-connect/token/introspect" \
    -H "Content-Type: application/x-www-form-urlencoded" \
    -d "client_id=$KEYCLOAK_CLIENT_ID" \
    -d "client_secret=$KEYCLOAK_CLIENT_SECRET" \
    -d "token=$token")

echo $introspect_response
