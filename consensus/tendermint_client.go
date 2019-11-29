package consensus

import (
	"fmt"
	"github.com/chenzhou9513/DecentralizedRedis/logger"
	"github.com/chenzhou9513/DecentralizedRedis/models"
	"github.com/chenzhou9513/DecentralizedRedis/utils"
	c "github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/types"
	"io/ioutil"
	"net/http"
	"net/url"
)

var tendermintHttpClient *c.HTTP

func initClient() {
	var host = "tcp://" + utils.Config.Tendermint.Url
	var wsEndpoint = "./websocket"
	tendermintHttpClient = c.NewHTTP(host, wsEndpoint)
}

func TestTendermintTime(tx []byte) (*ctypes.ResultBroadcastTxCommit){
	if tendermintHttpClient == nil {
		initClient()
	}
	resultBroadcastTxCommit, _ := tendermintHttpClient.BroadcastTxCommit(tx)
	return resultBroadcastTxCommit
}

func BroadcastTxCommit(op *models.TxCommitBody) (*ctypes.ResultBroadcastTxCommit) {

	if tendermintHttpClient == nil {
		initClient()
	}
	tx := types.Tx(utils.StructToJson(op))
	resultBroadcastTxCommit, err := tendermintHttpClient.BroadcastTxCommit(tx)
	if err != nil {
		logger.Error(err)
		return nil
	}
	return resultBroadcastTxCommit
}

func BroadcastTxCommitUseHttp(op *models.TxCommitBody) (*ctypes.ResultBroadcastTxCommit) {

	str := "http://" + utils.Config.Tendermint.Url + "/broadcast_tx_commit"
	u, _ := url.Parse(str)
	q, _ := url.ParseQuery(u.RawQuery)

	json := utils.StructToJson(op)
	hex := utils.ByteToHex(json)

	q.Add("tx", "\""+hex+"\"")

	u.RawQuery = q.Encode()
	req, _ := http.NewRequest("GET", fmt.Sprint(u), nil)
	res := utils.SendRequest(req)

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Error(err)
	}

	var obj map[string]interface{}
	utils.JsonToStruct(bytes, &obj)
	var result = &ctypes.ResultBroadcastTxCommit{}
	utils.JsonToStruct(utils.StructToJson(obj["result"]), result)

	if err != nil {
		logger.Error(err)
		return nil
	}
	return result
}


func ABCIDataQuery(path string, data []byte) *ctypes.ResultABCIQuery {
	if tendermintHttpClient == nil {
		initClient()
	}

	resultABCIQuery, err := tendermintHttpClient.ABCIQuery(path, data)
	if err != nil {
		logger.Error(err)
		return nil
	}
	return resultABCIQuery
}

func SearchTx(query string, page int, size int) *ctypes.ResultTxSearch {
	if tendermintHttpClient == nil {
		initClient()
	}

	resultTx, err := tendermintHttpClient.TxSearch(query, true, page, size)
	if err != nil {
		logger.Error(err)
		return nil
	}
	return resultTx
}

func GetTx(hash []byte) *ctypes.ResultTx {
	if tendermintHttpClient == nil {
		initClient()
	}

	resultTx, err := tendermintHttpClient.Tx(hash, true)
	if err != nil {
		logger.Error(err)
		return nil
	}
	return resultTx
}

func GetChainInfo(min int, max int) *ctypes.ResultBlockchainInfo {
	if tendermintHttpClient == nil {
		initClient()
	}
	minH := int64(min)
	maxH := int64(max)

	resultBlockchainInfo, err := tendermintHttpClient.BlockchainInfo(minH, maxH)
	if err != nil {
		logger.Error(err)
		return nil

	}
	return resultBlockchainInfo
}

func GetChainState() *ctypes.ResultStatus {
	if tendermintHttpClient == nil {
		initClient()
	}
	resultStatus, err := tendermintHttpClient.Status()
	if err != nil {
		logger.Error(err)
		return nil
	}
	return resultStatus
}

func GetBlockFromHeight(h int) *ctypes.ResultBlock {

	if tendermintHttpClient == nil {
		initClient()
	}
	height := int64(h)
	resultBlock, err := tendermintHttpClient.Block(&height)
	if err != nil {
		logger.Error(err)
		return nil
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
