package handlers

import (
	"github.com/chenzhou9513/DecentralizedRedis/models"
	s "github.com/chenzhou9513/DecentralizedRedis/service"
	"github.com/gin-gonic/gin"

	"net/http"
)

func GetTransactionByHash(c *gin.Context) {
	ginMsg := models.GinMsg{C: c}
	request := &models.TxHashRequest{}
	ginMsg.DecodeRequestBody(request)
	res := s.AppService.QueryTransaction(request.Hash)
	ginMsg.Response(http.StatusOK, res)
}

func GetBlockByHeight(c *gin.Context) {
	ginMsg := models.GinMsg{C: c}
	request := &models.BlockHeightRequest{}
	ginMsg.DecodeRequestBody(request)
	res := s.AppService.QueryBlock(request.Height)
	ginMsg.Response(http.StatusOK, res)
}

func GetChainState(c *gin.Context) {
	ginMsg := models.GinMsg{C: c}
	res := s.AppService.GetsChainState()
	ginMsg.Response(http.StatusOK, res)
}
