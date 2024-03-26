# ------------------------------------------------------------------------------
# Configuration - Make
# ------------------------------------------------------------------------------

# Some sensible Make defaults: https://tech.davis-hansson.com/p/make/
SHELL := bash
.SHELLFLAGS := -eu -o pipefail -c

# ------------------------------------------------------------------------------
# Configuration - Repository
# ------------------------------------------------------------------------------

MAKEFLAGS += --no-print-directory
REPO_URL ?= github.com/kong/kubernetes-ingress-controller
REPO_INFO ?= $(shell git config --get remote.origin.url)
GO_MOD_MAJOR_VERSION ?= $(subst $(REPO_URL)/,,$(shell go list -m))
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

TOOLS_VERSIONS_FILE = .tools_versions.yaml

MISE := $(shell which mise)
.PHONY: mise
mise:
	@mise -V >/dev/null || (echo "mise - https://github.com/jdx/mise - not found. Please install it." && exit 1)

.PHONY: tools
tools: controller-gen kustomize client-gen golangci-lint.download gotestsum crd-ref-docs skaffold looppointer.download staticcheck.download

export MISE_DATA_DIR = $(PROJECT_DIR)/bin/

CONTROLLER_GEN_VERSION = $(shell yq -ojson -r '.controller-tools' < $(TOOLS_VERSIONS_FILE))
CONTROLLER_GEN = $(PROJECT_DIR)/bin/installs/kube-controller-tools/$(CONTROLLER_GEN_VERSION)/bin/controller-gen
.PHONY: controller-gen
controller-gen: mise ## Download controller-gen locally if necessary.
	@$(MISE) plugin install --yes -q kube-controller-tools
	@$(MISE) install -q kube-controller-tools@$(CONTROLLER_GEN_VERSION)

KUSTOMIZE_VERSION = $(shell yq -ojson -r '.kustomize' < $(TOOLS_VERSIONS_FILE))
KUSTOMIZE = $(PROJECT_DIR)/bin/installs/kustomize/$(KUSTOMIZE_VERSION)/bin/kustomize
.PHONY: kustomize
kustomize: mise ## Download kustomize locally if necessary.
	@$(MISE) plugin install --yes -q kustomize
	@$(MISE) install -q kustomize@$(KUSTOMIZE_VERSION)

CLIENT_GEN_VERSION = $(shell yq -ojson -r '.code-generator' < $(TOOLS_VERSIONS_FILE))
CLIENT_GEN = $(PROJECT_DIR)/bin/installs/kube-code-generator/$(CLIENT_GEN_VERSION)/bin/client-gen
.PHONY: client-gen
client-gen: mise ## Download client-gen locally if necessary.
	@$(MISE) plugin install --yes -q kube-code-generator
	@$(MISE) install -q kube-code-generator@$(CLIENT_GEN_VERSION)

GOLANGCI_LINT_VERSION = $(shell yq -ojson -r '.golangci-lint' < $(TOOLS_VERSIONS_FILE))
GOLANGCI_LINT = $(PROJECT_DIR)/bin/installs/golangci-lint/$(GOLANGCI_LINT_VERSION)/bin/golangci-lint
.PHONY: golangci-lint.download
golangci-lint.download: mise ## Download golangci-lint locally if necessary.
	@$(MISE) plugin install --yes -q golangci-lint
	@$(MISE) install -q golangci-lint@$(GOLANGCI_LINT_VERSION)

GOTESTSUM_VERSION = $(shell yq -ojson -r '.gotestsum' < $(TOOLS_VERSIONS_FILE))
GOTESTSUM = $(PROJECT_DIR)/bin/installs/gotestsum/$(GOTESTSUM_VERSION)/bin/gotestsum
.PHONY: gotestsum
gotestsum: ## Download gotestsum locally if necessary.
	@$(MISE) plugin install --yes -q gotestsum https://github.com/pmalek/mise-gotestsum.git
	@$(MISE) install -q gotestsum

