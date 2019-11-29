package service

import (
	"fmt"
	"github.com/chenzhou9513/DecentralizedRedis/consensus"
	"github.com/chenzhou9513/DecentralizedRedis/models"
	"github.com/chenzhou9513/DecentralizedRedis/models/code"
	"github.com/chenzhou9513/DecentralizedRedis/utils"
	uuid "github.com/satori/go.uuid"
	"strconv"
	"time"
)

type ServiceImpl struct {
}

func MakeTxCommitBody(request *models.ExecuteRequest) *models.TxCommitBody {
	op := &models.TxCommitBody{}
	op.Data = &models.TxCommitData{}

	//Sequence
	u := uuid.NewV1()
	op.Data.Sequence = utils.ByteToHex(u.Bytes())

	//cmd
	op.Data.Operation = request.Cmd

	//Signature
	op.Signature = utils.ValidatorStringSign(utils.StructToJson(op.Data))

	//address
	op.Address = utils.ValidatorKey.Address.String()

	return op
}

func (s ServiceImpl) Execute(request *models.ExecuteRequest) *models.ExecuteResponse {
	op := MakeTxCommitBody(request)
	//timestamp
	timestamp := time.Now().UnixNano() / 1e6
	//tendermint response
	commitMsg := consensus.BroadcastTxCommit(op)

	//TODO 错误处理

	return &models.ExecuteResponse{Code: code.CodeTypeOK,
		CodeMsg:       code.Info(code.CodeTypeOK),
		Cmd:           request.Cmd,
		ExecuteResult: string(commitMsg.DeliverTx.Data),
		Signature:     op.Signature,
		Sequence:      op.Data.Sequence,
		TimeStamp:     strconv.FormatInt(timestamp, 10),
		Hash:          utils.ByteToHex(commitMsg.Hash),
		Height:        commitMsg.Height}
}

func (s ServiceImpl) ExecuteAsync(request *models.ExecuteRequest) *models.ExecuteAsyncResponse {
	op := MakeTxCommitBody(request)
	//timestamp
	timestamp := time.Now().UnixNano() / 1e6

	sync := consensus.BroadcastTxSync(op)

	fmt.Println(string(utils.StructToJson(sync)))

	return &models.ExecuteAsyncResponse{
		Code:      code.CodeTypeOK,
		CodeMsg:   code.Info(code.CodeTypeOK),
		Cmd:       op.Data.Operation,
		Signature: op.Signature,
		Sequence:  op.Data.Sequence,
		TimeStamp: strconv.FormatInt(timestamp, 10),
		Hash:      utils.ByteToHex(sync.Hash),
	}
}

func (s ServiceImpl) QueryTransaction(hash string) *models.Transaction {
	byteHash := utils.HexToByte(hash)
	tx := consensus.GetTx(byteHash)

	data := &models.TxCommitBody{}
	utils.JsonToStruct(tx.Tx, data)

	transaction := &models.Transaction{
		Hash:          utils.ByteToHex(tx.Hash),
		Height:        tx.Height,
		Index:         tx.Index,
		Data:          data,
		ExecuteResult: string(tx.TxResult.Data),
		Proof:         nil,
	}
	transaction.Proof = &models.TxProof{
		RootHash: utils.ByteToHex(tx.Proof.RootHash),
		Proof: &models.ProofDetail{
			Total:    tx.Proof.Proof.Total,
			Index:    tx.Proof.Proof.Index,
			LeafHash: "",
			Aunts:    nil,
		},
	}

	if tx.Proof.Proof.LeafHash != nil {
		transaction.Proof.Proof.LeafHash = utils.ByteToHex(tx.Proof.Proof.LeafHash)
	}

	if len(tx.Proof.Proof.Aunts) != 0 {
		transaction.Proof.Proof.Aunts = make([]string, 0)
		for _, v := range tx.Proof.Proof.Aunts {
			transaction.Proof.Proof.Aunts = append(transaction.Proof.Proof.Aunts, utils.ByteToHex(v))
		}
	}
	return transaction
}
