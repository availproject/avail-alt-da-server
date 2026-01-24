package avail

import (
	"avail-alt-da-server/avail/service"
	"avail-alt-da-server/flags"
	"context"
	"fmt"

	altda "github.com/ethereum-optimism/optimism/op-alt-da"
	log "github.com/ethereum/go-ethereum/log"
)

const AvailByte = 0x0a

type DAService interface {
	Get(ctx context.Context, key []byte) ([]byte, error)
	Put(ctx context.Context, value []byte) ([]byte, error)
}

type DAProvider struct {
	DAservice DAService
	DAByte    byte
}

func NewDAProvider(cfg flags.CLIConfig, l log.Logger) (DAProvider, error) {
	var daservice DAService
	var err error
	if cfg.TurboDA {
		daservice, err = service.NewTurboDAService(cfg.TurboDAURL, cfg.RPC, cfg.TurboDAKey, l)
		if err != nil {
			return DAProvider{}, fmt.Errorf("failed to create turbo da service: %w", err)
		}
	} else {
		daservice, err = service.NewAvailDAService(cfg.RPC, cfg.Seed, cfg.AppId, l)
		if err != nil {
			return DAProvider{}, fmt.Errorf("failed to create avail da service: %w", err)
		}
	}
	return DAProvider{DAservice: daservice, DAByte: AvailByte}, nil
}

func (d *DAProvider) Encode(comm []byte) []byte {
	return altda.GenericCommitment(append([]byte{byte(d.DAByte)}, comm...)).Encode()
}

func (d *DAProvider) Decode(comm []byte) ([]byte, error) {
	if comm[0] != 0x01 && comm[1] != d.DAByte {
		return nil, fmt.Errorf("invalid encoding")
	}
	return comm[2:], nil
}
