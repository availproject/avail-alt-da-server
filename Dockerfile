FROM golang:1.23.4 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN make da-server

FROM golang:1.23.4

WORKDIR /app

COPY --from=builder /app/bin/avail-da-server /app/bin/

CMD ["/app/bin/avail-da-server"]