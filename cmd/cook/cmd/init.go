package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"git.sr.ht/~rottenfishbone/go-cook/pkg/config"
	"git.sr.ht/~rottenfishbone/go-cook/pkg/seed"
	"github.com/spf13/cobra"
)

var (
	// A flag to determine if example recipes should be created.
	seedFlag bool
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

		created := config.ConfigInit(path, recipesPath, shoppingPath)
		if !created && !seedFlag {
			os.Stderr.WriteString("Config already exists...aborting.\n")
			os.Exit(1)
		} else if created {
			fmt.Printf("Created config at: %v\n", path)
		}

		// Populate config so we can figure out where recipes go
		if !config.LoadConfig(path) {
			os.Stderr.WriteString(
				"Could not open config file. Ensure the relevant permissions are set..\n")
			os.Exit(1)
		}

		if seedFlag {
			// Seed into the recipes folder
			recipesPath = config.GetConfig().Recipe.Dir
			fmt.Printf("Seeded recipes into: %v\n", recipesPath)
			seed.SeedToDir(recipesPath)
		}
	},
}

func init() {
	initCmd.Flags().BoolVarP(&seedFlag, "seed", "s", false,
		"Seeds the recipe directory with a few sample files (will not overwrite config)")

	rootCmd.AddCommand(initCmd)
}

func spawnDataDirs() {

}
