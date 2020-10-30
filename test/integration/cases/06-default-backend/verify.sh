#!/bin/bash
set -ex

[ "$(curl -sw '%{http_code}' -o /dev/null http://$SUT_HTTP_HOST/)" == 200 ]
[ "$(curl -sw '%{http_code}' -o /dev/null http://$SUT_HTTP_HOST/status/204)" == 204 ]
[ "$(curl -sw '%{http_code}' -o /dev/null http://$SUT_HTTP_HOST/foo/)" == 200 ]

