package types

import (
	"database/sql"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type NullTime struct {
	time.Time
}

func (t NullTime) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("null"), nil
	} else {
		return t.Time.MarshalJSON()
	}
}

type SqlTasksRow struct {
	Id            int
	Title         string
	Body          sql.NullString
	Due           sql.NullTime
	TimeCreated   time.Time
	Priority      int16
	Completed     bool
	CompletedDate sql.NullTime
}

type TaskReponse struct {
	Id            int        `json:"id"`
	Title         string     `json:"title"`
	Body          string     `json:"body"`
	Due           NullTime   `json:"due"`
	Priority      int16      `json:"priority"`
	Completed     bool       `json:"completed"`
	CompletedDate *time.Time `json:"completedDate,omitempty"`
}

type CompletedTaskResponse struct {
	Id            int       `json:"id"`
	Title         string    `json:"title"`
	Body          string    `json:"body"`
	Due           NullTime  `json:"due"`
	Priority      int16     `json:"priority"`
	CompletedDate time.Time `json:"completedDate,omitempty"`
}

type IncompleteTaskResponse struct {
	Id       int      `json:"id"`
	Title    string   `json:"title"`
	Body     string   `json:"body"`
	Due      NullTime `json:"due"`
	Priority int16    `json:"priority"`
}

type NewTaskPayload struct {
	Title    string    `json:"title"`
	Body     string    `json:"body"`
	Due      time.Time `json:"due"`
	Priority int16     `json:"priority,omitempty"`
}

type EditTaskPayload struct {
	Id       int        `json:"id"`
	Title    string     `json:"title,omitempty"`
	Body     string     `json:"body,omitempty"`
	Due      *time.Time `json:"due,omitempty"`
	Priority int16      `json:"priority,omitempty"`
}

type ApiKeyPayload struct {
	Label       string     `json:"label"`
	Description string     `json:"description,omitempty"`
	Expiration  *time.Time `json:"expiration,omitempty"`
}

type UserLoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type EmailResponse struct {
	Email string `json:"email"`
}

type User struct {
	UserLoginPayload
	AccountType string `json:"accountType,omitempty"`
}

type ChangeEmailAddressPayload struct {
	NewEmail string `json:"newEmail"`
}

type ApiKeyResponse struct {
	ApiKeyId string `json:"id"`
	ApiKey   string `json:"apikey"`
}

type ApiKeyMetadata struct {
	Id           string     `json:"id"`
	Label        string     `json:"label"`
	Description  string     `json:"description"`
	Expiration   time.Time  `json:"expiration"`
	LastAccessed *time.Time `json:"lastAccessed"`
	CreationDate time.Time  `json:"creationDate"`
}

type Email struct {
	Email string `json:"email"`
}

type JwtClaims struct {
	jwt.RegisteredClaims
	Uuid string `json:"uuid"`
}

type JwtResponse struct {
	Jwt string `json:"jwt"`
}

type SyncTimeResponse struct {
	SyncTime time.Time `json:"syncTime"`
}

type AccountCreationDateResponse struct {
	AccountCreationDate time.Time `json:"accountCreationDate"`
}

type ErrResponse struct {
	Message string `json:"message"`
}

type ProfilePhotoResponse struct {
	ProfilePhoto string `json:"profilePhoto"`
}

type ProfilePhotoPayload struct {
	ProfilePhoto int `json:"profilePhoto"`
}
