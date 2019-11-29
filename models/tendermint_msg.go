package models

type TxCommitBody struct {
	Operation string `json:"operation"`
	Sequence string `json:"sequence"`
	Signature string `json:"signature"`
	Address string `json:"address"`
}

type TxValidatorUpdate struct {
	PublicKey string `json:"public_key"'`
	Power string `json:"power"`
}
