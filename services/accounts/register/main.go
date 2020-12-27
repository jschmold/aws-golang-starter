package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jschmold/aws-golang-starter/modules/accounts"
	"github.com/jschmold/aws-golang-starter/modules/amazon"
	"github.com/jschmold/aws-golang-starter/modules/http"
)

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(arg events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	session := amazon.NewContext(&arg)
	conn, err := amazon.GetDBConnection()
	if err != nil {
		http.InternalServerException(session)
		return session.Response.Respond(), nil
	}

	defer conn.Close()

	deps := accounts.CreateRegistrationControllerDeps(conn)
	register := accounts.NewRegistrationController(deps)

	register.WithEmail(session)
	return session.Response.Respond(), nil
}

func main() {
	lambda.Start(Handler)
}
