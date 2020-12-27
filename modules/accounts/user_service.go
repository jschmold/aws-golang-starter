package accounts

import (
	"github.com/go-pg/pg/v10"
	"golang.org/x/crypto/bcrypt"
)

const hashCost int = 14

// IUserService is the interface that describes services performing actions on behalf of a UserModel
type IUserService interface {
	Create(CreateUserDTO) (*UserModel, error)
	Get(string) (*UserModel, error)
	GetByEmail(string) (*UserModel, error)
	SetPassword(userID string, pwd string) (*UserModel, error)
}

// userService exposes DB Functionality to this module
type userService struct {
	db *pg.DB
}

// NewUserService creates a new instance of the service responsible for User interactions
func NewUserService(db *pg.DB) IUserService {
	return &userService{db}
}

// Create a new user via their email with a password
func (svc userService) Create(reg CreateUserDTO) (user *UserModel, err error) {
	user = &UserModel{Email: reg.Email}

	if _, err = svc.db.Model(user).Returning("*").Insert(); err != nil {
		return
	}

	if len(reg.Password) > 0 {
		user, err = svc.SetPassword(user.ID, reg.Password)
	}

	return
}

// Get a user by their user id
func (svc *userService) Get(id string) (user *UserModel, err error) {
	user = &UserModel{}
	err = svc.db.Model(user).Where("id = ?", id).First()

	return
}

// GetByEmail finds a user by their email
func (svc *userService) GetByEmail(email string) (user *UserModel, err error) {
	user = &UserModel{}
	err = svc.db.Model(user).Where("email = ?", email).First()

	return
}

// SetPassword hashes a string password, and assigns it to the given userID's `password` field
func (svc *userService) SetPassword(userID string, password string) (user *UserModel, err error) {
	user, err = svc.Get(userID)
	if err != nil {
		return
	}

	pwd := []byte(password)
	user.Pwd, err = bcrypt.GenerateFromPassword(pwd, hashCost)

	svc.db.Model(user).WherePK().Update()

	return
}
