# Build the manager binary
FROM europe-docker.pkg.dev/kyma-project/prod/external/golang:1.22.2-alpine3.19 as builder

RUN apk update && apk upgrade && apk add --no-cache ca-certificates
RUN update-ca-certificates

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY main.go main.go
COPY handlers/ handlers/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o manager main.go

FROM scratch
WORKDIR /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /workspace/manager .

USER 65532:65532

ENTRYPOINT ["/manager"]