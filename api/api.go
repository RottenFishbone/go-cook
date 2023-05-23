package api

// TODO cache recipes in memory

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"git.sr.ht/~rottenfishbone/go-cook"
	"git.sr.ht/~rottenfishbone/go-cook/internal/pkg/common"
	"git.sr.ht/~rottenfishbone/go-cook/pkg/config"
)

// Ensures the config file is loaded. Panic on failure.
//
// An unloaded read should never occur at runtime, unless there is a bug.
func assertConfigLoaded() {
	if !config.IsLoaded() {
		panic("Attempted to read unloaded config")
	}
}

// Ensures a passed recipe name is `.cook` (or no ext) and within the
// recipe's directory.
func sanitizeRecipeName(path string) (string, error) {
	var err error
	rootDir := config.GetConfig().Recipe.Dir

	// Sanitize extension
	if path, err = sanitizeExt(path, ".cook"); err != nil {
		return "", err
	}
	// Sanitize relative paths
	if path, err = sanitizeRelPath(rootDir, path); err != nil {
		return "", err
	}

	// Path is now an absolute path with a `.cook` extension
	return path, nil
}

// Slices array using dynamically defined page sizes
func page[T any](arr []T, page uint64, size uint64) []T {
	start := page * size
	end := page*size + size
	cnt := uint64(len(arr))

	if cnt < start {
		return make([]T, 0)
	} else if cnt < end {
		return arr[start:]
	} else {
		return arr[start:end]
	}
}

// Ensures a relative path is not outside of a specified root directory
//
//   - Returns absolute path on success
//   - Returns empty string and an error if an illegal path is provided
func sanitizeRelPath(root string, path string) (string, error) {
	absPath, err := filepath.Abs(filepath.Join(root, path))
	if err != nil {
		return "", err
	}

	if !strings.HasPrefix(absPath, root) {
		errMsg :=
			fmt.Sprintf("Relative path escapes root directory:\n%s\n%s", absPath, root)
		return "", errors.New(errMsg)
	}

	return absPath, nil
}

// Ensures a file extension is either `ext` or empty (to which `ext` is added)
//
// Returns empty string and an error in every other case
func sanitizeExt(path string, ext string) (string, error) {
	curExt := filepath.Ext(path)

	if curExt == "" {
		return path + ext, nil
	} else if curExt == ext {
		return path, nil
	} else {
		errMsg := fmt.Sprintf("Illegal file extension: %s", curExt)
		return "", errors.New(errMsg)
	}
}

