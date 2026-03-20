FROM golang:alpine AS base

RUN apk add --no-cache git ca-certificates curl

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Development stage
FROM base AS dev
RUN go install github.com/air-verse/air@latest
CMD ["air", "-c", ".air.toml"]

FROM base AS builder
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o /app/server ./cmd/api

# Production stage
FROM alpine:3.19 AS prod

RUN apk add --no-cache ca-certificates tzdata curl

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/migrations ./migrations

USER appuser

EXPOSE 8080

ENTRYPOINT ["./server"]
