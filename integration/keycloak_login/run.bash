#!/bin/bash

# Usage: ./run.sh
#
# Environment
# - KEYCLOAK_URL
# - KEYCLOAK_REALM
# - KEYCLOAK_CLIENT_ID
# - KEYCLOAK_CLIENT_SECRET
# - USERNAME
# - PASSWORD

# Load env
set -a
[ -f .env ] && . .env
set +a

# Login
login_response=$(curl -s -X POST "$KEYCLOAK_URL/realms/$KEYCLOAK_REALM/protocol/openid-connect/token" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=password" \
  -d "client_id=$KEYCLOAK_CLIENT_ID" \
  -d "client_secret=$KEYCLOAK_CLIENT_SECRET" \
  -d "username=$USERNAME" \
  -d "password=$PASSWORD")

echo $login_response
