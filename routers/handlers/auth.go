package handlers

import (
	"github.com/chenzhou9513/redimint/logger"
	"github.com/chenzhou9513/redimint/models"
	"github.com/chenzhou9513/redimint/models/code"
	"github.com/chenzhou9513/redimint/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Login(c *gin.Context) {
	ginMsg := models.GinMsg{C: c}
	request := &models.LoginRequest{}
	ginMsg.DecodeRequestBody(request)

	if request.Name != utils.Config.App.AdminUser || request.Password != utils.Config.App.AdminPassword {
		c.JSON(http.StatusOK, gin.H{
			"code": code.CodeTypeIncorrectPassword,
			"msg":  "Incorrect Password",
		})
	}
	s, err := utils.GenerateToken(c.ClientIP(), "admin", request.Name, request.Password)
	if err != nil {
		logger.Log.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"code": code.CodeTypeInternalError,
			"msg":  err,
		})
	}
	ginMsg.Response(http.StatusOK, &models.TokenResponse{
		Code:    code.CodeTypeOK,
		CodeMsg: code.CodeTypeOKMsg,
		Token:   s,
	})
}
