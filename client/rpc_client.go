package client

import (
	"context"
	"github.com/chenzhou9513/DecentralizedRedis/rpc/proto"
	"google.golang.org/grpc"
)

type RpcClient struct {
	app proto.DecentralizedRedisClient
}

func NewRpcClient(address string) (*RpcClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &RpcClient{proto.NewDecentralizedRedisClient(conn)}, nil
}

func (r RpcClient) Query(cmd *proto.CommandRequest) (*proto.QueryResponse, error) {
	return r.app.Query(context.Background(), cmd)
}
