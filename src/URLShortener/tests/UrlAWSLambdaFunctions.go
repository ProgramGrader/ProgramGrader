package tests

import (
	"fmt"
	"log"
	"net/url"
	"strings"
)

// Unit Testing functions that will be used in aws lambda
// get s3url from bucket and verify that it exists.

// seems like we're going to need to make a Pre-Signed Url then store it into our dynamodb, and use that instead of the
// s3 object url, but this is something we need to fix in the put function for the local dynamodb

// create a temporary url --> s3url takes you to the obj which is a file with the url inside it, you need to convert
// the s3 object to a file then extract its contents. Then create the Url

// redirect user to said temp url

// getBucketAnKey given a s3 object url returns the bucket string and key. below method won't work because of authentication
func getBucketAnKey(s3ObjectUrl string) (bucket string, key string) {
	u, err := url.Parse(s3ObjectUrl)
	if err != nil {
		log.Fatal(err)
	}

	path := strings.SplitN(u.Path, "/", 3)
	bucket = path[1]
	key = path[2]
	fmt.Println(bucket)
	fmt.Println(key)

	return bucket, key
}

// Pre-signed Url needs to be created then stored into the local dynamodb because the url is temporarily available
// meaning if one is created manually testing will eventually not work because the URL will expire

// Creates and returns a Pre-signed URL given the
//func createPresignedURL(s3url string) {
//	// before getting a Pre-signed URL you must first create a Pre-signed client
//	fmt.Println("Create a pre-signed Client")
//	presignClient := s3.NewPre
//}

// looks like we have to use a presigned url to bypass access is denied :(
func downloadS3object(key string) {
	// The session the S3 Downloader will use
	//sess := session.Must(session.NewSession())
	//bucket, key := getBucketAnKey(get(key))
	//downloader := s3manager.NewDownloader(sess)
	//
	//downloadFile, error := os.Create("s3Url.txt")
	//if error != nil {
	//	log.Fatal("Failed to create new file ", error)
	//}
	//
	//_, error = downloader.Download(downloadFile, &s3.GetObjectInput{
	//	Bucket: aws.String(bucket),
	//	Key:    aws.String(key),
	//})
	//if error != nil {
	//	log.Fatal("Failed to download s3 object", error)
	//}
	//
	//err := downloadFile.Close()
	//if err != nil {
	//	return
	//}
}
