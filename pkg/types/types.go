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

type UserLoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	UserLoginPayload
	AccountType string `json:"accountType,omitempty"`
}

type ApiKey struct {
	ApiKey string `json:"key"`
}

type Email struct {
	Email string `json:"email"`
}
type JwtClaims struct {
	jwt.RegisteredClaims
	Uuid string `json:"uuid"`
	// Add more custom claims here
}

type JwtResponse struct {
	Jwt string `json:"jwt"`
}

type ErrResponse struct {
	Message string `json:"message"`
}
