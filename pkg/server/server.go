package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/senyc/jason/pkg/db"
)

type Server struct {
	db     *db.DB
	server *http.Server
}

func (s *Server) getAllTasks(w http.ResponseWriter, req *http.Request) {
	s.db.Query()
}

func (s *Server) Start() error {
	s.db = new(db.DB)

	err := s.db.Connect()

	if err != nil {
		return err
	}

	r := mux.NewRouter()
	r.HandleFunc("/", s.getAllTasks).Methods("GET")

	s.server = &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	err = s.server.ListenAndServe()
	return err
}

func (s *Server) Shutdown() error {
	// also close the db from here
	fmt.Println("shutting down")
	return s.server.Close()
}
