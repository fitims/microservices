package infrastructure

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"thaThrowdown/common/database/dgraph"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const awsAccessKey = "AFSDRF56FSFGDFGCVBDV"                     //process.env.AWS_ACCESS_KEY;
const awsSecretKey = "dsfsd4rffgd946yfggfghfhfSADF345Sxxvbdfrw" //process.env.AWS_SECRET_KEY;
const s3Bucket = "mybucket"                                     //process.env.S3_BUCKET;

// MediaType is the media type that will be uploaded to S3
type MediaType string

// BuildKey build the amazon S3 bucket key
func BuildKey(mediaID dgraph.UID, mediaType MediaType, filename string) string {
	path := fmt.Sprintf("%s/%d__%s", mediaType, mediaID, filename)
	log.Println("Amazon S3 Path : ", path)
	return path
}

func getS3Session() (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(awsAccessKey, awsSecretKey, ""),
	})

	return sess, err
}

// GetS3SignedRequest returns a valid S3 Token to be used by clients to upload media directly to Amazon S3
func GetS3SignedRequest(s3key string) (string, error) {

	sess, err := getS3Session()
	if err != nil {
		log.Println("error gettings AWS Session : ", err)
		return "", err
	}

	svc := s3.New(sess)

	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(s3key),
		ACL:    aws.String("public-read"),
	})

	urlStr, err := req.Presign(15 * time.Minute)

	if err != nil {
		log.Println("Failed to sign request", err)
		return "", err
	}
	return urlStr, nil
}

// DeleteS3Object deletes an object from Amazon S3 storage
func DeleteS3Object(s3key string) error {
	sess, err := getS3Session()
	if err != nil {
		log.Println("error gettings AWS Session : ", err)
		return err
	}

	svc := s3.New(sess)

	params := &s3.DeleteObjectInput{
		Bucket:       aws.String(s3Bucket),
		Key:          aws.String(s3key),
		RequestPayer: aws.String("RequestPayer"),
	}

	_, err = svc.DeleteObject(params)
	if err != nil {
		log.Println("infrastructure.DeleteS3Object - Error : ", err)
		return err
	}
	return nil
}

// UploadS3Object uploads a file to Amazon S3 Storage
func UploadS3Object(s3key string, content []byte) error {
	sess, err := getS3Session()
	if err != nil {
		log.Println("error gettings AWS Session : ", err)
		return err
	}

	svc := s3.New(sess)

	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(s3key),
		Body:   bytes.NewReader(content),
		ACL:    aws.String("public-read"),
	})

	if err != nil {
		log.Println("infrastructure.UploadS3Object - Error : ", err)
		return err
	}
	return nil
}

// DownloadS3Object downloads a file from Amazon S3 storage
func DownloadS3Object(s3key string) ([]byte, error) {
	sess, err := getS3Session()
	if err != nil {
		log.Println("error gettings AWS Session : ", err)
		return nil, err
	}

	svc := s3.New(sess)

	obj, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(s3key),
	})

	if err != nil {
		log.Println("There was an error getting the object : ", s3key)
		return nil, err
	}

	img, err := ioutil.ReadAll(obj.Body)
	if err != nil {
		log.Println("There was an error loading the content of the object  : ", err)
		return nil, err
	}
	return img, nil
}
