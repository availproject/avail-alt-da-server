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
| `--addr` | Address to bind the server to | `0.0.0.0` |
| `--host` | Host name for the server | `localhost` |
| `--port` | Port to run the server on | `8080` |
| `--avail.rpc` | Avail HTTP RPC URL (REQUIRED) | - |
| `--avail.seed` | Avail seed phrase for authentication | - |
| `--avail.appid` | Avail App ID for your application | `0` |
| `--avail.timeout` | Timeout for Avail operations | `100 * time.Second` |
| `--avail.turboda` | Enable Turbo DA service | `false` |
| `--avail.turboda.url` | Turbo DA service URL | - |
| `--avail.turboda.key` | Turbo DA API key | - |

### Environment Variables

Copy the example environment file:

```shell
cp .env.example .env
```

Edit `.env` with your configuration:

```bash
# Server Configuration
ADDR=0.0.0.0
PORT=8080

# Avail DA Configuration
AVAIL_RPC=https://turing-rpc.avail.so/rpc
AVAIL_SEED="your seed phrase here"
AVAIL_APPID=1

# Turbo DA Configuration (Optional)
TURBODA=false
TURBODA_URL=https://turing.turbo-api.availproject.org
TURBODA_KEY=your-api-key-here
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
  --addr=localhost \
  --port=8000 \
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
  --addr=localhost \
  --port=8000 \
  --avail.turboda=true \
  --avail.rpc=<Avail RPC URL> \
  --turboda.url=<Turbo DA URL> \
  --turboda.key=<Turbo DA API Key> \
```

## Development

### Running Tests

1. **Setup Test Environment**:

   ```shell
   cp .env.example .env
   # Fill in test values in .env
   ```

2. **Run Tests**:

   ```shell
   make test
   ```

3. **Run Specific Tests**:

   ```shell
   go test -v -tags avail ./...
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
docker run -p 8080:8080 --env-file .env avail-alt-da-server
```

## API Reference

The server exposes HTTP endpoints for data availability operations.

### Endpoints

#### POST /submit

Submit data to the DA layer.

**Request Body**:

```json
{
  "data": "base64-encoded-data",
  "namespace": "optional-namespace"
}
```

**Response**:

```json
{
  "commitment": "data-commitment-hash",
  "namespace": "data-namespace",
  "height": "block-height"
}
```

#### GET /retrieve/{commitment}

Retrieve data by commitment.

**Response**:

```json
{
  "data": "base64-encoded-data",
  "found": true
}
```

#### GET /health

Health check endpoint.

**Response**:

```json
{
  "status": "healthy",
  "timestamp": "2024-01-01T00:00:00Z"
}
```

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
