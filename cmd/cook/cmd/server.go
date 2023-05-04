package cmd

import (
	"git.sr.ht/~rottenfishbone/go-cook/internal/pkg/server"
	"github.com/spf13/cobra"
)

var noWebapp bool

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Hosts a local webserver to view/manage recipes.",
	Long:  ``,

	PreRun: func(cmd *cobra.Command, args []string) {
		initConfig()
	},

	Run: func(cmd *cobra.Command, args []string) {
		server.Start(6969, noWebapp)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().BoolVarP(&noWebapp, "no-webapp", "", false, "Host the API server without the web app")
}
