package cmd

import (
	"github.com/spf13/cobra"
)

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "restart redimint server",
	Long:  ``,
	Run:   restart,
}

func init() {
	rootCmd.AddCommand(restartCmd)
}

func restart(cmd *cobra.Command, args []string) {
	stop(cmd, args)
	daemon = true
	start(cmd, args)
}
