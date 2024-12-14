# Use the official Golang image as the base image
FROM golang:1.23 AS builder

# Set the working directory inside the container
WORKDIR /src

# Copy the Go modules manifests
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go srclication
RUN go build -o main ./src/main.go

# Use a minimal base image for the final build
FROM debian:bookworm-slim

# Install CA certificates
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*


# Set the working directory inside the container
WORKDIR /src

# Copy the built binary from the builder stage
COPY --from=builder /src/main .

# Copy any necessary static files or configuration (optional)
COPY ./src/config/environment.yaml ./config/

# Expose the port the src runs on (adjust if necessary)
EXPOSE 3000

# Command to run the executable
CMD ["./main"]
