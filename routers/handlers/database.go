package handlers

import (
	"fmt"
	"github.com/chenzhou9513/redimint/core"
	"github.com/chenzhou9513/redimint/logger"
	"github.com/chenzhou9513/redimint/models"
	"github.com/chenzhou9513/redimint/models/code"
	"github.com/chenzhou9513/redimint/utils/bench_mark"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func ExecuteCommand(c *gin.Context) {
	ginMsg := models.GinMsg{C: c}
	request := &models.ExecuteRequest{}
	ginMsg.DecodeRequestBody(request)
	var err error
	var res interface{}

	if strings.EqualFold(request.Mode, "async") {
		res, err = core.AppService.ExecuteAsync(&models.CommandRequest{request.Cmd})
	} else if strings.EqualFold(request.Mode, "commit") {
		res, err = core.AppService.Execute(&models.CommandRequest{request.Cmd})
	} else if strings.EqualFold(request.Mode, "private") {
		res, err = core.AppService.ExecuteWithPrivateKey(&models.CommandRequest{request.Cmd})
	} else {
		ginMsg.CommonResponse(http.StatusOK, code.CodeTypeInvalidExecuteMode, fmt.Sprintf("Invalid mode : %s", request.Mode))
	}

	if err != nil {
		ginMsg.Error(http.StatusOK, code.CodeTypeRedimintExecuteError, code.CodeTypeRedimintExecuteErrorMsg, err)
	}
	ginMsg.SuccessWithData(res)
}

func QueryCommand(c *gin.Context) {
	ginMsg := models.GinMsg{C: c}
	request := &models.CommandRequest{}
	ginMsg.DecodeRequestBody(request)
	res, err := core.AppService.Query(&models.CommandRequest{request.Cmd})
	if err != nil {
		ginMsg.Error(http.StatusOK, code.CodeTypeRedimintQueryError, code.CodeTypeRedimintQueryErrorMsg, err)
	}
	ginMsg.SuccessWithData(res)
}

func QueryPrivateCommand(c *gin.Context) {
	ginMsg := models.GinMsg{C: c}
	request := &models.CommandRequest{}
	ginMsg.DecodeRequestBody(request)
	addr := c.Query("address")
	var res *models.QueryResponse
	var err error
	if len(addr) != 0 {
		res, err = core.AppService.QueryPrivateDataWithAddress(&models.QueryPrivateWithAddrRequest{request.Cmd, strings.ToUpper(addr)})
	} else {
		res, err = core.AppService.QueryPrivateData(&models.CommandRequest{request.Cmd})
	}
	if err != nil {
		ginMsg.Error(http.StatusOK, code.CodeTypeRedimintQueryError, code.CodeTypeRedimintQueryErrorMsg, err)
	}
	ginMsg.Response(http.StatusOK, res)
}

func RestoreLocalDatabase(c *gin.Context) {
	ginMsg := models.GinMsg{C: c}
	err := core.AppService.RestoreLocalDatabase()
	if err != nil {
		ginMsg.Error(http.StatusOK, code.CodeTypeInternalError, code.CodeTypeInternalErrorMsg, err)
	}
	ginMsg.Success()
}

func BenchMarkTest(c *gin.Context) {
	ginMsg := models.GinMsg{C: c}
	request := &models.BenchMarkRequest{}
	ginMsg.DecodeRequestBody(request)
	fmt.Println(request)
	mark, err := bench_mark.NewBenchMark(request)
	if err != nil {
		logger.Log.Error(err)
		return
	}
	test := mark.StartTest()
	ginMsg.Response(http.StatusOK, test)
}
