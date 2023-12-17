package types 


type Task struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	Body     string `json:"body"`
	Due      string `json:"due"`
	Priority string `json:"priority"`
}
