package amazon

import (
	"encoding/base64"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	mgr "github.com/aws/aws-sdk-go/service/secretsmanager"
)

// GetAWSSecret gets a secret from AWS using the current stage
func GetAWSSecret(name, region string) (secret string, err error) {

	//Create a Secrets Manager client
	svc := mgr.New(session.New(), aws.NewConfig().WithRegion(region))

	input := &mgr.GetSecretValueInput{SecretId: aws.String(name)}
	fmt.Printf("SecretId: %s\n", *input.SecretId)

	// In this sample we only handle the specific exceptions for the 'GetSecretValue' API.
	// See https://docs.aws.amazon.com/mgr/latest/apireference/API_GetSecretValue.html

	result, err := svc.GetSecretValue(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case mgr.ErrCodeDecryptionFailure:
				// Secrets Manager can't decrypt the protected secret text using the provided KMS key.
				fmt.Println(mgr.ErrCodeDecryptionFailure, aerr.Error())

			case mgr.ErrCodeInternalServiceError:
				// An error occurred on the server side.
				fmt.Println(mgr.ErrCodeInternalServiceError, aerr.Error())

			case mgr.ErrCodeInvalidParameterException:
				// You provided an invalid value for a parameter.
				fmt.Println(mgr.ErrCodeInvalidParameterException, aerr.Error())

			case mgr.ErrCodeInvalidRequestException:
				// You provided a parameter value that is not valid for the current state of the resource.
				fmt.Println(mgr.ErrCodeInvalidRequestException, aerr.Error())

			case mgr.ErrCodeResourceNotFoundException:
				// We can't find the resource that you asked for.
				fmt.Println(mgr.ErrCodeResourceNotFoundException, aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	// Decrypts secret using the associated KMS CMK.
	// Depending on whether the secret is a string or binary, one of these fields will be populated.
	if result.SecretString != nil {
		secret = *result.SecretString
		return
	}

	var size int
	byteCount := base64.StdEncoding.DecodedLen(len(result.SecretBinary))
	decodedBinarySecretBytes := make([]byte, byteCount)
	size, err = base64.StdEncoding.Decode(decodedBinarySecretBytes, result.SecretBinary)

	if err != nil {
		fmt.Println("Base64 Decode Error:", err)
		return
	}
	secret = string(decodedBinarySecretBytes[:size])

	return
}
