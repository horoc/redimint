package client

import (
	"context"
	proto "github.com/chenzhou9513/redimint/grpc/proto"
	"google.golang.org/grpc"
)

type RpcClient struct {
	app proto.RedimintClient
}

func NewRpcClient(address string) (*RpcClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &RpcClient{proto.NewRedimintClient(conn)}, nil
}

func (r RpcClient) Query(cmd *proto.CommandRequest) (*proto.QueryResponse, error) {
	return r.app.Query(context.Background(), cmd)
}