// Collects the name of each (nested) subdirectory (and root) of the recipes directory
// as defined in the config.
func collectRecipeFolders() []string {
	assertConfigLoaded()

	// Use a stack to descend into possibly subfolders
	dirs := make([]string, 0)
	dirStack := make([]string, 0)
	recipeDir := config.GetConfig().Recipe.Dir
	dirStack = append(dirStack, recipeDir)

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

// Gathers all .cook recipes within the recipe root (possibly nested) and returns
// each as a relative filepath from recipe root.
//
// Returns (nil, err) on failure
func GetAllRecipeNames() ([]string, error) {
	//TODO use watcher to cache results and update on changes
	assertConfigLoaded()
	var err error

	// Scan every (sub)directory in the recipe dir for .cook files
	root := config.GetConfig().Recipe.Dir
	recipes := make([]string, 0)
	recipeDirs := collectRecipeFolders()
	for _, dir := range recipeDirs {
		var entries []fs.DirEntry
		if entries, err = os.ReadDir(dir); err != nil {
			return nil, err
		}

		for _, entry := range entries {
			if entry.IsDir() || filepath.Ext(entry.Name()) != ".cook" {
				continue
			}

			// Rebuild the relative path of the recipe (from recipe root)
			recipeName := strings.TrimSuffix(
				filepath.Join(dir, entry.Name())[len(root)+1:],
				".cook")
			recipes = append(recipes, recipeName)
		}
	}

	return recipes, nil
}

// Returns a JSON list of all recipe names (descended recursively) in the Recipes
// folder, as defined in the config.
//
// Recipes are represented as a path relative to the recipes root.
func GetAllRecipeNamesJSON() ([]byte, error) {
	var err error
	var recipes []string
	if recipes, err = GetAllRecipeNames(); err != nil {
		return nil, err
	}

	var jsonBytes []byte
	if jsonBytes, err = json.Marshal(recipes); err != nil {
		return nil, err
	}

	return jsonBytes, nil
}

// Returns raw byte data of recipe at provided relative filepath.
//
// e.g. "breakfast/eggs_benedict"
//
//		Considerations:
//		- byte array will be nil on failure
//	 	- This can only read `.cook` files within the recipe directory.
func GetRecipeSource(name string) ([]byte, error) {
	assertConfigLoaded()
	var err error
	var path string

	// Sanitize inputs
	if path, err = sanitizeRecipeName(name); err != nil {
		return nil, err
	}

	// Existence check
	if !common.FileExists(path) {
		return nil, errors.New("File not found.")
	}

	var raw []byte
	if raw, err = os.ReadFile(path); err != nil {
		return nil, err
	}

	return raw, nil
}

// Returns parsed and json econded recipe at provided relative filepath.
//
// e.g. "breakfast/eggs_benedict"
//
//		Considerations:
//		- byte array will be nil on failure
//	 	- This can only read `.cook` files within the recipe directory.
func GetRecipeJSON(name string) ([]byte, error) {
	var err error
	var raw []byte

	if raw, err = GetRecipeSource(name); err != nil {
		return nil, err
	}

	r := cook.ParseRecipe(name, &raw)
	var jsonData []byte
	if jsonData, err = json.Marshal(&r); err != nil {
		return nil, err
	}
	return jsonData, nil
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

	// Sanitize inputs
	if name, err = sanitizeRecipeName(name); err != nil {
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

	// Sanitize inputs
	if name, err = sanitizeRecipeName(name); err != nil {
		return err
	} else if target, err = sanitizeRecipeName(target); err != nil {
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
	rootDir := config.GetConfig().Recipe.Dir
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

// Creates a *new* recipe using the provided relative filepath.
//
// e.g. "breakfast/eggs_benedict"
//
//	Considerations:
//	- Returns nil on success, forwards error on failure
//	- Overwrites will return an error
//	- This can only create `.cook` files *within* the recipe directory.
func CreateRecipe(name string, contents *[]byte) error {
	var err error
	var file *os.File

	// Sanitize inputs
	if name, err = sanitizeRecipeName(name); err != nil {
		return err
	}

	// Prevent Overwrite
	if common.FileExists(name) {
		return errors.New("File already exists.")
	}

	// Create the new recipe file
	if file, err = os.Create(name); err != nil {
		return err
	}
	defer file.Close()

	// Populate with contents
	if _, err := file.Write(*contents); err != nil {
		return err
	}

	return nil
}

// Deletes a recipe based on its relative file path from the recipe root
// folder as defined in the config file.
//
// e.g. "breakfast/eggs_benedict"
//
//		Considerations:
//		- Returns nil on success, forwards error on failure
//		- This can only delete `.cook` files *within* the recipe directory.
//	 	- Deletes empty parent directories (automatic housekeeping)
func DeleteRecipe(name string) error {
	assertConfigLoaded()
	var err error
	var path string

	// Sanitize inputs
	if path, err = sanitizeRecipeName(name); err != nil {
		return err
	}

	// Existence check
	if !common.FileExists(path) {
		errMsg := fmt.Sprintf("File not found: %s", name+".cook")
		return errors.New(errMsg)
	}

	// Delete the file
	if err = os.Remove(path); err != nil {
		return err
	}

	// Try to clear empty directories, if the OS didn't
	rootDir := config.GetConfig().Recipe.Dir
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

// Builds a JSON list containing the `n`th page of recipes, using the
// provided page size.
//
// # Returns (jsonBytes, nil) on success, (nil, err) on failure
//
// Out of bounds requests will return an empty array.
func GetRecipeNamesPagedJSON(n uint64, size uint64) ([]byte, error) {
	var err error
	var recipes []string

	if recipes, err = GetAllRecipeNames(); err != nil {
		return nil, err
	}
	recipes = page(recipes, n, size)

	var jsonBytes []byte
	if jsonBytes, err = json.Marshal(recipes); err != nil {
		return nil, err
	}

	return jsonBytes, nil
}

// Returns a string list containing all recipes with names that match using
// the provided search pattern.
//
// Search patterns are defined as:
//   - comma delimited terms
//   - each term is a tag (except for filters)
//   - directory filters can be written with a colon. e.g. `:breakfast/sweet`
//
// each tag defined in the source under the metadata `tag` (comma delimited) are
// searched for.
func SearchRecipeNames(query string) ([]string, error) {
	// TODO optimize entire search system, as this should be fast
	// TODO in-recipe tags
	// TODO handle spaces in directory filters
	assertConfigLoaded()
	var err error
	var recipes []string
	query = strings.ToLower(query)
	queryTerms := strings.Split(query, ",")
	terms := make([]string, 0)

	// Find a directory filter (if it exists) and build list of search tags
	filter := ""
	for _, term := range queryTerms {
		term = strings.TrimSpace(term)
		if term == "" {
			continue
		}

		// Take any term that starts with ':' as filter
		// given duplicates, the last one will be taken as the filter
		if strings.HasPrefix(term, ":") {
			filter = term[1:]
			continue
		}

		terms = append(terms, term)
	}

	// Grab the recipe list
	if recipes, err = GetAllRecipeNames(); err != nil {
		return nil, err
	}

	matching := make([]string, 0)
	for _, rawRecipe := range recipes {
		// Normalize all recipe names
		recipe := filepath.ToSlash(strings.ToLower(rawRecipe))
		// Apply filter
		if filter != "" {
			if !strings.HasPrefix(recipe, filter) {
				continue
			}
		}

		// Build a map of tags so we don't have to iterate over it repeatedly
		recipeChunks := strings.Split(recipe, "/")
		name := recipeChunks[len(recipeChunks)-1]
		nameChunks := strings.Split(name, "_")
		recipeTags := append(recipeChunks, nameChunks...)
		tagMap := map[string]bool{}
		for _, tag := range recipeTags {
			tagMap[tag] = true
		}

		// Assert all search tags exist on recipe
		for _, searchTag := range terms {
			if _, found := tagMap[searchTag]; !found {
				// Check for partial name match otherwise
				if !strings.Contains(name, searchTag) {
					// Skip if no matches can be made
					goto outerContinue
				}
			}
		}
		matching = append(matching, rawRecipe)
	outerContinue:
	}

	return matching, nil
}

// Returns a JSON list containing all recipes with names that match using
// the provided search pattern.
//
// See `SearchRecipeNames()` for details.
func SearchRecipeNamesJSON(query string) ([]byte, error) {
	var err error
	var recipes []string

	if recipes, err = SearchRecipeNames(query); err != nil {
		return nil, err
	}

	var jsonBytes []byte
	if jsonBytes, err = json.Marshal(recipes); err != nil {
		return nil, err
	}

	return jsonBytes, nil
}

// Returns a JSON list containing all recipes with names that match using
// the provided search pattern.
//
// See `SearchRecipeNames` for details
func SearchRecipeNamesPagedJSON(query string, n uint64, count uint64) ([]byte, error) {
	var err error
	var recipes []string

	if recipes, err = SearchRecipeNames(query); err != nil {
		return nil, err
	}
	recipes = page(recipes, n, count)

	var jsonBytes []byte
	if jsonBytes, err = json.Marshal(recipes); err != nil {
		return nil, err
	}

	return jsonBytes, nil
}
