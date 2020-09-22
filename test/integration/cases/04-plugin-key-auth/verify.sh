set -ex

[ "$(curl -sw '%{http_code}' -o /dev/null http://$PROXY_IP/baz)" == 401 ]
[ "$(curl -H "apikey: my-sooper-secret-key" -sw '%{http_code}' -o /dev/null http://$PROXY_IP/baz)" == 200 ]
