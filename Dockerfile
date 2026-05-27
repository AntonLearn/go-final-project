FROM golang:1.25.0-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o todo-app .

FROM alpine:latest
RUN apk --no-cache update && apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/todo-app .
COPY --from=builder /app/web ./web

EXPOSE 7540

ENV TODO_PORT=7540
ENV TODO_DBFILE=scheduler.db
ENV TODO_PASSWORD=password

CMD ["./todo-app"]