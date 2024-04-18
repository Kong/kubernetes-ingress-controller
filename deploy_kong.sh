#!/bin/bash

docker run -d --name kong-gateway \
 -e "KONG_ADMIN_LISTEN=0.0.0.0:8001" \
 -e "KONG_DATABASE=off" \
 -p 8000:8000 \
 -p 8001:8001 \
 kong/kong-gateway:3.6.1.2
