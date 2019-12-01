package models

type TxHashRequest struct {
	Hash string `json:"hash"`
}

type BlockHeightRequest struct {
	Height int `json:"height"`
}
