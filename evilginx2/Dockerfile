FROM golang:1.18-alpine AS builder

WORKDIR /build

# Install build dependencies
RUN apk add --no-cache git make

# Copy source code
COPY . .

# Build frontend
RUN apk add --no-cache nodejs npm
RUN chmod +x build_frontend.sh
RUN ./build_frontend.sh

# Build Go application
RUN go build -o evilginx2 main.go

# Final stage
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy built application
COPY --from=builder /build/evilginx2 /app/evilginx2

# Copy frontend dist
COPY --from=builder /build/frontend/dist /app/frontend/dist

# Create directories
RUN mkdir -p /app/phishlets /app/redirectors

# Make executable
RUN chmod +x /app/evilginx2

# Expose ports
EXPOSE 8080 443 53/udp

ENTRYPOINT ["/app/evilginx2"]
