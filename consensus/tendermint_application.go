package consensus

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/chenzhou9513/DecentralizedRedis/database"
	"github.com/chenzhou9513/DecentralizedRedis/logger"
	"github.com/chenzhou9513/DecentralizedRedis/models"
	c "github.com/chenzhou9513/DecentralizedRedis/models/code"
	"github.com/chenzhou9513/DecentralizedRedis/utils"
	"github.com/dgraph-io/badger"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	tmtypes "github.com/tendermint/tendermint/types"
	"go.uber.org/atomic"
	"os"
	"strconv"
	"strings"
)

const (
	ValidatorSetChangePrefix string = "val:"
)

/*

	内存db的组织形式：
		1. logIndex : log == 1 : set x 1
		2. heightIndex : logList == 1 : set a 1 || set c 1 || ...

*/

type LogStoreApplication struct {
	db                *badger.DB
	currentBatch      *badger.Txn
	currentHeight     atomic.Int64
	currentHeightList []string

	valUpdates         []abcitypes.ValidatorUpdate
	valAddrToPubKeyMap map[string]abcitypes.PubKey
	valAddrVote        Vote

	logSize atomic.Int64
}

var _ abcitypes.Application = (*LogStoreApplication)(nil)
var LogStoreApp *LogStoreApplication

const SEP string = "||"
const SocketAddr string = "unix://tendermint.sock"
const BadgerPath string = "/tmp/badger"

func InitLogStoreApplication() {
	logDb, err := badger.Open(badger.DefaultOptions(BadgerPath))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open badger db: %v", err)
		os.Exit(1)
	}
	LogStoreApp = &LogStoreApplication{
		db:                 logDb,
		valAddrToPubKeyMap: make(map[string]abcitypes.PubKey),
		valAddrVote:        NewVote(),
	}
}

func (LogStoreApplication) IsValidatorUpdateTx(tx []byte) bool {
	commitBody := models.TxCommitBody{}
	utils.JsonToStruct(tx, &commitBody)

	if commitBody.Data != nil {
		return false
	}
	return true
}

func (LogStoreApplication) SetOption(req abcitypes.RequestSetOption) abcitypes.ResponseSetOption {
	return abcitypes.ResponseSetOption{}
}

func (LogStoreApplication) Info(req abcitypes.RequestInfo) abcitypes.ResponseInfo {
	return abcitypes.ResponseInfo{}
}

func (app *LogStoreApplication) isValid(tx []byte) (uint32, string) {
	// TODO check query or execution
	// TODO 语法检验
	// TODO 过滤一些不能用的命令

	var data []byte
	var sign string
	var address string

	if !app.IsValidatorUpdateTx(tx) {
		commitBody := models.TxCommitBody{}
		utils.JsonToStruct(tx, &commitBody)
		data = utils.StructToJson(commitBody.Data)
		sign = commitBody.Signature
		address = commitBody.Address
	} else {
		commitBody := models.ValidatorUpdateBody{}
		utils.JsonToStruct(tx, &commitBody)
		data = utils.StructToJson(commitBody.ValidatorUpdate)
		sign = commitBody.Signature
		address = commitBody.Address
	}

	if _, ok := app.valAddrToPubKeyMap[address]; !ok {
		return c.CodeTypeInvalidValidator, c.Info(c.CodeTypeInvalidValidator)
	}

	pubkey := ed25519.PubKeyEd25519{}
	copy(pubkey[:], app.valAddrToPubKeyMap[address].Data)

	if pubkey.VerifyBytes(data, utils.HexToByte(sign)) != true {
		return c.CodeTypeInvalidSign, c.Info(c.CodeTypeInvalidSign)
	}

	return c.CodeTypeOK, c.Info(c.CodeTypeOK)
}

func (app *LogStoreApplication) CheckTx(req abcitypes.RequestCheckTx) abcitypes.ResponseCheckTx {
	code, info := app.isValid(req.Tx)
	return abcitypes.ResponseCheckTx{Code: code, GasWanted: 1, Info: info}
}

