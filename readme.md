## Purpose

I wanted to be able to use AWS in a fashion that was simple, straightfoward, and easy to test. Attempting
to use Serverless was met with a lot of challenges, but what you see here in this repository solved the
following issues for me:

 - Easy local unit testing
 - Easy local integration testing
 - CI friendly testing pipeline
 - Local / Environment config so we can connect to Amazon RDS or a local Postgres
 - Easy migrations in SQL

## Shortcomings

This isn't a perfect solution, in that the setup process is not entirely automatic. You have to setup
RDS first so it is publicly accessible to GitHub (GitHub Actions), and create the database user yourself.
You must also setup the AWS Secrets Manager provide some DB info basics to this application by hand.

However, once the secrets manager and RDS are present, this is very easy to use.

## Setup

1. Setup your environment
  - Create an AWS account
  - Install the Serverless sdk
  - Setup your golang environment (I use gvm)
  - Sign into your AWS with Serverless
  - Install the AWS cli
2. Create a new RDS database (postgres, free tier is what I use)
  - During creation, make sure you set a master password for the `postgres` user.
3. Setup the Secrets Manager to draw from the created database. The name is _important_.
  - `DB_SECRET: ${opt:stage, 'dev'}-database-credentials` is the environment variable in the `services/accounts/serverless.yml`
  file. This means that your secret name _should match_ this pattern. I use `dev-database-credentials` for the secret name for
  this reason.
4. Create a GitHub repo, and set up the secrets. Reference the config section to see what you need.
  - Setup the AWS stuff in the "secrets" using the Variable and a valid value 
  - Note: The GitHub Actions do _NOT_ need access to a database for integration testing. The GitHub Actions makes their _own_ Postgres server and uses _that_ to integration test. You only need to provide the AWS stuff for the migrations portion of the actions.
5. Setup your local postgres. Your user should have permission to create extensions in the database.


## Running Tests

#### Unit testing

1. Generate the mocks automagically using `make mocks`. This automatically generates _all_ mocks for code within the `modules` folder, 
and cleans up the empty files that causes builds to fail.
2. Run `make test` for unit testing (this runs the <filename>_test.go files)

#### Integration testing

1. Create a `config.toml` file and add a `PG_URL="postgres://user:password@localhost:5432/example`
2. Run `make migrator` to create the automatic migrator
  - Note: this is config-aware, so if you have a `config.toml` file, it will read that
3. Run `./bin/migrator init` to setup the migration app
4. Run `./bin/migrator up` to run the migrations up to the latest version
5. Run `make e2e` to run the integration tests

## Config

See `config.go` in the root folder to see what is loaded application-wide.

| Variable | Value |
| :------- | :---- |
| AWS_REGION | This should be the AWS_REGION that you use with RDS. Required unless using PG_URL |
| DB_SECRET | The name of the secret to pull the database connection info from. Required unless using PG_URL |
| AWS_ACCESS_KEY_ID | This is the access key for the IAM user your app uses. Required unless using PG_URL |
| AWS_SECRET_ACCESS_KEY | This is the "password" for your IAM user for this app. Required unless PG_URL is in use. |
| PG_URL | The full DB address (including password) to connect to. I use this for local integration testing. |

## Commands

`make` will build the services in the `services` folder

`make e2e` will run the integration tests

`make test` will run all the unit tests in the modules folder

`make deploy` will deploy the services using your local install of serverless

`make migrator`