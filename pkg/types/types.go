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

type NewUser struct {
	FirstName   string `json:"firstName,omitempty"`
	LastName    string `json:"lastName,omitempty"`
	Email       string `json:"email"`
	AccountType string `json:"accountType,omitempty"`
}
