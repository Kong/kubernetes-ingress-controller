set -ex

[ "$(curl -sw '%{http_code}' -o /dev/null http://$SUT_HTTP_HOST/baz)" == 401 ]
[ "$(curl -H "apikey: my-sooper-secret-key" -sw '%{http_code}' -o /dev/null http://$SUT_HTTP_HOST/baz)" == 200 ]
