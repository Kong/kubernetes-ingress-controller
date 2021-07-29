#!/bin/bash
set -ex

[ "$(curl -sw '%{http_code}' -o /dev/null http://$SUT_HTTP_HOST/)" == 404 ]

# Match type: Prefix
[ "$(curl -sw '%{http_code}' -o /dev/null http://$SUT_HTTP_HOST/foo)" == 200 ]
[ "$(curl -sw '%{http_code}' -o /dev/null http://$SUT_HTTP_HOST/foo/)" == 200 ]
[ "$(curl -sw '%{http_code}' -o /dev/null http://$SUT_HTTP_HOST/fooo)" == 404 ]
[ "$(curl -sw '%{http_code}' -o /dev/null http://$SUT_HTTP_HOST/fooo/)" == 404 ]
[ "$(curl -sw '%{http_code}' -o /dev/null http://$SUT_HTTP_HOST/fooo/xxx)" == 404 ]
[ "$(curl -sw '%{http_code}' -o /dev/null http://$SUT_HTTP_HOST/foo/xxx)" == 200 ]

# Match type: Exact
[ "$(curl -sw '%{http_code}' -o /dev/null http://$SUT_HTTP_HOST/bar1)" == 200 ]
[ "$(curl -sw '%{http_code}' -o /dev/null http://$SUT_HTTP_HOST/bar1/)" == 404 ]
[ "$(curl -sw '%{http_code}' -o /dev/null http://$SUT_HTTP_HOST/bar1/xxx)" == 404 ]
[ "$(curl -sw '%{http_code}' -o /dev/null http://$SUT_HTTP_HOST/bar2)" == 404 ]
[ "$(curl -sw '%{http_code}' -o /dev/null http://$SUT_HTTP_HOST/bar2/)" == 200 ]
[ "$(curl -sw '%{http_code}' -o /dev/null http://$SUT_HTTP_HOST/bar2/xxx)" == 404 ]

# Match type: ImplementationSpecific
[ "$(curl -sw '%{http_code}' -o /dev/null http://$SUT_HTTP_HOST/baz)" == 404 ]
[ "$(curl -sw '%{http_code}' -o /dev/null http://$SUT_HTTP_HOST/baz/)" == 200 ]
[ "$(curl -sw '%{http_code}' -o /dev/null http://$SUT_HTTP_HOST/bazzz/)" == 200 ]
[ "$(curl -sw '%{http_code}' -o /dev/null http://$SUT_HTTP_HOST/bazzz/xxx)" == 200 ]
