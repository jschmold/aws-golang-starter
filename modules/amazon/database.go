package amazon

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/rds/rdsutils"
	"github.com/go-pg/pg/v10"
)

// DatabaseCredentials matches the type found in the secrets manager for Postgres secrets.
// This exposes serialization and AWS -> pg.Options functionality
type DatabaseCredentials struct {
	UserName string `json:"username"`
	Engine   string `json:"engine"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"dbname"`
	ID       string `json:"dbInstanceIdentifier"`
	Region   string
}

// GetDBConnection first called GetDatabaseCredentials and returns the error if there is one
// If it is successful, it calls connect, and returns the connection
func GetDBConnection() (db *pg.DB, err error) {
	creds, err := GetDatabaseCredentials()
	if err != nil {
		return
	}

	return creds.Connect()
}

// GetDatabaseCredentials automatically populates a new DatabaseCredentials struct based on the
// environment variables AWS_REGION and DB_SECRET. DB_SECRET is populated via the `serverless.yml`
// file, and the AWS_REGION is automatic via AWS Lambda. Combining the two, we get a
// DatabaseCredentials object, which we can then call `Connect` on to connect to postgres
func GetDatabaseCredentials() (creds DatabaseCredentials, err error) {
	region := os.Getenv("AWS_REGION")
	secretName := os.Getenv("DB_SECRET")

	creds = DatabaseCredentials{Region: region}
	secret, err := GetAWSSecret(secretName, region)

	if err != nil {
		fmt.Printf("Unable to get database credentials because:\n%s", err.Error())
		return
	}

	err = json.Unmarshal([]byte(secret), &creds)
	if err != nil {
		fmt.Printf("An error occurred unmarshalling DB credentials: \n%s", err.Error())
		return
	}

	return
}

// Endpoint returns the host:port for this credentials object
func (cr *DatabaseCredentials) Endpoint() string {
	return fmt.Sprintf("%s:%d", cr.Host, cr.Port)
}

// Pwd generates the AWS password using IAM
func (cr *DatabaseCredentials) Pwd() (string, error) {
	awsCreds := credentials.NewEnvCredentials()
	authToken, err := rdsutils.BuildAuthToken(cr.Endpoint(), cr.Region, cr.UserName, awsCreds)
	if err != nil {
		return "", err
	}

	return authToken, nil
}

// Connect takes these credentials and creates a Postgres connection out of them
func (cr *DatabaseCredentials) Connect() (db *pg.DB, err error) {
	pwd, err := cr.Pwd()
	if err != nil {
		return nil, err
	}

	opts := pg.Options{
		User:            cr.UserName,
		Addr:            cr.Endpoint(),
		Database:        cr.Database,
		Password:        pwd,
		ApplicationName: "Example",
		OnConnect:       onSuccessfulConnect,
		TLSConfig:       &tls.Config{InsecureSkipVerify: true},
	}

	db = pg.Connect(&opts)

	return
}

func onSuccessfulConnect(ctx context.Context, conn *pg.Conn) error {
	fmt.Println("Database: Successfully connected")
	return nil
}
