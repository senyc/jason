package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/senyc/jason/pkg/db"
	"github.com/senyc/jason/pkg/types"
)

type Server struct {
	db     *db.DB
	server *http.Server
}

func (s *Server) getAllTasks(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	tasks, err := s.db.GetAllTasksByUser(vars["userId"])
	if err != nil {
		panic(err)
	}
	j, err := json.Marshal(tasks)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(j)

	if err != nil {
		panic(err)
	}
}

func (s *Server) addNewTask(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	var newTask types.NewTask

	err := json.NewDecoder(req.Body).Decode(&newTask)
	if err != nil {
		panic(err)
	}

	err = s.db.AddNewTask(newTask, vars["userId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) getTaskById(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	task, err := s.db.GetTaskById(vars["userId"], vars["id"])
	if err != nil {
		panic(err)
	}
	j, err := json.Marshal(task)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(j)

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
