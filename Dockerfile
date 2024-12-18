# Use the official Go image for building
FROM golang:1.23.2 AS builder

# Set environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set the working directory for the builder stage
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download Go dependencies
RUN go mod download

# Copy the entire project directory into the container
COPY . .

# Navigate to the cmd directory and build the binary
WORKDIR /app/cmd
RUN go build -o /app/main .

# Use a minimal image for running the application
FROM alpine:latest

# Set the working directory in the final image
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Expose the application port
EXPOSE 8080

# Command to run the application
CMD ["./main"]
