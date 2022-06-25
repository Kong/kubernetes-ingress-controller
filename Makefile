# ------------------------------------------------------------------------------
# Configuration - Repository
# ------------------------------------------------------------------------------

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
	(cd third_party && go mod tidy ) && \
		GOBIN=$(PROJECT_DIR)/bin go install -modfile third_party/go.mod $(TOOL)

CONTROLLER_GEN = $(PROJECT_DIR)/bin/controller-gen
controller-gen: ## Download controller-gen locally if necessary.
	$(MAKE) _download_tool TOOL=sigs.k8s.io/controller-tools/cmd/controller-gen

KUSTOMIZE = $(PROJECT_DIR)/bin/kustomize
kustomize: ## Download kustomize locally if necessary.
	$(MAKE) _download_tool TOOL=sigs.k8s.io/kustomize/kustomize/v4

CLIENT_GEN = $(PROJECT_DIR)/bin/client-gen
client-gen: ## Download client-gen locally if necessary.
	$(MAKE) _download_tool TOOL=k8s.io/code-generator/cmd/client-gen

GOLANGCI_LINT = $(PROJECT_DIR)/bin/golangci-lint
golangci-lint: ## Download golangci-lint locally if necessary.
	$(MAKE) _download_tool TOOL=github.com/golangci/golangci-lint/cmd/golangci-lint

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
build: generate fmt vet lint
	go build -a -o bin/manager -ldflags "-s -w \
		-X github.com/kong/kubernetes-ingress-controller/v2/internal/metadata.Release=$(TAG) \
		-X github.com/kong/kubernetes-ingress-controller/v2/internal/metadata.Commit=$(COMMIT) \
		-X github.com/kong/kubernetes-ingress-controller/v2/internal/metadata.Repo=$(REPO_INFO)" internal/cmd/main.go

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: lint
lint: verify.tidy golangci-lint
	$(GOLANGCI_LINT) run -v

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
verify.manifests: verify.repo manifests manifests.single verify.diff

.PHONY: verify.generators
verify.generators: verify.repo generate verify.diff

# ------------------------------------------------------------------------------
# Build - Manifests
# ------------------------------------------------------------------------------

CRD_GEN_PATHS ?= ./...
CRD_OPTIONS ?= "+crd:allowDangerousTypes=true"

.PHONY: manifests
manifests: manifests.crds manifests.single

.PHONY: manifests.crds
manifests.crds: controller-gen ## Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects.
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=kong-ingress webhook paths="$(CRD_GEN_PATHS)" output:crd:artifacts:config=config/crd/bases

.PHONY: manifests.single
manifests.single: kustomize ## Compose single-file deployment manifests from building blocks
	./scripts/build-single-manifests.sh

# ------------------------------------------------------------------------------
# Build - Generators
# ------------------------------------------------------------------------------

.PHONY: generate
generate: generate.controllers generate.clientsets generate.gateway-api-crds-url

.PHONY: generate.controllers
generate.controllers: controller-gen
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="$(CRD_GEN_PATHS)"
	go generate ./...

# this will generate the custom typed clients needed for end-users implementing logic in Go to use our API types.
# TODO: we're hacking around client-gen for now to enable it for enabled go modules, should probably contribute upstream to improve this.
#       See: https://github.com/Kong/kubernetes-ingress-controller/issues/1254
.PHONY: generate.clientsets
generate.clientsets: client-gen
	@$(CLIENT_GEN) --go-header-file ./hack/boilerplate.go.txt \
		--clientset-name clientset \
		--input-base github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/  \
		--input configuration/v1,configuration/v1beta1 \
		--input-dirs github.com/kong/kubernetes-ingress-controller/pkg/apis/configuration/v1beta1/,github.com/kong/kubernetes-ingress-controller/pkg/apis/configuration/v1/ \
		--output-base client-gen-tmp/ \
		--output-package github.com/kong/kubernetes-ingress-controller/v2/pkg/
	@rm -rf pkg/clientset/
	@mv client-gen-tmp/github.com/kong/kubernetes-ingress-controller/v2/pkg/clientset pkg/
	@rm -rf client-gen-tmp/

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

.PHONY: container
debug-container:
	docker buildx build \
		-f Dockerfile \
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
KIND_CLUSTER_NAME ?= "integration-tests"
INTEGRATION_TEST_TIMEOUT ?= "45m"
E2E_TEST_TIMEOUT ?= "45m"
KONG_CONTROLLER_FEATURE_GATES ?= Gateway=true

.PHONY: test
test: test.unit

.PHONY: test.all
test.all: test.unit test.integration test.conformance

