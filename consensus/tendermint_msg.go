package consensus

import (
	"fmt"
	"github.com/chenzhou9513/DecentralizedRedis/logger"
	"github.com/chenzhou9513/DecentralizedRedis/utils"
	"net/http"
	"net/url"
	"os/exec"
)


type CommitBody struct {
	Operation string `json:"operation"`
	Sequence string `json:"sequence"`
	Signature string `json:"signature"`
}


func GetBroadcastTxCommitRequest(operation string) (*http.Request, error) {
	str := "http://" + utils.Config.Tendermint.Url + "/broadcast_tx_commit"
	u, _ := url.Parse(str)
	q, _ := url.ParseQuery(u.RawQuery)
	q.Add("tx", "\""+operation+"\"")
	u.RawQuery = q.Encode()
	request, err := http.NewRequest("GET", fmt.Sprint(u), nil)

	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return request, nil
}

func GetChainHeightBlockRequest(height string) (*http.Request, error) {
	str := "http://" + utils.Config.Tendermint.Url + "/block"
	u, _ := url.Parse(str)
	q, _ := url.ParseQuery(u.RawQuery)
	q.Add("height", height)
	u.RawQuery = q.Encode()
	request, err := http.NewRequest("GET", fmt.Sprint(u), nil)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return request, nil
}

func GetQueryRequest(operation string) (*http.Request, error) {
	str := "http://" + utils.Config.Tendermint.Url + "/abci_query"
	u, _ := url.Parse(str)
	q, _ := url.ParseQuery(u.RawQuery)
	q.Add("data", "\""+operation+"\"")
	u.RawQuery = q.Encode()
	request, err := http.NewRequest("GET", fmt.Sprint(u), nil)

	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return request, nil
}

func SendBroadcastTxCommitRequest(operation string) (string, error) {
	bytes, err := exec.Command("curl", "-s", "localhost:26657/broadcast_tx_commit?tx=\""+operation+"\"").Output()
	if err != nil {
		logger.Error(err)
		return "", err
	}
	return string(bytes), nil
}
