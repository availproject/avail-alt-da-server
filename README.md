# Avail Optimism Alt DA Server

## Introduction

This introduces a sidecar DA Server for Optimism that interacts with Avail DA for posting and retrieving data.

## Usage

- You will need to get a seed phrase, funded with Avail tokens, and an App ID. Steps to generate them can be found [here](https://docs.availproject.org/docs/end-user-guide)

#### Build

```
make da-server
```

#### Run Avail Server

```
go run ./cmd/avail  --addr=localhost --port=8000 --avail.rpc=<Avail RPC URL> --avail.seed="<seed phrase>" --avail.appid=<APP ID> --avail.timeout=<Timeout>
```

#### Run Tests

- Fill the following values inside `daclient_avail_test.go`

  ```
    const (
  	RPC     = ""                // RPC URL
  	SEED    = ""                // SEED PHRASE
  	APPID   = 0                 // APP ID
  	TIMEOUT = 100 * time.Second // TIMEOUT
  )
  ```

- Run the following command:
  ```
  make test
  ```
