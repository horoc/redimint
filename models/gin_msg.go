package models

import (
	"github.com/chenzhou9513/DecentralizedRedis/utils"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

type GinMsg struct {
	C *gin.Context
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (g *GinMsg) Response(httpCode int, data interface{}) {
	g.C.JSON(httpCode, data)
	//g.C.JSON(httpCode, Response{
	//	Code: errCode,
	//	Msg:  code.Info(errCode),
	//	Data: data,
	//})
	return
}

func (g *GinMsg) DecodeRequestBody(data interface{}) {
	body, _ := ioutil.ReadAll(g.C.Request.Body)
	utils.JsonToStruct(body, data)
}
