package handlers

import (
	"github.com/shishkebaber/user-api/data"
	"net/http"
)

// swagger:route POST /users users addUser
// Add user to DB
// responses:
//	201: noContentResponse
//  422: errorValidation
//  500: errorResponse
func (users *Users) Add(rw http.ResponseWriter, r *http.Request) {
	input := r.Context().Value(UserKey{}).(*data.User)

	users.logger.Info("Adding User")
	err := users.Db.AddUser(*input)
	if err != nil {
		users.logger.Error("User not created", err)

		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJson(&GenericError{Message: "User not created"}, rw)
		return
	}
	rw.WriteHeader(http.StatusCreated)
}
