FROM golang:1.16 AS build
WORKDIR /kong-ingress-controller
COPY go.mod ./
COPY go.sum ./
RUN go mod download
ADD . .
ARG TAG
ARG REPO_INFO
ARG COMMIT
RUN CGO_ENABLED=0 GOOS=linux go build -o kong-ingress-controller -ldflags "-s -w -X main.RELEASE=$TAG -X main.COMMIT=$COMMIT -X main.REPO=$REPO_INFO" ./cli/ingress-controller

# Build Alpine image
FROM alpine:3.11 AS alpine
ARG TAG
LABEL name="Kong Ingress Controller" \
      vendor="Kong" \
      version="$TAG" \
      release="1" \
      url="https://github.com/Kong/kubernetes-ingress-controller" \
      summary="Kong for Kubernetes Ingress" \
      description="Use Kong for Kubernetes Ingress. Configure plugins, health checking, load balancing and more in Kong for Kubernetes Services, all using Custom Resource Definitions (CRDs) and Kubernetes-native tooling."

RUN apk --no-cache add ca-certificates

# Create the user (ID 1000) and group that will be used in the
# running container to run the process as an unprivileged user.
RUN addgroup -S kic && \
    adduser -S kic -G kic -u 1000

# Import the compiled executable from the second stage.
COPY --from=build /kong-ingress-controller/kong-ingress-controller /bin
# Only for backwards compatibility
COPY --from=build /kong-ingress-controller/kong-ingress-controller .

# Perform any further action as an unprivileged user.
USER 1000

# Run the compiled binary.
ENTRYPOINT ["/kong-ingress-controller"]

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

# Import the compiled executable from the second stage.
COPY --from=build /kong-ingress-controller/kong-ingress-controller /bin
# Only for backwards compatibility
COPY --from=build /kong-ingress-controller/kong-ingress-controller .

# Perform any further action as an unprivileged user.
USER 1000

# Run the compiled binary.
ENTRYPOINT ["/kong-ingress-controller"]
