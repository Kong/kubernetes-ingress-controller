#!/bin/bash

cleanup() {
	kill $(jobs -p)
}
trap cleanup EXIT

CASES_DIR="$(dirname "$BASH_SOURCE")/../cases"
TEST_RUNNER="$(dirname "$BASH_SOURCE")/run-one-test.sh"

echo ">>> Obtaining Kong proxy IP..."
HTTP_PORT=27080
HTTPS_PORT=27443
kubectl port-forward -n kong svc/kong-proxy "$HTTP_PORT:80" "$HTTPS_PORT:443" &
export SUT_HTTP_HOST="127.0.0.1:$HTTP_PORT"
export SUT_HTTPS_HOST="127.0.0.1:$HTTPS_PORT"
echo ">>> Kong proxy host is '$SUT_HTTP_HOST' for HTTP and '$SUT_HTTPS_HOST' for HTTPS."

echo ">>> Setting up example services..."
setup_example_services() (
	set -ex

	kubectl apply -f https://bit.ly/sample-echo-service
	kubectl apply -f https://bit.ly/sample-httpbin-service

	kubectl wait --for=condition=Available deploy echo --timeout=300s
	kubectl wait --for=condition=Available deploy httpbin --timeout=300s
)

setup_example_services || { echo ">>> ERROR: Failed to set up example services."; exit 1; }

let TESTS_PASSED=0 TESTS_FAILED=0
for CASE_PATH in "$CASES_DIR"/*
do
	CASE_NAME="$(basename "$CASE_PATH")"

	if env \
		CASE_NAME="$CASE_NAME" \
		CASE_PATH="$CASE_PATH" \
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
