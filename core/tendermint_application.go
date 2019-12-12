package core

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
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/tendermint/tendermint/abci/example/code"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	tmtypes "github.com/tendermint/tendermint/types"
	"go.uber.org/atomic"
	"os"
	"strconv"
	"strings"
	"sync"
)

type LogStoreApplication struct {
	db           *badger.DB
	currentBatch *badger.Txn

	valUpdates             []abcitypes.ValidatorUpdate
	valAddrToPubKeyMap     map[string]abcitypes.PubKey
	valAddrVote            Vote
	committedValidatorVote Vote

	currentCommittedHeight int64
	currentHeight          int64

	privateTxSet *hashset.Set

	logSize atomic.Int64

	initFlag bool

	pauseFlag bool
	lock      sync.Mutex
	wg        sync.WaitGroup
}

var _ abcitypes.Application = (*LogStoreApplication)(nil)
var LogStoreApp *LogStoreApplication

const PrivateSep string = "_"
const VoteKeySep string = ":"

const SocketAddr string = "unix://tendermint.sock"
const BadgerPath string = "/tmp/badger"

func InitLogStoreApplication() {
	logDb, err := badger.Open(badger.DefaultOptions(BadgerPath))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open badger db: %v", err)
		os.Exit(1)
	}
	LogStoreApp = &LogStoreApplication{
		db:                     logDb,
		valAddrToPubKeyMap:     make(map[string]abcitypes.PubKey),
		valAddrVote:            NewVote(),
		committedValidatorVote: NewVote(),
		initFlag:               true,
		currentCommittedHeight: 1,
		privateTxSet:           hashset.New(),
	}
}

func (LogStoreApplication) SetOption(req abcitypes.RequestSetOption) abcitypes.ResponseSetOption {
	return abcitypes.ResponseSetOption{}
}

func (LogStoreApplication) Info(req abcitypes.RequestInfo) abcitypes.ResponseInfo {
	return abcitypes.ResponseInfo{}
}

