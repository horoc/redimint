package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"runtime"
)

var (
	Version  string
	Revision string
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get redimint version",
	Long:  ``,
	Run:   getVersion,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func getVersion(*cobra.Command, []string) {
	fmt.Printf("Redimint Version:      %s\nGit revision: %s\nGo version:   %s\n", Version, Revision, runtime.Version())
}
