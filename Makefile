REGISTRY?=kong-docker-kubernetes-ingress-controller.bintray.io
TAG?=0.4.0
REPO_INFO=$(shell git config --get remote.origin.url)
IMGNAME = kong-ingress-controller
IMAGE = $(REGISTRY)/$(IMGNAME)

ifndef COMMIT
  COMMIT := $(shell git rev-parse --short HEAD)
endif

export GO111MODULE=on

.PHONY: test-all
test-all: lint fmt vet test

.PHONY: test
test:
	go test -race ./...

.PHONY: vet
vet:
	go vet ./...

lint:
	golint -set_exit_status ./...

.PHONY: build
build:
	CGO_ENABLED=0 go build -o kong-ingress-controller ./cli/ingress-controller

.PHONY: fmt
fmt:
	bash -c "diff -u <(echo -n) <(gofmt -d -l -e -s .)"

.PHONY: verify-codegen
verify-codegen:
	./hack/verify-codegen.sh

.PHONY: update-codegen
update-codegen:
	./hack/update-codegen.sh

.PHONY: container
container:
	docker build \
    --build-arg TAG=${TAG} --build-arg COMMIT=${COMMIT} \
    --build-arg REPO_INFO=${REPO_INFO} \
    -t ${IMAGE}:${TAG} .
