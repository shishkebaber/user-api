package handlers

import (
	"context"
	"github.com/shishkebaber/user-api/data"
	"net/http"
)

func (u *Users) MiddlewareValidationUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")

		user := &data.User{}

		err := data.FromJson(user, r.Body)
		if err != nil {
			u.logger.Error("Error during deserialization User ", err)

			rw.WriteHeader(http.StatusBadRequest)
			data.ToJson(&GenericError{Message: err.Error()}, rw)
			return
		}

		//validate
		errs := u.v.Validate(user)
		if len(errs) != 0 {
			u.logger.Error("User validation error ", errs)

			rw.WriteHeader(http.StatusUnprocessableEntity)
			data.ToJson(&ValidationError{Messages: errs.Errors()}, rw)
			return
		}

		ctx := context.WithValue(r.Context(), UserKey{}, user)
		r = r.WithContext(ctx)

		//Proceed to the next handler
		next.ServeHTTP(rw, r)
	})
}

func (u *Users) MiddlewareValidationUpdateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")

		user := &data.UpdateUser{}

		err := data.FromJson(user, r.Body)
		if err != nil {
			u.logger.Error("Error during deserialization User for update ", err)

			rw.WriteHeader(http.StatusBadRequest)
			data.ToJson(&GenericError{Message: err.Error()}, rw)
			return
		}

		//validate
		errs := u.v.Validate(user)
		if len(errs) != 0 {
			u.logger.Error("User update validation error ", errs)

			rw.WriteHeader(http.StatusUnprocessableEntity)
			data.ToJson(&ValidationError{Messages: errs.Errors()}, rw)
			return
		}

		ctx := context.WithValue(r.Context(), UserUpdateKey{}, user)
		r = r.WithContext(ctx)

		//Proceed to the next handler
		next.ServeHTTP(rw, r)
	})
}
