package avail

import (
	"context"
	"fmt"

	"time"

	"avail-alt-da-server/scripts"
	"avail-alt-da-server/types"
	"avail-alt-da-server/utils"

	SDK "github.com/availproject/avail-go-sdk/sdk"
	"github.com/ethereum/go-ethereum/log"
	"github.com/vedhavyas/go-subkey/v2"
)

type AvailDAService struct {
	Account subkey.KeyPair
	RPCURL  string        `json:"api_url"`
	AppID   int           `json:"app_id"`
	Timeout time.Duration `json:"timeout"`
	log     log.Logger
}

func NewAvailDAService(rpcURL string, seed string, appID int, timeout time.Duration, log log.Logger) (*AvailDAService, error) {

	AppID := utils.EnsureValidAppID(appID)

	keyringPair, err := SDK.Account.NewKeyPair(seed)
	if err != nil {
		log.Warn("⚠️ cannot create LeyPair: error:%w", err)
		return nil, err
	}

	return &AvailDAService{
		Account: keyringPair,
		RPCURL:  rpcURL,
		AppID:   AppID,
		Timeout: timeout,
		log:     log,
	}, nil
}

func (s *AvailDAService) Get(ctx context.Context, comm []byte) ([]byte, error) {
	avail_blk_ref := types.AvailBlockRef{}
	err := avail_blk_ref.UnmarshalFromBinary(comm)
	if err != nil {
		s.log.Error("failed to unmarshal the ethereum tx data to avail block reference", "error", err)
		return []byte{}, fmt.Errorf("failed to unmarshal the ethereum tx data to avail block reference, error: %w", err)
	}

	input, err := scripts.GetBlockExtrinsicData(s.RPCURL, avail_blk_ref, s.log)

	if err != nil {
		s.log.Error("failed to get block extrinsic data", "error", err)
		return []byte{}, fmt.Errorf("failed to get block extrinsic data: %w", err)
	}

	return input, nil
}

func (s *AvailDAService) Put(ctx context.Context, value []byte) ([]byte, error) {

	if len(value) >= 512000 {
		return nil, fmt.Errorf("the length of input cannot be greater than 512kb")
	}

	avail_Blk_Ref, err := scripts.SubmitDataAndWatch(s.RPCURL, s.Account, s.AppID, ctx, value, s.log)

	if err != nil {
		s.log.Error("cannot submit data", "error", err)
		return nil, fmt.Errorf("cannot submit data:%w", err)
	}

	comm, err := avail_Blk_Ref.MarshalToBinary()

	if err != nil {
		s.log.Error("cannot get the binary form of avail block reference", "error", err)
		return nil, fmt.Errorf("cannot get the binary form of avail block reference:%w", err)
	}

	return comm, nil
}
