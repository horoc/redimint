package client

import (
	"context"
	"github.com/chenzhou9513/redimint/rpc/proto"
	"google.golang.org/grpc"
)

type RpcClient struct {
	app proto.redimintClient
}

func NewRpcClient(address string) (*RpcClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &RpcClient{proto.NewredimintClient(conn)}, nil
}

func (r RpcClient) Query(cmd *proto.CommandRequest) (*proto.QueryResponse, error) {
	return r.app.Query(context.Background(), cmd)
}
