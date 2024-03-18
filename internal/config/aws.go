package config

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/viper"
	"log"
)

func NewAws(v *viper.Viper) *s3.Client {

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			v.GetString("aws.id"),
			v.GetString("aws.secret"),
			"",
		)),
		config.WithRegion("ap-southeast-1"),
	)
	if err != nil {
		log.Fatalf("Failed connect aws: %v", err)
		return nil
	}

	client := s3.NewFromConfig(cfg)
	return client
}
