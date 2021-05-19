REGISTRY?=kong-docker-kubernetes-ingress-controller.bintray.io
TAG?=1.2.0
REPO_INFO=$(shell git config --get remote.origin.url)
IMGNAME?=kong-ingress-controller
IMAGE = $(REGISTRY)/$(IMGNAME)
# only for dev
DB?=false
RUN_VERSION?=20
KUBE_VERSION?=v1.20.2

ifndef COMMIT
  COMMIT := $(shell git rev-parse --short HEAD)
endif

export GO111MODULE=on

.PHONY: test-all
test-all: lint test

.PHONY: test
test:
	go test -race ./...

.PHONY: coverage
coverage:
	go test -race -v -count=1 -coverprofile=coverage.out.tmp ./...
	# ignoring generated code for coverage
	grep -E -v 'pkg/apis/|pkg/client/|generated.go|generated.deepcopy.go' coverage.out.tmp > coverage.out
	rm -f coverage.out.tmp

.PHONY: lint
lint: verify-tidy
	golangci-lint run ./...

.PHONY: build
build:
	CGO_ENABLED=0 go build -o kong-ingress-controller ./cli/ingress-controller

.PHONY: verify-manifests
verify-manifests:
	./hack/verify-manifests.sh

.PHONY: verify-codegen
verify-codegen:
	./hack/verify-codegen.sh

.PHONY: update-codegen
update-codegen:
	./hack/update-codegen.sh

.PHONY: verify-tidy
verify-tidy:
	./hack/verify-tidy.sh

.PHONY: container
container:
	docker build \
    --build-arg TAG=${TAG} --build-arg COMMIT=${COMMIT} \
    --build-arg REPO_INFO=${REPO_INFO} \
    -t ${IMAGE}:${TAG} .

.PHONY: railgun-container
railgun-container:
	docker build \
		-f Dockerfile.railgun \
    --build-arg TAG=${TAG} --build-arg COMMIT=${COMMIT} \
    --build-arg REPO_INFO=${REPO_INFO} \
    -t ${IMAGE}:${TAG} .

.PHONY: run
run:
	./hack/dev/start.sh ${DB} ${RUN_VERSION}

# ------------------------------------------------------------------------------
# Integration Tests
# ------------------------------------------------------------------------------

KIND_CLUSTER_NAME ?= "integration-tests"

# Create only the base cluster that would be used for integration tests.
# This can be helpful when developing new tests, as you can deploy the cluster
# and run the test suite setup but then run tests individually against the cluster:
#
#   $ make test.integration.cluster KIND_CLUSTER_NAME="integration-tests"
#   $ export KIND_CLUSTER="integration-tests"
#   $ go test -v -run 'TestTCPIngress' ./test/integration/
.PHONY: test.integration.cluster
test.integration.cluster:
	@./hack/setup-integration-tests.sh
	@go clean -testcache
	@KIND_CLUSTER_NAME="$(KIND_CLUSTER_NAME)" KIND_KEEP_CLUSTER="true" GOFLAGS="-tags=integration_tests" go test -race -v -run "SuiteOnly" ./test/integration/

# Our integration tests using all supported backends, with verbose output
.PHONY: test.integration
test.integration: test.integration.dbless test.integration.postgres

# Our integration tests using the dbless backend, with verbose output
.PHONY: test.integration.dbless
test.integration.dbless:
	@./hack/setup-integration-tests.sh
	@TEST_DATABASE_MODE="off" GOFLAGS="-tags=integration_tests" go test -race -v ./test/integration/

# Our integration tests using the postgres backend, with verbose output
# TODO: race checking has been temporarily turned off because of race conditions found with deck. This will be resolved in an upcoming Alpha release of KIC 2.0.
#       See: https://github.com/Kong/kubernetes-ingress-controller/issues/1324
.PHONY: test.integration.postgres
test.integration.postgres:
	@./hack/setup-integration-tests.sh
	@TEST_DATABASE_MODE="postgres" GOFLAGS="-tags=integration_tests" go test -v ./test/integration/

# Our integration tests using the legacy v1 controller manager
.PHONY: test.integration.legacy
test.integration.legacy:
	@./hack/setup-integration-tests.sh
	@KONG_LEGACY_CONTROLLER=1 GOFLAGS="-tags=integration_tests" go test -race -v ./test/integration/
