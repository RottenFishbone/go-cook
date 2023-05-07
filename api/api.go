package api

// TODO cache recipes in memory

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"git.sr.ht/~rottenfishbone/go-cook/internal/pkg/common"
	"git.sr.ht/~rottenfishbone/go-cook/pkg/config"
	"git.sr.ht/~rottenfishbone/go-cook/pkg/recipe"
)

// Ensures the config file is loaded. Panic on failure.
//
// An unloaded read should never occur at runtime, unless there is a bug.
func assertConfigLoaded() {
	if !config.IsLoaded() {
		panic("Attempted to read unloaded config")
	}
}

// Collects the name of each subdirectory (and root) of the recipes directory
// as defined in the config.
func collectRecipeFolders() []string {
	assertConfigLoaded()

	dirs := make([]string, 0)
	dirStack := make([]string, 0)
	dirStack = append(dirStack, config.Get(config.KeyRecipeDir))

	for len(dirStack) > 0 {
		// Pop the top of the stack
		head := dirStack[len(dirStack)-1]
		dirStack = dirStack[:len(dirStack)-1]
		// Push head into output dirs
		dirs = append(dirs, head)

		// Read contents of head
		entries, err := os.ReadDir(head)
		if err != nil {
			errstr := fmt.Sprintf("Error: %v\n", err.Error())
			os.Stderr.WriteString(errstr)
		}

		// Add each subdir to the stack
		for _, entry := range entries {
			if entry.IsDir() {
				dirStack = append(dirStack, filepath.Join(head, entry.Name()))
			}
		}
	}

	return dirs
}

// Returns a JSON list of all recipe names (descended recursively) in the Recipes
// folder, as defined in the config.
//
// Recipes are represented as a path relative to the recipes root.
func GetAllRecipeNames() string {
	assertConfigLoaded()
	root := config.Get(config.KeyRecipeDir)

	recipes := make([]string, 0)
	recipeDirs := collectRecipeFolders()
	for _, dir := range recipeDirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			errstr := fmt.Sprintf("Error: %v\n", err.Error())
			os.Stderr.WriteString(errstr)
		}

		for _, entry := range entries {
			if entry.IsDir() || filepath.Ext(entry.Name()) != ".cook" {
				continue
			}

			recipeName := strings.TrimSuffix(
				filepath.Join(dir, entry.Name())[len(root)+1:],
				".cook")
			recipes = append(recipes, recipeName)
		}
	}
	jsonBytes, err := json.Marshal(recipes)
	if err != nil {
		errstr := fmt.Sprintf("Error: %v\n", err.Error())
		os.Stderr.WriteString(errstr)
	}

	return string(jsonBytes)
}

// Returns recipe data by the specified filepath (without extension)
//
// `name` is a relative filepath from recipe root e.g.
//
//	"breakfast/eggs_benedict"
//
// This can only read `.cook` files within the recipe directory.
func GetRecipe(name string) string {
	assertConfigLoaded()

	root := config.Get(config.KeyRecipeDir)
	fullpath := filepath.Join(root, name+".cook")

	// Expand relative filepaths to prevent reading outside of root
	fullpath, err := filepath.Abs(fullpath)
	if err != nil {
		common.ShowError(err)
		return ""
	}
	// Ensure expanded path is within root
	if !strings.HasPrefix(fullpath, root) {
		os.Stderr.WriteString("Error: Attempted to access outside data directory\n")
		return ""
	}

	if !common.FileExists(fullpath) {
		return ""
	}

	r := recipe.LoadFromFile(fullpath)
	r.Name = name
	jsonStr, err := json.Marshal(*r)
	if err != nil {
		common.ShowError(err)
	}
	return string(jsonStr)
}

// Deletes a recipe based on its relative file path from the recipe root
// folder as defined in the config file.
// e.g.
//
//	"breakfast/eggs_benedict"
//
// # Returns true on success, false on delete
//
// This can only delete `.cook` files within the recipe directory.
func DeleteRecipe(name string) bool {
	assertConfigLoaded()

	root := config.Get(config.KeyRecipeDir)
	fullpath := filepath.Join(root, name+".cook")

	// Expand relative filepaths to prevent reading outside of root
	fullpath, err := filepath.Abs(fullpath)
	if err != nil {
		common.ShowError(err)
		return false
	}
	// Ensure expanded path is within root
	if !strings.HasPrefix(fullpath, root) {
		os.Stderr.WriteString("Error: Attempted to access outside data directory\n")
		return false
	}

	if !common.FileExists(fullpath) {
		return false
	}

	// Delete the file
	err = os.Remove(fullpath)
	if err != nil {
		common.ShowError(err)
		return false
	}

	return true
}

// TODO
// Returns a JSON list containing the `n`th page of recipes, using the
// provided page size.
func GetRecipesPage(n int, size int) string {
	return ""
}

// TODO
// Returns a JSON list containing all recipes.
func GetAllRecipes() string {
	return ""
}

// TODO
// Returns a JSON list containing all recipes with names that match using
// the provided regex pattern.
func GetRecipesMatching(regexPattern string) string {
	return ""
}
