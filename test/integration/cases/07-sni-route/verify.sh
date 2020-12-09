#!/bin/bash
set -ex
IFS=: read sut_hostname sut_port <<< "$SUT_HTTPS_HOST"

[ "$(curl -k -sw '%{http_code}' -o /dev/null https://example.com:$sut_port/foo --resolve example.com:$sut_port:$sut_hostname)" == 200 ]

[ "$(curl -k -sw '%{http_code}' -o /dev/null https://example.net:$sut_port/bar --resolve example.net:$sut_port:$sut_hostname)" == 200 ]

[ "$(curl -k -sw '%{http_code}' -o /dev/null https://example.net:$sut_port/foo --resolve example.net:$sut_port:$sut_hostname)" == 404 ]
