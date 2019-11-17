package handlers

import (
	"fmt"
	"github.com/chenzhou9513/DecentralizedRedis/benchmark"
	"github.com/chenzhou9513/DecentralizedRedis/logger"
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

func RestoreLocalDatabase(c *gin.Context) {
	ginMsg := models.GinMsg{C: c}
	err := s.AppService.RestoreLocalDatabase()
	if err != nil {
		logger.Error(err)
	}
	ginMsg.Response(http.StatusOK, nil)
}

func BenchMarkTest(c *gin.Context) {
	ginMsg := models.GinMsg{C: c}
	request := &models.BenchMarkRequest{}
	ginMsg.DecodeRequestBody(request)
	fmt.Println(request)
	mark, err := benchmark.NewBenchMark(request)
	if err != nil {
		logger.Info(err)
		return
	}
	test := mark.StartTest()
	ginMsg.Response(http.StatusOK, test)
}
