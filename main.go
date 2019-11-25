package main

import (
	"fmt"
	"github.com/chenzhou9513/DecentralizedRedis/consensus"
	"github.com/chenzhou9513/DecentralizedRedis/database"
	"github.com/chenzhou9513/DecentralizedRedis/network"
	"github.com/chenzhou9513/DecentralizedRedis/utils"
	abciserver "github.com/tendermint/tendermint/abci/server"
	"github.com/tendermint/tendermint/libs/log"
	logger"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)




func main() {



	logger.Info("Init service...")
	logger.Info("Init service...")
	logger.Info("Init service...")
	logger.Info("Init service...")

	utils.InitConfig()
	database.InitRedisClient()
	consensus.InitLogStoreApplication()
	tendermintlogger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	server := abciserver.NewSocketServer(consensus.SocketAddr, consensus.LogStoreApp)
	server.SetLogger(tendermintlogger)
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
