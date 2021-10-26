#!/bin/bash
set -ex

[ "$(curl -sw '%{http_code}' -o /dev/null http://$SUT_HTTP_HOST/foo)" == 200 ]
