# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app

# Copy go.mod only
COPY go.mod ./

# Generate go.sum and download dependencies
RUN go mod tidy

# Copy the rest of the source code
COPY . .

# Build the binary (normal build, not forcing static)
RUN go build -o evmon ./cmd/main.go

# Final stage
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/evmon .
EXPOSE 8080
CMD ["./evmon"]