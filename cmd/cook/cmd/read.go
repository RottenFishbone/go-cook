package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"git.sr.ht/~rottenfishbone/go-cook"
	"git.sr.ht/~rottenfishbone/go-cook/pkg/common"
	"git.sr.ht/~rottenfishbone/go-cook/pkg/config"
	"git.sr.ht/~rottenfishbone/go-cook/pkg/recipe"
	"github.com/spf13/cobra"
)

var readCmd = &cobra.Command{
	Use:   "read",
	Short: "Parses a recipe file and pretty prints it to stdout",
	Long: `Parses a .cook file and prints it to stdout. 

If the file does not exist at the passed location, the recipes folder will be searched
for it.`,

	// Print help if no arguments are passed
	PreRun: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}

		initConfig()
	},

	Run: func(cmd *cobra.Command, args []string) {
		recipeDir := config.Get(config.KeyRecipeDir)
		for _, path := range args {
			var r cook.Recipe

			// If not a local file, check recipes dir
			if !common.FileExists(path) {
				// Rebuild path
				newPath := filepath.Join(recipeDir, path)
				if filepath.Ext(newPath) == "" {
					newPath = newPath + ".cook"
				}

				// If not in recipes either, exit with an error
				if !common.FileExists(newPath) {
					errTxt := fmt.Sprintf("Recipe %v does not exist.\n", path)
					os.Stderr.WriteString(errTxt)
					os.Exit(1)
				}

				path = newPath
			}

			// Read file to bytes, parse and print
			data, err := os.ReadFile(path)
			if err != nil {
				panic(err)
			}
			r = cook.ParseRecipe(recipe.FilepathToName(path), &data)
			recipe.PrettyPrint(&r)
		}
	},
}

func init() {
	rootCmd.AddCommand(readCmd)
}
