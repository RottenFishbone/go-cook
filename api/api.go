package api

// TODO cache recipes in memory

import (
	"encoding/json"
	"errors"
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
	// TODO conform to new API style
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

// Replaces contents of specified recipe with `contents`
// Recipe is specified as a relative filepath from directory root
//
// e.g. "breakfast/eggs_benedict"
//
//	Considerations:
//	- Returns nil on success, otherwise forwards errors.
//	- Can only affect ".cook" files
//	- Refuses to write empty files
func UpdateRecipe(name string, contents *[]byte) error {
	if len(*contents) == 0 {
		return errors.New("Cannot update recipe to be empty file.")
	}
	assertConfigLoaded()
	var err error

	// Sanitize input
	rootDir := config.Get(config.KeyRecipeDir)
	if name, err = common.SanitizeRelPath(rootDir, name+".cook"); err != nil {
		return err
	}

	// Existence check
	if !common.FileExists(name) {
		return errors.New("File does not exist.")
	}

	var file *os.File
	if err = os.Rename(name, name+".bak"); err != nil {
		// ^ Backup original
		return err
	} else if file, err = os.Create(name); err != nil {
		// ^ Create new in original's place
		recoveryErr := os.Rename(name+".bak", name)
		if recoveryErr != nil {
			common.ShowError(err)
			common.ShowError(recoveryErr)
		}
		return err
	}
	defer file.Close()
	if _, err := file.Write(*contents); err != nil {
		// ^ Write bytes to new
		undoErr := os.Remove(name)
		recoveryErr := os.Rename(name+".bak", name)
		if undoErr != nil || recoveryErr != nil {
			common.ShowError(err)
			common.ShowError(undoErr)
			common.ShowError(recoveryErr)
		}
		return err
	}

	// Remove the backup, error is non-fatal
	if err = os.Remove(name + ".bak"); err != nil {
		common.ShowError(err)
	}
	return nil
}

// Rename a recipe file using relative filepaths as input.
//
// e.g. "breakfast/eggs_benedict" -> "lunch/deluxe_eggs_benedict"
//
//	Considerations:
//	- Returns nil on success, otherwise returns encountered error.
//	- Only affects ".cook" files within the recipe directory.
//	- Cleans up empty directories after a rename
func RenameRecipe(name string, target string) error {
	assertConfigLoaded()
	var err error

	// Sanitize both inputs
	rootDir := config.Get(config.KeyRecipeDir)
	if name, err = common.SanitizeRelPath(rootDir, name+".cook"); err != nil {
		return err
	} else if target, err = common.SanitizeRelPath(rootDir, target+".cook"); err != nil {
		return err
	}

	// Existence checks
	if !common.FileExists(name) {
		return errors.New("File does not exist.")
	} else if common.FileExists(target) {
		// No overwrite
		return errors.New("Target already exists.")
	}

	// Ensure target directory exists
	dir := filepath.Dir(target)
	if !common.FileExists(filepath.Dir(target)) {
		if err = os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}

	// Finally, we can rename
	if err = os.Rename(name, target); err != nil {
		return err
	}

	// Try to clear empty directories, if the OS didn't
	oldDir := strings.TrimPrefix(rootDir, filepath.Dir(name))
	dirs := strings.Split(filepath.ToSlash(oldDir), "/")
	if len(dirs) > 0 {
		oldDir = filepath.Join(rootDir, dirs[0])
	}
	if common.FileExists(oldDir) {
		if err = common.CleanupEmptyDir(oldDir); err != nil {
			common.ShowError(err)
			return nil // Client can ignore error, its just housekeeping
		}
	}

	return nil
}

// Deletes a recipe based on its relative file path from the recipe root
// folder as defined in the config file.
//
// e.g. "breakfast/eggs_benedict"
//
//	Considerations:
//	- Returns true on success, false on delete
//	- This can only delete `.cook` files within the recipe directory.
func DeleteRecipe(name string) bool {
	// TODO conform to new API style
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
