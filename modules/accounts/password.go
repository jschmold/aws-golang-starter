package accounts

import (
	"github.com/go-playground/validator/v10"
	"github.com/jschmold/aws-golang-starter/modules/common"
)

// PasswordResetModel describes exactly that: A password reset.
type PasswordResetModel struct {
	common.Verification

	tableName struct{} `pg:"accounts.password_resets"`
	ID        string   `pg:"id,type:uuid"`
	UserID    string   `pg:"user_id,type:uuid"`
	Email     string   `pg:"email,type:citext"`
}

// PasswordValidate ensures that passwords are at least 8 chars long, fewer than 64 chars long
// and is a custom validator for the `go-playground/validator/v10 package`
func PasswordValidate(field validator.FieldLevel) bool {
	val := field.Field().String()
	return len(val) > 7 && len(val) < 65
}
