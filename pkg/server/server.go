package server

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/senyc/jason/pkg/db"
)

type Server struct {
	db     *db.DB
	server *http.Server
	logger *log.Logger
}

func (s *Server) Start() error {
	s.db = new(db.DB)

	err := s.db.Connect()

	if err != nil {
		return err
	}
	s.logger = log.New(os.Stdout, "log:", log.LstdFlags|log.Lshortfile)

	r := mux.NewRouter()

	tasks := r.PathPrefix("/api/tasks/").Subrouter()
	user := r.PathPrefix("/api/user/").Subrouter()
	site := r.PathPrefix("/site/tasks/").Subrouter()

	tasks.Use(s.autorizationMiddleware)
	tasks.HandleFunc("/all", s.getAllTasks).Methods(http.MethodGet)
	tasks.HandleFunc("/complete", s.getCompletedTasks).Methods(http.MethodGet)
	tasks.HandleFunc("/incomplete", s.getIncompleteTasks).Methods(http.MethodGet)
	tasks.HandleFunc("/byId", s.getTaskById).Methods(http.MethodGet)
	tasks.HandleFunc("/markComplete", s.markAsCompleted).Methods(http.MethodPatch)
	tasks.HandleFunc("/markIncomplete", s.markAsIncomplete).Methods(http.MethodPatch)
	tasks.HandleFunc("/new", s.addNewTask).Methods(http.MethodPost)
	tasks.HandleFunc("/delete", s.deleteTask).Methods(http.MethodDelete)
	tasks.HandleFunc("/edit", s.editTask).Methods(http.MethodPatch)

	site.Use(s.jwtAuthorizationMiddleware)
	site.HandleFunc("/all", s.getAllTasks).Methods(http.MethodGet)
	site.HandleFunc("/complete", s.getCompletedTasks).Methods(http.MethodGet)
	site.HandleFunc("/incomplete", s.getIncompleteTasks).Methods(http.MethodGet)
	site.HandleFunc("/byId", s.getTaskById).Methods(http.MethodGet)
	site.HandleFunc("/markComplete", s.markAsCompleted).Methods(http.MethodPatch)
	site.HandleFunc("/markIncomplete", s.markAsIncomplete).Methods(http.MethodPatch)
	site.HandleFunc("/new", s.addNewTask).Methods(http.MethodPost)
	site.HandleFunc("/delete", s.deleteTask).Methods(http.MethodDelete)
	site.HandleFunc("/edit", s.editTask).Methods(http.MethodPatch)
	site.HandleFunc("/getEmail", s.getEmail).Methods(http.MethodGet)

	user.HandleFunc("/new", s.addNewUser).Methods(http.MethodPost)
	user.HandleFunc("/login", s.login).Methods(http.MethodPost)
	user.HandleFunc("/key/new", s.newApiKey).Methods(http.MethodPost)

	originsOk := handlers.AllowedOrigins([]string{"*"})
	headersOk := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "PATCH", "DELETE"})

	s.server = &http.Server{
		Addr:    ":8080",
		Handler: handlers.CORS(originsOk, headersOk, methodsOk)(r),
	}

	err = s.server.ListenAndServe()
	return err
}

func (s *Server) Shutdown() error {
	// also close the db from here
	return s.server.Close()
}
