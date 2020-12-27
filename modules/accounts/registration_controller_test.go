package accounts_test

import (
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	svmocks "github.com/jschmold/aws-golang-starter/mocks/accounts"
	httpmocks "github.com/jschmold/aws-golang-starter/mocks/http"
	"github.com/jschmold/aws-golang-starter/modules/accounts"
	"github.com/jschmold/aws-golang-starter/modules/common"
	"github.com/jschmold/aws-golang-starter/modules/http"
)

type RegistrationTestSetup struct {
	controller *accounts.RegistrationController
	mock       *gomock.Controller
	user       *svmocks.MockIUserService
	context    *httpmocks.MockIContext
	request    *httpmocks.MockIRequest
	response   *httpmocks.MockIResponse
}

func SetupRegistrationController(t *testing.T) RegistrationTestSetup {
	mockCtrl := gomock.NewController(t)
	userService := svmocks.NewMockIUserService(mockCtrl)
	deps := accounts.RegistrationControllerDeps{userService}

	ctrl := accounts.NewRegistrationController(deps)
	req := httpmocks.NewMockIRequest(mockCtrl)
	res := httpmocks.NewMockIResponse(mockCtrl)
	cont := httpmocks.NewMockIContext(mockCtrl)

	cont.EXPECT().GetRequest().AnyTimes().Return(req)
	cont.EXPECT().GetResponse().AnyTimes().Return(res)

	return RegistrationTestSetup{
		ctrl,
		mockCtrl,
		userService,
		cont,
		req,
		res,
	}
}

func fakeUserModel(email string) *accounts.UserModel {
	stamps := common.Timestamps{
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	}

	return &accounts.UserModel{
		ID:         "some random id",
		Email:      email,
		Pwd:        []byte("askdflklkasdfkllasdf"),
		Name:       "",
		Verified:   false,
		Timestamps: stamps,
	}
}

func TestRegister(suite *testing.T) {

	arg := accounts.CreateUserDTO{
		Email:    "me@jonathanschmold.ca",
		Password: "SomePassword",
	}

	suite.Run("Calls user service properly", func(t *testing.T) {
		setup := SetupRegistrationController(t)
		var register, mock, userService = setup.controller, setup.mock, setup.user

		defer mock.Finish()
		defer register.WithEmail(setup.context)

		mod := fakeUserModel(arg.Email)

		setup.request.
			EXPECT().
			GetBodyAs(gomock.Any()).
			Times(1).
			Do(func(body *accounts.CreateUserDTO) {
				body.Email = arg.Email
				body.Password = arg.Password
			})

		userService.EXPECT().
			GetByEmail(gomock.Any()).
			Return(nil, nil)

		userService.EXPECT().
			Create(gomock.Eq(arg)).
			Times(1).
			Return(mod, nil)

	})

	suite.Run("checks if email already in use", func(t *testing.T) {
		setup := SetupRegistrationController(t)
		mod := fakeUserModel(arg.Email)
		msg := "Someone has already registered with this email address"
		err := "Email in use"

		var register, mock, userService = setup.controller, setup.mock, setup.user
		defer mock.Finish()

		userService.EXPECT().
			GetByEmail(gomock.Eq(arg.Email)).
			Times(1).
			Return(mod, nil)

		userService.EXPECT().
			Create(gomock.Any()).
			Times(0)

		setup.request.EXPECT().
			GetBodyAs(gomock.Eq(&accounts.CreateUserDTO{})).
			Times(1).
			Do(func(body *accounts.CreateUserDTO) {
				body.Email = arg.Email
				body.Password = arg.Password
			})

		setup.response.
			EXPECT().
			SetError(409, msg, err)

		register.WithEmail(setup.context)
	})

	suite.Run("returns error with invalid body calls to controller", func(t *testing.T) {
		setup := SetupRegistrationController(t)
		val := validator.New()
		val.RegisterValidation("password", accounts.PasswordValidate)

		errors := val.Struct(&accounts.CreateUserDTO{}).(validator.ValidationErrors)
		err := common.ValidationErrorDict(errors)

		var register, mock = setup.controller, setup.mock

		defer mock.Finish()
		defer register.WithEmail(setup.context)

		response, request := setup.response, setup.request

		body := common.ValidationErrorResponse{
			Message:    http.FailedValidationExceptionMessage,
			Error:      http.FailedValidationExceptionError,
			Validation: err,
		}

		response.EXPECT().SetBody(body)
		response.EXPECT().SetCode(400)
		request.EXPECT().GetBodyAs(gomock.Any()).Do(func(arg interface{}) {})
	})
}
