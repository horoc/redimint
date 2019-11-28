package service

import "github.com/chenzhou9513/DecentralizedRedis/models"

type Service interface {

	Execute(request *models.ExecuteRequest) *models.ExecuteResponse

}