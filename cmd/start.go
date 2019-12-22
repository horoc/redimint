package cmd

import (
	"fmt"
	"github.com/chenzhou9513/redimint/core"
	"github.com/chenzhou9513/redimint/database"
	"github.com/chenzhou9513/redimint/ipfs"
	"github.com/chenzhou9513/redimint/logger"
	"github.com/chenzhou9513/redimint/network"
	"github.com/chenzhou9513/redimint/utils"
	"github.com/spf13/cobra"
	abciserver "github.com/tendermint/tendermint/abci/server"
	tlog "github.com/tendermint/tendermint/libs/log"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start redimint server",
	Long:  `Start redimint server, and includes tendermint server, redis server`,
	Run:   start,
}

const (
	TDPID_FILE = "./.tendermint_pid"
	RDPID_FILE = "./.redimint_pid"
	DBPID_FILE = "./.redis_pid"
)

var daemon bool

func init() {
	startCmd.Flags().BoolVarP(&daemon, "daemon", "d", false, "redimint start mode")
	rootCmd.AddCommand(startCmd)
}

func InitService() {
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

func start(cmd *cobra.Command, args []string) {
	if daemon {
		startTendermintDaemon()
		startRedisDaemon()
		startRedimintDaemon()
		os.Exit(0)
	}
	startTendermintDaemon()
	startRedisDaemon()
	InitService()
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

func startRedimintDaemon() {
	cmd := exec.Command("./redimint", "start")
	cmd.Start()
	fmt.Println("Redimint daemon process ID is : ", cmd.Process.Pid)
	savePID(cmd.Process.Pid, RDPID_FILE)
}

func startTendermintDaemon() {
	utils.DeleteFile("tendermint.sock")
	cmdStr := `nohup tendermint --home=../chain node --proxy_app=unix://tendermint.sock > ../log/tendermint.log 2>&1 &`
	cmd := exec.Command("bash", "-c", cmdStr)
	cmd.Start()
	fmt.Println("Tendermint daemon process ID is : ", cmd.Process.Pid)
	savePID(cmd.Process.Pid, TDPID_FILE)
}

func startRedisDaemon() {
	cmdStr := `nohup redis-server ../conf/redis.conf > ../log/redis.log 2>&1 &`
	cmd := exec.Command("bash", "-c", cmdStr)
	cmd.Start()
	fmt.Println("Redis daemon process ID is : ", cmd.Process.Pid)
	savePID(cmd.Process.Pid, DBPID_FILE)
}

func savePID(pid int, pidFile string) {

	file, err := os.Create(pidFile)
	if err != nil {
		log.Printf("Unable to create pid file : %v\n", err)
		os.Exit(1)
	}

	defer file.Close()

	_, err = file.WriteString(strconv.Itoa(pid))

	if err != nil {
		log.Printf("Unable to create pid file : %v\n", err)
		os.Exit(1)
	}

	file.Sync()
}
