package handlers

import (
	"github.com/chenzhou9513/DecentralizedRedis/logger"
	"github.com/chenzhou9513/DecentralizedRedis/models"
	"github.com/chenzhou9513/DecentralizedRedis/models/code"
	"github.com/chenzhou9513/DecentralizedRedis/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Login(c *gin.Context) {
	ginMsg := models.GinMsg{C: c}
	request := &models.LoginRequest{}
	ginMsg.DecodeRequestBody(request)

	if request.Name != utils.Config.App.Name || request.Password != utils.Config.App.Auth {
		c.JSON(http.StatusOK, gin.H{
			"code": code.CodeTypeIncorrectPassword,
			"msg":  "Incorrect Password",
		})
	}
	s, err := utils.GenerateToken(c.ClientIP())
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"code": code.CodeTypeInternalError,
			"msg":  err,
		})
	}
	ginMsg.Response(http.StatusOK, &models.TokenResponse{
		Code:    code.CodeTypeOK,
		CodeMsg: code.Info(code.CodeTypeOK),
		Token:   s,
	})
}