.PHONY: test.conformance
test.conformance:
	@./scripts/check-container-environment.sh
	@TEST_DATABASE_MODE="off" GOFLAGS="-tags=conformance_tests" go test -v -race \
		-timeout $(INTEGRATION_TEST_TIMEOUT) \
		-parallel $(NCPU) \
		-race \
		./test/conformance

.PHONY: test.integration
test.integration: test.integration.dbless test.integration.postgres

.PHONY: test.integration.enterprise
test.integration.enterprise: test.integration.enterprise.postgres

.PHONY: test.unit
test.unit:
	@go test -v -race \
		-covermode=atomic \
		-coverpkg=$(PKG_LIST) \
		-coverprofile=coverage.unit.out \
		./internal/... \
		./pkg/...

.PHONY: test.integration.dbless
test.integration.dbless:
	@./scripts/check-container-environment.sh
	@TEST_DATABASE_MODE="off" GOFLAGS="-tags=integration_tests" KONG_CONTROLLER_FEATURE_GATES=$(KONG_CONTROLLER_FEATURE_GATES) \
		go test -v -race \
		-timeout $(INTEGRATION_TEST_TIMEOUT) \
		-parallel $(NCPU) \
		-race \
		-covermode=atomic \
		-coverpkg=$(PKG_LIST) \
		-coverprofile=coverage.dbless.out \
		./test/integration

.PHONY: test.integration.postgres
test.integration.postgres:
	@./scripts/check-container-environment.sh
	@TEST_DATABASE_MODE="postgres" GOFLAGS="-tags=integration_tests" KONG_CONTROLLER_FEATURE_GATES=$(KONG_CONTROLLER_FEATURE_GATES) \
		go test -v \
		-timeout $(INTEGRATION_TEST_TIMEOUT) \
		-parallel $(NCPU) \
		-race \
		-covermode=atomic \
		-coverpkg=$(PKG_LIST) \
		-coverprofile=coverage.postgres.out \
		./test/integration

.PHONY: test.integration.enterprise.postgres
test.integration.enterprise.postgres:
	@./scripts/check-container-environment.sh
	@TEST_DATABASE_MODE="postgres" TEST_KONG_ENTERPRISE="true" GOFLAGS="-tags=integration_tests" KONG_CONTROLLER_FEATURE_GATES=$(KONG_CONTROLLER_FEATURE_GATES) \
		go test -v \
		-timeout $(INTEGRATION_TEST_TIMEOUT) \
		-parallel $(NCPU) \
		-race \
		-covermode=atomic \
		-coverpkg=$(PKG_LIST) \
		-coverprofile=coverage.enterprisepostgres.out \
		./test/integration

.PHONY: test.e2e
test.e2e:
	GOFLAGS="-tags=e2e_tests" go test -v \
		-race \
		-parallel $(NCPU) \
		-timeout $(E2E_TEST_TIMEOUT) \
		./test/e2e/...

# ------------------------------------------------------------------------------
# Operations - Local Deployment
# ------------------------------------------------------------------------------

# NOTE: the environment used for "make debug" or "make run" by default should
#       have a Kong Gateway deployed into the kong-system namespace, but these
#       defaults can be changed using the arguments below.
#
#       One easy way to create a testing/debugging environment that works with
#       these defaults is to use the Kong Kubernetes Testing Framework (KTF):
#
#       $ ktf envs create --addon metallb --addon kong --kong-disable-controller --kong-admin-service-loadbalancer
#
#       KTF can be installed by following the instructions at:
#
#       https://github.com/kong/kubernetes-testing-framework

KUBECONFIG ?= "${HOME}/.kube/config"
KONG_NAMESPACE ?= kong-system
KONG_PROXY_SERVICE ?= ingress-controller-kong-proxy
KONG_ADMIN_PORT ?= 8001
KONG_ADMIN_URL ?= "http://$(shell kubectl -n kong-system get service ingress-controller-kong-admin -o=go-template='{{range .status.loadBalancer.ingress}}{{.ip}}{{end}}'):$(KONG_ADMIN_PORT)"

debug: install
	dlv debug ./internal/cmd/main.go -- \
		--kong-admin-url $(KONG_ADMIN_URL) \
		--publish-service $(KONG_NAMESPACE)/$(KONG_PROXY_SERVICE) \
		--kubeconfig $(KUBECONFIG) \
		--feature-gates=$(KONG_CONTROLLER_FEATURE_GATES)