CRD_REF_DOCS_VERSION = $(shell yq -ojson -r '.crd-ref-docs' < $(TOOLS_VERSIONS_FILE))
CRD_REF_DOCS = $(PROJECT_DIR)/bin/crd-ref-docs
.PHONY: crd-ref-docs
crd-ref-docs: ## Download crd-ref-docs locally if necessary.
	GOBIN=$(PROJECT_DIR)/bin go install -v \
		github.com/elastic/crd-ref-docs@v$(CRD_REF_DOCS_VERSION)

SKAFFOLD_VERSION = $(shell yq -ojson -r '.skaffold' < $(TOOLS_VERSIONS_FILE))
SKAFFOLD = $(PROJECT_DIR)/bin/installs/skaffold/$(SKAFFOLD_VERSION)/bin/skaffold
.PHONY: skaffold
skaffold: mise ## Download skaffold locally if necessary.
	@$(MISE) plugin install --yes -q skaffold
	@$(MISE) install -q skaffold@$(SKAFFOLD_VERSION)

YQ_VERSION = $(shell yq -ojson -r '.yq' < $(TOOLS_VERSIONS_FILE))
YQ = $(PROJECT_DIR)/bin/installs/yq/$(YQ_VERSION)/bin/yq
.PHONY: yq
yq: mise # Download yq locally if necessary.
	@$(MISE) plugin install --yes -q yq
	@$(MISE) install -q yq@$(YQ_VERSION)

DLV = $(PROJECT_DIR)/bin/dlv
.PHONY: dlv
dlv: ## Download dlv locally if necessary.
	@$(MAKE) _download_tool TOOL=dlv

SETUP_ENVTEST_VERSION = $(shell yq -ojson -r '.setup-envtest' < $(TOOLS_VERSIONS_FILE))
SETUP_ENVTEST = $(PROJECT_DIR)/bin/installs/setup-envtest/$(SETUP_ENVTEST_VERSION)/bin/setup-envtest
.PHONY: setup-envtest
setup-envtest: mise ## Download setup-envtest locally if necessary.
	@$(MISE) plugin install --yes -q setup-envtest https://github.com/pmalek/mise-setup-envtest.git
	@$(MISE) install setup-envtest@v$(SETUP_ENVTEST_VERSION)

STATICCHECK = $(PROJECT_DIR)/bin/staticcheck
.PHONY: staticcheck.download
staticcheck.download: ## Download staticcheck locally if necessary.
	@$(MAKE) _download_tool TOOL=staticcheck

GOJUNIT_REPORT_VERSION = $(shell yq -ojson -r '.gojunit-report' < $(TOOLS_VERSIONS_FILE))
GOJUNIT_REPORT = $(PROJECT_DIR)/bin/installs/go-junit-report/$(GOJUNIT_REPORT_VERSION)/bin/go-junit-report
.PHONY: go-junit-report
go-junit-report: ## Download go-junit-report locally if necessary.
# TODO: Go back to using https://github.com/jwillker/asdf-go-junit-report when https://github.com/jwillker/asdf-go-junit-report/pull/4 merges.
	@$(MISE) plugin install --yes -q go-junit-report https://github.com/pmalek/asdf-go-junit-report.git
	@$(MISE) install go-junit-report@v$(GOJUNIT_REPORT_VERSION)

LOOPPOINTER= $(PROJECT_DIR)/bin/looppointer
.PHONY: looppointer.download
looppointer.download: ## Download looppointer locally if necessary.
	@$(MAKE) _download_tool TOOL=looppointer

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
build: generate _build

.PHONY: build.fips
build.fips: generate lint _build.fips

.PHONY: _build
_build:
	$(MAKE) _build.template MAIN=./internal/cmd/main.go

.PHONY: _build.fips
_build.fips:
	$(MAKE) _build.template MAIN=./internal/cmd/fips/main.go

