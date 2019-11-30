package models

import (
	"github.com/tendermint/tendermint/types"
	"github.com/tendermint/tendermint/version"
	"time"
)

//Execution
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

type ExecuteRequest struct {
	Cmd string `json:"command"`
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
type Block struct {
	BlockID    `json:"block_id"`
	Header     `json:"header"`
	Data       `json:"data"`
	Evidence   types.EvidenceData `json:"evidence"`
	LastCommit []*CommitSig             `json:"last_commit"`
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
	Height           int64                `json:"height"`
	Round            int                  `json:"round"`
	Timestamp        time.Time            `json:"timestamp"`
	ValidatorAddress string               `json:"validator_address"`
	ValidatorIndex   int                  `json:"validator_index"`
	Signature        string               `json:"signature"`
}