func (app *LogStoreApplication) InitChain(req abcitypes.RequestInitChain) abcitypes.ResponseInitChain {
	//TODO 清空Redis, 重启的是时候会初始化
	for _, v := range req.Validators {
		r := app.updateValidator(v)
		if r.IsErr() {
			logger.Error(r)
		}
	}
	return abcitypes.ResponseInitChain{}
}

func (app *LogStoreApplication) BeginBlock(req abcitypes.RequestBeginBlock) abcitypes.ResponseBeginBlock {
	app.currentHeightList = make([]string, 0)
	app.currentHeight.CAS(app.currentHeight.Load(), req.Header.Height)
	app.currentBatch = app.db.NewTransaction(true)

	//重置
	app.valUpdates = make([]abcitypes.ValidatorUpdate, 0)

	for _, ev := range req.ByzantineValidators {
		if ev.Type == tmtypes.ABCIEvidenceTypeDuplicateVote {
			// decrease voting power by 1
			if ev.TotalVotingPower == 0 {
				continue
			}
			app.updateValidator(abcitypes.ValidatorUpdate{
				PubKey: app.valAddrToPubKeyMap[string(ev.Validator.Address)],
				Power:  ev.TotalVotingPower - 1,
			})
		}
	}
	return abcitypes.ResponseBeginBlock{}
}

func (app *LogStoreApplication) DeliverTx(req abcitypes.RequestDeliverTx) abcitypes.ResponseDeliverTx {
	if app.IsValidatorUpdateTx(req.Tx) {
		return app.execValidatorTx(req.Tx)
	}

	response := abcitypes.ResponseDeliverTx{}
	code, msg := app.isValid(req.Tx)
	if code != c.CodeTypeOK {
		return abcitypes.ResponseDeliverTx{Code: code, Info: msg}
	}
	commitBody := models.TxCommitBody{}
	utils.JsonToStruct(req.Tx, &commitBody)

	app.logSize.Add(1)
	err := app.currentBatch.Set([]byte(strconv.FormatInt(app.logSize.Load(), 10)), []byte(commitBody.Data.Operation))
	app.currentHeightList = append(app.currentHeightList, commitBody.Data.Operation)
	//TODO 事务提交, 重试
	res, err := database.ExecuteCommand(commitBody.Data.Operation)
	response.Code = c.CodeTypeOK
	response.Data = []byte("Result:" + res)
	if err != nil {
		logger.Error(err)
		panic(err)
	}
	return response
}

func (app *LogStoreApplication) Commit() abcitypes.ResponseCommit {
	app.currentBatch.Commit()
	txn := app.db.NewTransaction(true)
	err := txn.Set([]byte(fmt.Sprintf("h%d", app.currentHeight.Load())), []byte(strings.Join(app.currentHeightList, SEP)))
	if err != nil {
		logger.Error(err)
		panic(err)
	}
	err = txn.Commit()
	if err != nil {
		logger.Error(err)
		panic(err)
	}
	return abcitypes.ResponseCommit{Data: []byte{}}
}

func (app *LogStoreApplication) EndBlock(req abcitypes.RequestEndBlock) abcitypes.ResponseEndBlock {
	return abcitypes.ResponseEndBlock{ValidatorUpdates: app.valUpdates}
}

//val:address1!power1  一次只允许更新一个
func (app *LogStoreApplication) execValidatorTx(tx []byte) abcitypes.ResponseDeliverTx {
	update := models.ValidatorUpdateBody{}
	utils.JsonToStruct(tx, &update)

	pubkeyS, powerS := update.ValidatorUpdate.PublicKey, update.ValidatorUpdate.Power

	// decode the pubkey
	pubkey, err := base64.StdEncoding.DecodeString(pubkeyS)
	if err != nil {
		return abcitypes.ResponseDeliverTx{
			Code: c.CodeTypeEncodingError,
			Log:  c.Info(c.CodeTypeEncodingError)}
	}

	// decode the power
	power, err := strconv.ParseInt(powerS, 10, 64)
	if err != nil {
		return abcitypes.ResponseDeliverTx{
			Code: c.CodeTypeEncodingError,
			Log:  c.Info(c.CodeTypeEncodingError)}
	}

	app.valAddrVote.addVote(pubkeyS, update.Address)
	if app.valAddrVote.getVoteNum(pubkeyS) >= len(app.valAddrToPubKeyMap) {
		app.updateValidator(abcitypes.Ed25519ValidatorUpdate(pubkey, power))
	}

	return abcitypes.ResponseDeliverTx{Code: 0}
}

