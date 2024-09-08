# Use an official Go image to build the Go binary
FROM golang:1.21.6 as builder

# Set the current working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum .env ./

# Download all the dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go binary with the Makefile
RUN make da-server

EXPOSE 8080
EXPOSE 433

# Set default values for environment variables
ENV ADDR=0.0.0.0
ENV PORT=8080
ENV AVAIL_RPC=http://localhost:9933
ENV AVAIL_SEED=""
ENV AVAIL_APPID=0
ENV AVAIL_TIMEOUT=100

EXPOSE ${PORT}

# Print environment variables and run the application
CMD echo "ADDR: ${ADDR}" && \
    echo "PORT: ${PORT}" && \
    echo "AVAIL_RPC: ${AVAIL_RPC}" && \
    echo "AVAIL_SEED: ${AVAIL_SEED}" && \
    echo "AVAIL_APPID: ${AVAIL_APPID}" && \
    echo "AVAIL_TIMEOUT: ${AVAIL_TIMEOUT}" && \
    ./bin/avail-da-server --addr=$ADDR --port=$PORT --avail.rpc="$AVAIL_RPC" --avail.seed="$AVAIL_SEED" --avail.appid=$AVAIL_APPID --avail.timeout=$AVAIL_TIMEOUT
