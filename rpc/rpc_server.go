package rpc

import (
	context "context"
	"fmt"
	"github.com/chenzhou9513/DecentralizedRedis/core"
	"github.com/chenzhou9513/DecentralizedRedis/logger"
	"github.com/chenzhou9513/DecentralizedRedis/models"
	"github.com/chenzhou9513/DecentralizedRedis/rpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type Server struct {
	port   string
	server *grpc.Server
	app    *DBService
}

func NewRpcServer(port string) *Server {
	s := &Server{
		server: grpc.NewServer(),
		app:    &DBService{},
		port:   port,
	}
	proto.RegisterDecentralizedRedisServer(s.server, s.app)
	reflection.Register(s.server)
	return s
}

func (s *Server) StartServer() {
	lis, err := net.Listen("tcp", "127.0.0.1:"+s.port)
	if err != nil {
		logger.Log.Error("failed to listen: %v", err)
		return
	}
	if err := s.server.Serve(lis); err != nil {
		logger.Log.Error("failed to serve: %v", err)
		return
	}

}

type DBService struct {
}

func (r DBService) Query(c context.Context, req *proto.CommandRequest) (*proto.QueryResponse, error) {
	fmt.Println("get resquest")
	queryResponse := core.AppService.Query(&models.CommandRequest{Cmd: req.Cmd})
	return &proto.QueryResponse{
		Code:    queryResponse.Code,
		CodeMsg: queryResponse.CodeMsg,
		Result:  queryResponse.Result,
	}, nil
}
