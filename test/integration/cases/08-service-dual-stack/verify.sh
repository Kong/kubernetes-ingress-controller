#!/bin/bash
set -ex

# For ipv4 service
TARGET_IPV4_ENDPOINT="target: 1.1.1.1:8080"
[ "$(curl -s -k https://$SUT_ADMIN_API_HOST/config | grep -o "$TARGET_IPV4_ENDPOINT")" == "$TARGET_IPV4_ENDPOINT" ]

# For ipv6 service
TARGET_IPV6_ENDPOINT="target: '[1::1]:8080'"
[ "$(curl -s -k https://$SUT_ADMIN_API_HOST/config/| grep -o -F "$TARGET_IPV6_ENDPOINT")" == "$TARGET_IPV6_ENDPOINT" ]
