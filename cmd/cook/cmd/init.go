package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"git.sr.ht/~rottenfishbone/cooklang-go/pkg/config"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Creates the default config file.",
	Long: `Creates a config file in $XDG_CONFIG_HOME, or at the passed location

A directory of where to store data such as recipes can optionally be passed as a second argument 
which will be added to the created config automatically.`,
	Run: func(cmd *cobra.Command, args []string) {
		var path string
		recipesPath := ""
		shoppingPath := ""
		switch len(args) {
		case 0:
			path = config.DefaultConfigPath()
		case 1:
			path = args[0]
		case 2:
			path = args[0]
			var err error
			dataPath, err := filepath.Abs(filepath.Join(args[1], "cook"))
			recipesPath = filepath.Join(dataPath, "recipes")
			shoppingPath = filepath.Join(dataPath, "shopping")
			if err != nil {
				panic(err)
			}
		default:
			os.Stderr.WriteString("Too many arguments. Exiting.\n")
			cmd.Usage()
			os.Exit(1)
		}

		path, err := filepath.Abs(path)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Creating config at: %v\n", path)
		config.ConfigInit(path, recipesPath, shoppingPath)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
