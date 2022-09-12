package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var RootCmd = &cobra.Command{
	Short: "fuse file system",
	RunE: func(cmd *cobra.Command, args []string) error {
		return ShowHelp(os.Stderr)(cmd, args)
	},
}

func init() {
	RootCmd.AddCommand(ServerCmd)
	RootCmd.AddCommand(ClientCmd)
}
