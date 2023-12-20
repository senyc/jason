package types

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
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	AccountType string `json:"accountType,omitempty"`
}

type ApiKey struct {
	ApiKey string `json:"key"`
}

type Email struct {
	Email string `json:"email"`
}
