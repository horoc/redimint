package consensus

import (
	"flag"
	"fmt"
	"github.com/chenzhou9513/DecentralizedRedis/database"
	"github.com/chenzhou9513/DecentralizedRedis/logger"
	"github.com/chenzhou9513/DecentralizedRedis/utils"
	"github.com/dgraph-io/badger"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	"go.uber.org/atomic"
	"os"
	"strconv"
	"strings"
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
	logSize           atomic.Int64
}

var _ abcitypes.Application = (*LogStoreApplication)(nil)

var LogStoreApp *LogStoreApplication
var SocketAddr string

const SEP string = "||"

func init() {
	flag.StringVar(&SocketAddr, "socket-addr", "unix://tendermint.sock", "Unix domain socket address")
}

func InitLogStoreApplication() {
	flag.Parse()
	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open badger db: %v", err)
		os.Exit(1)
	}
	LogStoreApp = NewLogStoreApplication(db)
}

func NewLogStoreApplication(db *badger.DB) *LogStoreApplication {
	return &LogStoreApplication{
		db: db,
	}
}

func (LogStoreApplication) Info(req abcitypes.RequestInfo) abcitypes.ResponseInfo {
	return abcitypes.ResponseInfo{}
}

func (LogStoreApplication) SetOption(req abcitypes.RequestSetOption) abcitypes.ResponseSetOption {
	return abcitypes.ResponseSetOption{}
}

func (app *LogStoreApplication) DeliverTx(req abcitypes.RequestDeliverTx) abcitypes.ResponseDeliverTx {
	response := abcitypes.ResponseDeliverTx{}
	code := app.isValid(req.Tx)
	if code != 0 {
		return abcitypes.ResponseDeliverTx{Code: code}
	}
	commitBody := CommitBody{}
	utils.JsonToStruct(req.Tx, &commitBody)

	//parts := bytes.Split(req.Tx, []byte(SEP))
	//seq, value := parts[0], parts[1]
	app.logSize.Add(1)
	err := app.currentBatch.Set([]byte(strconv.FormatInt(app.logSize.Load(), 10)), []byte(commitBody.Operation))
	app.currentHeightList = append(app.currentHeightList, commitBody.Operation)
	//TODO 事务提交
	res := database.ExecuteCommand(commitBody.Operation)
	response.Code = 0
	//info, err := json.Marshal(struct {
	//	Operation string
	//	Sequence  string
	//	Signature string
	//	Result    string
	//}{commitBody.Operation, commitBody.Sequence, commitBody.Signature, res})
	response.Info = "Result:"+res
	if err != nil {
		logger.Error(err)
		panic(err)
	}
	return response
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

func (LogStoreApplication) InitChain(req abcitypes.RequestInitChain) abcitypes.ResponseInitChain {
	//TODO 清空Redis, 重启的是时候会初始化
	return abcitypes.ResponseInitChain{}
}

func (app *LogStoreApplication) BeginBlock(req abcitypes.RequestBeginBlock) abcitypes.ResponseBeginBlock {
	app.currentHeightList = make([]string, 0)
	app.currentHeight.CAS(app.currentHeight.Load(), req.Header.Height)
	app.currentBatch = app.db.NewTransaction(true)
	return abcitypes.ResponseBeginBlock{}
}

func (LogStoreApplication) EndBlock(req abcitypes.RequestEndBlock) abcitypes.ResponseEndBlock {
	return abcitypes.ResponseEndBlock{}
}

func (app *LogStoreApplication) isValid(tx []byte) (code uint32) {
	// TODO check query or execution
	// TODO 语法检验
	// TODO 过滤一些不能用的命令
	//// check format
	//parts := bytes.Split(tx, []byte("="))
	//if len(parts) != 2 {
	//	return 1
	//}
	//
	//key, _ := parts[0], parts[1]
	//
	//// check if the same key=value already exists
	//err := app.db.View(func(txn *badger.Txn) error {
	//	item, err := txn.Get(key)
	//	if err != nil && err != badger.ErrKeyNotFound {
	//		return err
	//	}
	//	if err == nil {
	//		fmt.Println(item)
	//		var value []byte
	//		item.Value(func(val []byte) error {
	//			value = val
	//			return nil
	//		})
	//		fmt.Println(string(value))
	//		if item != nil {
	//			code = 2
	//		}
	//	}
	//	return nil
	//})
	//if err != nil {
	//	panic(err)
	//}

	return code
}

func (app *LogStoreApplication) CheckTx(req abcitypes.RequestCheckTx) abcitypes.ResponseCheckTx {
	code := app.isValid(req.Tx)
	return abcitypes.ResponseCheckTx{Code: code, GasWanted: 1}
}
