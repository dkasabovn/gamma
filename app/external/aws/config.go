package aws

import (
	"context"

	"gamma/app/system"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

var localStackResolver = aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
	return aws.Endpoint{
		PartitionID:   "aws",
		URL:           "http://localhost:4566",
		SigningRegion: "us-west-1",
	}, nil
})

func NewConfig() aws.Config {
	var (
		cfg aws.Config
		err error
	)

	// TODO: AWS Region
	if system.GetConfig().Environment != "prod" {
		cfg, err = config.LoadDefaultConfig(context.Background(), config.WithRegion(""), config.WithEndpointResolverWithOptions(localStackResolver))
	} else {
		cfg, err = config.LoadDefaultConfig(context.Background(), config.WithRegion(""))
	}

	if err != nil {
		panic(err)
	}
	return cfg
}
