# Avail Optimism Alt DA Server

## Introduction

This introduces a sidecar DA Server for Optimism that interacts with Avail DA for posting and retrieving data.

## Configuration

| Flag              | Default Value       |
| ----------------- | ------------------- |
| `--addr`          | `0.0.0.1`           |
| `--port`          | `3100`              |
| `--avail.rpc`     |                     |
| `--avail.seed`    |                     |
| `--avail.appid`   |                     |
| `--avail.timeout` | `100 * time.Second` |

## Usage

- You will need to get a seed phrase, funded with Avail tokens, and an App ID. Steps to generate them can be found [here](https://docs.availproject.org/docs/end-user-guide)

#### Build

```
make da-server
```

#### Run Avail Server

```
./bin/avail-da-server  --addr=localhost --port=8000 --avail.rpc=<Avail RPC URL> --avail.seed="<seed phrase>" --avail.appid=<APP ID> --avail.timeout=<(optional) Timeout>
```

#### Run Tests

- Copy `.env.example` to `.env`. Fill the values inside.

- Run the following command:
  ```
  make test
  ```

### Run using docker

- Copy `.env.example` to `.env`. Fill the values inside.

- Run the following commands:
  ```
  docker-compose build
  docker-compose up
  ```
