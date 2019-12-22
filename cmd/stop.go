package cmd

import (
	"fmt"
	"github.com/chenzhou9513/redimint/utils"
	"github.com/spf13/cobra"
	"os/exec"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "stop redimint server",
	Long:  ``,
	Run:   stop,
}

func init() {
	rootCmd.AddCommand(stopCmd)
}

func stop(cmd *cobra.Command, args []string) {
	pid := utils.ReadAll(DBPID_FILE)
	stopPID(pid)
	utils.DeleteFile(DBPID_FILE)

	pid = utils.ReadAll(TDPID_FILE)
	stopPID(pid)
	utils.DeleteFile(TDPID_FILE)

	pid = utils.ReadAll(RDPID_FILE)
	stopPID(pid)
	utils.DeleteFile(RDPID_FILE)
}

func stopPID(pid string) {
	cmd := exec.Command("kill", pid)
	cmd.Run()
	fmt.Println("Stop process ID is : ", pid)
}
