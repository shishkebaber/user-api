package handlers

import (
	"github.com/shishkebaber/user-api/data"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Users struct {
	logger *logrus.Logger
	db     data.UserDBI
	v      *data.Validation
}

type UserKey struct{}

type GenericError struct {
	Message string `json:"message"`
}

type ValidationError struct {
	Messages []string `json:"messages"`
}

func NewUsers(l *logrus.Logger, dbi data.UserDBI, v *data.Validation) *Users {
	return &Users{l, dbi, v}
}

func (users *Users) ListAll(rw http.ResponseWriter, r *http.Request) {

}
