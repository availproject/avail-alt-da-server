package main

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"path"
	"strconv"
	"time"

	altda "github.com/ethereum-optimism/optimism/op-alt-da"
	"github.com/ethereum-optimism/optimism/op-service/rpc"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/log"
)

type AvailStore interface {
	Get(ctx context.Context, key []byte) ([]byte, error)
	Put(ctx context.Context, value []byte) ([]byte, error)
}

type AvailDAServer struct {
	log            log.Logger
	endpoint       string
	store          AvailStore
	tls            *rpc.ServerTLSConfig
	httpServer     *http.Server
	listener       net.Listener
	useGenericComm bool
}

var ErrNotFound = errors.New("not found")

func NewAvailDAServer(host string, port int, store AvailStore, log log.Logger, useGenericComm bool) *AvailDAServer {
	endpoint := net.JoinHostPort(host, strconv.Itoa(port))
	return &AvailDAServer{
		log:      log,
		endpoint: endpoint,
		store:    store,
		httpServer: &http.Server{
			Addr: endpoint,
		},
		useGenericComm: useGenericComm,
	}
}

func (d *AvailDAServer) Start() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/get/", d.HandleGet)
	mux.HandleFunc("/put/", d.HandlePut)

	d.httpServer.Handler = mux

	listener, err := net.Listen("tcp", d.endpoint)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	d.listener = listener

	d.endpoint = listener.Addr().String()
	errCh := make(chan error, 1)
	go func() {
		if d.tls != nil {
			if err := d.httpServer.ServeTLS(d.listener, "", ""); err != nil {
				errCh <- err
			}
		} else {
			if err := d.httpServer.Serve(d.listener); err != nil {
				errCh <- err
			}
		}
	}()

	// verify that the server comes up
	tick := time.NewTimer(10 * time.Millisecond)
	defer tick.Stop()

	select {
	case err := <-errCh:
		return fmt.Errorf("http server failed: %w", err)
	case <-tick.C:
		return nil
	}
}

func (d *AvailDAServer) HandleGet(w http.ResponseWriter, r *http.Request) {
	d.log.Debug("GET", "url", r.URL)

	route := path.Dir(r.URL.Path)
	if route != "/get" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	key := path.Base(r.URL.Path)
	comm, err := hexutil.Decode(key)
	if err != nil {
		d.log.Error("Failed to decode commitment", "err", err, "key", key)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	input, err := d.store.Get(r.Context(), comm)
	if err != nil && errors.Is(err, ErrNotFound) {
		d.log.Error("Commitment not found", "key", key, "error", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		d.log.Error("Failed to read commitment", "err", err, "key", key)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(input); err != nil {
		d.log.Error("Failed to write pre-image", "err", err, "key", key)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (d *AvailDAServer) HandlePut(w http.ResponseWriter, r *http.Request) {
	d.log.Info("PUT", "url", r.URL)

	route := path.Dir(r.URL.Path)
	if route != "/put" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	input, err := io.ReadAll(r.Body)
	if err != nil {
		d.log.Error("Failed to read request body", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	comm, err := d.store.Put(r.Context(), input)
	if err != nil {
		d.log.Error("Failed to store commitment to the DA server", "err", err, "comm", comm)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	d.log.Info("stored commitment", "key", hex.EncodeToString(comm), "input_len", len(input))

	if _, err := w.Write(altda.GenericCommitment(comm).Encode()); err != nil {
		d.log.Error("Failed to write commitment request body", "err", err, "comm", comm)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (b *AvailDAServer) Endpoint() string {
	return b.listener.Addr().String()
}

func (b *AvailDAServer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = b.httpServer.Shutdown(ctx)
	return nil
}
