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

// DynamoDBClient abstracts & encapsulates interactions with DynamoDB service
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

// Opts is the type signature of functions modifying a DynamoDBClient struct
type Opts func(d *DynamoDBClient) *DynamoDBClient

// New creates and configures a new DynamoDBClient with the specified options
func New(opt ...Opts) *DynamoDBClient {
	d := &DynamoDBClient{}
	for _, fn := range opt {
		d = fn(d)
	}
	return d
}

// WithAWSConfig inserts the given region in the aws config
func WithAWSConfig(region string) Opts {
	return func(d *DynamoDBClient) *DynamoDBClient {
		d.config = aws.NewConfig().WithRegion(region)
		return d
	}
}

// WithAWSSession inserts the given AWS session object into DynamoDBClient
func WithAWSSession(sess *types.AWSSession) Opts {
	return func(d *DynamoDBClient) *DynamoDBClient {
		d.session = sess
		return d
	}
}

// WithSTSCreds configures the AWS config with STS creds using
// the injected role ARN, if running in production
func WithSTSCreds(isProd bool) Opts {
	return func(d *DynamoDBClient) *DynamoDBClient {
		if isProd {
			creds := stscreds.NewCredentials(d.session, DBRoleARN)
			d.config = d.config.WithCredentials(creds)
		}
		return d
	}
}

// WithConfigEndpoint configures the DynamoDB endpoint, if not running in prod
func WithConfigEndpoint(endpoint string, isProd bool) Opts {
	return func(d *DynamoDBClient) *DynamoDBClient {
		if !isProd {
			d.config = d.config.WithEndpoint(endpoint)
		}
		return d
	}
}

// NewClient creates a new DynamoDB client using the provided session and config
func NewClient() Opts {
	return func(d *DynamoDBClient) *DynamoDBClient {
		d.client = dynamodb.New(d.session, d.config)
		return d
	}
}

// WithClient inserts a dynamodb client into DynamoDBClient, useful in testing
func WithClient(m *types.AWSDynamoDB) Opts {
	return func(d *DynamoDBClient) *DynamoDBClient {
		d.client = m
		return d
	}
}

// WithPaginationLimit inserts the given pagination limit into DynamoDBClient
func WithPaginationLimit(limit int64) Opts {
	return func(d *DynamoDBClient) *DynamoDBClient {
		d.limit = limit
		return d
	}
}

// WithQuoteTable inserts the name of the QuoteTable into DynamoDBClient
func WithQuoteTable(t string) Opts {
	return func(d *DynamoDBClient) *DynamoDBClient {
		d.QuoteTable = t
		return d
	}
}

// WithTopicTable inserts the name of the TopicTable into DynamoDBClient
func WithTopicTable(t string) Opts {
	return func(d *DynamoDBClient) *DynamoDBClient {
		d.TopicTable = t
		return d
	}
}

// WithAuthorTable inserts the name of the AuthorTable into DynamoDBClient
func WithAuthorTable(t string) Opts {
	return func(d *DynamoDBClient) *DynamoDBClient {
		d.AuthorTable = t
		return d
	}
}
