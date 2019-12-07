package core

import (
	"fmt"
	"github.com/chenzhou9513/DecentralizedRedis/logger"
	"github.com/chenzhou9513/DecentralizedRedis/models"
	"github.com/chenzhou9513/DecentralizedRedis/utils"
	c "github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/types"
)

var tendermintHttpClient *c.HTTP

func InitClient() {
	var host = "tcp://" + utils.Config.Tendermint.Url
	var wsEndpoint = "./websocket"
	tendermintHttpClient = c.NewHTTP(host, wsEndpoint)
}

func BroadcastTxCommit(op *models.TxCommitBody) (*ctypes.ResultBroadcastTxCommit, error) {

	tx := types.Tx(utils.StructToJson(op))
	resultBroadcastTxCommit, err := tendermintHttpClient.BroadcastTxCommit(tx)
	if err != nil {
		err = fmt.Errorf("BroadcastTxCommit command error : %s, %s", tx, err)
		logger.Error(err)
		return nil, err
	}
	return resultBroadcastTxCommit, nil
}

func BroadcastTxSync(op *models.TxCommitBody) (*ctypes.ResultBroadcastTx, error) {

	tx := types.Tx(utils.StructToJson(op))
	resultBroadcastTxCommit, err := tendermintHttpClient.BroadcastTxSync(tx)
	if err != nil {
		err = fmt.Errorf("BroadcastTxSync command error : %s, %s", tx, err)
		logger.Error(err)
		return nil, err
	}
	return resultBroadcastTxCommit, nil
}

func ABCIDataQuery(path string, data []byte) *ctypes.ResultABCIQuery {

	resultABCIQuery, err := tendermintHttpClient.ABCIQuery(path, data)
	if err != nil {
		logger.Error(err)
		return nil
	}
	return resultABCIQuery
}

func GetTx(hash []byte) (*ctypes.ResultTx, error) {

	resultTx, err := tendermintHttpClient.Tx(hash, true)
	if err != nil {
		err = fmt.Errorf("get transaction by hash error : %s, %s", utils.ByteToHex(hash), err)
		logger.Error(err)
		return nil, err
	}
	return resultTx, nil
}

func GetChainInfo(min int, max int) (*ctypes.ResultBlockchainInfo, error) {

	minH := int64(min)
	maxH := int64(max)

	resultBlockchainInfo, err := tendermintHttpClient.BlockchainInfo(minH, maxH)
	if err != nil {
		err = fmt.Errorf("get chain info error : %s", err)
		logger.Error(err)
		return nil, err
	}
	return resultBlockchainInfo, nil
}

func GetChainState() (*ctypes.ResultStatus, error) {

	resultStatus, err := tendermintHttpClient.Status()
	if err != nil {
		err = fmt.Errorf("get chain state error : %s", err)
		logger.Error(err)
		return nil, err
	}
	return resultStatus, nil
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
