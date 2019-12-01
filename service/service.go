package service

import "github.com/chenzhou9513/DecentralizedRedis/models"

type Service interface {
	Query(request *models.CommandRequest) *models.QueryResponse
	Execute(request *models.CommandRequest) *models.ExecuteResponse
	ExecuteAsync(request *models.CommandRequest) *models.ExecuteAsyncResponse

	QueryTransaction(hash string) *models.Transaction
	QueryBlock(height int) *models.Block
	GetChainState() *models.ChainState
	GetChainInfo(min int, max int) *models.ChainInfo
}
