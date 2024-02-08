# Use a multi-stage build to keep the final image light
# Step 1: Build the feedback-service Go app
FROM golang:1.22.0 as feedback-builder

# Set the working directory inside the container
WORKDIR /feedback-app

# Set environment variables for Go build
ENV CGO_ENABLED=0
ENV GOOS=linux

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the application
RUN go build -o feedback-service .

# Step 2: Use a lightweight base image
FROM alpine:latest

# Set the working directory in the container
WORKDIR /root/

# Install ca-certificates for HTTPS requests, and any other dependencies your Go program might require
RUN apk --no-cache add ca-certificates

# Copy the built binary from the builder stage
COPY --from=feedback-builder /feedback-app/feedback-service .

# Expose the port that feedback-service runs on
EXPOSE 5001

# Command to run the binary
CMD ["./feedback-service"]