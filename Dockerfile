FROM golang:1.27.1-alpine AS builder

WORKDIR /app

# Cache dependencies first (improves build times)
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Compile static binary with symbol stripping (-s -w reduces size by ~30%)
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /app/server ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server /app/server

# Non-root user (UID 65532) for defense-in-depth
USER nonroot:nonroot

EXPOSE 8080

ENTRYPOINT ["/app/server"]
