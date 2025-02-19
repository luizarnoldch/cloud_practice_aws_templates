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

//	Replicate the EndpointResolverV2 interface structure:
//	https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/dynamodb#EndpointResolverV2

//	type EndpointResolverV2 interface {
//		// ResolveEndpoint attempts to resolve the endpoint with the provided options,
//		// returning the endpoint if found. Otherwise an error is returned.
//		ResolveEndpoint(ctx context.Context, params EndpointParameters) (
//			smithyendpoints.Endpoint, error,
//		)
//	}

type CustomEndpointResolverV2 struct{}

func (*CustomEndpointResolverV2) ResolveEndpoint(ctx context.Context, params dynamodb.EndpointParameters) (smithy.Endpoint, error) {
	resolverV2, err := dynamodb.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
	if err != nil {
		log.Printf("Error resolving endpoint: %v\n", err)
		return smithy.Endpoint{}, err
	}

	if aws.ToString(params.Region) == AWS_REGION {
		parsedURL, err := url.Parse(resolverV2.URI.String())
		if err != nil {
			log.Printf("Error parsing resolved endpoint URI: %v\n", err)
			return smithy.Endpoint{}, err
		}

		parsedURL.Host = "localhost:8000"
		parsedURL.Scheme = "http"
		resolverV2.URI = *parsedURL
	}

	return resolverV2, nil
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
		o.EndpointResolverV2 = &CustomEndpointResolverV2{}
	})

	log.Printf("Local DynamoDB client configured successfully")
	return client, nil
}

func NewDynamoDBClient(ctx context.Context) (*dynamodb.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		defer ctx.Done()
		log.Printf("Error:%s", err.Error())
		return nil, err
	}
	log.Printf("DynamoDB connected successfully")
	return dynamodb.NewFromConfig(cfg), nil
}
