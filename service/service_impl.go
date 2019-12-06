package service

import (
	"fmt"
	"github.com/chenzhou9513/DecentralizedRedis/core"
	"github.com/chenzhou9513/DecentralizedRedis/database"
	"github.com/chenzhou9513/DecentralizedRedis/logger"
	"github.com/chenzhou9513/DecentralizedRedis/models"
	"github.com/chenzhou9513/DecentralizedRedis/models/code"
	"github.com/chenzhou9513/DecentralizedRedis/utils"
	uuid "github.com/satori/go.uuid"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/types"
	"strconv"
	"time"
)

var AppService Service

type ServiceImpl struct {
}

const PRIVATE_SEP string = "_"

func InitService() {
	AppService = ServiceImpl{}
}

func (s ServiceImpl) MakeTxCommitBody(request *models.CommandRequest) *models.TxCommitBody {
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

func (s ServiceImpl) Execute(request *models.CommandRequest) *models.ExecuteResponse {
	op := s.MakeTxCommitBody(request)
	//timestamp
	timestamp := time.Now().UnixNano() / 1e6
	//tendermint response
	commitMsg := core.BroadcastTxCommit(op)

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

func (s ServiceImpl) ExecuteAsync(request *models.CommandRequest) *models.ExecuteAsyncResponse {
	op := s.MakeTxCommitBody(request)
	//timestamp
	timestamp := time.Now().UnixNano() / 1e6

	sync := core.BroadcastTxSync(op)

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

func (s ServiceImpl) ExecuteWithPrivateKey(request *models.CommandRequest) *models.ExecuteResponse {

	key, err := database.GetKey(request.Cmd)
	if err != nil {
		logger.Error(err)
	}
	key = utils.ValidatorKey.Address.String() + PRIVATE_SEP + key
	cmd, err := database.ReplaceKey(request.Cmd, key)
	if err != nil {
		logger.Error(err)
	}
	request.Cmd = cmd

	op := s.MakeTxCommitBody(request)

	//timestamp
	timestamp := time.Now().UnixNano() / 1e6

	sync := core.BroadcastTxSync(op)

	return &models.ExecuteResponse{
		Code:      code.CodeTypeOK,
		CodeMsg:   code.Info(code.CodeTypeOK),
		Cmd:       op.Data.Operation,
		ExecuteResult: string(sync.Data),
		Signature: op.Signature,
		Sequence:  op.Data.Sequence,
		TimeStamp: strconv.FormatInt(timestamp, 10),
		Hash:      utils.ByteToHex(sync.Hash),
	}
}

func (s ServiceImpl) Query(request *models.CommandRequest) *models.QueryResponse {
	result, err := database.ExecuteCommand(request.Cmd)
	if err != nil {
		return &models.QueryResponse{
			Code:    code.CodeTypeRedisExecutionError,
			CodeMsg: code.Info(code.CodeTypeRedisExecutionError) + ": " + err.Error(),
			Result:  "",
		}
	}
	return &models.QueryResponse{
		Code:    code.CodeTypeOK,
		CodeMsg: code.Info(code.CodeTypeOK),
		Result:  result,
	}
}

func (s ServiceImpl) QueryPrivateKey(request *models.CommandRequest, address string) *models.QueryResponse {

	key, err := database.GetKey(request.Cmd)
	if err != nil {
		logger.Error(err)
	}
	key = utils.ValidatorKey.Address.String() + PRIVATE_SEP + key
	cmd, err := database.ReplaceKey(request.Cmd, key)
	if err != nil {
		logger.Error(err)
	}
	request.Cmd = cmd
	result, err := database.ExecuteCommand(request.Cmd)
	if err != nil {
		return &models.QueryResponse{
			Code:    code.CodeTypeRedisExecutionError,
			CodeMsg: code.Info(code.CodeTypeRedisExecutionError) + ": " + err.Error(),
			Result:  "",
		}
	}
	return &models.QueryResponse{
		Code:    code.CodeTypeOK,
		CodeMsg: code.Info(code.CodeTypeOK),
		Result:  result,
	}
}

func (s ServiceImpl) RestoreLocalDatabase() error {
	database.StopRedis()
	database.StartRedisServer()

	return nil
}

func (s ServiceImpl) QueryTransaction(hash string) *models.Transaction {
	byteHash := utils.HexToByte(hash)
	tx := core.GetTx(byteHash)
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

func (s ServiceImpl) QueryBlock(height int) *models.Block {
	originBlock := core.GetBlockFromHeight(height)
	return s.ConvertBlock(originBlock)
}

func (s ServiceImpl) GetChainInfo(min int, max int) *models.ChainInfo {
	info := core.GetChainInfo(min, max)
	res := &models.ChainInfo{
		LastHeight: info.LastHeight,
		BlockMetas: make([]*models.BlockMeta, 0),
	}
	metas := info.BlockMetas
	for _, v := range metas {
		res.BlockMetas = append(res.BlockMetas, &models.BlockMeta{
			BlockID: *s.ConvertBlockID(&v.BlockID),
			Header:  *s.ConvertBlockHeader(&v.Header),
		})
	}
	return res
}

func (s ServiceImpl) GetChainState() *models.ChainState {
	originState := core.GetChainState()
	state := &models.ChainState{}
	utils.JsonToStruct(utils.StructToJson(originState), state)
	fmt.Println(originState.ValidatorInfo.PubKey)
	state.ValidatorInfo.PubKey = originState.ValidatorInfo.PubKey.Bytes()
	return state
}

func (s ServiceImpl) ConvertBlockID(b *types.BlockID) *models.BlockID {
	blockID := models.BlockID{}
	blockID.Hash = utils.ByteToHex(b.Hash)
	blockID.PartsHeader = models.PartSetHeader{
		Total: b.PartsHeader.Total,
		Hash:  utils.ByteToHex(b.PartsHeader.Hash),
	}
	return &blockID
}

func (s ServiceImpl) ConvertBlockHeader(b *types.Header) *models.Header {
	return &models.Header{
		Version:            b.Version,
		ChainID:            b.ChainID,
		Height:             b.Height,
		Time:               time.Time{},
		NumTxs:             b.NumTxs,
		TotalTxs:           b.TotalTxs,
		LastBlockID:        *s.ConvertBlockID(&b.LastBlockID),
		LastCommitHash:     utils.ByteToHex(b.LastCommitHash),
		DataHash:           utils.ByteToHex(b.DataHash),
		ValidatorsHash:     utils.ByteToHex(b.ValidatorsHash),
		NextValidatorsHash: utils.ByteToHex(b.NextValidatorsHash),
		ConsensusHash:      utils.ByteToHex(b.ConsensusHash),
		AppHash:            utils.ByteToHex(b.AppHash),
		LastResultsHash:    utils.ByteToHex(b.LastCommitHash),
		EvidenceHash:       utils.ByteToHex(b.EvidenceHash),
		ProposerAddress:    utils.ByteToHex(b.ProposerAddress),
	}
}

func (s ServiceImpl) ConvertBlockData(b *types.Data) *models.Data {
	data := &models.Data{
		Txs:  make([]string, 0),
		Hash: utils.ByteToHex(b.Hash()),
	}
	for _, v := range b.Txs {
		data.Txs = append(data.Txs, fmt.Sprintf("%x", v.Hash()))
	}
	return data
}

func (s ServiceImpl) ConvertCommitSign(b *types.CommitSig) *models.CommitSig {

	return &models.CommitSig{
		Type:             b.Type,
		Height:           b.Height,
		Round:            b.Round,
		Timestamp:        b.Timestamp,
		ValidatorAddress: utils.ByteToHex(b.ValidatorAddress),
		ValidatorIndex:   b.ValidatorIndex,
		Signature:        utils.ByteToHex(b.Signature),
	}
}

func (s ServiceImpl) ConvertBlock(b *ctypes.ResultBlock) *models.Block {

	blockID := s.ConvertBlockID(&(b.BlockMeta.BlockID))

	header := s.ConvertBlockHeader(&(b.Block.Header))

	data := s.ConvertBlockData(&(b.Block.Data))

	lastCommit := make([]*models.CommitSig, 0)

	for _, v := range b.Block.LastCommit.Precommits {
		lastCommit = append(lastCommit, s.ConvertCommitSign(v))
	}

	block := &models.Block{
		BlockID:    *blockID,
		Header:     *header,
		Data:       *data,
		Evidence:   b.Block.Evidence,
		LastCommit: lastCommit,
	}

	return block
}
