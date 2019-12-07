package models

type ExecuteRequest struct {
	Cmd  string `json:"cmd"`
	Mode string `json:"mode"`
}

type TxHashRequest struct {
	Hash string `json:"hash"`
}

type BlockHeightRequest struct {
	Height int `json:"height"`
}

type CommittedTxListRequest struct {
	Begin int `json:"begin_height"`
	End   int `json:"end_height"`
}

type ChainInfoRequest struct {
	Min int `json:"min"`
	Max int `json:"max"`
}
