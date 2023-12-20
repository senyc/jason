package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/senyc/jason/pkg/auth"

	"github.com/gorilla/mux"
	"github.com/senyc/jason/pkg/types"
)

var (
	noContext error = errors.New("Failure obtaining userId from context")
)

func (s *Server) getAllTasks(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	uid, ok := ctx.Value("userId").(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		panic(noContext)
	}

	tasks, err := s.db.GetAllTasksByUser(uid)
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
	ctx := req.Context()
	uid, ok := ctx.Value("userId").(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		panic(noContext)
	}
	var newTask types.NewTask

	err := json.NewDecoder(req.Body).Decode(&newTask)
	if err != nil {
		panic(err)
	}

	err = s.db.AddNewTask(newTask, uid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) getTaskById(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	ctx := req.Context()
	uid, ok := ctx.Value("userId").(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		panic(noContext)

	}

	task, err := s.db.GetTaskById(uid, vars["id"])
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

func (s *Server) getApiKey(w http.ResponseWriter, req *http.Request) {
	var userAuth types.UserLogin
	err := json.NewDecoder(req.Body).Decode(&userAuth)
	if err != nil {
		panic(err)
	}
	apiKey, err := s.db.GetApiKey(userAuth)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		panic(err)
	}
	response := types.ApiKey{ApiKey: apiKey}

	j, err := json.Marshal(response)
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
	var newUser types.User
	err := json.NewDecoder(req.Body).Decode(&newUser)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	// Generate new api key for user
	apiKey, err := auth.GetApiKey()

	// Encrypt password
	newUser.Password, err = auth.EncryptPassword(newUser.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	err = s.db.AddNewUser(newUser, apiKey)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) markAsCompleted(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	ctx := req.Context()
	uid, ok := ctx.Value("userId").(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		panic(noContext)

	}

	err := s.db.MarkTaskCompleted(uid, vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) markAsIncomplete(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	ctx := req.Context()
	uid, ok := ctx.Value("userId").(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		panic(noContext)

	}

	err := s.db.MarkTaskIncomplete(uid, vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) login(w http.ResponseWriter, req *http.Request) {
	var userAuth types.UserLogin

	err := json.NewDecoder(req.Body).Decode(&userAuth)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}
	encryptedPass, err := s.db.GetPasswordFromLogin(userAuth.Email)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		panic(err)
	}

	if err := auth.IsAuthorized(userAuth.Password, encryptedPass); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		panic(err)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
