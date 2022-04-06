package common

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

func GetAWSSecret(secretName string) map[string]interface{} {
	// Use this code snippet in your app.
	// If you need more information about configurations or implementing the sample code, visit the email docs:
	// https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/setting-up.html

	//Create a Secrets Manager client
	sess, err := session.NewSession()
	if err != nil {
		// Handle session creation error
		fmt.Println(err.Error())
		CheckIfErrorWithMessage(err, "Could not get session")
	}
	svc := secretsmanager.New(sess,
		aws.NewConfig().WithRegion(AWSRegion))

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	// In this sample we only handle the specific exceptions for the 'GetSecretValue' API.
	// See https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html

	result, err := svc.GetSecretValue(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeDecryptionFailure:
				// Secrets Manager can't decrypt the protected secret text using the provided KMS key.
				fmt.Println(secretsmanager.ErrCodeDecryptionFailure, aerr.Error())

			case secretsmanager.ErrCodeInternalServiceError:
				// An error occurred on the server side.
				fmt.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())

			case secretsmanager.ErrCodeInvalidParameterException:
				// You provided an invalid value for a parameter.
				fmt.Println(secretsmanager.ErrCodeInvalidParameterException, aerr.Error())

			case secretsmanager.ErrCodeInvalidRequestException:
				// You provided a parameter value that is not valid for the current state of the resource.
				fmt.Println(secretsmanager.ErrCodeInvalidRequestException, aerr.Error())

			case secretsmanager.ErrCodeResourceNotFoundException:
				// We can't find the resource that you asked for.
				fmt.Println(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		CheckIfErrorWithMessage(err, "Error getting secret")
	}

	// Decrypts secret using the associated KMS key.
	// Depending on whether the secret is a string or binary, one of these fields will be populated.
	var secretString string
	//var decodedBinarySecret string
	if result.SecretString != nil {
		secretString = *result.SecretString
	} else {
		decodedBinarySecretBytes := make([]byte, base64.StdEncoding.DecodedLen(len(result.SecretBinary)))
		valueLen, err := base64.StdEncoding.Decode(decodedBinarySecretBytes, result.SecretBinary)

		CheckIfErrorWithMessage(err, "Base64 Decode Error")

		secretString = string(decodedBinarySecretBytes[:valueLen])
	}

	// https://irshadhasmat.medium.com/golang-simple-json-parsing-using-empty-interface-and-without-struct-in-go-language-e56d0e69968
	var conversion map[string]interface{}
	err = json.Unmarshal([]byte(secretString), &conversion)
	CheckIfError(err)

	return conversion //fmt.Sprintf("%v", conversion["Token"])
}
