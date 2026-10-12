# Runtime base is distroless static: no shell, no libc, runs as nonroot (65532).
ARG GO_IMAGE=golang:1.27
ARG RUNTIME_IMAGE=gcr.io/distroless/static-debian12:nonroot

FROM ${GO_IMAGE} AS builder
WORKDIR /src

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags "-s -w" -o /out/rate-limiter ./cmd/server

FROM ${RUNTIME_IMAGE}
COPY --from=builder /out/rate-limiter /rate-limiter
EXPOSE 8080
USER 65532:65532
ENTRYPOINT ["/rate-limiter"]
