package handlers

import (
	"fmt"
	"github.com/chenzhou9513/DecentralizedRedis/benchmark"
	"github.com/chenzhou9513/DecentralizedRedis/core"
	"github.com/chenzhou9513/DecentralizedRedis/logger"
	"github.com/chenzhou9513/DecentralizedRedis/models"
	"github.com/chenzhou9513/DecentralizedRedis/models/code"
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
		res := core.AppService.ExecuteAsync(&models.CommandRequest{request.Cmd})
		ginMsg.Response(http.StatusOK, res)
	} else if strings.EqualFold(request.Mode, "commit") {
		res := core.AppService.Execute(&models.CommandRequest{request.Cmd})
		ginMsg.Response(http.StatusOK, res)
	} else if strings.EqualFold(request.Mode, "private") {
		res := core.AppService.ExecuteWithPrivateKey(&models.CommandRequest{request.Cmd})
		ginMsg.Response(http.StatusOK, res)
	} else {
		ginMsg.ErrorResponse(http.StatusOK, code.CodeTypeInvalidExecuteMode, fmt.Sprintf("Invalid mode : %s", request.Mode))
	}
}

func QueryCommand(c *gin.Context) {
	ginMsg := models.GinMsg{C: c}
	request := &models.CommandRequest{}
	ginMsg.DecodeRequestBody(request)
	res := core.AppService.Query(&models.CommandRequest{request.Cmd})
	ginMsg.Response(http.StatusOK, res)
}

func QueryPrivateCommand(c *gin.Context) {
	ginMsg := models.GinMsg{C: c}
	request := &models.CommandRequest{}
	ginMsg.DecodeRequestBody(request)
	addr := c.Query("address")
	var res *models.QueryResponse
	if len(addr) != 0 {
		res = core.AppService.QueryPrivateDataWithAddress(&models.CommandRequest{request.Cmd}, strings.ToUpper(addr))
	} else {
		res = core.AppService.QueryPrivateDataWithAddress(&models.CommandRequest{request.Cmd}, utils.ValidatorKey.Address.String())
	}
	ginMsg.Response(http.StatusOK, res)
}

func RestoreLocalDatabase(c *gin.Context) {
	ginMsg := models.GinMsg{C: c}
	err := core.AppService.RestoreLocalDatabase()
	if err != nil {
		logger.Log.Error(err)
		ginMsg.Response(http.StatusOK, gin.H{
			"code": code.CodeTypeInternalError,
			"msg":  err,
		})
	}
	ginMsg.Response(http.StatusOK, gin.H{
		"code": code.CodeTypeOK,
		"msg":  code.Info(code.CodeTypeOK),
	})
}

func BenchMarkTest(c *gin.Context) {
	ginMsg := models.GinMsg{C: c}
	request := &models.BenchMarkRequest{}
	ginMsg.DecodeRequestBody(request)
	fmt.Println(request)
	mark, err := benchmark.NewBenchMark(request)
	if err != nil {
		logger.Log.Error(err)
		return
	}
	test := mark.StartTest()
	ginMsg.Response(http.StatusOK, test)
}
