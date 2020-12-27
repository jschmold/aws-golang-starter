// Use this file to simplify the writing of error messages to the response object.
// The example "EmailInUseException" demonstrates how to go about creating these exception shortcuts.
package accounts

import (
	"github.com/jschmold/aws-golang-starter/modules/common"
	"github.com/jschmold/aws-golang-starter/modules/http"
)

// EmailInUseException Indicate that the email someone is registering with is already used
func EmailInUseException(res http.IResponse) {
	err := common.ErrorResponse{
		Message: "Sorry, that email has already been registered. Try logging in?",
		Error:   "Email in use",
	}
	res.SetBody(err)
}
