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

CONTROLLER_GEN = $(shell pwd)/bin/controller-gen
controller-gen: ## Download controller-gen locally if necessary.
	$(call go-get-tool,$(CONTROLLER_GEN),sigs.k8s.io/controller-tools/cmd/controller-gen@v0.7.0)

KUSTOMIZE = $(shell pwd)/bin/kustomize
kustomize: ## Download kustomize locally if necessary.
	$(call go-get-tool,$(KUSTOMIZE),sigs.k8s.io/kustomize/kustomize/v4@v4.3.0)

CLIENT_GEN = $(shell pwd)/bin/client-gen
client-gen: ## Download client-gen locally if necessary.
	$(call go-get-tool,$(CLIENT_GEN),k8s.io/code-generator/cmd/client-gen@v0.21.3)

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
lint: verify.tidy
	golangci-lint run -v

.PHONY: verify.tidy
verify.tidy:
	./hack/verify-tidy.sh

.PHONY: verify.repo
verify.repo:
	./hack/verify-repo.sh

.PHONY: verify.diff
verify.diff:
	./hack/verify-diff.sh

.PHONY: verify.versions
verify.versions:
	./hack/verify-versions.sh $(TAG)

.PHONY: verify.manifests
verify.manifests: verify.repo manifests manifests.single verify.diff

.PHONY: verify.generators
verify.generators: verify.repo generate verify.diff

# ------------------------------------------------------------------------------
# Build - Manifests
# ------------------------------------------------------------------------------

CRD_OPTIONS ?= "+crd:allowDangerousTypes=true"

.PHONY: manifests
manifests: manifests.crds manifests.single

.PHONY: manifests.crds
manifests.crds: controller-gen ## Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects.
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=kong-ingress webhook paths="./..." output:crd:artifacts:config=config/crd/bases

.PHONY: manifests.single
manifests.single: kustomize ## Compose single-file deployment manifests from building blocks
	./hack/deploy/build-single-manifests.sh

# ------------------------------------------------------------------------------
# Build - Generators
# ------------------------------------------------------------------------------

.PHONY: generate
generate: generate.controllers generate.clientsets

.PHONY: generate.controllers
generate.controllers: controller-gen
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."
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
    --build-arg TAG=${TAG} --build-arg COMMIT=${COMMIT} \
    --build-arg REPO_INFO=${REPO_INFO} \
    -t ${IMAGE}:${TAG} .

.PHONY: container
debug-container:
	docker buildx build \
    -f Dockerfile \
    --target debug \
    --build-arg TAG=${TAG}-debug --build-arg COMMIT=${COMMIT} \
    --build-arg REPO_INFO=${REPO_INFO} \
    -t ${IMAGE}:${TAG} .

# ------------------------------------------------------------------------------
# Testing
# ------------------------------------------------------------------------------

NCPU ?= $(shell getconf _NPROCESSORS_ONLN)
PKG_LIST = ./pkg/...,./internal/...
KIND_CLUSTER_NAME ?= "integration-tests"
INTEGRATION_TEST_TIMEOUT ?= "20m"
E2E_TEST_TIMEOUT ?= "30m"

.PHONY: test
test: test.unit

.PHONY: test.all
test.all: test.unit test.integration

.PHONY: test.integration
test.integration: test.integration.enterprise.postgres  test.integration.dbless test.integration.postgres

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
	@TEST_DATABASE_MODE="off" GOFLAGS="-tags=integration_tests" go test -v -race \
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
	@TEST_DATABASE_MODE="postgres" GOFLAGS="-tags=integration_tests" go test -v \
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
	@TEST_DATABASE_MODE="postgres" TEST_KONG_ENTERPRISE="true" GOFLAGS="-tags=integration_tests" go test -v \
		-timeout $(INTEGRATION_TEST_TIMEOUT) \
		-parallel $(NCPU) \
		-race \
		-covermode=atomic \
		-coverpkg=$(PKG_LIST) \
		-coverprofile=coverage.enterprisepostgres.out \
		./test/integration

.PHONY: test.integration.legacy
test.integration.legacy: container
	KIC_IMAGE="${IMAGE}:${TAG}" KUBE_VERSION=${KUBE_VERSION} ./hack/legacy/test/test.sh

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

run: manifests generate fmt vet ## Run a controller from your host.
	go run ./internal/cmd/main.go

install: manifests kustomize ## Install CRDs into the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/crd | kubectl apply -f -

uninstall: manifests kustomize ## Uninstall CRDs from the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/crd | kubectl delete -f -

deploy: manifests kustomize ## Deploy controller to the K8s cluster specified in ~/.kube/config.
	cd config/manager && $(KUSTOMIZE) edit set image controller=${IMAGE}
	$(KUSTOMIZE) build config/default | kubectl apply -f -

undeploy: ## Undeploy controller from the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/default | kubectl delete -f -
