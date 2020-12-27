package accounts

import "github.com/jschmold/aws-golang-starter/modules/common"

// CreateUserDTO is the DTO for creating a new user via registration with Email/Password
type CreateUserDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
}

// UserModel describes the table `accounts.users`
type UserModel struct {
	common.Timestamps

	tableName struct{} `pg:"accounts.users"`
	ID        string   `pg:"id,type:uuid"`
	Name      string   `pg:"name,type:varchar(256)"`
	Email     string   `pg:"email,type:citext"`
	Verified  bool     `pg:"email_verified,type:boolean"`
	Pwd       []byte   `pg:"password,type:bytea"`
}

// UserConfirmationModel is the struct that represents the `accounts.user_confirmations` table.
// A UserConfirmation is created when someone registers or is invited.
type UserConfirmationModel struct {
	common.Verification

	tableName struct{} `pg:"accounts.user_confirmations"`
	ID        string   `pg:"id,type:uuid"`
	UserID    string   `pg:"user_id,type:uuid"`
}
