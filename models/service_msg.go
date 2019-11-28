package models

type ExecuteResponse struct {
	Code          uint32    `json:"code"`
	CodeMsg       string `json:"code_info"`
	Cmd           string `json:"command"`
	ExecuteResult string `json:"execute_result"`
	Signature     string `json:"signature"`
	Sequence      string `json:"sequence"`
	TimeStamp     string `json:"time_stamp"`
	Hash          []byte `json:"hash"`
	Height        int64  `json:"height"`
}

type ExecuteRequest struct {
	Cmd string `json:"command"`
}
