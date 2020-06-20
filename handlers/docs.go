// Package classification Users API.
//
// Documentation for User API
//
// Schemes:http
// BasePath:/
// Version:1.0.0
//
// Consumes:
//	-application/json
// Produces:
//	-application/json
//swagger:meta
package handlers

import "github.com/shishkebaber/user-api/data"

// List of Users returns in the response
//swagger:response usersResponse
type usersResponseWrapper struct {
	//All users
	//In:body
	Body []data.UpdateUser
}

// swagger:parameters addUser
type userCreateParamsWrapper struct {
	// User data structure to Create.
	// Note: the id field is ignored by  create operations
	// in: body
	// required: true
	Body data.User
}

// swagger:parameters updateUser
type userUpdateParamsWrapper struct {
	// User data structure to Update.
	// Note: the password field is ignored by  update operations
	// in: body
	// required: true
	Body data.UpdateUser
}

// No content is returned by this API endpoint
// swagger:response noContentResponse
type noContentResponseWrapper struct {
}

// Generic error message returned as a string
// swagger:response errorResponse
type errorResponseWrapper struct {
	// Description of the error
	// in: body
	Body GenericError
}

// Validation errors defined as an array of strings
// swagger:response errorValidation
type errorValidationWrapper struct {
	// Collection of the errors
	// in: body
	Body ValidationError
}
