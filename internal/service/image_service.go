package service

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

type ImageService struct {
	AwsClient *s3.Client
	Validate  *validator.Validate
	Log       *logrus.Logger
}

func NewImageService(s *s3.Client, validate *validator.Validate, log *logrus.Logger) *ImageService {
	return &ImageService{AwsClient: s, Validate: validate, Log: log}
}

func (s *ImageService) Upload(ctx context.Context, file io.Reader, filename string) (string, error) {
	input := &s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("S3_BUCKET_NAME")),
		Key:    aws.String(filename),
		ACL:    types.ObjectCannedACLPublicRead,
		Body:   file,
	}

	_, err := s.AwsClient.PutObject(ctx, input)
	if err != nil {
		s.Log.Error("Unable to upload image to S3: ", err)
		return "", err
	}

	url := "https://" + *input.Bucket + ".s3.amazonaws.com/" + *input.Key

	return url, nil
}
