# Build stage - builds the CLI for multiple platforms and the server
FROM golang:1.24.1-alpine AS builder

WORKDIR /build

# Install git for dependency management
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the go-find CLI for multiple platforms
RUN mkdir -p /binaries/linux/amd64 /binaries/linux/arm64 /binaries/linux/arm \
    /binaries/darwin/amd64 /binaries/darwin/arm64 \
    /binaries/windows/amd64 /binaries/windows/arm64

# Linux builds
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /binaries/linux/amd64/go-find .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o /binaries/linux/arm64/go-find .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -ldflags="-s -w" -o /binaries/linux/arm/go-find .

# macOS builds
RUN CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o /binaries/darwin/amd64/go-find .
RUN CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o /binaries/darwin/arm64/go-find .

# Windows builds
RUN CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o /binaries/windows/amd64/go-find.exe .
RUN CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -ldflags="-s -w" -o /binaries/windows/arm64/go-find.exe .

# Build the HTTP server
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /server ./server

# Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the server binary
COPY --from=builder /server .

# Copy all platform binaries
COPY --from=builder /binaries ./binaries

# Create a non-root user for security
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser && \
    chown -R appuser:appuser /app

USER appuser

# Expose port (Easypanel will set PORT env var)
EXPOSE 8080

# Run the HTTP server
CMD ["/app/server"]
