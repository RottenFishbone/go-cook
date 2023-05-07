package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"git.sr.ht/~rottenfishbone/go-cook/api"
	"git.sr.ht/~rottenfishbone/go-cook/pkg/config"
)

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
		handleRecipeByNameGET(name, raw, w)
	case http.MethodPost:
		rename := r.URL.Query().Get("rename")
		handleRecipeByNamePOST(name, rename, &body, w)
	case http.MethodPut:
		handleRecipeByNamePUT(name, &body, w)
	case http.MethodDelete:
		if api.DeleteRecipe(name) {
			return
		} else {
			http.Error(w, "File not found.", http.StatusNotFound)
			return
		}
	default:
		http.Error(w, "Method not supported.", http.StatusMethodNotAllowed)
		return
	}
}

// Helper function to hangleGET requests for endpoint `recipes/byName`
func handleRecipeByNameGET(name string, raw string, w http.ResponseWriter) {
	var recipeData []byte
	if raw != "" && raw != "true" && raw != "false" { // Can only have 1 of 3 vals
		http.Error(w, "Malformed Query, invalid `raw` parameter.", http.StatusUnprocessableEntity)
		return
	}
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
}

// Helper function to handle PUT requests for endpoint `recipes/byName`
func handleRecipeByNamePUT(name string, body *[]byte, w http.ResponseWriter) {
	// TODO Implement
}

// Helper function to handle POST requests for endpoint `recipes/byName`
func handleRecipeByNamePOST(name string, rename string, body *[]byte, w http.ResponseWriter) {
	// TODO recipe validations
	var err error

	// If rename is defined, rename the recipe file and set name to rename
	if rename != "" {
		if err = api.RenameRecipe(name, rename); err != nil {
			errMsg := fmt.Sprintf("Failed to rename: %s", err)
			http.Error(w, errMsg, http.StatusInternalServerError)
			return
		}
		name = rename
	}

	if body != nil && len(*body) > 0 {
		// Update contents with provided body
		if err = api.UpdateRecipe(name, body); err != nil {
			errMsg := fmt.Sprintf("Failed to update: %s", err)
			http.Error(w, errMsg, http.StatusInternalServerError)
			return
		}
	}
}
