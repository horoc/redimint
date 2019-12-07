package models

import (
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/types"
	"github.com/tendermint/tendermint/version"
	"time"
)

//Execution
type CommandRequest struct {
	Cmd string `json:"cmd"`
}

type ExecuteResponse struct {
	Code          uint32 `json:"code"`
	CodeMsg       string `json:"code_info"`
	Cmd           string `json:"command"`
	ExecuteResult string `json:"execute_result"`
	Signature     string `json:"signature"`
	Sequence      string `json:"sequence"`
	TimeStamp     string `json:"time_stamp"`
	Hash          string `json:"hash"`
	Height        int64  `json:"height"`
}

type ExecuteAsyncResponse struct {
	Code      uint32 `json:"code"`
	CodeMsg   string `json:"code_info"`
	Cmd       string `json:"command"`
	Signature string `json:"signature"`
	Sequence  string `json:"sequence"`
	TimeStamp string `json:"time_stamp"`
	Hash      string `json:"hash"`
}

type QueryResponse struct {
	Code    uint32 `json:"code"`
	CodeMsg string `json:"code_info"`
	Result  string `json:"result"`
}

//Transaction
type Transaction struct {
	Hash          string        `json:"hash"`
	Height        int64         `json:"height"`
	Index         uint32        `json:"index"`
	Data          *TxCommitBody `json:"data"`
	ExecuteResult string        `json:"execute_result"`
	Proof         *TxProof      `json:"proof,omitempty"`
}

type TransactionList struct {
	Height int64         `json:"height"`
	Total  int64         `json:"total"`
	Txs    []Transaction `json:"txs"`
}

type TxProof struct {
	RootHash string       `json:"root_hash"`
	Proof    *ProofDetail `json:"proof"`
}

type ProofDetail struct {
	Total    int      `json:"total"`     // Total number of items.
	Index    int      `json:"index"`     // Index of item to prove.
	LeafHash string   `json:"leaf_hash"` // Hash of item value.
	Aunts    []string `json:"aunts"`     // Hashes from leaf's sibling to a root's child.
}

//Block
type BlockMeta struct {
	BlockID `json:"block_id"`
	Header  `json:"header"`
}

type Block struct {
	BlockID    `json:"block_id"`
	Header     `json:"header"`
	Data       `json:"data"`
	Evidence   types.EvidenceData `json:"evidence"`
	LastCommit []*CommitSig       `json:"last_commit"`
}

type BlockID struct {
	Hash        string        `json:"hash"`
	PartsHeader PartSetHeader `json:"parts"`
}

type PartSetHeader struct {
	Total int    `json:"total"`
	Hash  string `json:"hash"`
}

type Header struct {
	// basic block info
	Version  version.Consensus `json:"version"`
	ChainID  string            `json:"chain_id"`
	Height   int64             `json:"height"`
	Time     time.Time         `json:"time"`
	NumTxs   int64             `json:"num_txs"`
	TotalTxs int64             `json:"total_txs"`

	// prev block info
	LastBlockID BlockID `json:"last_block_id"`

	// hashes of block data
	LastCommitHash string `json:"last_commit_hash"` // commit from validators from the last block
	DataHash       string `json:"data_hash"`        // transactions

	// hashes from the app output from the prev block
	ValidatorsHash     string `json:"validators_hash"`      // validators for the current block
	NextValidatorsHash string `json:"next_validators_hash"` // validators for the next block
	ConsensusHash      string `json:"consensus_hash"`       // consensus params for current block
	AppHash            string `json:"app_hash"`             // state after txs from the previous block
	LastResultsHash    string `json:"last_results_hash"`    // root hash of all results from the txs from the previous block

	// consensus info
	EvidenceHash    string `json:"evidence_hash"`    // evidence included in the block
	ProposerAddress string `json:"proposer_address"` // original proposer of the block
}

type Data struct {
	// Txs that will be applied by state @ block.Height+1.
	// NOTE: not all txs here are valid.  We're just agreeing on the order first.
	// This means that block.AppHash does not include these txs.
	Txs []string `json:"txs"`

	Hash string `json:"hash"`
}

type CommitSig struct {
	Type             types.SignedMsgType `json:"type"`
	Height           int64               `json:"height"`
	Round            int                 `json:"round"`
	Timestamp        time.Time           `json:"timestamp"`
	ValidatorAddress string              `json:"validator_address"`
	ValidatorIndex   int                 `json:"validator_index"`
	Signature        string              `json:"signature"`
}

//chain
type ChainInfo struct {
	Code       uint32       `json:"code"`
	CodeMsg    string       `json:"code_info"`
	LastHeight int64        `json:"last_height"`
	BlockMetas []*BlockMeta `json:"block_metas"`
}

type NodeInfo struct {
	ProtocolVersion p2p.ProtocolVersion `json:"protocol_version"`
	ID              p2p.ID              `json:"id"`          // authenticated identifier
	ListenAddr      string              `json:"listen_addr"` // accepting incoming
	Network         string              `json:"network"`     // network/chain ID
	Version         string              `json:"version"`     // major.minor.revision
	Channels        HexBytes            `json:"channels"`    // channels this node knows about¬
	Moniker         string              `json:"moniker"`     // arbitrary moniker
}

type SyncInfo struct {
	LatestBlockHash   HexBytes  `json:"latest_block_hash"`
	LatestAppHash     HexBytes  `json:"latest_app_hash"`
	LatestBlockHeight int64     `json:"latest_block_height"`
	LatestBlockTime   time.Time `json:"latest_block_time"`
	CatchingUp        bool      `json:"catching_up"`
}

type ValidatorInfo struct {
	Address     HexBytes `json:"address"`
	PubKey      HexBytes `json:"pub_key"`
	VotingPower int64    `json:"voting_power"`
}

type ChainState struct {
	Code          uint32        `json:"code"`
	CodeMsg       string        `json:"code_info"`
	NodeInfo      NodeInfo      `json:"node_info"`
	SyncInfo      SyncInfo      `json:"sync_info"`
	ValidatorInfo ValidatorInfo `json:"validator_info"`
}

//bench mark
type BenchMarkRequest struct {
	TxNums       int    `json:"tx_nums"`
	TxSendPerSec int    `json:"tx_send_per_sec"`
	Connections  int    `json:"connections"`
	Mode         string `json:"mode"`
	Cmd          string `json:"cmd"`
}

type BenchMarkResponse struct {
	Latency *BenchMarkDetail
	Tps     *BenchMarkDetail
}

type BenchMarkDetail struct {
	Avg   string
	Max   string
	Stdev string
}
