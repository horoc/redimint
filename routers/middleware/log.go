package middleware

import (
	"bytes"
	"github.com/chenzhou9513/DecentralizedRedis/logger"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"strings"
	"time"
)

func Flutter(str string) string {
	str = strings.Replace(str, " ", "", -1)
	str = strings.Replace(str, "\n", "", -1)
	str = strings.Replace(str, "\t", "", -1)
	return str
}

func Log() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()

		method := c.Request.Method
		path := c.Request.URL.Path

		bodyBytes, _ := ioutil.ReadAll(c.Request.Body)
		c.Request.Body.Close()
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		request := Flutter(string(bodyBytes))
		logger.Log.Infof(" %s  %s | body : %s |",
			method,
			path,
			request,
		)

		c.Next()
		end := time.Now()
		latency := end.Sub(start)

		statusCode := c.Writer.Status()
		logger.Log.Infof(" %s  %s | %3d | %v |",
			method, path,
			statusCode,
			latency,
		)
	}

}
