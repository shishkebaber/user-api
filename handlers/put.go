package handlers

import (
	"github.com/shishkebaber/user-api/data"
	"net/http"
)

// swagger:route PUT /users users updateUser
// Update user details
// responses:
//	201: noContentResponse
//  404: errorResponse
//  422: errorValidation
func (users *Users) Update(rw http.ResponseWriter, r *http.Request) {
	users.logger.Info("Updating User")
	rw.Header().Add("Content-Type", "application/json")

	input := r.Context().Value(UserUpdateKey{}).(*data.UpdateUser)

	users.logger.Info("Adding User")
	err := users.Db.UpdateUser(*input)
	if err != nil {
		users.logger.Error("User not found", err)

		rw.WriteHeader(http.StatusNotFound)
		data.ToJson(&GenericError{Message: "User not found in database"}, rw)
		return
	}
}
