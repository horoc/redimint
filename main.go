package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/dgraph-io/badger"
	abciserver "github.com/tendermint/tendermint/abci/server"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	"os"
	"os/signal"
	"syscall"
)

var _ abcitypes.Application = (*KVStoreApplication)(nil)

type KVStoreApplication struct {
	db           *badger.DB
	currentBatch *badger.Txn
}

func NewKVStoreApplication(db *badger.DB) *KVStoreApplication {
	return &KVStoreApplication{
		db: db,
	}
}

func (KVStoreApplication) Info(req abcitypes.RequestInfo) abcitypes.ResponseInfo {
	return abcitypes.ResponseInfo{}
}

func (KVStoreApplication) SetOption(req abcitypes.RequestSetOption) abcitypes.ResponseSetOption {
	return abcitypes.ResponseSetOption{}
}

func (app *KVStoreApplication) DeliverTx(req abcitypes.RequestDeliverTx) abcitypes.ResponseDeliverTx {
	code := app.isValid(req.Tx)
	if code != 0 {
		return abcitypes.ResponseDeliverTx{Code: code}
	}

	parts := bytes.Split(req.Tx, []byte("="))
	key, value := parts[0], parts[1]

	err := app.currentBatch.Set(key, value)
	if err != nil {
		panic(err)
	}

	return abcitypes.ResponseDeliverTx{Code: 0}
}

func (app *KVStoreApplication) Commit() abcitypes.ResponseCommit {
	app.currentBatch.Commit()
	return abcitypes.ResponseCommit{Data: []byte{}}
}
func (app *KVStoreApplication) Query(reqQuery abcitypes.RequestQuery) (resQuery abcitypes.ResponseQuery) {
	resQuery.Key = reqQuery.Data
	err := app.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(reqQuery.Data)
		if err != nil && err != badger.ErrKeyNotFound {
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
		panic(err)
	}
	return
}

func (KVStoreApplication) InitChain(req abcitypes.RequestInitChain) abcitypes.ResponseInitChain {
	return abcitypes.ResponseInitChain{}
}

func (app *KVStoreApplication) BeginBlock(req abcitypes.RequestBeginBlock) abcitypes.ResponseBeginBlock {
	app.currentBatch = app.db.NewTransaction(true)
	return abcitypes.ResponseBeginBlock{}
}

func (KVStoreApplication) EndBlock(req abcitypes.RequestEndBlock) abcitypes.ResponseEndBlock {
	return abcitypes.ResponseEndBlock{}
}

func (app *KVStoreApplication) isValid(tx []byte) (code uint32) {
	// check format
	parts := bytes.Split(tx, []byte("="))
	if len(parts) != 2 {
		return 1
	}

	key, _ := parts[0], parts[1]

	// check if the same key=value already exists
	err := app.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil && err != badger.ErrKeyNotFound {
			return err
		}
		if err == nil {
			fmt.Println(item)
			var value []byte
			item.Value(func(val []byte) error{
				value = val
				return nil
			})
			fmt.Println(string(value))
			if item != nil{
				code = 2
			}
			//return nil
			//
			//return item.Value(func(val []byte) error {
			//	if bytes.Equal(val, value) {
			//		code = 2
			//	}
			//	return nil
			//})
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	return code
}

func (app *KVStoreApplication) CheckTx(req abcitypes.RequestCheckTx) abcitypes.ResponseCheckTx {
	code := app.isValid(req.Tx)
	return abcitypes.ResponseCheckTx{Code: code, GasWanted: 1}
}


var socketAddr string

func init() {
	flag.StringVar(&socketAddr, "socket-addr", "unix://example.sock", "Unix domain socket address")
}

func main() {
	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open badger db: %v", err)
		os.Exit(1)
	}
	defer db.Close()
	app := NewKVStoreApplication(db)

	flag.Parse()

	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	server := abciserver.NewSocketServer(socketAddr, app)
	server.SetLogger(logger)
	if err := server.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "error starting socket server: %v", err)
		os.Exit(1)
	}
	defer server.Stop()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	os.Exit(0)
}
