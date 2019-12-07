package main

import (
	"fmt"
	"github.com/chenzhou9513/DecentralizedRedis/core"
	"github.com/chenzhou9513/DecentralizedRedis/database"
	"github.com/chenzhou9513/DecentralizedRedis/ipfs"
	"github.com/chenzhou9513/DecentralizedRedis/logger"
	"github.com/chenzhou9513/DecentralizedRedis/network"
	"github.com/chenzhou9513/DecentralizedRedis/utils"
	abciserver "github.com/tendermint/tendermint/abci/server"
	"github.com/tendermint/tendermint/libs/log"

	"os"
	"os/signal"
	"syscall"
)

func Init() {
	utils.InitKey()
	utils.InitConfig()

	ipfs.InitIPFS()
	logger.InitLogger()
	core.InitClient()
	core.InitService()
	database.InitRedis()
	core.InitLogStoreApplication()
}

func main() {

	Init()
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
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
