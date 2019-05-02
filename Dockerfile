FROM golang:1.12 AS build
WORKDIR /kong-ingress-controller
COPY go.mod ./
COPY go.sum ./
RUN go mod download
ADD . .
ARG TAG
ARG REPO_INFO
ARG COMMIT
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o kong-ingress-controller -ldflags "-s -w -X main.RELEASE=$TAG -X main.COMMIT=$COMMIT -X main.REPO=$REPO_INFO" ./cli/ingress-controller

FROM alpine:3.9
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=build /kong-ingress-controller .
ENTRYPOINT ["./kong-ingress-controller"]
