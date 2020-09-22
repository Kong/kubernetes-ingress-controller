set -ex

curl -v http://$PROXY_IP/baz/ | grep "my-request-id"
