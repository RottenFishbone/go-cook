package server

import (
	"fmt"
	"log"
	"net/http"
)

var handlerFuncs = map[string]func(http.ResponseWriter, *http.Request){
	"/": rootHandler,
}

func Start(port int) {
	if port < 0 || port > 65535 {
		panic("Attempted to start server with invalid port")
	}

	for k, v := range handlerFuncs {
		http.HandleFunc(k, v)
	}

	fmt.Printf("Starting server at: 127.0.0.1:%v...\n", port)
	addr := ":" + fmt.Sprint(port)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Henlo!")
}
