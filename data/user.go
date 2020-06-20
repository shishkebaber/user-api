package data

import "golang.org/x/crypto/bcrypt"

type User struct {
	Id        int    `json:"-"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Nickname  string `json:"nickname" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
	Country   string `json:"country" validate:"required"`
}

type UpdateUser struct {
	Id        int    `json:"id" validate:"required"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Nickname  string `json:"nickname" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Country   string `json:"country" validate:"required"`
}

type UserDBI interface {
	AddUser(user User) error
	UpdateUser(user UpdateUser) error
	GetUsers(filters map[string][]string) ([]*UpdateUser, error)
	DeleteUser(id int) (int64, error)
}

func hashSaltPassword(pwd string) (*string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pwd), 8)
	if err != nil {
		return nil, err
	}
	result := string(hashedPassword)
	return &result, nil
}
