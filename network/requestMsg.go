package network

type ExecutionRequest struct {
	Operation  string `json:"operation"`
}


type QueryBlockRequest struct {
	Key  string `json:"key"`
}