.PHONY: _build.template
_build.template:
	go build -o bin/manager -ldflags "-s -w \
		-X $(REPO_URL)/$(GO_MOD_MAJOR_VERSION)/internal/manager/metadata.Release=$(TAG) \
		-X $(REPO_URL)/$(GO_MOD_MAJOR_VERSION)/internal/manager/metadata.Commit=$(COMMIT) \
		-X $(REPO_URL)/$(GO_MOD_MAJOR_VERSION)/internal/manager/metadata.Repo=$(REPO_INFO)" ${MAIN}

.PHONY: _build.debug
_build.debug:
	$(MAKE) _build.template.debug MAIN=./internal/cmd/main.go

.PHONY: _build.template.debug
_build.template.debug:
	go build -o bin/manager-debug -trimpath -gcflags=all="-N -l" -ldflags " \
		-X $(REPO_URL)/$(GO_MOD_MAJOR_VERSION)/internal/manager/metadata.Release=$(TAG) \
		-X $(REPO_URL)/$(GO_MOD_MAJOR_VERSION)/internal/manager/metadata.Commit=$(COMMIT) \
		-X $(REPO_URL)/$(GO_MOD_MAJOR_VERSION)/internal/manager/metadata.Repo=$(REPO_INFO)" ${MAIN}

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: lint
lint: verify.tidy golangci-lint staticcheck looppointer

.PHONY: golangci-lint
golangci-lint: golangci-lint.download
	$(GOLANGCI_LINT) run --verbose --config $(PROJECT_DIR)/.golangci.yaml $(GOLANGCI_LINT_FLAGS)

.PHONY: staticcheck
staticcheck: staticcheck.download
	# Workaround for staticcheck not supporting nolint directives, see: https://github.com/dominikh/go-tools/issues/822.
	go list ./... | \
		grep -F -e internal/konnect/controlplanes -v | \
		xargs $(STATICCHECK) -tags envtest,e2e_tests,integration_tests,istio_tests,conformance_tests -f stylish

looppointer: looppointer.download
	$(LOOPPOINTER) -v ./internal/... ./test/...

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

CRD_GEN_PATHS ?= ./pkg/apis/configuration/...
CRD_INCUBATOR_GEN_PATHS ?= ./pkg/apis/incubator/...
CRD_OPTIONS ?= "+crd:allowDangerousTypes=true"

.PHONY: manifests
manifests: manifests.crds manifests.rbac manifests.webhook manifests.single

.PHONY: manifests.crds
manifests.crds: controller-gen ## Generate WebhookConfiguration and CustomResourceDefinition objects.
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=kong-ingress webhook paths="$(CRD_INCUBATOR_GEN_PATHS)" output:crd:artifacts:config=config/crd/incubator
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=kong-ingress webhook paths="$(CRD_GEN_PATHS)" output:crd:artifacts:config=config/crd/bases

.PHONY: manifests.rbac ## Generate ClusterRole objects.
manifests.rbac: controller-gen
	$(CONTROLLER_GEN) rbac:roleName=kong-ingress paths="./internal/controllers/configuration/" paths="./controllers/license/"
	$(CONTROLLER_GEN) rbac:roleName=kong-ingress-gateway paths="./internal/controllers/gateway/" output:rbac:artifacts:config=config/rbac/gateway
	$(CONTROLLER_GEN) rbac:roleName=kong-ingress-crds paths="./internal/controllers/crds/" output:rbac:artifacts:config=config/rbac/crds

.PHONY: manifests.webhook
manifests.webhook: controller-gen ## Generate ValidatingWebhookConfiguration.
	$(CONTROLLER_GEN) webhook paths="./internal/admission/..." output:webhook:artifacts:config=config/webhook

.PHONY: manifests.single
manifests.single: kustomize ## Compose single-file deployment manifests from building blocks
	./scripts/build-single-manifests.sh $(KUSTOMIZE)

# ------------------------------------------------------------------------------
# Build - Generators
# ------------------------------------------------------------------------------

