package types

import (
	"math/big"

	"github.com/ubiq/go-ubiq/common"
)

const (
	BloomByteLength = 256
	BloomBitLength  = 8 * BloomByteLength
)

type BlockNonce [8]byte

type Bloom [BloomByteLength]byte

type Header struct {
	ParentHash  common.Hash    `json:"parentHash"       gencodec:"required"`
	UncleHash   common.Hash    `json:"sha3Uncles"       gencodec:"required"`
	Coinbase    common.Address `json:"miner"            gencodec:"required"`
	Root        common.Hash    `json:"stateRoot"        gencodec:"required"`
	TxHash      common.Hash    `json:"transactionsRoot" gencodec:"required"`
	ReceiptHash common.Hash    `json:"receiptsRoot"     gencodec:"required"`
	Bloom       Bloom          `json:"logsBloom"        gencodec:"required"`
	Difficulty  *big.Int       `json:"difficulty"       gencodec:"required"`
	Number      *big.Int       `json:"number"           gencodec:"required"`
	GasLimit    uint64         `json:"gasLimit"         gencodec:"required"`
	GasUsed     uint64         `json:"gasUsed"          gencodec:"required"`
	Time        uint64         `json:"timestamp"        gencodec:"required"`
	Extra       []byte         `json:"extraData"        gencodec:"required"`
	MixDigest   common.Hash    `json:"mixHash"`
	Nonce       BlockNonce     `json:"nonce"`
	BaseFee     *big.Int       `json:"baseFeePerGas"`
}

type Transaction struct {
	AccountNonce uint64          `json:"nonce"    gencodec:"required"`
	Price        *big.Int        `json:"gasPrice" gencodec:"required"`
	GasLimit     uint64          `json:"gas"      gencodec:"required"`
	Sender       common.Address  `json:"from"     gencodec:"required"`
	Recipient    *common.Address `json:"to"       rlp:"nil"` // nil means contract creation
	Amount       *big.Int        `json:"value"    gencodec:"required"`
	Payload      []byte          `json:"input"    gencodec:"required"`
}

type Logs struct {
	Address     common.Address `json:"address" gencodec:"required"`
	Topics      []common.Hash  `json:"topics" gencodec:"required"`
	Data        []byte         `json:"data" gencodec:"required"`
	BlockNumber uint64         `json:"blockNumber" rlp:"-"`
	TxHash      common.Hash    `json:"transactionHash" gencodec:"required" rlp:"-"`
	TxIndex     uint           `json:"transactionIndex" rlp:"-"`
	BlockHash   common.Hash    `json:"blockHash" rlp:"-"`
	Index       uint           `json:"logIndex" rlp:"-"`
	Removed     bool           `json:"removed" rlp:"-"`
}

type Receipt struct {
	PostStateOrStatus []byte
	CumulativeGasUsed uint64
	TxHash            common.Hash
	ContractAddress   common.Address
	Logs              []*Logs
	GasUsed           uint64
}
type BlockResult struct {
	Hash            common.Hash
	TotalDifficulty *big.Int
	Header          *Header
	Transactions    []*Transaction
	Uncles          []*Header
	Receipts        []*Receipt
	Senders         []interface{}
}
