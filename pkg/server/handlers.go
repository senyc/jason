package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/senyc/jason/pkg/auth"
	"github.com/senyc/jason/pkg/contact"
	"github.com/senyc/jason/pkg/db"
	"github.com/senyc/jason/pkg/dbconv"

	"github.com/senyc/jason/pkg/types"
)

var (
	noContext         error = errors.New("Failure obtaining userId from context")
	noIdFound         error = errors.New("No identification provided")
	incorrectPassword error = errors.New("Incorrect password, please try again")
	noLoginExists     error = errors.New("No login exists for this email address, please try again")
	noPasswordExists     error = errors.New("Password has been reset, please enter a new password")
)

func (s *Server) getCompletedTasks(w http.ResponseWriter, req *http.Request) {
	var res []types.CompletedTaskResponse
	ctx := req.Context()
	uuid, ok := ctx.Value("userId").(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(noContext)
	}

	completedTasks, err := s.db.GetCompletedTasks(uuid)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(err)
	}

	for _, row := range completedTasks {
		task, _ := dbconv.ToCompletedTaskResponse(row)
		res = append(res, task)
	}

	err = s.db.UpdateLastAccessedToNow(uuid)
	if err != nil {
		s.logger.Panic(err)
	}

	j, err := json.Marshal(res)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(j)

	if err != nil {
		s.logger.Panic(err)
	}
}

func (s *Server) getIncompleteTasks(w http.ResponseWriter, req *http.Request) {
	var res []types.IncompleteTaskResponse
	ctx := req.Context()
	uuid, ok := ctx.Value("userId").(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(noContext)
	}

	incompleteTasks, err := s.db.GetIncompleteTasks(uuid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(err)
	}

	for _, row := range incompleteTasks {
		task, _ := dbconv.ToIncompleteTaskResponse(row)
		res = append(res, task)
	}

	err = s.db.UpdateLastAccessedToNow(uuid)
	if err != nil {
		s.logger.Panic(err)
	}

	j, err := json.Marshal(res)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(j)

	if err != nil {
		s.logger.Panic(err)
	}
}

func (s *Server) getAllTasks(w http.ResponseWriter, req *http.Request) {
	var res []types.TaskReponse
	ctx := req.Context()
	uuid, ok := ctx.Value("userId").(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(noContext)
	}

	allTasks, err := s.db.GetAllTasksByUser(uuid)
	if err != nil {
		s.logger.Panic(err)
	}

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(err)
	}

	for _, row := range allTasks {
		task, _ := dbconv.ToTaskResponse(row)
		res = append(res, task)
	}

	j, err := json.Marshal(res)

	err = s.db.UpdateLastAccessedToNow(uuid)
	if err != nil {
		s.logger.Panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(j)

	if err != nil {
		s.logger.Panic(err)
	}
}

