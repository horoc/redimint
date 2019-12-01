package network

import (
	"fmt"
	"github.com/chenzhou9513/DecentralizedRedis/logger"
	"github.com/chenzhou9513/DecentralizedRedis/routers"
	"github.com/chenzhou9513/DecentralizedRedis/service"
	"github.com/chenzhou9513/DecentralizedRedis/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

var AppServer *Server

type Server struct {
	host       string
	port       string
	httpServer *http.Server
	appService service.Service
}

func NewServer(host string, port string) *Server {

	gin.SetMode(utils.Config.Server.RunMode)

	routersInit := routers.InitRouter()
	endPoint := fmt.Sprintf(":%d", utils.Config.Server.Port)

	httpServer := &http.Server{
		Addr:    endPoint,
		Handler: routersInit,
	}

	logger.Info("[info] start http server listening ", endPoint)
	server := &Server{host, port, httpServer, service.AppService}
	AppServer = server
	//server.setRoute()
	return server
}

func (server *Server) Start() {
	fmt.Printf("Server will be started at %s:%s...\n", server.host, server.port)
	if err := server.httpServer.ListenAndServe(); err != nil {
		logger.Error(err)
		return
	}
}
