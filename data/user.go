package data

type User struct {
	Id        int    `json:"id"`
	FirstName string `json:"first-name" validate:"required"`
	LastName  string `json:"last-name" validate:"required"`
	Nickname  string `json:"nickname" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
	Country   string `json:"country" validate:"required"`
}

type UserDBI interface {
	AddUser(user User) error
	UpdateUser(user User) error
	GetUsers() (users []*User)
	DeleteUser(id int) error
}
