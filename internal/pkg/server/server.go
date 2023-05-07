package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"git.sr.ht/~rottenfishbone/go-cook"
	"git.sr.ht/~rottenfishbone/go-cook/api"
	"git.sr.ht/~rottenfishbone/go-cook/internal/pkg/common"
	"git.sr.ht/~rottenfishbone/go-cook/internal/web"
	"git.sr.ht/~rottenfishbone/go-cook/pkg/config"
)

// The manifest of each API endpoint mapped to its handler
var apiHandlerFuncs = map[string]func(http.ResponseWriter, *http.Request){
	"recipes/parse":  apiRecipeParse,
	"recipes/names":  apiRecipeNames,
	"recipes/byName": apiRecipeByName,
}

func Start(port int, onlyApi bool) {
	if port < 0 || port > 65535 {
		panic("Attempted to start server with invalid port")
	}

	for k, v := range apiHandlerFuncs {
		http.HandleFunc("/api/0/"+k, v)
	}

	if !onlyApi {
		// Ensure the webserver was correctly embedded at compile time
		entries, _ := fs.ReadDir(web.WebDist, "dist")
		if len(entries) == 1 {
			os.Stderr.WriteString("Webapp was not compiled alongside 'cook', see build instructions to enable the server.\n")
			os.Exit(1)
		}

		// Fetch the web server files and serve 'dist' as root
		filesys, _ := fs.Sub(web.WebDist, "dist")
		http.Handle("/", http.FileServer(http.FS(filesys)))
	}

	fmt.Printf("Starting server at: 127.0.0.1:%v...\n", port)
	addr := "0.0.0.0:" + fmt.Sprint(port)
	log.Fatal(http.ListenAndServe(addr, nil))
}

// Handles parse requests as POST bodies.
func apiRecipeParse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	// Read the body from the POST
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request.", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Parse body and encode to JSON
	recipe := cook.ParseRecipe("", &body)
	jsonBytes, jsonErr := json.Marshal(&recipe)
	if jsonErr != nil {
		http.Error(w, "Error parsing recipe text", http.StatusBadRequest)
		return
	}

	// Write the recipe back
	_, err = w.Write(jsonBytes)
	if err != nil {
		common.ShowError(err)
	}
}

// Handles requests for recipe name lists (these are preferred to full recipes for lists).
func apiRecipeNames(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	// Pull the name from the request
	page := r.URL.Query().Get("page")
	count := r.URL.Query().Get("count")
	pageNum, pErr := strconv.ParseUint(page, 10, 64)
	countNum, cErr := strconv.ParseUint(count, 10, 64)

	if count == "" {
		http.Error(w, "Malformed Query, missing `count` parameter.", http.StatusUnprocessableEntity)
		return
	} else if countNum == 0 {
		fmt.Fprintf(w, "%s", api.GetAllRecipeNames())
		return
	} else if page == "" {
		http.Error(w, "Malformed Query, missing `page` parameter.", http.StatusUnprocessableEntity)
		return
	} else {
		if pErr != nil {
			http.Error(w, "Malformed Query, failed to parse `page` parameter.", http.StatusUnprocessableEntity)
			return
		}
		if cErr != nil {
			http.Error(w, "Malformed Query, failed to parse `count` parameter.", http.StatusUnprocessableEntity)
			return
		}

		// TODO handle paginated requests
		var _ = pageNum
	}
}

