# ------------------------------------------------------------------------------
# Configuration - Repository
# ------------------------------------------------------------------------------

MAKEFLAGS += --no-print-directory
REPO_URL ?= github.com/kong/kubernetes-ingress-controller
REPO_INFO ?= $(shell git config --get remote.origin.url)
TAG ?= $(shell git describe --tags)

ifndef COMMIT
  COMMIT := $(shell git rev-parse --short HEAD)
endif

# ------------------------------------------------------------------------------
# Configuration - Golang
# ------------------------------------------------------------------------------

export GO111MODULE=on

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# ------------------------------------------------------------------------------
# Configuration - Tooling
# ------------------------------------------------------------------------------

PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))

.PHONY: _download_tool
_download_tool:
	(cd third_party && go mod tidy && \
		GOBIN=$(PROJECT_DIR)/bin go generate -tags=third_party ./$(TOOL).go )

.PHONY: tools
tools: controller-gen kustomize client-gen golangci-lint gotestsum crd-ref-docs

CONTROLLER_GEN = $(PROJECT_DIR)/bin/controller-gen
.PHONY: controller-gen
controller-gen: ## Download controller-gen locally if necessary.
	@$(MAKE) _download_tool TOOL=controller-gen

KUSTOMIZE = $(PROJECT_DIR)/bin/kustomize
.PHONY: kustomize
kustomize: ## Download kustomize locally if necessary.
	@$(MAKE) _download_tool TOOL=kustomize

CLIENT_GEN = $(PROJECT_DIR)/bin/client-gen
.PHONY: client-gen
client-gen: ## Download client-gen locally if necessary.
	@$(MAKE) _download_tool TOOL=client-gen

GOLANGCI_LINT = $(PROJECT_DIR)/bin/golangci-lint
.PHONY: golangci-lint
golangci-lint: ## Download golangci-lint locally if necessary.
	@$(MAKE) _download_tool TOOL=golangci-lint

GOTESTSUM = $(PROJECT_DIR)/bin/gotestsum
.PHONY: gotestsum
gotestsum: ## Download gotestsum locally if necessary.
	@$(MAKE) _download_tool TOOL=gotestsum

CRD_REF_DOCS = $(PROJECT_DIR)/bin/crd-ref-docs
.PHONY: crd-ref-docs
crd-ref-docs: ## Download crd-ref-docs locally if necessary.
	@$(MAKE) _download_tool TOOL=crd-ref-docs

DLV = $(PROJECT_DIR)/bin/dlv
.PHONY: dlv
dlv: ## Download dlv locally if necessary.
	@$(MAKE) _download_tool TOOL=dlv

SKAFFOLD = $(PROJECT_DIR)/bin/skaffold
.PHONY: skaffold
skaffold: ## Download skaffold locally if necessary.
# NOTE: this step is not idempotent like other tool download steps because for
# some reason skaffold doesn't want to be included in imports or installed via
# go install:
# go: github.com/GoogleContainerTools/skaffold@v2.0.4: invalid version: module contains a go.mod file, so module path must match major version ("github.com/GoogleContainerTools/skaffold/v2")
ifeq ($(wildcard $(SKAFFOLD)),)
	curl -Lo skaffold https://storage.googleapis.com/skaffold/releases/v2.1.0/skaffold-$(shell go env GOOS)-$(shell go env GOARCH)
	@chmod +x skaffold
	@mv skaffold ./bin/
endif

STATICCHECK = $(PROJECT_DIR)/bin/staticcheck
.PHONY: staticcheck
staticcheck.download: ## Download staticcheck locally if necessary.
	@$(MAKE) _download_tool TOOL=staticcheck

# ------------------------------------------------------------------------------
# Build
# ------------------------------------------------------------------------------

all: build

