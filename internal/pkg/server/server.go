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
	"recipes/parse":  apiRecipeParse,
	"recipes/names":  apiRecipeNames,
	"recipes/byName": apiRecipeByName, // `by_name.go`
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
		http.Error(w, "Malformed Query, missing `count` parameter.",
			http.StatusUnprocessableEntity)
		return
	} else if countNum == 0 {
		fmt.Fprintf(w, "%s", api.GetAllRecipeNames())
		return
	} else if page == "" {
		http.Error(w, "Malformed Query, missing `page` parameter.",
			http.StatusUnprocessableEntity)
		return
	} else {
		if pErr != nil {
			http.Error(w, "Malformed Query, failed to parse `page` parameter.",
				http.StatusUnprocessableEntity)
			return
		}
		if cErr != nil {
			http.Error(w, "Malformed Query, failed to parse `count` parameter.",
				http.StatusUnprocessableEntity)
			return
		}

		// TODO handle paginated requests
		var _ = pageNum
	}
}
