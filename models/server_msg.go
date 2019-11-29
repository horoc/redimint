package models


type ExecutionRequest struct {
	Operation string `json:"operation"`
}

type TxHashRequest struct{
	Hash string `json:"hash"`
}