## Build stage
FROM golang:1.25-alpine3.23 AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

# Build fully static binary with CGO disabled
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main .

## Package stage
FROM alpine:3.23.3

# Install root CA certificates for outbound HTTPS calls to bank APIs
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user for security
RUN adduser -D -g '' appuser

WORKDIR /home/appuser

COPY --from=builder /app/main .

# Transfer ownership to non-root user
RUN chown appuser:appuser main
USER appuser

EXPOSE 8080

CMD ["./main"]
