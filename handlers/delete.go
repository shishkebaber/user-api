package handlers

import (
	"github.com/gorilla/mux"
	"github.com/shishkebaber/user-api/data"
	"net/http"
	"strconv"
)

// swagger:route DELETE /users/{id} users deleteUser
// Delete user from DB
// responses:
//	200: noContentResponse
//  404: errorResponse
//  500: errorResponse
func (users *Users) Delete(rw http.ResponseWriter, r *http.Request) {
	users.logger.Info("Deleting User")
	rw.Header().Add("Content-Type", "application/json")

	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		users.logger.Error("Internal error", err)
	}
	users.logger.Info("Deleting User")
	deletedCount, err := users.Db.DeleteUser(id)
	if err != nil {
		users.logger.Error("User not deleted", err)

		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJson(&GenericError{Message: "User not deleted"}, rw)
		return
	}
	if deletedCount == 0 {
		users.logger.Error("User not found")

		rw.WriteHeader(http.StatusNotFound)
		data.ToJson(&GenericError{Message: "User not found"}, rw)
		return
	}
	rw.WriteHeader(http.StatusOK)
}
