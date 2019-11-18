package main

import (
	"flag"
	"fmt"
	"github.com/chenzhou9513/DecentralizedRedis/consensus"
	"github.com/chenzhou9513/DecentralizedRedis/database"
	"github.com/chenzhou9513/DecentralizedRedis/network"
	"github.com/chenzhou9513/DecentralizedRedis/utils"
	"github.com/dgraph-io/badger"
	abciserver "github.com/tendermint/tendermint/abci/server"
	"github.com/tendermint/tendermint/libs/log"
	"os"
	"os/signal"
	"syscall"
)



var socketAddr string

func init() {
	flag.StringVar(&socketAddr, "socket-addr", "unix://tendermint.sock", "Unix domain socket address")
}

func main() {
	utils.InitConfig()

	database.InitRedisClient()

	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open badger db: %v", err)
		os.Exit(1)
	}
	defer db.Close()
	app := consensus.NewLogStoreApplication(db)

	flag.Parse()

	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	server := abciserver.NewSocketServer(socketAddr, app)
	server.SetLogger(logger)
	if err := server.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "error starting socket server: %v", err)
		os.Exit(1)
	}
	defer server.Stop()

	httpServer := network.NewServer("0.0.0.0","30001")
	httpServer.Start()



	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	os.Exit(0)
}
