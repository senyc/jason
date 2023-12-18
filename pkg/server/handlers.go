package server

import (
	"encoding/json"
	"net/http"
	"github.com/senyc/jason/pkg/auth"

	"github.com/gorilla/mux"
	"github.com/senyc/jason/pkg/types"
)

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

func (s *Server) addNewUser(w http.ResponseWriter, req *http.Request) {
	var newUser types.NewUser
	err := json.NewDecoder(req.Body).Decode(&newUser)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}
	err, apiKey := auth.GetApiKey()

	err = s.db.AddNewUser(newUser, apiKey)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) markAsCompleted(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	err := s.db.MarkTaskCompleted(vars["userId"], vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) markAsIncomplete(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	err := s.db.MarkTaskIncomplete(vars["userId"], vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
}
