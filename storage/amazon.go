package storage

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ectrc/snow/aid"
)

type AmazonClient struct {
	client *s3.Client
	ClientSettingsBucket	string
}

func NewAmazonClient(bucketURI, accessKeyID, secretAccessKey, clientSettingsBucket string) *AmazonClient {
	s := &AmazonClient{}

	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: bucketURI,
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")),
	)
	if err != nil {
		panic(err.Error())
	}

	s.client = s3.NewFromConfig(cfg)
	s.ClientSettingsBucket = clientSettingsBucket
	aid.Print("(amazon) client is ready")

	return s
}

func (a *AmazonClient) GetAllUserFiles() ([]string, error) {
	listObjectsOutput, err := a.client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(a.ClientSettingsBucket),
	})
	if err != nil {
		return nil, err
	}

	var files []string
	for _, object := range listObjectsOutput.Contents {
		fmt.Println(*object.Key)
		files = append(files, *object.Key)
	}

	return files, nil
}

func (a *AmazonClient) CreateUserFile(fileName string, data []byte) error {
	_, err := a.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(a.ClientSettingsBucket),
		Key:    aws.String("client/"+fileName),
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		return err
	}

	return nil
}

func (a *AmazonClient) GetUserFile(fileName string) ([]byte, error) {
	getObjectOutput, err := a.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(a.ClientSettingsBucket),
		Key:    aws.String("client/"+fileName),
	})
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(getObjectOutput.Body)

	bytes := buf.Bytes()
	decoded := make([]byte, base64.StdEncoding.DecodedLen(len(bytes)))
	base64.StdEncoding.Decode(decoded, bytes)

	return buf.Bytes(), nil
}