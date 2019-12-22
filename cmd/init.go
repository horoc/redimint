package cmd

import (
	"github.com/chenzhou9513/redimint/utils"
	"github.com/spf13/cobra"
	"os/exec"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialization redimint service",
	Long:  `Initialization redimint service, init all basic file under chain directory`,
	Run:   initRedimint,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

/*
	├── HomeDir
	│   ├── bin
	│   │   ├── redimint
	│   ├── conf
	│   │   ├── redis.conf
	│   │   ├── configuration.yaml
	│   ├── chain
	│   │   ├── config
	│   │   │   ├── genesis.json
	│   │   │   ├── config.toml
	│   │   │   ├── ... ...
	│   │   ├── data
	│   │   │   ├── ... ...
*/
func initRedimint(cmd *cobra.Command, args []string) {
	initTendermint()
}

func initTendermint() {
	utils.DeleteFile("../chain")
	utils.DeleteFile("./tendermint.sock")
	cmd := exec.Command("tendermint", "init", "--home=../chain")
	cmd.Run()
}
