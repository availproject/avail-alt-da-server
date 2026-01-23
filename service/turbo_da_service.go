package avail

import (
	"avail-alt-da-server/scripts"
	"avail-alt-da-server/types"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/availproject/avail-go-sdk/primitives"
	SDK "github.com/availproject/avail-go-sdk/sdk"
	"github.com/ethereum/go-ethereum/log"
)

const (
	MAX_ATTEMPTS = 10
)

type TurboDAService struct {
	SDK    *SDK.SDK
	APIURL string
	key    string
	log    log.Logger
}

func NewTurboDAService(apiURL string, rpcURL string, key string, log log.Logger) (*TurboDAService, error) {
	sdk, err := SDK.NewSDK(rpcURL)
	if err != nil {
		log.Error("AvailDAError: ❌ failed to create avail sdk", "error", err)
		return nil, fmt.Errorf("failed to create avail sdk: %w", err)
	}
	return &TurboDAService{
		SDK:    &sdk,
		APIURL: apiURL,
		key:    key,
		log:    log,
	}, nil
}

func (s *TurboDAService) Get(ctx context.Context, comm []byte) ([]byte, error) {
	s.log.Info("AvailDAInfo: 📥 Received Get request", "comm", comm)
	blobPointer := &types.BlobPointer{}
	if err := blobPointer.UnmarshalFromBinary(comm); err != nil {
		s.log.Error("AvailDAError: ❌ failed to decode BlobPointer", "error", err)
		return nil, fmt.Errorf("failed to decode BlobPointer: %w", err)
	}
	data, err := scripts.GetDatafromAvail(s.SDK, blobPointer.BlockHeight, blobPointer.ExtrinsicIndex)
	if err != nil {
		s.log.Error("AvailDAError: ❌ failed to retrieve blob data", "error", err)
		return []byte{}, fmt.Errorf("failed to retrieve blob data: %w", err)
	}
	return data, nil
}
func (s *TurboDAService) Put(ctx context.Context, value []byte) ([]byte, error) {
	if len(value) >= 512000 {
		return nil, fmt.Errorf("the length of input cannot be greater than 512kb")
	}

	txDetails, err := submitDataToTurboDA(ctx, s.log, s.APIURL, s.key, value, MAX_ATTEMPTS)
	if err != nil {
		s.log.Error("AvailDAError: ⚠️ cannot submit data", "error", err)
		return nil, fmt.Errorf("cannot submit data:%w", err)
	}

	blobPointer := types.NewBlobPointer(txDetails.BlockNumber, txDetails.TxIndex, txDetails.Commitment)
	payload, err := blobPointer.MarshalToBinary()
	if err != nil {
		s.log.Error("AvailDAError: ❌ failed to encode blob pointer", "error", err)
		return nil, fmt.Errorf("encode blob pointer failed: %w", err)
	}

	return payload, nil
}

type TurboDAResponse struct {
	SubmissionID string `json:"submission_id"`
}

type TurboDAStatusResponse struct {
	Data  *TurboDAStatusData `json:"data"`
	Error interface{}        `json:"error"`
	ID    string             `json:"id"`
	State string             `json:"state"`
}

type TurboDAStatusData struct {
	AmountData  string `json:"amount_data"`
	BlockHash   string `json:"block_hash"`
	BlockNumber int64  `json:"block_number"`
	CreatedAt   string `json:"created_at"`
	DataBilted  string `json:"data_bilted"`
	DataHash    string `json:"data_hash"`
	Fees        string `json:"fees"`
	TxHash      string `json:"tx_hash"`
	TxIndex     int64  `json:"tx_index"`
	UserID      string `json:"user_id"`
}

func submitDataToTurboDA(ctx context.Context, logger log.Logger, url string, apiKey string, data []byte, maxAttempts int) (types.TransactionDetails, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url+"/v1/submit_raw_data", bytes.NewReader(data))
	if err != nil {
		return types.TransactionDetails{}, fmt.Errorf("failed to create request to Turbo DA: %w", err)
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("x-api-key", apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return types.TransactionDetails{}, fmt.Errorf("failed to post data to Turbo DA: %w", err)
	}
	defer resp.Body.Close()

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return types.TransactionDetails{}, fmt.Errorf("failed to read response from Turbo DA: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return types.TransactionDetails{}, fmt.Errorf("turbo da error: status code %d", resp.StatusCode)
	}

	var postResp TurboDAResponse
	if err := json.Unmarshal(respData, &postResp); err != nil {
		return types.TransactionDetails{}, fmt.Errorf("failed to unmarshal response from Turbo DA: %w", err)
	}

	getURL := fmt.Sprintf("%s/v1/get_submission_info?submission_id=%s", url, postResp.SubmissionID)

	var statusResp TurboDAStatusResponse
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		select {
		case <-ctx.Done():
			return types.TransactionDetails{}, ctx.Err()
		default:
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
		if err != nil {
			return types.TransactionDetails{}, fmt.Errorf("failed to create status request: %w", err)
		}
		req.Header.Set("x-api-key", apiKey)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			logger.Warn("AvailDAWarn: ⚠️ Failed to fetch submission status", "error", err)
			logger.Debug("AvailDAinfo", "attempt", attempt, "maxAttempts", maxAttempts, "error", err)
		} else {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()

			if err := json.Unmarshal(body, &statusResp); err != nil {
				logger.Warn("AvailDAWarn: ⚠️ Invalid JSON from Turbo DA", "error", err)
			} else {
				logger.Debug("AvailDAInfo: ⏳ Attempt info", "attempt", attempt, "maxAttempts", maxAttempts, "status", statusResp.State)
				if statusResp.State == "Finalized" {
					logger.Debug("AvailDAInfo: ✅ Turbo DA finalized submission", "submissionID", postResp.SubmissionID)
					blockHash, err := primitives.NewBlockHashFromHexString(statusResp.Data.BlockHash)
					if err != nil {
						return types.TransactionDetails{}, fmt.Errorf("invalid block hash from Turbo DA: %w", err)
					}
					return types.TransactionDetails{BlockNumber: uint32(statusResp.Data.BlockNumber), BlockHash: blockHash, TxIndex: uint32(statusResp.Data.TxIndex)}, nil
				}
			}
		}

		// exponential backoff with cap
		sleep := time.Duration(minInt(30, 2<<attempt)) * time.Second
		logger.Debug("AvailDAInfo: 🕒 Waiting before the next check...", "sleep", sleep)
		time.Sleep(sleep)
	}

	return types.TransactionDetails{}, fmt.Errorf("submission %s did not finalize within %d attempts", postResp.SubmissionID, maxAttempts)
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
