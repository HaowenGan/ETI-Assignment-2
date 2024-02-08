# Use a multi-stage build to keep the final image light
# Step 1: Build the Go app
FROM golang:1.22.0 as builder

# Set the working directory inside the container
WORKDIR /app

# Set environment variables for Go build
ENV CGO_ENABLED=0
ENV GOOS=linux

# Copy go mod and sum files
COPY user-service/go.mod user-service/go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY user-service .

# Build the application
RUN go build -o main .

# Step 2: Use a lightweight base image
FROM alpine:latest

# Set the working directory in the container
WORKDIR /root/

# Install ca-certificates for HTTPS requests, and any other dependencies your Go program might require
RUN apk --no-cache add ca-certificates

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

# Create a directory to serve as the 'parent' of the working directory and copy the front-end assets and HTML files into it
RUN mkdir /root/parent
COPY /front-end/css /root/parent/front-end/css/
COPY /front-end/js /root/parent/front-end/js/
COPY /*.html /root/parent/front-end/

# Expose port 5000 for the application
EXPOSE 5000

# Command to run the binary
CMD ["./main"]