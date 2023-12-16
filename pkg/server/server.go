package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func handleGet(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "get\n")
}

func handlePost(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "post\n")
}

func Start() {
	r := mux.NewRouter()
	r.HandleFunc("/", handleGet).Methods("GET")
	r.HandleFunc("/", handlePost).Methods("POST")

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	server.ListenAndServe()
}
