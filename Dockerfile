# Build the manager binary
ARG BUILDPLATFORM
ARG TARGETPLATFORM
FROM --platform=${BUILDPLATFORM} docker.io/library/golang:1.27.1 AS builder


ARG VERSION=dev
ARG COMMIT
ARG BUILD_TIME
ARG PROJECT=github.com/americancode/adcs-issuer
ARG TARGETOS
ARG TARGETARCH

WORKDIR /workspace


# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum


# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download


# Copy the go source
COPY main.go main.go
COPY api/ api/
COPY controllers/ controllers/
COPY issuers/ issuers/
COPY adcs/ adcs/
COPY healthcheck/ healthcheck/
COPY version/ version/
# Build
#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o manager main.go


RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} GO111MODULE=on go build \
		-ldflags "-s -w -X ${PROJECT}/version.Release=${VERSION} \
		-X ${PROJECT}/version.Commit=${COMMIT} -X ${PROJECT}/version.BuildTime=${BUILD_TIME}" \
		-o manager main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM --platform=${TARGETPLATFORM} gcr.io/distroless/static:nonroot
WORKDIR /

COPY --from=builder /workspace/manager .
USER nonroot:nonroot

ENTRYPOINT ["/manager"]
