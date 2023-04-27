package cmd

import (
	"git.sr.ht/~rottenfishbone/go-cook/pkg/server"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Hosts a local webserver to view/manage recipes.",
	Long:  ``,

	PreRun: func(cmd *cobra.Command, args []string) {
		initConfig()
	},

	Run: func(cmd *cobra.Command, args []string) {
		server.Start(6969)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
