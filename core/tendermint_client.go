package core

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


func InitClient() {
	var host = "tcp://" + utils.Config.Tendermint.Url
	var wsEndpoint = "./websocket"
	tendermintHttpClient = c.NewHTTP(host, wsEndpoint)
}

func BroadcastTxCommit(op *models.TxCommitBody) (*ctypes.ResultBroadcastTxCommit) {

	tx := types.Tx(utils.StructToJson(op))
	resultBroadcastTxCommit, err := tendermintHttpClient.BroadcastTxCommit(tx)
	if err != nil {
		logger.Error(err)
		return nil
	}
	return resultBroadcastTxCommit
}

func BroadcastTxSync(op *models.TxCommitBody) (*ctypes.ResultBroadcastTx) {

	tx := types.Tx(utils.StructToJson(op))
	resultBroadcastTxCommit, err := tendermintHttpClient.BroadcastTxSync(tx)
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

	resultABCIQuery, err := tendermintHttpClient.ABCIQuery(path, data)
	if err != nil {
		logger.Error(err)
		return nil
	}
	return resultABCIQuery
}

func SearchTx(query string, page int, size int) *ctypes.ResultTxSearch {

	resultTx, err := tendermintHttpClient.TxSearch(query, true, page, size)
	if err != nil {
		logger.Error(err)
		return nil
	}
	return resultTx
}

func GetTx(hash []byte) *ctypes.ResultTx {

	resultTx, err := tendermintHttpClient.Tx(hash, true)
	if err != nil {
		logger.Error(err)
		return nil
	}
	return resultTx
}

func GetChainInfo(min int, max int) *ctypes.ResultBlockchainInfo {

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

	resultStatus, err := tendermintHttpClient.Status()
	if err != nil {
		logger.Error(err)
		return nil
	}
	return resultStatus
}

func GetBlockFromHeight(h int) *ctypes.ResultBlock {

	height := int64(h)
	resultBlock, err := tendermintHttpClient.Block(&height)
	if err != nil {
		logger.Error(err)
		return nil
	}
	return resultBlock
}
