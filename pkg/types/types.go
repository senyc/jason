package types

import "github.com/golang-jwt/jwt/v5"

type Task struct {
	Id        string `json:"id"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	Due       string `json:"due"`
	Priority  string `json:"priority"`
	Completed bool   `json:"completed"`
}

type NewTask struct {
	Title    string `json:"title"`
	Body     string `json:"body,omitempty"`
	Priority string `json:"priority,omitempty"`
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
