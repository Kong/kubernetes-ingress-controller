set -ex

[ "$(curl -sw '%{http_code}' -o /dev/null http://$PROXY_IP/)" == 200 ]
[ "$(curl -sw '%{http_code}' -o /dev/null http://$PROXY_IP/status/204)" == 204 ]
[ "$(curl -sw '%{http_code}' -o /dev/null http://$PROXY_IP/foo/)" == 200 ]

