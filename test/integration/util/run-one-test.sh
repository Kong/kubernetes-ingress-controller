#!/bin/bash

fail_usage() {
	echo ">>> ERR: Required environment variable $1 not set."
	exit 100
}

cleanup() {
	echo ">>> Test $CASE_NAME: cleanup"
	for kind in ingress kongplugin kongconsumer
	do
	kubectl delete "$kind" --all-namespaces --all
	done
}
trap cleanup EXIT

[ -n "$PROXY_IP" ] || fail_usage PROXY_IP
[ -n "$CASE_NAME" ] || fail_usage CASE_NAME
[ -n "$CASE_PATH" ] || fail_usage CASE_PATH

echo ">>> Test $CASE_NAME: apply manifests"
kubectl apply -f "$CASE_PATH" || exit 1

echo ">>> Test $CASE_NAME: wait"
sleep 6

echo ">>> Test $CASE_NAME: verify"
"$CASE_PATH/verify.sh"
STATUS=$?

if [ $STATUS != 0 ]
then
	echo ">>> Test $CASE_NAME: FAIL (exit code $STATUS)"
	exit $STATUS
fi

echo ">>> Test $CASE_NAME: PASS"
exit 0