func (app *LogStoreApplication) isValid(tx []byte) (uint32, string) {
	var data []byte
	var sign string
	var address string

	if !app.IsValidatorUpdateTx(tx) {
		commitBody := models.TxCommitBody{}
		utils.JsonToStruct(tx, &commitBody)
		data = utils.StructToJson(commitBody.Data)
		sign = commitBody.Signature
		address = commitBody.Address
		if !database.IsValidCmd(commitBody.Data.Operation) {
			return c.CodeTypeInvalidTx, fmt.Sprintf("Invalid redis command : %s", commitBody.Data.Operation)
		}
		if database.IsQueryCmd(commitBody.Data.Operation) {
			return c.CodeTypeInvalidTx, fmt.Sprintf("Read only command can not commit to tendermint : %s", commitBody.Data.Operation)
		}
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

//#################### InitChain ####################
func (app *LogStoreApplication) InitChain(req abcitypes.RequestInitChain) abcitypes.ResponseInitChain {
	//TODO 清空Redis, 重启的是时候会初始化
	for _, v := range req.Validators {
		r := app.updateValidator(v)
		if r.IsErr() {
			logger.Log.Error(r)
		}
	}
	return abcitypes.ResponseInitChain{}
}

//#################### CheckTx ####################
func (app *LogStoreApplication) CheckTx(req abcitypes.RequestCheckTx) abcitypes.ResponseCheckTx {

	if app.pauseFlag == true {
		app.wg.Wait()
	}

	app.initFlag = false
	code, info := app.isValid(req.Tx)
	if code != c.CodeTypeOK {
		return abcitypes.ResponseCheckTx{Code: code, GasWanted: 1, Info: info}
	}

	if !app.IsValidatorUpdateTx(req.Tx) {
		commitBody := models.TxCommitBody{}
		utils.JsonToStruct(req.Tx, &commitBody)
		if app.IsPrivateCommand(commitBody) && strings.EqualFold(commitBody.Address, utils.ValidatorKey.Address.String()) {
			res, err := database.ExecuteCommand(commitBody.Data.Operation)
			if err != nil {
				logger.Log.Error(err)
				panic(err)
			}
			app.privateTxSet.Add(commitBody.Data.Sequence)
			return abcitypes.ResponseCheckTx{Code: code, GasWanted: 1, Info: info, Data: []byte("Result:" + res)}
		}
	}

	return abcitypes.ResponseCheckTx{Code: code, GasWanted: 1, Info: info}
}

//#################### BeginBlock ####################
func (app *LogStoreApplication) BeginBlock(req abcitypes.RequestBeginBlock) abcitypes.ResponseBeginBlock {

	app.currentHeight = req.Header.Height

	if app.pauseFlag == true {
		app.wg.Wait()
	}

	//reset valUpdates
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

//#################### DeliverTx ####################
func (app *LogStoreApplication) DeliverTx(req abcitypes.RequestDeliverTx) abcitypes.ResponseDeliverTx {
	if app.IsValidatorUpdateTx(req.Tx) {
		return app.execValidatorTx(req.Tx)
	}

	code, info := app.isValid(req.Tx)
	if code != c.CodeTypeOK {
		return abcitypes.ResponseDeliverTx{Code: code, Info: info}
	}

	commitBody := models.TxCommitBody{}
	utils.JsonToStruct(req.Tx, &commitBody)
	if app.IsPrivateCommand(commitBody) {
		if app.initFlag || !strings.EqualFold(commitBody.Address, utils.ValidatorKey.Address.String()) {
			if !app.privateTxSet.Contains(commitBody.Data.Sequence) {
				_, err := database.ExecuteCommand(commitBody.Data.Operation)
				if err != nil {
					logger.Log.Error(err)
					panic(err)
				}
			} else {
				app.privateTxSet.Remove(commitBody.Data.Sequence)
			}
		}
		return abcitypes.ResponseDeliverTx{Code: c.CodeTypeOK, Info: c.Info(c.CodeTypeOK)}
	} else {
		app.logSize.Add(1)
		res, err := database.ExecuteCommand(commitBody.Data.Operation)
		if err != nil {
			logger.Log.Error(err)
			panic(err)
		}
		return abcitypes.ResponseDeliverTx{Code: c.CodeTypeOK, Data: []byte("Result:" + res)}
	}
}

//#################### Commit ####################
func (app *LogStoreApplication) Commit() abcitypes.ResponseCommit {
	return abcitypes.ResponseCommit{Data: []byte{}}
}

//#################### EndBlock ####################
func (app *LogStoreApplication) EndBlock(req abcitypes.RequestEndBlock) abcitypes.ResponseEndBlock {
	app.currentCommittedHeight = app.currentHeight
	return abcitypes.ResponseEndBlock{ValidatorUpdates: app.valUpdates}
}

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

	//publicKey:power
	voteKey := pubkeyS + VoteKeySep + powerS

	app.valAddrVote.addVote(voteKey, update.Address, update)
	count := app.valAddrVote[voteKey]
	if app.valAddrVote.getVoteNum(voteKey) >= len(app.valAddrToPubKeyMap)*2/3 {
		app.updateValidator(abcitypes.Ed25519ValidatorUpdate(pubkey, power))
		app.committedValidatorVote[voteKey] = count
		delete(app.valAddrVote, voteKey)
	}

	return abcitypes.ResponseDeliverTx{Code: code.CodeTypeOK, Data: utils.StructToJson(count)}
}

func (app *LogStoreApplication) QueryVotingValidators() *Vote {
	return &app.valAddrVote
}

func (app *LogStoreApplication) Query(reqQuery abcitypes.RequestQuery) (resQuery abcitypes.ResponseQuery) {
	resQuery.Key = reqQuery.Data
	err := app.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(reqQuery.Data)
		if err != nil && err != badger.ErrKeyNotFound {
			logger.Log.Error(err)
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
		logger.Log.Error(err)
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

func (app *LogStoreApplication) GetCurrentHeight() int64 {
	return app.currentCommittedHeight
}

func (app *LogStoreApplication) Pause() {
	app.lock.Lock()
	app.pauseFlag = true
	app.wg.Add(1)
	app.lock.Unlock()
}

func (app *LogStoreApplication) Continue() {
	app.lock.Lock()
	app.pauseFlag = false
	app.wg.Done()
	app.lock.Unlock()
}

func (LogStoreApplication) IsValidatorUpdateTx(tx []byte) bool {
	commitBody := models.TxCommitBody{}
	utils.JsonToStruct(tx, &commitBody)

	if commitBody.Data != nil {
		return false
	}
	return true
}

func (app *LogStoreApplication) IsPrivateCommand(commitBody models.TxCommitBody) bool {
	key, err := database.GetKey(commitBody.Data.Operation)
	if err != nil {
		logger.Log.Error(err)
		return false
	}
	split := strings.Split(key, PrivateSep)
	if _, ok := app.valAddrToPubKeyMap[split[0]]; ok && strings.EqualFold(split[0], commitBody.Address) {
		return true
	}
	return false
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
