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
	site := r.PathPrefix("/site/user/").Subrouter()

	tasks.Use(s.autorizationMiddleware)
	tasks.HandleFunc("/all", s.getAllTasks).Methods(http.MethodGet)
	tasks.HandleFunc("/byId/{id}", s.getTaskById).Methods(http.MethodGet)
	tasks.HandleFunc("/markComplete/{id}", s.markAsCompleted).Methods(http.MethodPatch)
	tasks.HandleFunc("/markIncomplete/{id}", s.markAsIncomplete).Methods(http.MethodPatch)
	tasks.HandleFunc("/new", s.addNewTask).Methods(http.MethodPost)

	site.Use(s.jwtAuthorizationMiddleware)
	site.HandleFunc("/all", s.getAllTasks).Methods(http.MethodGet)
	site.HandleFunc("/byId/{id}", s.getTaskById).Methods(http.MethodGet)
	site.HandleFunc("/markComplete/{id}", s.markAsCompleted).Methods(http.MethodPatch)
	site.HandleFunc("/markIncomplete/{id}", s.markAsIncomplete).Methods(http.MethodPatch)
	site.HandleFunc("/new", s.addNewTask).Methods(http.MethodPost)

	user.HandleFunc("/new", s.addNewUser).Methods(http.MethodPost)
	user.HandleFunc("/login", s.login).Methods(http.MethodPost)
	user.HandleFunc("/key/new", s.newApiKey).Methods(http.MethodPost)

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