.PHONY: generate
generate: generate.controllers generate.clientsets generate.gateway-api-consts generate.docs fmt

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
		--input-base $(REPO_URL)/$(GO_MOD_MAJOR_VERSION)/pkg/apis/ \
		--input configuration/v1,configuration/v1beta1,configuration/v1alpha1,incubator/v1alpha1 \
		--input-dirs $(REPO_URL)/pkg/apis/configuration/v1alpha1/,$(REPO_URL)/pkg/apis/configuration/v1beta1/,$(REPO_URL)/pkg/apis/configuration/v1/,$(REPO_URL)/pkg/apis/incubator/v1alpha1 \
		--output-base pkg/ \
		--output-package $(REPO_URL)/$(GO_MOD_MAJOR_VERSION)/pkg/ \
		--trim-path-prefix pkg/$(REPO_URL)/$(GO_MOD_MAJOR_VERSION)/

.PHONY: generate.docs
generate.docs: generate.apidocs generate.cli-arguments-docs

.PHONY: generate.apidocs
generate.apidocs: crd-ref-docs
	./scripts/apidocs-gen/generate.sh $(CRD_REF_DOCS)

.PHONY: generate.cli-arguments
generate.cli-arguments-docs:
	go run ./scripts/cli-arguments-docs-gen/main.go > ./docs/cli-arguments.md

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
KONG_CLUSTER_VERSION ?= v1.28.0
JUNIT_REPORT ?= /dev/null

.PHONY: bench
bench:
	@go test -count 1 -bench=. -benchmem -run=^$$ ./internal/...

.PHONY: test.all
test.all: test.unit test.envtest test.integration test.conformance

.PHONY: test.conformance
test.conformance: _check.container.environment go-junit-report
	@TEST_DATABASE_MODE="off" \
		TEST_KONG_HELM_CHART_VERSION="$(TEST_KONG_HELM_CHART_VERSION)" \
		GOFLAGS="-tags=conformance_tests" \
		go test \
		-ldflags " \
		-X $(REPO_URL)/$(GO_MOD_MAJOR_VERSION)/internal/manager/metadata.ProjectURL=$(REPO_URL) \
		-X $(REPO_URL)/$(GO_MOD_MAJOR_VERSION)/internal/manager/metadata.Release=$(TAG) \
		-X $(REPO_URL)/$(GO_MOD_MAJOR_VERSION)/internal/manager/metadata.Repo=$(REPO_INFO)" \
		-v \
		-race $(GOTESTFLAGS) \
		-timeout $(INTEGRATION_TEST_TIMEOUT) \
		-parallel $(NCPU) \
		./test/conformance | \
	$(GOJUNIT_REPORT) -iocopy -out $(JUNIT_REPORT) -parser gotest

.PHONY: test.integration
test.integration: test.integration.dbless test.integration.postgres

.PHONY: test.integration.enterprise
test.integration.enterprise: test.integration.enterprise.postgres test.integration.enterprise.dbless

.PHONY: _test.unit
.ONESHELL: _test.unit
_test.unit: gotestsum
	GOTESTSUM_FORMAT=$(GOTESTSUM_FORMAT) \
		$(GOTESTSUM) -- \
		-race $(GOTESTFLAGS) \
		-tags envtest \
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

.PHONY: test.golden.update
test.golden.update:
	@go test -v -run TestTranslator_GoldenTests ./internal/dataplane/translator -update


.PHONY: use-setup-envtest
use-setup-envtest:
	$(SETUP_ENVTEST) use

ENVTEST_TIMEOUT ?= 5m

.PHONY: _test.envtest
.ONESHELL: _test.envtest
_test.envtest: gotestsum setup-envtest use-setup-envtest
	KUBEBUILDER_ASSETS="$(shell $(SETUP_ENVTEST) use -p path)" \
		GOTESTSUM_FORMAT=$(GOTESTSUM_FORMAT) \
		$(GOTESTSUM) \
		--hide-summary output \
		-- \
		-race $(GOTESTFLAGS) \
		-tags envtest \
		-covermode=atomic \
		-timeout $(ENVTEST_TIMEOUT) \
		-coverpkg=$(PKG_LIST) \
		-coverprofile=coverage.envtest.out \
		./test/envtest/...

