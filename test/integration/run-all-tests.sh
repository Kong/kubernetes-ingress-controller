#!/bin/bash

CASES_DIR="$(dirname "$BASH_SOURCE")/cases"
TEST_RUNNER="$(dirname "$BASH_SOURCE")/util/run-one-test.sh"

echo ">>> Obtaining Kong proxy IP..."
PROXY_IP=$(kubectl get service --namespace kong kong-proxy -o jsonpath={.spec.clusterIP})
echo ">>> Kong proxy IP is '$PROXY_IP'."

echo ">>> Setting up example services..."
setup_example_services() (
	set -ex

	kubectl apply -f https://bit.ly/sample-echo-service
	kubectl apply -f https://bit.ly/sample-httpbin-service

	kubectl wait --for=condition=Available deploy echo --timeout=120s
	kubectl wait --for=condition=Available deploy httpbin --timeout=120s
)

setup_example_services || { echo ">>> ERROR: Failed to set up example services."; exit 1; }

let TESTS_PASSED=0 TESTS_FAILED=0
for CASE_PATH in "$CASES_DIR"/*
do
	CASE_NAME="$(basename "$CASE_PATH")"

	if env \
		CASE_NAME="$CASE_NAME" \
		CASE_PATH="$CASE_PATH" \
		PROXY_IP="$PROXY_IP" \
		$TEST_RUNNER
	then
		let TESTS_PASSED++
	else
		echo ">>> Test $CASE_NAME exited with status $?"
		let TESTS_FAILED++
	fi
done

echo ">>> Overall tests PASSED: $TESTS_PASSED"
echo ">>> Overall tests FAILED: $TESTS_FAILED"

[[ $TESTS_FAILED == 0 ]]
