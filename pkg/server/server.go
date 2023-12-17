package server

import (
	"fmt"
	"net/http"
	"github.com/senyc/jason/pkg/db"
	"github.com/gorilla/mux"
)

func getAllTasks(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "get you did it\n")
	db.Connect()
}

func handlePost(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "post\n")
}

func Start() {
	r := mux.NewRouter()
	r.HandleFunc("/", getAllTasks).Methods("GET")
	r.HandleFunc("/", handlePost).Methods("POST")

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
