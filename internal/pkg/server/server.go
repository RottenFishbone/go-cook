package server

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"git.sr.ht/~rottenfishbone/go-cook"
	"git.sr.ht/~rottenfishbone/go-cook/api"
	"git.sr.ht/~rottenfishbone/go-cook/internal/pkg/common"
	"git.sr.ht/~rottenfishbone/go-cook/internal/web"
)

// The manifest of each API endpoint mapped to its handler
var apiHandlerFuncs = map[string]func(http.ResponseWriter, *http.Request){
	"recipes/parse": apiRecipeParse,
	"recipes/names": apiRecipeNames,
	"recipes/":      apiRecipe,
}

func Start(port int, onlyApi bool) {
	if port < 0 || port > 65535 {
		panic("Attempted to start server with invalid port")
	}

	// Iterate over handler map and add to server
	for k, v := range apiHandlerFuncs {
		http.HandleFunc("/api/0/"+k, v)
	}

	// Host webserver, if requested
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

// Handles request to the recipes root endpoint
func apiRecipe(w http.ResponseWriter, r *http.Request) {
	// Pull the name from the request
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Malformed Query, missing `name` parameter.", http.StatusUnprocessableEntity)
		return
	}

	apiRecipeByName(name, w, r)
}

// Handles requests for recipe name lists (these are preferred to full recipes for lists).
//
// Only accepts GET requests.
//
// params:
//   - q <query> -- searches recipes
//   - page <uint> -- `n`th page of size `count`
//   - count <uint> -- size of each page
func apiRecipeNames(w http.ResponseWriter, r *http.Request) {
	var err error
	var names []byte

	if r.Method != http.MethodGet {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	// Grab URL query
	q := r.URL.Query().Get("q")
	count := r.URL.Query().Get("count")
	if count == "" {
		http.Error(w, "Malformed Query, missing `count` parameter.",
			http.StatusUnprocessableEntity)
		return
	}

	page := r.URL.Query().Get("page")
	if page == "" {
		http.Error(w, "Malformed Query, missing `page` parameter.",
			http.StatusUnprocessableEntity)
		return
	}

	// Parse into numbers
	var pageNum, countNum uint64
	if pageNum, err = strconv.ParseUint(page, 10, 64); err != nil {
		http.Error(w, "Malformed Query, failed to parse `page` parameter.",
			http.StatusUnprocessableEntity)
		return
	}
	if countNum, err = strconv.ParseUint(count, 10, 64); err != nil {
		http.Error(w, "Malformed Query, failed to parse `count` parameter.",
			http.StatusUnprocessableEntity)
		return
	}

	if q == "" {
		if names, err = api.GetRecipeNamesPagedJSON(pageNum, countNum); err != nil {
			http.Error(w, "Failed to fetch page.", http.StatusInternalServerError)
			return
		}
	} else {
		if names, err = api.SearchRecipeNamesPagedJSON(q, pageNum, countNum); err != nil {
			http.Error(w, "Failed to fetch page.", http.StatusInternalServerError)
			return
		}
	}

	w.Write(names)
}
