#!/bin/bash
az keyvault create \
    --name wayrudev-gateway-kv \
    --resource-group Wayru-dev \
    --location centralus \
    --enable-purge-protection true \
    --retention-days 90 \
    --sku standard
