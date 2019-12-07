package handlers

import (
	"fmt"
	"github.com/chenzhou9513/DecentralizedRedis/benchmark"
	"github.com/chenzhou9513/DecentralizedRedis/logger"
	"github.com/chenzhou9513/DecentralizedRedis/models"
	"github.com/chenzhou9513/DecentralizedRedis/models/code"
	s "github.com/chenzhou9513/DecentralizedRedis/service"
	"github.com/chenzhou9513/DecentralizedRedis/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func ExecuteCommand(c *gin.Context) {
	ginMsg := models.GinMsg{C: c}
	request := &models.ExecuteRequest{}
	ginMsg.DecodeRequestBody(request)
	if strings.EqualFold(request.Mode, "async") {
		res := s.AppService.ExecuteAsync(&models.CommandRequest{request.Cmd})
		ginMsg.Response(http.StatusOK, res)
	} else if strings.EqualFold(request.Mode, "commit") {
		res := s.AppService.Execute(&models.CommandRequest{request.Cmd})
		ginMsg.Response(http.StatusOK, res)
	} else if strings.EqualFold(request.Mode, "private") {
		res := s.AppService.ExecuteWithPrivateKey(&models.CommandRequest{request.Cmd})
		ginMsg.Response(http.StatusOK, res)
	} else {
		ginMsg.ErrorResponse(http.StatusOK, code.CodeTypeInvalidExecuteMode, fmt.Sprintf("Invalid mode : %s", request.Mode))
	}
}

func QueryCommand(c *gin.Context) {
	ginMsg := models.GinMsg{C: c}
	request := &models.CommandRequest{}
	ginMsg.DecodeRequestBody(request)
	res := s.AppService.Query(&models.CommandRequest{request.Cmd})
	ginMsg.Response(http.StatusOK, res)
}

func QueryPrivateCommand(c *gin.Context) {
	ginMsg := models.GinMsg{C: c}
	request := &models.CommandRequest{}
	ginMsg.DecodeRequestBody(request)
	addr := c.Query("address")
	var res *models.QueryResponse
	if len(addr) != 0 {
		res = s.AppService.QueryPrivateDataWithAddress(&models.CommandRequest{request.Cmd}, strings.ToUpper(addr))
	} else {
		res = s.AppService.QueryPrivateDataWithAddress(&models.CommandRequest{request.Cmd}, utils.ValidatorKey.Address.String())
	}
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
