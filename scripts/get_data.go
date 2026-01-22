package scripts

import (
	"fmt"

	"github.com/availproject/avail-go-sdk/primitives"
	SDK "github.com/availproject/avail-go-sdk/sdk"
	"github.com/ethereum/go-ethereum/log"
)

func GetDatafromAvail(sdk *SDK.SDK, blockNumber uint32, index uint32) ([]byte, error) {
	blockHash, err := sdk.Client.BlockHash(blockNumber)
	if err != nil {
		return nil, fmt.Errorf("❎ Cannot get block hash: %w", err)
	}

	block, err := SDK.NewBlock(sdk.Client, blockHash)
	if err != nil {
		return nil, fmt.Errorf("❎ Cannot get block: %w", err)
	}

	var blob SDK.DataSubmission

	blobs := block.DataSubmissions(SDK.Filter{}.WTxIndex(index))
	if len(blobs) == 0 {
		return nil, fmt.Errorf("❎ No blobs found for transaction index %d in block %d", index, blockNumber)
	}
	blob = blobs[0]

	signerAddress, err := primitives.NewAccountIdFromMultiAddress(blob.TxSigner)
	if err != nil {
		log.Warn("AvailDAWarn:‼️ Unable to extract the signer address for the blob")
	}

	log.Debug("AvailDADebug: ✅ Data retrieved from Avail chain signer: %s, appID: %d, extrinsicHash: %s",
		signerAddress.ToHuman(),
		blob.AppId,
		blob.TxHash,
	)

	return blob.Data, nil
}
