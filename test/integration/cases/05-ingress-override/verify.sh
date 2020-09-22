set -ex

[ "$(curl -XPOST -sw '%{http_code}' -o /dev/null http://$PROXY_IP/foo)" == 404 ]
[ "$(curl -XGET -sw '%{http_code}' -o /dev/null http://$PROXY_IP/foo)" == 200 ]

[ "$(curl -XPOST -sw '%{http_code}' -o /dev/null http://$PROXY_IP/bar)" == 200 ]
[ "$(curl -XGET -sw '%{http_code}' -o /dev/null http://$PROXY_IP/bar)" == 200 ]
