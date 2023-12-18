package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/senyc/jason/pkg/db"
)

type Server struct {
	db     *db.DB
	server *http.Server
}

func (s *Server) Start() error {
	s.db = new(db.DB)

	err := s.db.Connect()

	if err != nil {
		return err
	}

	r := mux.NewRouter()
	r.HandleFunc("/{userId}", s.getAllTasks).Methods("GET")
	r.HandleFunc("/{userId}/{id}", s.getTaskById).Methods("GET")
	r.HandleFunc("/newTask/{userId}", s.addNewTask).Methods("POST")

	s.server = &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	err = s.server.ListenAndServe()
	return err
}

func (s *Server) Shutdown() error {
	// also close the db from here
	return s.server.Close()
}
