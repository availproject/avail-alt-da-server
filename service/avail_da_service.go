package avail

import (
	"context"
	"fmt"

	"time"

	"avail-alt-da-server/scripts"
	"avail-alt-da-server/types"

	SDK "github.com/availproject/avail-go-sdk/sdk"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/vedhavyas/go-subkey/v2"
)

const (
	AvailNetworkID = 42
)

type AvailDAService struct {
	SDK     *SDK.SDK
	Account subkey.KeyPair
	RPCURL  string        `json:"api_url"`
	AppID   int           `json:"app_id"`
	Timeout time.Duration `json:"timeout"`
	log     log.Logger
}

func NewAvailDAService(rpcURL string, seed string, appID int, timeout time.Duration, log log.Logger) (*AvailDAService, error) {

	sdk, err := SDK.NewSDK(rpcURL)
	if err != nil {
		log.Error("failed to create SDK", "error", err)
		return nil, err
	}

	AppID := validateAppID(appID)

	keyringPair, err := SDK.Account.NewKeyPair(seed)
	if err != nil {
		log.Warn("⚠️ cannot create KeyPair: error:%w", err)
		return nil, err
	}

	return &AvailDAService{
		SDK:     &sdk,
		Account: keyringPair,
		RPCURL:  rpcURL,
		AppID:   AppID,
		Timeout: timeout,
		log:     log,
	}, nil
}

func (s *AvailDAService) Get(ctx context.Context, comm []byte) ([]byte, error) {
	s.log.Info("AvailDAInfo: 📥 Received Get request", "comm", comm)
	blobPointer := &types.BlobPointer{}
	if err := blobPointer.UnmarshalFromBinary(comm); err != nil {
		return nil, fmt.Errorf("failed to decode BlobPointer: %w", err)
	}
	data, err := scripts.GetDatafromAvail(s.SDK, blobPointer.BlockHeight, blobPointer.ExtrinsicIndex)
	if err != nil {
		s.log.Error("failed to retrieve blob data", "error", err)
		return []byte{}, fmt.Errorf("failed to retrieve blob data: %w", err)
	}
	return data, nil
}

func (s *AvailDAService) Put(ctx context.Context, value []byte) ([]byte, error) {

	if len(value) >= 512000 {
		return nil, fmt.Errorf("the length of input cannot be greater than 512kb")
	}

	txDetails, err := submitDataToAvailDA(ctx, s.SDK, s.Account, s.AppID, value, s.log)
	if err != nil {
		s.log.Error("AvailError: ⚠️ cannot submit data", "error", err)
		return nil, fmt.Errorf("cannot submit data:%w", err)
	}

	blobPointer := types.NewBlobPointer(txDetails.BlockNumber, txDetails.TxIndex, txDetails.Commitment)
	payload, err := blobPointer.MarshalToBinary()
	if err != nil {
		return nil, fmt.Errorf("encode blob pointer failed: %w", err)
	}

	return payload, nil
}

func submitDataToAvailDA(ctx context.Context, sdk *SDK.SDK, acc subkey.KeyPair, appID int, data []byte, log log.Logger) (types.TransactionDetails, error) {

	resultCh := make(chan struct {
		details types.TransactionDetails
		err     error
	}, 1)

	// Run the blocking SDK call in a goroutine
	go func() {
		log.Info("AvailDAInfo: 📤 Submitting data to Avail chain")
		tx := sdk.Tx.DataAvailability.SubmitData(data)
		txDetails, err := tx.ExecuteAndWatchFinalization(
			acc,
			SDK.NewTransactionOptions().WithAppId(uint32(appID)),
		)
		if err == nil {
			status := txDetails.IsSuccessful().Unwrap()
			if !status {
				err = fmt.Errorf("extrinsic failed on avail chain, status: %v", status)
			}
		}

		resultCh <- struct {
			details types.TransactionDetails
			err     error
		}{types.TransactionDetails{BlockNumber: txDetails.BlockNumber, BlockHash: txDetails.BlockHash, TxIndex: txDetails.TxIndex, Commitment: crypto.Keccak256Hash(data)}, err}
	}()

	// Now wait for either SDK result or context cancellation
	select {
	case <-ctx.Done():
		return types.TransactionDetails{}, ctx.Err()
	case res := <-resultCh:
		if res.err != nil {
			return types.TransactionDetails{}, fmt.Errorf("⚠️ extrinsic got rejected: %w", res.err)
		}

		log.Debug("AvailDADebug: ✅ Data is included in Avail chain address=%s appID=%d block_number=%d block_hash=%s tx_index=%d",
			acc.SS58Address(AvailNetworkID),
			appID,
			res.details.BlockNumber,
			res.details.BlockHash,
			res.details.TxIndex,
		)
		log.Info("AvailDAInfo: 📤 Data submitted to Avail chain")
		return types.TransactionDetails{BlockNumber: res.details.BlockNumber, BlockHash: res.details.BlockHash, TxIndex: res.details.TxIndex}, nil
	}
}

func validateAppID(appID int) int {
	if appID > 0 {
		return appID
	}
	return 0
}
