package models

type User struct {
	ID             int     `json:"id"`
	Username       string  `json:"usernmame"`
	Password       string  `json:"password"`
	Email          string  `json:"email"`
	IsAdmin        bool    `json:"is_admin"`
	ProfilePicture *string `json:"profile_picture"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}
