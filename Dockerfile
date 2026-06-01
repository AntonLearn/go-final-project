# ===============================================
# Multi-stage build for Go Scheduler Application
# ===============================================

# ------------------- Builder Stage -------------------
FROM golang:1.25.0-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -a \
    -installsuffix cgo \
    -ldflags="-s -w" \
    -o todo-app ./cmd/.

# ------------------- Final Stage -------------------
FROM alpine:latest

RUN apk --no-cache update && \
    apk --no-cache add ca-certificates tzdata && \
    rm -rf /var/cache/apk/*

RUN adduser -D -g '' appuser

WORKDIR /app

# Copy binary
COPY --from=builder /app/todo-app .

# Copy web directory
COPY --from=builder /app/web/ ./web/

# Copy default database (if exists) using RUN to avoid linter issues
RUN cp -f /app/scheduler.db ./scheduler.db 2>/dev/null || true

# Set correct permissions
RUN chown -R appuser:appuser /app

USER appuser

EXPOSE 7540

ENV TODO_PORT=7540
ENV TODO_DBFILE=scheduler.db
ENV TODO_PASSWORD=""

HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:7540/health || exit 1

CMD ["./todo-app"]