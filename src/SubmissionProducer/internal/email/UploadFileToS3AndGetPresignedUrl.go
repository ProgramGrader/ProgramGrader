package email

import (
	"SubmissionProducer/internal/common"
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"net/http"
	"os"
	"strings"
)

// S3PutObjectAPI defines the interface for the PutObject function.
// We use this interface to test the function using a mocked service.
type S3PutObjectAPI interface {
	PutObject(ctx context.Context,
		params *s3.PutObjectInput,
		optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

// PutFile uploads a file to an Amazon Simple Storage Service (Amazon S3) bucket
// Inputs:
//     c is the context of the method call, which includes the AWS Region
//     api is the interface that defines the method call
//     input defines the input arguments to the service call.
// Output:
//     If success, a PutObjectOutput object containing the result of the service call and nil
//     Otherwise, nil and an error from the call to PutObject
func PutFile(c context.Context, api S3PutObjectAPI, input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	return api.PutObject(c, input)
}

// S3ListObjectsAPI defines the interface for the ListObjectsV2 function.
// We use this interface to test the function using a mocked service.
type S3ListObjectsAPI interface {
	ListObjectsV2(ctx context.Context,
		params *s3.ListObjectsV2Input,
		optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
}

// GetObjects retrieves the objects in an Amazon Simple Storage Service (Amazon S3) bucket
// Inputs:
//     c is the context of the method call, which includes the AWS Region
//     api is the interface that defines the method call
//     input defines the input arguments to the service call.
// Output:
//     If success, a ListObjectsV2Output object containing the result of the service call and nil
//     Otherwise, nil and an error from the call to ListObjectsV2
func GetObjects(c context.Context, api S3ListObjectsAPI, input *s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error) {
	return api.ListObjectsV2(c, input)
}

// S3PresignGetObjectAPI defines the interface for the PresignGetObject function.
// We use this interface to test the function using a mocked service.
type S3PresignGetObjectAPI interface {
	PresignGetObject(
		ctx context.Context,
		params *s3.GetObjectInput,
		optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error)
}

// GetPresignedURL retrieves a presigned URL for an Amazon S3 bucket object.
// Inputs:
//     c is the context of the method call, which includes the AWS Region.
//     api is the interface that defines the method call.
//     input defines the input arguments to the service call.
// Output:
//     If successful, the presigned URL for the object and nil.
//     Otherwise, nil and an error from the call to PresignGetObject.
func GetPresignedURL(c context.Context, api S3PresignGetObjectAPI, input *s3.GetObjectInput) (*v4.PresignedHTTPRequest, error) {
	return api.PresignGetObject(c, input)
}

func UploadFileToS3AndGetPresignedUrl(pathAndFileName string, classNumber string, attachmentName string) string {

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := s3.NewFromConfig(cfg)

	bucket := common.BucketName

	prefix := fmt.Sprintf("emailAttachments/%s/", classNumber)

	ListInput := &s3.ListObjectsV2Input{
		Bucket: &bucket,
		Prefix: &prefix,
	}

	resp, err := GetObjects(context.TODO(), client, ListInput)
	if err != nil {
		fmt.Println("Got error retrieving list of objects:")
		fmt.Println(err)
	}

	var fileExists = false
	var numberOfFileExist int

	// List all items in emailAttachments/Class#/
	for _, item := range resp.Contents {
		s3FileNameWithExtention := strings.Split(*item.Key, "/")[len(strings.Split(*item.Key, "/"))-1]

		nameWithNumber := strings.Split(s3FileNameWithExtention, ".")[0]

		name := strings.Join(strings.Split(nameWithNumber, "-")[:len(strings.Split(nameWithNumber, "-"))-1], "")

		if name == attachmentName {
			fileExists = true
			numberOfFileExist = numberOfFileExist + 1
		}
	}

	// if not exists <name>-1 else increament by how many exists
	var filename string
	if fileExists {
		filename = fmt.Sprintf("emailAttachments/%s/%v-%v.zip", classNumber, attachmentName, numberOfFileExist)
	} else {
		filename = fmt.Sprintf("emailAttachments/%s/%v-%v.zip", classNumber, attachmentName, 1)
	}

	// add file to S3

	file, err := os.Open(pathAndFileName)
	common.CheckIfErrorWithMessage(err, "Unable to open file to put in S3")

	defer func(file *os.File) {
		err := file.Close()
		common.CheckIfError(err)
	}(file)

	//PutInput := &s3.PutObjectInput{
	//	Bucket: &bucket,
	//	Key:    &filename,
	//	Body:   file,
	//}
	//
	//_, err = PutFile(context.TODO(), client, PutInput)

	upFileInfo, _ := file.Stat()
	var fileSize int64 = upFileInfo.Size()
	fileBuffer := make([]byte, fileSize)
	file.Read(fileBuffer)

	inputPut := s3.PutObjectInput{
		Bucket:        aws.String(bucket),
		Key:           aws.String(filename),
		ACL:           types.ObjectCannedACL("private"),
		Body:          bytes.NewReader(fileBuffer),
		ContentLength: fileSize,
		ContentType:   aws.String(http.DetectContentType(fileBuffer)),
	}

	// Put the file object to s3 with the file name
	_, err = client.PutObject(context.Background(), &inputPut)
	common.CheckIfErrorWithMessage(err, "error uploading file")

	// return presigned url
	input := &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &filename,
	}

	psClient := s3.NewPresignClient(client)

	urlResp, err := GetPresignedURL(context.TODO(), psClient, input)
	common.CheckIfErrorWithMessage(err, "error retrieving pre-signed object")

	return urlResp.URL

}
