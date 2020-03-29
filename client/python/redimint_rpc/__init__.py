from .redimint_pb2 import (
    AuthRequest,
    Token,
    CommandRequest,
    QueryResponse,
    QueryPrivateWithAddrRequest,
    ExecuteResponse,
    ExecuteAsyncResponse
)

from .redimint_pb2_grpc import (
    RedimintStub
)

import grpc


def get_client(server):
    conn = grpc.insecure_channel(server)
    return RedimintStub(channel=conn)
