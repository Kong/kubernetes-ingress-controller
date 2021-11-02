# Build the manager binary
FROM golang:1.17 as builder

WORKDIR /workspace

COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

COPY pkg/ pkg/
COPY pkg/ pkg/
COPY internal/ internal/

# Build
ARG TAG
ARG COMMIT
ARG REPO_INFO
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o manager -ldflags "-s -w -X github.com/kong/kubernetes-ingress-controller/v2/internal/metadata.Release=$TAG -X github.com/kong/kubernetes-ingress-controller/v2/internal/metadata.Commit=$COMMIT -X github.com/kong/kubernetes-ingress-controller/v2/internal/metadata.Repo=$REPO_INFO" ./internal/cmd/main.go

# Build a manager binary with debug symbols and download Delve
FROM builder as builder-delve

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o manager-debug -gcflags=all="-N -l" -ldflags "-X github.com/kong/kubernetes-ingress-controller/v2/internal/metadata.Release=$TAG -X github.com/kong/kubernetes-ingress-controller/v2/internal/metadata.Commit=$COMMIT -X github.com/kong/kubernetes-ingress-controller/v2/internal/metadata.Repo=$REPO_INFO" ./internal/cmd/main.go

# Create an image that runs a debug build with a Delve remote server on port 2345
FROM golang:1.17 AS debug
RUN go install github.com/go-delve/delve/cmd/dlv@latest
# We want all source so Delve file location operations work
COPY --from=builder-delve /workspace/ /workspace/
USER 65532:65532

ENTRYPOINT ["/go/bin/dlv"]
CMD ["exec", "--continue", "--accept-multiclient",  "--headless", "--api-version=2", "--listen=:2345", "--log", "/workspace/manager-debug"]

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot AS distroless
WORKDIR /
COPY --from=builder /workspace/manager .
USER 65532:65532

ENTRYPOINT ["/manager"]

# Build UBI image
FROM registry.access.redhat.com/ubi8/ubi AS redhat
ARG TAG

LABEL name="Kong Ingress Controller" \
      vendor="Kong" \
      version="$TAG" \
      release="1" \
      url="https://github.com/Kong/kubernetes-ingress-controller" \
      summary="Kong for Kubernetes Ingress" \
      description="Use Kong for Kubernetes Ingress. Configure plugins, health checking, load balancing and more in Kong for Kubernetes Services, all using Custom Resource Definitions (CRDs) and Kubernetes-native tooling."

# Create the user (ID 1000) and group that will be used in the
# running container to run the process as an unprivileged user.
RUN groupadd --system kic && \
    adduser --system kic -g kic -u 1000

COPY --from=builder /workspace/manager .
COPY LICENSE /licenses/

# Perform any further action as an unprivileged user.
USER 1000

# Run the compiled binary.
ENTRYPOINT ["/manager"]
