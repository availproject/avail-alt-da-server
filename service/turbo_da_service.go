package avail

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/log"
)

type TurboDAService struct {
	APIURL  string
	key     string
	Timeout time.Duration
	log     log.Logger
}

func NewTurboDAService(apiURL string, key string, timeout time.Duration, log log.Logger) (*TurboDAService, error) {
	return &TurboDAService{
		APIURL:  apiURL,
		key:     key,
		Timeout: timeout,
		log:     log,
	}, nil
}

func (s *TurboDAService) Get(ctx context.Context, key []byte) ([]byte, error) {
	// Implement Turbo DA Get logic
	s.log.Info("Turbo DA Get called", "key", string(key))
	return nil, fmt.Errorf("TurboDAService.Get not implemented")
}
func (s *TurboDAService) Put(ctx context.Context, value []byte) ([]byte, error) {
	s.log.Info("Turbo DA Put called", "value_size", len(value))
	return nil, fmt.Errorf("TurboDAService.Put not implemented")
}
