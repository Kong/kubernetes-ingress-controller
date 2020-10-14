set -ex

curl -v http://$SUT_HTTP_HOST/baz/ | grep "my-request-id"
