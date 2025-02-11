package lib

import (
	"context"
	"log"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	smithy "github.com/aws/smithy-go/endpoints"
)

var (
	AWS_REGION        = "us-east-1"
	AWS_CLIENT_DOMAIN = "http://localhost:8000"
)

type resolverV2 struct{}

func (*resolverV2) ResolveEndpoint(ctx context.Context, params dynamodb.EndpointParameters) (smithy.Endpoint, error) {
	if aws.ToString(params.Region) == AWS_REGION {
		u, err := url.Parse(AWS_CLIENT_DOMAIN)
		if err != nil {
			return smithy.Endpoint{}, nil
		}
		return smithy.Endpoint{
			URI: *u,
		}, nil
	}

	return dynamodb.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
}

func NewLocalDynamoDBClient(ctx context.Context) (*dynamodb.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(AWS_REGION),
		config.WithRetryer(func() aws.Retryer {
			return retry.AddWithMaxAttempts(retry.NewStandard(), 5)
		}),
	)
	if err != nil {
		return nil, err
	}

	client := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.EndpointResolverV2 = &resolverV2{}
	})

	log.Printf("Local DynamoDB client configured successfully")
	return client, nil
}
