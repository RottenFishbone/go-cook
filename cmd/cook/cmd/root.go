package cmd

import (
	"os"

	"git.sr.ht/~rottenfishbone/go-cook/pkg/config"
	"github.com/spf13/cobra"
)

// cook
var (
	cfgPath   string
	cfgUseEnv bool

	rootCmd = &cobra.Command{
		Use:   "cook",
		Short: "cook provides a CLI to interact with recipes written in `cooklang`",
		Long:  "",

		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}
		},
	}
)

func Execute() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgPath, "config", "",
		"Path to config toml. (defaults to $COOK_HOME then $XDG_DATA_HOME/cook/config.toml)")
	rootCmd.PersistentFlags().BoolVarP(&cfgUseEnv, "env", "", false,
		"Use environment variables instead of a config file.")
	rootCmd.MarkFlagsMutuallyExclusive("env", "config")
}

func initConfig() {
	if cfgUseEnv {
		config.LoadConfigEnv()
		return
	}

	// Load config file into config.Vars
	if !config.LoadConfig(cfgPath) {
		if cfgPath == "" {
			_, err := os.Stat(config.DefaultConfigPath())
			// Test for its existence so we can figure out why it didn't open
			if os.IsNotExist(err) {
				os.Stderr.WriteString(
					"No config file found, run `cook init` to create one automatically.")
			} else if os.IsPermission(err) {
				os.Stderr.WriteString(
					"Could not open config file. Ensure relevant permissions are set.\n")
			}
			os.Exit(1)
		}
	} else {
		config.EnsureDataDirInit()
	}
}
