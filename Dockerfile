# ===============================================
# Multi-stage build for Go Scheduler Application
# ===============================================

# ------------------- Builder Stage -------------------
FROM golang:1.25.5-alpine AS builder

WORKDIR /app

# Cache dependency layers
COPY go.mod go.sum ./
RUN go mod download

# Copy application source code
COPY . .

# Compile the binary in accordance with the local console build command
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-s -w" -o todo-app ./cmd/.

# ------------------- Final Stage -------------------
FROM alpine:3.20

# Update security certificates and provision timezone data properties
RUN apk --no-cache update && apk --no-cache add ca-certificates tzdata && rm -rf /var/cache/apk/*

# Create a secure, non-privileged system application user
RUN adduser -D -g '' appuser

WORKDIR /app

# Copy the compiled binary execution context from the builder stage
COPY --from=builder /app/todo-app .

# Copy user-facing web directory assets
COPY --from=builder /app/web/ ./web/

# Assign application directory ownership privileges to the unprivileged user
RUN chown -R appuser:appuser /app

USER appuser

# Expose the designated network interface port
EXPOSE 7540

# Define container internal health status verification scripts
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 CMD wget --no-verbose --tries=1 --spider http://localhost:7540/ || exit 1

# Execute the application binary target
CMD ["./todo-app"]