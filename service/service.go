package service

import "github.com/chenzhou9513/DecentralizedRedis/models"

type Service interface {

	RestoreLocalDatabase() error
	Query(request *models.CommandRequest) *models.QueryResponse
	QueryPrivateDataWithAddress(request *models.CommandRequest, address string) *models.QueryResponse

	Execute(request *models.CommandRequest) *models.ExecuteResponse
	ExecuteAsync(request *models.CommandRequest) *models.ExecuteAsyncResponse
	ExecuteWithPrivateKey(request *models.CommandRequest) *models.ExecuteResponse

	QueryTransaction(hash string) *models.Transaction
	QueryCommittedTxList(beginHeight int, endHeight int) *models.TransactionCommittedList
	QueryBlock(height int) *models.Block

	GetChainState() *models.ChainState
	GetChainInfo(min int, max int) *models.ChainInfo
}
