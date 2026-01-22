package scripts

import (
	"fmt"

	"avail-alt-da-server/types"

	"github.com/availproject/avail-go-sdk/metadata"
	"github.com/availproject/avail-go-sdk/primitives"
	SDK "github.com/availproject/avail-go-sdk/sdk"
	"github.com/ethereum/go-ethereum/log"
)

func GetBlockExtrinsicData(apiUrl string, avail_blk_ref types.AvailBlockRef, log log.Logger) ([]byte, error) {

	Hash := avail_blk_ref.BlockHash
	Address := avail_blk_ref.Sender
	Nonce := avail_blk_ref.Nonce

	avail_blk, err := fetchBlock(apiUrl, Hash)
	if err != nil {
		log.Error("cannot fetch block", "error", err)
		return []byte{}, fmt.Errorf("cannot fetch block: %w", err)
	}

	return extractExtrinsic(Address, Hash, Nonce, avail_blk)
}

func fetchBlock(apiURL string, Hash string) (SDK.Block, error) {

	sdk, err := SDK.NewSDK(apiURL)
	if err != nil {
		log.Error("cannot create sdk", "error", err)
		return SDK.Block{}, fmt.Errorf("cannot create sdk: %w", err)
	}

	blockHash, err := primitives.NewBlockHashFromHexString(Hash)
	if err != nil {
		log.Error("unable to convert string hash into types.hash", "error", err)
		return SDK.Block{}, fmt.Errorf("unable to convert string hash into types.hash, error:%w", err)
	}
	block, err := SDK.NewBlock(sdk.Client, blockHash)
	if err != nil {
		log.Error("unable to create block", "error", err)
		return SDK.Block{}, fmt.Errorf("unable to create block, error:%w", err)
	}

	return block, nil
}

func extractExtrinsic(Address string, Hash string, Nonce int64, avail_blk SDK.Block) ([]byte, error) {
	accountId, err := metadata.NewAccountIdFromAddress(Address)
	if err != nil {
		log.Error("unable to create account id from address", "error", err)
		return []byte{}, fmt.Errorf("unable to create account id from address: %v, error: %w", Address, err)
	}

	for _, blob := range avail_blk.Block.Extrinsics {

		ext_Nonce := blob.Signed.Unwrap().Nonce

		if sameSignature(&blob, accountId) && ext_Nonce == uint32(Nonce) {
			extrinsic, ok := SDK.NewDataSubmission(&blob)
			if !ok {
				log.Error("unable to create data submission from extrinsic")
				return []byte{}, fmt.Errorf("unable to create data submission from extrinsic")
			}
			return extrinsic.Data, nil
		}
	}

	log.Error("didn't find any extrinsic data for address:%v in block having hash:%v", Address, Hash)
	return []byte{}, fmt.Errorf("didn't find any extrinsic data for address:%v in block having hash:%v", Address, Hash)
}

func sameSignature(tx *primitives.DecodedExtrinsic, accountId metadata.AccountId) bool {
	txAccountId := tx.Signed.Unwrap().Address.Id.Unwrap()
	if accountId.Value != txAccountId {
		return false
	}

	return true
}
