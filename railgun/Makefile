##@ Versioning

ifndef COMMIT
  COMMIT := $(shell git rev-parse --short HEAD)
endif
REPO_INFO=$(shell git config --get remote.origin.url)

##@ Container Image Vars

# Image URL to use all building/pushing image targets
IMG ?= controller:latest
# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
# TODO: investigate `go-kong` changes that would allow us to remove `allowDangerousTypes=true`.
#       See: https://github.com/Kong/go-kong/issues/54

##@ Testing Vars

PKG_LIST = ./...
COVERAGE_PROFILE=coverage.out
COVERAGE_INTEGRATION_PROFILE=coverage-integration.out

KIND_CLUSTER_NAME ?= "integration-tests"

##@ Build Vars

CRD_OPTIONS ?= "crd:trivialVersions=true,preserveUnknownFields=false,allowDangerousTypes=true"
# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

all: build

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

manifests: controller-gen ## Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects.
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases

generate: controller-gen ## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

fmt: ## Run go fmt against code.
	go fmt ./...

vet: ## Run go vet against code.
	go vet ./...

ENVTEST_ASSETS_DIR=$(shell pwd)/testbin
test: manifests generate fmt vet ## Run tests.
	mkdir -p ${ENVTEST_ASSETS_DIR}
	test -f ${ENVTEST_ASSETS_DIR}/setup-envtest.sh || curl -sSLo ${ENVTEST_ASSETS_DIR}/setup-envtest.sh https://raw.githubusercontent.com/kubernetes-sigs/controller-runtime/v0.7.0/hack/setup-envtest.sh
	source ${ENVTEST_ASSETS_DIR}/setup-envtest.sh; fetch_envtest_tools $(ENVTEST_ASSETS_DIR); setup_envtest_env $(ENVTEST_ASSETS_DIR); go test -race -covermode=atomic -coverpkg=$(PKG_LIST) $(PKG_LIST) -coverprofile=$(COVERAGE_PROFILE)

##@ Build

build: generate fmt vet ## Build manager binary.
	go build -a -o bin/manager -ldflags "-s -w \
		-X github.com/kong/kubernetes-ingress-controller/railgun/manager.Release=$(TAG) \
		-X github.com/kong/kubernetes-ingress-controller/railgun/manager.Commit=$(COMMIT) \
		-X github.com/kong/kubernetes-ingress-controller/railgun/manager.Repo=$(REPO_INFO)" main.go

run: manifests generate fmt vet ## Run a controller from your host.
	go run ./main.go

docker-push: ## Push docker image with the manager.
	docker push ${IMG}

##@ Deployment

install: manifests kustomize ## Install CRDs into the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/crd | kubectl apply -f -

uninstall: manifests kustomize ## Uninstall CRDs from the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/crd | kubectl delete -f -

deploy: manifests kustomize ## Deploy controller to the K8s cluster specified in ~/.kube/config.
	cd config/manager && $(KUSTOMIZE) edit set image controller=${IMG}
	$(KUSTOMIZE) build config/default | kubectl apply -f -

undeploy: ## Undeploy controller from the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/default | kubectl delete -f -

CONTROLLER_GEN = $(shell pwd)/bin/controller-gen
controller-gen: ## Download controller-gen locally if necessary.
	$(call go-get-tool,$(CONTROLLER_GEN),sigs.k8s.io/controller-tools/cmd/controller-gen@v0.4.1)

KUSTOMIZE = $(shell pwd)/bin/kustomize
kustomize: ## Download kustomize locally if necessary.
	$(call go-get-tool,$(KUSTOMIZE),sigs.k8s.io/kustomize/kustomize/v3@v3.8.7)

# go-get-tool will 'go get' any package $2 and install it to $1.
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-get-tool
@[ -f $(1) ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
go mod init tmp ;\
echo "Downloading $(2)" ;\
GOBIN=$(PROJECT_DIR)/bin go get $(2) ;\
rm -rf $$TMP_DIR ;\
}
endef

