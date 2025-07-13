# ─────── Stage 1: Build binary ───────
FROM --platform=$BUILDPLATFORM golang:1.24-alpine3.22 AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ ./cmd/
COPY internal/ ./internal/


# Build statically for the target OS/ARCH
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o main ./cmd/main.go

# ─────── Stage 2: Minimal final image ───────
FROM scratch

COPY --from=builder /app/main /main

ENTRYPOINT ["/main"]