package http

import (
	"github.com/go-playground/validator/v10"
	"github.com/jschmold/aws-golang-starter/modules/common"
)

// InternalServerException takes an IContext pointer, assigns the "something bad happened" error,
// and then sets the error code for you.
func InternalServerException(cont IContext) {
	msg := "Uh oh! Something went wrong with our server. Contact support if it persists."
	error := "An internal error occurred"
	cont.GetResponse().SetError(500, msg, error)
}

// FailedValidationException receives a set of validation errors and transforms it into something
// that is a little easier to parse for clients
func FailedValidationException(cont IContext, errors validator.ValidationErrors) {
	msg := FailedValidationExceptionMessage
	err := FailedValidationExceptionError
	valid := common.ValidationErrorDict(errors)

	respBody := common.ValidationErrorResponse{
		Error:      err,
		Message:    msg,
		Validation: valid,
	}

	resp := cont.GetResponse()
	resp.SetBody(respBody)
	resp.SetCode(400)
}