# Kong Custom Targets

# this will generate the custom typed clients needed for end-users implementing logic in Go to use our API types.
# TODO: we're hacking around client-gen for now to enable it for enabled go modules, should probably contribute upstream to improve this.
#       See: https://github.com/Kong/kubernetes-ingress-controller/issues/1254
generate.clientsets:
	@client-gen --go-header-file ./hack/boilerplate.go.txt \
		--clientset-name clientset \
		--input-base github.com/kong/kubernetes-ingress-controller/railgun/apis/  \
		--input configuration/v1,configuration/v1beta1 \
		--input-dirs github.com/kong/kubernetes-ingress-controller/railgun/apis/configuration/v1beta1/,github.com/kong/kubernetes-ingress-controller/railgun/apis/configuration/v1/ \
		--output-base client-gen-tmp/ \
		--output-package github.com/kong/kubernetes-ingress-controller/railgun/pkg/
	@rm -rf pkg/clientset/
	@mv client-gen-tmp/github.com/kong/kubernetes-ingress-controller/railgun/pkg/clientset pkg/
	@rm -rf client-gen-tmp/

# Some of our controllers are simple enough to be generated, this generates them.
generate.controllers:
	go generate ./...

# Cleanup any deployed test proxies
clean.test.proxy:
	helm uninstall --namespace kong-system kong-test-proxy

# Create only the base cluster that would be used for integration tests.
# This can be helpful when developing new tests, as you can deploy the cluster
# and run the test suite setup but then run tests individually against the cluster:
#
#   $ make test.integration.cluster KIND_CLUSTER_NAME="integration-tests"
#   $ export KIND_CLUSTER="integration-tests"
#   $ go test -v -run 'TestTCPIngress' ./test/integration/
.PHONY: test.integration.cluster
test.integration.cluster:
	@./scripts/check-container-environment.sh
	@go clean -testcache
	@KIND_CLUSTER="" KIND_CLUSTER_NAME="$(KIND_CLUSTER_NAME)" KIND_KEEP_CLUSTER="true" GOFLAGS="-tags=integration_tests" go test -race -v -run "SuiteOnly" ./test/integration/

# Our integration tests using all supported backends, with verbose output
.PHONY: test.integration
test.integration: test.integration.dbless test.integration.postgres

# Our integration tests using the dbless backend, with verbose output
.PHONY: test.integration.dbless
test.integration.dbless:
	@./scripts/check-container-environment.sh
	@TEST_DATABASE_MODE="off" GOFLAGS="-tags=integration_tests" go test -timeout 15m -race -v -count=1 -covermode=atomic -coverpkg=$(PKG_LIST) -coverprofile=$(COVERAGE_INTEGRATION_PROFILE) ./test/integration/

# Our integration tests using the postgres backend, with verbose output
# TODO: race checking has been temporarily turned off because of race conditions found with deck. This will be resolved in an upcoming Alpha release of KIC 2.0.
#       See: https://github.com/Kong/kubernetes-ingress-controller/issues/1324
.PHONY: test.integration.postgres
test.integration.postgres:
	@./scripts/check-container-environment.sh
	@TEST_DATABASE_MODE="postgres" GOFLAGS="-tags=integration_tests" go test -timeout 15m -v -count=1 -covermode=atomic -coverpkg=$(PKG_LIST) -coverprofile=$(COVERAGE_INTEGRATION_PROFILE) ./test/integration/

# Our integration tests using the legacy v1 controller manager
.PHONY: test.integration.legacy
test.integration.legacy:
	@./scripts/check-container-environment.sh
	@KONG_LEGACY_CONTROLLER=1 GOFLAGS="-tags=integration_tests" go test -timeout 15m -race -v -count=1 -covermode=atomic -coverpkg=$(PKG_LIST) -coverprofile=$(COVERAGE_INTEGRATION_PROFILE) ./test/integration/
