# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o evmon ./cmd/main.go

# Final stage
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/evmon .
EXPOSE 8080
CMD ["./evmon"]