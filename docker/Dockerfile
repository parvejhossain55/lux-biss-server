# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy dependency files and download
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application as a static binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/api/main.go

# Run stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Create a non-privileged user
RUN adduser -D -g '' appuser

WORKDIR /app

# Copy the binary and config from the builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

# Use the non-privileged user
USER appuser

EXPOSE 8080

CMD ["./main"]
