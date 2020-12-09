#!/bin/bash
set -ex
IFS=: read sut_hostname sut_port <<< "$SUT_HTTPS_HOST"

# Expect HTTPS host+path routes to match.
[ "$(curl -k -sw '%{http_code}' -o /dev/null https://example.com:$sut_port/foo --resolve example.com:$sut_port:$sut_hostname)" == 200 ]
[ "$(curl -k -sw '%{http_code}' -o /dev/null https://example.net:$sut_port/bar --resolve example.net:$sut_port:$sut_hostname)" == 200 ]

# Expect HTTPS paths with wrong SNI NOT to match.
[ "$(curl -k -sw '%{http_code}' -o /dev/null https://example.net:$sut_port/foo --resolve example.net:$sut_port:$sut_hostname)" == 404 ]
[ "$(curl -k -sw '%{http_code}' -o /dev/null https://example.com:$sut_port/bar --resolve example.com:$sut_port:$sut_hostname)" == 404 ]

# Expect HTTPS paths with wrong SNI (but the right Host header) NOT to match.
[ "$(curl --http1.1 -k -sw '%{http_code}' -o /dev/null -H "Host: example.com" https://example.net:$sut_port/foo --resolve example.net:$sut_port:$sut_hostname)" == 404 ]
[ "$(curl --http1.1 -k -sw '%{http_code}' -o /dev/null -H "Host: example.net" https://example.com:$sut_port/bar --resolve example.com:$sut_port:$sut_hostname)" == 404 ]

# Expect plaintext HTTP paths NOT to match.
[ "$(curl -sw '%{http_code}' -H 'Host: example.com' -o /dev/null http://$SUT_HTTP_HOST/foo)" == 426 ]
[ "$(curl -sw '%{http_code}' -H 'Host: example.net' -o /dev/null http://$SUT_HTTP_HOST/bar)" == 426 ]

# For different SNIs expect different server certificates.
[ "$(openssl s_client -servername example.com -connect "$SUT_HTTPS_HOST" </dev/null | sed -n '/^subject/s/.*CN = \([^,]*\).*/\1/p')" == "example.com" ]
[ "$(openssl s_client -servername example.net -connect "$SUT_HTTPS_HOST" </dev/null | sed -n '/^subject/s/.*CN = \([^,]*\).*/\1/p')" == "example.net" ]
