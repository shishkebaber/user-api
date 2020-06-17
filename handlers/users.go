package handlers

import (
	"github.com/gorilla/mux"
	"github.com/shishkebaber/user-api/data"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Users struct {
	logger *logrus.Logger
	Db     data.UserDBI
	v      *data.Validation
}

type UserKey struct{}

type GenericError struct {
	Message string `json:"message"`
}

type ValidationError struct {
	Messages []string `json:"messages"`
}

func NewUsersHandler(l *logrus.Logger, dbi data.UserDBI, v *data.Validation) *Users {
	return &Users{l, dbi, v}
}

func (users *Users) ListAll(rw http.ResponseWriter, r *http.Request) {

}

func InitHandlers(usersHandlers *Users) *mux.Router {
	sMux := mux.NewRouter()
	getR := sMux.Methods(http.MethodGet).Subrouter()
	getR.HandleFunc("/users", usersHandlers.ListAll)
	return sMux
}
