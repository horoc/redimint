package handlers

import (
	"github.com/chenzhou9513/DecentralizedRedis/models"
	s "github.com/chenzhou9513/DecentralizedRedis/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ExecuteCommand(c *gin.Context) {
	ginMsg := models.GinMsg{C: c}
	request := &models.CommandRequest{}
	ginMsg.DecodeRequestBody(request)
	res := s.AppService.Execute(&models.CommandRequest{request.Cmd})
	ginMsg.Response(http.StatusOK, res)
}

func ExecuteCommandAsync(c *gin.Context) {
	ginMsg := models.GinMsg{C: c}
	request := &models.CommandRequest{}
	ginMsg.DecodeRequestBody(request)
	res := s.AppService.ExecuteAsync(&models.CommandRequest{request.Cmd})
	ginMsg.Response(http.StatusOK, res)
}

func QueryCommand(c *gin.Context) {
	ginMsg := models.GinMsg{C: c}
	request := &models.CommandRequest{}
	ginMsg.DecodeRequestBody(request)
	res := s.AppService.Query(&models.CommandRequest{request.Cmd})
	ginMsg.Response(http.StatusOK, res)
}