.PHONY: clean
clean:
	@rm -rf build/
	@rm -rf testbin/
	@rm -rf bin/*
	@rm -f coverage*.out

.PHONY: build
build: generate fmt vet lint _build

.PHONY: build.fips
build.fips: generate fmt vet lint _build.fips

.PHONY: _build
_build:
	$(MAKE) _build.template MAIN=./internal/cmd/main.go

.PHONY: _build.fips
_build.fips:
	$(MAKE) _build.template MAIN=./internal/cmd/fips/main.go

.PHONY: _build.template
_build.template:
	go build -o bin/manager -ldflags "-s -w \
		-X $(REPO_URL)/v2/internal/manager/metadata.Release=$(TAG) \
		-X $(REPO_URL)/v2/internal/manager/metadata.Commit=$(COMMIT) \
		-X $(REPO_URL)/v2/internal/manager/metadata.Repo=$(REPO_INFO)" ${MAIN}

.PHONY: _build.debug
_build.debug:
	$(MAKE) _build.template.debug MAIN=./internal/cmd/main.go

.PHONY: _build.template.debug
_build.template.debug:
	go build -o bin/manager-debug -trimpath -gcflags=all="-N -l" -ldflags " \
		-X $(REPO_URL)/v2/internal/manager/metadata.Release=$(TAG) \
		-X $(REPO_URL)/v2/internal/manager/metadata.Commit=$(COMMIT) \
		-X $(REPO_URL)/v2/internal/manager/metadata.Repo=$(REPO_INFO)" ${MAIN}

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: lint
lint: verify.tidy golangci-lint staticcheck
	$(GOLANGCI_LINT) run -v

.PHONY: staticcheck
staticcheck: staticcheck.download
	$(STATICCHECK) -tags e2e_tests,integration_tests,istio_tests,conformance_tests \
		-f stylish \
		./...

.PHONY: verify.tidy
verify.tidy:
	./scripts/verify-tidy.sh

.PHONY: verify.repo
verify.repo:
	./scripts/verify-repo.sh

.PHONY: verify.diff
verify.diff:
	./scripts/verify-diff.sh

.PHONY: verify.versions
verify.versions:
	./scripts/verify-versions.sh $(TAG)

.PHONY: verify.manifests
verify.manifests: verify.repo manifests verify.diff

.PHONY: verify.generators
verify.generators: verify.repo generate verify.diff

# ------------------------------------------------------------------------------
# Build - Manifests
# ------------------------------------------------------------------------------

CRD_GEN_PATHS ?= ./...
CRD_OPTIONS ?= "+crd:allowDangerousTypes=true"

.PHONY: manifests
manifests: manifests.crds manifests.rbac manifests.single

.PHONY: manifests.crds
manifests.crds: controller-gen ## Generate WebhookConfiguration and CustomResourceDefinition objects.
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=kong-ingress webhook paths="$(CRD_GEN_PATHS)" output:crd:artifacts:config=config/crd/bases

.PHONY: manifests.rbac ## Generate ClusterRole objects.
manifests.rbac: controller-gen
	$(CONTROLLER_GEN) rbac:roleName=kong-ingress paths="./internal/controllers/configuration/"
	$(CONTROLLER_GEN) rbac:roleName=kong-ingress-knative paths="./internal/controllers/knative/" output:rbac:artifacts:config=config/rbac/knative
	$(CONTROLLER_GEN) rbac:roleName=kong-ingress-gateway paths="./internal/controllers/gateway/" output:rbac:artifacts:config=config/rbac/gateway

.PHONY: manifests.single
manifests.single: kustomize ## Compose single-file deployment manifests from building blocks
	./scripts/build-single-manifests.sh

# ------------------------------------------------------------------------------
# Build - Generators
# ------------------------------------------------------------------------------

.PHONY: generate
generate: generate.controllers generate.clientsets generate.gateway-api-urls generate.docs fmt

.PHONY: generate.controllers
generate.controllers: controller-gen
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="$(CRD_GEN_PATHS)"
	go generate $(PROJECT_DIR)/internal/cmd
# TODO: Previously this didn't have build tags assigned so technically nothing really
# happened upon go generate invocation for fips binary.
# Unfortunately this requires a bit more code to change the generation code since
# github.com/kong/kubernetes-ingress-controller/v2/hack/generators/controllers/networking
# relies on a relative path to boilerplate.go.txt which breaks if accessed from internal/cmd/fips.
# Related issue: https://github.com/Kong/kubernetes-ingress-controller/issues/2853
# go generate --tags fips $(PROJECT_DIR)/internal/cmd/fips

# this will generate the custom typed clients needed for end-users implementing logic in Go to use our API types.
.PHONY: generate.clientsets
generate.clientsets: client-gen
	$(CLIENT_GEN) \
		--go-header-file ./hack/boilerplate.go.txt \
		--logtostderr \
		--clientset-name clientset \
		--input-base $(REPO_URL)/v2/pkg/apis/  \
		--input configuration/v1,configuration/v1beta1,configuration/v1alpha1 \
		--input-dirs $(REPO_URL)/pkg/apis/configuration/v1alpha1/,$(REPO_URL)/pkg/apis/configuration/v1beta1/,$(REPO_URL)/pkg/apis/configuration/v1/ \
		--output-base pkg/ \
		--output-package $(REPO_URL)/v2/pkg/ \
		--trim-path-prefix pkg/$(REPO_URL)/v2/

.PHONY: generate.docs
generate.docs: crd-ref-docs
	./scripts/apidocs-gen/generate.sh $(CRD_REF_DOCS)

# ------------------------------------------------------------------------------
# Build - Container Images
# ------------------------------------------------------------------------------

REGISTRY ?= kong
IMGNAME ?= kubernetes-ingress-controller
IMAGE ?= $(REGISTRY)/$(IMGNAME)

.PHONY: container
container:
	docker buildx build \
		-f Dockerfile \
		--target distroless \
		--build-arg TAG=${TAG} \
		--build-arg COMMIT=${COMMIT} \
		--build-arg REPO_INFO=${REPO_INFO} \
		-t ${IMAGE}:${TAG} .

.PHONY: container.debug
container.debug:
	docker buildx build \
		-f Dockerfile.debug \
		--target debug \
		--build-arg TAG=${TAG}-debug \
		--build-arg COMMIT=${COMMIT} \
		--build-arg REPO_INFO=${REPO_INFO} \
		-t ${IMAGE}:${TAG} .

# ------------------------------------------------------------------------------
# Testing
# ------------------------------------------------------------------------------

NCPU ?= $(shell getconf _NPROCESSORS_ONLN)
PKG_LIST = ./pkg/...,./internal/...
INTEGRATION_TEST_TIMEOUT ?= "45m"
E2E_TEST_TIMEOUT ?= "45m"
E2E_TEST_RUN ?= ""
KONG_CONTROLLER_FEATURE_GATES ?= GatewayAlpha=true
GOTESTSUM_FORMAT ?= standard-verbose

.PHONY: test
test: test.unit

.PHONY: test.all
test.all: test.unit test.integration test.conformance

.PHONY: test.conformance
test.conformance: gotestsum
	@./scripts/check-container-environment.sh
	@TEST_DATABASE_MODE="off" GOFLAGS="-tags=conformance_tests" \
	GOTESTSUM_FORMAT=$(GOTESTSUM_FORMAT) \
	$(GOTESTSUM) -- -race \
		-timeout $(INTEGRATION_TEST_TIMEOUT) \
		-parallel $(NCPU) \
		./test/conformance

.PHONY: test.integration
test.integration: test.integration.dbless test.integration.postgres test.integration.cp

.PHONY: test.integration.enterprise
test.integration.enterprise: test.integration.enterprise.postgres

.PHONY: _test.unit
_test.unit: gotestsum
	GOTESTSUM_FORMAT=$(GOTESTSUM_FORMAT) \
	$(GOTESTSUM) -- -race $(GOTESTFLAGS) \
		-covermode=atomic \
		-coverpkg=$(PKG_LIST) \
		-coverprofile=coverage.unit.out \
		./internal/... \
		./pkg/...

.PHONY: test.unit
test.unit:
	@$(MAKE) _test.unit GOTESTFLAGS="$(GOTESTFLAGS)"

.PHONY: test.unit.pretty
test.unit.pretty:
	@$(MAKE) GOTESTSUM_FORMAT=pkgname _test.unit

.PHONY: _check.container.environment
_check.container.environment:
	@./scripts/check-container-environment.sh

.PHONY: _test.integration
_test.integration: _check.container.environment gotestsum
	TEST_DATABASE_MODE="$(DBMODE)" \
		GOFLAGS="-tags=$(GOTAGS)" \
		KONG_CONTROLLER_FEATURE_GATES=$(KONG_CONTROLLER_FEATURE_GATES) \
		GOTESTSUM_FORMAT=$(GOTESTSUM_FORMAT) \
		$(GOTESTSUM) -- $(GOTESTFLAGS) \
		-timeout $(INTEGRATION_TEST_TIMEOUT) \
		-parallel $(NCPU) \
		-race \
		-covermode=atomic \
		-coverpkg=$(PKG_LIST) \
		-coverprofile=$(COVERAGE_OUT) \
		./test/integration

.PHONY: test.integration.dbless.knative
test.integration.dbless.knative:
	@$(MAKE) _test.integration \
		GOTAGS="integration_tests,knative" \
		GOTESTFLAGS="-run TestKnative" \
		DBMODE=off \
		COVERAGE_OUT=coverage.dbless.knative.out

.PHONY: test.integration.dbless
test.integration.dbless:
	@$(MAKE) _test.integration \
		GOTAGS="integration_tests" \
		DBMODE=off \
		COVERAGE_OUT=coverage.dbless.out

.PHONY: test.integration.dbless.pretty
test.integration.dbless.pretty:
	@$(MAKE) GOTESTSUM_FORMAT=pkgname _test.integration \
		GOTAGS="integration_tests" \
		DBMODE=off \
		GOTESTFLAGS="-json" \
		COVERAGE_OUT=coverage.dbless.out

.PHONY: test.integration.postgres.knative
test.integration.postgres.knative:
	@$(MAKE) _test.integration \
		GOTAGS="integration_tests,knative" \
		GOTESTFLAGS="-run TestKnative" \
		DBMODE=postgres \
		COVERAGE_OUT=coverage.postgres.knative.out

.PHONY: test.integration.postgres
test.integration.postgres:
	@$(MAKE) _test.integration \
		GOTAGS="integration_tests" \
		DBMODE=postgres \
		COVERAGE_OUT=coverage.postgres.out

.PHONY: test.integration.postgres.pretty
test.integration.postgres.pretty:
	@$(MAKE) GOTESTSUM_FORMAT=pkgname _test.integration \
		GOTAGS="integration_tests" \
		DBMODE=postgres \
		GOTESTFLAGS="-json" \
		COVERAGE_OUT=coverage.postgres.out

.PHONY: test.integration.enterprise.postgres
test.integration.enterprise.postgres:
	@TEST_KONG_ENTERPRISE="true" \
		GOTAGS="integration_tests" \
		$(MAKE) _test.integration \
		DBMODE=postgres \
		COVERAGE_OUT=coverage.enterprisepostgres.out

.PHONY: test.integration.enterprise.postgres.pretty
test.integration.enterprise.postgres.pretty:
	@TEST_KONG_ENTERPRISE="true" \
		GOTAGS="integration_tests" \
		GOTESTSUM_FORMAT=pkgname \
		$(MAKE) _test.integration \
		DBMODE=postgres \
		GOTESTFLAGS="-json" \
		COVERAGE_OUT=coverage.enterprisepostgres.out

.PHONY: test.e2e
test.e2e: gotestsum
	GOFLAGS="-tags=e2e_tests" \
	GOTESTSUM_FORMAT=$(GOTESTSUM_FORMAT) \
	$(GOTESTSUM) -- $(GOTESTFLAGS) \
		-race \
		-run $(E2E_TEST_RUN) \
		-parallel $(NCPU) \
		-timeout $(E2E_TEST_TIMEOUT) \
		./test/e2e/...

.PHONY: test.istio
test.istio: gotestsum
	ISTIO_TEST_ENABLED="true" \
	GOTESTSUM_FORMAT=$(GOTESTSUM_FORMAT) \
	GOFLAGS="-tags=istio_tests" $(GOTESTSUM) -- $(GOTESTFLAGS) \
		-race \
		-parallel $(NCPU) \
		-timeout $(E2E_TEST_TIMEOUT) \
		./test/e2e/...

# ------------------------------------------------------------------------------
# Operations - Local Deployment
# ------------------------------------------------------------------------------

# NOTE:
# The environment used for "make debug" or "make run" by default should
# have a Kong Gateway deployed into the kong-system namespace, but these
# defaults can be changed using the arguments below.
#
# One easy way to create a testing/debugging environment that works with
# these defaults is to use the Kong Kubernetes Testing Framework (KTF):
#
# $ ktf envs create --addon metallb --addon kong --kong-disable-controller --kong-admin-service-loadbalancer
#
# KTF can be installed by following the instructions at:
# https://github.com/kong/kubernetes-testing-framework
#
# Alternatively one can use Kong's helm chart to deploy it on the cluster, using
# for example the following set of flags:
#   helm upgrade --create-namespace --install --namespace kong-system kong kong/kong \
#       --set ingressController.enabled=false \
#       --set admin.enabled=true \
#       --set admin.type=LoadBalancer \
#       --set admin.http.enabled=true \
#       --set admin.tls.enabled=false
#
# https://github.com/Kong/charts/tree/main/charts/kong

KUBECONFIG ?= "${HOME}/.kube/config"
KONG_NAMESPACE ?= kong-system
KONG_PROXY_SERVICE ?= ingress-controller-kong-proxy
KONG_PROXY_UDP_SERVICE ?= ingress-controller-kong-udp-proxy
KONG_ADMIN_SERVICE ?= ingress-controller-kong-admin
KONG_ADMIN_PORT ?= 8001
KONG_ADMIN_URL ?= "http://$(shell kubectl -n $(KONG_NAMESPACE) get service $(KONG_ADMIN_SERVICE) -o=go-template='{{range .status.loadBalancer.ingress}}{{.ip}}{{end}}'):$(KONG_ADMIN_PORT)"

.PHONY: _ensure-namespace
_ensure-namespace:
	@kubectl create ns $(KONG_NAMESPACE) 2>/dev/null || true

.PHONY: debug
debug: install _ensure-namespace
	$(DLV) debug ./internal/cmd/main.go -- \
		--anonymous-reports=false \
		--kong-admin-url $(KONG_ADMIN_URL) \
		--publish-service $(KONG_NAMESPACE)/$(KONG_PROXY_SERVICE) \
		--publish-service-udp $(KONG_NAMESPACE)/$(KONG_PROXY_UDP_SERVICE) \
		--kubeconfig $(KUBECONFIG) \
		--feature-gates=$(KONG_CONTROLLER_FEATURE_GATES)

# By default dlv will look for a config in:
# > If $XDG_CONFIG_HOME is set, then configuration and command history files are
# > located in $XDG_CONFIG_HOME/dlv.
# > Otherwise, they are located in $HOME/.config/dlv on Linux and $HOME/.dlv on other systems.
#
# ref: https://github.com/go-delve/delve/blob/master/Documentation/cli/README.md#configuration-and-command-history
# 
# This sets the XDG_CONFIG_HOME to this project's subdirectory so that project
# specific substitution paths can be isolated to this project only and not shared
# across projects under $HOME or common XDG_CONFIG_HOME.
.PHONY: debug.connect
debug.connect:
	XDG_CONFIG_HOME="$(PROJECT_DIR)/.config" $(DLV) connect localhost:40000

SKAFFOLD_DEBUG_PROFILE ?= debug

# This will port-forward 40000 from KIC's debugger to localhost. Connect to that
# port with debugger/IDE of your choice
.PHONY: debug.skaffold
debug.skaffold:
	TAG=$(TAG)-debug REPO_INFO=$(REPO_INFO) COMMIT=$(COMMIT) \
		CMD=debug \
		SKAFFOLD_PROFILE=$(SKAFFOLD_DEBUG_PROFILE) \
		$(MAKE) _skaffold

# This will port-forward 40000 from KIC's debugger to localhost. Connect to that
# port with debugger/IDE of your choice.
#
# To make it work with Konnect, you must provide following files under ./config/variants/konnect/debug:
#   * `konnect.env` with CONTROLLER_KONNECT_RUNTIME_GROUP_ID env variable set
#     to the UUID of a Runtime Group you have created in Konnect.
#   * `tls.crt` and `tls.key` with TLS client cerificate and its key (generated by Konnect).
.PHONY: debug.skaffold.konnect
debug.skaffold.konnect:
	SKAFFOLD_DEBUG_PROFILE=debug-konnect \
		$(MAKE) debug.skaffold

# This will port-forward 40000 from KIC's debugger to localhost. Connect to that
# port with debugger/IDE of your choice
.PHONY: debug.skaffold.sync
debug.skaffold.sync:
	$(MAKE) debug.skaffold SKAFFOLD_FLAGS="--auto-build --auto-deploy --auto-sync"

SKAFFOLD_RUN_PROFILE ?= dev

.PHONY: run.skaffold
run.skaffold:
	TAG=$(TAG) REPO_INFO=$(REPO_INFO) COMMIT=$(COMMIT) \
		CMD=dev \
		SKAFFOLD_PROFILE=$(SKAFFOLD_RUN_PROFILE) \
		$(MAKE) _skaffold

.PHONY: _skaffold
_skaffold: skaffold
	$(SKAFFOLD) $(CMD) --port-forward=pods --profile=$(SKAFFOLD_PROFILE) $(SKAFFOLD_FLAGS)

.PHONY: run
run: install _ensure-namespace
	@$(MAKE) _run

# This target can be used to skip all the precondition checks, code generation
# and other logic around running the controller.
# It should be run only after the cluster has been already prepared to run with KIC.
.PHONY: _run
_run:
	go run ./internal/cmd/main.go \
		--anonymous-reports=false \
		--kong-admin-url $(KONG_ADMIN_URL) \
		--publish-service $(KONG_NAMESPACE)/$(KONG_PROXY_SERVICE) \
		--publish-service-udp $(KONG_NAMESPACE)/$(KONG_PROXY_UDP_SERVICE) \
		--kubeconfig $(KUBECONFIG) \
		--feature-gates=$(KONG_CONTROLLER_FEATURE_GATES)

# ------------------------------------------------------------------------------
# Gateway API
# ------------------------------------------------------------------------------

# GATEWAY_API_VERSION will be processed by kustomize and therefore accepts
# only branch names, tags, or full commit hashes, i.e. short hashes or go
# pseudo versions are not supported [1].
# Please also note that kustomize fails silently when provided with an
# unsupported ref and downloads the manifests from the main branch.
#
# [1]: https://github.com/kubernetes-sigs/kustomize/blob/master/examples/remoteBuild.md#remote-directories
GATEWAY_API_PACKAGE ?= sigs.k8s.io/gateway-api
GATEWAY_API_RELEASE_CHANNEL ?= experimental
GATEWAY_API_VERSION ?= $(shell go list -m -f '{{ .Version }}' $(GATEWAY_API_PACKAGE))
GATEWAY_API_CRDS_LOCAL_PATH = $(shell go env GOPATH)/pkg/mod/$(GATEWAY_API_PACKAGE)@$(GATEWAY_API_VERSION)/config/crd
GATEWAY_API_REPO ?= github.com/kubernetes-sigs/gateway-api
GATEWAY_API_RAW_REPO ?= https://raw.githubusercontent.com/kubernetes-sigs/gateway-api
GATEWAY_API_CRDS_URL = $(GATEWAY_API_REPO)/config/crd/$(GATEWAY_API_RELEASE_CHANNEL)?ref=$(GATEWAY_API_VERSION)
GATEWAY_API_RAW_REPO_URL = $(GATEWAY_API_RAW_REPO)/$(GATEWAY_API_VERSION)

.PHONY: print-gateway-api-crds-url
print-gateway-api-crds-url:
	@echo $(GATEWAY_API_CRDS_URL)

.PHONY: print-gateway-api-raw-repo-url
print-gateway-api-raw-repo-url:
	@echo $(GATEWAY_API_RAW_REPO_URL)

.PHONY: generate.gateway-api-urls
generate.gateway-api-urls:
	CRDS_STANDARD_URL=$(shell GATEWAY_API_RELEASE_CHANNEL="" $(MAKE) print-gateway-api-crds-url)\
		CRDS_EXPERIMENTAL_URL=$(shell GATEWAY_API_RELEASE_CHANNEL="experimental" $(MAKE) print-gateway-api-crds-url) \
		RAW_REPO_URL=$(shell $(MAKE) print-gateway-api-raw-repo-url) \
		INPUT=$(shell pwd)/test/internal/cmd/generate-gateway-api-urls/gateway_consts.tmpl \
		OUTPUT=$(shell pwd)/test/consts/zz_generated_gateway.go \
		go generate -tags=generate_gateway_api_urls ./test/internal/cmd/generate-gateway-api-urls

.PHONY: go-mod-download-gateway-api
go-mod-download-gateway-api:
	@go mod download $(GATEWAY_API_PACKAGE)

.PHONY: install-gateway-api-crds
install-gateway-api-crds: go-mod-download-gateway-api kustomize
	$(KUSTOMIZE) build $(GATEWAY_API_CRDS_LOCAL_PATH) | kubectl apply -f -
	$(KUSTOMIZE) build $(GATEWAY_API_CRDS_LOCAL_PATH)/experimental | kubectl apply -f -

.PHONY: uninstall-gateway-api-crds
uninstall-gateway-api-crds: go-mod-download-gateway-api kustomize
	$(KUSTOMIZE) build $(GATEWAY_API_CRDS_LOCAL_PATH) | kubectl delete -f -
	$(KUSTOMIZE) build $(GATEWAY_API_CRDS_LOCAL_PATH)/experimental | kubectl delete -f -

# Install CRDs into the K8s cluster specified in $KUBECONFIG.
.PHONY: install
install: manifests install-gateway-api-crds
	$(KUSTOMIZE) build config/crd | kubectl apply -f -

# Uninstall CRDs from the K8s cluster specified in $KUBECONFIG.
.PHONY: uninstall
uninstall: manifests uninstall-gateway-api-crds
	$(KUSTOMIZE) build config/crd | kubectl delete -f -

deploy: manifests kustomize ## Deploy controller to the K8s cluster specified in $KUBECONFIG.
	cd config/manager && $(KUSTOMIZE) edit set image controller=${IMAGE}
	$(KUSTOMIZE) build config/default | kubectl apply -f -

undeploy: ## Undeploy controller from the K8s cluster specified in $KUBECONFIG.
	$(KUSTOMIZE) build config/default | kubectl delete -f -
