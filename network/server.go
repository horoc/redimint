package network

import (
	"fmt"
	"github.com/chenzhou9513/DecentralizedRedis/logger"
	"github.com/chenzhou9513/DecentralizedRedis/routers"
	"github.com/chenzhou9513/DecentralizedRedis/rpc"
	"github.com/chenzhou9513/DecentralizedRedis/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"net/http"
)

var AppServer *Server

type Server struct {
	httpPort   string
	rpcPort    string
	rpcServer  *rpc.Server
	httpServer *http.Server
}

func NewServer() *Server {

	server := &Server{
		httpPort:   strconv.Itoa(utils.Config.HttpServer.Port),
		rpcPort:    strconv.Itoa(utils.Config.Rpc.Port),
		rpcServer:  nil,
		httpServer: nil,
	}

	gin.SetMode(utils.Config.HttpServer.RunMode)
	routersInit := routers.InitRouter()
	endPoint := fmt.Sprintf(":%d", utils.Config.HttpServer.Port)

	server.httpServer = &http.Server{
		Addr:    endPoint,
		Handler: routersInit,
	}

	server.rpcServer = rpc.NewRpcServer(strconv.Itoa(utils.Config.Rpc.Port))
	AppServer = server

	return server
}

func (server *Server) Start() {

	fmt.Printf("Rpc Server will be started at :%s...\n", server.rpcPort)
	go server.rpcServer.StartServer()

	fmt.Printf("Http Server will be started at :%s...\n", server.httpPort)
	if err := server.httpServer.ListenAndServe(); err != nil {
		logger.Error(err)
		return
	}
}

