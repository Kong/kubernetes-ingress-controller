set -ex

[ "$(curl -sw '%{http_code}' -o /dev/null http://$PROXY_IP/foo)" == 200 ]
[ "$(curl -k -sw '%{http_code}' -o /dev/null https://$PROXY_IP/foo)" == 200 ]

[ "$(curl -sw '%{http_code}' -o /dev/null http://$PROXY_IP/bar)" == 426 ]
[ "$(curl -k -sw '%{http_code}' -o /dev/null https://$PROXY_IP/bar)" == 200 ]