func (s *Server) addNewTask(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	uuid, ok := ctx.Value("userId").(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(noContext)
	}
	var newTask types.NewTaskPayload

	err := json.NewDecoder(req.Body).Decode(&newTask)
	if err != nil {
		s.logger.Panic(err)
	}
	err = s.db.AddNewTask(newTask, uuid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(err)
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) getTaskById(w http.ResponseWriter, req *http.Request) {
	var res types.TaskReponse
	id := req.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		res := types.ErrResponse{Message: noIdFound.Error()}
		j, err := json.Marshal(res)
		if err != nil {
			s.logger.Panic(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(j)
		return
	}
	ctx := req.Context()
	uuid, ok := ctx.Value("userId").(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(noContext)

	}

	task, err := s.db.GetTaskById(uuid, id)
	if err != nil {
		if err == db.NoTasksFoundError {
			w.WriteHeader(http.StatusBadRequest)

			res := types.ErrResponse{Message: err.Error()}
			j, err := json.Marshal(res)
			if err != nil {
				s.logger.Panic(err)
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(j)
			return
		} else {
			s.logger.Panic(err)
		}
	}

	res, _ = dbconv.ToTaskResponse(task)
	j, err := json.Marshal(res)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(j)

	if err != nil {
		s.logger.Panic(err)
	}
}

func (s *Server) newApiKey(w http.ResponseWriter, req *http.Request) {
	var apiKeyPayload types.ApiKeyPayload
	ctx := req.Context()
	uuid, ok := ctx.Value("userId").(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(noContext)
	}

	err := json.NewDecoder(req.Body).Decode(&apiKeyPayload)
	if err != nil {
		s.logger.Panic(err)
	}

	apiKey, err := auth.GetNewApiKey()
	if err != nil {
		s.logger.Panic(err)
	}

	err = s.db.AddApiKey(uuid, auth.EncryptApiKey(apiKey), apiKeyPayload)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		s.logger.Panic(err)
	}

	// Gets the id of the api key in case the user wants to immediately delete it
	keyMetadata, err := s.db.GetApiKeyMetadata(auth.EncryptApiKey(apiKey))
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		s.logger.Panic(err)
	}

	response := types.ApiKeyResponse{ApiKey: apiKey, ApiKeyId: keyMetadata.Id}

	j, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(j)
	if err != nil {
		s.logger.Panic(err)
	}
}

func (s *Server) addNewUser(w http.ResponseWriter, req *http.Request) {
	var newUser types.User
	err := json.NewDecoder(req.Body).Decode(&newUser)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(err)
	}

	// Encrypt password
	newUser.Password, err = auth.EncryptPassword(newUser.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(err)
	}

	err = s.db.AddNewUser(newUser)
	if err != nil {
		if err == db.NewUserUniquenessConstraintError {
			w.WriteHeader(http.StatusBadRequest)
			errResponse := types.ErrResponse{Message: err.Error()}
			j, err := json.Marshal(errResponse)
			if err != nil {
				s.logger.Panic(err)
			}
			w.Header().Set("Content-Type", "application/json")
			_, err = w.Write(j)
			if err != nil {
				s.logger.Panic(err)
			}
		}
		s.logger.Panic(err)
	}

	uuid, err := s.db.GetUuidFromEmail(newUser.Email)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		s.logger.Panic(err)
	}
	sendJwt(w, uuid)
}

