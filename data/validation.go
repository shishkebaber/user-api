package data

import (
	"fmt"
	"github.com/go-playground/validator"
)

type ValidationError struct {
	validator.FieldError
}

type ValidationErrors []ValidationError

func (v ValidationError) Error() string {
	return fmt.Sprintf(
		"Error: '%s' field validation failed on the '%s' tag",
		v.Field(),
		v.Tag(),
	)
}

func (v ValidationErrors) Errors() []string {
	errs := []string{}
	for _, err := range v {
		errs = append(errs, err.Error())
	}

	return errs
}

type Validation struct {
	validate *validator.Validate
}

func NewValidation() *Validation {
	validate := validator.New()
	return &Validation{validate}
}

func (v *Validation) Validate(i interface{}) ValidationErrors {
	errs := v.validate.Struct(i)
	var resultErrs ValidationErrors
	if errs != nil {
		errs := errs.(validator.ValidationErrors)
		if len(errs) == 0 {
			return nil
		}

		for _, err := range errs {
			e := ValidationError{err.(validator.FieldError)}
			resultErrs = append(resultErrs, e)
		}
	}

	return resultErrs
}
