package middleware

import (
	"github.com/chenzhou9513/DecentralizedRedis/models/code"
	"github.com/chenzhou9513/DecentralizedRedis/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		respCode := code.CodeTypeOK
		token := c.GetHeader("token")
		address := c.ClientIP()
		var claims *utils.Claims
		var err error
		if token == "" {
			respCode = code.CodeTypePermissionDenied
		} else {
			claims, err = utils.ParseToken(token)
			if err != nil {
				switch err.(*jwt.ValidationError).Errors {
				case jwt.ValidationErrorExpired:
					respCode = code.CodeTypeTokenTimeoutError
				default:
					respCode = code.CodeTypeTokenInvalidError
				}
			}
		}

		if respCode != code.CodeTypeOK || strings.EqualFold(claims.Address, address) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": respCode,
				"msg":  "Invalid Token",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
