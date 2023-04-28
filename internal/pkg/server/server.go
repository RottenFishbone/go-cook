package server

import (
	"fmt"
	"log"
	"net/http"
	"io/fs"
	
	"git.sr.ht/~rottenfishbone/go-cook/web"
)

var handlerFuncs = map[string]func(http.ResponseWriter, *http.Request){
}

func Start(port int) {
	if port < 0 || port > 65535 {
		panic("Attempted to start server with invalid port")
	}

	for k, v := range handlerFuncs {
		http.HandleFunc(k, v)
	}
	
	// Fetch the web server files and serve 'dist' as root
	fs, _ := fs.Sub(web.WebDist, "dist")
	http.Handle("/", http.FileServer(http.FS(fs)))

	// Serve the API as well
	http.HandleFunc("/api", apiHandler)

	fmt.Printf("Starting server at: 127.0.0.1:%v...\n", port)
	addr := ":" + fmt.Sprint(port)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Henlo!")
}
