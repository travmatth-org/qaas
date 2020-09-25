package db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/travmatth-org/qaas/internal/types"
)

type DynamoDBClient struct {
	client      *types.AWSDynamoDB
	QuoteTable  string
	TopicTable  string
	AuthorTable string
	limit       int
	expr        *expression.Expression
}

func New() *DynamoDBClient {
	return &DynamoDBClient{}
}

func (d *DynamoDBClient) Init(sess *session.Session, config *aws.Config) *DynamoDBClient {
	d.client = dynamodb.New(sess, config)
	return d
}

func (d *DynamoDBClient) Mock(m *types.AWSDynamoDB) *DynamoDBClient {
	d.client = m
	return d
}

func (d *DynamoDBClient) SetPaginationLimit(limit int) *DynamoDBClient {
	d.limit = limit
	return d
}

func (d *DynamoDBClient) SetQuoteTable(t string) *DynamoDBClient {
	d.QuoteTable = t
	return d
}

func (d *DynamoDBClient) SetTopicTable(t string) *DynamoDBClient {
	d.TopicTable = t
	return d
}

func (d *DynamoDBClient) SetAuthorTable(t string) *DynamoDBClient {
	d.AuthorTable = t
	return d
}