.PHONY: test.envtest
test.envtest:
	$(MAKE) _test.envtest GOTESTSUM_FORMAT=standard-verbose

.PHONY: test.envtest.pretty
test.envtest.pretty:
	$(MAKE) _test.envtest GOTESTSUM_FORMAT=testname

.PHONY: _check.container.environment
_check.container.environment:
	@./scripts/check-container-environment.sh

TEST_KONG_HELM_CHART_VERSION ?= $(shell yq -ojson -r '.integration.helm.kong' < .github/test_dependencies.yaml)

# Integration tests don't use gotestsum because there's a data race issue
# when go toolchain is writing to os.Stderr which is being read in go-kong
# https://github.com/Kong/go-kong/blob/c71247b5c8aae2/kong/client.go#L182
# which in turn produces a data race becuase gotestsum needs go test invoked with
# -json which enables the problematic branch in go toolchain (that writes to os.Stderr).
#
# Related issue: https://github.com/Kong/kubernetes-ingress-controller/issues/3754
.PHONY: _test.integration
_test.integration: _check.container.environment go-junit-report
	KONG_CLUSTER_VERSION="$(KONG_CLUSTER_VERSION)" \
		TEST_KONG_HELM_CHART_VERSION="$(TEST_KONG_HELM_CHART_VERSION)" \
		TEST_DATABASE_MODE="$(DBMODE)" \
		GOFLAGS="-tags=$(GOTAGS)" \
		KONG_CONTROLLER_FEATURE_GATES="$(KONG_CONTROLLER_FEATURE_GATES)" \
		go test $(GOTESTFLAGS) \
		-v \
		-timeout $(INTEGRATION_TEST_TIMEOUT) \
		-parallel $(NCPU) \
		-race \
		-covermode=atomic \
		-coverpkg=$(PKG_LIST) \
		-coverprofile=$(COVERAGE_OUT) \
		./test/integration | \
	$(GOJUNIT_REPORT) -iocopy -out $(JUNIT_REPORT) -parser gotest

.PHONY: _test.integration.isolated
_test.integration.isolated: _check.container.environment go-junit-report
	KONG_CLUSTER_VERSION="$(KONG_CLUSTER_VERSION)" \
		TEST_KONG_HELM_CHART_VERSION="$(TEST_KONG_HELM_CHART_VERSION)" \
		TEST_DATABASE_MODE="$(DBMODE)" \
		GOFLAGS="-tags=$(GOTAGS)" \
		KONG_CONTROLLER_FEATURE_GATES="$(KONG_CONTROLLER_FEATURE_GATES)" \
		go test $(GOTESTFLAGS) \
		-v \
		-timeout $(INTEGRATION_TEST_TIMEOUT) \
		-parallel $(NCPU) \
		-race \
		-covermode=atomic \
		-coverpkg=$(PKG_LIST) \
		-coverprofile=$(COVERAGE_OUT) \
		./test/integration/isolated -args --parallel $(E2E_FRAMEWORK_FLAGS) | \
	$(GOJUNIT_REPORT) -iocopy -out $(JUNIT_REPORT) -parser gotest

.PHONY: test.integration.isolated.dbless
test.integration.isolated.dbless:
	@$(MAKE) _test.integration.isolated \
		GOTAGS="integration_tests" \
		DBMODE=off \
		COVERAGE_OUT=coverage.dbless.out

.PHONY: test.integration.isolated.postgres
test.integration.isolated.postgres:
	@$(MAKE) _test.integration.isolated \
		GOTAGS="integration_tests" \
		DBMODE=postgres \
		COVERAGE_OUT=coverage.postgres.out

