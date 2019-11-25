package consensus

import (
	"fmt"
	"github.com/chenzhou9513/DecentralizedRedis/utils"
	c "github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/types"
)

var tendermintHttpClient *c.HTTP

func initClient() {
	var host = "http://" + utils.Config.Tendermint.Url
	var wsEndpoint = "/websocket"
	tendermintHttpClient = c.NewHTTP(host, wsEndpoint)
}

func BroadcastTxCommit(op string) (*ctypes.ResultBroadcastTxCommit) {

	if tendermintHttpClient == nil {
		initClient()
	}
	tx := types.Tx(op)
	fmt.Println("hash")
	fmt.Println(fmt.Sprintf("%x", tx.Hash()))

	resultBroadcastTxCommit, e := tendermintHttpClient.BroadcastTxCommit(tx)
	if e!=nil{
		fmt.Println(e)
	}
	return resultBroadcastTxCommit
}

func ABCIDataQuery(path string,  data []byte) *ctypes.ResultABCIQuery{
	if tendermintHttpClient == nil {
		initClient()
	}

	resultABCIQuery, e := tendermintHttpClient.ABCIQuery(path, data)
	if e != nil {
		fmt.Println(e)
	}
	return resultABCIQuery
}

func SearchTx(query string, page int, size int) *ctypes.ResultTxSearch {
	if tendermintHttpClient == nil {
		initClient()
	}

	resultTx, e := tendermintHttpClient.TxSearch(query, true, page, size)
	if e != nil {
		fmt.Println(e)
	}
	return resultTx
}

func GetTx(hash []byte) *ctypes.ResultTx {
	if tendermintHttpClient == nil {
		initClient()
	}

	resultTx, e := tendermintHttpClient.Tx(hash, true)
	if e != nil {
		fmt.Println(e)
	}
	return resultTx
}

func GetChainInfo(min int, max int) *ctypes.ResultBlockchainInfo{
	if tendermintHttpClient == nil {
		initClient()
	}
	minH := int64(min)
	maxH := int64(max)

	resultBlockchainInfo, e := tendermintHttpClient.BlockchainInfo(minH, maxH)
	if e != nil {
		fmt.Println(e)
	}
	return resultBlockchainInfo
}

func GetChainState() *ctypes.ResultStatus{
	if tendermintHttpClient == nil {
		initClient()
	}
	resultStatus, e := tendermintHttpClient.Status()
	if e != nil {
		fmt.Println(e)
	}
	return resultStatus
}


func GetBlockFromHeight(h int) *ctypes.ResultBlock {

	if tendermintHttpClient == nil {
		initClient()
	}
	height := int64(h)
	resultBlock, e := tendermintHttpClient.Block(&height)
	if e != nil {
		fmt.Println(e)
	}

	//str := "http://" + utils.Config.Tendermint.Url + "/block"
	//u, _ := url.Parse(str)
	//q, _ := url.ParseQuery(u.RawQuery)
	//q.Add("height", h)
	//u.RawQuery = q.Encode()
	//request, e := http.NewRequest("GET", fmt.Sprint(u), nil)
	//if e != nil {
	//	fmt.Println(e)
	//}
	//response := utils.SendRequest(request)
	//res := new(ctypes.ResultBlock)
	//bodyBytes, e := ioutil.ReadAll(response.Body)
	//if response.Body != nil{
	//	response.Body.Close()
	//}
	//if e != nil {
	//	fmt.Println(e)
	//}
	//e = json.Unmarshal(bodyBytes, res)
	//if e!=nil{
	//	fmt.Println(e)
	//}
	//fmt.Println(res)
	return resultBlock
}
