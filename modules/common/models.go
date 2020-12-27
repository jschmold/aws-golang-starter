package common

import "time"

// Verification is the type that all verifications branch from
type Verification struct {
	CreatedAt  time.Time `pg:"created_at,type:timestamptz,notnull"`
	Code       string    `pg:"code,type:varchar(32),notnull"`
	VerifiedAt time.Time `pg:"verified_at,type:timestamptz"`
	ExpiresAt  time.Time `pg:"expires_at,type:timestamptz"`
}

// Timestamps is the type that is used in any table that has the `common.timestamps`
type Timestamps struct {
	CreatedAt time.Time `pg:"created_at,type:timestamptz"`
	UpdatedAt time.Time `pg:"updated_at,type:timestamptz"`
	DeletedAt time.Time `pg:"deleted_at,type:timestamptz"`
}
