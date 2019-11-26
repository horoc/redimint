package network

type ExecutionRequest struct {
	Operation string `json:"operation"`
}

type QueryResponse struct {
	Operation string `json:"operation"`
	Result string `json:"result"`
}

type QueryBlockRequest struct {
	Key string `json:"key"`
}

type QueryLogResponse struct {
	Height int
	Logs   []string
}

type TpsTest struct {
	TotalTx   int
	TotalTime string
	Details   []TestDetail
}

type TestDetail struct {
	Time     string
	Response string
}
