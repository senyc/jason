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

    tasks := r.PathPrefix("/api/tasks/").Subrouter()
    user := r.PathPrefix("/api/user/").Subrouter()

	tasks.Use(s.autorizationMiddleware)
	tasks.HandleFunc("/all", s.getAllTasks).Methods("GET")
	tasks.HandleFunc("/byId/{id}", s.getTaskById).Methods("GET")
	tasks.HandleFunc("/markComplete/{id}", s.markAsCompleted).Methods("PATCH")
	tasks.HandleFunc("/markIncomplete/{id}", s.markAsIncomplete).Methods("PATCH")
	tasks.HandleFunc("/new", s.addNewTask).Methods("POST")

	user.HandleFunc("/new", s.addNewUser).Methods("POST")
	user.HandleFunc("/getApiKey", s.getApiKey).Methods("GET")

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
