package handlers

import (
	"github.com/go-openapi/runtime/middleware"
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

type UserUpdateKey struct{}

type GenericError struct {
	Message string `json:"message"`
}

type ValidationError struct {
	Messages []string `json:"messages"`
}

func NewUsersHandler(l *logrus.Logger, dbi data.UserDBI, v *data.Validation) *Users {
	return &Users{l, dbi, v}
}

func InitHandlers(usersHandlers *Users) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	//var api1 = router.PathPrefix("/v1").Subrouter()

	getR := router.Methods(http.MethodGet).Subrouter()
	getR.HandleFunc("/users", usersHandlers.ListAll)

	postR := router.Methods(http.MethodPost).Subrouter()
	postR.HandleFunc("/users", usersHandlers.Add)
	postR.Use(usersHandlers.MiddlewareValidationUser)

	putR := router.Methods(http.MethodPut).Subrouter()
	putR.HandleFunc("/users", usersHandlers.Update)
	putR.Use(usersHandlers.MiddlewareValidationUpdateUser)

	deleteR := router.Methods(http.MethodDelete).Subrouter()
	deleteR.HandleFunc("/users/{id:[0-9]}", usersHandlers.Delete)

	// handler for documentation
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)

	getR.StrictSlash(true).Handle("/docs", sh)
	getR.StrictSlash(true).Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	return router
}
