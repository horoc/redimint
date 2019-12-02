package rpc

import (
	context "context"
	"fmt"
	"github.com/chenzhou9513/DecentralizedRedis/logger"
	"github.com/chenzhou9513/DecentralizedRedis/models"
	"github.com/chenzhou9513/DecentralizedRedis/rpc/proto"
	s "github.com/chenzhou9513/DecentralizedRedis/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type RpcServer struct {
	port   string
	server *grpc.Server
	app    *RpcDBService
}

func NewRpcServer(port string) *RpcServer {
	s := &RpcServer{
		server: grpc.NewServer(),
		app:    &RpcDBService{},
		port:   port,
	}
	proto.RegisterDecentralizedRedisServer(s.server, s.app)
	reflection.Register(s.server)
	return s
}

func (s *RpcServer) StartServer() {
	lis, err := net.Listen("tcp", "127.0.0.1:"+s.port)
	if err != nil {
		logger.Error("failed to listen: %v", err)
		return
	}
	if err := s.server.Serve(lis); err != nil {
		logger.Error("failed to serve: %v", err)
		return
	}
}

type RpcDBService struct {
}

func (r RpcDBService) Query(c context.Context, req *proto.CommandRequest) (*proto.QueryResponse, error) {
	fmt.Println("get resquest")
	queryResponse := s.AppService.Query(&models.CommandRequest{Cmd: req.Cmd})
	return &proto.QueryResponse{
		Code:    queryResponse.Code,
		CodeMsg: queryResponse.CodeMsg,
		Result:  queryResponse.Result,
	}, nil
}
