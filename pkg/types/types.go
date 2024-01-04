package types

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Task struct {
	Id        int       `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Due       time.Time `json:"due,omitempty"`
	Priority  int16     `json:"priority"`
	Completed bool      `json:"completed"`
}

type CompletedTask struct {
	Id            int       `json:"id"`
	Title         string    `json:"title"`
	Body          string    `json:"body"`
	Due           time.Time `json:"due,omitempty"`
	Priority      int16     `json:"priority"`
	CompletedDate string    `json:"completedDate"`
}

type IncompleteTask struct {
	Id       int       `json:"id"`
	Title    string    `json:"title"`
	Body     string    `json:"body"`
	Due      time.Time `json:"due,omitempty"`
	Priority int16     `json:"priority"`
}

type NewTask struct {
	Title    string    `json:"title"`
	Body     string    `json:"body,omitempty"`
	Due      time.Time `json:"due,omitempty"`
	Priority int16     `json:"priority,omitempty"`
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	UserLogin
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
