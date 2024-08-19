### Standard binary
# Build the manager binary
FROM --platform=$BUILDPLATFORM golang:1.23.0 AS builder

ARG GOPATH
ARG GOCACHE

ARG TARGETPLATFORM
ARG TARGETOS
ARG TARGETARCH

RUN printf "Building for TARGETPLATFORM=${TARGETPLATFORM}" \
    && printf ", TARGETARCH=${TARGETARCH}" \
    && printf ", TARGETOS=${TARGETOS}" \
    && printf ", TARGETVARIANT=${TARGETVARIANT} \n" \
    && printf "With 'uname -s': $(uname -s) and 'uname -m': $(uname -m)"

WORKDIR /workspace

# Use cache mounts to cache Go dependencies and bind mounts to avoid unnecessary
# layers when using COPY instructions for go.mod and go.sum.
# https://docs.docker.com/build/guide/mounts/
RUN --mount=type=cache,target=$GOPATH/pkg/mod \
    --mount=type=cache,target=$GOCACHE \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

COPY controllers/ controllers/
COPY pkg/ pkg/
COPY internal/ internal/
COPY Makefile .

# Build
ARG TAG
ARG COMMIT
ARG REPO_INFO

# Use cache mounts to cache Go dependencies and bind mounts to avoid unnecessary
# layers when using COPY instructions for go.mod and go.sum.
# https://docs.docker.com/build/guide/mounts/
RUN --mount=type=cache,target=$GOPATH/pkg/mod \
    --mount=type=cache,target=$GOCACHE \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    CGO_ENABLED=0 GOOS=linux GOARCH="${TARGETARCH}" GO111MODULE=on \
    make _build

### FIPS 140-2 binary
# Build the manager binary
# https://github.com/golang/go/tree/dev.boringcrypto/misc/boring#building-from-docker
FROM us-docker.pkg.dev/google.com/api-project-999119582588/go-boringcrypto/golang:1.18.10b7 AS builder-fips

ARG TARGETPLATFORM
ARG TARGETOS
ARG TARGETARCH

RUN printf "Building for TARGETPLATFORM=${TARGETPLATFORM}" \
    && printf ", TARGETARCH=${TARGETARCH}" \
    && printf ", TARGETOS=${TARGETOS}" \
    && printf ", TARGETVARIANT=${TARGETVARIANT} \n" \
    && printf "With 'uname -s': $(uname -s) and 'uname -m': $(uname -m)"

WORKDIR /workspace

COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

COPY pkg/ pkg/
COPY internal/ internal/

# Build
ARG TAG
ARG COMMIT
ARG REPO_INFO

RUN CGO_ENABLED=0 GOOS=linux GOARCH="${TARGETARCH}" GO111MODULE=on make _build.fips

### distroless FIPS 140-2
FROM gcr.io/distroless/static:nonroot AS distroless-fips
WORKDIR /
COPY --from=builder-fips /workspace/manager .
USER 1000:1000

ENTRYPOINT ["/manager"]

### Distroless/default
# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot AS distroless
ARG TAG
ARG TARGETPLATFORM
ARG TARGETOS
ARG TARGETARCH

LABEL name="Kong Ingress Controller" \
      vendor="Kong" \
      version="$TAG" \
      release="1" \
      url="https://github.com/Kong/kubernetes-ingress-controller" \
      summary="Kong for Kubernetes Ingress" \
      description="Use Kong for Kubernetes Ingress. Configure plugins, health checking, load balancing and more in Kong for Kubernetes Services, all using Custom Resource Definitions (CRDs) and Kubernetes-native tooling."

WORKDIR /
COPY --from=builder /workspace/bin/manager .
USER 1000:1000

ENTRYPOINT ["/manager"]