.PHONY: test.integration.dbless
test.integration.dbless:
	@$(MAKE) _test.integration \
		GOTAGS="integration_tests" \
		DBMODE=off \
		COVERAGE_OUT=coverage.dbless.out

.PHONY: test.integration.postgres
test.integration.postgres:
	@$(MAKE) _test.integration \
		GOTAGS="integration_tests" \
		DBMODE=postgres \
		COVERAGE_OUT=coverage.postgres.out

.PHONY: test.integration.enterprise.postgres
test.integration.enterprise.postgres:
	@TEST_KONG_ENTERPRISE="true" \
		GOTAGS="integration_tests" \
		$(MAKE) _test.integration \
		DBMODE=postgres \
		COVERAGE_OUT=coverage.enterprisepostgres.out

.PHONY: test.integration.enterprise.dbless
test.integration.enterprise.dbless:
	@TEST_KONG_ENTERPRISE="true" \
		GOTAGS="integration_tests" \
		$(MAKE) _test.integration \
		DBMODE=off \
		COVERAGE_OUT=coverage.enterprisedbless.out

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

.PHONY: test.performance
test.performance: gotestsum
	GOFLAGS="-tags=performance_tests" \
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

.PHONY: test.kongintegration
test.kongintegration:
	$(MAKE) _test.kongintegration GOTESTSUM_FORMAT=standard-verbose

.PHONY: test.kongintegration.pretty
test.kongintegration.pretty:
	$(MAKE) _test.kongintegration GOTESTSUM_FORMAT=testname

.PHONY: _test.kongintegration
_test.kongintegration: gotestsum go-junit-report
	# Disable testcontainer's reaper (Ryuk). It's needed because Ryuk requires
	# privileged mode to run, which is not desired and could cause issues in CI.
	TESTCONTAINERS_RYUK_DISABLED="true" \
	GOTESTSUM_FORMAT=$(GOTESTSUM_FORMAT) \
	$(GOTESTSUM) -- $(GOTESTFLAGS) \
		-race \
		-parallel $(NCPU) \
		-coverpkg=$(PKG_LIST) \
		-coverprofile=coverage.kongintegration.out \
		./test/kongintegration | \
	$(GOJUNIT_REPORT) -iocopy -out $(JUNIT_REPORT) -parser gotest

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

SKAFFOLD_DEBUG_PROFILE ?= debug_multi_gw

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
#   * `konnect.env` with CONTROLLER_KONNECT_CONTROL_PLANE_ID env variable set
#     to the UUID of a Control Plane you have created in Konnect.
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

# NOTE: We're using the --keep-running-on-failure=true to allow deployments like
# postgres multigateway to eventually stabilize.
# TODO: verify if --keep-running-on-failure=true is still needed when
# https://github.com/Kong/kubernetes-ingress-controller/issues/5116 is implemented.
.PHONY: _skaffold
_skaffold: skaffold
	$(SKAFFOLD) $(CMD) --keep-running-on-failure=true --port-forward=pods --profile=$(SKAFFOLD_PROFILE) $(SKAFFOLD_FLAGS)

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

.PHONY: generate.gateway-api-consts
generate.gateway-api-consts:
	GATEWAY_API_VERSION=$(GATEWAY_API_VERSION) \
		CRDS_STANDARD_URL=$(shell GATEWAY_API_RELEASE_CHANNEL="" $(MAKE) print-gateway-api-crds-url) \
		CRDS_EXPERIMENTAL_URL=$(shell GATEWAY_API_RELEASE_CHANNEL="experimental" $(MAKE) print-gateway-api-crds-url) \
		RAW_REPO_URL=$(shell $(MAKE) print-gateway-api-raw-repo-url) \
		INPUT=$(shell pwd)/test/internal/cmd/generate-gateway-api-consts/gateway_consts.tmpl \
		OUTPUT=$(shell pwd)/test/consts/zz_generated_gateway.go \
		go generate -tags=generate_gateway_api_consts ./test/internal/cmd/generate-gateway-api-consts

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
