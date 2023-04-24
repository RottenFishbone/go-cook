package cmd

import (
	"fmt"
	"os"

	"git.sr.ht/~rottenfishbone/cooklang-go/pkg/config"
	"github.com/spf13/cobra"
)

// cook
var (
    cfgPath string
    cfgUseEnv bool

    rootCmd = &cobra.Command{
        Use:   "cook",
        Short: "cook provides a CLI to interact with recipes written in `cooklang`",
        Long: "",
        Run: func(cmd *cobra.Command, args []string) {
            if len(args) == 0 {
                cmd.Help()
                os.Exit(0)
            } else {
                initConfig()
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
    rootCmd.PersistentFlags().StringVar(&cfgPath, "config", "", "Path to config toml. (defaults to $COOK_HOME then $XDG_DATA_HOME/cook/config.toml)")
    rootCmd.PersistentFlags().BoolVarP(&cfgUseEnv, "env", "", false, "Use environment variables instead of a config file.")
}

func initConfig() {
    if cfgUseEnv {
        // Prevent conflicting config arguments
        if cfgPath != "" {
            os.Stderr.WriteString("Cannot use --config and --env flag simultaneously")
            os.Exit(1)
        }

        config.LoadConfigEnv()
        return 
    } 

    if !config.LoadConfig(cfgPath) {
        if cfgPath == "" {
            fmt.Println("No config file found, run `cook init` to create one automatically.")
        }
    }
}


