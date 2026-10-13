package todo

type Todo struct {
	ID        int    `json:"id"`
	Task      string `json:"task"`
	Status    string `json:"status"`
	Priority  string `json:"priority"`
	UpdatedAt string `json:"updated_at"`
	CreatedAt string `json:"created_at"`
}
