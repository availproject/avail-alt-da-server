# Use an official Go image to build the Go binary
FROM golang:1.21.6 as builder

# Set the current working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download all the dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go binary with the Makefile
RUN make da-server

# Copy the .env file
COPY .env .

# Set default values for environment variables
ENV ADDR=0.0.0.1
ENV PORT=3100
ENV AVAIL_RPC=http://localhost:9933
ENV AVAIL_SEED=""
ENV AVAIL_APPID=0
ENV AVAIL_TIMEOUT=100

ENTRYPOINT ["/bin/bash"]
