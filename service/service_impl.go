package service

import (
	"encoding/binary"
	"fmt"
	"github.com/chenzhou9513/DecentralizedRedis/consensus"
	"github.com/chenzhou9513/DecentralizedRedis/models"
	"github.com/chenzhou9513/DecentralizedRedis/models/code"
	"github.com/chenzhou9513/DecentralizedRedis/utils"
	uuid "github.com/satori/go.uuid"
	"time"
)

type ServiceImpl struct {
}

func (s ServiceImpl) Execute(request *models.ExecuteRequest) *models.ExecuteResponse {
	op := &models.TxCommitBody{}


	//Sequence
	u := uuid.NewV4()
	u1 := binary.BigEndian.Uint64(u[0:8])
	u2 := binary.BigEndian.Uint64(u[8:16])
	op.Sequence = fmt.Sprintf("%x%x", u1, u2)

	//Signature
	op.Signature = utils.ValidatorStringSign([]byte(request.Cmd))

	//cmd
	op.Operation = request.Cmd

	//address
	op.Address = utils.ValidatorKey.Address.String()

	//timestamp
	timestamp := time.Now().UnixNano() / 1e6

	//tendermint response
	commitMsg := consensus.BroadcastTxCommitUseHttp(op)


	//TODO 错误处理

	return &models.ExecuteResponse{code.CodeTypeOK,
		code.Info(code.CodeTypeOK),
		request.Cmd,
		string(commitMsg.DeliverTx.Data),
		op.Signature,
		op.Sequence,
		string(timestamp),
		commitMsg.Hash,
		commitMsg.Height}
}
