package db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/travmatth-org/qaas/internal/types"
)

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

func New() *DynamoDBClient {
	return &DynamoDBClient{}
}

func (d *DynamoDBClient) WithAWSConfig(region string) *DynamoDBClient {
	d.config = aws.NewConfig().WithRegion(region)
	return d
}

func (d *DynamoDBClient) WithAWSSession(sess *types.AWSSession) *DynamoDBClient {
	d.session = sess
	return d
}

func (d *DynamoDBClient) WithSTSCreds(arn string) *DynamoDBClient {
	creds := stscreds.NewCredentials(d.session, arn)
	d.config = d.config.WithCredentials(creds)
	return d
}

func (d *DynamoDBClient) WithConfigEndpoint(endpoint string) *DynamoDBClient {
	d.config = d.config.WithEndpoint(endpoint)
	return d
}

func (d *DynamoDBClient) NewClient() *DynamoDBClient {
	d.client = dynamodb.New(d.session, d.config)
	return d
}

func (d *DynamoDBClient) Mock(m *types.AWSDynamoDB) *DynamoDBClient {
	d.client = m
	return d
}

func (d *DynamoDBClient) WithPaginationLimit(limit int64) *DynamoDBClient {
	d.limit = limit
	return d
}

func (d *DynamoDBClient) WithQuoteTable(t string) *DynamoDBClient {
	d.QuoteTable = t
	return d
}

func (d *DynamoDBClient) WithTopicTable(t string) *DynamoDBClient {
	d.TopicTable = t
	return d
}

func (d *DynamoDBClient) WithAuthorTable(t string) *DynamoDBClient {
	d.AuthorTable = t
	return d
}
