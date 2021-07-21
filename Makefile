TAG?=1.3.1
RG_TAG?=2.0.0-alpha.2
REGISTRY?=kong
REPO_INFO=$(shell git config --get remote.origin.url)
IMGNAME?=kubernetes-ingress-controller
IMAGE = $(REGISTRY)/$(IMGNAME)
# only for dev
DB?=false
RUN_VERSION?=20
KUBE_VERSION?=v1.20.7

PKG_LIST := ./...

ifndef COMMIT
  COMMIT := $(shell git rev-parse --short HEAD)
endif

export GO111MODULE=on

.PHONY: test-all
test-all: lint test

.PHONY: test
test:
	go test -race -covermode=atomic -coverpkg=$(PKG_LIST) $(PKG_LIST)

.PHONY: coverage
coverage:
	go test -race -v -count=1 -covermode=atomic -coverpkg=$(PKG_LIST) -coverprofile=coverage.out.tmp $(PKG_LIST)
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

.PHONY: container-alpine
container-alpine:
	docker build \
	--build-arg TAG=${TAG} --build-arg COMMIT=${COMMIT} \
	--build-arg REPO_INFO=${REPO_INFO} \
	--target alpine \
	-t ${IMAGE}:${TAG}-alpine .

.PHONY: container-redhat
container-redhat:
	docker build \
	--build-arg TAG=${TAG} --build-arg COMMIT=${COMMIT} \
	--build-arg REPO_INFO=${REPO_INFO} \
	--target redhat \
	-t ${IMAGE}:${TAG}-redhat .

.PHONY: container
container: container-alpine
	docker tag "${IMAGE}:${TAG}-alpine" "${IMAGE}:${TAG}"

.PHONY: railgun-container
railgun-container:
	docker build \
    -f Dockerfile.railgun \
    --build-arg TAG=${RG_TAG} --build-arg COMMIT=${COMMIT} \
    --build-arg REPO_INFO=${REPO_INFO} \
    -t ${IMAGE}:${RG_TAG} .

.PHONY: run
run:
	./hack/dev/start.sh ${DB} ${RUN_VERSION}


.PHONY: integration-test
integration-test: container
	KIC_IMAGE="${IMAGE}:${TAG}" KUBE_VERSION=${KUBE_VERSION} ./test/integration/test.sh
