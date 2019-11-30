package service

import "github.com/chenzhou9513/DecentralizedRedis/models"

type Service interface {
	Execute(request *models.ExecuteRequest) *models.ExecuteResponse
	ExecuteAsync(request *models.ExecuteRequest) *models.ExecuteAsyncResponse

	QueryTransaction(hash string) *models.Transaction
	QueryBlock(height int) *models.Block
}
