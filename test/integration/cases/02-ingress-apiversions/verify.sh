set -ex

[ "$(curl -sw '%{http_code}' -o /dev/null http://$PROXY_IP/)" == 404 ]

[ "$(curl -sw '%{http_code}' -o /dev/null http://$PROXY_IP/bar/)" == 200 ]
[ "$(curl -sw '%{http_code}' -o /dev/null http://$PROXY_IP/baz/)" == 200 ]

[ "$(curl -sw '%{http_code}' -o /dev/null http://$PROXY_IP/qux/)" == 404 ]
