# Build the manager binary
FROM golang:1.16 as builder

WORKDIR /workspace

COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

COPY pkg/ pkg/
COPY railgun/pkg/ railgun/pkg/
COPY railgun/apis/ railgun/apis/
COPY railgun/internal/ railgun/internal/
COPY railgun/main.go railgun/main.go

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o manager -ldflags "-s -w -X manager.Release=$TAG -X manager.Commit=$COMMIT -X manager.Repo=$REPO_INFO" ./railgun/main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/manager .
USER 65532:65532

ENTRYPOINT ["/manager"]
