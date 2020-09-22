set -ex

[ "$(curl -sw '%{http_code}' -o /dev/null http://$PROXY_IP/)" == 404 ]

# Match type: Prefix
[ "$(curl -sw '%{http_code}' -o /dev/null http://$PROXY_IP/foo)" == 200 ]
[ "$(curl -sw '%{http_code}' -o /dev/null http://$PROXY_IP/foo/)" == 200 ]
[ "$(curl -sw '%{http_code}' -o /dev/null http://$PROXY_IP/fooo)" == 404 ]
[ "$(curl -sw '%{http_code}' -o /dev/null http://$PROXY_IP/fooo/)" == 404 ]
[ "$(curl -sw '%{http_code}' -o /dev/null http://$PROXY_IP/fooo/xxx)" == 404 ]
[ "$(curl -sw '%{http_code}' -o /dev/null http://$PROXY_IP/foo/xxx)" == 200 ]

# Match type: Exact
[ "$(curl -sw '%{http_code}' -o /dev/null http://$PROXY_IP/bar1)" == 200 ]
[ "$(curl -sw '%{http_code}' -o /dev/null http://$PROXY_IP/bar1/)" == 404 ]
[ "$(curl -sw '%{http_code}' -o /dev/null http://$PROXY_IP/bar1/xxx)" == 404 ]
[ "$(curl -sw '%{http_code}' -o /dev/null http://$PROXY_IP/bar2)" == 404 ]
[ "$(curl -sw '%{http_code}' -o /dev/null http://$PROXY_IP/bar2/)" == 200 ]
[ "$(curl -sw '%{http_code}' -o /dev/null http://$PROXY_IP/bar2/xxx)" == 404 ]

# Match type: ImplementationSpecific
[ "$(curl -sw '%{http_code}' -o /dev/null http://$PROXY_IP/baz)" == 404 ]
[ "$(curl -sw '%{http_code}' -o /dev/null http://$PROXY_IP/baz/)" == 200 ]
[ "$(curl -sw '%{http_code}' -o /dev/null http://$PROXY_IP/bazzz/)" == 200 ]
[ "$(curl -sw '%{http_code}' -o /dev/null http://$PROXY_IP/bazzz/xxx)" == 200 ]
