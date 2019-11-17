package models

type TxHashRequest struct {
	Hash string `json:"hash"`
}

type BlockHeightRequest struct {
	Height int `json:"height"`
}

type ChainInfoRequest struct {
	Min int `json:"min"`
	Max int `json:"max"`
}