package handlers

import (
	"github.com/shishkebaber/user-api/data"
	"net/http"
)

// swagger:route GET /users users listUsers
// Return a list of users from the database, possible filters by fields [first-name, last-name, nickname, email,  country]
// responses:
//	200: usersResponse
func (users *Users) ListAll(rw http.ResponseWriter, r *http.Request) {
	users.logger.Info("Getting users")
	rw.Header().Add("Content-Type", "application/json")
	args := r.URL.Query()

	result, err := users.Db.GetUsers(args)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJson(&GenericError{Message: err.Error()}, rw)
		return
	}

	err = data.ToJson(result, rw)
	if err != nil {
		// we should never be here but log the error just incase
		users.logger.Error("Unable to serializing user", "error", err)
	}
}
