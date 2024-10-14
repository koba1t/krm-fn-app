# Build the manager binary
FROM public.ecr.aws/docker/library/golang:1.22.7-bullseye as builder

WORKDIR /workspace

# Download module
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

# Copy the go source
COPY main.go main.go
# COPY pkg/ pkg/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o main main.go

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/main .
USER nonroot:nonroot

ENTRYPOINT ["/main"]
