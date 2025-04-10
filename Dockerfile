FROM golang:1.23.4 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN make da-server

FROM golang:1.23.4

WORKDIR /app

COPY --from=builder /app/bin/avail-da-server /app/bin/

ENV ADDR=0.0.0.0
ENV PORT=8080
ENV AVAIL_RPC=http://localhost:9933
ENV AVAIL_SEED=""
ENV AVAIL_APPID=0
ENV AVAIL_TIMEOUT=100s

EXPOSE ${PORT}
EXPOSE 8080
EXPOSE 433

CMD ["sh", "-c", "/app/bin/avail-da-server --addr=$ADDR --port=$PORT --avail.rpc=\"$AVAIL_RPC\" --avail.seed=\"$AVAIL_SEED\" --avail.appid=$AVAIL_APPID --avail.timeout=$AVAIL_TIMEOUT"]
