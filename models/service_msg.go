package models

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

type Transaction struct {
	Hash          string        `json:"hash"`
	Height        int64         `json:"height"`
	Index         uint32        `json:"index"`
	Data          *TxCommitBody `json:"data"`
	ExecuteResult string        `json:"execute_result"`
	Proof         *TxProof       `json:"proof,omitempty"`
}

type TxProof struct {
	RootHash string      `json:"root_hash"`
	Proof    *ProofDetail `json:"proof"`
}

type ProofDetail struct {
	Total    int      `json:"total"`     // Total number of items.
	Index    int      `json:"index"`     // Index of item to prove.
	LeafHash string   `json:"leaf_hash"` // Hash of item value.
	Aunts    []string `json:"aunts"`     // Hashes from leaf's sibling to a root's child.
}

type ExecuteRequest struct {
	Cmd string `json:"command"`
}


