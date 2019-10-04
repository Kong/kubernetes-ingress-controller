FROM golang:1.13 AS build
WORKDIR /kong-ingress-controller
COPY go.mod ./
COPY go.sum ./
RUN go mod download
ADD . .
ARG TAG
ARG REPO_INFO
ARG COMMIT
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o kong-ingress-controller -ldflags "-s -w -X main.RELEASE=$TAG -X main.COMMIT=$COMMIT -X main.REPO=$REPO_INFO" ./cli/ingress-controller

# Final stage: the running container.
FROM alpine:3.9
RUN apk --no-cache add ca-certificates

# Create the user (ID 1000) and group that will be used in the running container to
# run the process as an unprivileged user.
RUN addgroup -S kong && \
    adduser -S kong -G kong -u 1000

# Import the compiled executable from the second stage.
COPY --from=build /kong-ingress-controller/kong-ingress-controller .

# Perform any further action as an unprivileged user.
USER 1000

# Run the compiled binary.
ENTRYPOINT ["/kong-ingress-controller"]

