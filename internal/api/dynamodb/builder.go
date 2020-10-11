package db

import (
	"github.com/travmatth-org/qaas/internal/types"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

// DBRoleARN injected during build process
var DBRoleARN string

type DynamoDBClient struct {
	session     *types.AWSSession
	config      *types.AWSConfig
	client      *types.AWSDynamoDB
	limit       int64
	expr        *expression.Expression
	QuoteTable  string
	TopicTable  string
	AuthorTable string
}

type dynamodbOpts func(d *DynamoDBClient) *DynamoDBClient

func New(opts ...dynamodbOpts) *DynamoDBClient {
	d := &DynamoDBClient{}
	for _, opt := range opts {
		d = opt(d)
	}
	return d
}

func WithAWSConfig(region string) dynamodbOpts {
	return func(d *DynamoDBClient) *DynamoDBClient {
		d.config = aws.NewConfig().WithRegion(region)
		return d
	}
}

func WithAWSSession(sess *types.AWSSession) dynamodbOpts {
	return func(d *DynamoDBClient) *DynamoDBClient {
		d.session = sess
		return d
	}
}

func WithSTSCreds(isProd bool) dynamodbOpts {
	return func(d *DynamoDBClient) *DynamoDBClient {
		if isProd {
			creds := stscreds.NewCredentials(d.session, DBRoleARN)
			d.config = d.config.WithCredentials(creds)
		}
		return d
	}
}

func WithConfigEndpoint(endpoint string, isProd bool) dynamodbOpts {
	return func(d *DynamoDBClient) *DynamoDBClient {
		if !isProd {
			d.config = d.config.WithEndpoint(endpoint)
		}
		return d
	}
}

func NewClient() dynamodbOpts {
	return func(d *DynamoDBClient) *DynamoDBClient {
		d.client = dynamodb.New(d.session, d.config)
		return d
	}
}

func MockClient(m *types.AWSDynamoDB) dynamodbOpts {
	return func(d *DynamoDBClient) *DynamoDBClient {
		d.client = m
		return d
	}
}

func WithPaginationLimit(limit int64) dynamodbOpts {
	return func(d *DynamoDBClient) *DynamoDBClient {
		d.limit = limit
		return d
	}
}

func WithQuoteTable(t string) dynamodbOpts {
	return func(d *DynamoDBClient) *DynamoDBClient {
		d.QuoteTable = t
		return d
	}
}

func WithTopicTable(t string) dynamodbOpts {
	return func(d *DynamoDBClient) *DynamoDBClient {
		d.TopicTable = t
		return d
	}
}

func WithAuthorTable(t string) dynamodbOpts {
	return func(d *DynamoDBClient) *DynamoDBClient {
		d.AuthorTable = t
		return d
	}
}
