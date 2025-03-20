# Use the official Golang image as the base image for building
FROM golang:1.23-alpine AS builder

# Install necessary build tools for CGO
RUN apk add --no-cache gcc musl-dev

# Set the working directory inside the container
WORKDIR /src

# Copy go.mod and go.sum, download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . ./

# Build the application
RUN go build -o main ./src/main.go

# Use a minimal base image for the final image
FROM alpine:latest

# Install necessary certificates and runtime dependencies
RUN apk --no-cache add ca-certificates libc6-compat

# Set the working directory
WORKDIR /src

# Copy the built binary from the builder stage
COPY --from=builder /src/main ./main

# Copy configuration files if needed
COPY ./src/config/environment.yaml ./config/environment.yaml

# Expose the port the application runs on
EXPOSE 3000

# Run the application
CMD ["./main"]

# File name: Dockerfile.Prod