func (s *Server) markAsCompleted(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		res := types.ErrResponse{Message: noIdFound.Error()}
		j, err := json.Marshal(res)
		if err != nil {
			s.logger.Panic(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(j)
		return
	}
	ctx := req.Context()
	uuid, ok := ctx.Value("userId").(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(noContext)
	}

	err := s.db.MarkTaskCompleted(uuid, id)
	if err != nil {
		if err == db.NoTasksFoundError {
			w.WriteHeader(http.StatusBadRequest)

			res := types.ErrResponse{Message: err.Error()}
			j, err := json.Marshal(res)
			if err != nil {
				s.logger.Panic(err)
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(j)
			return
		} else {
			s.logger.Panic(err)
		}
	}
}

func (s *Server) markAsIncomplete(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		res := types.ErrResponse{Message: noIdFound.Error()}
		j, err := json.Marshal(res)
		if err != nil {
			s.logger.Panic(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(j)
		return
	}
	ctx := req.Context()
	uuid, ok := ctx.Value("userId").(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(noContext)
	}

	err := s.db.MarkTaskIncomplete(uuid, id)
	if err != nil {
		if err == db.NoTasksFoundError {
			w.WriteHeader(http.StatusBadRequest)

			res := types.ErrResponse{Message: err.Error()}
			j, err := json.Marshal(res)
			if err != nil {
				s.logger.Panic(err)
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(j)
		} else {
			s.logger.Panic(err)
		}
	}
}

func sendJwt(w http.ResponseWriter, uuid string) error {
	token, err := auth.GetNewJWT(uuid)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return err
	}
	res := types.JwtResponse{Jwt: token}
	j, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
	return nil
}

func (s *Server) login(w http.ResponseWriter, req *http.Request) {
	var userAuth types.UserLoginPayload

	err := json.NewDecoder(req.Body).Decode(&userAuth)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(err)
	}
	encryptedPass, err := s.db.GetPasswordFromLogin(userAuth.Email)
	if err != nil {
		errResponse := types.ErrResponse{Message: noLoginExists.Error()}
		j, err := json.Marshal(errResponse)
		if err != nil {
			s.logger.Panic(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(j)
		return
	} else if encryptedPass == "" {
		errResponse := types.ErrResponse{Message: noPasswordExists.Error()}
		j, err := json.Marshal(errResponse)
		if err != nil {
			s.logger.Panic(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(j)
		return
	}
	
	err = auth.IsAuthorized(userAuth.Password, encryptedPass)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		errResponse := types.ErrResponse{Message: incorrectPassword.Error()}
		j, err := json.Marshal(errResponse)
		if err != nil {
			s.logger.Panic(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(j)
		return
	}
	uuid, err := s.db.GetUuidFromEmail(userAuth.Email)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		s.logger.Panic(err)
	}
	err = sendJwt(w, uuid)
	if err != nil {
		s.logger.Panic(err)
	}
}

func (s *Server) deleteTask(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		res := types.ErrResponse{Message: noIdFound.Error()}
		j, err := json.Marshal(res)
		if err != nil {
			s.logger.Panic(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(j)
		return
	}
	ctx := req.Context()
	uuid, ok := ctx.Value("userId").(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(noContext)
	}

	err := s.db.DeleteTask(uuid, id)
	if err != nil {
		if err == db.NoTasksFoundError {
			w.WriteHeader(http.StatusBadRequest)

			res := types.ErrResponse{Message: err.Error()}
			j, err := json.Marshal(res)
			if err != nil {
				s.logger.Panic(err)
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(j)
			return
		} else {
			s.logger.Panic(err)
		}
	}
}

func (s *Server) editTask(w http.ResponseWriter, req *http.Request) {
	var editPayload types.EditTaskPayload
	ctx := req.Context()
	uuid, ok := ctx.Value("userId").(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(noContext)
	}
	err := json.NewDecoder(req.Body).Decode(&editPayload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(err)
	}
	err = s.db.EditTask(uuid, editPayload)
	if err != nil {
		if err == db.NoTasksFoundError {
			w.WriteHeader(http.StatusBadRequest)

			res := types.ErrResponse{Message: err.Error()}
			j, err := json.Marshal(res)
			if err != nil {
				s.logger.Panic(err)
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(j)
		} else {
			s.logger.Panic(err)
		}
	}
}

func (s *Server) getEmail(w http.ResponseWriter, req *http.Request) {
	var emailResponse types.EmailResponse
	ctx := req.Context()
	uuid, ok := ctx.Value("userId").(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(noContext)
	}
	email, err := s.db.GetEmailAddress(uuid)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(err)
	}
	emailResponse.Email = email
	j, err := json.Marshal(emailResponse)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func (s *Server) getSyncTime(w http.ResponseWriter, req *http.Request) {
	var syncTimeResponse types.SyncTimeResponse
	ctx := req.Context()
	uuid, ok := ctx.Value("userId").(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(noContext)
	}
	syncTime, err := s.db.GetLastAccessed(uuid)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(err)
	}
	syncTimeResponse.SyncTime = syncTime
	j, err := json.Marshal(syncTimeResponse)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func (s *Server) getAccountCreationDate(w http.ResponseWriter, req *http.Request) {
	var syncTimeResponse types.AccountCreationDateResponse
	ctx := req.Context()
	uuid, ok := ctx.Value("userId").(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(noContext)
	}
	accountDate, err := s.db.GetAccountCreationDate(uuid)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(err)
	}
	syncTimeResponse.AccountCreationDate = accountDate
	j, err := json.Marshal(syncTimeResponse)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func (s *Server) getAllApiKeys(w http.ResponseWriter, req *http.Request) {
	var apiKeyMetadataResponse []types.ApiKeyMetadata
	ctx := req.Context()
	uuid, ok := ctx.Value("userId").(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(noContext)
	}
	apiKeyMetadataResponse, err := s.db.GetAllApiKeyMetadata(uuid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(err)
	}
	j, err := json.Marshal(apiKeyMetadataResponse)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func (s *Server) revokeAllApiKeys(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	uuid, ok := ctx.Value("userId").(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(noContext)
	}
	err := s.db.RevokeAllApiKeys(uuid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(err)
	}
}

func (s *Server) revokeApiKey(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		res := types.ErrResponse{Message: noIdFound.Error()}
		j, err := json.Marshal(res)
		if err != nil {
			s.logger.Panic(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(j)
		return
	}
	ctx := req.Context()
	uuid, ok := ctx.Value("userId").(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(noContext)
	}

	err := s.db.RevokeApiKey(uuid, id)
	if err != nil {
		if err == db.NoTasksFoundError {
			w.WriteHeader(http.StatusBadRequest)

			res := types.ErrResponse{Message: err.Error()}
			j, err := json.Marshal(res)
			if err != nil {
				s.logger.Panic(err)
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(j)
			return
		} else {
			s.logger.Panic(err)
		}
	}
}

func (s *Server) changeEmailAddress(w http.ResponseWriter, req *http.Request) {
	var payload types.ChangeEmailAddressPayload
	ctx := req.Context()
	uuid, ok := ctx.Value("userId").(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(noContext)
	}

	err := json.NewDecoder(req.Body).Decode(&payload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(err)
	}

	err = s.db.ChangeEmailAddress(uuid, payload.NewEmail)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(err)
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) deleteAccount(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	uuid, ok := ctx.Value("userId").(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(noContext)
	}

	err := s.db.DeleteAllTasks(uuid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(err)
	}

	err = s.db.DeleteAllApiKeys(uuid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(err)
	}

	err = s.db.DeleteUser(uuid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(err)
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) getProfilePhoto(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	uuid, ok := ctx.Value("userId").(string)

	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(noContext)
	}

	profilePhoto, err := s.db.GetProfilePhoto(uuid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(err)
	}
	j, err := json.Marshal(types.ProfilePhotoResponse{ProfilePhoto: profilePhoto})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func (s *Server) changeProfilePhoto(w http.ResponseWriter, req *http.Request) {
	var payload types.ProfilePhotoPayload
	ctx := req.Context()
	uuid, ok := ctx.Value("userId").(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(noContext)
	}

	err := json.NewDecoder(req.Body).Decode(&payload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(err)
	}

	err = s.db.ChangeProfilePhoto(uuid, payload.ProfilePhoto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(err)
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) resetUserPassword(w http.ResponseWriter, req *http.Request) {
	var passwordRequestToken types.ResetPasswordRequestPayload

	ctx := req.Context()
	uuid, ok := ctx.Value("userId").(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(noContext)
	}

	err := json.NewDecoder(req.Body).Decode(&passwordRequestToken)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(err)
	}

	userToken, err := s.db.GetResetPasswordToken(uuid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(err)
	}
	
	if !auth.DeletionRequestValid(passwordRequestToken.ResetToken, userToken) {
		w.WriteHeader(http.StatusUnauthorized)
		s.logger.Panic(err)
	} 
	err = s.db.ClearUserPassword(uuid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(err)
	}
	// Empties the token so they can't be reused if it is compromised
	err = s.db.SetForgotPasswordToken(uuid, "")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(err)
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) setNewUserPassword(w http.ResponseWriter, req *http.Request) {
	var resetPasswordPayload types.ResetPasswordPayload
	ctx := req.Context()
	uuid, ok := ctx.Value("userId").(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(noContext)
	}

	err := json.NewDecoder(req.Body).Decode(&resetPasswordPayload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(err)
	}
	newPassword, _ := auth.EncryptPassword(resetPasswordPayload.NewPassword)
	err = s.db.SetNewPassword(uuid, newPassword)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(err)
	}
}

func (s *Server) sendForgotPasswordRequest(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	uuid, ok := ctx.Value("userId").(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(noContext)
	}

	userToken, _ := auth.GetNewApiKey()
	err := s.db.SetForgotPasswordToken(uuid, userToken)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(err)
	}

	email, err := s.db.GetEmailAddress(uuid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(err)
	}
	err = contact.SendResetEmail(email, userToken)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Panic(err)
	}

	// user will be rerouted to a default forgot password page 
	// after filling out the information the user will have their password reset
	// we will need to add a check for an empty password though to redirect people without passwords to the forgot password page
}
