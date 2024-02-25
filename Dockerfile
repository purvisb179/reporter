# Stage 1: build the Go binary
FROM --platform=linux/amd64 golang:1.20 as builder

WORKDIR /app

# Fetch dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binary with GOARCH=amd64 to target AMD architecture
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o myapp main.go

# Stage 2: copy the Go binary to an Alpine container and update certs
# Using --platform=linux/amd64 ensures the Alpine image is suitable for AMD64 architecture
FROM --platform=linux/amd64 alpine:latest

# Update certificates in Alpine
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/myapp /app/myapp

# Copy the templates directory
COPY --from=builder /app/internal/templates /app/internal/templates

# The command to start the app
CMD ["/app/myapp", "serve"]
