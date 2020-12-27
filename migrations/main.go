package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/go-pg/migrations/v8"
	"github.com/go-pg/pg/v10"
	server "github.com/jschmold/aws-golang-starter"
	"github.com/jschmold/aws-golang-starter/modules/amazon"
)

func connect() (db *pg.DB, err error) {

	// first try to connect locally, because if PG_URL is set
	// this means that it's not a secret, so we can just use it
	connString := os.Getenv(server.PgURLEnv)
	if connString != "" {
		fmt.Printf("Using %s to connect to postgres\n", server.PgURLEnv)
		return localConnection(connString)
	}

	region := os.Getenv(server.AWSRegionEnv)
	if region != "" {
		fmt.Println("Using amazon to connect to postgres")
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

func main() {
	err := server.LoadConfig()
	if err != nil {
		panic(err)
	}

	flag.Usage = usage
	flag.Parse()

	db, err := connect()

	if err != nil {
		panic(err)
	}

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	migrations.DefaultCollection.DiscoverSQLMigrations(dir + "/migrations")

	oldVersion, newVersion, err := migrations.Run(db, flag.Args()...)
	if err != nil {
		exitf(err.Error())
	}

	if newVersion != oldVersion {
		fmt.Printf("migrated from version %d to %d\n", oldVersion, newVersion)
	} else {
		fmt.Printf("version is already %d\n", oldVersion)
	}
}