func (app *LogStoreApplication) GetLogFromHeight(h int) []string {
	var logs = make([]string, 0)
	err := app.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(fmt.Sprintf("h%d", h)))
		if err != nil && err != badger.ErrKeyNotFound {
			logger.Error(err)
			return err
		}
		if err == badger.ErrKeyNotFound {
			return nil
		} else {
			return item.Value(func(val []byte) error {
				split := strings.Split(string(val), SEP)
				for _, v := range split {
					logs = append(logs, v)
				}
				return nil
			})
		}

	})
	if err != nil {
		logger.Error(err)
		panic(err)
	}
	return logs
}

func (app *LogStoreApplication) Query(reqQuery abcitypes.RequestQuery) (resQuery abcitypes.ResponseQuery) {
	resQuery.Key = reqQuery.Data
	err := app.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(reqQuery.Data)
		if err != nil && err != badger.ErrKeyNotFound {
			logger.Error(err)
			return err
		}
		if err == badger.ErrKeyNotFound {
			resQuery.Log = "does not exist"
		} else {
			return item.Value(func(val []byte) error {
				resQuery.Log = "exists"
				resQuery.Value = val
				return nil
			})
		}
		return nil
	})
	if err != nil {
		logger.Error(err)
		panic(err)
	}
	return
}

func (app *LogStoreApplication) updateValidator(v abcitypes.ValidatorUpdate) abcitypes.ResponseDeliverTx {

	pubkey := ed25519.PubKeyEd25519{}
	copy(pubkey[:], v.PubKey.Data)
	val := app.getBadgerVal([]byte("val:" + string(v.PubKey.Data)))

	if v.Power == 0 {
		// remove validator
		if val != nil {
			pubStr := base64.StdEncoding.EncodeToString(v.PubKey.Data)
			return abcitypes.ResponseDeliverTx{
				Code: c.CodeTypeEncodingError,
				Log:  c.InfoWithDetail(c.CodeTypeEncodingError, pubStr)}
		}
		app.deleteBadgerKey([]byte("val:" + string(v.PubKey.Data)))
		delete(app.valAddrToPubKeyMap, pubkey.Address().String())
	} else {
		// add or update validator
		value := bytes.NewBuffer(make([]byte, 0))
		if err := abcitypes.WriteMessage(&v, value); err != nil {
			return abcitypes.ResponseDeliverTx{
				Code: c.CodeTypeEncodingError,
				Log:  c.InfoWithDetail(c.CodeTypeEncodingError, err.Error())}
		}
		app.updateBadgerVal([]byte("val:"+string(v.PubKey.Data)), value.Bytes())
		app.valAddrToPubKeyMap[pubkey.Address().String()] = v.PubKey
	}

	// we only update the changes array if we successfully updated the tree
	app.valUpdates = append(app.valUpdates, v)

	return abcitypes.ResponseDeliverTx{Code: 0}
}

func (app *LogStoreApplication) updateBadgerVal(key []byte, val []byte) error {
	txn := app.db.NewTransaction(true)
	if err := txn.Set(key, val); err == nil {
		_ = txn.Commit()
	} else {
		return err
	}
	return nil
}

func (app *LogStoreApplication) deleteBadgerKey(key []byte) error {
	txn := app.db.NewTransaction(true)
	if err := txn.Delete(key); err == nil {
		_ = txn.Commit()
	} else {
		return err
	}
	return nil
}

func (app *LogStoreApplication) getBadgerVal(key []byte) []byte {
	var value []byte
	app.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		} else {
			return item.Value(func(val []byte) error {
				value = val
				return nil
			})
		}
		return nil
	})
	return value
}
