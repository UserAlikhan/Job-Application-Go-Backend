package models

type Job struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Location    string `json:"location"`
	Company     string `json:"company"`
	Salary      string `json:"salary"`
	UserID      int    `json:"user_id"`
	User        *User  `json:"user,omitempty"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
