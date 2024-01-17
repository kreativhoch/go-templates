package s3client

import (
	"bytes"
	ctx "context"
	"io"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.uber.org/zap"
)

type ClientInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

type S3Client struct {
	s3Client *s3.Client
	logger   *zap.SugaredLogger
	context  ctx.Context
}

func NewClient(logger *zap.SugaredLogger, httpClient httpclient.ClientInterface) *S3Client {
	return &S3Client{
		s3Client: s3.New(s3.Options{
			BaseEndpoint: aws.String(os.Getenv("S3_ENDPOINT")),
			Region:       os.Getenv("S3_REGION"),
			Credentials:  aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(os.Getenv("S3_KEY"), os.Getenv("S3_SECRET"), "")),
			HTTPClient:   httpClient,
		}),
		logger:  logger,
		context: ctx.Background(),
	}
}

func (s *S3Client) GetObject(key string) ([]byte, error) {
	object, err := s.s3Client.GetObject(s.context, &s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("S3_BUCKET")),
		Key:    aws.String(key),
	})

	if err != nil {
		return nil, err
	}

	responseBytes, readErr := io.ReadAll(object.Body)

	if readErr != nil {
		return nil, err
	}

	return responseBytes, err
}

func (s *S3Client) PutObject(key string, body []byte) error {
	_, err := s.s3Client.PutObject(s.context, &s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("S3_BUCKET")),
		Key:    aws.String(key),
		Body:   bytes.NewReader(body),
	})

	return err
}
