package rpc

import (
	context "context"
	"github.com/chenzhou9513/redimint/core"
	"github.com/chenzhou9513/redimint/logger"
	"github.com/chenzhou9513/redimint/models"
	proto "github.com/chenzhou9513/redimint/rpc/proto"

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
	proto.RegisterRedimintServer(s.server, s.app)
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
	queryResponse, err := core.AppService.Query(&models.CommandRequest{Cmd: req.Cmd})
	if err != nil {
		return nil, err
	}
	return &proto.QueryResponse{
		Result: queryResponse.Result,
	}, nil
}
