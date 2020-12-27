package accounts

import (
	"github.com/go-pg/pg/v10"
	"github.com/go-playground/validator/v10"
	"github.com/jschmold/aws-golang-starter/modules/http"
)

// RegistrationControllerDeps describes the dependencies for the
// Registration controller, used in `NewRegistration`
type RegistrationControllerDeps struct {
	UserService IUserService
}

// RegistrationController governs registering and confirming new accounts
type RegistrationController struct {
	user    IUserService
	checker *validator.Validate
}

// CreateRegistrationControllerDeps takes any external dependencies, and generates direct dependencies
// that can be used to build a registration controller
func CreateRegistrationControllerDeps(db *pg.DB) RegistrationControllerDeps {
	user := NewUserService(db)
	return RegistrationControllerDeps{UserService: user}
}

// NewRegistrationController returns a new instance of a Registration controller for the endpoints
// to use, and serve to clients
func NewRegistrationController(deps RegistrationControllerDeps) (ctrl *RegistrationController) {
	validate := validator.New()
	validate.RegisterValidation("password", PasswordValidate)
	ctrl = &RegistrationController{user: deps.UserService, checker: validate}

	return
}

// WithEmail registers a new account with an email / password combo
func (regis *RegistrationController) WithEmail(ctx http.IContext) {

	req, res := ctx.GetRequest(), ctx.GetResponse()
	body := &CreateUserDTO{}
	req.GetBodyAs(body)

	err := regis.checker.Struct(body)

	switch err.(type) {

	case validator.ValidationErrors:
		http.FailedValidationException(ctx, err.(validator.ValidationErrors))
		return
	case nil:
		break

	default:
		http.InternalServerException(ctx)
		return
	}

	existing, err := regis.user.GetByEmail(body.Email)

	if err != nil {
		http.InternalServerException(ctx)
		return
	}

	if existing != nil {
		msg := "Someone has already registered with this email address"
		errorMsg := "Email in use"
		res.SetError(409, msg, errorMsg)
		return
	}

	_, err = regis.user.Create(*body)
	if err != nil {
		http.InternalServerException(ctx)
		return
	}

	return
}