// Handles requests to get/change individual recipes via the `name` URL parameter
//   - GET :returns the parsed recipe as JSON
//     [param `raw=<true/false>` will return the raw recipe file, unparsed.]
//   - DELETE: deletes the recipe from the server
//   - POST: update the file with the POST body as text (UNIMPL.)
//     [param `rename=<string>` will move the recipe to the passed string.
//     POST body can be left blank to simply rename (no changes)]
//   - PUT: add new recipe (fails on overwrite, use POST to overwrite) (UNIMPL.)
func apiRecipeByName(w http.ResponseWriter, r *http.Request) {
	// Pull the name from the request
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Malformed Query, missing `name` parameter.", http.StatusUnprocessableEntity)
		return
	}

	// Read the body for relevant methods
	var body []byte
	if r.Method == http.MethodPost || r.Method == http.MethodPut {
		var err error
		if body, err = ioutil.ReadAll(r.Body); err != nil {
			http.Error(w, "Error reading request.", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
	}

	// Handle the request, based on the method requested
	switch r.Method {
	case http.MethodGet:
		// Try to grab raw param if it exists
		raw := r.URL.Query().Get("raw")
		if raw != "" && raw != "true" && raw != "false" { // Can only have 1 of 3 vals
			http.Error(w, "Malformed Query, invalid `raw` parameter.", http.StatusUnprocessableEntity)
			return
		}

		var recipeData []byte
		if raw != "true" {
			recipeData = []byte(api.GetRecipe(name))
		} else {
			path := filepath.Join(config.Get(config.KeyRecipeDir), name+".cook")
			recipeData, _ = os.ReadFile(path)
		}
		if recipeData != nil && len(recipeData) > 0 {
			w.Write(recipeData)
		} else {
			http.Error(w, "File not found.", http.StatusNotFound)
			return
		}

	case http.MethodPost:
		rename := r.URL.Query().Get("rename")
		handleRecipeByNamePOST(name, rename, &body, w)

	case http.MethodPut:
		handleRecipeByNamePUT(name, &body, w)

	case http.MethodDelete:
		if api.DeleteRecipe(name) {
			w.WriteHeader(http.StatusOK)
		} else {
			http.Error(w, "File not found.", http.StatusNotFound)
			return
		}

	default:
		http.Error(w, "Method not supported.", http.StatusMethodNotAllowed)
		return
	}
}

// Helper function to handle PUT requests for endpoint `recipes/byName`
func handleRecipeByNamePUT(name string, body *[]byte, w http.ResponseWriter) {
	// TODO Implement
}

// Helper function to handle POST requests for endpoint `recipes/byName`
func handleRecipeByNamePOST(name string, rename string, body *[]byte, w http.ResponseWriter) {
	// TODO recipe validations
	var err error
	rootDir := config.Get(config.KeyRecipeDir)

	var path string
	if path, err = common.SanitizeRelPath(rootDir, name+".cook"); err != nil {
		// Cannot access outside of recipes
		common.ShowError(err)
		http.Error(w, "Attempted to access outside directory.", http.StatusBadRequest)
		return
	} else if !common.FileExists(path) {
		// Cannot modify non-existent files
		http.Error(w, "File not found.", http.StatusNotFound)
		return
	}

	if rename != "" { // Handle renaming
		if rename, err = common.SanitizeRelPath(rootDir, rename+".cook"); err != nil {
			// Prevent out of bounds rename target
			common.ShowError(err)
			http.Error(w, "Attempted to write outside directory.", http.StatusBadRequest)
			return
		} else if common.FileExists(rename) {
			// Prevent overwrites
			http.Error(w, "Attempted to overwrite.", http.StatusBadRequest)
			return
		}

		if len(*body) == 0 {
			// Simple rename
			err = os.Rename(path, rename)
			if err != nil {
				common.ShowError(err)
				http.Error(w, "Failed to rename file.", http.StatusInternalServerError)
			} else {
				w.WriteHeader(http.StatusOK)
			}
			return
		} else {
			// Rename and edit
			var file *os.File
			if file, err = os.Create(rename); err != nil {
				common.ShowError(err)
				http.Error(w, "Failed to create new file.", http.StatusInternalServerError)
				return
			}
			defer file.Close()

			if _, err = file.Write(*body); err != nil {
				http.Error(w, "Failed to write to new file.", http.StatusInternalServerError)
			} else if err = os.Remove(path); err != nil {
				http.Error(w, "Failed to remove original file.", http.StatusInternalServerError)
			} else {
				w.WriteHeader(http.StatusOK)
				return
			}

			// Show whatever error broke the chain
			common.ShowError(err)
		}
	} else { // Handle update (without rename)
		if len(*body) == 0 {
			common.ShowError(errors.New("Tried to update recipe with to empty file."))
			http.Error(w, "Body is empty.", http.StatusBadRequest)
			return
		}

		var file *os.File
		if err = os.Rename(path, path+".bak"); err != nil {
			// Backup recipe
			common.ShowError(err)
			http.Error(w, "Failed to modify file.", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		if file, err = os.Create(path); err != nil {
			http.Error(w, "Failed to create new (temp) file.", http.StatusInternalServerError)
		} else if _, err = file.Write(*body); err != nil {
			http.Error(w, "Failed to write to file.", http.StatusInternalServerError)
		} else if err = os.Remove(path + ".bak"); err != nil {
			http.Error(w, "Failed to remove backup file.", http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Show whatever error broke the chain
		common.ShowError(err)
	}
}
