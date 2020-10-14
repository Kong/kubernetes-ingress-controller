set -ex

[ "$(curl -sw '%{http_code}' -o /dev/null http://$SUT_HTTP_HOST/foo)" == 200 ]
[ "$(curl -k -sw '%{http_code}' -o /dev/null https://$SUT_HTTPS_HOST/foo)" == 200 ]

[ "$(curl -sw '%{http_code}' -o /dev/null http://$SUT_HTTP_HOST/bar)" == 426 ]
[ "$(curl -k -sw '%{http_code}' -o /dev/null https://$SUT_HTTPS_HOST/bar)" == 200 ]

