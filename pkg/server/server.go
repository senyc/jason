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
	vars := mux.Vars(req)
	_, err := s.db.GetAllTasksByUser(vars["userId"])
	if err != nil {
		panic(err)
	}
}

func (s *Server) getTaskById(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	task, err := s.db.GetTaskById(vars["userId"], vars["id"])
	if err != nil {
		panic(err)
	}
	_, err = fmt.Fprintf(w, fmt.Sprint(task))
	if err != nil {
		panic(err)
	}
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
