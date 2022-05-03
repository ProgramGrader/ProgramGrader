package tests

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"log"
)

// CreateLocalClient returns the config associated with local docker container
func CreateLocalClient() (*dynamodb.Client, error) {
	awsEndpoint := "http://localhost:8000"
	awsRegion := "us-east-2"

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(awsRegion),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: awsEndpoint}, nil
			})),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: "dummy", SecretAccessKey: "dummy", SessionToken: "dummy",
				Source: "Hard-coded credentials; values are irrelevant for local DynamoDB",
			},
		}),
	)
	if err != nil {
		panic(err)
	}

	return dynamodb.NewFromConfig(cfg), err
}

var localClientConfig, _ = CreateLocalClient()

func TableExists(client *dynamodb.Client, name string) bool {
	tables, err := client.ListTables(context.TODO(), &dynamodb.ListTablesInput{})
	if err != nil {
		log.Fatal("error listing tables", err)

	}
	for _, n := range tables.TableNames {
		if n == name {
			return true
		}
	}
	return false
}

// CreateTable creates tables if it does not exist
func CreateTable(client *dynamodb.Client, tableName string) {
	if !TableExists(client, tableName) {
		var input = dynamodb.CreateTableInput{
			AttributeDefinitions: []types.AttributeDefinition{
				{
					AttributeName: aws.String("key"),
					AttributeType: types.ScalarAttributeTypeS,
				},
			},
			KeySchema: []types.KeySchemaElement{
				{
					AttributeName: aws.String("key"),
					KeyType:       types.KeyTypeHash,
				},
			},

			BillingMode: types.BillingModePayPerRequest,
			TableName:   aws.String(tableName),
		}
		_, err := client.CreateTable(context.TODO(), &input)
		if err != nil {
			log.Fatal("Error creating table ", err)
		}
	} else {
		print("Table exists\n")
	}
}

// Get given hash returns value
func Get(tableName string, hash string) string {

	getItemInput := &dynamodb.GetItemInput{
		TableName:            aws.String(tableName),
		ConsistentRead:       aws.Bool(true),
		ProjectionExpression: aws.String("s3object"),

		Key: map[string]types.AttributeValue{
			"key": &types.AttributeValueMemberS{Value: hash},
		},
	}

	output, err := localClientConfig.GetItem(context.TODO(), getItemInput)
	if err != nil {
		log.Fatalf("Failed to get item, %v", err)
	}

	if output.Item == nil {
		log.Fatal("Item not found: ", hash)
	}

	var value string

	err = attributevalue.Unmarshal(output.Item["s3object"], &value)
	if err != nil {
		log.Fatalf("unmarshal failed, %v", err)
	}

	return value

}

// Put creates/update a new entry in the Dynamodb
func Put(tableName string, hash string, s3object string) {

	var itemInput = dynamodb.PutItemInput{
		TableName: aws.String(tableName),

		Item: map[string]types.AttributeValue{
			"key":      &types.AttributeValueMemberS{Value: hash},
			"s3object": &types.AttributeValueMemberS{Value: s3object},
		},
	}

	_, err := localClientConfig.PutItem(context.TODO(), &itemInput)
	if err != nil {
		log.Fatal("Error inserting value ", err)
	}
}

// Delete removes a item from the table given the key
func Delete(tableName string, key string) error {

	deleteInput := dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"key": &types.AttributeValueMemberS{Value: key},
		},
	}

	_, err := localClientConfig.DeleteItem(context.TODO(), &deleteInput)
	if err != nil {
		panic(err)
	}

	return err
}

// DeleteAll for testing purposes
func DeleteAll(tableName string) {
	scan := dynamodb.NewScanPaginator(localClientConfig, &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	})

	for scan.HasMorePages() {
		out, err := scan.NextPage(context.TODO())
		if err != nil {
			print("Page error")
			panic(err)
		}

		for _, item := range out.Items {
			_, err = localClientConfig.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
				TableName: aws.String(tableName),
				Key: map[string]types.AttributeValue{
					"key": item["key"],
				},
			})
			if err != nil {
				print("Error Deleting Item")
				panic(err)
			}

		}
	}
}
