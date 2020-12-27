package test

import (
	"errors"
	"os"

	"github.com/go-pg/pg/v10"
	server "github.com/jschmold/aws-golang-starter"
	"github.com/jschmold/aws-golang-starter/modules/amazon"
)

// NewDB initializes a new Postgres connection using the PG_URL environment variable
func NewDB() (db *pg.DB, err error) {
	// first try to connect locally, because if PG_URL is set
	// this means that it's not a secret, so we can just use it
	connString := os.Getenv(server.PgURLEnv)
	if connString != "" {
		return localConnection(connString)
	}

	region := os.Getenv(server.AWSRegionEnv)
	if region != "" {
		return amazon.GetDBConnection()
	}

	return nil, errors.New("Neither PG_URL or amazon keys were set, unable to connect")
}

func localConnection(url string) (db *pg.DB, err error) {
	opts, err := pg.ParseURL(url)
	if err != nil {
		return nil, err
	}

	return pg.Connect(opts), nil
}
