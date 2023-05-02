package server

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"git.sr.ht/~rottenfishbone/go-cook/api"
	"git.sr.ht/~rottenfishbone/go-cook/internal/web"
)

// The manifest of each API endpoint mapped to its handler
var apiHandlerFuncs = map[string]func(http.ResponseWriter, *http.Request){
	"recipes/allNames": apiRecipeAllNames,
	"recipes/byName":   apiRecipeByName,
}

func Start(port int) {
	if port < 0 || port > 65535 {
		panic("Attempted to start server with invalid port")
	}

	for k, v := range apiHandlerFuncs {
		http.HandleFunc("/api/0/"+k, v)
	}

	// Fetch the web server files and serve 'dist' as root
	fs, _ := fs.Sub(web.WebDist, "dist")
	http.Handle("/", http.FileServer(http.FS(fs)))

	fmt.Printf("Starting server at: 127.0.0.1:%v...\n", port)
	addr := ":" + fmt.Sprint(port)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func apiRecipeAllNames(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "%s", api.GetAllRecipeNames())
}

func apiRecipeByName(w http.ResponseWriter, r *http.Request) {
	// Pull the name from the request
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Malformed Query, missing `name` parameter.", http.StatusUnprocessableEntity)
	}

	// Handle the request, based on the method requested
	switch r.Method {
	case "GET":
		recipeStr := api.GetRecipe(name)
		if recipeStr != "" {
			w.Write([]byte(recipeStr))
		} else {
			http.Error(w, "File not found.", http.StatusNotFound)
		}
	case "DELETE":
		if api.DeleteRecipe(name) {
			w.WriteHeader(http.StatusOK)
		} else {
			http.Error(w, "File not found.", http.StatusNotFound)
		}
	default:
		http.Error(w, "Method not supported.", http.StatusMethodNotAllowed)
	}
}
