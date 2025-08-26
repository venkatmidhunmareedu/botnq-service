# Start from official Go image
FROM golang:1.22 AS builder

# Set working directory inside container
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the Go binary
RUN go build -o main ./cmd

# Small final image
FROM gcr.io/distroless/base-debian12

WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/main .

# Expose the port your Go app listens on
EXPOSE 8080

# Run the app
CMD ["./main"]
