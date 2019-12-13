package main

import (
	"fmt"
	"github.com/chenzhou9513/redimint/core"
	"github.com/chenzhou9513/redimint/database"
	"github.com/chenzhou9513/redimint/ipfs"
	"github.com/chenzhou9513/redimint/logger"
	"github.com/chenzhou9513/redimint/network"
	"github.com/chenzhou9513/redimint/utils"
	abciserver "github.com/tendermint/tendermint/abci/server"
	tlog "github.com/tendermint/tendermint/libs/log"

	"os"
	"os/signal"
	"syscall"
)

func Init() {
	utils.InitKey()
	utils.InitFiles()
	utils.InitConfig()

	ipfs.InitIPFS()
	core.InitClient()
	core.InitService()
	logger.InitLogger()
	database.InitRedis()
	core.InitLogStoreApplication()
}

func main() {

	Init()
	logger := tlog.NewTMLogger(tlog.NewSyncWriter(os.Stdout))
	server := abciserver.NewSocketServer(core.SocketAddr, core.LogStoreApp)
	server.SetLogger(logger)
	if err := server.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "error starting socket server: %v", err)
		os.Exit(1)
	}
	defer server.Stop()

	appServer := network.NewServer()
	appServer.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	os.Exit(0)
}
