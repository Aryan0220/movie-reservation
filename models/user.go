package models

type User struct {
	ID	int			`json:"id"`
	Name	string	`json:"name"`
	Email 	string	`json:"email"`
	Role	bool	`json:"admin"`
	Password string `json:"password,omitempty"`
}

