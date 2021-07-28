#!/bin/bash
set -ex

[ "$(curl -XPOST -sw '%{http_code}' -o /dev/null http://$SUT_HTTP_HOST/foo)" == 404 ]
[ "$(curl -XGET -sw '%{http_code}' -o /dev/null http://$SUT_HTTP_HOST/foo)" == 200 ]

[ "$(curl -XPOST -sw '%{http_code}' -o /dev/null http://$SUT_HTTP_HOST/bar)" == 200 ]
[ "$(curl -XGET -sw '%{http_code}' -o /dev/null http://$SUT_HTTP_HOST/bar)" == 200 ]
