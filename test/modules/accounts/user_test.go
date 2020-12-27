package user

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/go-pg/pg/v10"
	server "github.com/jschmold/aws-golang-starter"
	"github.com/jschmold/aws-golang-starter/modules/accounts"
	"github.com/jschmold/aws-golang-starter/test"
	"golang.org/x/crypto/bcrypt"
)

var database *pg.DB

func init() {
	server.LoadConfig()
}

func cleanTest(email string) {
	search := &accounts.UserModel{Email: email}
	database.Model(search).Where("email = ?", email).Delete()
}

func makeUser(email string) (mdl *accounts.UserModel, err error) {
	mdl = &accounts.UserModel{Email: email}
	_, err = database.Model(mdl).Insert()
	return
}

func TestMain(m *testing.M) {
	err := server.LoadConfig()

	if err != nil {
		errMsg := fmt.Sprintf("Unable to load config because: %s", err.Error())
		panic(errMsg)
	}

	fmt.Println("Tests running")

	database, err = test.NewDB()
	if err != nil {
		errMsg := fmt.Sprintf("Could not establish postgres connection: %s", err.Error())
		_ = fmt.Errorf("%s", errMsg)
		panic(errMsg)
	}
	defer database.Close()

	code := m.Run()
	os.Exit(code)
}

func TestUserCreate(suite *testing.T) {
	svc := accounts.NewUserService(database)

	email := "user_create@e2e-tests.example.com"

	dto := accounts.CreateUserDTO{Email: email, Password: "SomeThing"}

	suite.Run("Sets up the right details", func(t *testing.T) {
		defer cleanTest(email)
		user, err := svc.Create(dto)
		if err != nil {
			t.Error(err)
		}

		if user.Email != email {
			t.Errorf("Email does not match created, expected %s, got %s", email, user.Email)
		}

		if user.CreatedAt.IsZero() || user.UpdatedAt.IsZero() {
			t.Error("Timestamps were not created properly")
		}

		if user.Verified {
			t.Errorf("Created users should not start verified")
		}
	})

	suite.Run("hashes the password", func(t *testing.T) {
		defer cleanTest(dto.Email)
		user, err := svc.Create(dto)
		if err != nil {
			t.Error(err)
		}

		if string(user.Pwd) == dto.Password || len(user.Pwd) == 0 {
			t.Errorf("Service did not hash the provided password")
		}

		err = bcrypt.CompareHashAndPassword(user.Pwd, []byte(dto.Password))
		if err != nil {
			t.Errorf("Password was not hashed with Bcrypt")
		}
	})
}

func TestUserGet(t *testing.T) {
	svc := accounts.NewUserService(database)
	email := "user_get@e2e-tests.example.com"

	defer cleanTest(email)
	ref, err := makeUser(email)

	if err != nil {
		t.Error(err)
	}

	user, err := svc.Get(ref.ID)
	if err != nil {
		t.Error(err)
	}

	if reflect.DeepEqual(user, ref) != true {
		t.Errorf("User and ref are not deeply equal")
	}
}

func TestUserGetByEmail(suite *testing.T) {
	svc := accounts.NewUserService(database)
	email := "user_get_by_email@e2e-tests.example.com"

	suite.Run("Gets existing", func(t *testing.T) {
		defer cleanTest(email)
		ref, err := makeUser(email)

		if err != nil {
			suite.Error(err)
		}

		user, err := svc.GetByEmail(email)
		if err != nil {
			suite.Error(err)
		}

		if reflect.DeepEqual(user, ref) != true {
			suite.Errorf("User and ref are not deeply equal")
		}
	})

	suite.Run("Returns nil if no results", func(t *testing.T) {
		_, err := svc.GetByEmail(email)
		if err == nil {
			suite.Errorf("Expected an error, did not get one")
		}
	})
}
