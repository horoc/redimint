package models

type TxCommitBody struct {
	Data      *TxCommitData `json:"data"`
	Signature string `json:"signature"`
	Address   string `json:"address"`
}

type TxCommitData struct {
	Operation string `json:"operation"`
	Sequence  string `json:"sequence"`
}

type TxValidatorUpdate struct {
	PublicKey string `json:"public_key"'`
	Power     string `json:"power"`
}