run: install
	go run ./internal/cmd/main.go \
		--kong-admin-url $(KONG_ADMIN_URL) \
		--publish-service $(KONG_NAMESPACE)/$(KONG_PROXY_SERVICE) \
		--kubeconfig $(KUBECONFIG) \
		--feature-gates=$(KONG_CONTROLLER_FEATURE_GATES)

# ------------------------------------------------------------------------------
# Gateway API
# ------------------------------------------------------------------------------

GATEWAY_API_PACKAGE ?= sigs.k8s.io/gateway-api
# TODO: Below hardcoded ref is a workaround for the fact that we're using an untagged version
#       of sigs.k8s.io/gateway-api in go.mod - that occurred after v0.4.0 (which was tagged on master)
#       but which contains a breaking change w.r.t to the file structure in said repo - and the
#       fact that kustomize accepts only branch names, tags, or full commit hashes, i.e. short
#       hashes or go pseudo versions are not supported [1].
#       Please also note that kustomize fails silently when provided with an unsupported ref
#       and downloads the manifests from the main branch.
#
#       [1]: https://github.com/kubernetes-sigs/kustomize/blob/master/examples/remoteBuild.md#remote-directories
#
#       This causes a problem where we cannot use go pseudo version from go.mod i.e.
#       v0.4.1-0.20220306235253-71fee1c2808f and where we cannot update to a newer version
#       sigs.k8s.io/gateway-api because v0.5.0 hasn't been released yet and v0.4.x versions
#       do not contain the change in file structure that some of the code in this repo already
#       relies on.
#
#       In order to avoid unnecessary work we're just hardcoding the full SHA that
#       corresponds to what's in go.mod - v0.4.1-0.20220306235253-71fee1c2808f - until
#       v0.5.0 is released which we can then use in go.mod and scrape via go list ...
# 
#       Whenever the above happens the hardcoded SHA can be replaced with:
#       $(shell go list -m -f "{{.Version}}" $(GATEWAY_API_PACKAGE))
#
#       Related issue: https://github.com/Kong/kubernetes-ingress-controller/issues/2595
GATEWAY_API_VERSION ?= 71fee1c2808fa19a5f19d952d155fc072cf9324c
GATEWAY_API_CRDS_LOCAL_PATH = $(shell go env GOPATH)/pkg/mod/$(GATEWAY_API_PACKAGE)@$(GATEWAY_API_VERSION)/config/crd
GATEWAY_API_REPO ?= github.com/kubernetes-sigs/gateway-api
GATEWAY_API_CRDS_URL = $(GATEWAY_API_REPO)/config/crd?ref=$(GATEWAY_API_VERSION)

.PHONY: print-gateway-api-crds-url
print-gateway-api-crds-url:
	@echo $(GATEWAY_API_CRDS_URL)

.PHONY: generate.gateway-api-crds-url
generate.gateway-api-crds-url:
	URL=$(shell $(MAKE) print-gateway-api-crds-url) \
		INPUT=$(shell pwd)/test/internal/cmd/generate-gateway-api-crds-url/gateway_consts.tmpl \
		OUTPUT=$(shell pwd)/test/consts/gateway.go \
		go generate ./test/internal/cmd/generate-gateway-api-crds-url

.PHONY: go-mod-download-gateway-api
go-mod-download-gateway-api:
	@go mod download $(GATEWAY_API_PACKAGE)

.PHONY: install-gateway-api-crds
install-gateway-api-crds: go-mod-download-gateway-api
	$(KUSTOMIZE) build $(GATEWAY_API_CRDS_LOCAL_PATH) | kubectl apply -f -

.PHONY: uninstall-gateway-api-crds
uninstall-gateway-api-crds: go-mod-download-gateway-api
	$(KUSTOMIZE) build $(GATEWAY_API_CRDS_LOCAL_PATH) | kubectl delete -f -

# Install CRDs into the K8s cluster specified in $KUBECONFIG.
.PHONY: install
install: manifests kustomize install-gateway-api-crds
	$(KUSTOMIZE) build config/crd | kubectl apply -f -

# Uninstall CRDs from the K8s cluster specified in $KUBECONFIG.
.PHONY: uninstall
uninstall: manifests kustomize uninstall-gateway-api-crds
	$(KUSTOMIZE) build config/crd | kubectl delete -f -

deploy: manifests kustomize ## Deploy controller to the K8s cluster specified in $KUBECONFIG.
	cd config/manager && $(KUSTOMIZE) edit set image controller=${IMAGE}
	$(KUSTOMIZE) build config/default | kubectl apply -f -

undeploy: ## Undeploy controller from the K8s cluster specified in $KUBECONFIG.
	$(KUSTOMIZE) build config/default | kubectl delete -f -
