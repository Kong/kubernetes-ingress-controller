set -ex

[ "$(curl -sw '%{http_code}' -o /dev/null http://$SUT_HTTP_HOST/)" == 404 ]

[ "$(curl -sw '%{http_code}' -o /dev/null http://$SUT_HTTP_HOST/bar/)" == 200 ]
[ "$(curl -sw '%{http_code}' -o /dev/null http://$SUT_HTTP_HOST/baz/)" == 200 ]

[ "$(curl -sw '%{http_code}' -o /dev/null http://$SUT_HTTP_HOST/qux/)" == 404 ]
