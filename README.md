# Avail Optimism Alt DA Server

[![Go Version](https://img.shields.io/badge/Go-1.23.4+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

## Introduction

This is a sidecar Data Availability (DA) Server for Optimism that integrates with Avail DA for posting and retrieving data. The server acts as an alternative DA layer, providing secure and scalable data availability solutions for Optimism rollups.

### Features

- **Dual DA Support**: Compatible with both Avail DA and Turbo DA services
- **Optimism Integration**: Seamlessly integrates with Optimism rollup architecture
- **High Performance**: Optimized for fast data submission and retrieval
- **Flexible Configuration**: Extensive configuration options for different deployment scenarios
- **Docker Support**: Containerized deployment for easy scaling

## Prerequisites

- **Go**: Version 1.23.4 or higher
- **Docker**: Optional, for containerized deployment
- **Avail Account**: Funded wallet with AVAIL tokens and App ID (for Avail DA)
- **Turbo DA Account**: URL and API key (for Turbo DA, optional)

## Installation

### Clone Repository

```shell
git clone <repository-url>
cd avail-alt-da-server
```

### Build Binary

```shell
make da-server
```

The binary will be created at `./bin/avail-da-server`.

## Configuration

The server can be configured via command-line flags or environment variables.

### Command Line Flags

| Flag | Description | Default Value |
|------|-------------|---------------|
| `--addr` | Address to bind the server to | `127.0.0.1` |
| `--port` | Port to run the server on | `3100` |
| `--avail.rpc` | Avail HTTP RPC URL (REQUIRED) | - |
| `--avail.seed` | Avail seed phrase for authentication | - |
| `--avail.appid` | Avail App ID for your application | `0` |
| `--turboda` | Enable Turbo DA service | `false` |
| `--turboda.url` | Turbo DA service URL | - |
| `--turboda.key` | Turbo DA API key | - |

### Environment Variables

Copy the example environment file:

```shell
cp .env.example .env
```

Edit `.env` with your configuration:

```bash
# Server Configuration
OP_PLASMA_AVAIL_DA_SERVER_ADDR=127.0.0.1
OP_PLASMA_AVAIL_DA_SERVER_PORT=3100

# Avail DA Configuration
OP_PLASMA_AVAIL_DA_SERVER_AVAIL_RPC=https://turing-rpc.avail.so/rpc
OP_PLASMA_AVAIL_DA_SERVER_AVAIL_SEED="your seed phrase here"
OP_PLASMA_AVAIL_DA_SERVER_AVAIL_APPID=1

# Turbo DA Configuration (Optional)
OP_PLASMA_AVAIL_DA_SERVER_TURBODA=false
OP_PLASMA_AVAIL_DA_SERVER_TURBODA_URL=https://turing.turbo-api.availproject.org
OP_PLASMA_AVAIL_DA_SERVER_TURBODA_KEY=your-api-key-here
```

## Usage

### Quick Start

1. **Configure Environment**:

   ```shell
   cp .env.example .env
   # Edit .env with your settings
   ```

2. **Build the Server**:

   ```shell
   make da-server
   ```

### Service Modes

#### Avail DA Service

For direct integration with Avail blockchain:

**Prerequisites**:

- Funded Avail wallet with AVAIL tokens
- Valid App ID from Avail network

**Setup Instructions**: [Avail Account Setup Guide](https://docs.availproject.org/docs/end-user-guide)

**Run Server**:

```shell
./bin/avail-da-server \
  --addr=127.0.0.1 \
  --port=3100 \
  --avail.rpc=<Avail RPC URL> \
  --avail.seed="<seed phrase>" \
  --avail.appid=<APP ID> \
```

#### Turbo DA Service

For high-performance data availability via Turbo DA:

**Prerequisites**:

- Turbo DA service URL
- API key from Turbo DA provider

**Setup Instructions**: [Turbo DA Setup Guide](https://docs.availproject.org/da/build-with-avail/turbo-da)

**Run Server**:

```shell
./bin/avail-da-server \
  --addr=127.0.0.1 \
  --port=3100 \
  --turboda=true \
  --avail.rpc=<Avail RPC URL> \
  --turboda.url=<Turbo DA URL> \
  --turboda.key=<Turbo DA API Key> \
```

## Development

### Running Tests

1. **Setup Test Environment**:

   ```shell
   cp .env.example .env
   # Fill in test values in .env - tests require AVAIL_RPC, AVAIL_SEED, and AVAIL_APPID
   ```

2. **Run Tests**:

   ```shell
   make test
   ```

3. **Run Specific Tests**:

   ```shell
   go test -v -tags avail ./test/
   ```

### Building for Different Platforms

```shell
# Build for Linux
TARGETOS=linux TARGETARCH=amd64 make da-server

# Build for macOS
TARGETOS=darwin TARGETARCH=amd64 make da-server

# Build for ARM64
TARGETOS=linux TARGETARCH=arm64 make da-server
```

## Docker Deployment

### Using Docker Compose

1. **Configure Environment**:

   ```shell
   cp .env.example .env
   # Edit .env with your configuration
   ```

2. **Build and Run**:

   ```shell
   docker-compose build
   docker-compose up
   ```

3. **Run in Background**:

   ```shell
   docker-compose up -d
   ```

4. **View Logs**:

   ```shell
   docker-compose logs -f
   ```

5. **Stop Services**:

   ```shell
   docker-compose down
   ```

### Manual Docker Build

```shell
# Build image
docker build -t avail-alt-da-server .

# Run container
docker run -p 3100:3100 --env-file .env avail-alt-da-server
```

## API Reference

The server exposes HTTP endpoints for data availability operations.

### Endpoints

#### POST /put

Submit data to the DA layer.

**Request Body**: Raw binary data

**Response**: Hex-encoded commitment hash

#### GET /get/{commitment}

Retrieve data by hex-encoded commitment.

**Response**: Raw binary data if found, 404 if not found

## Troubleshooting

### Common Issues

#### Connection Errors

- **Issue**: Failed to connect to Avail RPC
- **Solution**: Verify RPC URL is accessible and network connectivity

#### Authentication Errors

- **Issue**: Invalid seed phrase or App ID
- **Solution**: Ensure wallet is funded and App ID is correctly configured

#### Timeout Issues

- **Issue**: Operations timing out
- **Solution**: Increase `--avail.timeout` value or check network conditions

#### Turbo DA Errors

- **Issue**: Invalid Turbo DA credentials
- **Solution**: Verify URL and API key are correct and active

### Debug Mode

Enable verbose logging:

```shell
./bin/avail-da-server --log.level=debug
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Requirements

- Follow Go coding standards
- Add tests for new features
- Update documentation as needed
- Ensure all tests pass before submitting

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- **Documentation**: [Avail Project Docs](https://docs.availproject.org)
- **Issues**: Create an issue in the GitHub repository
- **Community**: Join the Avail Discord community
