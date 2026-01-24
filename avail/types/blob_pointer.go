package types

import (
	"fmt"

	"github.com/availproject/avail-go-sdk/primitives"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// ---------------------- Transaction Details ----------------------

type TransactionDetails struct {
	BlockNumber uint32
	BlockHash   primitives.H256
	TxIndex     uint32
	Commitment  common.Hash
}

// -------------------- ABI Types --------------------
var (
	unit8Type  = abi.Type{T: abi.UintTy, Size: 8}
	byte32Type = abi.Type{T: abi.FixedBytesTy, Size: 32}
	uint32Type = abi.Type{Size: 32, T: abi.UintTy}
)

// -------------------- BlobPointer --------------------
// BlobPointer version
const (
	BLOBPOINTER_VERSION0 = 0x00
	BLOBPOINTER_VERSION1 = 0x01
	BLOBPOINTER_VERSION2 = 0x02
	BLOBPOINTER_VERSION3 = 0x03
	BLOBPOINTER_VERSION4 = 0x04
)

// BlobPointer contains the reference to the data blob on Avail
type BlobPointer struct {
	Version            uint8
	BlockHeight        uint32      // Block height for avail chain in which data in being included
	ExtrinsicIndex     uint32      // extrinsic index in the block height
	BlobDataKeccak265H common.Hash // Keccak256(blobData) to verify the originality of proof (it will work as preimage of the commitment)
}

var blobPointerArguments = abi.Arguments{
	{Type: unit8Type}, {Type: uint32Type}, {Type: uint32Type}, {Type: byte32Type},
}

func NewBlobPointer(blockHeight uint32, extrinsicIndex uint32, dataCommitment common.Hash) *BlobPointer {
	return &BlobPointer{
		Version:            BLOBPOINTER_VERSION0,
		BlockHeight:        blockHeight,
		ExtrinsicIndex:     extrinsicIndex,
		BlobDataKeccak265H: dataCommitment,
	}
}

func (b *BlobPointer) MarshalToBinary() ([]byte, error) {
	packedData, err := blobPointerArguments.PackValues([]interface{}{b.Version, b.BlockHeight, b.ExtrinsicIndex, b.BlobDataKeccak265H})
	if err != nil {
		return []byte{}, fmt.Errorf("unable to covert the blobPointer into array of bytes and getting error:%w", err)
	}
	return packedData, nil
}

func (b *BlobPointer) UnmarshalFromBinary(data []byte) error {
	unpackedData, err := blobPointerArguments.UnpackValues(data)
	if err != nil {
		return fmt.Errorf("unable to covert the data bytes into blobPointer and getting error:%w", err)
	}
	b.Version = unpackedData[0].(uint8)
	b.BlockHeight = unpackedData[1].(uint32)
	b.ExtrinsicIndex = unpackedData[2].(uint32)
	b.BlobDataKeccak265H = unpackedData[3].([32]uint8)
	return nil
}

// Method to convert BlobPointer to string
func (bp *BlobPointer) String() string {
	return fmt.Sprintf(
		"BlockHeight: %d,  ExtrinsicIndex: %d,  BlobDataKeccak265H: %s",
		bp.BlockHeight,
		bp.ExtrinsicIndex,
		bp.BlobDataKeccak265H.Hex(),
	)
}
